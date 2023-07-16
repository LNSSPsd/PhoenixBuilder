package BuiltlnFn

import (
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/omega/defines"
	"phoenixbuilder/omega/mainframe/lang_support/lua_frame/definition"
	omgApi "phoenixbuilder/omega/mainframe/lang_support/lua_frame/omgcomponentapi"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

// 内置函数加载器
type BuiltlnFn struct {
	//omg组件
	OmegaFrame       *omgApi.OmgApi
	Listener         sync.Map
	mainframe        defines.MainFrame
	PackageChanSlice []*definition.PackageChan
}

// 将结构体塞入通道内
func (m *BuiltlnFn) RegisterPackage(packageChan *definition.PackageChan) {
	m.PackageChanSlice = append(m.PackageChanSlice, packageChan)
}

// 删除指定元素
func (m *BuiltlnFn) RemovePackageChan(num int) {
	if num < 0 || num >= len(m.PackageChanSlice) {
		return
	}
	m.PackageChanSlice = append(m.PackageChanSlice[:num], m.PackageChanSlice[num+1:]...)
}

// 注入消息
func (m *BuiltlnFn) PackageInjectIntoChan(data interface{}, PackageType string) {
	if m.PackageChanSlice == nil {
		return
	}
	for k, v := range m.PackageChanSlice {
		if v.PackageType == PackageType {
			select {
			case v.PackageMsgChan <- data:
			default:
			}
			m.RemovePackageChan(k)
		}
	}
}

// 分发监听包
func (m *BuiltlnFn) PackageHandler() {
	//注册然后分发
	//获取终端消息
	m.OmegaFrame.MainFrame.SetBackendCmdInterceptor(func(cmds []string) (stop bool) {
		m.PackageInjectIntoChan(cmds, definition.BACK_END_TYPE)
		return false
	})
	m.OmegaFrame.MainFrame.GetGameListener().SetGameChatInterceptor(func(chat *defines.GameChat) (stop bool) {
		m.PackageInjectIntoChan(chat, definition.MSG_TYPE)
		return false
	})
	//登进
	m.OmegaFrame.MainFrame.GetGameListener().AppendLoginInfoCallback(func(entry protocol.PlayerListEntry) {
		m.PackageInjectIntoChan(entry, definition.LOGIN_TYPE)
	})
	//登出
	m.OmegaFrame.MainFrame.GetGameListener().AppendLogoutInfoCallback(func(entry protocol.PlayerListEntry) {
		m.PackageInjectIntoChan(entry, definition.LOGOUT_TYPE)
	})
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
