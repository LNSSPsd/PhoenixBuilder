// +build with_v8

package script_kickstarter

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"phoenixbuilder/fastbuilder/script"
	"phoenixbuilder/wayland_v8/host"
	v8 "rogchap.com/v8go"
	"strings"
)

func LoadScript(scriptPath string, hb script.HostBridge) (func(),error) {
	iso := v8.NewIsolate()
	global := v8.NewObjectTemplate(iso)
	scriptPath = strings.TrimSpace(scriptPath)
	if scriptPath == "" {
		return nil,fmt.Errorf("Empty script path!")
	}
	fmt.Printf("Loading script: %s", scriptPath)
	fmt.Printf("JS engine vesion: %v\n",host.JSVERSION)
	_, scriptName := path.Split(scriptPath)
	file, err := os.OpenFile(scriptPath, os.O_RDONLY, 0755)
	if err != nil {
		return nil,err
	}
	scriptData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil,err
	}
	script:=string(scriptData)
	identifyStr:= ""//script.GetStringSha(script)
	stopFunc:=host.InitHostFns(iso,global,hb,scriptName,identifyStr,scriptPath)
	ctx := v8.NewContext(iso, global)
	host.CtxFunctionInject(ctx)
	go func() {
		finalVal, err := ctx.RunScript(script, scriptName)
		if err != nil {
			fmt.Printf("Script %s ran into a runtime error: %v\n",scriptPath,err.Error())
		}
		fmt.Printf("Script %s completed: %v\n",scriptPath,finalVal)
	}()
	return stopFunc,nil
}
