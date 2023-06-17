package luaComponent

import (
	"fmt"

	"github.com/pterm/pterm"
	lua "github.com/yuin/gopher-lua"
)

//关于lua全生命周期的活动都将由此

type Monitor struct {
	L *lua.LState
}

// 开始加载程序 返回通讯通道
func (m *Monitor) Start() {
	// 创建一个新的 Lua 虚拟机实例
	L := lua.NewState()
	defer L.Close()

	// 为 Lua 虚拟机提供一个安全的环境
	L.PreloadModule("communicator", registerCommunicatorModule)

	// 加载 Lua 脚本文件
	err := L.DoFile("lua/center.lua")
	if err != nil {
		fmt.Println("Error loading script:", err)
		return
	}
	//获取对象
	pterm.Info.Println("成功启动lua交互器")
	m.L = L

}

// lua插件状况
func (m *Monitor) IsIn() {

}

// 停止lua交互器
func (m *Monitor) Stop() {

}
