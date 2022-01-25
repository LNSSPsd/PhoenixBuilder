package main

import (
	"plugin_example/plugin_structs"
	"unsafe"
)

// func PluginInit(bridgeif unsafe.Pointer,mainref interface{}) string {
// 	bridge:=*(*plugin_structs.PluginBridge)(bridgeif)
// 	mainfunc:=mainref.(func(plugin_structs.PluginBridge)string)
// 	return mainfunc(bridge)
// }

func PluginInit(conn unsafe.Pointer,mainref interface{}) string {
	bridge:=*(*plugin_structs.PluginBridge)(conn)
	// assert to function with args and ret then call it.
	mainfunc:=mainref.(func(plugin_structs.PluginBridge)string)
	return mainfunc(bridge)
}