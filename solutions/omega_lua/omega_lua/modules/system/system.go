package system

import (
	"phoenixbuilder/solutions/omega_lua/omega_lua/pollers"
	"time"

	lua "github.com/yuin/gopher-lua"
)

type SystemModule struct {
	LuaGoSystem
	pollers.LuaAsyncInvoker
}

type LuaGoSystem interface {
	Print(string)
	UserInputChan() <-chan string
}

// 获得一个集成了go实现部分 和lua交互实现部分的结构体
func NewSystemModule(goImplements LuaGoSystem, luaAsyncInvoker pollers.LuaAsyncInvoker) *SystemModule {
	return &SystemModule{
		//继承luaGoSystem
		LuaGoSystem: goImplements,
		//继承luaAsyncInvoker
		LuaAsyncInvoker: luaAsyncInvoker,
	}
}

func (m *SystemModule) MakeLValue(L *lua.LState) (lua.LValue, map[lua.LValue]pollers.LuaEventDataChanMaker) {
	//创建一个表
	luaModule := L.NewTable()
	//获取现在开始时间
	startTime := float64(time.Now().UnixMilli()) / 1000
	//向内注入函数
	luaModule = L.SetFuncs(luaModule, map[string]lua.LGFunction{
		"print": m.luaGoSystemPrint,
		"os":    m.luaGoSystemOs,
		"cwd":   m.luaGoSystemCwd,
		"now": func(l *lua.LState) int {
			l.Push(lua.LNumber((float64(time.Now().UnixMilli()) / 1000) - startTime))
			return 1
		},
	})
	//看起来像是特殊的实现
	// poller flags for sleep and input
	flagSleep := L.NewFunction(m.luaGoSleep)
	flagInput := L.NewFunction(m.luaGoInput)
	//看上去像是函数以及对应的资源管理器
	//这个资源管理器仿佛是负责内部环境资源协调的
	pollerFlags := map[lua.LValue]pollers.LuaEventDataChanMaker{
		flagSleep: goSleepSourceMaker,
		flagInput: m.goInputSourceMaker,
	}
	// inject block_sleep and block_input flags into module
	luaModule.RawSetString("block_sleep", flagSleep)
	luaModule.RawSetString("block_input", flagInput)
	// inject start_time into module
	luaModule.RawSetString("start_time", lua.LNumber(startTime))
	return luaModule, pollerFlags
}
