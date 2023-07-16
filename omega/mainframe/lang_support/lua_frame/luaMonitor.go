package luaFrame

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"phoenixbuilder/omega/mainframe/lang_support/lua_frame/BuiltlnFn"
	"phoenixbuilder/omega/mainframe/lang_support/lua_frame/LuaUtils"
	"phoenixbuilder/omega/mainframe/lang_support/lua_frame/definition"
	"phoenixbuilder/omega/mainframe/lang_support/lua_frame/luaConfig"
	omgApi "phoenixbuilder/omega/mainframe/lang_support/lua_frame/omgcomponentapi"
	"phoenixbuilder/solutions/omega_lua/monk"
	"phoenixbuilder/solutions/omega_lua/omega_lua/concurrent"
	"phoenixbuilder/solutions/omega_lua/omega_lua/modules/command"
	"phoenixbuilder/solutions/omega_lua/omega_lua/modules/listen"
	"phoenixbuilder/solutions/omega_lua/omega_lua/modules/packets_utils"
	"phoenixbuilder/solutions/omega_lua/omega_lua/modules/system"
	submodule_holder "phoenixbuilder/solutions/omega_lua/omega_lua/modules_holder"
	"sync"

	lua "github.com/yuin/gopher-lua"
)

// 插件监测器
type Monitor struct {
	//每个插件拥有自己的lua运行环境 并且每个插件的名字都将是这个插件唯一的指示标志
	//在运行的初期就会初始化所有的插件 并且根据产生的配置文件决定是否开启 这与omg普通插件没有区别
	//区别点在于lua的优势导致 这个插件能够热重载以及能够修改其中的主要逻辑
	ComponentPoll map[string]*LuaComponent
	//omg框架
	LuaComponentData map[string]LuaUtils.Result
	OmgFrame         *omgApi.OmgApi
	FileControl      *LuaUtils.FileControl
	BuiltlnFner      *BuiltlnFn.BuiltlnFn
	Enver            *Env
}
type Env struct {
	L *lua.LState
}

// 插件
type LuaComponent struct {
	L *lua.LState
	//排队中的消息
	Msg map[string]string
	//是否运行
	Running bool
	//插件的配置
	Config luaConfig.LuaCommpoentConfig
}

func NewMonitor(lc *omgApi.OmgApi) *Monitor {
	return &Monitor{
		ComponentPoll: make(map[string]*LuaComponent),
		//获取omg框架
		OmgFrame: lc,
		BuiltlnFner: &BuiltlnFn.BuiltlnFn{
			OmegaFrame:       lc,
			Listener:         sync.Map{},
			PackageChanSlice: []*definition.PackageChan{},
		},
		FileControl: &LuaUtils.FileControl{},
		Enver:       &Env{},
	}

}

// 创建lua环境 内置函数 上下文管理器 刷新环境
func (m *Monitor) CreateLuaEnv(ctx context.Context) (ac concurrent.AsyncCtrl) {
	L := lua.NewState()
	//获得一个上下文管理器
	ac = concurrent.NewAsyncCtrl(ctx)
	m.Enver.L = L
	goSystem := monk.NewMonkSystem()
	goPackets := monk.NewMonkPackets(128)
	goCmdSender := monk.NewMonkCmdSender()
	// lua wrapper
	systemModule := system.NewSystemModule(goSystem, ac)
	luaSystemModule, systemPollerFlags := systemModule.MakeLValue(L)
	packetsModule := packets_utils.NewOmegaPacketsModule(goPackets)
	luaPacketsModule := packetsModule.MakeLValue(L)
	cmdModule := command.NewCmdModule(goCmdSender, packetsModule.NewGamePacket)
	luaCmdModule := cmdModule.MakeLValue(L, ac)

	// pollers
	ListenModule := listen.NewListenModule(ac,
		goPackets, packetsModule.NewGamePacket,
		systemPollerFlags)
	luaListenModule := ListenModule.MakeLValue(L)

	// load modules
	L.PreloadModule("omega", submodule_holder.NewSubModuleHolder(map[string]lua.LValue{
		"system":  luaSystemModule,
		"listen":  luaListenModule,
		"packets": luaPacketsModule,
		"cmds":    luaCmdModule,
	}).Loader)
	return ac
}

// 开启所有可运行插件 并且监听所有错误
func (m *Monitor) StartAllComponent(ctx concurrent.AsyncCtrl, CodeAndConfigs map[string]omgApi.CodeAndConfig) error {
	errChans := []<-chan error{}
	for _, v := range CodeAndConfigs {
		errchan := concurrent.FireLuaCodeInGoRoutine(ctx, m.Enver.L, string(v.Code))
		errChans = append(errChans, errchan)
	}
	//如果遇到错误
	go func() {
		for _, v := range errChans {
			select {
			case err := <-v:
				panic(err)
			default:
				continue
			}
		}
	}()
	return nil
}

