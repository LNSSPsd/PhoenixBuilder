package main

import (
	_ "embed"
	"fmt"
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
	script:=string(testScript)
	identifyStr:= host.GetStringSha(script)
	host.InitHostFns(iso,global,hb,scriptName,identifyStr)
	ctx := v8.NewContext(iso, global)
	host.CtxFunctionInject(ctx)
	_, err := ctx.RunScript(script, scriptName)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("脚本已经执行完毕")
	c:=make(chan struct{})
	<-c
}

