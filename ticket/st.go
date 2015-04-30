package ticket

import (
	"fmt"
	"net/http"
	"net/url"
	"text/template"
	"time"

	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/util"
	"gopkg.in/mgo.v2/bson"
)

type ServiceTicket struct {
	Ticket   string    `bson:"_id"`
	Service  string    `bson:"service"`
	Tgt      string    `bson:"tgt"`
	Validity time.Time `bson:"validity"`
	FromSso  bool      `bson:"from_sso"`
}

func NewServiceTicket(tgt string, svc string, sso bool) ServiceTicket {
	st := generateTicket("ST", 32)
	t := time.Unix(time.Now().Unix()+int64(config.Get().TicketValidity.ServiceTicket), 0)
	tkt := ServiceTicket{
		Service:  svc,
		Tgt:      tgt,
		Ticket:   st,
		Validity: t,
		FromSso:  sso,
	}
	util.GetPersistence("st").Insert(tkt)
	return tkt
}

func (st ServiceTicket) Validate() bool {
	if st.Service != "" && len(config.Get().Services) > 0 {
		matched := false
		u, err := url.Parse(st.Service)
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

func (st ServiceTicket) Serve(w http.ResponseWriter, r *http.Request) {
	if !st.Validate() {
		lt := NewEmptyLoginTicket()
		w.WriteHeader(http.StatusForbidden)
		lt.Serve(w, "template/login.tmpl", util.LoginRequestorData{
			Config:  config.Get(),
			Message: util.LoginRequestorMessage{Type: "danger", Message: "The service that asked for authentication is not authorized to do so."}})
		return
	}

	url := fmt.Sprintf("%s?ticket=%s", st.Service, st.Ticket)
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
	t.Execute(w, util.LoginRequestorData{Config: config.Get(), Session: util.LoginRequestorSession{Service: st.Service, Username: tkt.Username, Url: url}, ShowForm: false})
}

func (st ServiceTicket) GetTicketGrantingTicket() TicketGrantingTicket {
	var tgt TicketGrantingTicket
	util.GetPersistence("tgt").Find(bson.M{"_id": st.Tgt}).One(&tgt)

	return tgt
}
