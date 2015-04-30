package main

import (
	"flag"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gocas/authenticator"
	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/interceptor"
	"github.com/apognu/gocas/protocol/cas"
	"github.com/apognu/gocas/protocol/oauth"
	"github.com/apognu/gocas/util"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type Protocol func(*mux.Router)

var AvailableProtocols = map[string]Protocol{
	"cas":   cas.New,
	"oauth": oauth.New,
}

var (
	c = flag.String("config", "/etc/gocas.yaml", "path to GoCAS configuration file")
)

func redirect(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", util.Url("/login"))
	w.WriteHeader(http.StatusFound)
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	config.Set(*c)
	if AvailableProtocols[config.Get().Protocol] == nil {
		logrus.Fatalf("unknown protocol: %s", config.Get().Protocol)
	}
	if authenticator.AvailableAuthenticators[config.Get().Authenticator] == nil {
		logrus.Fatalf("unknown authenticator: %s", config.Get().Authenticator)
	}

	if stat, err := os.Stat(config.Get().TemplatePath); !os.IsNotExist(err) {
		if !stat.IsDir() {
			logrus.Fatalf("template path %s is not a directory", config.Get().TemplatePath)
		}
	} else {
		logrus.Fatalf("template path does not exist: %s", err)
	}

	r := mux.NewRouter().StrictSlash(true)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(config.Get().TemplatePath))))
	r.HandleFunc("/", redirect)

	sr := r
	if config.Get().UrlPrefix != "" {
		sr = r.PathPrefix(config.Get().UrlPrefix).Subrouter()
		sr.HandleFunc("/", redirect)
	}

	sr.HandleFunc("/validate", validateHandler).Methods("GET")
	sr.HandleFunc("/serviceValidate", serviceValidateHandler).Methods("GET")
	sr.HandleFunc("/logout", logoutHandler).Methods("GET")

	sr.HandleFunc("/proxy", proxyHandler).Methods("GET")
	sr.HandleFunc("/proxyValidate", proxyValidateHandler).Methods("GET")

	if config.Get().RestApi {
		sr.HandleFunc("/v1/tickets", restGetTicketGrantingTicketHandler).Methods("POST")
		sr.HandleFunc("/v1/tickets/{ticket}", restGetServiceTicketHandler).Methods("POST")
		sr.HandleFunc("/v1/tickets/{ticket}", restLogoutHandler).Methods("DELETE")
	}

	AvailableProtocols[config.Get().Protocol](sr)

	logrus.Infof("started gocas CAS server, %s", time.Now())

	n := negroni.New()
	for _, in := range interceptor.AvailableInterceptors {
		in.Init()
		n.UseFunc(in.Intercept)
	}
	n.UseHandler(r)

	logrus.Fatalf("could not start server: %s", http.ListenAndServe(config.Get().Listen, n))
}
