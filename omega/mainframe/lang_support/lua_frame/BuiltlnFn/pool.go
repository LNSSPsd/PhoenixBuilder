package BuiltlnFn

import lua "github.com/yuin/gopher-lua"

type BuiltlnFunctionDic map[string]func(L *lua.LState) int

// 获取内置函数
func (b *BuiltlnFn) GetSkynetBuiltlnFunction() BuiltlnFunctionDic {
	return map[string]func(L *lua.LState) int{
		"GetListener": b.BuiltlnListner,
		"GetControl":  b.BuiltGameContrler,
	}
}
