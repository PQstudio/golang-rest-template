package main

import (
	"net/http"

	mw "bitbucket.com/aria.pqstudio.pl-api/middleware"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"bitbucket.com/aria.pqstudio.pl-api/utils/web"

	"github.com/gorilla/context"
	"github.com/zenazn/goji"
	router "github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"bitbucket.com/aria.pqstudio.pl-api/oauth2"
	"bitbucket.com/aria.pqstudio.pl-api/user"
)

func setupRoutes() {
	goji.Use(mw.ContentTypeJSON)

	main := router.New()
	main.Use(middleware.SubRouter)

	goji.Handle("/v1/*", main)
	goji.Handle("/oauth2/*", oauth2.Routes)

	// TODO: until https://github.com/zenazn/goji/issues/65 is fixed so we can use routes without trailing slash
	main.Handle("/users", http.StripPrefix("/users", user.Routes))
	main.Handle("/users/*", http.StripPrefix("/users", user.Routes))

	goji.Use(context.ClearHandler)

	goji.NotFound(NotFound)
	main.NotFound(NotFound)

	user.Routes.NotFound(NotFound)
	oauth2.Routes.NotFound(NotFound)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	error := &web.Error{Message: "url_not_found"}

	Log.Info("URL not found: %s", r.RequestURI)
	w.WriteHeader(http.StatusNotFound)
	web.ToJSON(w, error)
}
