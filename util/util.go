package util

import (
	"crypto/md5"
	"fmt"
	"strings"
	"time"

	"github.com/apognu/gocas/config"
	"gopkg.in/mgo.v2/bson"
)

func Url(path string) string {
	return fmt.Sprintf("%s%s", config.Get().UrlPrefix, path)
}

func ResolveTemplate(tmpl string) string {
	return fmt.Sprintf("%s/%s.tmpl", config.Get().TemplatePath, tmpl)
}

func GetRemoteAddr(raddr string) string {
	i := strings.LastIndex(raddr, ":")
	if i == -1 {
		return raddr
	}
	return raddr[:i]
}

func IncrementFailedLogin(ip string, u string) {
	t := time.Now()
	flip := FailedLogin{
		Id:        fmt.Sprintf("%x", md5.Sum([]byte(GetRemoteAddr(ip)))),
		Ip:        GetRemoteAddr(ip),
		UpdatedAt: t,
		Count:     1,
	}
	flu := FailedLogin{
		Id:        fmt.Sprintf("%x", md5.Sum([]byte(u))),
		Username:  u,
		UpdatedAt: t,
		Count:     1,
	}

	err := GetPersistence("failed").Update(bson.M{"_id": flip.Id}, bson.M{"$inc": bson.M{"count": 1}, "$set": bson.M{"updated_at": t}})
	if err != nil {
		GetPersistence("failed").Insert(flip)
	}

	if u != "" {
		err = GetPersistence("failed").Update(bson.M{"_id": flu.Id}, bson.M{"$inc": bson.M{"count": 1}, "$set": bson.M{"updated_at": t}})
		if err != nil {
			GetPersistence("failed").Insert(flu)
		}
	}
}
