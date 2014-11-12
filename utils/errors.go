package utils

import (
	"net/http"

	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"database/sql"
)

type handleErr func(http.ResponseWriter, *http.Request) error

// handle all errors from application here
func (fn handleErr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		switch err {
		case sql.ErrNoRows:
			Log.Debug(err.Error())
			http.Error(w, "", 404)
		default:
			Log.Critical(err.Error())
			http.Error(w, err.Error(), 500)
		}
		return
	}
}
