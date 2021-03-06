package model

import (
	"github.com/jamieomatthews/validation"
	"github.com/martini-contrib/binding"
	"time"
)

type User struct {
	UID       string    `json:"uid"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Salt      string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

func (user User) Validate() binding.Errors {
	var errors binding.Errors

	v := validation.NewValidation(&errors, user)
	v.KeyTag("json")

	v.Validate(&user.Email).Message("required").Required()
	v.Validate(&user.Email).Message("incorrect").Email()

	v.Validate(&user.Password).Message("required").Required()
	v.Validate(&user.Password).Message("range").Range(6, 60)

	return *v.Errors.(*binding.Errors)
}
