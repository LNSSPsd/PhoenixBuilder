package mainframe

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"

	lua "github.com/yuin/gopher-lua"
)

// 内置函数
type BuiltlnFner interface {
	//实现与lua对接
	NewFrame(L *lua.LState) int
	GetListener(L *lua.LState) int
	GetControl(L *lua.LState) int
	LoadFn(l *lua.LState) error
}

// 实现BuiltFner
type BuiltlnFn struct {
	OmgFrame *LuaComponenter
}

// 写入
func (b *BuiltlnFn) LoadFn(L *lua.LState) error {

	// 创建一个Lua table

	skynet := L.NewTable()

	//注入方法 GetListener GetControl
	L.SetField(skynet, "GetListener", L.NewFunction(b.GetListener))
	L.SetField(skynet, "GetControl", L.NewFunction(b.GetControl))
	// 将table命名为ComplexStruct，并将其设为全局变量
	L.SetGlobal("skynet", skynet)
	return nil
}

func (b *BuiltlnFn) GetListener(L *lua.LState) int {
	listener := L.NewTable()
	//listener的方法 listen("可变参数") 获取参数  listenPackage(Id)

	L.SetField(listener, "listenMsg", L.NewFunction(func(l *lua.LState) int {

		return 1
	}))
	L.SetField(listener, "listenPackage", L.NewFunction(func(l *lua.LState) int {

		return 1
	}))
	//返回listener对象
	L.Push(listener)
	return 1
}
func (b *BuiltlnFn) GetControl(L *lua.LState) int {
	GameControl := L.NewTable()
	L.SetField(GameControl, "SendWsCmd", L.NewFunction(func(l *lua.LState) int {
		if l.GetTop() == 1 {
			args := L.CheckString(1)

			b.OmgFrame.mainFrame.GetGameControl().SendCmdAndInvokeOnResponse(args, func(output *packet.CommandOutput) {
				fmt.Println("测试", output.OutputMessages)
			})

		}
		return 1
	}))
	L.SetField(GameControl, "SendCmdAndInvokeOnResponse", L.NewFunction(func(l *lua.LState) int {
		if l.GetTop() == 1 {
			args := L.CheckString(1)
			b.OmgFrame.mainFrame.GetGameControl().SendCmdAndInvokeOnResponse(args, func(output *packet.CommandOutput) {
				cmdBack := L.NewTable()
				if output.SuccessCount > 0 {
					L.SetField(cmdBack, "Success", lua.LBool(true))
				} else {
					L.SetField(cmdBack, "Success", lua.LBool(false))
				}
				L.SetField(cmdBack, "outputmsg", lua.LString(fmt.Sprintf("%v", output.OutputMessages)))
				L.Push(cmdBack)

			})
		}
		return 1
	}))
	L.Push(GameControl)
	return 1
}

// 指令返回信息
type CmdInvokeResponse struct {
	isSuccess bool
	BackMsg   string
}

// 消息返回信息
type MsgResponse struct {
	playerName string
	Msg        []string
}

// 监听
type Listener interface {
}

// 实现Listener
type Listen struct {
}
type GameControler interface {
	//占位
	GameControler() string
	SendWsCmd(str string)
	SendPlayerCmd(str string)
	SendWsCmdAndInvokeOnResponse(str string, callBack func(*CmdInvokeResponse) bool)
	SetOnParamMsg(playerName string, callBack func(*MsgResponse) bool)
}

// 游行行为 实现了gamecontroler
type GameControl struct {
}

// 让gamecontrol实现gamecontroler
func (g *GameControl) GameControler() string {
	return "this is gamecontrol"
}
func (g *GameControl) SendWsCmd(str string) {
	//todo
}
func (g *GameControl) SendPlayerCmd(str string) {
	//todo
}
func (g *GameControl) SendWsCmdAndInvokeOnResponse(str string, callBack func(*CmdInvokeResponse) bool) {
	//to do
}
func (g *GameControl) SetOnParamMsg(playerName string, callBack func(*MsgResponse) bool) {
	//todo
}
