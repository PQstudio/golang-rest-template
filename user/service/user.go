package service

import (
	"time"

	"bitbucket.com/aria.pqstudio.pl-api/security"
	"bitbucket.com/aria.pqstudio.pl-api/user/datastore"
	"bitbucket.com/aria.pqstudio.pl-api/user/model"

	"bitbucket.com/aria.pqstudio.pl-api/utils"
)

func GetUser(uid string) (*model.User, error) {
	user, err := datastore.GetUser(uid)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	user, err := datastore.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// TODO: check if email is unique - validation
func CreateUser(user *model.User) error {
	password, salt, err := security.GenerateHashAndSalt(user.Password, 10)
	if err != nil {
		return err
	}

	user.UID = utils.NewUUID()
	user.CreatedAt = time.Now().UTC()
	user.Password = password
	user.Salt = salt

	err = datastore.CreateUser(user)
	if err != nil {
		return err
	}

	return nil
}
