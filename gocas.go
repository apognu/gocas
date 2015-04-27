package main

import (
	"flag"
	"math/rand"
	"net/http"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/apognu/gocas/util"
	"github.com/gorilla/mux"
)

var (
	config = flag.String("config", "/etc/gocas.yaml", "path to GoCAS configuration file")
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	util.SetConfig(*config)

	r := mux.NewRouter()
	prefix := util.GetConfig().UrlPrefix
	sr := r
	if prefix != "" {
		sr = r.PathPrefix(prefix).Subrouter()
	}
	sr.HandleFunc("/login", loginRequestorHandler).Methods("GET")
	sr.HandleFunc("/login", loginAcceptorHandler).Methods("POST")
	sr.HandleFunc("/validate", validateHandler).Methods("GET")
	sr.HandleFunc("/serviceValidate", serviceValidateHandler).Methods("GET")
	sr.HandleFunc("/logout", logoutHandler).Methods("GET")

	logrus.Infof("started gocas CAS server, %s", time.Now())
	http.ListenAndServe("0.0.0.0:8080", r)
}
