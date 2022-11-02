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

const FindMigrationEntry string = `
	SELECT name FROM migration WHERE name = ?
`

const DeleteMigrationEntry string = `
	DELETE FROM migration WHERE name = ?
`

type SQLMigration struct {
	name      string
	upQuery   string
	downQuery string
}

func NewSQLMigration(name string, upQuery string, downQuery string) SQLMigration {
	return SQLMigration{
		name:      name,
		upQuery:   upQuery,
		downQuery: downQuery,
	}
}
