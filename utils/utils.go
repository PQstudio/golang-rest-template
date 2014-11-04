package utils

import (
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"github.com/PuerkitoBio/throttled"
	"github.com/PuerkitoBio/throttled/store"
	"github.com/gorilla/context"
	router "github.com/zenazn/goji/web"
	"net/http"
	"time"
)

func RateLimit(i int, minutes time.Duration) *throttled.Throttler {
	// TODO: Change to redis store
	return throttled.RateLimit(throttled.Q{i, minutes * time.Minute}, &throttled.VaryBy{RemoteAddr: true, Path: true}, store.NewMemStore(1000))
}

func R(f http.Handler) router.Handler {
	fn := func(c router.C, w http.ResponseWriter, r *http.Request) {
		for k, v := range c.URLParams {
			context.Set(r, "URL"+k, v)
		}
		f.ServeHTTP(w, r)
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
