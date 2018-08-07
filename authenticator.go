package infobip

import (
	"net/http"
)

// Auth contains token to be sent with every request after authentication.
type Auth struct {
	Token string
}

func (a *Auth) setAuth(r *http.Request) {
	r.Header.Add("Authorization", "IBSSO "+a.Token)
}
