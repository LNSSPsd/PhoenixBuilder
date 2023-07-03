package BuiltlnFn

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/yuin/gopher-lua"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var conn *websocket.Conn

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	// 升级HTTP连接为WebSocket连接
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("升级连接为WebSocket失败：", err)
		return
	}
	defer conn.Close()

	// 处理WebSocket连接
	for {
		// 读取客户端发送的消息
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息失败：", err)
			break
		}
		conn.WriteMessage(websocket.TextMessage, []byte("hello"))
		fmt.Println("收到消息：", string(message))
	}
}
func BuiltnWebSokcet(L *lua.LState) int {

	// 将WebSocket服务器函数暴露给Lua 返回开启成功与否
	L.SetGlobal("startWebSocketServer", L.NewFunction(func(L *lua.LState) int {
		if L.GetTop() == 1 {
			addr := L.CheckString(1)
			http.HandleFunc("/ws", websocketHandler)
			go http.ListenAndServe(addr, nil)
			fmt.Println("开始监听")

			L.Push(lua.LBool(true))
		} else {
			fmt.Println("错误 websocket连接应该有一个参数 描述你的端口")
			L.Push(lua.LBool(false))
		}
		return 1
	}))

	// 将发送信息的函数暴露给Lua
	L.SetGlobal("sendMessage", L.NewFunction(func(L *lua.LState) int {
		msg := L.ToString(1)
		if conn == nil {
			L.ArgError(1, "connection not established")
			return 0
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
			log.Println(err)
		}
		return 0
	}))

	// 将接收信息的函数暴露给Lua
	L.SetGlobal("receiveMessage", L.NewFunction(func(L *lua.LState) int {
		if conn == nil {
			L.ArgError(1, "connection not established")
			return 0
		}
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			L.Push(lua.LNil)
		} else {
			L.Push(lua.LString(string(msg)))
		}
		return 1
	}))

	// 将断开连接的函数暴露给Lua
	L.SetGlobal("closeConnection", L.NewFunction(func(L *lua.LState) int {
		if conn != nil {
			if err := conn.Close(); err != nil {
				log.Println(err)
			}
			conn = nil
		}
		return 0
	}))

	// 将检查连接状态的函数暴露给Lua
	L.SetGlobal("isConnected", L.NewFunction(func(L *lua.LState) int {
		L.Push(lua.LBool(conn != nil))
		return 1
	}))
	return 0
}
