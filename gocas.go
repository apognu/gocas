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

	mux := mux.NewRouter()
	mux.HandleFunc("/login", loginRequestorHandler).Methods("GET")
	mux.HandleFunc("/login", loginAcceptorHandler).Methods("POST")
	mux.HandleFunc("/validate", validateHandler).Methods("GET")
	mux.HandleFunc("/serviceValidate", serviceValidateHandler).Methods("GET")
	mux.HandleFunc("/logout", logoutHandler).Methods("GET")

	logrus.Infof("started gocas CAS server, %s", time.Now())
	http.ListenAndServe("0.0.0.0:8080", mux)
}
