package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	//环境变量
	// os.Clearenv()
	os.Setenv("version", "1.0.0")
	version := os.Getenv("version")
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		for key, v := range r.Header {
			value := strings.Join(v, " ")
			rw.Header().Add(key, value)
		}
		rw.Header().Add("version", version)
		rw.WriteHeader(200)
		// fmt.Fprintf(rw, "success")
		ip, err := getIP(r)
		if err != nil {
			fmt.Printf("getIperr:%s", err.Error())
			return
		}
		fmt.Printf("clientIp:%s HttpStatusCode %d \n", ip, 200)
	})
	http.HandleFunc("/healthz", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "ok")
		rw.WriteHeader(200)
	})
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err.Error())
	}
}

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("No valid ip found")
}
