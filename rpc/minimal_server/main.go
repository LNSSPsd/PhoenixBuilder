package main

import (
	"RPC/channel"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "home.html")
}

const (
	ReadBufferSize  = 1024
	WriteBufferSize = 1024
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func readPump(connectC2S chan []byte, wsConn *websocket.Conn) {
	defer func() {
		wsConn.Close()
	}()
	wsConn.SetReadLimit(maxMessageSize)
	wsConn.SetReadDeadline(time.Now().Add(pongWait))
	wsConn.SetPongHandler(func(string) error { wsConn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		msgType, message, err := wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if msgType == websocket.BinaryMessage {
			connectC2S <- message
		}
	}
}

func writePump(connectS2C chan []byte, wsConn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		wsConn.Close()
	}()
	for {
		select {
		case message, ok := <-connectS2C:
			wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				wsConn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			err := wsConn.WriteMessage(websocket.BinaryMessage, message)
			if err != nil {
				return
			}
		case <-ticker.C:
			wsConn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := wsConn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func main() {
	address := ":8888"
	config := &MajorConfig{}
	majorFunc := InitMajorFunction(config)

	for pattern, handleFn := range majorFunc.PageServes {
		http.HandleFunc(pattern, handleFn)
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := (&websocket.Upgrader{
			ReadBufferSize:  ReadBufferSize,
			WriteBufferSize: WriteBufferSize,
		}).Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		connectS2C := make(chan []byte)
		connectC2S := make(chan []byte)
		go readPump(connectC2S, wsConn)
		go writePump(connectS2C, wsConn)
		baseConnection := &channel.ChanChanel{
			ConnectW: connectS2C,
			ConnectR: connectC2S,
			IsClosed: false,
		}
		encryptionConn := &channel.EncryptedChannel{
			Connect:    baseConnection,
			IsInitator: true,
		}
		err = encryptionConn.Init()
		if err != nil {
			log.Printf("Encryption Init Error (%v)", err)
			wsConn.Close()
			return
		}
		connMux := channel.NewMux(encryptionConn)
		majorFunc.Connect(connMux, func() { wsConn.Close() })
	})
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
