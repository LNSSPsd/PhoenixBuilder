package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
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

func main() {
	address := ":8888"

	http.HandleFunc("/ws_test", func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := (&websocket.Upgrader{
			ReadBufferSize:  ReadBufferSize,
			WriteBufferSize: WriteBufferSize,
		}).Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		go func() {
			for  {
				msgType, data, err :=wsConn.ReadMessage()
				if err!=nil{
					println(err)
					return
				}
				fmt.Printf("%v %v\n",msgType,data)
				err = wsConn.WriteMessage(msgType, data)
				if err!=nil{
					println(err)
					return
				}
			}

		}()
	})
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

