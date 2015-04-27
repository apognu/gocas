package authenticator

type Authenticator interface {
	Auth(u string, p string) (bool, string)
}

var AvailableAuthenticators = map[string]Authenticator{
	"dummy": Dummy{},
	"ldap":  Ldap{},
}
