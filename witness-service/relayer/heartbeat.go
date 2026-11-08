// Copyright (c) 2026 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package relayer

import (
	"context"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"

	"github.com/iotexproject/ioTube/witness-service/db"
	"github.com/iotexproject/ioTube/witness-service/util"
)

// maxHeartbeatVersionLen bounds the stored build strings.
const maxHeartbeatVersionLen = 64

// heartbeatVersionColumns are the witness build columns of witness_heartbeats,
// added to tables created before they existed.
var heartbeatVersionColumns = []string{"version", "commitID"}

// heartbeatMigrationLockWait (seconds) bounds how long adding the columns waits
// for the table's metadata lock. While an ALTER waits for it, every later query
// on the table queues behind it, including heartbeats of relayers sharing the
// DB; better to give up and run without versions.
var heartbeatMigrationLockWait = 5

// migrateHeartbeatVersionColumns adds the build columns to witness_heartbeats
// and reports whether they are present. Relayers sharing the DB may race to add
// them; a duplicate column is fine.
func migrateHeartbeatVersionColumns(store *db.SQLStore) (bool, error) {
	if store.DriverName() != "mysql" {
		return false, nil
	}
	schema, table, quoted, err := mysqlTable(store, witnessHeartbeatTableName)
	if err != nil {
		return false, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(4*heartbeatMigrationLockWait)*time.Second)
	defer cancel()
	conn, err := store.DB().Conn(ctx)
	if err != nil {
		return false, errors.Wrap(err, "failed to get a connection")
	}
	defer conn.Close()
	if _, err := conn.ExecContext(ctx, fmt.Sprintf("SET SESSION lock_wait_timeout = %d", heartbeatMigrationLockWait)); err != nil {
		return false, errors.Wrap(err, "failed to set lock_wait_timeout")
	}
	defer func() {
		// the connection goes back to the pool; drop it rather than leave the
		// short lock timeout on it
		_ = conn.Raw(func(any) error { return driver.ErrBadConn })
	}()
	for _, col := range heartbeatVersionColumns {
		var n int
		if err := conn.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA=? AND TABLE_NAME=? AND COLUMN_NAME=?",
			schema, table, col,
		).Scan(&n); err != nil {
			return false, errors.Wrapf(err, "failed to inspect column %s", col)
		}
		if n > 0 {
			continue
		}
		_, err := conn.ExecContext(ctx, fmt.Sprintf(
			"ALTER TABLE %s ADD COLUMN `%s` varchar(%d) NOT NULL DEFAULT ''", quoted, col, maxHeartbeatVersionLen,
		))
		var me *mysql.MySQLError
		if err != nil && !(errors.As(err, &me) && me.Number == 1060) { // 1060: duplicate column name
			return false, errors.Wrapf(err, "failed to add column %s", col)
		}
	}
	return true, nil
}

// heartbeatVersionsSupported runs the migration at recorder start. Heartbeats
// are best-effort monitoring, so a failed migration (e.g. no ALTER privilege)
// only stops recording builds; liveness keeps being recorded.
func heartbeatVersionsSupported(store *db.SQLStore) bool {
	ok, err := migrateHeartbeatVersionColumns(store)
	if err != nil {
		log.Printf("WARNING witness heartbeats will not record witness versions: %v\n", err)
	}
	return ok
}

// heartbeatVersionField returns s if it is a plausible build string (at most
// maxHeartbeatVersionLen of [0-9A-Za-z.+_-]), otherwise "".
func heartbeatVersionField(s string) string {
	if len(s) > maxHeartbeatVersionLen {
		return ""
	}
	for _, r := range s {
		if (r < 'a' || r > 'z') && (r < 'A' || r > 'Z') && (r < '0' || r > '9') && r != '.' && r != '+' && r != '_' && r != '-' {
			return ""
		}
	}
	return s
}

// isHeartbeatWitness reports whether an EVM heartbeat comes from an address in
// the validator's active witness set. When the set cannot be loaded the
// heartbeat is not recorded.
func isHeartbeatWitness(validator TransferValidator, witnessAddr []byte) bool {
	if len(witnessAddr) != common.AddressLength {
		return false
	}
	member, loaded := validator.IsActiveWitness(common.BytesToAddress(witnessAddr))
	return member && loaded
}

// maxUnverifiedHeartbeatRows is the row limit of witness_heartbeats for relayers
// without a witness set (Solana): beyond it, only rows already in the table are
// updated.
var maxUnverifiedHeartbeatRows = 1024

// heartbeatRowAllowed reports whether a heartbeat of an unverified witness may
// be written: its row exists, or the table is below maxUnverifiedHeartbeatRows.
// The caller serializes check and insert; relayers sharing the DB may still
// overshoot the cap by one row each.
func heartbeatRowAllowed(store *db.SQLStore, witnessAddr []byte, cashier util.Address) (bool, error) {
	var n int
	if err := store.DB().QueryRow(
		fmt.Sprintf("SELECT COUNT(*) FROM `%s` WHERE `witness`=? AND `cashier`=?", witnessHeartbeatTableName),
		"0x"+hex.EncodeToString(witnessAddr), cashier.String(),
	).Scan(&n); err != nil {
		return false, err
	}
	if n > 0 {
		return true, nil
	}
	if err := store.DB().QueryRow(fmt.Sprintf("SELECT COUNT(*) FROM `%s`", witnessHeartbeatTableName)).Scan(&n); err != nil {
		return false, err
	}
	return n < maxUnverifiedHeartbeatRows, nil
}

// upsertWitnessHeartbeat records that a witness pinged this relayer for a
// route, with its scan tip and, if withVersions, its self-reported build.
func upsertWitnessHeartbeat(store *db.SQLStore, withVersions bool, witnessAddr []byte, cashier util.Address, tipHeight uint64, version, commit string) error {
	if !validWitnessHeartbeatAddr(witnessAddr) {
		return nil
	}
	if !withVersions {
		_, err := store.DB().Exec(
			fmt.Sprintf(
				"INSERT INTO `%s` (`witness`, `cashier`, `tipHeight`) VALUES (?, ?, ?) "+
					"ON DUPLICATE KEY UPDATE `tipHeight`=VALUES(`tipHeight`), `updateTime`=CURRENT_TIMESTAMP",
				witnessHeartbeatTableName,
			),
			"0x"+hex.EncodeToString(witnessAddr), cashier.String(), tipHeight,
		)
		return err
	}
	_, err := store.DB().Exec(
		fmt.Sprintf(
			"INSERT INTO `%s` (`witness`, `cashier`, `tipHeight`, `version`, `commitID`) VALUES (?, ?, ?, ?, ?) "+
				"ON DUPLICATE KEY UPDATE `tipHeight`=VALUES(`tipHeight`), `version`=VALUES(`version`), "+
				"`commitID`=VALUES(`commitID`), `updateTime`=CURRENT_TIMESTAMP",
			witnessHeartbeatTableName,
		),
		"0x"+hex.EncodeToString(witnessAddr), cashier.String(), tipHeight,
		heartbeatVersionField(version), heartbeatVersionField(commit),
	)
	return err
}
