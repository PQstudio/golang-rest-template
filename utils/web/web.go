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

func (e *Error) Error() string {
	return e.Message
}

func Bind(w http.ResponseWriter, r io.ReadCloser, obj Model) error {
	err := FromJSON(w, r, obj)
	if err != nil {
		return err
	}

	if err := obj.Validate(); err != nil {
		return &Error{
			Message: "validation_error",
			Errors:  err,
		}
	}
	return nil
}

func FromJSON(w http.ResponseWriter, r io.ReadCloser, obj Model) error {
	defer r.Close()
	if err := json.NewDecoder(r).Decode(obj); err != nil {
		return &Error{Message: "deserialization_error"}
	}
	return nil
}

func FromJSONStrict(w http.ResponseWriter, r io.ReadCloser, obj interface{}) error {
	defer r.Close()
	if err := json.NewDecoder(r).Decode(obj); err != nil {
		return &Error{Message: "deserialization_error"}
	}
	return nil
}

func ToJSON(w http.ResponseWriter, obj interface{}) error {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		return &Error{Message: "serialization_error"}
	}
	return nil
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
