package ticket

import (
	"math/rand"
	"time"

	"github.com/apognu/gocas/config"
)

type LoginTicket struct {
	Ticket   string    `bson:"_id"`
	Service  string    `bson:"service"`
	Validity time.Time `bson:"validity"`
}

func NewLoginTicket(svc string) LoginTicket {
	var TicketRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	lt := make([]rune, 32)
	for i := range lt {
		lt[i] = TicketRunes[rand.Intn(len(TicketRunes))]
	}

	t := time.Unix(time.Now().Unix()+int64(config.Get().TicketValidity.LoginTicket), 0)
	return LoginTicket{
		Service:  svc,
		Ticket:   "LT-" + string(lt),
		Validity: t,
	}
}
