package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"micro-services/support-server/pkg/config"
	"net/http"
	"sync"
	"time"
)

var connections = struct {
	sync.RWMutex
	m map[string]*websocket.Conn
}{m: make(map[string]*websocket.Conn)}

// 存储连接的内存映射,这个是客服自己
//var supportConnections = make(map[string]*websocket.Conn)

//// 存储连接的内存映射，这个是用户客服对话的
//var userSppConnections = make(map[string]*websocket.Conn)

type Message struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

// WebSocket连接升级器
var upgrade = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许所有来源
		return true
	},
}

// Connect 连接用户和客服，
func Connect(c *gin.Context) {
	userID := c.Query("userID")
	supportID := c.Query("supportID")

	// 升级为WebSocket连接
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  "升级WS失败！",
		})
		return
	}

	// 创建聊天室的唯一ID，可以通过 userID 和 supportID 组合生成
	roomID := userID + "_" + supportID

	// 将连接存入 内存映射
	//userSppConnections[roomID] = conn
	connections.Lock()
	connections.m[roomID] = conn
	connections.Unlock()

	connections.RLock()
	sConn := connections.m[supportID]
	connections.RUnlock()

	if sConn != nil {
		msg := Message{Sender: "attention", Text: roomID}
		jsonData, _ := json.Marshal(msg)
		sConn.WriteMessage(websocket.TextMessage, jsonData)
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from user:", err)
			break
		}
		targetConn := connections.m[roomID+"_support"]
		targetConn.WriteMessage(websocket.TextMessage, msg)
	}
}

func SupportConnect(c *gin.Context) {
	supportID := c.Query("supportID")
	fmt.Println("supportID:", supportID)

	// 升级为WebSocket连接
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade Error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"code": 500,
			"msg":  "升级WS失败！",
		})
		return
	}
	defer conn.Close()

	// 将连接存储到内存映射中
	connections.Lock()
	connections.m[supportID] = conn
	connections.Unlock()

	// 在 Redis 中设置状态
	config.RdClient.HMSet(config.Ctx, supportID, "status", "1")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message from support %s: %v", supportID, err)
			}
			// 出错，删除内存映射，更新 Redis 状态
			connections.Lock()
			delete(connections.m, supportID)
			connections.Unlock()

			config.RdClient.HDel(config.Ctx, supportID, "status")

			log.Printf("WebSocket connection closed for support %s", supportID)
			break
		}

		if messageType == websocket.CloseMessage {
			log.Printf("Received close message from support %s", supportID)
			// 删除内存映射，更新 Redis 状态
			connections.Lock()
			delete(connections.m, supportID)
			connections.Unlock()

			config.RdClient.HDel(config.Ctx, supportID, "status")

			log.Printf("WebSocket connection closed for support %s", supportID)
			break
		}

		// 把消息发送回去，这里发的是用户和客服的房间号 userID_supportID
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Error sending message to support %s: %v", supportID, err)
			break
		}
	}
}
func ConnectRoom(c *gin.Context) {
	roomID := c.Query("roomID")

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade Error in room %s: %v", roomID, err)
		return
	}
	defer conn.Close()

	connections.RLock()
	userConn := connections.m[roomID]
	connections.RUnlock()

	a := Message{
		Sender: "tip",
		Text:   "房间不存在",
	}
	j, err := json.Marshal(a)
	if err != nil {
		log.Println("json.Marshal err:", err)
	}
	if userConn == nil {
		conn.WriteMessage(websocket.CloseMessage, j)
		return
	}

	connections.Lock()
	connections.m[roomID+"_support"] = conn
	connections.Unlock()

	defer func() {
		connections.Lock()
		delete(connections.m, roomID+"_support")
		connections.Unlock()
	}()

	go func() {
		for {
			select {
			case <-time.After(30 * time.Second):
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					log.Printf("Error sending ping message in room %s: %v", roomID, err)
					return
				}
			}
		}
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message from support in room %s: %v", roomID, err)
			break
		}

		if err := userConn.WriteMessage(messageType, message); err != nil {
			log.Printf("Error forwarding message to user in room %s: %v", roomID, err)
			userConn.Close()
			break
		}
	}
}
