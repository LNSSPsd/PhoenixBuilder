package BuiltlnFn

import lua "github.com/yuin/gopher-lua"

type BuiltDataControler struct {
	*BuiltlnFn
}

// 获取config 获取data信息 写入data信息 写入config信息
func (b *BuiltDataControler) BuiltFunc(L *lua.LState) int {
	DataControler := L.NewTable()
	L.Push(DataControler)
	return 1
}
