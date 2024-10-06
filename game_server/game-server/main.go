package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var Upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type USER struct {
	user_id int `json:"user_id"`
	health  int `json:"health"`
}

type Move struct {
	X        int `json:"x"`
	Y        int `json:"y"`
	playerID int `json:"playerID"`
}

type Shoot struct {
	Playerid int `json:"playerhurt"`
	damage   int `json:"damage"`
}

type Block struct {
	Index     int    `json:"index"`
	Timestamp string `json:"timestamp"`
	Data      string `json:"data"`
	PrevHash  string `json:"prev_hash"`
	Hash      string `json:"hash"`
}

var Blockchain []Block

var mu sync.Mutex
var clients = make(map[int]*websocket.Conn)

func CalculateHash(block Block) string {
	record := strconv.Itoa(block.Index) + block.Timestamp + block.Data + block.PrevHash
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateBlock(oldBlock Block, data string) Block {
	newBlock := Block{
		Index:     oldBlock.Index + 1,
		Timestamp: time.Now().String(),
		Data:      data,
		PrevHash:  oldBlock.Hash,
		Hash:      "",
	}
	newBlock.Hash = CalculateHash(newBlock)
	return newBlock
}

func AddBlock(newBlock Block) {
	mu.Lock()
	Blockchain = append(Blockchain, newBlock)
	mu.Unlock()
}

func broadcastMessage(message []byte) {
	mu.Lock()
	defer mu.Unlock()
	for _, conn := range clients {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error sending message:", err)
			conn.Close()
		}
	}
}

func sra_damage(player *USER, damage int) []byte {
	player.health -= damage
	var message []byte
	if player.health <= 0 {
		player.health = 0
		message = []byte(fmt.Sprintf("User %d is dead", player.user_id))
	} else {
		message = []byte(fmt.Sprintf("User %d took %d damage", player.user_id, damage))
	}
	log.Printf("User %d took %d damage", player.user_id, damage)

	data := fmt.Sprintf("User %d took %d damage", player.user_id, damage)
	newBlock := GenerateBlock(Blockchain[len(Blockchain)-1], data)
	AddBlock(newBlock)

	return message
}

func sra_move(player *USER, move Move) []byte {
	var message []byte
	log.Printf("User %d moved to %d, %d", player.user_id, move.X, move.Y)
	message = []byte(fmt.Sprintf("User %d moved to %d, %d", player.user_id, move.X, move.Y))

	data := fmt.Sprintf("User %d moved to %d, %d", player.user_id, move.X, move.Y)
	newBlock := GenerateBlock(Blockchain[len(Blockchain)-1], data)
	AddBlock(newBlock)

	return message
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrade.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	var user USER
	if err := conn.ReadJSON(&user); err != nil {
		log.Println("Error reading user:", err)
		return
	}

	mu.Lock()
	clients[user.user_id] = conn
	mu.Unlock()

	initialMessage := []byte("Hello, World!")
	if err := conn.WriteMessage(websocket.TextMessage, initialMessage); err != nil {
		log.Println("Error writing message:", err)
	}

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if messageType == websocket.PingMessage || messageType == websocket.PongMessage {
			continue
		}

		var message struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(p, &message); err != nil {
			continue
		}

		switch message.Type {
		case "move":
			var move Move
			if err := json.Unmarshal(p, &move); err != nil {
				continue
			}
			player := &USER{user_id: move.playerID, health: 100}
			msg := sra_move(player, move)
			broadcastMessage(msg)

		case "shoot":
			var shoot Shoot
			if err := json.Unmarshal(p, &shoot); err != nil {
				continue
			}
			player := &USER{user_id: shoot.Playerid, health: 100}
			msg := sra_damage(player, shoot.damage)
			broadcastMessage(msg)

		case "join":
			mu.Lock()
			clients[user.user_id] = conn
			mu.Unlock()
		}

		chainBytes, _ := json.Marshal(Blockchain)
		broadcastMessage(chainBytes)
	}
}

func main() {
	genesisBlock := Block{Index: 0, Timestamp: time.Now().String(), Data: "Genesis Block", PrevHash: "", Hash: ""}
	genesisBlock.Hash = CalculateHash(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	http.HandleFunc("/ws", wshandler)
	log.Println("WebSocket server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe error:", err)
	}
}
