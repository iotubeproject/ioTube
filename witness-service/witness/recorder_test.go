// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package witness

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"github.com/iotexproject/ioTube/witness-service/db"
)

func TestRecorder(t *testing.T) {
	require := require.New(t)
	dbFile := "test.db"
	if _, err := os.Stat(dbFile); err == nil {
		require.NoError(os.Remove(dbFile))
	}
	store := db.NewStore("sqlite3", dbFile)
	recorder := NewRecorder(store, "records")
	require.NoError(recorder.Start(context.Background()))
	defer func() {
		require.NoError(recorder.Stop(context.Background()))
		if _, err := os.Stat(dbFile); err == nil {
			require.NoError(os.Remove(dbFile))
		}
	}()
	stmt, err := recorder.store.DB().Prepare(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS records (
               id BIGINT NOT NULL,
			   sender VARCHAR(42) NOT NULL,
			   recipient VARCHAR(42) NOT NULL,
               amount VARCHAR(28) NOT NULL,
               creationTime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
               updateTime DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
               status TINYINT NOT NULL DEFAULT %d,
               txHash VARCHAR(64) NULL,
               notes VARCHAR(45) NULL,
               PRIMARY KEY (id));
               CREATE TRIGGER records_update_trigger
               AFTER UPDATE ON records
               FOR EACH ROW
               WHEN NEW.updateTime <= OLD.updateTime
               BEGIN
                       UPDATE records SET updateTime=CURRENT_TIMESTAMP WHERE id=OLD.id;
               END`,
		New,
	))
	defer stmt.Close()
	require.NoError(err)
	_, err = stmt.Exec()
	require.NoError(err)

	m, err := recorder.NextIDToFetch()
	require.NoError(err)
	require.Equal(0, big.NewInt(0).Cmp(m))
	r1 := &TxRecord{id: big.NewInt(1), recipient: "recipient1", amount: big.NewInt(10)}
	r2 := &TxRecord{id: big.NewInt(2), recipient: "recipient2", amount: big.NewInt(20)}
	require.NoError(recorder.Create(r1))
	rs, err := recorder.NewRecords(10)
	require.NoError(err)
	require.Equal(1, len(rs))
	require.Error(recorder.Create(r1))
	require.NoError(recorder.Create(r2))
	rs, err = recorder.NewRecords(10)
	require.NoError(err)
	require.Equal(2, len(rs))
	require.Equal(r1, rs[0])
	require.Equal(r2, rs[1])
	require.NoError(recorder.StartProcess(r1))
	rs, err = recorder.NewRecords(10)
	require.NoError(err)
	require.Equal(1, len(rs))
	require.NoError(recorder.StartProcess(r2))
	require.Error(recorder.StartProcess(r2))
	require.NoError(recorder.MarkAsSubmitted(r1, "tx1"))
	require.Error(recorder.MarkAsSubmitted(r1, "tx1"))
	require.NoError(recorder.MarkAsSubmitted(r2, "tx2"))
	rs, err = recorder.NewRecords(1)
	require.NoError(err)
	require.Equal(0, len(rs))
	rs, err = recorder.RecordsToConfirm(10, 1)
	require.NoError(err)
	require.Equal(0, len(rs))
	rs, err = recorder.RecordsToConfirm(0, 10)
	require.NoError(err)
	require.Equal(2, len(rs))
	r1.txhash = "tx1"
	require.Equal(r1, rs[0])
	r2.txhash = "tx2"
	require.Equal(r2, rs[1])
	require.NoError(recorder.Confirm(r1))
	rs, err = recorder.RecordsToConfirm(0, 10)
	require.NoError(err)
	require.Equal(1, len(rs))
	require.NoError(recorder.Confirm(r2))
	rs, err = recorder.RecordsToConfirm(0, 1)
	require.NoError(err)
	require.Equal(0, len(rs))
	require.Error(recorder.Confirm(r2))
	m, err = recorder.NextIDToFetch()
	require.NoError(err)
	require.Equal(0, big.NewInt(3).Cmp(m))
}
