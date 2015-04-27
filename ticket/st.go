package ticket

import (
	"math/rand"
	"time"

	"github.com/apognu/gocas/util"
	"gopkg.in/mgo.v2/bson"
)

type ServiceTicket struct {
	Ticket   string    `bson:"_id"`
	Service  string    `bson:"service"`
	Tgt      string    `bson:"tgt"`
	Validity time.Time `bson:"validity"`
}

func NewServiceTicket(tgt string, svc string) ServiceTicket {
	var TicketRunes = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	st := make([]rune, 32)
	for i := range st {
		st[i] = TicketRunes[rand.Intn(len(TicketRunes))]
	}

	t := time.Unix(time.Now().Unix()+int64(util.GetConfig().TicketValidity.ServiceTicket), 0)
	return ServiceTicket{
		Service:  svc,
		Tgt:      tgt,
		Ticket:   "ST-" + string(st),
		Validity: t,
	}
}

func (st ServiceTicket) GetTicketGrantingTicket() TicketGrantingTicket {
	var tgt TicketGrantingTicket
	util.GetPersistence("tgt").Find(bson.M{"_id": st.Tgt}).One(&tgt)

	return tgt
}
