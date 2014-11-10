package oauth2

import (
	//. "bitbucket.com/aria.pqstudio.pl-api/utils/logger"
	"bitbucket.com/aria.pqstudio.pl-api/oauth2/storage"
	"github.com/RangelReale/osin"
	"net/http"
)

var Server *osin.Server

type Data struct {
	GrantType    string `json:"grant_type,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
	Scope        string `json:"scope,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

func Init() {
	sconfig := osin.NewServerConfig()
	sconfig.AllowedAuthorizeTypes = osin.AllowedAuthorizeType{osin.TOKEN}
	sconfig.AllowedAccessTypes = osin.AllowedAccessType{osin.REFRESH_TOKEN, osin.PASSWORD, osin.ASSERTION}
	sconfig.AllowGetAccessRequest = true
	Server = osin.NewServer(sconfig, oauth2.NewMySQLStorage())
}

func AccessToken(r *http.Request) string {
	auth := r.Header["Authorization"]
	var token string
	if len(auth) > 0 {
		token = auth[0][7:]
	}

	return token
}
