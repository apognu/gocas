package main

import (
	"fmt"
	"net/http"

	"github.com/apognu/gocas/authenticator"
	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/ticket"
	"github.com/apognu/gocas/util"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func restGetTicketGrantingTicketHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("username") == "" || r.FormValue("password") == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auth, u := authenticator.AvailableAuthenticators[config.Get().Authenticator].Auth(r)
	if !auth {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	tgt := ticket.NewTicketGrantingTicket(u, util.GetRemoteAddr(r.RemoteAddr))
	w.WriteHeader(http.StatusCreated)
	w.Header().Add("Location", fmt.Sprintf("%s%s/%s", config.Get().Url, r.RequestURI, tgt.Ticket))
}

func restGetServiceTicketHandler(w http.ResponseWriter, r *http.Request) {
	tgt := mux.Vars(r)["ticket"]
	svc := r.FormValue("service")

	if tgt == "" || svc == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var tkt ticket.TicketGrantingTicket
	util.GetPersistence("tgt").Find(bson.M{"_id": tgt, "client_ip": util.GetRemoteAddr(r.RemoteAddr)}).One(&tkt)
	if tgt != tkt.Ticket {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	st := ticket.NewServiceTicket(tkt.Ticket, svc, true)
	if !st.Validate() {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.Write([]byte(st.Ticket))
}

func restLogoutHandler(w http.ResponseWriter, r *http.Request) {
	tgt := mux.Vars(r)["ticket"]

	util.GetPersistence("tgt").Remove(bson.M{"_id": tgt, "client_ip": util.GetRemoteAddr(r.RemoteAddr)})
}
