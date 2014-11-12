package user

import (
	"net/http"

	"bitbucket.com/aria.pqstudio.pl-api/middleware"
	"bitbucket.com/aria.pqstudio.pl-api/utils"
	"bitbucket.com/aria.pqstudio.pl-api/utils/web"

	router "github.com/zenazn/goji/web"

	"bitbucket.com/aria.pqstudio.pl-api/user/model"
	"bitbucket.com/aria.pqstudio.pl-api/user/service"
)

var Routes *router.Mux = router.New()

func init() {
	limit := utils.RateLimit(1000, 60)

	Routes.Get("/me", utils.M(utils.M1, limit.Throttle, middleware.Auth, GetCurrent))
	Routes.Get("/:uid", utils.M(limit.Throttle, middleware.Auth, GetOne))
	Routes.Post("", utils.M(limit.Throttle, Post))
}

func GetCurrent(w http.ResponseWriter, r *http.Request) error {
	userID := web.ContextS(r, "userID")
	user, err := service.GetUser(userID)
	if err != nil {
		return err
	}

	res := web.Whitelist(*user, "UID", "Email", "CreatedAt")
	web.ToJSON(w, &res)
	return nil
}

func GetOne(w http.ResponseWriter, r *http.Request) error {
	uid := web.ContextS(r, "URLuid")

	user, err := service.GetUser(uid)
	if err != nil {
		return err
	}

	res := web.Whitelist(*user, "UID", "Email", "CreatedAt")
	web.ToJSON(w, &res)

	return nil
}

func Post(w http.ResponseWriter, r *http.Request) error {
	var user model.User

	err := web.Bind(w, r.Body, &user)
	if err != nil {
		return err
	}

	err = service.CreateUser(&user)
	if err != nil {
		return err
	}

	res := web.Whitelist(user, "UID", "Email", "CreatedAt")
	web.ToJSON(w, &res)

	return nil
}
