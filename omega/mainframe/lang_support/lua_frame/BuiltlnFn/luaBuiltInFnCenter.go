package BuiltlnFn

import (
	"phoenixbuilder/omega/defines"
	omgApi "phoenixbuilder/omega/mainframe/lang_support/lua_frame/omgcomponentapi"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

// 内置函数加载器
type BuiltlnFn struct {
	//omg组件
	OmegaFrame *omgApi.OmgApi
	Listener   sync.Map
	mainframe  defines.MainFrame
}

// 载入
func (b *BuiltlnFn) LoadFn(L *lua.LState) error {
	// 创建一个Lua table
	//注册skynet
	skynet := L.NewTable()
	//注入
	SkynetBuiltFnDic := b.GetSkynetBuiltlnFunction()
	for k, v := range SkynetBuiltFnDic {
		L.SetField(skynet, k, L.NewFunction(v.BuiltFunc))
	}
	// 将table命名为ComplexStruct，并将其设为全局变量
	L.SetGlobal("skynet", skynet)
	//内置websocket连接
	BuiltnWebSokcet(L)
	return nil
}
