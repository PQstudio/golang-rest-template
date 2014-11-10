package web

import (
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/context"
	"github.com/martini-contrib/binding"
)

type Model interface {
	Validate() binding.Errors
}

type Error struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"errors,omitempty"`
}

func Bind(w http.ResponseWriter, r io.ReadCloser, obj Model) bool {
	if ok := FromJSON(w, r, obj); !ok {
		return false
	}

	if err := obj.Validate(); err != nil {
		// go doesn't have 422 status code so yeah, teapot
		error := &Error{
			Message: "validation_error",
			Errors:  err,
		}

		HttpError(w, error, http.StatusTeapot)
		return false
	}
	return true
}

func FromJSON(w http.ResponseWriter, r io.ReadCloser, obj Model) bool {
	defer r.Close()
	if err := json.NewDecoder(r).Decode(obj); err != nil {
		error := &Error{Message: "deserialization_error"}
		HttpError(w, error, http.StatusBadRequest)

		return false
	}
	return true
}

func FromJSONStrict(w http.ResponseWriter, r io.ReadCloser, obj interface{}) bool {
	defer r.Close()
	if err := json.NewDecoder(r).Decode(obj); err != nil {
		error := &Error{Message: "deserialization_error"}
		HttpError(w, error, http.StatusBadRequest)

		return false
	}
	return true
}

func ToJSON(w http.ResponseWriter, obj interface{}) bool {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		error := &Error{Message: "serialization_error"}
		HttpError(w, error, http.StatusBadRequest)

		return false
	}
	return true
}

func HttpError(w http.ResponseWriter, err *Error, status int) {
	if err != nil {
		Log.Error(err.Message)
		w.WriteHeader(status)
		ToJSON(w, err)
	} else {
		w.WriteHeader(status)
	}
}

func Whitelist(from Model, fields ...string) map[string]interface{} {
	out := make(map[string]interface{})
	obj := reflect.ValueOf(from)
	t := reflect.TypeOf(from)
	for _, v := range fields {
		val := obj.FieldByName(v).Interface()
		if tag, ok := t.FieldByName(v); ok {
			name := strings.Split(tag.Tag.Get("json"), ",")[0]
			out[name] = val
		}
	}

	return out
}

func ContextS(r *http.Request, key string) string {
	return string(context.Get(r, key).(string))
}
