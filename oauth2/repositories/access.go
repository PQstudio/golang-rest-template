package oauth2

import (
	. "bitbucket.com/aria.pqstudio.pl-api/utils"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/db"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"

	"database/sql"
	"errors"
	"github.com/RangelReale/osin"
)

const (
	accessTable string = "access_data"
)

func GetAccessByToken(token string) (*osin.AccessData, string, error) {
	access := &osin.AccessData{}
	var clientID, uid string

	// get access
	err := DB.QueryRow("SELECT lower(hex(uid)), lower(hex(clientID)), accessToken, refreshToken, expiresIn, scope, redirectURI, createdAt FROM "+accessTable+" WHERE accessToken = ?", token).Scan(&uid, &clientID, &access.AccessToken, &access.RefreshToken, &access.ExpiresIn, &access.Scope, &access.RedirectUri, &access.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, "", err
	} else if err != nil {
		Log.Error(err.Error())
		return nil, "", err
	}

	c, _, err := GetClientByUID(clientID)
	if err == sql.ErrNoRows {
		return nil, "", errors.New("Client not found")
	} else if err != nil {
		return nil, "", errors.New("Server error")
	}
	access.Client = c

	access.AuthorizeData = nil
	access.AccessData = nil

	return access, uid, nil
}

func CreateAccess(data *osin.AccessData) (string, error) {
	stmt, err := DB.Prepare("INSERT " + accessTable + " SET uid=unhex(?),clientID=unhex(?),authorizeDataID=unhex(?),accessDataID=unhex(?),accessToken=?,refreshToken=?,expiresIn=?,scope=?,redirectURI=?,createdAt=?")
	if err != nil {
		Log.Error(err.Error())
		return "", err
	}
	defer stmt.Close()

	var clientID string

	if data.Client == nil {
		return "", errors.New("Client missing")
	}

	_, clientID, err = GetClientByID(data.Client.GetId())
	if err == sql.ErrNoRows {
		return "", errors.New("Client not found")
	} else if err != nil {
		return "", errors.New("Server error")
	}

	uid := NewUUID()
	if data.RefreshToken != "" {
		_, err = stmt.Exec(uid, clientID, nil, nil, data.AccessToken, data.RefreshToken, data.ExpiresIn, data.Scope, data.RedirectUri, data.CreatedAt)
		if err != nil {
			Log.Error(err.Error())
			return "", err
		}

	} else {
		_, err = stmt.Exec(uid, clientID, nil, nil, data.AccessToken, nil, data.ExpiresIn, data.Scope, data.RedirectUri, data.CreatedAt)
		if err != nil {
			Log.Error(err.Error())
			return "", err
		}
	}

	return uid, nil
}

func DeleteAccessByToken(token string) error {
	stmt, err := DB.Prepare("DELETE FROM " + accessTable + " WHERE accessToken=?")
	if err != nil {
		Log.Error(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(token)
	if err != nil {
		Log.Error(err.Error())
		return err
	}

	return nil
}

func GetAccessByRefresh(token string) (*osin.AccessData, string, error) {
	access := &osin.AccessData{}
	var clientUID, uid string

	// get access
	err := DB.QueryRow("SELECT lower(hex(uid)), lower(hex(clientID)), accessToken, refreshToken, expiresIn, scope, redirectURI, createdAt FROM "+accessTable+" WHERE refreshToken = ?", token).Scan(&uid, &clientUID, &access.AccessToken, &access.RefreshToken, &access.ExpiresIn, &access.Scope, &access.RedirectUri, &access.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, "", err
	} else if err != nil {
		Log.Error(err.Error())
		return nil, "", err
	}

	c, _, err := GetClientByUID(clientUID)
	if err == sql.ErrNoRows {
		return nil, "", errors.New("Client not found")
	} else if err != nil {
		return nil, "", errors.New("Server error")
	}
	access.Client = c

	access.AuthorizeData = nil
	access.AccessData = nil

	return access, uid, nil
}
