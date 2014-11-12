package middleware

import (
	"net/http"

	"github.com/gorilla/context"

	"bitbucket.com/aria.pqstudio.pl-api/oauth2"
	"bitbucket.com/aria.pqstudio.pl-api/utils/web"

	"bitbucket.com/aria.pqstudio.pl-api/oauth2/datastore"
	"bitbucket.com/aria.pqstudio.pl-api/user/service"
)

func Auth(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		token := oauth2.AccessToken(r)
		access, err := datastore.GetAccessByToken(token)
		if err != nil {
			web.HttpError(w, nil, http.StatusForbidden)
			return
		}

		user, err := service.GetUser(access.UserID)
		if err != nil {
			web.HttpError(w, nil, http.StatusForbidden)
			return
		}
		context.Set(r, "userID", user.UID)
		context.Set(r, "userEmail", user.Email)

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
