package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/ticket"
	"github.com/apognu/gocas/util"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
)

const template = "template/oauth_login.tmpl"

var oauthConfig *oauth2.Config

func New(r *mux.Router) {
	oauthConfig = &oauth2.Config{
		ClientID:     config.Get().Oauth.ClientID,
		ClientSecret: config.Get().Oauth.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.Get().Oauth.AuthURL,
			TokenURL: config.Get().Oauth.TokenURL,
		},
		RedirectURL: config.Get().Oauth.RedirectURL,
		Scopes:      config.Get().Oauth.Scopes,
	}

	r.HandleFunc("/login", loginRequestorHandler).Methods("GET")
	r.HandleFunc("/callback", loginCallbackHandler).Methods("GET")
}

func loginRequestorHandler(w http.ResponseWriter, r *http.Request) {
	svc := r.FormValue("service")
	tgt, err := r.Cookie("CASTGC")
	if err == nil {
		var tkt ticket.TicketGrantingTicket
		util.GetPersistence("tgt").Find(bson.M{"_id": tgt.Value, "client_ip": util.GetRemoteAddr(r.RemoteAddr)}).One(&tkt)

		// TGT is valid
		if tgt.Value == tkt.Ticket && time.Now().Before(tkt.Validity) {
			if svc != "" {
				st := ticket.NewServiceTicket(tkt.Ticket, svc, true)
				st.Serve(w, r)
				return
			} else {
				lt := ticket.NewLoginTicket(svc)
				lt.Serve(w, template, util.LoginRequestorData{
					Config:  config.Get(),
					Session: util.LoginRequestorSession{Service: svc, Username: tkt.Username}})
				return
			}
		}
	}

	lt := ticket.NewLoginTicket(svc)
	url := oauthConfig.AuthCodeURL(lt.Ticket)
	lt.Serve(w, template, util.LoginRequestorData{
		Config:  config.Get(),
		Session: util.LoginRequestorSession{Url: url}})
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code, lt := r.FormValue("code"), r.FormValue("state")

	var tkt ticket.LoginTicket
	util.GetPersistence("lt").Find(bson.M{"_id": lt}).One(&tkt)
	util.GetPersistence("lt").Remove(bson.M{"_id": tkt.Ticket})
	if lt == "" || tkt.Ticket != lt {
		lt := ticket.NewLoginTicket(tkt.Service)
		lt.Serve(w, template, util.LoginRequestorData{
			Config:  config.Get(),
			Message: util.LoginRequestorMessage{Type: "danger", Message: "Form submission token was incorrect."}})
		return
	}

	c := context.TODO()
	token, err := oauthConfig.Exchange(c, code)
	if err != nil {
		w.Header().Add("Location", config.Get().UrlPrefix)
		w.WriteHeader(http.StatusFound)
		return
	}

	cl := oauthConfig.Client(c, token)
	resp, err := cl.Get(config.Get().Oauth.UserinfoURL)
	if err != nil {
		w.Header().Add("Location", config.Get().UrlPrefix)
		w.WriteHeader(http.StatusFound)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Header().Add("Location", config.Get().UrlPrefix)
		w.WriteHeader(http.StatusFound)
		return
	}
	var info map[string]interface{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		fmt.Println(err)
		w.Header().Add("Location", config.Get().UrlPrefix)
		w.WriteHeader(http.StatusFound)
		return
	}

	tgt := ticket.NewTicketGrantingTicket(info["name"].(string), util.GetRemoteAddr(r.RemoteAddr))
	util.GetPersistence("tgt").Insert(tgt)
	http.SetCookie(w, &http.Cookie{Name: "CASTGC", Value: tgt.Ticket, Path: "/"})

	if tkt.Service != "" {
		st := ticket.NewServiceTicket(tkt.Ticket, tkt.Service, false)
		st.Serve(w, r)
		return
	}

	w.Header().Add("Location", config.Get().UrlPrefix)
	w.WriteHeader(http.StatusFound)
}
