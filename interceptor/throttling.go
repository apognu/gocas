package interceptor

import (
	"net/http"
	"strings"
	"time"

	"github.com/apognu/gocas/config"
	"github.com/apognu/gocas/ticket"
	"github.com/apognu/gocas/util"
	"gopkg.in/mgo.v2/bson"
)

type ThrottlingInterceptor struct{}

func (ThrottlingInterceptor) getFailureCount(ip string, u string) (int, int) {
	ipt, _ := util.GetPersistence("failed").Find(bson.M{"$and": []bson.M{bson.M{"ip": ip}, bson.M{"count": bson.M{"$gte": config.Get().Throttling.MaxFailuresByIp}}}}).Count()
	ut, _ := util.GetPersistence("failed").Find(bson.M{"$and": []bson.M{bson.M{"username": u}, bson.M{"count": bson.M{"$gte": config.Get().Throttling.MaxFailuresByUsername}}}}).Count()

	return ipt, ut
}

func (ThrottlingInterceptor) Init() {
	go func() {
		var r []util.FailedLogin
		for {
			t := time.Now().Add(-config.Get().Throttling.DecrementInterval)
			util.GetPersistence("failed").Find(bson.M{"$and": []bson.M{bson.M{"count": bson.M{"$gt": 0}}, bson.M{"updated_at": bson.M{"$lt": t}}}}).All(&r)

			for _, item := range r {
				util.GetPersistence("failed").UpdateId(item.Id, bson.M{"$inc": bson.M{"count": -1}, "$set": bson.M{"updated_at": time.Now()}})
			}

			time.Sleep(10 * time.Second)
		}
	}()
}

func (i ThrottlingInterceptor) Intercept(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ipt, ut := i.getFailureCount(util.GetRemoteAddr(r.RemoteAddr), r.FormValue("username"))

	if strings.HasPrefix(r.RequestURI, "/static/") {
		next(w, r)
		return
	}

	if ipt > 0 || ut > 0 {
		w.WriteHeader(http.StatusForbidden)
		ticket.NewEmptyLoginTicket().Serve(w, util.ResolveTemplate("throttling"), util.LoginRequestorData{Config: config.Get()})
		return
	}

	next(w, r)
}
