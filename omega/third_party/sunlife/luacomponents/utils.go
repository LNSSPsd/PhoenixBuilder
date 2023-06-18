package luaComponent

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/pterm/pterm"
	lua "github.com/yuin/gopher-lua"
)

const (
	HEADLUA    = "luas"
	HEADRELOAD = "reload"
	HEADSTART  = "start"
)
const (
	BINDINGFILE = "Binding.json"
)

// 指令信息 必须遵循 HEAD BEHAVIOR
type CmdMsg struct {
	isCmd    bool
	Head     string
	Behavior string
	args     []string
}
type PrintMsg struct {
	Type string
	Body interface{}
}

// 绑定函数 "名字":"逻辑实现的文件名"
type MappedBinding struct {
	Map map[string]string `json:"绑定"`
}

// 打印指定消息
func printInfo(str PrintMsg) {
	pterm.Info.Printfln("[%v][%v]: %v ", time.Now().YearDay(), str.Type, str.Body)
}

// 构造一个输出函数
func newPrintMsg(typeName string, BodyString interface{}) PrintMsg {
	return PrintMsg{
		Type: typeName,
		Body: BodyString,
	}
}

// 检查各级目录是否完好
func checkFilePath() {
	rootPath := getRootPath() + SEPA + "lua"
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		os.MkdirAll(rootPath, os.ModePerm)
	}
	configPath := rootPath + SEPA + "config"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		os.MkdirAll(configPath, os.ModePerm)
	}
}

// 获取data的相对位置omega_storage\\data
func getRootPath() string {
	return OMGPATH
}

// 针对binding.json文件进行的各种包装
// 获取binding.json的路径
func getBindingPath() string {
	return getRootPath() + SEPA + "lua" + SEPA + BINDINGFILE
}

// 获取插件路径绝对路径 文件名字/插件名字
func getComponentPath() []string {
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

// 获取data/lua/config
func getConfigPath() string {
	return getRootPath() + SEPA + "lua" + SEPA + "config"
}

// 获取bindingJson内容
func getBindingJson() MappedBinding {
	bindingPath := getBindingPath() //不出意外就是data/lua/Binding.json
	file, err := os.Open(bindingPath)
	if err == nil {
		// 文件存在，解析 JSON 数据到结构体
		defer file.Close()

		data, _ := ioutil.ReadAll(file)
		var maps MappedBinding
		json.Unmarshal(data, &maps)
		return maps

	} else {
		bindingMap := MappedBinding{
			Map: make(map[string]string),
		}
		jsonData, err := json.MarshalIndent(bindingMap, "", "  ")
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
		}
		// 创建文件
		file, err := os.Create(bindingPath)
		if err != nil {
			fmt.Println("Error creating file:", err)
		}
		defer file.Close()

		// 将 JSON 数据写入文件
		_, err = file.Write(jsonData)
		if err != nil {
			fmt.Println("Error writing JSON to file:", err)
		}
		return bindingMap
	}
}

// 向binding写入绑定
func writeBindingJson(name string, path string) error {
	maps := getBindingJson()
	if !checkCompoentduplicates(name) {
		maps.Map[name] = path
		data, err := json.Marshal(maps)
		bindingPath := getBindingPath()
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
		}
		// 创建文件
		file, err := os.Create(bindingPath)
		if err != nil {
			fmt.Println("Error creating file:", err)
		}
		defer file.Close()

		// 将 JSON 数据写入文件
		_, err = file.Write(data)
		if err != nil {
			fmt.Println("Error writing JSON to file:", err)
		}
		return nil
	}
	return errors.New("已经有该名字的插件了 可以重新取名字 或者输入lua luas delect [插件名字]删除现有插件")

}

// 首先确定的是 配置在data/lua/config下 实现逻辑在data/lua/下 绑定它们的在data/lua/Binding.json
// 其次应该在程序一开始便开始检查这些
func checkCompoentduplicates(name string) bool {
	maps := getBindingJson().Map
	if _, ok := maps[name]; ok {
		return true
	}
	return false
}

// 删除插件
func delectCompoent(name string) error {
	//检查插件是否存在
	maps := getBindingJson().Map
	if _, ok := maps[name]; !ok {
		return errors.New(fmt.Sprintf("警告! 你正在想要删除%d但是我们并没有在插件中找到该名字的插件", name))
	}
	//先是在config中删除对应的文件
	configPath := getConfigPath() + SEPA + name + ".json"
	compentPath := maps[name]
	//删除文件关联中心

	delete(maps, name)
	data, _ := json.Marshal(maps)
	bindingPath := getBindingPath()
	err := os.MkdirAll(bindingPath, os.ModePerm)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(bindingPath, data, os.ModePerm)
	if err != nil {
		return err
	}
	//删除config
	err = delectFile(configPath)
	if err != nil {
		return err
	}
	//删除主要实现逻辑的
	err = delectFile(compentPath)
	if err != nil {
		return err
	}
	return nil
}

// 安全地删除指定文件
func delectFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	} else {
		// 文件存在，删除文件
		err := os.Remove(path)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

// 格式化处理指令
func formateCmd(str string) CmdMsg {

	words := strings.Fields(str)
	if len(words) < 3 {
		return CmdMsg{isCmd: false}
	}
	if words[0] != "lua" {
		return CmdMsg{isCmd: false}
	}
	head := words[1]
	//如果不属于任何指令则返回空cmdmsg
	if head != HEADLUA && head != HEADRELOAD && head != HEADSTART {
		return CmdMsg{isCmd: false}
	}
	behavior := words[2]
	args := []string{}
	if len(words) >= 3 {
		args = words[3:]
	}
	return CmdMsg{
		Head:     head,
		Behavior: behavior,
		args:     args,
		isCmd:    true,
	}
}

// 检查插件是否符合规范
func checkCompoent(l *lua.LState) error {
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
