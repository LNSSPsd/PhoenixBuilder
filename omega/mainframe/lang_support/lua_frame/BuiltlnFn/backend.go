package BuiltlnFn

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"phoenixbuilder/omega/mainframe/lang_support/lua_frame/definition"
)

// 终端输入输出接口
type BuiltBackEnder struct {
	*BuiltlnFn
	BackEndMsgChan chan string
}

// 获取终端操控中心
func (b *BuiltBackEnder) BuiltFunc(L *lua.LState) int {

	//发送消息
	BackEnder := L.NewTable()
	L.SetField(BackEnder, "GetMsg", L.NewFunction(b.GetBackEndMsg))
	L.Push(BackEnder)
	return 1
}

// 获取消息
func (b *BuiltBackEnder) GetBackEndMsg(L *lua.LState) int {
	//创造一个用于传输数据的通道
	msgChan := make(chan interface{}, 1)
	msgType := definition.BACK_END_TYPE
	b.RegisterPackage(&definition.PackageChan{
		PackageType:    msgType,
		PackageMsgChan: msgChan,
	})
	msg := <-msgChan
	switch v := msg.(type) {
	case []string:
		// 获取命令行参数
		reMsg := ""
		for _, key := range v {
			reMsg += key + " "
		}
		// 将字符串数组压入栈顶
		L.Push(lua.LString(reMsg))
		return 1
	default:
		fmt.Println("backendmsg获取意料之外的msg")
	}
	L.Push(lua.LString(""))
	return 1
}
