package util

import "strings"

func GetRemoteAddr(raddr string) string {
	i := strings.LastIndex(raddr, ":")
	if i == -1 {
		return raddr
	}
	return raddr[:i]
}
