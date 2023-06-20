package mainframe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/pterm/pterm"
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
type FileControl struct {
	//文件锁
	FileLock *FileLock
}

// 文件锁类型
type FileLock struct {
	mu sync.RWMutex
}

// 获取文件锁
func (lock *FileLock) Lock() {
	lock.mu.Lock()
}

// 释放文件锁
func (lock *FileLock) Unlock() {
	lock.mu.Unlock()
}

// 获取文件读锁
func (lock *FileLock) RLock() {
	lock.mu.RLock()
}

// 释放文件读锁
func (lock *FileLock) RUnlock() {
	lock.mu.RUnlock()
}

// 创建一个新的文件锁
func NewFileLock() *FileLock {
	return &FileLock{}
}

// 安全写入文件
func (f *FileControl) Write(filename string, data []byte) error {
	// 获取写锁
	lock := f.FileLock
	lock.Lock()
	defer lock.Unlock()

	// 写入数据
	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return err
	}

	return nil
}

// 安全读取文件
func (f *FileControl) Read(filename string) ([]byte, error) {
	// 获取读锁
	lock := f.FileLock
	lock.RLock()
	defer lock.RUnlock()

	// 读取数据
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// 获取插件路径绝对路径 文件名字/插件名字
func GetComponentPath() []string {
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

// 检查各级目录是否完好
func (f *FileControl) CheckFilePath() {
	rootPath := GetRootPath() + SEPA + "lua"
	if _, err := os.Stat(rootPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		os.MkdirAll(rootPath, os.ModePerm)
	}
	configPath := rootPath + SEPA + "config"
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		os.MkdirAll(configPath, os.ModePerm)
	}
	f.GetBindingJson()

}

// 打印指定消息
func PrintInfo(str PrintMsg) {
	pterm.Info.Printfln("[%v][%v]: %v ", time.Now().YearDay(), str.Type, str.Body)
}

// 构造一个输出函数
func NewPrintMsg(typeName string, BodyString interface{}) PrintMsg {
	return PrintMsg{
		Type: typeName,
		Body: BodyString,
	}
}

// 获取data的相对位置omega_storage\\data
func GetRootPath() string {
	return OMGPATH
}

// 针对binding.json文件进行的各种包装
// 获取binding.json的路径
func GetBindingPath() string {
	return GetRootPath() + SEPA + "lua" + SEPA + BINDINGFILE
}

// 获取data/lua/config
func GetConfigPath() string {
	return GetRootPath() + SEPA + "lua" + SEPA + "config"
}

// 获取bindingJson内容
func (f *FileControl) GetBindingJson() MappedBinding {
	bindingPath := GetBindingPath() //不出意外就是data/lua/Binding.json
	_, err := os.Stat(bindingPath)
	if err != nil {
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
	data, err := f.Read(bindingPath)
	if err != nil {
		PrintInfo(NewPrintMsg("警告", err))
		return MappedBinding{}
	}
	var maps MappedBinding
	json.Unmarshal(data, &maps)
	return maps

}

// 向binding写入绑定
func (f *FileControl) WriteBindingJson(name string, path string) error {
	maps := f.GetBindingJson()
	if !f.CheckCompoentduplicates(name) {
		maps.Map[name] = path
		data, err := json.Marshal(maps)
		bindingPath := GetBindingPath()
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
		}
		//写入json
		f.Write(bindingPath, data)

		return nil
	}
	return errors.New("已经有该名字的插件了 可以重新取名字 或者输入lua luas delect [插件名字]删除现有插件")

}

// 首先确定的是 配置在data/lua/config下 实现逻辑在data/lua/下 绑定它们的在data/lua/Binding.json
// 其次应该在程序一开始便开始检查这些
func (f *FileControl) CheckCompoentduplicates(name string) bool {
	maps := f.GetBindingJson().Map
	if _, ok := maps[name]; ok {
		return true
	}
	return false
}

// 删除插件
func (f *FileControl) DelectCompoent(name string) error {
	//检查插件是否存在
	maps := f.GetBindingJson().Map
	PrintInfo(NewPrintMsg("数据", maps))
	if _, ok := maps[name]; !ok {
		return errors.New(fmt.Sprintf("警告! 你正在想要删除%d但是我们并没有在插件中找到该名字的插件", name))
	}
	//先是在config中删除对应的文件
	configPath := GetConfigPath() + SEPA + name + ".json"
	compentPath := maps[name]
	//删除文件关联中心
	delete(maps, name)
	data, _ := json.Marshal(maps)
	bindingPath := GetBindingPath()
	file, err := os.Create(bindingPath)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(data)
	//删除config
	err = DelectFile(configPath)
	if err != nil {
		return err
	}
	//删除主要实现逻辑的
	err = DelectFile(compentPath)
	if err != nil {
		return err
	}
	PrintInfo(NewPrintMsg("提示", fmt.Sprintf("%v已经删除 干净了", name)))
	return nil
}

// 安全地删除指定文件
func DelectFile(path string) error {
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
func FormateCmd(str string) CmdMsg {

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

/*
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
*/
