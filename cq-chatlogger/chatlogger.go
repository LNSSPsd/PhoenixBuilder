package cqchat

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol/packet"
	"strconv"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pterm/pterm"
)

// 将http升级为websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// CQMessages 来自go-cq端的消息
var CQMessages chan IMessage

var MCMessages chan *packet.Text
var Has_Connected bool
var Conn *minecraft.Conn
var ServerID string

func GlobalConn(conn *minecraft.Conn, SID string) {
	Conn = conn
	ServerID = SID
}

// receiveMessage 接收go-cq的消息 转发到main
func receiveMessage(cqconn *websocket.Conn) {
	fmt.Println("receive start")
	for {
		// todo 优化
		// msgType为0时 消息正常接收 其他未知
		msgType, data, err := cqconn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			cqconn.Close()
		}
		if msgType != 0 {
			fmt.Println(string(data))
		}
		// 先解析出事件种类(event或message)
		post, err := ParseMetaPost(data)
		if post.PostType == "meta_event" && post.MetaEventType == "lifecycle"{
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
		continue
	}
}

// SendMessage 接收游戏内消息 发送至go-cq
func sendMessage(cqconn *websocket.Conn) {
	for {
		msg, ok := <- MCMessages
		fmt.Println(Conn.IdentityData())
		if !ok || len(msg.Message) == 0 || msg.SourceName == Conn.IdentityData().DisplayName{
			continue
		}
		fmsg := FormatMCMessage(*msg)
		echo, _ := uuid.NewUUID()//标识符 
		fmt.Println(fmsg)
		qmsg := QMessage{
			Action: "send_group_msg",
			Params: struct{
				GroupID int64 `json:"group_id"`
				Message string `json:"message"`
			}{
				GroupID: Setting.DefaultGroupID,
				Message: fmsg,
			},
			Echo: echo.String(),
		}
		fmt.Println(qmsg)
		data, _ := json.Marshal(qmsg)
		cqconn.WriteMessage(1, data)
	}
}
//websockect链接
func handleFunc(w http.ResponseWriter, r *http.Request) {
	cqconn, _ := upgrader.Upgrade(w, r, nil)//websocket链接通道?
	// 这个conn是和gocq连接的
	// 所有和游戏内的conn连接(发包)的操作全部在main.go里
	go receiveMessage(cqconn)
	go sendMessage(cqconn)
}

func Run() {
	pterm.Println(pterm.Yellow("尝试启动chatlogger."))
	http.HandleFunc("/fastbuilder/chatlogger", handleFunc)
	err := http.ListenAndServe(Setting.Port, nil)	
	
	if err != nil {
		log.Println("CONNECTION ERROR:", err)
	}
}
