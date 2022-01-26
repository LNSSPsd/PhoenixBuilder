package main

import (
	"hash/adler32"
	"plugin_example/plugin_structs"
	"unsafe"
)

// func PluginInit(bridgeif unsafe.Pointer,mainref interface{}) string {
// 	bridge:=*(*plugin_structs.PluginBridge)(bridgeif)
// 	mainfunc:=mainref.(func(plugin_structs.PluginBridge)string)
// 	return mainfunc(bridge)
// }

type a interface{
	// aa()
}

func PluginInit(bridgeif unsafe.Pointer,mainref interface{}) string {
	// code where it called.
	// 只有plugin.go调用了这里.
	// name:=mainfunc.(func(unsafe.Pointer,interface{})string)(unsafe.Pointer(conn),mainref)

	// () operates after *,   so it deferences the unsafe.Pointer(it points the *Minecraft.conn and becomes deferencable.)
	// Now bridge is deferenced and bridgeif is transformed the interface called deferenced PluginBridge.
	// 此时PluginBridge (接口类型)被解引用, 且bridgeif这个Pointer约定为此接口, 最后解引用bridgeif, 以便传给mainfunc.
	bridge:=*(*plugin_structs.PluginBridge)(bridgeif)
	// assert to function with args and ret then call it.
	// 插件的Main(现在仍是接口类)随即被断言为具体类型的函数.
	mainfunc:=mainref.(func(plugin_structs.PluginBridge)string)
	return mainfunc(bridge)
}