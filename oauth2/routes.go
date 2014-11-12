package oauth2

import (
	"bitbucket.com/aria.pqstudio.pl-api/utils"
	"bitbucket.com/aria.pqstudio.pl-api/utils/web"

	"database/sql"
	"net/http"

	"github.com/RangelReale/osin"
	router "github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"

	"bitbucket.com/aria.pqstudio.pl-api/user/model"
	userService "bitbucket.com/aria.pqstudio.pl-api/user/service"

	accessService "bitbucket.com/aria.pqstudio.pl-api/oauth2/service"

	"bitbucket.com/aria.pqstudio.pl-api/security"
)

var Routes *router.Mux = router.New()

func init() {
	Routes.Use(middleware.SubRouter)

	limit := utils.RateLimit(1000, 60)

	Routes.Post("/token", utils.M(limit.Throttle, Token))
	Routes.Put("/token/invalidate", utils.M(limit.Throttle, TokenInvalidate))
	Routes.Post("/token/check", utils.M(limit.Throttle, TokenCheck))
}

func TokenCheck(w http.ResponseWriter, r *http.Request) error {
	token := AccessToken(r)
	res := make(map[string]string, 1)

	_, err := Server.Storage.LoadAccess(token)
	if err != nil {
		res["status"] = "unauthorized"
	} else {
		res["status"] = "authorized"
	}

	web.ToJSON(w, &res)
	return nil
}

// TODO: needs authentication
func TokenInvalidate(w http.ResponseWriter, r *http.Request) error {
	token := AccessToken(r)

	Server.Storage.RemoveAccess(token)

	web.HttpError(w, nil, http.StatusNoContent)

	return nil
}

func Token(w http.ResponseWriter, r *http.Request) error {
	var data Data

	err := web.FromJSONStrict(w, r.Body, &data)
	if err != nil {
		return err
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
	var user *model.User

	if ar := Server.HandleAccessRequest(resp, r); ar != nil {
		switch ar.Type {
		case osin.PASSWORD:
			user, err = userService.GetUserByEmail(ar.Username)
			if err == sql.ErrNoRows {
				ar.Authorized = false
			} else if err != nil {
				return err
			} else if security.CompareHashAndSalt(ar.Password, user.Salt, user.Password) {
				ar.Authorized = true
			}

		case osin.REFRESH_TOKEN:
			ar.Authorized = true
		}
		Server.FinishAccessRequest(resp, r, ar)
	}
	if resp.IsError && resp.InternalError != nil {
		return resp.InternalError
	}

	if !resp.IsError {
		accessService.UpdateAccessByToken(resp.Output["access_token"].(string), user.UID)
	}
	osin.OutputJSON(resp, w, r)

	return nil
}
