package authenticator

import (
	"net/http"
	"strings"
)

type Dummy struct{}

func (Dummy) Auth(r *http.Request) (bool, string) {
	u, p := r.FormValue("username"), r.FormValue("password")
	if strings.TrimSpace(u) != "" && u == p {
		return true, u
	}
	return false, ""
}
