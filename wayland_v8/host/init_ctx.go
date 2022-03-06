package host

import (
	"encoding/json"
	"fmt"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"phoenixbuilder/minecraft/protocol/packet"
	"rogchap.com/v8go"
	"strings"
)

func AllowPath(path string) bool {
	if strings.Contains(path,"fbtoken"){
		return false
	}
	if strings.Contains(path,"fb_script_permission"){
		return false
	}
	return true
}

func LoadPermission(hb HostBridge,identifyStr string) map[string]bool{
	permission:= map[string]bool{}
	fullPermission:=map[string]map[string]bool{}
	file, err := hb.LoadFile("fb_script_permission.json")
	if err!= nil {
		return permission
	}
	err = json.Unmarshal([]byte(file), &fullPermission)
	if err != nil {
		return permission
	}
	if savedPermission,ok:=fullPermission[identifyStr];ok{
		return savedPermission
	}
	return permission
}

func SavePermission(hb HostBridge,identifyStr string,permission map[string]bool){
	fullPermission:=map[string]map[string]bool{}
	file, err := hb.LoadFile("fb_script_permission.json")
	dataToSave:=[]byte{}
	if err== nil {
		json.Unmarshal([]byte(file), &fullPermission)
	}
	fullPermission[identifyStr]=permission
	dataToSave, _ = json.Marshal(fullPermission)
	hb.SaveFile("fb_script_permission.json",string(dataToSave))
}

