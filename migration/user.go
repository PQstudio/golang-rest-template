package migration

import (
	"github.com/BurntSushi/migration"
)

func SetupUser(tx migration.LimitedTx) error {
	var stmts = []string{
		userTable,
	}
	for _, stmt := range stmts {
		_, err := tx.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

var userTable = `
CREATE TABLE IF NOT EXISTS users (
	 uid      BINARY(16) NOT NULL,
     email    VARCHAR(255),
     password VARCHAR(1000),
     salt     VARCHAR(255),
     createdAt TIMESTAMP,

     PRIMARY KEY (uid),
     UNIQUE (email)
) ENGINE=InnoDB;
`
