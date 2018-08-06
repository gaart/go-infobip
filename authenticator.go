package infobip

import (
	"net/http"
)

type Auth struct {
	Token string
}

func (a *Auth) SetAuth(r *http.Request) {
	r.Header.Add("Authorization", "IBSSO "+a.Token)
}
