package luaComponent

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	lua "github.com/yuin/gopher-lua"
)

const (
	COMPONENT_INIT_FN   = "init"
	COMPONENT_ACTIVE_FN = "active"
	COMPONENT__SAVE_FN  = "save"
	OMGPATH             = "omega_storage" + SEPA + "data"
	SEPA                = string(filepath.Separator)
)

type Monitor struct {
	L       *lua.LState
	running bool     //是否运行
	msg     []string //传入消息列表
	//每个插件拥有自己的lua运行环境 并且每个插件的名字都将是这个插件唯一的指示标志
	//在运行的初期就会初始化所有的插件 并且根据产生的配置文件决定是否开启 这与omg普通插件没有区别
	//区别点在于lua的优势导致 这个插件能够热重载以及能够修改其中的主要逻辑
	ComponentPoll map[string]*LuaComponent
}
type LuaComponent struct {
	L *lua.LState
	//排队中的消息
	Msg map[string]string
	//是否运行
	Running bool
	//插件的配置
	Config LuaCommpoentConfig
}

// 开始加载程序 返回通讯通道
func (m *Monitor) Start() {
	//检查lua插件所需要的所有目录结构
	checkFilePath()
	err := m.RegistrationPlugins()

	if err != nil {
		printInfo(newPrintMsg("警告", err))
	}
	println(m.ComponentPoll)
}

// 单独加载某个插件
func (m *Monitor) Load(name string) error {
	bindingMap := getBindingJson().Map
	path := bindingMap[name]
	//检查是否已经存在
	if m.ComponentPoll == nil {
		m.ComponentPoll = make(map[string]*LuaComponent)
	}
	if _, ok := m.ComponentPoll[name]; ok {
		m.ComponentPoll[name].L.Close()
		delete(m.ComponentPoll, name)
	}
	L := lua.NewState()
	// 为 Lua 虚拟机提供一个安全的环境 提供基础的方法
	L.PreloadModule("communicator", m.registerCommunicatorModule)
	// 加载 Lua 脚本文件

	_, err := L.LoadFile(path)
	if err != nil {
		fmt.Println("Error loading script:", err, "插件:", path)
	}
	if err := checkCompoent(L); err != nil {
		m.ComponentPoll[name] = &LuaComponent{
			L:       L,
			Msg:     make(map[string]string),
			Running: false,
			Config:  LuaCommpoentConfig{},
		}
	} else {

		printInfo(newPrintMsg("警告", "插件不符合规范"+err.Error()))
	}
	return nil
}

// 读取 并且检查插件 返回插件列表
func (m *Monitor) RegistrationPlugins() error {
	names := getBindingJson().Map
	for k, _ := range names {
		m.Load(k)
	}
	return nil
}

// lua插件状况
func (m *Monitor) IsIn() {

}

// 停止lua交互器
func (m *Monitor) Stop() {

}

// 注册一个名为 "communicator" 的 Lua 模块，提供双向通信功能
func (m *Monitor) registerCommunicatorModule(L *lua.LState) int {
	communicator := map[string]lua.LGFunction{
		"send_message": sendMessage,
		"get_message":  getMessage,
	}

	mod := L.SetFuncs(L.NewTable(), communicator)
	L.Push(mod)
	return 1
}

//应该有几个基础方法也就是:sendcmd(cmd)直接发送指令到游戏内
//SetOnParamMsg(name ,fnName)绑定玩家与对应处理函数
//CustomResourcePack(packageid,fnName)绑定自定义包
//SendCmdAndInvokeOnResponse(cmd,fnName)发送指令并且绑定处理函数
//getNewMsg(fn)

// sendMessage 是一个 Go 函数，用于从 Lua 脚本发送消息
// 每次lua发送消息后便移交给信息中心处理
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

// 传入json信息 并且将消息分类后 按照类别分别移交给插件注册中心 包绑定中心 行为逻辑中心
func (m *Monitor) HandleMsg(msg string) {
	cmdmsg := formateCmd(msg)
	if !cmdmsg.isCmd {

	}
}

// 接受指令处理并且执行
func (m *Monitor) CmdCenter(msg string) error {

	CmdMsg := formateCmd(msg)
	if !CmdMsg.isCmd {
		return errors.New(fmt.Sprintf("很显然%v并不是指令的任何一种 请输入lua luas help寻求帮助", msg))
	}

	switch CmdMsg.Head {
	case HEADLUA:
		//lua指令
		if err := m.luaCmdHandler(&CmdMsg); err != nil {
			printInfo(newPrintMsg("警告", err))
		}
	case HEADRELOAD:
		if err := m.Reload(&CmdMsg); err != nil {
			printInfo(newPrintMsg("警告", err))
		}
	case HEADSTART:
		if err := m.StartCmdHandler(&CmdMsg); err != nil {
			printInfo(newPrintMsg("警告", err))
		}
	}
	return nil
}

