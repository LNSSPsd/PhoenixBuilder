package main

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
)

type Runnable interface {
	// start running until an error occour or complete
	// if no error, final result will be returned in a json string
	Run() (string,error)
}

type OttoKeeper interface {
	// try to compile script and then insert host(golang) func
	LoadNewScript(script string,name string) (Runnable,error)
	// actually is a pack of VM.Set, but each Runnable will be automatically set this
	SetInitFn(func(vm *otto.Otto,name string))
}

type RunnableAlpha struct {
	Name string
	VM             *otto.Otto
	Script string
	OnResultCallback func(Result string,err error)
}

func (runnable *RunnableAlpha) Run() (string,error) {
	finalVal, err := runnable.VM.Run(runnable.Script)
	Errorf:=func(fmtStr string,a ...interface{}) error {
		return fmt.Errorf("JS-Script(%v): "+fmtStr,runnable.Name,a)
	}
	if err != nil {
		return "", Errorf("Runtime Error (%v)",err)
	}else{
		err := runnable.VM.Set("finalVal", finalVal)
		if err != nil {
			return "", Errorf("cannot set final result (%v)",err)
		}
		jsonResult, err :=runnable.VM.Run("JSON.stringify(finalVal)")
		if err != nil {
			return "", Errorf("cannot stringify final result (%v)",err)
		}
		strResult, err := jsonResult.ToString()
		if err != nil {
			return"", Errorf("cannot get final result (%v)",err)
		}
		return strResult,nil
	}
}

func (Runnable *RunnableAlpha) RunInRoutine(){
	go func() {
		result, err := Runnable.Run()
		if err != nil {
			Runnable.OnResultCallback("",fmt.Errorf("RuntimeError"))
		}else{
			Runnable.OnResultCallback(result,err)
		}
	}()
}

type OttoKeeperAlpha struct {
	initFn func(vm *otto.Otto,name string)
}

func (oa *OttoKeeperAlpha)LoadNewScript(script string,name string) Runnable {
	vm:=otto.New()
	oa.initFn(vm,name)
	return &RunnableAlpha{Name: name,VM: vm,Script:script,OnResultCallback: func(Result string, err error) {}}
}

func MakeWaitConnect() (func(),func(call otto.FunctionCall) otto.Value){
	initC:=make(chan struct{})
	return func() {
		close(initC)
	},
	func(call otto.FunctionCall) otto.Value {
		<-initC
		return otto.Value{}
	}
}

func MakeGeneralUserInput()(chan string,func(fbCmd otto.FunctionCall) otto.Value){
	cmdChan:=make(chan string)
	return cmdChan, func(call otto.FunctionCall) otto.Value {
		fbCmd, _ := call.Argument(0).ToString()
		cmdChan<-fbCmd
		return otto.Value{}
	}
}

func main() {
	ottoHostWaitConnect,ottoVMWaitConnect:=MakeWaitConnect()
	initFn:=func(vm *otto.Otto,name string) {
		err := vm.Set("FB_WaitConnect", ottoVMWaitConnect)
		if err!=nil{
			panic(err)
		}
	}
	ottoKeeper:=&OttoKeeperAlpha{initFn}
	file, err := os.OpenFile("in.js",os.O_RDONLY,0755)
	if err != nil {
		panic(err)
	}
	all, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}
	script := ottoKeeper.LoadNewScript(string(all), "测试脚本")
 	finalResult, err :=script.Run()
	 if err!=nil{
		 panic(err)
	 }
	 fmt.Println(finalResult)
}

