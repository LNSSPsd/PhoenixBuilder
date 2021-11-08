package cqchat

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/pterm/pterm"
	"log"
	"net/http"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol/packet"
	"strconv"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// CQMessages 来自go-cq端的消息
var CQMessages chan IMessage

var MCMessages chan packet.Text

var Conn minecraft.Conn

// receiveMessage 接收并处理协议端的消息
func receiveMessage(conn *websocket.Conn) {
	fmt.Println("receive start")
	for {
		// todo 优化
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			conn.Close()
		}
		if msgType != 0 {
			fmt.Println(string(data))
		}
		post, err := ParseMetaPost(data)
		if post.PostType == "meta_event" && post.MetaEventType == "lifecycle" {
			pterm.Println(pterm.Yellow("已成功连接: " + strconv.Itoa(post.SelfID)))
		}
		if post.PostType == "message" && err == nil {
			action, err := GetMessageData(data)
			if err != nil || action == nil {
				continue
			}
			fmt.Println(action)
			CQMessages <- action
		}

	}
}

// SendMessage
func sendMessage(conn *websocket.Conn) {
	for {
		msg := <-MCMessages
		if len(msg.Message) == 0 {
			continue
		}
		// todo
	}
}

func handleFunc(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)
	go receiveMessage(conn)
}

func Run() {
	pterm.Println(pterm.Yellow("尝试启动chatlogger."))
	http.HandleFunc("/fastbuilder/chatlogger", handleFunc)
	err := http.ListenAndServe(Setting.Port, nil)
	if err != nil {
		log.Println("CONNECTION ERROR:", err)
	}
}
