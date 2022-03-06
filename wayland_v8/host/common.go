package host

import (
	"fmt"
	v8 "rogchap.com/v8go"
)

func AddWebSocket(iso *v8.Isolate) *v8.FunctionTemplate{
	return v8.NewFunctionTemplate(iso, func(info *v8.FunctionCallbackInfo) *v8.Value {

	}
}