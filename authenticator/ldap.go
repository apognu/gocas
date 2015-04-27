package authenticator

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gocas/config"
	"github.com/mqu/openldap"
)

type Ldap struct{}

func (Ldap) Auth(u string, p string) (bool, string) {
	ldap, err := openldap.Initialize(config.Get().Ldap.Host)
	if err != nil {
		logrus.Errorf("cannot connect to LDAP server: %s", err)
		return false, ""
	}
	ldap.SetOption(openldap.LDAP_OPT_PROTOCOL_VERSION, openldap.LDAP_VERSION3)
	err = ldap.Bind(fmt.Sprintf("%s=%s,%s", config.Get().Ldap.Dn, u, config.Get().Ldap.Base), p)
	if err != nil {
		return false, ""
	}
	ldap.Close()
	return true, ""
}
