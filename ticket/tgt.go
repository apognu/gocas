package ticket

import (
	"time"

	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/util"
)

type TicketGrantingTicket struct {
	Ticket   string    `bson:"_id"`
	Username string    `bson:"username"`
	ClientIP string    `bson:"client_ip"`
	Validity time.Time `bson:"validity"`
}

func NewTicketGrantingTicket(u string, ip string) TicketGrantingTicket {
	tgt := generateTicket("TGT", 32)
	t := time.Unix(time.Now().Unix()+int64(config.Get().TicketValidity.TicketGrantingTicket), 0)
	tkt := TicketGrantingTicket{
		Ticket:   tgt,
		Username: u,
		ClientIP: ip,
		Validity: t,
	}
	util.GetPersistence("tgt").Insert(tkt)
	return tkt
}