func InitHostFns(iso *v8go.Isolate,global *v8go.ObjectTemplate,hb HostBridge,_scriptName string,identifyStr string) {
	scriptName:=_scriptName
	permission:=LoadPermission(hb,identifyStr)
	updatePermission:= func() {
		SavePermission(hb,identifyStr,permission)
	}

	throwException:= func(funcName string,str string) *v8go.Value {
		value, _ := v8go.NewValue(iso, "Crashed In Host Function ["+funcName+"], because: "+str)
		iso.ThrowException(value)
		return nil
	}
	printException:= func(funcName string,str string) *v8go.Value {
		fmt.Println("Crashed In Host Function ["+funcName+"], because: "+str)
		return nil
	}
	throwNotConnectedException:= func(funcName string) *v8go.Value  {
		return throwException(funcName,"connection to MC not established")
	}
	hasStrIn := func(info *v8go.FunctionCallbackInfo,pos int,argName string) (string,bool) {
		if len(info.Args())<pos+1{
			return fmt.Sprintf("no arg %v provided in pos %v",argName,pos),false
		}
		if !info.Args()[pos].IsString(){
			return fmt.Sprintf("arg %v in pos %v is not a string (you set: %v)",argName,pos,info.Args()[pos].String()),false
		}
		return info.Args()[pos].String(),true
	}
	hasFuncIn := func(info *v8go.FunctionCallbackInfo,pos int,argName string) (string,*v8go.Function) {
		if len(info.Args())<pos+1{
			return fmt.Sprintf("no arg %v provided in pos %v",argName,pos),nil
		}
		function, err := info.Args()[pos].AsFunction()
		if err != nil {
			return fmt.Sprintf("arg %v in pos %v is not a function (you set: %v)",argName,pos,info.Args()[pos].String()),nil
		}
		return "",function
	}
	t:=&Terminator{
		c:           make(chan struct{}),
		isTeminated: false,
		TerminateHook: make([]func(),0),
	}

	// function FB_SetName(scriptName string) None
	if err:=global.Set("FB_SetName",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"FB_SetName[scriptName]"); !ok{
				throwException("FB_SetName",str)
			}else{
				hb.Println("脚本["+scriptName+"]正在将自己命名为["+str+"]",t,scriptName)
				scriptName=str
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_WaitConnect() None
	if err:=global.Set("FB_WaitConnect",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			hb.WaitConnect(t)
			return nil
		})); err!=nil{panic(err)}

	// fuunction FB_WaitConnectAsync(cb func()) None
	if err:=global.Set("FB_WaitConnectAsync",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if errStr,cbFn:=hasFuncIn(info,0,"FB_WaitConnectAsync[cb]"); cbFn==nil{
				throwException("FB_WaitConnectAsync",errStr)
			}else{
				go func() {
					hb.WaitConnect(t)
					cbFn.Call(info.This())
				}()

			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_Println(msg string) None
	if err:=global.Set("FB_Println",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"FB_Println[msg]"); !ok{
				throwException("FB_Println",str)
			}else{
				hb.Println(str,t,scriptName)
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_GeneralCmd(fbCmd string) None
	if err:=global.Set("FB_GeneralCmd",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected(){
				throwNotConnectedException("FB_GeneralCmd")
			}
			if str,ok:= hasStrIn(info,0,"FB_GeneralCmd[fbCmd]"); !ok{
				throwException("FB_GeneralCmd",str)
			}else{
				hb.FBCmd(str,t)
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_SendMCCmd(mcCmd string) None
	if err:=global.Set("FB_SendMCCmd",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected(){
				throwNotConnectedException("FB_SendMCCmd")
			}
			if str,ok:= hasStrIn(info,0,"FB_SendMCCmd[mcCmd]"); !ok{
				throwException("FB_SendMCCmd",str)
			}else{
				hb.MCCmd(str,t,false)
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_SendMCCmdAndGetResult(mcCmd string) jsObject
	if err:=global.Set("FB_SendMCCmdAndGetResult",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected(){
				throwNotConnectedException("FB_SendMCCmdAndGetResult")
			}
			if str,ok:= hasStrIn(info,0,"FB_SendMCCmdAndGetResult[mcCmd]"); !ok{
				throwException("FB_SendMCCmdAndGetResult",str)
			}else{
				pk:= hb.MCCmd(str,t,true)
				strPk, err := json.Marshal(pk)
				if err != nil {
					return throwException("FB_SendMCCmdAndGetResult","Cannot convert host packet to Json Str: "+str)
				}
				value, err := v8go.JSONParse(info.Context(), string(strPk))
				if err != nil {
					return throwException("FB_SendMCCmdAndGetResult",str)
				}else{
					return value
				}
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_SendMCCmdAndGetResultAsync(mcCmd string, onResult func(jsObject)) None
	// jsObject=null, if cannot get result in callback
	if err:=global.Set("FB_SendMCCmdAndGetResultAsync",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected(){
				throwNotConnectedException("FB_SendMCCmdAndGetResultAsync")
			}
			if str,ok:= hasStrIn(info,0,"FB_SendMCCmdAndGetResultAsync[mcCmd]"); !ok{
				throwException("FB_SendMCCmdAndGetResultAsync",str)
			}else{
				if errStr,cbFn:=hasFuncIn(info,1,"FB_SendMCCmdAndGetResultAsync[onResult]"); cbFn==nil{
					throwException("FB_WaitConnectAsync",errStr)
				}else{
					ctx:=info.Context()
					go func() {
						pk:= hb.MCCmd(str,t,true)
						strPk, err := json.Marshal(pk)
						if err != nil {
							printException("FB_SendMCCmdAndGetResultAsync","Cannot convert host packet to Json Str: "+str)
							cbFn.Call(info.This(),v8go.Null(iso))
							return
						}
						val, err := v8go.JSONParse(ctx,string(strPk))
						if err != nil {
							printException("FB_SendMCCmdAndGetResultAsync","Cannot Parse Json Packet in Host: "+str)
							cbFn.Call(info.This(),v8go.Null(iso))
							return
						}else {
							cbFn.Call(info.This(),val)
						}
					}()
				}
				return nil
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_RequireUserInput(hint string) string
	if err:=global.Set("FB_RequireUserInput",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"FB_RequireUserInput[hint]"); !ok{
				throwException("FB_RequireUserInput",str)
			}else{
				userInput:= hb.GetInput(str,t,scriptName)
				value,_:=v8go.NewValue(iso,userInput)
				return value
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_RequireUserInputAsync(hint,onInput func(string)) None
	if err:=global.Set("FB_RequireUserInputAsync",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"FB_RequireUserInputAsync[hint]"); !ok{
				throwException("FB_RequireUserInputAsync",str)
			}else{
				if errStr,cbFn:=hasFuncIn(info,1,"FB_RequireUserInputAsync[onInput]"); cbFn==nil{
					throwException("FB_RequireUserInputAsync",errStr)
				}else{
					go func() {
						userInput:= hb.GetInput(str,t,scriptName)
						value,_:=v8go.NewValue(iso,userInput)
						cbFn.Call(info.This(),value)
					}()
				}

			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_RegPacketCallBack(packetType,onPacket func(jsObject)) deRegFn
	// when deRegFn is called, onPacket function will no longer be called
	if err:=global.Set("FB_RegPackCallBack",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"FB_RegPackCallBack[packetType]"); !ok{
				throwException("FB_RegPackCallBack",str)
			}else{
				if errStr,cbFn:=hasFuncIn(info,1,"FB_RegPackCallBack[onPacket]"); cbFn==nil{
					throwException("FB_RegPackCallBack",errStr)
				}else{
					ctx:=info.Context()
					deRegFn, err := hb.RegPacketCallBack(str, func(pk packet.Packet) {
						strPk, err := json.Marshal(pk)
						if err!=nil{
							printException("FB_RegPackCallBack","Cannot convert host packet to Json Str: "+err.Error())
							cbFn.Call(info.This(),v8go.Null(iso))
						}else{
							val, err := v8go.JSONParse(ctx,string(strPk))
							if err != nil {
								printException("FB_RegPackCallBack","Cannot Parse Json Packet in Host: "+str)
								cbFn.Call(info.This(),v8go.Null(iso))
								return
							}else {
								cbFn.Call(info.This(),val)
							}
						}
					},t)
					if err != nil {
						return throwException("FB_RegPackCallBack",err.Error())
					}
					jsCbFn:=v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
						deRegFn()
						return nil
					})
					t.TerminateHook=append(t.TerminateHook,deRegFn)
					return jsCbFn.GetFunction(ctx).Value
				}
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_RegChat(onMsg func(name,msg)) deRegFn
	if err:=global.Set("FB_RegChat",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if errStr,cbFn:=hasFuncIn(info,0,"FB_RegChat[onMsg]"); cbFn==nil{
				throwException("FB_RegChat",errStr)
			}else{
				ctx:=info.Context()
				deRegFn, err := hb.RegPacketCallBack("IDText", func(pk packet.Packet) {
					p := pk.(*packet.Text)
					SourceName, err := v8go.NewValue(iso,p.SourceName)
					if err != nil {
						printException("FB_RegChat",err.Error())
						cbFn.Call(info.This(),v8go.Null(iso),v8go.Null(iso))
						return
					}
					Message, err := v8go.NewValue(iso,p.Message)
					if err != nil {
						printException("FB_RegChat",err.Error())
						cbFn.Call(info.This(),v8go.Null(iso),v8go.Null(iso))
						return
					}
					cbFn.Call(info.This(),SourceName,Message)
				},t)
				if err != nil {
					return throwException("FB_RegChat",err.Error())
				}
				jsCbFn:=v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
					deRegFn()
					return nil
				})
				t.TerminateHook=append(t.TerminateHook,deRegFn)
				return jsCbFn.GetFunction(ctx).Value
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_Query(info string) string
	if err:=global.Set("FB_Query",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"FB_Query[info]"); !ok{
				throwException("FB_Query",str)
			}else{
				if str=="script_sha"{
					value,_:=v8go.NewValue(iso,identifyStr)
					return value
				}
				userInput:= hb.Query(str)
				value,_:=v8go.NewValue(iso,userInput)
				return value
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_GetAbsPath(path string) string
	if err:=global.Set("FB_GetAbsPath",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"FB_GetAbsPath[path]"); !ok{
				throwException("FB_GetAbsPath",str)
			}else{
				absPath:=hb.GetAbsPath(str)
				value,_:=v8go.NewValue(iso,absPath)
				return value
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_RequireFilePermission(hint,path) isSuccess
	if err:=global.Set("FB_RequireFilePermission",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if hint,ok:= hasStrIn(info,1,"FB_RequireFilePermission[hint]"); !ok{
				throwException("FB_RequireUserInput",hint)
				return nil
			}else{
				if dir,ok:= hasStrIn(info,0,"FB_RequireFilePermission[hint]"); !ok{
					throwException("FB_RequireUserInput",dir)
					return nil
				}else{
					dir=hb.GetAbsPath(dir)
					permissionKey:="VisitDir:"+dir
					if hasPermission,ok:= permission[permissionKey];ok&&hasPermission{
						value,_ :=v8go.NewValue(iso,true)
						return value
					}else{
						for{
							warning:="脚本["+scriptName+"]["+_scriptName+"]想访问文件夹"+dir+"的所有内容\n"+
								"理由是:"+hint+"\n"+
								"(警告，恶意脚本可能会删除，篡改，威胁这个文件夹下的所有文件！)\n"+
								"是否允许? 输入[是/否/Y/y/N/n]:"
							choose:=hb.GetInput(warning,t,scriptName)
							if choose=="是" || choose=="Y" || choose=="y"{
								value,_ :=v8go.NewValue(iso,true)
								permission[permissionKey]=true
								updatePermission()
								return value
							}else if choose=="否" || choose =="N"|| choose=="n"{
								value,_ :=v8go.NewValue(iso,false)
								return value
							}
							hb.Println("无效输入，请输入[是/否/Y/y/N/n]其中之一",t,scriptName)
						}
					}
				}
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_ReadFile(path) string
	// if permission is not granted or read fail, "" is returned
	if err:=global.Set("FB_ReadFile",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"FB_ReadFile[path]"); !ok{
				throwException("FB_ReadFile",str)
			}else{
				p:=hb.GetAbsPath(str)

			}
			return nil
		})); err!=nil{panic(err)}

	// fetch
	if err := fetch.InjectTo(iso, global); err != nil {
		panic(err)
	}

	// websocket

	// special function here

	// function FB_ScriptCrash(string reason) None
	if err:=global.Set("FB_ScriptCrash",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {

			if str,ok:= hasStrIn(info,0,"FB_ScriptCrash[reson]"); !ok{
				throwException("FB_ScriptCrash",str)
			}else{
				throwException("FB_ScriptCrash",str)
				t.Terminate()
			}
			return nil
		})); err!=nil{panic(err)}
	//// FB_Block
	//if err := global.Set("FB_Block",
	//	v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
	//		hb.Block(t)
	//		return nil
	//	})); err!=nil{
	//	panic(err)
	//}
}