// 插件行为 重加载某个插件 如果参数为all则全部插件重加载 记住reload和startComponent是有区别的
// reload是再次扫描对应的插件然后默认不开启 而startCompent是直接在插件池子里面开启插件
func (m *Monitor) Reload(cmdmsg *CmdMsg) error {
	switch cmdmsg.Behavior {
	case "component":
		args := cmdmsg.args
		if len(args) != 1 {
			return errors.New("lua reload compoent指令后面应该有且仅有一个参数")
		}
		componentName := args[0]
		if args[0] == "all" {
			//关闭插件
			for _, v := range m.ComponentPoll {
				v.L.Close()
			}
			//读取新的插件
			for _, v := range getComponentPath() {
				if err := m.Load(v); err != nil {
					printInfo(newPrintMsg("警告", err))
				}
			}
		}
		if err := m.Load(componentName); err != nil {
			return err
		}
	}
	return nil
}
func (m *Monitor) StartCmdHandler(CmdMsg *CmdMsg) error {
	args := CmdMsg.args
	switch CmdMsg.Behavior {
	case "component":
		if len(args) != 1 {
			return errors.New("lua start compoent指令后面应该有且仅有一个参数")
		}
		componentName := args[0]
		m.Run(componentName)

	}
	return nil
}

// 启动已有插件
func (m *Monitor) Run(name string) error {
	if _, ok := m.ComponentPoll[name]; !ok {
		return errors.New("我们并没有在当前的插件池中找到该名字的插件 请你确定有该插件 或者说请尝试重加载一次插件:lua reload component all")
	}
	maps := getBindingJson().Map
	// 调用 Lua 函数
	go m.ComponentPoll[name].L.DoFile(maps[name])
	m.ComponentPoll[name].Running = true
	printInfo(newPrintMsg("启动", fmt.Sprintf("%d插件启动成功 ", name)))
	return nil
}

// lua指令类执行
func (m *Monitor) luaCmdHandler(CmdMsg *CmdMsg) error {
	args := CmdMsg.args
	switch CmdMsg.Behavior {
	case "help":
		warning := []string{
			"lua luas help 寻求指令帮助\n",
			"lua reload component [重加载的插件名字] 加载/重加载指定插件 如果参数是all就是全部插件重载\n",
			"lua start component [需要开启的插件名字] 开启插件 参数为all则开启所有插件\n",
			"lua luas new [新插件名字] [描述]创建一个自定义空白插件[描述为选填]\n",
			"lua luas delect [插件名字]\n",
		}
		msg := ""
		for _, v := range warning {
			msg += v
		}
		printInfo(newPrintMsg("提示", msg))
	case "new":
		//to do
		if len(args) != 1 && len(args) != 2 {
			return errors.New("lua luas new后面应该加上[插件名字]或者说[插件名字] [用途]")
		}
		componentName := args[0]
		componentUsage := ""
		if len(args) == 2 {
			componentUsage = args[1]
		}
		//检查是否重合
		BindingMaps := getBindingJson().Map
		if _, ok := BindingMaps[componentName]; ok {
			return errors.New("已经存在同样名字的插件了 请重新命名 或者说删除原有的插件lua luas delect [插件名字]")
		}
		if err := m.newComponent(componentName, componentUsage); err != nil {
			printInfo(newPrintMsg("警告", err))
		}

	case "delect":
		if len(args) != 1 {
			return errors.New("lua luas delect指令后面应该加上需要删除的插件名字")
		}
		delectCompoent(args[0])
	default:
		return errors.New("未知指令 请输入lua luas help寻求帮助")
	}
	return nil
}

// 创建一个插件
func (m *Monitor) newComponent(componentName string, componentUsage string) error {
	//to do
	//开始创建config
	luaConfig := LuaCommpoentConfig{
		Name:     componentName,
		Usage:    componentUsage,
		Disabled: true,
		Version:  "0.0.1",
		Author:   "author",
		Config:   make(map[string]interface{}),
	}

	// 将 JSON 数据写入文件
	getConfigPath()
	filePath := OMGPATH + SEPA + "lua" + SEPA + "config" + SEPA + componentName + ".json" // 替换为实际文件路径

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 文件不存在，创建新文件

		jsonData, err := json.Marshal(luaConfig)
		if err != nil {
			// 处理错误
			return err
		}
		file, err := os.Create(filePath)
		if err != nil {
			// 处理错误
			fmt.Println(err)
		}

		_, err = file.Write(jsonData)
		if err != nil {
			// 处理错误
			return err
		}
		file.Close()
	}

	//创建逻辑区域
	// 指定目录和文件名
	dir := OMGPATH + SEPA + "lua" + SEPA
	filename := componentName + ".lua"

	// 创建文件的完整路径
	filepath := filepath.Join(dir, filename)

	// 创建文件
	file, errs := os.Create(filepath)
	if errs != nil {
		fmt.Println("Error creating file:", errs)
		return errs
	}
	file.Close()

	fmt.Printf("File %s created in %s\n", filename, dir)
	//绑定写入
	writeBindingJson(componentName, filepath)

	printInfo(newPrintMsg("提示", "创建完成 请重加载组件"))
	return nil
}
