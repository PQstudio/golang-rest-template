package utils

import (
	"database/sql"
	"net/http"
	"strings"

	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"bitbucket.com/aria.pqstudio.pl-api/utils/web"
)

type handleErr func(http.ResponseWriter, *http.Request) error

// handle all errors from application here
func (fn handleErr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {

		errMessage := err.Error()
		if err == sql.ErrNoRows {
			Log.Debug(err.Error())
			HttpError(w, nil, 404)
			return
		}
		if strings.HasPrefix(errMessage, web.ValidationErr) {
			Log.Debug(err.Error())
			HttpError(w, err, 422)
			return
		}
		if strings.HasPrefix(errMessage, web.SerializationErr) {
			Log.Debug(err.Error())
			HttpError(w, nil, 400)
			return
		}
		if strings.HasPrefix(errMessage, web.NotFoundErr) {
			Log.Debug(err.Error())
			HttpError(w, nil, 404)
			return
		}

		Log.Critical(err.Error())
		HttpError(w, err, 500)

		return
	}
}

func HttpError(w http.ResponseWriter, err error, status int) {
	if err != nil {
		w.WriteHeader(status)
		web.ToJSON(w, err)
	} else {
		w.WriteHeader(status)
	}
}
