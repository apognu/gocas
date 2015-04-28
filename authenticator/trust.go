package authenticator

import "net/http"

type Trust struct{}

func (Trust) Auth(r *http.Request) (bool, string) {
	if r.Header.Get("REMOTE_USER") != "" {
		return true, r.Header.Get("REMOTE_USER")
	}
	if r.Header.Get("REMOTE-USER") != "" {
		return true, r.Header.Get("REMOTE-USER")
	}
	return false, ""
}
