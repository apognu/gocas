package ticket

import (
	"net/http"
	"text/template"
	"time"

	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/util"
)

type LoginTicket struct {
	Ticket   string    `bson:"_id"`
	Service  string    `bson:"service"`
	Validity time.Time `bson:"validity"`
}

func NewLoginTicket(svc string) LoginTicket {
	lt := generateTicket("LT", 32)
	t := time.Unix(time.Now().Unix()+int64(config.Get().TicketValidity.LoginTicket), 0)
	tkt := LoginTicket{
		Service:  svc,
		Ticket:   lt,
		Validity: t,
	}
	util.GetPersistence("lt").Insert(tkt)
	return tkt
}

func NewEmptyLoginTicket() LoginTicket {
	return LoginTicket{}
}

func (lt LoginTicket) Serve(w http.ResponseWriter, tmpl string, data util.LoginRequestorData) {
	data.Session.Ticket = lt.Ticket

	t, err := template.ParseFiles(tmpl)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	t.Execute(w, data)
}
