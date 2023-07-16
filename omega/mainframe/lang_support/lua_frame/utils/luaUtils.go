package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"phoenixbuilder/omega/mainframe/lang_support/lua_frame/definition"
	"phoenixbuilder/omega/mainframe/lang_support/lua_frame/luaConfig"
	"runtime"
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
	IsCmd    bool
	Head     string
	Behavior string
	Args     []string
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

// 获取omega_storage位置
func GetRootPath() string {
	if runtime.GOOS == "android" {
		return filepath.Join(GetAndroidDownPath(), definition.OMGSTIRAGEPATH)
	}
	return definition.OMGSTIRAGEPATH
}

// 获取安卓的下载目录
func GetAndroidDownPath() string {
	downloadDir := filepath.Join(os.Getenv("EXTERNAL_STORAGE"), "Download")
	return downloadDir
}

// 获取"omega_storage\\data"
func GetDataPath() string {
	return filepath.Join(GetRootPath(), "data")
}

// "omega_storage\\配置"
func GetOmgConfigPath() string {
	return filepath.Join(GetRootPath(), "配置")
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
		return CmdMsg{IsCmd: false}
	}
	if words[0] != "lua" {
		return CmdMsg{IsCmd: false}
	}
	head := words[1]
	//如果不属于任何指令则返回空cmdmsg
	if head != HEADLUA && head != HEADRELOAD && head != HEADSTART {
		return CmdMsg{IsCmd: false}
	}
	behavior := words[2]
	args := []string{}
	if len(words) >= 3 {
		args = words[3:]
	}
	return CmdMsg{
		Head:     head,
		Behavior: behavior,
		Args:     args,
		IsCmd:    true,
	}
}

/*

文件管理系统

*/

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
	// 使用 os.Open 打开文件。
	file, err := os.Open(filename)
	if err != nil {
		return []byte{}, err
	}
	defer file.Close()

	// 获取文件信息，以确定要读取的字节数。
	_, err = file.Stat()
	if err != nil {
		return []byte{}, err
	}

	// 使用 ioutil.ReadAll 从文件中读取内容。
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return content, err
	}
	return content, nil
}

// 读取并返回结构体
func (f *FileControl) ReadConfig(path string) (luaConfig.LuaCommpoentConfig, error) {
	newConfig := luaConfig.LuaCommpoentConfig{
		Disabled: true, //默认关闭
	}
	data, err := f.Read(path)
	if err != nil {
		return newConfig, err
	}

	err = json.Unmarshal(data, &newConfig)
	if err != nil {
		return newConfig, err
	}

	return newConfig, nil
}

// 删除插件
func (f *FileControl) DelectCompoentFile(name string) error {
	f.DeleteSubDir(name)
	//关闭相关内容
	PrintInfo(NewPrintMsg("提示", fmt.Sprintf("%v已经删除 干净了", name)))

	return nil
}

// deleteSubDir 函数接受一个父目录路径和一个子目录名称作为参数，
// 并安全地删除指定的子目录及其所有子文件。
func (f *FileControl) DeleteSubDir(subDirName string) error {
	parentDir := GetOmgConfigPath()
	subDir := filepath.Join(parentDir, subDirName)
	// 检查子目录是否存在。
	if !f.fileExists(subDir) {
		return nil
	}

	// 删除子目录及其所有子文件。
	err := os.RemoveAll(subDir)
	if err != nil {
		return err
	}

	return nil
}

// Result 结构体用于存储 JSON 文件和 Lua 文件的路径。
type Result struct {
	JsonFile   string
	LuaFile    string
	JsonConfig luaConfig.LuaCommpoentConfig
}

