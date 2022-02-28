package main

import (
	"RPC/channel"
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/websocket"
)

func readPump(connectS2C chan []byte, wsConn *websocket.Conn) {
	defer wsConn.Close()
	for {
		msgType, message, err := wsConn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		if msgType == websocket.BinaryMessage {
			connectS2C <- message
		}
	}
}

func writePump(connectC2S chan []byte, wsConn *websocket.Conn) {
	defer wsConn.Close()
	for {
		message, ok := <-connectC2S
		if !ok {
			return
		}
		err := wsConn.WriteMessage(websocket.BinaryMessage, message)
		if err != nil {
			return
		}
	}
}

func initEncryptionConnection(wsConn *websocket.Conn) (*channel.Mux, error) {
	connectS2C := make(chan []byte)
	connectC2S := make(chan []byte)
	go readPump(connectS2C, wsConn)
	go writePump(connectC2S, wsConn)
	baseConnection := &channel.ChanChanel{
		ConnectW: connectC2S,
		ConnectR: connectS2C,
		IsClosed: false,
	}
	encryptionConn := &channel.EncryptedChannel{
		Connect:    baseConnection,
		IsInitator: false,
	}
	err := encryptionConn.Init()
	if err != nil {
		return nil, fmt.Errorf("Encryption Init Error (%v)")
	}
	connMux := channel.NewMux(encryptionConn)
	return connMux, nil

}

func main() {
	host := "ws://127.0.0.1:8888/ws"
	wsConn, _, err := websocket.DefaultDialer.Dial(host, nil)
	if err != nil {
		panic(err)
	}
	connMux, err := initEncryptionConnection(wsConn)
	if err != nil {
		panic(err)
	}
	channel3 := connMux.GetSubChannel(3)
	msgGet := channel3.Get()
	fmt.Println(string(msgGet))
	err = channel3.Send(msgGet)
	if err != nil {
		panic(err)
	}
	bufio.NewReader(os.Stdin).ReadLine()
}
