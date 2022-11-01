package db

const CreateMigrationTable string = `
	CREATE TABLE IF NOT EXISTS migration (
		name TEXT NOT NULL PRIMARY KEY,
		created_at INTEGER NOT NULL
	);
`

const NewMigrationEntry string = `
	INSERT INTO migration (name, created_at) VALUES (?, ?)
`

type SQLMigration struct {
	name        string
	createQuery string
	dropQuery   string
}

func NewSQLMigration(name string, createQuery string, dropQuery string) SQLMigration {
	return SQLMigration{
		name:        name,
		createQuery: createQuery,
		dropQuery:   dropQuery,
	}
}
