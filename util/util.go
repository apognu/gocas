package util

import (
	"fmt"
	"strings"
)

func Url(path string) string {
	return fmt.Sprintf("%s%s", GetConfig().UrlPrefix, path)
}

func GetRemoteAddr(raddr string) string {
	i := strings.LastIndex(raddr, ":")
	if i == -1 {
		return raddr
	}
	return raddr[:i]
}
