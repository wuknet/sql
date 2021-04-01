package fun

import (
	//"fmt"
	"net"
	"net/http"
	"strings"
)

func GetIP(r *http.Request) string {
	// 尝试从 X-Forwarded-For 中获取
	//fmt.Println(r.RemoteAddr)
	//http://pv.sohu.com/cityjson?ie=utf-8
	ip := strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	err := net.ParseIP(ip)
	if err != nil {
		return ip
	}
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip = strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		err := net.ParseIP(ip)
		if err != nil {
			return ip
		}
	}
	return ""
}
