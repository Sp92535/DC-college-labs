package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/{n}", func(w http.ResponseWriter, r *http.Request) {

		path := strings.TrimPrefix(r.URL.Path, "/")
		n, err := strconv.Atoi(path)
		if err != nil {
			http.Error(w, "invalid number", http.StatusBadRequest)
			return
		}

		w.Write(fmt.Appendf(nil, "response from -> %s\n", getLocalIp()))
		w.Write(fmt.Append(nil, fib(n)))

	})

	log.Printf("Listening on port 6969...")
	err := http.ListenAndServe(":6969", mux)
	if err != nil {
		log.Fatalf("error listneing %v", err)
	}

}

func getLocalIp() string {
	// dial google dns
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatalf("error getting local ip: %v", err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP.String()
}

func fib(n int) int {
	if n <= 1 {
		return n
	}

	return fib(n-1) + fib(n-2)
}
