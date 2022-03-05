package host

import (
	"fmt"
	v8 "rogchap.com/v8go"
)

func AddPrint(iso *v8.Isolate) *v8.FunctionTemplate{
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		outStr:=""
		for _,a:=range info.Args(){
			outStr+=fmt.Sprintf("%v",a)
		}
		fmt.Printf(outStr)
		return nil
	})
	return printfn
}

func AddPrintln(iso *v8.Isolate) *v8.FunctionTemplate{
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		outStr:=""
		for _,a:=range info.Args(){
			outStr+=fmt.Sprintf("%v",a)
		}
		fmt.Println(outStr)
		return nil
	})
	return printfn
}

func AddBlock(iso *v8.Isolate) *v8.FunctionTemplate{
	printfn := v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {
		c:=make(chan struct{})
		<-c
		return nil
	})
	return printfn
}