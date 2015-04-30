package authenticator

import (
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gocas/config"
	"github.com/kirves/goradius"
)

type Radius struct{}

func (Radius) Auth(r *http.Request) (bool, string) {
	u, p := r.FormValue("username"), r.FormValue("password")

	fmt.Println(config.Get().Radius.Host)
	rad := goradius.Authenticator(config.Get().Radius.Host, config.Get().Radius.Port, config.Get().Radius.Secret)
	ok, err := rad.Authenticate(u, p)
	if err != nil {
		logrus.Errorf("could not authenticate into RADIUS server: %s", err)
		return false, u
	}
	if ok {
		return true, u
	}

	return false, u
}
