package internal

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
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

// SupportConnect 建立客服与服务端连接，用于推送用户连接房间
// 把客服id存到redis，还有conn连接
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

	//supportConnections[supportID] = conn
	connections.Lock()
	connections.m[supportID] = conn
	connections.Unlock()

	if err != nil {
		log.Println("Error storing support connection in Redis:", err)
		return
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from support:", err)
			// 出错，删除内存映射，中断连接
			//delete(supportConnections, supportID)
			break
		}
		// 把消息发送回去，这里发的是用户和客服的房间号 userID_supportID
		err = conn.WriteMessage(websocket.TextMessage, message)
	}
}

func ConnectRoom(c *gin.Context) {
	roomID := c.Query("roomID")

	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket Upgrade Error: %v", err)
		return
	}

	connections.RLock()
	userConn := connections.m[roomID]
	connections.RUnlock()

	if userConn == nil {
		conn.WriteMessage(websocket.CloseMessage, []byte("房间不存在"))
		conn.Close()
		return
	}

	// 存储客服端连接
	connections.Lock()
	connections.m[roomID+"_support"] = conn
	connections.Unlock()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message from support:", err)
			// 出错，删除内存映射，中断连接
			//delete(supportConnections, supportID)
			break
		}
		userConn.WriteMessage(websocket.TextMessage, message)
	}
}
