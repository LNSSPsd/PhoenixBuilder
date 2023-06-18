package luaComponent

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
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

// 描述了一个lua插件该有的东西
type LuaCommpoentConfig struct {
	Name     string                 `json:"插件名字"`
	Usage    string                 `json:"插件用途"`
	Disabled bool                   `json:"是否禁用"`
	Version  string                 `json:"版本号"`
	Author   string                 `json:"作者"`
	Config   map[string]interface{} `json:"配置"`
}

// 开始加载程序 返回通讯通道
func (m *Monitor) Start() {
	m.ComponentPoll = m.RegistrationPlugins()
	println(m.ComponentPoll)

}

// 读取 并且检查插件 返回插件列表
func (m *Monitor) RegistrationPlugins() map[string]*LuaComponent {
	paths := m.getComponentPath()
	pool := make(map[string]*LuaComponent)
	for i := 0; i < len(paths); i++ {
		path := paths[i]
		L := lua.NewState()
		// 为 Lua 虚拟机提供一个安全的环境 提供基础的方法
		L.PreloadModule("communicator", m.registerCommunicatorModule)
		// 加载 Lua 脚本文件

		_, err := L.LoadFile(path)
		if err != nil {
			fmt.Println("Error loading script:", err, "插件:", path)
		}
		if err := m.checkCompoent(L); err != nil {
			pool[path] = &LuaComponent{
				L:       L,
				Msg:     make(map[string]string),
				Running: false,
				Config:  LuaCommpoentConfig{},
			}
		} else {
			m.printInfo("插件不符合规范")
		}

	}
	return pool
}

// 打印指定消息
func (m *Monitor) printInfo(str string) {
	pterm.Info.Printfln("[lua插件] %d", str)
}

// 检查插件是否符合规范
func (m *Monitor) checkCompoent(l *lua.LState) error {
	// 检查函数是否存在Init函数 active函数 以及选填的[save函数] [getData函数]
	ComponentFn := []string{
		COMPONENT_INIT_FN,
		COMPONENT_ACTIVE_FN,
	}
	for _, v := range ComponentFn {
		if l.GetGlobal(v) == lua.LNil {
			return errors.New("错误 该插件不含有" + v + "函数")
		}
	}
	return nil

	//errors.New("错误")
}

// 读取运行并且返回运行插件列表
func (m *Monitor) startLuaCompoent() {

}

// 获取插件路径 文件名字/插件名字
func (m *Monitor) getComponentPath() []string {
	nameList := []string{}
	dirPath := OMGPATH + SEPA + "lua"
	fileExt := ".lua"
	// 读取目录下的所有文件
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 如果目录不存在，则创建它
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			panic(err)
		}
		fmt.Println("Directory created:", dirPath)
	}
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		panic(err)
	}

	// 遍历目录下的所有文件名
	for _, file := range files {
		// 如果文件后缀名为 .lua，则打印文件名（去掉后缀名）
		if strings.HasSuffix(file.Name(), fileExt) {
			nameList = append(nameList, file.Name())
		}

	}
	return nameList
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

}