// GetLuaComponentPath返回一个包含同名字 JSON 文件和 Lua 文件路径的字典。
func (f *FileControl) GetLuaComponentData() (map[string]Result, error) {
	dir := GetOmgConfigPath()
	results := make(map[string]Result)

	// 使用 filepath.Walk 遍历指定目录及其子目录。
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 如果当前路径是一个目录，则检查是否存在与目录名同名的 JSON 和 Lua 文件。
		if info.IsDir() {
			dirName := info.Name()

			jsonFile := filepath.Join(path, dirName+".json")
			luaFile := filepath.Join(path, dirName+".lua")

			// 如果找到 JSON 和 Lua 文件，将它们的路径添加到结果字典中。
			if f.fileExists(jsonFile) && f.fileExists(luaFile) {
				//读取json文件
				config, err := f.ReadConfig(jsonFile)
				if err != nil {
					PrintInfo(NewPrintMsg("警告", err))
				}
				results[dirName] = Result{
					JsonFile:   jsonFile,
					LuaFile:    luaFile,
					JsonConfig: config,
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return results, nil
}

type LuaResult struct {
	Config luaConfig.LuaCommpoentConfig
	Code   []byte
}

// 获取指定配置和code
func (f *FileControl) GetConfigAndCode(name string) (LuaResult, error) {
	Path := filepath.Join(GetOmgConfigPath(), name)
	//检查是否存在
	info, err := os.Stat(Path)
	if err != nil {
		return LuaResult{}, errors.New("不存在该目录")
	}
	//检查是否为目录
	if !info.IsDir() {
		return LuaResult{}, errors.New("不应该在配置文件里面存在一个lua插件名字的文件")
	}
	//获取路径
	jsonPath := filepath.Join(Path, name+".json")
	LuaPath := filepath.Join(Path, name+".lua")
	jsonConfig, err := f.ReadConfig(jsonPath)
	if err != nil {
		return LuaResult{}, err
	}
	luaCode, err := f.Read(LuaPath)
	if err != nil {
		return LuaResult{}, err
	}
	return LuaResult{
		jsonConfig,
		luaCode,
	}, nil
}

// fileExists 函数接受一个文件路径作为参数，如果文件存在则返回 true，否则返回 false。
func (f *FileControl) fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// 在指定目录下创建具有指定名称的子目录，并在子目录中创建同名的 JSON 和 Lua 文件。
func (f *FileControl) CreateDirAndFiles(name string) error {
	// 创建子目录。
	dir := GetOmgConfigPath()
	data := luaConfig.LuaCommpoentConfig{
		Name:     name,
		Usage:    "",
		Version:  "0.0.1",
		Source:   "Lua-Component",
		Disabled: true,
		Author:   "",
		Config:   make(map[string]interface{}),
	}
	luaCode := `--根据注释初步了解如何书写代码
	--gameCtrol = skynet.GetControl()初始化操作机器人游戏行为的
	--gameCtrol.SendWsCmd("/say hellow") 发送指令
	--gameCtrol.SendCmdAndInvokeOnResponse("/say hellow") 发送指令并且返回一个表 内有是否成功 和返回信息两个值
	--h =gameCtrol.SendCmdAndInvokeOnResponse("/say hellow")
	--print(h.Success,"是否成功")
	--print(h.outputmsg,"输出信息")
	--listener = skynet.GetListener()
	--MsgListener = listener.GetMsgListner()
	--while true do
	--    print(MsgListener:NextMsg())   nextMsg可以读取玩家说话 如果玩家没有说话 那么就会堵塞直到玩家说话
	--end
	
	`

	subDir := filepath.Join(dir, name)
	// 检查目录是否已经存在，如果存在则返回错误。
	if _, err := os.Stat(subDir); !os.IsNotExist(err) {
		return errors.New("该名字的lua组件已经存在")
	}
	err := os.Mkdir(subDir, 0755)
	if err != nil {
		return err
	}

	// 创建 JSON 文件并根据指定结构体进行初始化。
	jsonFile := filepath.Join(subDir, name+".json")
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(jsonFile, jsonData, 0644)
	if err != nil {
		return err
	}

	// 创建 Lua 文件并根据指定字符串进行初始化。
	luaFile := filepath.Join(subDir, name+".lua")
	err = ioutil.WriteFile(luaFile, []byte(luaCode), 0644)
	if err != nil {
		return err
	}

	return nil
}

// RemoveAt 从列表中删除指定下标的元素，返回删除后的列表和删除的元素
func RemoveSlice(list []interface{}, index int) ([]interface{}, interface{}) {
	if index < 0 || index >= len(list) {
		return list, nil
	}
	value := list[index]
	return append(list[:index], list[index+1:]...), value
}
