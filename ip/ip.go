package ip

import (
	"fmt"
	"github.com/sebest/xff"
	"net"
	"net/http"
	"os"
)

//MustGetPublic returns ip from request
func MustGetPublicFromRequest(r *http.Request) string {
	ip, _, err := net.SplitHostPort(xff.GetRemoteAddr(r))

	if !xff.IsPublicIP(net.ParseIP(ip)) || err != nil {
		if os.Getenv("IP_DEBUG") == "" {
			panic("IP NOT PUBLIC")
		} else {
			fmt.Println("WARNING: Using env IP_DEBUG for getting public IP address")
		}
	}

	return ip
}
