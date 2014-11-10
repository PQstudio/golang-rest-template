package oauth2

import (
	. "bitbucket.com/aria.pqstudio.pl-api/utils"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/db"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"

	"database/sql"
	"github.com/RangelReale/osin"
)

const (
	clientTable string = "client_data"
)

func GetClientByUID(id string) (*osin.DefaultClient, string, error) {
	client := &osin.DefaultClient{}

	var uid string

	Log.Debug("%#v", id)
	// get client
	err := DB.QueryRow("SELECT lower(hex(uid)), clientID, clientSecret, redirectURI FROM "+clientTable+" WHERE uid = unhex(?)", id).Scan(&uid, &client.Id, &client.Secret, &client.RedirectUri)
	if err == sql.ErrNoRows {
		return nil, "", err
	} else if err != nil {
		Log.Error(err.Error())
		return nil, "", err
	}

	return client, uid, nil
}

func GetClientByID(id string) (*osin.DefaultClient, string, error) {
	client := &osin.DefaultClient{}

	var uid string

	// get client
	err := DB.QueryRow("SELECT lower(hex(uid)), clientID, clientSecret, redirectURI FROM "+clientTable+" WHERE clientID = ?", id).Scan(&uid, &client.Id, &client.Secret, &client.RedirectUri)
	if err == sql.ErrNoRows {
		return nil, "", err
	} else if err != nil {
		Log.Error(err.Error())
		return nil, "", err
	}

	return client, uid, nil
}

func CreateClient(client *osin.DefaultClient) error {
	stmt, err := DB.Prepare("INSERT " + clientTable + " SET uid=unhex(?),clientID=?,clientSecret=?,redirectURI=?")
	if err != nil {
		Log.Error(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(NewUUID(), client.Id, client.Secret, client.RedirectUri)
	if err != nil {
		Log.Error(err.Error())
		return err
	}

	return nil
}
