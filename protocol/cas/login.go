package cas

import (
	"net/http"
	"time"

	"github.com/apognu/gocas/authenticator"
	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/ticket"
	"github.com/apognu/gocas/util"
	"gopkg.in/mgo.v2/bson"
)

const template = "template/login.tmpl"

func forbidden(w http.ResponseWriter, svc string, msg string) {
	lt := ticket.NewLoginTicket(svc)
	w.WriteHeader(http.StatusForbidden)
	lt.Serve(w, template, util.LoginRequestorData{
		Config:   config.Get(),
		Session:  util.LoginRequestorSession{Service: svc},
		Message:  util.LoginRequestorMessage{Type: "danger", Message: msg},
		ShowForm: true})
}

func loginRequestorHandler(w http.ResponseWriter, r *http.Request) {
	tgt, err := r.Cookie("CASTGC")
	svc := r.FormValue("service")
	renew, gateway := r.FormValue("renew"), r.FormValue("gateway")
	lt := ticket.NewLoginTicket(svc)

	// The client sent us a TGT, do not display login form
	if err == nil && renew != "true" {
		var tkt ticket.TicketGrantingTicket
		util.GetPersistence("tgt").Find(bson.M{"_id": tgt.Value, "client_ip": util.GetRemoteAddr(r.RemoteAddr)}).One(&tkt)

		// TGT is valid
		if tgt.Value == tkt.Ticket && time.Now().Before(tkt.Validity) {
			if svc != "" {
				st := ticket.NewServiceTicket(tkt.Ticket, svc, true)
				st.Serve(w, r)
				return
			} else {
				lt.Serve(w, template, util.LoginRequestorData{
					Config:  config.Get(),
					Session: util.LoginRequestorSession{Service: svc, Username: tkt.Username}})
				return
			}
		}
	}

	if gateway == "true" {
		w.WriteHeader(http.StatusForbidden)
		lt.Serve(w, template, util.LoginRequestorData{
			Config:  config.Get(),
			Message: util.LoginRequestorMessage{Type: "danger", Message: "This service requires a pre-established SSO session."}})
		return
	}

	lt.Serve(w, template, util.LoginRequestorData{
		Config:   config.Get(),
		Session:  util.LoginRequestorSession{Service: svc},
		ShowForm: true})
}

func loginAcceptorHandler(w http.ResponseWriter, r *http.Request) {
	svc := r.FormValue("service")
	lt := r.FormValue("lt")
	u, p := r.FormValue("username"), r.FormValue("password")

	var tkt ticket.LoginTicket
	util.GetPersistence("lt").Find(bson.M{"_id": lt}).One(&tkt)
	util.GetPersistence("lt").Remove(bson.M{"_id": tkt.Ticket})

	if lt == "" || tkt.Ticket != lt {
		forbidden(w, svc, "Form submission token was incorrect.")
		return
	}
	if tkt.Validity.Before(time.Now()) {
		forbidden(w, svc, "Form submission token has expired.")
		return
	}
	if svc != tkt.Service {
		forbidden(w, svc, "Form submission token reused in another context.")
		return
	}

	auth, redirect := authenticator.AvailableAuthenticators[config.Get().Authenticator].Auth(u, p)
	if !auth && redirect != "" {
		w.Header().Add("Location", redirect)
		w.WriteHeader(http.StatusFound)
		return
	}
	if !auth {
		forbidden(w, svc, "The credential you provided were incorrect.")
		return
	}
	tgt := ticket.NewTicketGrantingTicket(u, util.GetRemoteAddr(r.RemoteAddr))
	util.GetPersistence("tgt").Insert(tgt)

	http.SetCookie(w, &http.Cookie{Name: "CASTGC", Value: tgt.Ticket})

	if svc != "" {
		st := ticket.NewServiceTicket(tkt.Ticket, svc, false)
		st.Serve(w, r)
		return
	}

	nlt := ticket.NewLoginTicket(svc)
	nlt.Serve(w, template, util.LoginRequestorData{
		Config:  config.Get(),
		Session: util.LoginRequestorSession{Service: svc, Username: tgt.Username}})
}
