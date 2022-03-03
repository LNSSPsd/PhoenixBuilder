package main

import (
	"fmt"
	"github.com/robertkrimen/otto"
)

func main() {

	vm:=otto.New()
	vm.Set("getCB",func(call otto.FunctionCall) otto.Value {
		cb,_:=vm.ToValue(func(cbCall otto.FunctionCall) otto.Value {
			fmt.Printf("Cb is called %v\n",cbCall)
			return otto.Value{}
		})
		return cb
	})
	var jsCB otto.Value
	vm.Set("webSocketConnect",func(call otto.FunctionCall) otto.Value {
		if !call.Argument(0).IsFunction(){
			panic("is not function")
		}
		jsCB=call.Argument(0)
		return otto.Value{}
	})
	vm.Run(`
cb=getCB()
cb("hello")


setCB(function (data) {
    console.log("js callback "+data+" is called")
})`)
	jsVal,_:=otto.ToValue("hi")
	jsCB.Call(otto.UndefinedValue(),jsVal)
}

