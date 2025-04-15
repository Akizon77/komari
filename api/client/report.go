package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"komari/database"
	"komari/ws"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func UploadReport(c *gin.Context) {
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println("Failed to read request body:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var data map[string]interface{}
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	err = database.SaveReport(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%v", err)})
		return
	}
	uuid, err := database.GetClientUUIDByToken(data["token"].(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("%v", err)})
		return
	}
	delete(data, "token")
	ws.LatestReport[uuid] = data
	//log.Println(string(bodyBytes))
	c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the body for further use
	c.JSON(200, gin.H{"status": "success"})
}

func WebSocketReport(c *gin.Context) {
	if !websocket.IsWebSocketUpgrade(c.Request) {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Require WebSocket upgrade"})
		return
	}
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": "Failed to upgrade to WebSocket"})
		return
	}
	defer conn.Close()

	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message:", err)
		return
	}

	data := map[string]interface{}{}
	err = json.Unmarshal(message, &data)
	if err != nil {
		conn.WriteJSON(gin.H{"status": "error", "error": "Invalid JSON"})
		return
	}
	// it should ok,token was verfied in the middleware
	token := ""
	var errMsg string

	// 优先检查查询参数中的 token
	if token_, success := c.Params.Get("token"); success && token_ != "" {
		token = token_
	} else if data != nil && data["token"] != nil {
		if t, ok := data["token"].(string); ok && t != "" {
			token = t
		} else {
			errMsg = "Invalid token format in data"
		}
	} else {
		errMsg = "Token not provided"
	}

	// 如果 token 为空，返回错误
	if token == "" {
		conn.WriteJSON(gin.H{"status": "error", "error": errMsg})
		return
	}

	// Check if a connection with the same token already exists
	if _, exists := ws.ConnectedClients[token]; exists {
		conn.WriteJSON(gin.H{"status": "error", "error": "Token already in use"})
		return
	}
	ws.ConnectedClients[token] = conn
	defer func() {
		delete(ws.ConnectedClients, token)
	}()

	clientUUID, err := database.GetClientUUIDByToken(token)
	if err != nil {
		conn.WriteJSON(gin.H{"status": "error", "error": fmt.Sprintf("%v", err)})
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		err = json.Unmarshal(message, &data)
		if err != nil {
			break
		}

		report := database.ParseReport(data)

		err = database.SaveClientReport(clientUUID, report)
		if err != nil {
			conn.WriteJSON(gin.H{"status": "error", "error": fmt.Sprintf("%v", err)})
		}

		uuid, err := database.GetClientUUIDByToken(token)
		if err != nil {
			conn.WriteJSON(gin.H{"status": "error", "error": fmt.Sprintf("%v", err)})
			return
		}
		delete(data, "token")
		ws.LatestReport[uuid] = data
	}
}
