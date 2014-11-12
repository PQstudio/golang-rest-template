package db

import (
	"database/sql"
	"github.com/BurntSushi/migration"

	migrate "bitbucket.com/aria.pqstudio.pl-api/migration"
)

var DB *sql.DB

func Connect(provider string, dsn string) error {
	migration.DefaultGetVersion = migrate.GetVersion
	migration.DefaultSetVersion = migrate.SetVersion

	database, err := migration.Open(provider, dsn, migrate.Migrations)
	DB = database
	return err
}
