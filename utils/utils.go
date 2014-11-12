package utils

import (
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"net/http"
	"strings"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/PuerkitoBio/throttled"
	"github.com/PuerkitoBio/throttled/store"
	"github.com/gorilla/context"
	router "github.com/zenazn/goji/web"
)

func RateLimit(i int, minutes time.Duration) *throttled.Throttler {
	// TODO: Change to redis store
	return throttled.RateLimit(throttled.Q{i, minutes * time.Minute}, &throttled.VaryBy{RemoteAddr: true, Path: true}, store.NewMemStore(1000))
}

// Last function must be the actual handler.
func M(middlewares ...interface{}) http.Handler {
	// get last function
	fh := middlewares[len(middlewares)-1]
	h := handleErr(fh.(func(http.ResponseWriter, *http.Request) error))

	var final http.Handler
	if h != nil {
		final = h
	} else {
		final = http.DefaultServeMux
	}

	// count without last function
	for i := len(middlewares) - 2; i >= 0; i-- {
		final = middlewares[i].(func(http.Handler) http.Handler)(final)
	}

	// TODO: until https://github.com/zenazn/goji/issues/76 will be implemented
	fn := func(c router.C, w http.ResponseWriter, r *http.Request) {
		for k, v := range c.URLParams {
			context.Set(r, "URL"+k, v)
		}
		final.ServeHTTP(w, r)
	}
	return router.HandlerFunc(fn)
}

func M1(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		Log.Debug("middleware")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func NewUUID() string {
	uuid := uuid.New()

	return strings.Replace(uuid, "-", "", -1)
}
