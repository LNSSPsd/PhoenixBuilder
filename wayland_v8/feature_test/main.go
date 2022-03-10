package main

import (
	_ "embed"
	"fmt"
	"phoenixbuilder/fastbuilder/script"
	"phoenixbuilder/wayland_v8/host"

	v8 "rogchap.com/v8go"
)

//go:embed test.js
var testScript []byte
//go:embed test2.js
var testScript2 []byte
func main() {
	var err error
	hb:= script.NewHostBridge()
	iso := v8.NewIsolate()
	global := v8.NewObjectTemplate(iso)


	scriptName:="test.js"
	script:=string(testScript)
	identifyStr:= host.GetStringSha(script)
	host.InitHostFns(iso,global,hb,scriptName,identifyStr,"scriptPath")
	ctx := v8.NewContext(iso, global)
	host.CtxFunctionInject(ctx)
	_, err = ctx.RunScript(script, scriptName)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("脚本已经执行完毕")

	iso2 := v8.NewIsolate()
	global2 := v8.NewObjectTemplate(iso2)

	scriptName2:="test2.js"
	script2:=string(testScript2)
	identifyStr2:= host.GetStringSha(script2)
	host.InitHostFns(iso2,global2,hb,scriptName2,identifyStr2,"scriptPath")
	ctx2 := v8.NewContext(iso2, global2)
	host.CtxFunctionInject(ctx2)
	_, err = ctx2.RunScript(script2, scriptName2)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("脚本2已经执行完毕")
	c:=make(chan struct{})
	<-c
}

