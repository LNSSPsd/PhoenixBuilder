package host

import (
	"encoding/json"
	"fmt"
	//"github.com/gorilla/websocket"
	"go.kuoruan.net/v8go-polyfills/base64"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"go.kuoruan.net/v8go-polyfills/timers"
	"go.kuoruan.net/v8go-polyfills/url"
	//"os"
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
	//permission:=LoadPermission(hb,identifyStr)
	/*updatePermission:= func() {
		SavePermission(hb,identifyStr,permission)
	}*/

	throwException:= func(funcName string,str string) *v8go.Value {
		value, _ := v8go.NewValue(iso, "Script crashed at ["+funcName+"] due to "+str)
		iso.ThrowException(value)
		return nil
	}
	printException:= func(funcName string,str string) *v8go.Value {
		fmt.Println("Script triggered an exception at ["+funcName+"] due to "+str)
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
	engine:=v8go.NewObjectTemplate(iso)
	global.Set("engine",engine)
	// function engine.setName(scriptName string)
	if err:=engine.Set("setName",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"engine.setName[scriptName]"); !ok{
				throwException("engine.setName: No arguments assigned",str)
			}else{
				hb.Println("Script \""+scriptName+"\" is naming itself as \""+str+"\"",t,scriptName)
				scriptName=str
			}
			return nil
		}),
	); err!=nil{panic(err)}

	// function engine.waitConnectionSync()
	if err:=engine.Set("waitConnectionSync",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			hb.WaitConnect(t)
			return nil
		}),
	); err!=nil{panic(err)}
	
	// function engine.waitConnection(cb)
	if err:=engine.Set("waitConnection",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			_args:=info.Args()
			if(len(_args)==0) {
				throwException("engine.waitConnection(cb)"," No arguments assigned")
			}
			first_arg:=_args[0]
			if(!first_arg.IsFunction()) {
				throwException("engine.waitConnection(cb)"," Callback should be a function")
			}
			f, e:=first_arg.AsFunction()
			if(e!=nil) {
				throwException("engine.waitConnection(cb)"," Callback should be a function, but got function.")
			}
			go func() {
				hb.WaitConnect(t)
				f.Call(info.Context().Global())
			} ()
			return nil
		}),
	); err!=nil{panic(err)}

	// function engine.message(msg string) None
	if err:=engine.Set("message",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"engine.message[msg]"); !ok{
				throwException("engine.message",str)
			}else{
				hb.Println(str,t,scriptName)
			}
			return nil
		}),
	); err!=nil{panic(err)}
	game:=v8go.NewObjectTemplate(iso)
	global.Set("game",game)
	// One shot command
	// function game.eval(fbCmd string) None
	if err:=game.Set("eval",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected(){
				throwNotConnectedException("game.eval")
			}
			if str,ok:= hasStrIn(info,0,"game.eval[fbCmd]"); !ok{
				throwException("game.eval",str)
			}else{
				hb.FBCmd(str,t)
			}
			return nil
		}),
	); err!=nil{panic(err)}

	// function game.oneShotCommand(mcCmd string) None
	if err:=game.Set("oneShotCommand",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected(){
				throwNotConnectedException("game.oneShotCommand")
			}
			if str,ok:= hasStrIn(info,0,"game.oneShotCommand[mcCmd]"); !ok{
				throwException("game.oneShotCommand",str)
			}else{
				hb.MCCmd(str,t,false)
			}
			return nil
		}),
	); err!=nil{panic(err)}

	// function game.sendCommandSync(mcCmd string) jsObject
	if err:=game.Set("sendCommandSync",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected(){
				throwNotConnectedException("game.sendCommandSync")
			}
			if str,ok:= hasStrIn(info,0,"game.sendCommandSync[mcCmd]"); !ok{
				throwException("game.sendCommandSync",str)
			}else{
				pk:= hb.MCCmd(str,t,true)
				strPk, err := json.Marshal(pk)
				if err != nil {
					return throwException("game.sendCommandSync","Cannot convert host packet to Json Str: "+str)
				}
				value, err := v8go.JSONParse(info.Context(), string(strPk))
				if err != nil {
					return throwException("game.sendCommandSync",str)
				}else{
					return value
				}
			}
			return nil
		}),
	); err!=nil{panic(err)}

	// function game.sendCommand(mcCmd string, onResult function(jsObject)) None
	// jsObject=null, if cannot get result in callback
	if err:=game.Set("sendCommand",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected(){
				throwNotConnectedException("game.sendCommand")
			}
			if str,ok:= hasStrIn(info,0,"game.sendCommand[mcCmd]"); !ok{
				throwException("game.sendCommand",str)
			}else{
				if _,cbFn:=hasFuncIn(info,1,"game.sendCommand[onResult]"); cbFn==nil{
					hb.MCCmd(str,t,false)
					return nil
				}else{
					ctx:=info.Context()
					go func() {
						pk:= hb.MCCmd(str,t,true)
						strPk, err := json.Marshal(pk)
						if err != nil {
							printException("game.sendCommand","Cannot convert host packet to Json Str: "+str)
							cbFn.Call(info.Context().Global(),v8go.Null(iso))
							return
						}
						val, err := v8go.JSONParse(ctx,string(strPk))
						if err != nil {
							printException("game.sendCommand","Cannot Parse Json Packet in Host: "+str)
							cbFn.Call(info.Context().Global(),v8go.Null(iso))
							return
						}else {
							cbFn.Call(info.Context().Global(),val)
						}
					}()
				}
				return nil
			}
			return nil
		}),
	); err!=nil{panic(err)}

	// function engine.questionSync(hint string) string
	if err:=engine.Set("questionSync",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"engine.questionSync[hint]"); !ok{
				throwException("engine.questionSync",str)
			}else{
				userInput:= hb.GetInput(str,t,scriptName)
				value,_:=v8go.NewValue(iso,userInput)
				return value
			}
			return nil
		}),
	); err!=nil{panic(err)}

	// function question(hint,cb) None
	if err:=engine.Set("question",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"engine.question[hint,cb]"); !ok{
				throwException("engine.question",str)
			}else{
				if errStr,cbFn:=hasFuncIn(info,1,"engine.question[hint,cb]"); cbFn==nil{
					throwException("engine.question",errStr)
				}else{
					go func() {
						userInput:= hb.GetInput(str,t,scriptName)
						value,_:=v8go.NewValue(iso,userInput)
						cbFn.Call(info.Context().Global(),value)
					}()
				}

			}
			return nil
		}),
	); err!=nil{panic(err)}

	// function game.listenPacket(packetType,onPacketCb) deRegFn
	// when deRegFn is called, onPacketCb function will no longer be called
	if err:=game.Set("listenPacket",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str,ok:= hasStrIn(info,0,"game.listenPacket[packetType]"); !ok{
				throwException("game.listenPacket",str)
			}else{
				if errStr,cbFn:=hasFuncIn(info,1,"game.listenPacket[onPacketCb]"); cbFn==nil{
					throwException("game.listenPacket",errStr)
				}else{
					ctx:=info.Context()
					deRegFn, err := hb.RegPacketCallBack(str, func(pk packet.Packet) {
						strPk, err := json.Marshal(pk)
						if err!=nil{
							printException("game.listenPacket","Cannot convert host packet to Json Str: "+err.Error())
							cbFn.Call(info.This(),v8go.Null(iso))
						}else{
							val, err := v8go.JSONParse(ctx,string(strPk))
							if err != nil {
								printException("game.listenPacket","Cannot Parse Json Packet in Host: "+str)
								cbFn.Call(ctx.Global(),v8go.Null(iso))
								return
							}else {
								cbFn.Call(ctx.Global(),val)
							}
						}
					},t)
					if err != nil {
						return throwException("game.listenPacket",err.Error())
					}
					jsCbFn:=v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
						deRegFn()
						return nil
					})
					return jsCbFn.GetFunction(ctx).Value
				}
			}
			return nil
		}),
	); err!=nil{panic(err)}

	// function game.listenChat(onMsg function(name,msg)) deRegFn
	if err:=game.Set("listenChat",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if errStr,cbFn:=hasFuncIn(info,0,"game.listenChat[onMsg]"); cbFn==nil{
				throwException("game.listenChat",errStr)
			}else{
				ctx:=info.Context()
				deRegFn, err := hb.RegPacketCallBack("IDText", func(pk packet.Packet) {
					p := pk.(*packet.Text)
					SourceName, err := v8go.NewValue(iso,p.SourceName)
					if err != nil {
						printException("game.listenChat",err.Error())
						cbFn.Call(info.Context().Global(),v8go.Null(iso),v8go.Null(iso))
						return
					}
					Message, err := v8go.NewValue(iso,p.Message)
					if err != nil {
						printException("game.listenChat",err.Error())
						cbFn.Call(info.Context().Global(),v8go.Null(iso),v8go.Null(iso))
						return
					}
					cbFn.Call(info.Context().Global(),SourceName,Message)
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
		}),
	); err!=nil{panic(err)}
	
	consts:=v8go.NewObjectTemplate(iso)
	s256v,_:=v8go.NewValue(iso,identifyStr)
	consts.Set("script_sha256",s256v)
	/*
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
		}),
	); err!=nil{panic(err)}

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
		})); err!=nil{panic(err)}*/


	// function engine.crash(string reason) None
	if err:=engine.Set("crash",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {

			if str,ok:= hasStrIn(info,0,"engine.crash[reason]"); !ok{
				throwException("engine.crash",str)
			}else{
				throwException("engine.crash",str)
				t.Terminate()
			}
			return nil
		})); err!=nil{panic(err)}
	/*
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
		})); err!=nil{panic(err)}*/


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