package util

import (
	"time"

	"github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var (
	db *mgo.Session
)

func GetPersistence(c string) *mgo.Collection {
	if db == nil {
		var err error
		db, err = mgo.DialWithTimeout(GetConfig().Mongo.Host, 2*time.Second)
		if err != nil {
			logrus.Fatalf("error connecting to MongoDB: %s", err)
		}
	}

	return db.DB("gocas").C(c)
}
