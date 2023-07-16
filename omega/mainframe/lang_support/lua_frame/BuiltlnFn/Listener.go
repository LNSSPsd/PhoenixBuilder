package BuiltlnFn

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/omega/defines"
	"phoenixbuilder/omega/mainframe/lang_support/lua_frame/definition"
)

// 模拟消息
type Message struct {
	Type    string
	Content string
}
type BuiltListener struct {
	*BuiltlnFn
}

/*
// Listener实现

	func (b *BuiltListener) BuiltlnListner(L *lua.LState) int {
		// 注册Listener类型

		mt := L.NewTypeMetatable("listener")
		L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
			"NextMsg": NextMsg,
		}))
		listener := L.NewTable()
		//listener的方法 listen("可变参数") 获取参数  listenPackage(Id)

		L.SetField(listener, "GetMsgListner", L.NewFunction(b.GetMsgListener))
		L.SetField(listener, "listenPackage", L.NewFunction(func(l *lua.LState) int {

			return 1
		}))
		//返回listener对象
		L.Push(listener)
		return 1
	}
*/
func (b *BuiltListener) BuiltFunc(L *lua.LState) int {
	// 注册Listener类型
	listener := L.NewTable()
	//listener的方法 listen("可变参数") 获取参数  listenPackage(Id)

	L.SetField(listener, "GetMsgListner", L.NewFunction(b.GetMsgListener))
	L.SetField(listener, "listenPackage", L.NewFunction(func(l *lua.LState) int {

		return 1
	}))
	L.SetField(listener, "GetLogInfoer", L.NewFunction(b.GetLoger))
	//返回listener对象
	L.Push(listener)
	return 1
}

// 获取玩家登录登出
func (b *BuiltListener) GetLoger(l *lua.LState) int {
	LogInfoer := l.NewTable()
	l.SetField(LogInfoer, "GetLoginInfo", l.NewFunction(b.GetLoginInfo))
	l.SetField(LogInfoer, "GetLogoutInfo", l.NewFunction(b.GetLogoutInfo))
	l.Push(LogInfoer)
	return 1
}

// 登出
func (b *BuiltListener) GetLogoutInfo(l *lua.LState) int {
	msgChan := make(chan interface{}, 1)
	b.RegisterPackage(&definition.PackageChan{
		PackageType:    definition.LOGOUT_TYPE,
		PackageMsgChan: msgChan,
	})
	logoutMsg := <-msgChan
	logoutTable := l.NewTable()
	switch v := logoutMsg.(type) {
	case protocol.PlayerListEntry:
		l.SetField(logoutTable, "Name", lua.LString(v.Username))
	default:
		fmt.Println("登陆包解析失败")
	}
	l.Push(logoutTable)
	return 1
}

// 登进
func (b *BuiltListener) GetLoginInfo(l *lua.LState) int {
	msgChan := make(chan interface{}, 1)
	b.RegisterPackage(&definition.PackageChan{
		PackageType:    definition.LOGIN_TYPE,
		PackageMsgChan: msgChan,
	})
	logoutMsg := <-msgChan
	logoutTable := l.NewTable()
	switch v := logoutMsg.(type) {
	case protocol.PlayerListEntry:
		l.SetField(logoutTable, "Name", lua.LString(v.Username))
	default:
		fmt.Println("登出包解析失败")
	}
	l.Push(logoutTable)
	return 1
}

// NextMsg 用于从监听器的消息通道中获取下一个消息
func (f *BuiltListener) NextMsg(L *lua.LState) int {
	msgChan := make(chan interface{}, 1)
	msgType := definition.MSG_TYPE
	f.RegisterPackage(&definition.PackageChan{
		PackageType:    msgType,
		PackageMsgChan: msgChan,
	})
	msg := <-msgChan
	switch v := msg.(type) {
	case *defines.GameChat:
		Name := v.Name
		newMsg := ""
		for _, key := range v.Msg {
			newMsg += key + " "
		}
		L.Push(lua.LString(Name))
		L.Push(lua.LString(newMsg))
		return 2
	default:
		fmt.Println("无法解析")
		L.ArgError(1, "无法解析新的玩家消息")
		L.Push(lua.LString(""))
		L.Push(lua.LString(""))
	}
	return 2
}

// GetListener 创建一个新的监听器并返回其引用
func (f *BuiltListener) GetMsgListener(L *lua.LState) int {
	listener := L.NewTable()
	L.SetField(listener, "NextMsg", L.NewFunction(f.NextMsg))
	L.Push(listener)
	return 1
}
