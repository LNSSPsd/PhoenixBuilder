package BuiltlnFn

import lua "github.com/yuin/gopher-lua"

type BuiltlnFunctionDic map[string]BuiltlnFunctioner
type BuiltlnFunctioner interface {
	BuiltFunc(L *lua.LState) int
}

// 获取内置函数 方便注入
func (b *BuiltlnFn) GetSkynetBuiltlnFunction() BuiltlnFunctionDic {
	return map[string]BuiltlnFunctioner{
		"GetListener":   &BuiltListener{b},
		"GetControl":    &BuiltGameControler{b},
		"loadComponent": &LoadSide{b},
	}
}
