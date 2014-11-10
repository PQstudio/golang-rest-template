package migrations

import (
	"github.com/BurntSushi/migration"
)

func SetupOAuth2(tx migration.LimitedTx) error {
	var stmts = []string{
		clientTable,
		authorizeTable,
		accessTable,
		//refreshTable,
	}
	for _, stmt := range stmts {
		_, err := tx.Exec(stmt)
		if err != nil {
			return err
		}
	}
	return nil
}

var accessTable = `
CREATE TABLE IF NOT EXISTS access_data (
	 uid             BINARY(16) NOT NULL,
     clientID        BINARY(16),
     authorizeDataID BINARY(16),
     accessDataID    BINARY(16),
     accessToken     VARCHAR(500),
     refreshToken    VARCHAR(500),
     expiresIn       INT,
     scope           VARCHAR(1000),
     redirectURI     VARCHAR(1000),
     createdAt       TIMESTAMP,

     PRIMARY KEY (uid),
     FOREIGN KEY(authorizeDataID) REFERENCES authorize_data(uid)
) ENGINE=InnoDB;
`

//var refreshTable = `
//CREATE TABLE IF NOT EXISTS refresh_data (
//uid             BINARY(16) NOT NULL,
//clientID        BINARY(16),
//accessToken     VARCHAR(500),
//expiresIn       INT,
//createdAt       TIMESTAMP,

//PRIMARY KEY (uid),
//FOREIGN KEY(clientID) REFERENCES client_data(uid),
//) ENGINE=InnoDB;
//`

// not used right now
var authorizeTable = `
CREATE TABLE IF NOT EXISTS authorize_data (
uid         BINARY(16) NOT NULL,
clientID    BINARY(16),
code        VARCHAR(500),
state       VARCHAR(500),
expiresIn   INT,
scope       VARCHAR(1000),
redirectURI VARCHAR(1000),
createdAt   TIMESTAMP,

PRIMARY KEY (uid),
UNIQUE (code),
FOREIGN KEY(clientID) REFERENCES client_data(uid)
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
