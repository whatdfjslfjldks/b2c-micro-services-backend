package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"sync"
)

var (
	userConnections    = make(map[string]*websocket.Conn) // 用户 WebSocket 连接
	supportConnections = make(map[string]*websocket.Conn) // 客服 WebSocket 连接
	userSupportMapping = make(map[string]string)          // 映射：用户ID -> 客服ID
	mutex              sync.Mutex
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleChat(w http.ResponseWriter, r *http.Request) {
	// 获取用户ID
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	// 保存用户连接
	mutex.Lock()
	userConnections[userID] = conn
	mutex.Unlock()

	log.Printf("User %s connected", userID)

	// 随机选一个客服连接
	var connections []*websocket.Conn
	var keys []string
	for key, conn := range supportConnections {
		fmt.Println("111:", key)
		keys = append(keys, key)
		connections = append(connections, conn)
	}
	fmt.Println("keys:", keys)

	// 如果没有客服连接，返回错误
	if len(connections) == 0 {
		http.Error(w, "No support available", http.StatusServiceUnavailable)
		return
	}

	// 使用随机数选择一个客服连接
	idx := rand.Intn(len(connections)) // 生成一个随机索引
	supportConn := connections[idx]
	supportId := keys[idx]
	fmt.Println("supportId:", supportId)

	// 映射用户与客服
	mutex.Lock()
	userSupportMapping[userID] = supportId
	mutex.Unlock()

	// 发送已连接消息给用户
	type Message struct {
		Sender string `json:"sender"`
		Text   string `json:"text"`
	}

	message := Message{
		Sender: "tip",
		Text:   fmt.Sprintf("已经连接到客服: support_%d", idx),
	}

	// 将结构体转换为 JSON 格式
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.Fatal("Error marshaling JSON:", err)
	}

	// 返给用户一个已连接到客服的消息
	err = conn.WriteMessage(websocket.TextMessage, jsonData)
	if err != nil {
		log.Println("Error sending message to user:", err)
		return
	}

	// 监听用户消息
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from user:", err)
			break
		}

		// 将用户的消息转发给指定的客服
		mutex.Lock()
		err = supportConn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("Error sending message to support:", err)
		}
		mutex.Unlock()
	}
}

func handleSupport(w http.ResponseWriter, r *http.Request) {
	// 获取客服ID
	supportID := r.URL.Query().Get("supportID")
	if supportID == "" {
		http.Error(w, "supportID is required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error during connection upgrade:", err)
		return
	}
	defer conn.Close()

	// 保存客服连接
	mutex.Lock()
	supportConnections[supportID] = conn
	mutex.Unlock()

	log.Printf("Support %s connected", supportID)

	// 监听客服消息
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from support:", err)
			break
		}

		// 将客服的消息转发给指定的用户
		mutex.Lock()
		for userID, assignedSupportID := range userSupportMapping {
			fmt.Println("assignedSupportID:", assignedSupportID)
			fmt.Println("userID:", userID)
			fmt.Println("supportID:", supportID)
			if assignedSupportID == supportID {
				// 发送消息给该用户
				userConn := userConnections[userID]
				err := userConn.WriteMessage(websocket.TextMessage, msg)
				if err != nil {
					log.Println("Error sending message to user:", err)
				}
			}
		}
		mutex.Unlock()
	}
}

func main() {
	http.HandleFunc("/chat", handleChat)       // 用户连接
	http.HandleFunc("/support", handleSupport) // 客服连接

	log.Println("Server starting on :8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
