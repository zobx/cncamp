package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	//环境变量
	// os.Clearenv()
	// os.Setenv("version", "1.0.0")
	version := os.Getenv("VERSION")
	r := http.NewServeMux()
	r.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		for key, v := range r.Header {
			value := strings.Join(v, " ")
			rw.Header().Add(key, value)
		}
		rw.Header().Add("version", version)
		rw.WriteHeader(200)
		fmt.Fprintf(rw, "success"+version)
		ip, err := getIP(r)
		if err != nil {
			fmt.Printf("getIperr:%s", err.Error())
			return
		}
		fmt.Printf("clientIp:%s HttpStatusCode %d \n", ip, 200)
	})
	r.HandleFunc("/healthz", func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprint(rw, "ok")
		rw.WriteHeader(200)
	})
	server := &http.Server{
		Addr:    ":80",
		Handler: r,
	}
	go server.ListenAndServe()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigs:
		server.Shutdown(context.Background())
		log.Println("http shutdown")
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
