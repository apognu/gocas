package ticket

import (
	"time"

	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/util"
	"gopkg.in/mgo.v2/bson"
)

type ProxyGrantingTicketIOU struct {
	Ticket string `bson:"_id"`
}

func NewProxyGrantingTicketIOU() ProxyGrantingTicketIOU {
	pgtiou := generateTicket("PGTIOU", 32)
	tkt := ProxyGrantingTicketIOU{
		Ticket: pgtiou,
	}
	return tkt
}

type ProxyGrantingTicket struct {
	Ticket   string    `bson:"_id"`
	Service  string    `bson:"service"`
	Username string    `bson:"username"`
	ClientIP string    `bson:"client_ip"`
	Validity time.Time `bson:"validity"`
}

func NewProxyGrantingTicket(svc string, pgtiou string, u string, ip string) ProxyGrantingTicket {
	pgt := generateTicket("PGT", 32)
	// TODO: Config setting for PGT expiration
	t := time.Unix(time.Now().Unix()+int64(config.Get().TicketValidity.ProxyGrantingTicket), 0)
	tkt := ProxyGrantingTicket{
		Ticket:   pgt,
		Service:  svc,
		Username: u,
		ClientIP: ip,
		Validity: t,
	}
	util.GetPersistence("pgt").Insert(tkt)
	return tkt
}

type ProxyTicket struct {
	Ticket   string    `bson:"_id"`
	Pgt      string    `bson:"pgt"`
	Validity time.Time `bson:"validity"`
}

func NewProxyTicket(pgt string) ProxyTicket {
	pt := generateTicket("PT", 32)
	// TODO: Config setting for PT expiration
	t := time.Unix(time.Now().Unix()+int64(config.Get().TicketValidity.ProxyTicket), 0)
	tkt := ProxyTicket{
		Ticket:   pt,
		Pgt:      pgt,
		Validity: t,
	}
	util.GetPersistence("pgt").Insert(tkt)
	return tkt
}

func (pt ProxyTicket) GetProxyGrantingTicket() ProxyGrantingTicket {
	var pgt ProxyGrantingTicket
	util.GetPersistence("pgt").Find(bson.M{"_id": pt.Pgt}).One(&pgt)

	return pgt
}
