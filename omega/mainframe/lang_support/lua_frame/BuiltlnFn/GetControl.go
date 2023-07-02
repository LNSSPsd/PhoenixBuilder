package BuiltlnFn

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/omega/defines"
)

type BuiltGameControler struct {
	*BuiltlnFn
}

func (b *BuiltGameControler) BuiltFunc(L *lua.LState) int {
	GameControl := L.NewTable()
	L.SetField(GameControl, "SendWsCmd", L.NewFunction(b.SendWsCmd))
	L.SetField(GameControl, "SendCmdAndInvokeOnResponse", L.NewFunction(b.SendCmdAndInvokeOnResponse))
	//等待说话
	L.SetField(GameControl, "SetOnParamMsg", L.NewFunction(b.SetOnParamMsg))
	L.Push(GameControl)
	return 1
}

/*
	func (b *BuiltGameControler) BuiltGameContrler(L *lua.LState) int {
		GameControl := L.NewTable()
		L.SetField(GameControl, "SendWsCmd", L.NewFunction(b.SendWsCmd))
		L.SetField(GameControl, "SendCmdAndInvokeOnResponse", L.NewFunction(b.SendCmdAndInvokeOnResponse))
		//等待说话
		L.SetField(GameControl, "SetOnParamMsg", L.NewFunction(b.SetOnParamMsg))
		L.Push(GameControl)
		return 1
	}
*/
func (b *BuiltGameControler) SendWsCmd(L *lua.LState) int {
	args := L.CheckString(1)
	b.OmegaFrame.MainFrame.GetGameControl().SendCmd(args)

	return 1
}
func (b *BuiltGameControler) SendCmdAndInvokeOnResponse(L *lua.LState) int {
	if L.GetTop() == 1 {
		args := L.CheckString(1)
		ch := make(chan bool)
		b.OmegaFrame.MainFrame.GetGameControl().SendCmdAndInvokeOnResponse(args, func(output *packet.CommandOutput) {
			cmdBack := L.NewTable()
			if output.SuccessCount > 0 {
				L.SetField(cmdBack, "Success", lua.LBool(true))
			} else {
				L.SetField(cmdBack, "Success", lua.LBool(false))
			}
			L.SetField(cmdBack, "outputmsg", lua.LString(fmt.Sprintf("%v", output.OutputMessages)))
			L.Push(cmdBack)
			ch <- true
		})
		<-ch
	} else {
		fmt.Println("参数应该仅有一个")
	}
	return 1
}
func (b *BuiltGameControler) SetOnParamMsg(L *lua.LState) int {
	if L.GetTop() == 1 {
		name := L.CheckString(1)
		ch := make(chan bool)
		b.OmegaFrame.MainFrame.GetGameControl().SetOnParamMsg(name, func(chat *defines.GameChat) (catch bool) {
			msg := ""
			for _, v := range chat.Msg {
				msg += v + " "
			}
			L.Push(lua.LString(msg))
			ch <- true
			return false
		})
		<-ch
	} else {
		fmt.Println("参数应该仅有一个")
	}

	return 1
}
