package main

import (
	_ "embed"
	"phoenixbuilder/wayland_v8/host"

	v8 "rogchap.com/v8go"
)

//go:embed test.js
var testScript []byte

func main() {
	iso := v8.NewIsolate()
	global := v8.NewObjectTemplate(iso)

	hb:= host.NewHostBridge()
	scriptName:="test.js"
	host.InitHostFns(iso,global,hb,scriptName)
	ctx := v8.NewContext(iso, global)
	ctx.RunScript(string(testScript), scriptName)

	c:=make(chan struct{})
	<-c
}

