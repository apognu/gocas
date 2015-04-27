package main

import (
	"net/http"

	"github.com/apognu/gocas/util"
	"gopkg.in/mgo.v2/bson"
)

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	tgt, err := r.Cookie("CASTGC")
	if err == nil {
		util.GetPersistence("tgt").Remove(bson.M{"_id": tgt.Value})
		http.SetCookie(w, &http.Cookie{
			Name:   "CASTGC",
			Value:  "",
			MaxAge: -1,
		})
	}

	w.Header().Add("Location", util.Url("/login"))
	w.WriteHeader(http.StatusFound)
}
