package main

import (
	"net/http"
	"strings"

	"github.com/apognu/gocas/ticket"
	"github.com/apognu/gocas/util"
	"gopkg.in/mgo.v2/bson"
)

func proxyHandler(w http.ResponseWriter, r *http.Request) {
	pgt := r.FormValue("ticket")
	svc := r.FormValue("targetService")

	if pgt == "" || svc == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write(util.NewCASProxyFailureResponse("INVALID_REQUEST", "'pgt' and 'targetService' parameters are both required"))
	}

	var tkt ticket.ProxyGrantingTicket
	util.GetPersistence("pgt").Find(bson.M{"_id": pgt}).One(&tkt)
	if tkt.Ticket != pgt {
		w.WriteHeader(http.StatusForbidden)
		w.Write(util.NewCASProxyFailureResponse("UNAUTHORIZED_SERVICE", "unknown PGT"))
		return
	}
	if tkt.Service != svc {
		w.WriteHeader(http.StatusForbidden)
		w.Write(util.NewCASProxyFailureResponse("UNAUTHORIZED_SERVICE", "the given service is not authorized to request Proxy Ticket from this PGT"))
		return
	}

	pt := ticket.NewProxyTicket(pgt)
	w.Write(util.NewCASProxySuccessResponse(pt.Ticket))
}

func proxyValidateHandler(w http.ResponseWriter, r *http.Request) {
	svc, pt := r.FormValue("service"), r.FormValue("ticket")

	// /proxyValidate must implement /serviceValidate if given a Service Ticket
	if strings.HasPrefix(pt, "ST-") {
		serviceValidateHandler(w, r)
		return
	}

	var tkt ticket.ProxyTicket
	util.GetPersistence("pgt").Find(bson.M{"_id": pt}).One(&tkt)
	util.GetPersistence("pgt").Remove(bson.M{"_id": pt})
	if tkt.Ticket != pt {
		w.Write(util.NewCASFailureResponse("INVALID_TICKET", "Ticket not recognized"))
		return
	}

	pgt := tkt.GetProxyGrantingTicket()
	if pgt.Service != svc {
		w.Write(util.NewCASFailureResponse("INVALID_SERVICE", "Ticket was used for another service than it was generated for"))
		return
	}

	w.Write(util.NewCASSuccessResponse(pgt.Username, ""))
}
