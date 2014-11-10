package oauth2

import (
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"errors"
	"github.com/RangelReale/osin"

	"bitbucket.com/aria.pqstudio.pl-api/oauth2/repositories"
	"database/sql"
)

type MySQLStorage struct {
}

func NewMySQLStorage() *MySQLStorage {
	r := &MySQLStorage{}

	return r
}

func (s *MySQLStorage) Clone() osin.Storage {
	return s
}

func (s *MySQLStorage) Close() {
}

func (s *MySQLStorage) GetClient(id string) (osin.Client, error) {
	Log.Notice("OAuth2, get client: %s\n", id)
	c, _, err := oauth2.GetClientByID(id)

	if err == sql.ErrNoRows {
		return nil, errors.New("Client not found")
	} else if err != nil {
		return nil, errors.New("Server error")
	}
	return c, nil
}

func (s *MySQLStorage) SetClient(id string, client osin.Client) error {
	Log.Notice("OAuth2, set client: %s\n", id)

	c := &osin.DefaultClient{
		Id:          client.GetId(),
		Secret:      client.GetSecret(),
		RedirectUri: client.GetRedirectUri(),
	}

	err := oauth2.CreateClient(c)
	if err != nil {
		return errors.New("Server error")
	}
	return nil
}

func (s *MySQLStorage) SaveAuthorize(data *osin.AuthorizeData) error {
	return errors.New("Not implemented")
}

func (s *MySQLStorage) LoadAuthorize(code string) (*osin.AuthorizeData, error) {
	return nil, errors.New("Not implemented")
}

func (s *MySQLStorage) RemoveAuthorize(code string) error {
	return errors.New("Not implemented")
}

func (s *MySQLStorage) SaveAccess(data *osin.AccessData) error {
	Log.Notice("OAuth2, save access: %s\n", data.AccessToken)

	_, err := oauth2.CreateAccess(data)
	if err != nil {
		return errors.New("Server error")
	}

	return nil
}

func (s *MySQLStorage) LoadAccess(code string) (*osin.AccessData, error) {
	Log.Notice("OAuth2, load access: %s\n", code)

	a, _, err := oauth2.GetAccessByToken(code)
	if err == sql.ErrNoRows {
		return nil, errors.New("Access not found")
	} else if err != nil {
		return nil, errors.New("Server error")
	}

	return a, nil
}

func (s *MySQLStorage) RemoveAccess(code string) error {
	Log.Notice("OAuth2, remove access: %s\n", code)

	err := oauth2.DeleteAccessByToken(code)
	if err != nil {
		return errors.New("Server error")
	}

	return nil
}

func (s *MySQLStorage) LoadRefresh(code string) (*osin.AccessData, error) {
	Log.Notice("OAuth2, load refresh: %s\n", code)
	a, _, err := oauth2.GetAccessByRefresh(code)
	Log.Debug("%#v", a)
	if err == sql.ErrNoRows {
		return nil, errors.New("Refresh not found")
	} else if err != nil {
		return nil, errors.New("Server error")
	}

	return a, nil
}

func (s *MySQLStorage) RemoveRefresh(code string) error {
	Log.Notice("OAuth2, remove refresh: %s\n", code)
	return nil
}
