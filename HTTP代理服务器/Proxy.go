package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
)

func main() {
	// 创建一个HTTP代理服务器
	proxy := &http.Server{
		Addr:    ":8080", // 设置代理服务器的监听地址和端口
		Handler: http.HandlerFunc(handleRequest), // 设置请求处理函数
	}

	// 启动代理服务器
	go func() {
		log.Fatal(proxy.ListenAndServe())
	}()

	fmt.Println("Proxy server is running on port 8080")

	// 阻止主函数退出
	select {}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	// 解析请求
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Println(err)
		return
	}

	// 打印请求信息
	fmt.Println(string(requestDump))

	// 创建一个新的请求对象，并复制原始请求的其他信息
	clientReq := &http.Request{
		Method:        r.Method,
		URL:           r.URL,
		Proto:         r.Proto,
		ProtoMajor:    r.ProtoMajor,
		ProtoMinor:    r.ProtoMinor,
		Header:        r.Header,
		Body:          r.Body,
		ContentLength: r.ContentLength,
		Host:          r.Host,
	}

	// 创建一个新的HTTP客户端
	client := &http.Client{}

	// 发送请求
	resp, err := client.Do(clientReq)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	// 解析响应
	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		log.Println(err)
		return
	}

	// 打印响应信息
	fmt.Println(string(responseDump))

	// 将响应写回到客户端
	resp.Write(w)
}
