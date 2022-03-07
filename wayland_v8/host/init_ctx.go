package host

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"go.kuoruan.net/v8go-polyfills/base64"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"go.kuoruan.net/v8go-polyfills/timers"
	"go.kuoruan.net/v8go-polyfills/url"
	"os"
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

func InitHostFns(iso *v8go.Isolate,global *v8go.ObjectTemplate,hb HostBridge,_scriptName string,identifyStr string,scriptPath string) func() {
	scriptName:=_scriptName
	permission:=LoadPermission(hb,identifyStr)
	updatePermission:= func() {
		SavePermission(hb,identifyStr,permission)
	}

	throwException:= func(funcName string,str string) *v8go.Value {
		value, _ := v8go.NewValue(iso, "脚本崩溃于函数 ["+funcName+"], 原因是: "+str)
		iso.ThrowException(value)
		return nil
	}
	printException:= func(funcName string,str string) *v8go.Value {
		fmt.Println("脚本在函数 ["+funcName+"] 出现错误, 原因是: "+str)
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
		isTerminated: false,
		TerminateHook: make([]func(),0),
	}
	t.TerminateHook=append(t.TerminateHook, func() {
		iso.TerminateExecution()
	})

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
				}else if str=="script_path"{
					value,_:=v8go.NewValue(iso,scriptPath)
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
					dir=hb.GetAbsPath(dir)+string(os.PathSeparator)
					if !AllowPath(dir){
						throwException("FB_ReadFile","脚本正在试图访问禁止访问的路径，可能为恶意脚本！")
						t.Terminate()
						return nil
					}
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

	// function FB_ReadFile(path string) string
	// if permission is not granted or read fail, "" is returned
	if err:=global.Set("FB_ReadFile",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"FB_ReadFile[path]"); !ok{
				throwException("FB_ReadFile",str)
			}else{
				p:=hb.GetAbsPath(str)
				hasPermission:=false
				for permissionName,_:=range permission{
					if strings.HasPrefix(permissionName,"VisitDir:"){
						if strings.HasPrefix(p,permissionName[len("VisitDir:"):]){
							hasPermission=true
							break
						}
					}
				}
				if !hasPermission{
					throwException("FB_ReadFile","脚本正在试图访问无权限的路径，可能为恶意脚本！")
					t.Terminate()
					return nil
				}
				if !AllowPath(p){
					throwException("FB_ReadFile","脚本正在试图访问禁止访问的路径，可能为恶意脚本！")
					t.Terminate()
					return nil
				}
				data, err := hb.LoadFile(p)
				if err != nil {
					value,_:=v8go.NewValue(iso,"")
					return value
				}
				value,_:=v8go.NewValue(iso,data)
				return value
			}
			return nil
		})); err!=nil{panic(err)}

	// function FB_SaveFile(path string,data string) isSuccess
	if err:=global.Set("FB_SaveFile",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if p,ok:= hasStrIn(info,0,"FB_SaveFile[path]"); !ok{
				throwException("FB_SaveFile",p)
			}else{
				if data,ok:= hasStrIn(info,1,"FB_SaveFile[data]"); !ok{
					throwException("FB_SaveFile",data)
				}else{
					p:=hb.GetAbsPath(p)
					hasPermission:=false
					for permissionName,_:=range permission{
						if strings.HasPrefix(permissionName,"VisitDir:"){
							if strings.HasPrefix(p,permissionName[len("VisitDir:"):]){
								hasPermission=true
								break
							}
						}
					}
					if !hasPermission{
						throwException("FB_ReadFile","脚本正在试图访问无权限的路径，可能为恶意脚本！")
						t.Terminate()
						return nil
					}
					if !AllowPath(p){
						throwException("FB_ReadFile","脚本正在试图访问禁止访问的路径，可能为恶意脚本！")
						t.Terminate()
						return nil
					}
					err := hb.SaveFile(p,data)
					if err != nil {
						value,_:=v8go.NewValue(iso,false)
						return value
					}
					value,_:=v8go.NewValue(iso,true)
					return value
				}
			}
			return nil
		})); err!=nil{panic(err)}


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

	// function FB_WaitConnect() None
	if err:=global.Set("FB_AutoRestart",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			hb.RequireAutoRestart()
			return nil
		})); err!=nil{panic(err)}

	// FB_websocketConnectV2(address string,onNewMessage func(msgType int,data string)) func SendMsg(msgType int, data string)
	// 一般情况下，MessageType 为1(Text Messsage),即字符串类型，或者 0 byteArray (也被以字符串的方式传递)
	// onNewMessage 在连接关闭时会读取到两个null值
	if err:=global.Set("FB_WebSocketConnectV2",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if address,ok:= hasStrIn(info,0,"FB_WebSocketConnectV2[address]"); !ok{
				throwException("FB_WebSocketConnectV2",address)
			}else{
				if errStr,cbFn:=hasFuncIn(info,1,"FB_WebSocketConnectV2[onNewMessage]"); cbFn==nil{
					throwException("FB_WebSocketConnectV2",errStr)
				}else{
					ctx:=info.Context()
					conn, _, err := websocket.DefaultDialer.Dial(address, nil)
					if err != nil {
						return throwException("FB_WebSocketConnectV2",err.Error())
					}
					jsWriteFn:=v8go.NewFunctionTemplate(iso, func(writeInfo *v8go.FunctionCallbackInfo) *v8go.Value {
						if t.isTerminated{
							return nil
						}
						if len(writeInfo.Args())<2{
							throwException("SendMsg returned by FB_websocketConnectV2","not enough arguments")
							return nil
						}
						if !writeInfo.Args()[1].IsString(){
							throwException("SendMsg returned by FB_websocketConnectV2","SendMsg[data] should be string")
						}
						msgType:=int(writeInfo.Args()[0].Number())
						err := conn.WriteMessage(msgType, []byte(writeInfo.Args()[1].String()))
						if err != nil {
							return throwException("SendMsg returned by FB_websocketConnectV2","write fail")
						}
						return nil
					})
					go func() {
						msgType, data, err := conn.ReadMessage()
						if t.isTerminated{
							return
						}
						if err != nil {
							cbFn.Call(info.This(),v8go.Null(iso),v8go.Null(iso))
							return
						}
						jsMsgType,err:=v8go.NewValue(iso,int32(msgType))
						jsMsgData,err:=v8go.NewValue(iso,string(data))
						cbFn.Call(info.This(),jsMsgType,jsMsgData)
					}()
					return jsWriteFn.GetFunction(ctx).Value
				}
			}
			return nil
		})); err!=nil{panic(err)}


	// fetch
	if err := fetch.InjectTo(iso, global); err != nil {
		panic(err)
	}
	// setTimeout, clearTimeout, setInterval and clearInterval
	if err :=timers.InjectTo(iso, global);err!=nil{
		panic(err)
	}
	//  atob and btoa
	if err:=base64.InjectTo(iso,global);err!=nil{
		panic(err)
	}
	return func() {
		t.Terminate()
	}
}

func CtxFunctionInject(ctx *v8go.Context){
	// URL and URLSearchParams
	if err:=url.InjectTo(ctx);err!=nil{
		panic(err)
	}
}