package user

import (
	. "bitbucket.com/aria.pqstudio.pl-api/utils"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/db"
	. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"bitbucket.com/aria.pqstudio.pl-api/utils/web"

	"database/sql"
	"net/http"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
	router "github.com/zenazn/goji/web"

	. "bitbucket.com/aria.pqstudio.pl-api/users/models"
)

var Routes *router.Mux = router.New()

func init() {
	limit := RateLimit(1000, 60)
	middleware := alice.New(M1, limit.Throttle)

	Routes.Get("/me", R(middleware.ThenFunc(GetCurrent)))

	Routes.Get("/:id", R(middleware.ThenFunc(GetOne)))

	Routes.Post("", R(middleware.ThenFunc(Post)))
}

func GetCurrent(w http.ResponseWriter, r *http.Request) {
	var user User = User{
		Email: "gregory90@gmail.com",
	}
	Log.Debug("CurrentUser: %+v", user)
	Log.Debug("%+v", context.Get(r, "URLParams"))

	res := web.Whitelist(user, "Email")
	web.ToJSON(w, &res)
}

func GetOne(w http.ResponseWriter, r *http.Request) {
	id := context.Get(r, "URLid")
	stmt, err := DB.Prepare("SELECT username, departname FROM userinfo WHERE uid = ?")
	if err != nil {
		Log.Error("%+v", err)
	}

	var username, departname string
	err = stmt.QueryRow(id).Scan(&username, &departname)
	if err == sql.ErrNoRows {
		error := web.Error{Message: "not_found"}
		web.HttpError(w, error, http.StatusNotFound)
		return
	} else if err != nil {
		Log.Error("%+v", err)
	}

	user := User{
		Email:    username,
		Password: departname,
	}

	res := web.Whitelist(user, "Email", "Password")
	web.ToJSON(w, &res)
}

func Post(w http.ResponseWriter, r *http.Request) {
	var user User

	if ok := web.Bind(w, r.Body, &user); !ok {
		return
	}

	stmt, err := DB.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
	if err != nil {
		Log.Debug("%+v", err)
	}

	result, err := stmt.Exec(user.Email, user.Password, "2014-10-09")
	if err != nil {
		Log.Debug("%+v", err)
	}
	Log.Debug("%+v", result)

	res := web.Whitelist(user, "Email")
	web.ToJSON(w, &res)
}
