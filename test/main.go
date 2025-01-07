package main

import (
	"fmt"
	"net/http"
	"strings"
)

func getIPFromRequest(req *http.Request) string {
	// 1. 尝试从 "X-Forwarded-For" 获取用户的 IP 地址
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		// 如果有多个代理，"X-Forwarded-For" 会包含多个 IP 地址，逗号分隔，取第一个即为客户端的真实 IP
		ips := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(ips[0])
	}

	// 2. 如果没有 "X-Forwarded-For"，则直接获取 "RemoteAddr"
	return req.RemoteAddr
}

func handler(w http.ResponseWriter, req *http.Request) {
	// 获取客户端的 IP 地址
	clientIP := getIPFromRequest(req)

	// 返回 IP 地址
	fmt.Fprintf(w, "Your IP Address is: %s", clientIP)
}

func main() {
	// 设置路由
	http.HandleFunc("/get-ip", handler)

	// 启动 HTTP 服务器
	fmt.Println("Server started at :8080")
	http.ListenAndServe(":8080", nil)
}
