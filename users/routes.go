package user

import (
	. "bitbucket.com/aria.pqstudio.pl-api/utils"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"bitbucket.com/aria.pqstudio.pl-api/utils/web"

	"database/sql"
	"net/http"

	"github.com/justinas/alice"
	router "github.com/zenazn/goji/web"

	"bitbucket.com/aria.pqstudio.pl-api/users/models"
	repo "bitbucket.com/aria.pqstudio.pl-api/users/repositories"
)

var Routes *router.Mux = router.New()

func init() {
	limit := RateLimit(1000, 60)
	middleware := alice.New(M1, limit.Throttle)

	Routes.Get("/me", R(middleware.ThenFunc(GetCurrent)))
	Routes.Get("/:uid", R(middleware.ThenFunc(GetOne)))
	Routes.Post("", R(middleware.ThenFunc(Post)))
}

func GetCurrent(w http.ResponseWriter, r *http.Request) {
	user := user.User{
		Email: "gregory90@gmail.com",
	}
	Log.Debug("CurrentUser: %+v", user)

	res := web.Whitelist(user, "Email")
	web.ToJSON(w, &res)
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	uid := web.ContextS(r, "URLuid")

	user, err := repo.GetOne(uid)
	// TODO: handle database errors somewhere else
	if err == sql.ErrNoRows {
		web.HttpError(w, nil, http.StatusNotFound)
		return
	} else if err != nil {
		web.HttpError(w, nil, http.StatusInternalServerError)
		return
	}

	res := web.Whitelist(*user, "UID", "Email", "CreatedAt")
	web.ToJSON(w, &res)
}

func Post(w http.ResponseWriter, r *http.Request) {
	var user user.User

	if ok := web.Bind(w, r.Body, &user); !ok {
		return
	}

	err := repo.Create(&user)
	if err != nil {
		web.HttpError(w, nil, http.StatusInternalServerError)
		return
	}

	res := web.Whitelist(user, "UID", "Email", "CreatedAt")
	web.ToJSON(w, &res)
}
