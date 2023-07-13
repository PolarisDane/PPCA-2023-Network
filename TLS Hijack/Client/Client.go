package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
)

func main() {
	caCert, err := ioutil.ReadFile("localhost.crt")
	if err != nil {
		log.Fatalf("client: read ca cert: %s", err)
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	conf := &tls.Config{
		RootCAs: caCertPool,
	}

	conn, err := tls.Dial("tcp", "localhost:8080", conf)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	n, err := conn.Write([]byte("CONNECT www.baidu.com:443 HTTP/1.1\nHost: www.baidu.com:443\nProxy-Connection: Keep-Alive\nUser-Agent: curl/7.81.0"))
	if err != nil {
		log.Println(n, err)
		return
	}

	buf := make([]byte, 512)
	n, err = conn.Read(buf)
	if err != nil {
		log.Println(n, err)
		return
	}

	fmt.Println(string(buf[:n]))
}
