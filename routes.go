package main

import (
	"bitbucket.com/aria.pqstudio.pl-api/middlewares"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"bitbucket.com/aria.pqstudio.pl-api/utils/web"

	"bitbucket.com/aria.pqstudio.pl-api/users"
	"github.com/zenazn/goji"
	router "github.com/zenazn/goji/web"
	"net/http"
)

func setupRoutes() {
	main := router.New()
	goji.Use(middleware.ContentTypeJSON)

	goji.Handle("/v1*", http.StripPrefix("/v1", main))

	main.Handle("/users*", http.StripPrefix("/users", user.Routes))

	goji.NotFound(NotFound)
	user.Routes.NotFound(NotFound)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	error := &web.Error{Message: "url_not_found"}

	Log.Info("URL not found: %s", r.RequestURI)
	w.WriteHeader(http.StatusNotFound)
	web.ToJSON(w, error)
}
