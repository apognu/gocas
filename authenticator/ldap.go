package authenticator

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gocas/config"
	"github.com/mqu/openldap"
)

type Ldap struct{}

func (Ldap) Auth(r *http.Request) (bool, string) {
	u, p := r.FormValue("username"), r.FormValue("password")
	ldap, err := openldap.Initialize(config.Get().Ldap.Host)
	if err != nil {
		logrus.Errorf("cannot connect to LDAP server: %s", err)
		return false, u
	}
	ldap.SetOption(openldap.LDAP_OPT_PROTOCOL_VERSION, openldap.LDAP_VERSION3)
	err = ldap.Bind(fmt.Sprintf("%s=%s,%s", config.Get().Ldap.Dn, u, config.Get().Ldap.Base), p)
	if err != nil {
		logrus.Errorf("cannot connect to LDAP server: %s", err)
		return false, u
	}
	ldap.Close()
	return true, u
}
