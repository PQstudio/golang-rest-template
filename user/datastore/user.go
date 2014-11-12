package datastore

import (
	"bitbucket.com/aria.pqstudio.pl-api/user/model"

	. "bitbucket.com/aria.pqstudio.pl-api/utils/db"
)

const (
	table string = "users"
)

func GetUser(uid string) (*model.User, error) {
	user := &model.User{}
	err := DB.QueryRow("SELECT lower(hex(uid)), email, password, salt, createdAt FROM "+table+" WHERE uid = unhex(?)", uid).Scan(&user.UID, &user.Email, &user.Password, &user.Salt, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := DB.QueryRow("SELECT lower(hex(uid)), email, password, salt, createdAt FROM "+table+" WHERE email = ?", email).Scan(&user.UID, &user.Email, &user.Password, &user.Salt, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// TODO: check if email is unique - validation
func CreateUser(user *model.User) error {
	stmt, err := DB.Prepare("INSERT " + table + " SET uid=unhex(?),email=?,password=?,salt=?,createdAt=?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(user.UID, user.Email, user.Password, user.Salt, user.CreatedAt)
	return err
}
