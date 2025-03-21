package utils

import (
	"log"
	"net"
	"net/http"
)

func GetIP(r *http.Request) string {
	remoteAddr := r.RemoteAddr
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}
	return host
}

func GetPrivateIP() net.IP { // https://gosamples.dev/local-ip-address/
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}