// 接受指令处理并且执行
func (m *Monitor) CmdCenter(msg string) error {

	CmdMsg := LuaUtils.FormateCmd(msg)
	if !CmdMsg.IsCmd {
		return errors.New(fmt.Sprintf("很显然%v并不是指令的任何一种 请输入lua luas help寻求帮助", msg))
	}

	switch CmdMsg.Head {
	case LuaUtils.HEADLUA:
		//lua指令
		if err := m.luaCmdHandler(&CmdMsg); err != nil {
			LuaUtils.PrintInfo(LuaUtils.NewPrintMsg("警告", err))
		}
		/*case HEADRELOAD:
		go func() {
			if err := m.Reload(&CmdMsg); err != nil {
				PrintInfo(NewPrintMsg("警告", err))
			}
		}()*/
		/*
			case HEADSTART:
				go func() {
					if err := m.StartCmdHandler(&CmdMsg); err != nil {
						PrintInfo(NewPrintMsg("警告", err))
					}
				}()
		*/
	}
	return nil
}

/*
// 插件行为 重加载某个插件 如果参数为all则全部插件重加载 记住reload和startComponent是有区别的
// reload是再次扫描对应的插件然后默认不开启 而startCompent是直接在插件池子里面开启插件
func (m *Monitor) Reload(cmdmsg *LuaUtils.CmdMsg) error {

	switch cmdmsg.Behavior {
	case "component":
		args := cmdmsg.Args
		if len(args) != 1 {
			return errors.New("lua reload compoent指令后面应该有且仅有一个参数")
		}
		componentName := args[0]
		//检查

		if args[0] == "all" {
			nameDic, err := m.FileControl.GetLuaComponentData()
			if err != nil {
				return err
			}
			//开启所有插件
			for name, _ := range nameDic {
				if newErr := m.RunComponent(name); newErr != nil {
					return newErr
				}
			}
			return nil
		}
		//单独开启
		if NewErr := m.RunComponent(componentName); NewErr != nil {
			return NewErr
		}
		LuaUtils.PrintInfo(LuaUtils.NewPrintMsg("提示", fmt.Sprintf("%v已经重新加载", componentName)))
		return nil
	default:
		LuaUtils.PrintInfo(LuaUtils.NewPrintMsg("警告", "无效指令"))

	}
	return nil
}
*/
/*
// 处理cmd
func (m *Monitor) StartCmdHandler(CmdMsg *CmdMsg) error {
	args := CmdMsg.args
	switch CmdMsg.Behavior {
	case "component":
		if len(args) != 1 {
			return errors.New("lua start compoent指令后面应该有且仅有一个参数")
		}
		if args[0] == "all" {
			//to do
			PrintInfo(NewPrintMsg("提示", fmt.Sprintf("全部插件已经开启")))
			return nil
		}
		// to do (修改)
		//componentName := args[0]
		//if err := m.Run(componentName); err != nil {
		//PrintInfo(NewPrintMsg("警告", err))
		//} else {
		//PrintInfo(NewPrintMsg("提示", fmt.Sprintf("%v插件已经开启", componentName)))
		//}

	default:
		PrintInfo(NewPrintMsg("警告", "这不是一个合理的指令"))
	}
	return nil
}*/
// 安全地关闭组件并且从配置文件中删除
/*
func (m *Monitor) CloseLua(name string) error {

	if v, ok := m.ComponentPoll[name]; ok {
		//如果为nil
		if v.L == nil {
			delete(m.ComponentPoll, name)
			return nil
		}
		v.L.Close()
		v.L = nil
		v.Running = false
		delete(m.ComponentPoll, name)
		return nil
	}

	return nil
}*/

