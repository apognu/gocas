package ticket

import (
	"math/rand"
	"time"

	"github.com/apognu/gocas/config"
)

type TicketGrantingTicket struct {
	Ticket   string    `bson:"_id"`
	Username string    `bson:"username"`
	ClientIP string    `bson:"client_ip"`
	Validity time.Time `bson:"validity"`
}

func NewTicketGrantingTicket(u string, ip string) TicketGrantingTicket {
	var TicketRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	tgt := make([]rune, 32)
	for i := range tgt {
		tgt[i] = TicketRunes[rand.Intn(len(TicketRunes))]
	}

	t := time.Unix(time.Now().Unix()+int64(config.Get().TicketValidity.TicketGrantingTicket), 0)
	return TicketGrantingTicket{
		Ticket:   "TGT-" + string(tgt),
		Username: u,
		ClientIP: ip,
		Validity: t,
	}
}
