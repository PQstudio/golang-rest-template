package migration

import (
	"github.com/BurntSushi/migration"
)

func SetupOAuth2(tx migration.LimitedTx) error {
	var stmts = []string{
		clientTable,
		accessTable,
	}
	for _, stmt := range stmts {
		_, err := tx.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: deal with foreign keys
var accessTable = `
CREATE TABLE IF NOT EXISTS access_data (
	 uid             BINARY(16) NOT NULL,
     clientID        BINARY(16),
     authorizeDataID BINARY(16),
     accessDataID    BINARY(16),
     userID          BINARY(16),
     accessToken     VARCHAR(500),
     refreshToken    VARCHAR(500),
     expiresIn       INT,
     scope           VARCHAR(1000),
     redirectURI     VARCHAR(1000),
     createdAt       TIMESTAMP,

     PRIMARY KEY (uid)
) ENGINE=InnoDB;
`

var clientTable = `
CREATE TABLE IF NOT EXISTS client_data (
	 uid          BINARY(16) NOT NULL,
     clientID     VARCHAR(500),
     clientSecret VARCHAR(1000),
     redirectURI  VARCHAR(1000),

     PRIMARY KEY (uid),
     UNIQUE (clientID)
) ENGINE=InnoDB;
`
