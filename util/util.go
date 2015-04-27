package util

import (
	"fmt"
	"strings"

	"github.com/apognu/gocas/config"
)

func Url(path string) string {
	return fmt.Sprintf("%s%s", config.Get().UrlPrefix, path)
}

func GetRemoteAddr(raddr string) string {
	i := strings.LastIndex(raddr, ":")
	if i == -1 {
		return raddr
	}
	return raddr[:i]
}
