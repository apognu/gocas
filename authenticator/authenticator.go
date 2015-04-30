package authenticator

import "net/http"

type Authenticator interface {
	Auth(*http.Request) (bool, string)
}

var AvailableAuthenticators = map[string]Authenticator{
	"dummy":  Dummy{},
	"ldap":   Ldap{},
	"radius": Radius{},
	"trust":  Trust{},
}