// lua指令类执行
func (m *Monitor) luaCmdHandler(CmdMsg *LuaUtils.CmdMsg) error {
	args := CmdMsg.Args
	switch CmdMsg.Behavior {
	case "help":
		warning := []string{
			"lua luas help 寻求指令帮助\n",
			//"lua reload component [重加载的插件名字] 加载/重加载指定插件 如果参数是all就是全部插件重载\n",
			"lua luas new [新插件名字] [描述]创建一个自定义空白插件[描述为选填]\n",
			"lua luas delect [插件名字]\n",
			"lua luas list 列出当前正在运行的插件\n",
			//"lua luas stop [插件名字] 暂停插件运行 参数为all则暂停所有插件运行",
		}
		msg := ""
		for _, v := range warning {
			msg += v
		}
		LuaUtils.PrintInfo(LuaUtils.NewPrintMsg("提示", msg))
	case "new":
		//参数检查
		if len(args) != 1 && len(args) != 2 {
			return errors.New("lua luas new后面应该加上[插件名字]或者说[插件名字]")
		}
		componentName := args[0]
		/*
			componentUsage := ""
			if len(args) == 2 {
				componentUsage = args[1]
			}*/
		//检查当前是否有
		if _, ok := m.LuaComponentData[componentName]; ok {
			return errors.New(fmt.Sprintf("已经含有%v插件 无法创立", componentName))
		}
		//如果没有则创建文件
		if err := m.FileControl.CreateDirAndFiles(componentName); err != nil {
			return err
		}
		LuaUtils.PrintInfo(LuaUtils.NewPrintMsg("提示", fmt.Sprintf("已经创建文件基本结构请到目录%v 修改", filepath.Join(LuaUtils.GetOmgConfigPath(), componentName))))

	case "delect":
		if len(args) != 1 {
			return errors.New("lua luas delect指令后面应该加上需要删除的插件名字")
		}
		//从运行中删除
		//m.CloseLua(args[0])
		//文件删除
		m.FileControl.DelectCompoentFile(args[0]) //DelectCompoent(args[0])
	case "list":
		msg := ""
		for k, v := range m.ComponentPoll {
			if v.Running {
				msg += fmt.Sprintf("[%v]", k)
			}
		}
		LuaUtils.PrintInfo(LuaUtils.NewPrintMsg("信息", msg+"处于开启状态"))
	/*case "stop":
	if len(args) != 1 {
		return errors.New("lua luas stop指令后面应该加上需要删除的插件名字")
	}
	name := args[0]
	if name == "all" {
		for k, _ := range m.ComponentPoll {
			m.CloseLua(k)
			PrintInfo(NewPrintMsg("提示", fmt.Sprintf("%v插件关闭成功", k)))
		}
		PrintInfo(NewPrintMsg("提示", "全部组件已经关闭"))
		return nil
	}
	if _, ok := m.ComponentPoll[name]; !ok {
		return errors.New(fmt.Sprintf("我们并没有在加载的插件池子中找到%v", name))
	}
	m.CloseLua(name)
	PrintInfo(NewPrintMsg("提示", fmt.Sprintf("%v插件关闭成功", name)))
	*/
	default:
		return errors.New("未知指令 请输入lua luas help寻求帮助")
	}
	return nil
}

/*
// 安全地启动插件
func (m *Monitor) RunComponent(name string) error {
	//安全关闭
	if v, ok := m.ComponentPoll[name]; ok {
		if v.L == nil {
			LuaUtils.PrintInfo(LuaUtils.NewPrintMsg("警告", fmt.Sprintf("%v存在于运行池子中 但是lua解释器为nil", name)))
			//删除
			delete(m.ComponentPoll, name)
		}
		//判断是否存在某个方法
		method := reflect.ValueOf(v.L).MethodByName("DoString")
		if !method.IsValid() {
			LuaUtils.PrintInfo(LuaUtils.NewPrintMsg("警告", fmt.Sprintf("%v检测如果调用loadfille会触发错误", name)))
			delete(m.ComponentPoll, name)
		}
		//如果启动状态则关闭删除
		v.L.Close()
		delete(m.ComponentPoll, name)

	}
	//获取配置
	data, err := m.FileControl.GetConfigAndCode(name)
	if err != nil {
		return err
	}
	//判断配置是否开启
	if data.Config.Disabled {
		m.OmgFrame.Omega.GetBackendDisplay().Write(pterm.Warning.Sprintf("\t跳过加载组件  [%v] %v@%v", data.Config.Source, name, data.Config.Version))
		return nil
	}
	m.OmgFrame.Omega.GetBackendDisplay().Write(pterm.Success.Sprintf("\t正在加载组件 [%v] %v@%v", data.Config.Source, name, data.Config.Version))

	//另外开线程
	go func(newName string) {
		//如果没有
		L := lua.NewState()
		// 为 Lua 虚拟机提供一个安全的环境 提供基础的方法
		configString, ConfigErr := json.Marshal(data.Config)
		if ConfigErr != nil {
			LuaUtils.PrintInfo(LuaUtils.NewPrintMsg("警告", ConfigErr))
		}
		L.SetGlobal("ComponentConfig", lua.LString(configString))
		if err := m.BuiltlnFner.LoadFn(L); err != nil {
			fmt.Println(err)
			return
		}
		m.ComponentPoll[name] = &LuaComponent{
			L:       L,
			Msg:     make(map[string]string),
			Running: false, //初始化完成但是未运行
			Config:  data.Config,
		}
		defer m.CloseLua(newName)

		if _, ok := m.ComponentPoll[name]; !ok {
			LuaUtils.PrintInfo(LuaUtils.NewPrintMsg("警告", fmt.Sprintf("%v不存在该组件名字", name)))
			return
		}
		if err := m.ComponentPoll[name].L.DoString(string(data.Code)); err != nil {
			LuaUtils.PrintInfo(LuaUtils.NewPrintMsg("lua代码报错", err))
		}
	}(name)
	return nil
}*/
