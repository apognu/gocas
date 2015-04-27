package oauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"text/template"
	"time"

	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/ticket"
	"github.com/apognu/gocas/util"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
)

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

func showLoginForm(w http.ResponseWriter, data util.LoginRequestorData) {
	lt := ticket.NewLoginTicket(data.Session.Service)
	data.Session.Ticket = lt.Ticket
	util.GetPersistence("lt").Insert(lt)

	t, err := template.ParseFiles("template/oauth_login.tmpl")
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
	w.Header().Add("Location", url)
	w.WriteHeader(http.StatusFound)
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

	lt := ticket.NewLoginTicket(svc)
	util.GetPersistence("lt").Insert(lt)
	url := oauthConfig.AuthCodeURL(lt.Ticket)
	t, err := template.ParseFiles("template/oauth_login.tmpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.Execute(w, util.LoginRequestorData{Config: config.Get(), Session: util.LoginRequestorSession{Url: url}})
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code, lt := r.FormValue("code"), r.FormValue("state")

	var tkt ticket.LoginTicket
	util.GetPersistence("lt").Find(bson.M{"_id": lt}).One(&tkt)
	util.GetPersistence("lt").Remove(bson.M{"_id": tkt.Ticket})
	if lt == "" || tkt.Ticket != lt {
		showLoginForm(w, util.LoginRequestorData{
			Config:  config.Get(),
			Message: util.LoginRequestorMessage{Type: "danger", Message: "Form submission token was incorrect."}})
		return
	}

	c := context.TODO()
	token, err := oauthConfig.Exchange(c, code)
	if err != nil {
		w.Header().Add("Location", config.Get().UrlPrefix)
		w.WriteHeader(http.StatusFound)
	}

	cl := oauthConfig.Client(c, token)
	resp, err := cl.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		w.Header().Add("Location", config.Get().UrlPrefix)
		w.WriteHeader(http.StatusFound)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Header().Add("Location", config.Get().UrlPrefix)
		w.WriteHeader(http.StatusFound)
	}
	var info map[string]string
	err = json.Unmarshal(body, &info)
	if err != nil {
		w.Header().Add("Location", config.Get().UrlPrefix)
		w.WriteHeader(http.StatusFound)
	}

	tgt := ticket.NewTicketGrantingTicket(info["name"], util.GetRemoteAddr(r.RemoteAddr))
	util.GetPersistence("tgt").Insert(tgt)
	http.SetCookie(w, &http.Cookie{Name: "CASTGC", Value: tgt.Ticket, Path: "/"})

	if tkt.Service != "" {
		serveServiceTicket(w, r, tgt.Ticket, tkt.Service, false)
		return
	}

	w.Header().Add("Location", config.Get().UrlPrefix)
	w.WriteHeader(http.StatusFound)
}
