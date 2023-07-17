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
		Addr:    ":8080",                         // 设置代理服务器的监听地址和端口
		Handler: http.HandlerFunc(HandleRequest), // 设置请求处理函数
	}

	// 启动代理服务器
	go func() {
		log.Fatal(proxy.ListenAndServe())
	}()

	fmt.Println("Proxy server is running on port 8080")

	// 阻止主函数退出
	select {}
}

func GetResponse(Writer http.ResponseWriter, Request *http.Request) (err error) {
	// 创建一个新的HTTP客户端
	client := &http.Client{}

	// 发送请求
	resp, err := client.Do(Request)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// 解析响应
	responseDump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 打印响应信息
	fmt.Println(string(responseDump))

	// 将响应写回到客户端
	resp.Write(Writer)

	return err
}

func HandleRequest(Writer http.ResponseWriter, Request *http.Request) {
	// 解析请求
	requestDump, err := httputil.DumpRequest(Request, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 打印请求信息
	fmt.Println(string(requestDump))

	fmt.Println(Request.Method)

	fmt.Println(Request.URL)

	fmt.Println(Request.Proto)

	fmt.Println(Request.ProtoMajor)

	fmt.Println(Request.ProtoMinor)

	fmt.Println(Request.Header)

	fmt.Println(Request.Body)

	fmt.Println(Request.ContentLength)

	fmt.Println(Request.Host)

	fmt.Println()

	// 创建一个新的请求对象，并复制原始请求的其他信息
	clientReq := &http.Request{
		Method:        Request.Method,
		URL:           Request.URL,
		Proto:         Request.Proto,
		ProtoMajor:    Request.ProtoMajor,
		ProtoMinor:    Request.ProtoMinor,
		Header:        Request.Header,
		Body:          Request.Body,
		ContentLength: Request.ContentLength,
		Host:          Request.Host,
	}

	err = GetResponse(Writer, clientReq)

	if err != nil {
		fmt.Println(err)
		return
	}

}
