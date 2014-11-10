package user

import (
	. "bitbucket.com/aria.pqstudio.pl-api/utils"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/db"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"

	"bitbucket.com/aria.pqstudio.pl-api/security"
	"bitbucket.com/aria.pqstudio.pl-api/users/models"
	"database/sql"
	"time"
)

const (
	table string = "users"
)

func GetOne(uid string) (*user.User, error) {
	user := &user.User{}
	err := DB.QueryRow("SELECT lower(hex(uid)), email, password, salt, createdAt FROM "+table+" WHERE uid = unhex(?)", uid).Scan(&user.UID, &user.Email, &user.Password, &user.Salt, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		Log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func GetOneByEmail(email string) (*user.User, error) {
	user := &user.User{}
	err := DB.QueryRow("SELECT lower(hex(uid)), email, password, salt, createdAt FROM "+table+" WHERE email = ?", email).Scan(&user.UID, &user.Email, &user.Password, &user.Salt, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, err
	} else if err != nil {
		Log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

// TODO: check if email is unique - validation
func Create(user *user.User) error {
	stmt, err := DB.Prepare("INSERT " + table + " SET uid=unhex(?),email=?,password=?,salt=?,createdAt=?")
	if err != nil {
		Log.Error(err.Error())
		return err
	}
	defer stmt.Close()

	password, salt, err := security.GenerateHashAndSalt(user.Password, 10)
	if err != nil {
		Log.Error(err.Error())
		return err
	}

	user.UID = NewUUID()
	user.CreatedAt = time.Now().UTC()
	_, err = stmt.Exec(user.UID, user.Email, password, salt, user.CreatedAt)
	if err != nil {
		Log.Error(err.Error())
		return err
	}

	return nil
}
