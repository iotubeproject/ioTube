package relayer

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/ioTube/witness-service/util"
)

func TestHeartbeatVersionColumnsMySQL(t *testing.T) {
	store := mysqlIntegrationStore(t)
	exec := func(q string, args ...any) {
		t.Helper()
		_, err := store.DB().Exec(q, args...)
		require.NoError(t, err)
	}
	// witness_heartbeats as created before the build columns existed
	legacy := func() {
		exec("DROP TABLE IF EXISTS witness_heartbeats")
		exec("CREATE TABLE witness_heartbeats (" +
			"`witness` varchar(66) NOT NULL, `cashier` varchar(64) NOT NULL, `tipHeight` bigint(20) NOT NULL DEFAULT 0," +
			"`creationTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP," +
			"`updateTime` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP," +
			"PRIMARY KEY (`witness`, `cashier`), KEY `updateTime_index` (`updateTime`)) ENGINE=InnoDB DEFAULT CHARSET=latin1")
	}
	t.Cleanup(func() { exec("DROP TABLE IF EXISTS witness_heartbeats") })
	read := func(witness string) (tip uint64, ver, commit string) {
		t.Helper()
		require.NoError(t, store.DB().QueryRow(
			"SELECT `tipHeight`, `version`, `commitID` FROM witness_heartbeats WHERE `witness`=?", witness,
		).Scan(&tip, &ver, &commit))
		return
	}

	legacy()
	exec("INSERT INTO witness_heartbeats (`witness`, `cashier`, `tipHeight`) VALUES ('0xold', 'c', 1)")
	for i := 0; i < 2; i++ {
		ok, err := migrateHeartbeatVersionColumns(store)
		require.NoError(t, err)
		require.True(t, ok)
	}
	_, ver, commit := read("0xold")
	require.Equal(t, "", ver)
	require.Equal(t, "", commit)

	addr := common.HexToAddress("0x5160").Bytes()
	witness := "0x" + common.Bytes2Hex(addr)
	cashier := util.ETHAddressToAddress(common.HexToAddress("0xc1"))
	require.NoError(t, upsertWitnessHeartbeat(store, true, addr, cashier, 7, "v7.1.0", "abcdef0"))
	tip, ver, commit := read(witness)
	require.Equal(t, uint64(7), tip)
	require.Equal(t, "v7.1.0", ver)
	require.Equal(t, "abcdef0", commit)

	require.NoError(t, upsertWitnessHeartbeat(store, true, addr, cashier, 8, "v7.1.0'; --", "abcdef0"))
	tip, ver, _ = read(witness)
	require.Equal(t, uint64(8), tip)
	require.Equal(t, "", ver)

	// a relayer that predates the columns keeps writing heartbeats
	require.NoError(t, upsertWitnessHeartbeat(store, false, addr, cashier, 9, "ignored", "ignored"))
	tip, _, _ = read(witness)
	require.Equal(t, uint64(9), tip)

	// relayers sharing the DB may migrate at the same time
	legacy()
	var wg sync.WaitGroup
	errs := make([]error, 4)
	for i := range errs {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_, errs[i] = migrateHeartbeatVersionColumns(store)
		}(i)
	}
	wg.Wait()
	for _, err := range errs {
		require.NoError(t, err)
	}
	require.NoError(t, upsertWitnessHeartbeat(store, true, addr, cashier, 10, "v7.1.0", "abcdef0"))

	// a long transaction holding the table's metadata lock makes the
	// migration give up quickly; heartbeats keep working without versions
	legacy()
	tx, err := store.DB().Begin()
	require.NoError(t, err)
	var rows int
	require.NoError(t, tx.QueryRow("SELECT COUNT(*) FROM witness_heartbeats").Scan(&rows))
	oldWait := heartbeatMigrationLockWait
	heartbeatMigrationLockWait = 1
	start := time.Now()
	require.False(t, heartbeatVersionsSupported(store))
	require.Less(t, time.Since(start), 10*time.Second)
	heartbeatMigrationLockWait = oldWait
	require.NoError(t, upsertWitnessHeartbeat(store, false, addr, cashier, 12, "", ""))
	require.NoError(t, tx.Commit())
	var lockWait int
	require.NoError(t, store.DB().QueryRow("SELECT @@SESSION.lock_wait_timeout").Scan(&lockWait))
	require.NotEqual(t, 1, lockWait, "a pooled connection kept the migration's lock_wait_timeout")

	// the Solana recorder (no witness set) caps new rows but keeps updating
	// existing ones
	legacy()
	sol := NewSolRecorder(store, "", "", util.NewETHAddressDecoder(), util.NewSOLAddressDecoder())
	sol.heartbeatVersions, _ = migrateHeartbeatVersionColumns(store)
	oldMax := maxUnverifiedHeartbeatRows
	maxUnverifiedHeartbeatRows = 3
	t.Cleanup(func() { maxUnverifiedHeartbeatRows = oldMax })
	solCashier := util.ETHAddressToAddress(common.HexToAddress("0xc5"))
	for i := 1; i <= 5; i++ {
		require.NoError(t, sol.UpsertWitnessHeartbeat(common.BigToAddress(big.NewInt(int64(i))).Bytes(), solCashier, uint64(i), "v7.1.0", "abc"))
	}
	var total int
	require.NoError(t, store.DB().QueryRow("SELECT COUNT(*) FROM witness_heartbeats").Scan(&total))
	require.Equal(t, 3, total)
	// concurrent heartbeats of new witnesses at the cap boundary
	exec("DELETE FROM witness_heartbeats WHERE `tipHeight` = 3")
	var cwg sync.WaitGroup
	for i := 10; i < 30; i++ {
		cwg.Add(1)
		go func(i int) {
			defer cwg.Done()
			require.NoError(t, sol.UpsertWitnessHeartbeat(common.BigToAddress(big.NewInt(int64(i))).Bytes(), solCashier, uint64(i), "", ""))
		}(i)
	}
	cwg.Wait()
	require.NoError(t, store.DB().QueryRow("SELECT COUNT(*) FROM witness_heartbeats").Scan(&total))
	require.Equal(t, 3, total)
	require.NoError(t, sol.UpsertWitnessHeartbeat(common.BigToAddress(big.NewInt(1)).Bytes(), solCashier, 99, "v7.1.1", "def"))
	tip, ver, _ = read("0x" + common.Bytes2Hex(common.BigToAddress(big.NewInt(1)).Bytes()))
	require.Equal(t, uint64(99), tip)
	require.Equal(t, "v7.1.1", ver)
	maxUnverifiedHeartbeatRows = oldMax

	// a starting recorder migrates and records builds
	legacy()
	suffix := time.Now().UnixNano()
	recorder := NewRecorder(store, nil, fmt.Sprintf("relayer_transfers_%d", suffix), fmt.Sprintf("relayer_witnesses_%d", suffix), "", "")
	require.NoError(t, recorder.Start(context.Background()))
	t.Cleanup(func() {
		exec(fmt.Sprintf("DROP TABLE IF EXISTS relayer_witnesses_%d", suffix))
		exec(fmt.Sprintf("DROP TABLE IF EXISTS relayer_transfers_%d", suffix))
	})
	require.True(t, recorder.heartbeatVersions)
	require.NoError(t, recorder.UpsertWitnessHeartbeat(addr, cashier, 11, "v7.1.0-1-g1234567", "1234567"))
	tip, ver, commit = read(witness)
	require.Equal(t, uint64(11), tip)
	require.Equal(t, "v7.1.0-1-g1234567", ver)
	require.Equal(t, "1234567", commit)
}
