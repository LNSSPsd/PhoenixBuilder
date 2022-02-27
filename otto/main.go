package main

import (
	"fmt"
	"github.com/robertkrimen/otto"
)

type Runnable interface {
	// start running until an error occour or complete
	// if no error, final result will be returned in a json string
	Run() (string,error)
}

type OttoKeeper interface {
	// try to compile script and then insert host(golang) func
	LoadNewScript(script string,name string) (Runnable,error)
	// actually is a pack of vm.Set, but each Runnable will be automatically set this
	SetHostEnv(name string, ottoValue interface{})
}

type RunnableAlpha struct {
	name string
	vm *otto.Otto
	compiledScript string
}

func (runnable *RunnableAlpha) Run() (string,error) {
	finalVal, err := runnable.vm.Run(runnable.compiledScript)
	Errorf:=func(fmtStr string,a ...interface{}) error {
		return fmt.Errorf("JS-Script(%v): "+fmtStr,runnable.name,a)
	}
	if err != nil {
		return "", Errorf("Runtime Error (%v)",err)
	}else{
		err := runnable.vm.Set("finalVal", finalVal)
		if err != nil {
			return "", Errorf("cannot set final result (%v)",err)
		}
		jsonResult, err :=runnable.vm.Run("JSON.stringify(finalVal)")
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

type OttoKeeperAlpha struct {

}

func LoadNewScript(script string,name string) (Runnable,error){

}

func main() {

}

