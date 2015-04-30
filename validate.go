package main

import (
	"fmt"
	"net/http"

	"github.com/apognu/gocas/ticket"
	"github.com/apognu/gocas/util"
	"gopkg.in/mgo.v2/bson"
)

func validateHandler(w http.ResponseWriter, r *http.Request) {
	svc, st := r.FormValue("service"), r.FormValue("ticket")
	renew := r.FormValue("renew")

	var tkt ticket.ServiceTicket
	util.GetPersistence("st").Find(bson.M{"_id": st, "service": svc}).One(&tkt)
	util.GetPersistence("st").Remove(bson.M{"_id": st})
	if tkt.Service != svc || tkt.Ticket != st {
		w.Write([]byte("no\n"))
		return
	}
	if renew == "true" && tkt.FromSso {
		w.Write([]byte("no\n"))
		return
	}

	w.Write([]byte(fmt.Sprintf("yes\n%s\n", tkt.GetTicketGrantingTicket().Username)))
}

func serviceValidateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/xml")
	svc, st := r.FormValue("service"), r.FormValue("ticket")
	renew := r.FormValue("renew")

	if svc == "" || st == "" {
		w.Write(util.NewCASFailureResponse("INVALID_REQUEST", "Both ticket and service parameters must be given"))
		return
	}

	var tkt ticket.ServiceTicket
	util.GetPersistence("st").Find(bson.M{"_id": st}).One(&tkt)
	util.GetPersistence("st").Remove(bson.M{"_id": st})
	if tkt.Ticket != st {
		w.Write(util.NewCASFailureResponse("INVALID_TICKET", "Ticket not recognized"))
		return
	}
	if tkt.Service != svc {
		w.Write(util.NewCASFailureResponse("INVALID_SERVICE", "Ticket was used for another service than it was generated for"))
		return
	}
	if renew == "true" && tkt.FromSso {
		w.Write(util.NewCASFailureResponse("INVALID_TICKET_SPEC", "The service requires a service ticket obtained through principal validation"))
		return
	}

	pgtUrl := r.FormValue("pgtUrl")
	if pgtUrl != "" {
		pgtiou := ticket.NewProxyGrantingTicketIOU()
		pgt := ticket.NewProxyGrantingTicket(svc, pgtiou.Ticket, tkt.GetTicketGrantingTicket().Username, util.GetRemoteAddr(r.RemoteAddr))

		// TODO: Only accept HTTPS URLs
		resp, err := http.Get(fmt.Sprintf("%s?pgtIou=%s&pgtId=%s", pgtUrl, pgtiou.Ticket, pgt.Ticket))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(util.NewCASFailureResponse("INTERNAL_ERROR", "Internal Server Error"))
			return
		}
		if resp.StatusCode != 200 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write(util.NewCASFailureResponse("INVALID_PROXY_CALLBACK", "Error requesting proxy callback URL"))
			return
		}

		w.Write(util.NewCASSuccessResponse(tkt.GetTicketGrantingTicket().Username, pgtiou.Ticket))
		return
	}

	w.Write(util.NewCASSuccessResponse(tkt.GetTicketGrantingTicket().Username, ""))
}
