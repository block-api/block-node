package db

import (
	"database/sql"
	"time"
)

// SQLite struct definition
type SQLite struct {
	Db         *sql.DB
	migrations []SQLMigration
}

// Migrations return migrations array
func (s *SQLite) Migrations() []SQLMigration {
	return s.migrations
}

// AddMigration adds SQLMigration to migrations, can pass array of migrations
func (s *SQLite) AddMigration(migration ...SQLMigration) {
	s.migrations = append(s.migrations, migration...)
}

// RunMigrations run migrations - from first one to last
func (s *SQLite) RunMigrations() error {
	if len(s.migrations) < 1 {
		return nil
	}

	for _, migration := range s.migrations {
		if _, err := s.Db.Exec(migration.createQuery); err != nil {
			return err
		}

		_, _ = s.Db.Exec(NewMigrationEntry, migration.name, time.Now().Unix())
	}

	return nil
}

// RevertMigrations revert all migrations in reverse order, from the last one to first
func (s *SQLite) RevertMigrations() error {
	if len(s.migrations) < 1 {
		return nil
	}

	for i := len(s.migrations) - 1; i >= 0; i-- {
		if _, err := s.Db.Exec(s.migrations[i].dropQuery); err != nil {
			return err
		}
	}

	return nil
}
