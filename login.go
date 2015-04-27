package main

import (
	"fmt"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/ticket"
	"github.com/apognu/gocas/util"
	"gopkg.in/mgo.v2/bson"
)

func showLoginForm(w http.ResponseWriter, data util.LoginRequestorData) {
	lt := ticket.NewLoginTicket(data.Session.Service)
	data.Session.Ticket = lt.Ticket
	util.GetPersistence("lt").Insert(lt)

	t, err := template.ParseFiles("template/login.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}

func serveServiceTicket(w http.ResponseWriter, r *http.Request, tgt string, svc string, sso bool) {
	st := ticket.NewServiceTicket(tgt, svc, sso)
	util.GetPersistence("st").Insert(st)

	url := fmt.Sprintf("%s?ticket=%s", svc, st.Ticket)
	if r.FormValue("warn") != "true" {
		w.Header().Add("Location", url)
		w.WriteHeader(http.StatusFound)
		return
	}

	tkt := st.GetTicketGrantingTicket()
	t, err := template.ParseFiles("template/warn.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.Execute(w, util.LoginRequestorData{Config: config.Get(), Session: util.LoginRequestorSession{Service: svc, Username: tkt.Username}})
}

func isServiceWhitelisted(svc string) bool {
	if svc != "" && len(config.Get().Services) > 0 {
		matched := false
		u, err := url.Parse(svc)
		if err == nil {
			for _, s := range config.Get().Services {
				if s == u.Host {
					matched = true
				}
			}
		}
		return matched
	}
	return true
}

func loginRequestorHandler(w http.ResponseWriter, r *http.Request) {
	tgt, err := r.Cookie("CASTGC")
	svc := r.FormValue("service")
	renew, gateway := r.FormValue("renew"), r.FormValue("gateway")

	if !isServiceWhitelisted(svc) {
		showLoginForm(w, util.LoginRequestorData{
			Config:   config.Get(),
			Session:  util.LoginRequestorSession{Service: svc},
			Message:  util.LoginRequestorMessage{Type: "danger", Message: fmt.Sprintf("Service <b>%s</b> is not allowed to use the SSO.", svc)},
			ShowForm: false})
		return
	}

	// The client sent us a TGT, do not display login form
	if err == nil && renew != "true" {
		var tkt ticket.TicketGrantingTicket
		util.GetPersistence("tgt").Find(bson.M{"_id": tgt.Value, "client_ip": util.GetRemoteAddr(r.RemoteAddr)}).One(&tkt)

		// TGT is valid
		if tgt.Value == tkt.Ticket && time.Now().Before(tkt.Validity) {
			if svc != "" {
				serveServiceTicket(w, r, tkt.Ticket, svc, true)
				return
			} else {
				showLoginForm(w, util.LoginRequestorData{
					Config:  config.Get(),
					Session: util.LoginRequestorSession{Service: svc, Username: tkt.Username}})
				return
			}
		}
	}

	if gateway == "true" {
		showLoginForm(w, util.LoginRequestorData{
			Config:  config.Get(),
			Message: util.LoginRequestorMessage{Type: "danger", Message: "This service requires a pre-established SSO session."}})
		return
	}

	showLoginForm(w, util.LoginRequestorData{
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
		showLoginForm(w, util.LoginRequestorData{
			Config:   config.Get(),
			Session:  util.LoginRequestorSession{Service: svc},
			Message:  util.LoginRequestorMessage{Type: "danger", Message: "Form submission token was incorrect."},
			ShowForm: true})
		return
	}
	if tkt.Validity.Before(time.Now()) {
		showLoginForm(w, util.LoginRequestorData{
			Config:   config.Get(),
			Session:  util.LoginRequestorSession{Service: svc},
			Message:  util.LoginRequestorMessage{Type: "danger", Message: "Form submission token has expired."},
			ShowForm: true})
		return
	}
	if svc != tkt.Service {
		showLoginForm(w, util.LoginRequestorData{
			Config:   config.Get(),
			Session:  util.LoginRequestorSession{Service: svc},
			Message:  util.LoginRequestorMessage{Type: "danger", Message: "Form submission token reused in another context."},
			ShowForm: true})
		return
	}

	auth := util.AvailableAuthenticators[config.Get().Authenticator]
	if !auth.Auth(u, p) {
		showLoginForm(w, util.LoginRequestorData{
			Config:   config.Get(),
			Session:  util.LoginRequestorSession{Service: svc},
			Message:  util.LoginRequestorMessage{Type: "danger", Message: "The credential you provided were incorrect."},
			ShowForm: true})
		return
	}
	tgt := ticket.NewTicketGrantingTicket(u, util.GetRemoteAddr(r.RemoteAddr))
	util.GetPersistence("tgt").Insert(tgt)

	http.SetCookie(w, &http.Cookie{Name: "CASTGC", Value: tgt.Ticket})

	if svc != "" {
		serveServiceTicket(w, r, tgt.Ticket, svc, false)
		return
	}

	showLoginForm(w, util.LoginRequestorData{
		Config:  config.Get(),
		Session: util.LoginRequestorSession{Service: svc, Username: tgt.Username}})
}
