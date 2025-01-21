package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许任何来源的连接
		return true
	},
}

// WebSocket 连接管理结构体
type WebSocketServer struct {
	connections map[*websocket.Conn]bool
	mu          sync.Mutex
}

func (server *WebSocketServer) addConnection(conn *websocket.Conn) {
	server.mu.Lock()
	defer server.mu.Unlock()
	server.connections[conn] = true
}

func (server *WebSocketServer) removeConnection(conn *websocket.Conn) {
	server.mu.Lock()
	defer server.mu.Unlock()
	delete(server.connections, conn)
}

func (server *WebSocketServer) broadcastMessage(message []byte) {
	server.mu.Lock()
	defer server.mu.Unlock()

	// 向所有连接的客户端发送消息
	for conn := range server.connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error sending message to client:", err)
			conn.Close()
			delete(server.connections, conn)
		}
	}
}

func handleWebSocketConnection(w http.ResponseWriter, r *http.Request, server *WebSocketServer) {
	// 升级 HTTP 连接为 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket Upgrade Error:", err)
		return
	}
	defer conn.Close()

	// 添加连接到服务器的连接池中
	server.addConnection(conn)
	defer server.removeConnection(conn)

	// 双向通信：服务端和客户端都可以发送和接收消息
	for {
		// 接收来自客户端的消息
		messageType, msg, err := conn.ReadMessage()
		fmt.Println("messageType:", messageType)
		if err != nil {
			log.Println("WebSocket Read Error:", err)
			break
		}
		fmt.Printf("Received from client: %s\n", msg)

		// 将收到的消息广播给所有连接的客户端
		server.broadcastMessage(msg)
	}
}

func main() {
	server := &WebSocketServer{
		connections: make(map[*websocket.Conn]bool),
	}

	// 启动 WebSocket 服务
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocketConnection(w, r, server)
	})

	// 启动 HTTP 服务，监听 8080 端口
	log.Println("Starting WebSocket server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe failed: ", err)
	}
}
