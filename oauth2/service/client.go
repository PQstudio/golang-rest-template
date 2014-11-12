package service

import (
	"bitbucket.com/aria.pqstudio.pl-api/utils"
	//. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"

	"bitbucket.com/aria.pqstudio.pl-api/oauth2/datastore"
	"bitbucket.com/aria.pqstudio.pl-api/oauth2/model"

	"github.com/RangelReale/osin"
)

func GetClientByUID(uid string) (*osin.DefaultClient, error) {
	client, err := datastore.GetClientByUID(uid)
	if err != nil {
		return nil, err
	}

	c := &osin.DefaultClient{
		Id:          client.Id,
		Secret:      client.Secret,
		RedirectUri: client.RedirectUri,
	}

	return c, nil
}

func GetClientByID(id string) (*osin.DefaultClient, error) {
	client, err := datastore.GetClientByID(id)
	if err != nil {
		return nil, err
	}

	c := &osin.DefaultClient{
		Id:          client.Id,
		Secret:      client.Secret,
		RedirectUri: client.RedirectUri,
	}

	return c, nil
}

func CreateClient(client *osin.DefaultClient) error {
	c := &model.Client{
		UID:         utils.NewUUID(),
		Id:          client.Id,
		Secret:      client.Secret,
		RedirectUri: client.RedirectUri,
	}
	err := datastore.CreateClient(c)
	return err
}
