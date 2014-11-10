package db

import (
	"database/sql"
	"github.com/BurntSushi/migration"

	"bitbucket.com/aria.pqstudio.pl-api/migrations"
)

var DB *sql.DB

func Connect(provider string, dsn string) error {
	migration.DefaultGetVersion = migrations.GetVersion
	migration.DefaultSetVersion = migrations.SetVersion

	database, err := migration.Open(provider, dsn, migrations.Migrations)
	DB = database
	return err
}
