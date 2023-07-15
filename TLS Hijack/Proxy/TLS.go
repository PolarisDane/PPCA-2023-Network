package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	ProxyServer := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: http.HandlerFunc(HandleRequest),
	}
	go func() {
		log.Println(ProxyServer.ListenAndServe())
	}()

	fmt.Println("Proxy server is running on port 8080, TLS hijack is enabled")

	select {}
}

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodConnect {
		fmt.Println("CONNECT METHOD REQUIRED")
		handleTunneling(w, r)
	} else {
		fmt.Println("HTTP METHOD REQUIRED")
		handleHTTP(w, r)
	}
}

func handleTunneling(w http.ResponseWriter, r *http.Request) {
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking not supported", http.StatusInternalServerError)
		return
	}

	clientConn, _, err := hijacker.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	cert, err := tls.LoadX509KeyPair("localhost.crt", "localhost.key")
	if err != nil {
		log.Fatalf("proxy: loadkeys: %s", err)
	}

	clientConn.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	tlsClientConn := tls.Server(clientConn, tlsConfig)

	destConn, err := tls.Dial("tcp", r.Host, &tls.Config{
		InsecureSkipVerify: true, // you might want to verify the server certificate in a real-world situation
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	connect(destConn, tlsClientConn)
}

func handleHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()
	copyHeader(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func Forward(tarConn, srcConn net.Conn) {
	buf := make([]byte, 32*1024)
	for {
		rcount, err := srcConn.Read(buf[:32*1024])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		if rcount > 0 {
			wcount, err := tarConn.Write(buf[:rcount])
			fmt.Print(string(buf[:rcount]))
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			if wcount < rcount {
				fmt.Println("End of forward")
				return
			}
		}
	}
}

func connect(destConn, clientConn net.Conn) {
	defer destConn.Close()
	defer clientConn.Close()
	go Forward(clientConn, destConn)
	io.Copy(destConn, clientConn)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
