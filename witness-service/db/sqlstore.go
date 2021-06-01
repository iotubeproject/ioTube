// Copyright (c) 2019 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package db

import (
	"context"
	"database/sql"
	"time"
)

type (
	// Config is the config of database
	Config struct {
		URI    string `json:"uri" yaml:"uri"`
		Driver string `json:"driver" yaml:"driver"`
	}

	// SQLStore is local sqlite3
	SQLStore struct {
		db  *sql.DB
		cfg Config
	}
)

// NewStore instantiates a store
func NewStore(cfg Config) *SQLStore {
	return &SQLStore{db: nil, cfg: cfg}
}

// Start opens the SQL (creates new file if not existing yet)
func (s *SQLStore) Start(_ context.Context) error {
	if s.db != nil {
		return nil
	}
	db, err := sql.Open(s.cfg.Driver, s.cfg.URI)
	if err != nil {
		return err
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)
	db.SetConnMaxLifetime(time.Minute * 60)
	s.db = db
	return nil
}

// Stop closes the SQL
func (s *SQLStore) Stop(_ context.Context) error {
	if s.db != nil {
		if err := s.db.Close(); err != nil {
			return err
		}
	}
	s.db = nil
	return nil
}

// DB returns the sql db
func (s *SQLStore) DB() *sql.DB {
	return s.db
}

// DriverName returns the name of the driver
func (s *SQLStore) DriverName() string {
	return s.cfg.Driver
}
