package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"time"
)

func handshakeTLS(tar string) {
	cert, err := tls.LoadX509KeyPair("server.crt", "server.key")
	if err != nil {
		log.Fatalf("proxy: loadkeys: %s", err)
	}
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}
	tlsListener, err := tls.Listen("tcp", "127.0.0.1:9000", tlsConfig)
	tlsConn, err := tlsListener.Accept()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	tlsConfig = &tls.Config{
		InsecureSkipVerify: true,
	}
	serverConn, err := tls.Dial("tcp", tar, tlsConfig)
	go io.Copy(tlsConn, serverConn)
	io.Copy(serverConn, tlsConn)
}

func generateCert(host string) (tls.Certificate, error) {

	host, _, _ = net.SplitHostPort(host)

	// 读取根证书
	rootCertPEM, err := ioutil.ReadFile("ca_root.crt") // 改变文件名
	if err != nil {
		log.Fatalf("Failed to read root certificate: %v", err)
	}
	block, _ := pem.Decode(rootCertPEM)
	if block == nil {
		log.Fatalf("Failed to decode PEM block containing the certificate")
	}
	rootCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse certificate: %v", err)
	}

	// 读取私钥
	rootKeyPEM, err := ioutil.ReadFile("ca_private.key") // 改变文件名
	if err != nil {
		log.Fatalf("Failed to read private key: %v", err)
	}
	block, _ = pem.Decode(rootKeyPEM)
	if block == nil {
		log.Fatalf("Failed to decode PEM block containing the private key")
	}
	rootKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Now you can use rootCert and rootKey to sign a new certific
	// 这里假设你已经有了根证书和根证书的私钥
	// rootCert 是 *x509.Certificate 对象
	// rootKey 是 crypto.PrivateKey 对象

	// 生成新的私钥
	priv, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	// 创建证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(2), // 为了确保每个证书的序列号是唯一的，你可能需要动态生成这个值
		Subject: pkix.Name{
			CommonName:   host, // 这里设置你要签发的主机名
			Organization: []string{"My Company"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180), // 180 days
		KeyUsage: x509.KeyUsageKeyEncipherment |
			x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{host}, // 这里设置你要签发的主机名
	}

	// 使用根证书签发新的证书
	derBytes, _ := x509.CreateCertificate(rand.Reader, &template, rootCert, &priv.PublicKey, rootKey)

	// 将新的证书和私钥编码为 PEM 格式
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyBytes, _ := x509.MarshalECPrivateKey(priv)
	keyPem := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyBytes})

	// Now you can use certPem and keyPem for tls.X509KeyPair()

	cert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		log.Fatalf("Creating certificate: %s", err)
	}

	return cert, nil
}
