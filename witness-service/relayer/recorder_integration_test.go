package relayer

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
)

func mysqlIntegrationStore(t *testing.T) *db.SQLStore {
	t.Helper()
	dsn := os.Getenv("IOTUBE_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("IOTUBE_TEST_MYSQL_DSN is not set")
	}
	store := db.NewSQLStoreFactory().NewStore(db.Config{Driver: "mysql", URI: dsn})
	require.NoError(t, store.Start(context.Background()))
	t.Cleanup(func() {
		require.NoError(t, store.Stop(context.Background()))
	})
	return store
}

func TestMigrateTransferPrimaryKeyMySQL(t *testing.T) {
	store := mysqlIntegrationStore(t)
	table := fmt.Sprintf("relayer_legacy_%d", time.Now().UnixNano())
	witnessTable := table + "_witnesses"
	quotedTable, err := mysqlIdentifier(table)
	require.NoError(t, err)
	quotedWitnessTable, err := mysqlIdentifier(witnessTable)
	require.NoError(t, err)
	_, err = store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE %s (`cashier` varchar(64) NOT NULL, `token` varchar(42) NOT NULL, `tidx` bigint NOT NULL, `id` varchar(66) NOT NULL, PRIMARY KEY (`cashier`,`token`,`tidx`), UNIQUE KEY `id_UNIQUE` (`id`)) ENGINE=InnoDB",
		quotedTable,
	))
	require.NoError(t, err)
	_, err = store.DB().Exec(fmt.Sprintf(
		"CREATE TABLE %s (`transferId` varchar(66) NOT NULL, PRIMARY KEY (`transferId`), FOREIGN KEY (`transferId`) REFERENCES %s (`id`) ON DELETE CASCADE) ENGINE=InnoDB",
		quotedWitnessTable,
		quotedTable,
	))
	require.NoError(t, err)
	t.Cleanup(func() {
		_, dropErr := store.DB().Exec("DROP TABLE IF EXISTS " + quotedWitnessTable)
		require.NoError(t, dropErr)
		_, dropErr = store.DB().Exec("DROP TABLE IF EXISTS " + quotedTable)
		require.NoError(t, dropErr)
	})

	require.NoError(t, migrateTransferPrimaryKey(store, table))
	require.NoError(t, migrateTransferPrimaryKey(store, table))
	_, err = store.DB().Exec(
		fmt.Sprintf("INSERT INTO %s (`cashier`,`token`,`tidx`,`id`) VALUES (?,?,?,?),(?,?,?,?)", quotedTable),
		"cashier", "token", 1, "id-1",
		"cashier", "token", 1, "id-2",
	)
	require.NoError(t, err)

	var count int
	require.NoError(t, store.DB().QueryRow("SELECT COUNT(*) FROM "+quotedTable).Scan(&count))
	require.Equal(t, 2, count)
}

func TestTransfersWithWitnessQuorumMySQL(t *testing.T) {
	store := mysqlIntegrationStore(t)
	suffix := time.Now().UnixNano()
	transferTable := fmt.Sprintf("relayer_transfers_%d", suffix)
	witnessTable := fmt.Sprintf("relayer_witnesses_%d", suffix)
	recorder := NewRecorder(store, nil, transferTable, witnessTable, "", "")
	require.NoError(t, recorder.Start(context.Background()))
	t.Cleanup(func() {
		quotedWitness, quoteErr := mysqlIdentifier(witnessTable)
		require.NoError(t, quoteErr)
		quotedTransfer, quoteErr := mysqlIdentifier(transferTable)
		require.NoError(t, quoteErr)
		_, dropErr := store.DB().Exec("DROP TABLE IF EXISTS " + quotedWitness)
		require.NoError(t, dropErr)
		_, dropErr = store.DB().Exec("DROP TABLE IF EXISTS " + quotedTransfer)
		require.NoError(t, dropErr)
	})

	validatorAddr := common.HexToAddress("0x9000000000000000000000000000000000000009")
	cashierAddr := common.HexToAddress("0x1000000000000000000000000000000000000001")
	proposal := func(recipient common.Address) *Transfer {
		transfer := &Transfer{
			cashier:     util.ETHAddressToAddress(cashierAddr),
			token:       util.ETHAddressToAddress(common.HexToAddress("0x2000000000000000000000000000000000000002")),
			index:       7,
			sender:      util.ETHAddressToAddress(common.HexToAddress("0x3000000000000000000000000000000000000003")),
			recipient:   util.ETHAddressToAddress(recipient),
			amount:      big.NewInt(10),
			fee:         big.NewInt(0),
			timestamp:   time.Now(),
			gasPrice:    big.NewInt(0),
			txSender:    util.ETHAddressToAddress(common.HexToAddress("0x5000000000000000000000000000000000000005")),
			blockHeight: 100,
		}
		transfer.GenID(validatorAddr)
		return transfer
	}
	blockedProposal := proposal(common.HexToAddress("0x4000000000000000000000000000000000000004"))
	readyProposal := proposal(common.HexToAddress("0x4100000000000000000000000000000000000004"))
	activeWitnesses := []common.Address{
		common.HexToAddress("0xa000000000000000000000000000000000000001"),
		common.HexToAddress("0xa000000000000000000000000000000000000002"),
		common.HexToAddress("0xa000000000000000000000000000000000000003"),
		common.HexToAddress("0xa000000000000000000000000000000000000004"),
	}
	require.NoError(t, recorder.AddTransferAndWitness(
		util.ETHAddressToAddress(validatorAddr),
		blockedProposal,
		&Witness{addr: activeWitnesses[0].Bytes(), signature: []byte{1}},
		nil,
		nil,
	))
	for i := 0; i < 3; i++ {
		require.NoError(t, recorder.AddTransferAndWitness(
			util.ETHAddressToAddress(validatorAddr),
			readyProposal,
			&Witness{addr: activeWitnesses[i].Bytes(), signature: []byte{byte(i + 1)}},
			nil,
			nil,
		))
	}

	transfers, err := recorder.TransfersWithWitnessQuorum(
		10,
		activeWitnesses,
		CashiersQueryOption([]string{util.ETHAddressToAddress(cashierAddr).String()}),
	)
	require.NoError(t, err)
	require.Len(t, transfers, 1)
	require.Equal(t, readyProposal.ID(), transfers[0].ID())
}
