package BuiltlnFn

import (
	"fmt"
	lua "github.com/yuin/gopher-lua"
	"path/filepath"
	"phoenixbuilder/omega/mainframe/lang_support/lua_frame/LuaUtils"
)

type BuiltFileControler struct {
	*BuiltlnFn
}

func (b *BuiltFileControler) BuiltFunc(L *lua.LState) int {
	fileControler := L.NewTable()
	L.SetField(fileControler, "GetData", L.NewFunction(b.GetData))
	L.SetField(fileControler, "GetConfigPath", L.NewFunction(func(l *lua.LState) int {
		l.Push(lua.LString(LuaUtils.GetOmgConfigPath()))
		return 1
	}))
	//获取配置文件的坐标位置
	L.SetField(fileControler, "GetDataPath", L.NewFunction(func(l *lua.LState) int {
		l.Push(lua.LString(filepath.Join(LuaUtils.GetRootPath(), "data")))
		return 1
	}))

	L.Push(fileControler)
	return 1
}

// 获取配置文件 返回string
func (b *BuiltFileControler) GetData(l *lua.LState) int {
	if l.GetTop() != 1 {
		l.ArgError(1, "请传入文件名字来获取指定信息")
	}
	filePath := l.CheckString(1)
	if data, err := b.OmegaFrame.MainFrame.GetFileData(filePath); err != nil {
		l.ArgError(1, fmt.Sprintf("%v", err))
		l.Push(lua.LString(""))
		return 1
	} else {
		l.Push(lua.LString(string(data)))
	}
	return 1
}
