package middleware

import (
	"net/http"
)

type HandleErr func(http.ResponseWriter, *http.Request) error

func (fn HandleErr) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
