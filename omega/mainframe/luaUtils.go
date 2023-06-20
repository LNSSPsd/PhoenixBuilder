package mainframe

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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

// file 结构体表示一个文件，包括文件名、文件内容和一个互斥锁（mutex），用于防止同时操作一个文件
type file struct {
	name    string
	content string
	mutex   sync.Mutex // 用于防止同时操作一个文件
}

// fileCenter 结构体表示文件处理中心，包括所有文件的集合和一个互斥锁，用于防止同时操作文件集合
type FileCenter struct {
	files sync.Map // 所有文件的集合，使用 sync.Map 来实现并发安全
}

// getFile 方法用于申请文件操作请求，参数包括文件名和回调函数。该方法首先获取或创建一个 file 对象，然后对该对象加锁，防止其他操作同时修改该文件。接着调用回调函数，将文件内容作为参数传入，获取回调函数的返回值。最后返回回调函数的返回值。
func (fc *FileCenter) GetFile(name string, callback func(io.Reader) (interface{}, error)) (interface{}, error) {
	value, _ := fc.files.LoadOrStore(name, &file{name: name, mutex: sync.Mutex{}}) // 获取或创建文件对象

	fileObj := value.(*file) // 将 sync.Map 中的对象转换为 file 类型

	fileObj.mutex.Lock() // 保证同一时间只有一个请求可以修改该文件
	defer fileObj.mutex.Unlock()
	reader := strings.NewReader(fileObj.content)

	return callback(reader) // 调用回调函数，将文件内容作为参数传入，获取回调函数的返回值，并返回该值
}

// writeFile 方法用于申请文件写入操作请求，参数包括文件名、待写入的内容和回调函数。该方法首先获取或创建一个 file 对象，然后对该对象加锁，防止其他操作同时修改该文件。接着调用回调函数，将文件作为参数传入，获取回调函数的返回值，并将待写入的内容作为参数传入 file 的 write 方法。最后返回回调函数的返回值。
func (fc *FileCenter) WriteFile(name string, content string) error {
	value, _ := fc.files.LoadOrStore(name, &file{name: name, mutex: sync.Mutex{}}) // 获取或创建文件对象

	fileObj := value.(*file) // 将 sync.Map 中的对象转换为 file 类型

	fileObj.mutex.Lock() // 保证同一时间只有一个请求可以修改该文件
	defer fileObj.mutex.Unlock()

	fileObj.write(content) // 安全地写入文件内容

	return nil // 调用回调函数，将文件作为参数传入，获取回调函数的返回值，并返回该值
}

// write 方法用于安全地写入文件内容，参数为待写入的内容。该方法使用互斥锁保证同一时间只有一个请求可以修改该文件。
func (f *file) write(content string) error {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	file, err := os.Create(f.name)
	if err != nil {
		fmt.Printf("打开文件 %s 失败： %v\n", f.name, err)
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		fmt.Printf("写入文件 %s 失败： %v\n", f.name, err)
		return err
	}

	f.content = content
	return nil
}

// 获取插件路径绝对路径 文件名字/插件名字
func (f *FileCenter) GetComponentPath() []string {
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
func (f *FileCenter) CheckFilePath() {
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
func (f *FileCenter) GetBindingJson() MappedBinding {
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
	maps, err := f.GetFile(bindingPath, func(filer io.Reader) (interface{}, error) {
		data, _ := ioutil.ReadAll(filer)
		var maps MappedBinding
		json.Unmarshal(data, &maps)
		return maps, nil
	})
	if err != nil {
		PrintInfo(NewPrintMsg("警告", err))
		return MappedBinding{}
	}
	if k, v := maps.(MappedBinding); v {
		return k
	}
	return MappedBinding{}

}

// 向binding写入绑定
func (f *FileCenter) WriteBindingJson(name string, path string) error {
	maps := f.GetBindingJson()
	if !f.CheckCompoentduplicates(name) {
		maps.Map[name] = path
		data, err := json.Marshal(maps)
		bindingPath := GetBindingPath()
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
		}

		f.WriteFile(bindingPath, string(data))

		return nil
	}
	return errors.New("已经有该名字的插件了 可以重新取名字 或者输入lua luas delect [插件名字]删除现有插件")

}

// 首先确定的是 配置在data/lua/config下 实现逻辑在data/lua/下 绑定它们的在data/lua/Binding.json
// 其次应该在程序一开始便开始检查这些
func (f *FileCenter) CheckCompoentduplicates(name string) bool {
	maps := f.GetBindingJson().Map
	if _, ok := maps[name]; ok {
		return true
	}
	return false
}

// 删除插件
func (f *FileCenter) DelectCompoent(name string) error {
	//检查插件是否存在
	maps := f.GetBindingJson().Map
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
