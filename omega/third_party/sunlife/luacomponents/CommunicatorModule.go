package luaComponent

import (
	"fmt"

	lua "github.com/yuin/gopher-lua"
)

// 注册一个名为 "communicator" 的 Lua 模块，提供双向通信功能
func registerCommunicatorModule(L *lua.LState) int {
	communicator := map[string]lua.LGFunction{
		"send_message": sendMessage,
		"get_message":  getMessage,
	}

	mod := L.SetFuncs(L.NewTable(), communicator)
	L.Push(mod)
	return 1
}

// sendMessage 是一个 Go 函数，用于从 Lua 脚本发送消息
func sendMessage(L *lua.LState) int {
	message := L.Get(1).String()
	fmt.Println("Lua script sent a message:", message)
	return 0
}

// getMessage 是一个 Go 函数，用于向 Lua 脚本发送消息
func getMessage(L *lua.LState) int {
	// 在这个示例中，我们只是简单地返回一个静态字符串
	// 实际应用中，您可能需要根据需要从其他来源获取消息
	message := "Hello from Go!"
	L.Push(lua.LString(message))
	return 1
}
