package oauth2

import (
	. "bitbucket.com/aria.pqstudio.pl-api/utils"
	//. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"bitbucket.com/aria.pqstudio.pl-api/utils/web"

	"database/sql"
	"net/http"

	"github.com/RangelReale/osin"
	"github.com/justinas/alice"
	router "github.com/zenazn/goji/web"

	//"bitbucket.com/aria.pqstudio.pl-api/users/models"
	userRepo "bitbucket.com/aria.pqstudio.pl-api/users/repositories"

	"bitbucket.com/aria.pqstudio.pl-api/security"
)

var Routes *router.Mux = router.New()

func init() {
	limit := RateLimit(1000, 60)
	middleware := alice.New(M1, limit.Throttle)

	Routes.Post("/token", R(middleware.ThenFunc(Token)))
	Routes.Put("/token/invalidate", R(middleware.ThenFunc(TokenInvalidate)))
	Routes.Post("/token/check", R(middleware.ThenFunc(TokenCheck)))
}

func TokenCheck(w http.ResponseWriter, r *http.Request) {
	token := AccessToken(r)
	res := make(map[string]string, 1)

	_, err := Server.Storage.LoadAccess(token)
	if err != nil {
		res["status"] = "unauthorized"
	} else {
		res["status"] = "authorized"
	}

	web.ToJSON(w, &res)
}

// TODO: needs authentication
func TokenInvalidate(w http.ResponseWriter, r *http.Request) {
	token := AccessToken(r)

	Server.Storage.RemoveAccess(token)

	web.HttpError(w, nil, http.StatusNoContent)
}

func Token(w http.ResponseWriter, r *http.Request) {
	var data Data

	if ok := web.FromJSONStrict(w, r.Body, &data); !ok {
		return
	}
	r.ParseForm()
	r.Form.Add("grant_type", data.GrantType)
	r.Form.Add("username", data.Username)
	r.Form.Add("password", data.Password)
	r.Form.Add("scope", data.Scope)
	r.Form.Add("refresh_token", data.RefreshToken)

	r.SetBasicAuth(data.ClientID, data.ClientSecret)

	resp := Server.NewResponse()
	defer resp.Close()
	if ar := Server.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.PASSWORD:
			user, err := userRepo.GetOneByEmail(ar.Username)
			if err == sql.ErrNoRows {
				ar.Authorized = false
			} else if err != nil {
				web.HttpError(w, nil, http.StatusInternalServerError)
				return
			} else if security.CompareHashAndSalt(ar.Password, user.Salt, user.Password) {
				ar.Authorized = true
			}

		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		}
		Server.FinishAccessRequest(resp, r, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		web.HttpError(w, nil, http.StatusInternalServerError)
	}
	osin.OutputJSON(resp, w, r)
}
