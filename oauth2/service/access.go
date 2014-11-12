package service

import (
	"bitbucket.com/aria.pqstudio.pl-api/utils"
	//. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"

	"bitbucket.com/aria.pqstudio.pl-api/oauth2/datastore"
	"bitbucket.com/aria.pqstudio.pl-api/oauth2/model"

	"github.com/RangelReale/osin"
)

func GetAccessByToken(token string) (*osin.AccessData, error) {
	access, err := datastore.GetAccessByToken(token)
	if err != nil {
		return nil, err
	}

	client, err := GetClientByUID(access.ClientID)
	if err != nil {
		return nil, err
	}

	a := &osin.AccessData{
		Client:        client,
		AuthorizeData: nil,
		AccessData:    nil,
		AccessToken:   access.AccessToken,
		RefreshToken:  access.RefreshToken,
		ExpiresIn:     access.ExpiresIn,
		Scope:         access.Scope,
		RedirectUri:   access.RedirectUri,
		CreatedAt:     access.CreatedAt,
	}

	return a, nil
}

func CreateAccess(data *osin.AccessData) error {
	client, err := datastore.GetClientByID(data.Client.GetId())
	if err != nil {
		return err
	}

	access := &model.AccessData{
		UID:             utils.NewUUID(),
		ClientID:        client.UID,
		AuthorizeDataID: "",
		AccessDataID:    "",
		AccessToken:     data.AccessToken,
		RefreshToken:    data.RefreshToken,
		ExpiresIn:       data.ExpiresIn,
		Scope:           data.Scope,
		RedirectUri:     data.RedirectUri,
		CreatedAt:       data.CreatedAt,
	}
	err = datastore.CreateAccess(access)
	return err
}

func DeleteAccessByToken(token string) error {
	err := datastore.DeleteAccessByToken(token)
	return err
}

func GetAccessByRefresh(token string) (*osin.AccessData, error) {
	access, err := datastore.GetAccessByRefresh(token)
	if err != nil {
		return nil, err
	}

	client, err := GetClientByUID(access.ClientID)
	if err != nil {
		return nil, err
	}

	a := &osin.AccessData{
		Client:        client,
		AuthorizeData: nil,
		AccessData:    nil,
		AccessToken:   access.AccessToken,
		RefreshToken:  access.RefreshToken,
		ExpiresIn:     access.ExpiresIn,
		Scope:         access.Scope,
		RedirectUri:   access.RedirectUri,
		CreatedAt:     access.CreatedAt,
	}

	return a, nil
}

func UpdateAccessByToken(token string, userUID string) error {
	err := datastore.UpdateAccessByToken(token, userUID)
	return err
}
