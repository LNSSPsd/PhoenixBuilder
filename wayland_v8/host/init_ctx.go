package host

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"rogchap.com/v8go"
	"time"
)

type Terminator struct {
	c chan struct{}
	isTeminated bool
}

type HostBridge interface{
	WaitConnect(t *Terminator)
	IsConnected() bool
	//Block(t *Terminator)
	Println(str string,t *Terminator,scriptName string)
	FBCmd(fbCmd string,t *Terminator)
	MCCmd(mcCmd string,t *Terminator,waitResult bool) *packet.CommandOutput
	GetInput(hint string,t *Terminator,scriptName string) string
	RegPacketCallBack(packetType string,onPacket func(packet.Packet),t *Terminator) (func(),error)
	Query(info string) string
}

type HostBridgeBeta struct {
	isConnected bool
	connetWaiter chan struct{}
	// cb funcs
	vmCbsCount map[uint32]uint64
	vmCbs map[uint32]map[uint64]func(packet.Packet)
	// query
	HostQueryExpose map[string]func()string
}

func NewHostBridge() *HostBridgeBeta {
	return &HostBridgeBeta{
		connetWaiter:make(chan struct{}),
		vmCbsCount: map[uint32]uint64{},
		vmCbs: map[uint32]map[uint64]func(packet.Packet){},
		HostQueryExpose: map[string]func() string{
			"user_name": func() string {
				return "2401PT"
			},
			"sha_token": func() string {
				return "sha_token12asjkdao23201"
			},
		},
	}
}

func (hb *HostBridgeBeta) WaitConnect(t *Terminator)  {
	if !hb.isConnected{
		timer:=time.NewTimer(time.Second*1)
		go func() {
			<-timer.C
			hb.isConnected=true
			close(hb.connetWaiter)
		}()
	}
	select {
	case <-hb.connetWaiter:
	case <-t.c:
	}
}

func (hb *HostBridgeBeta) IsConnected() bool {
	return hb.isConnected
}

func (hb *HostBridgeBeta) Println(str string,t *Terminator,scriptName string)  {
	if t.isTeminated{
		return
	}
	fmt.Println("["+scriptName+"]: "+str)
}

func (hb *HostBridgeBeta) FBCmd(fbCmd string,t *Terminator)  {
	if t.isTeminated{
		return
	}
	fmt.Println("[FBCmd]: "+fbCmd)
}

func (hb *HostBridgeBeta) MCCmd(mcCmd string,t *Terminator,waitResult bool) *packet.CommandOutput {
	if t.isTeminated{
		return nil
	}
	fmt.Println("[MCCmd]: "+mcCmd)
	if waitResult{
		return &packet.CommandOutput{
			CommandOrigin:  protocol.CommandOrigin{
				Origin:         1,
				UUID:           uuid.UUID{1,2,3,4,5,6,7,83,2,13},
				RequestID:      "RequestID",
				PlayerUniqueID: 5,
			},
			OutputType:     0,
			SuccessCount:   1,
			OutputMessages: []protocol.CommandOutputMessage{{
				Success:    true,
				Message:    "hello!",
				Parameters: nil,
			}},
			DataSet:        "",
		}
	}else{
		return nil
	}
}

func (hb *HostBridgeBeta) GetInput(hint string,t *Terminator,scriptName string) string{
	if t.isTeminated{
		return ""
	}

	fmt.Println("[scriptName]: "+hint)
	if t.isTeminated{
		return ""
	}

	return "test_input"
}

func (hb *HostBridgeBeta) RegPacketCallBack(packetType string,onPacket func(packet.Packet),t *Terminator) (func(),error){
	packetID,ok:=PacketNameMap[packetType]
	if !ok{
		return nil,fmt.Errorf("no such packet type "+packetType)
	}
	_c,ok:=hb.vmCbsCount[packetID]
	c:=_c
	if !ok{
		hb.vmCbsCount[packetID]=0
		hb.vmCbs[packetID]=make(map[uint64]func(packet.Packet))
		c=0
	}
	c+=1
	hb.vmCbsCount[packetID]++
	hb.vmCbs[packetID][c]=onPacket
	go func() {
		<-t.c
		if _,ok:=hb.vmCbs[packetID][c];ok{
			delete(hb.vmCbs[packetID],c)
		}

	}()
	go func() {
		for{
			if cb,ok:=hb.vmCbs[packetID][c];!ok{
				return
			}else{
				cb(&packet.Text{
					TextType:         0,
					NeedsTranslation: false,
					SourceName:       "fakeUser",
					Message:          "hello from routine",
					Parameters:       nil,
					XUID:             "",
					PlatformChatID:   "",
					PlayerRuntimeID:  "",
				})
				time.Sleep(3*time.Second)
			}
		}
	}()
	return func(){
		fmt.Println("DeReg called!")
		delete(hb.vmCbs[packetID],c)
	},nil
}

func (hb *HostBridgeBeta) Query(info string) string {
	if fn,ok :=hb.HostQueryExpose[info]; ok{
		return fn()
	} else{
		return ""
	}
}

func InitHostFns(iso *v8go.Isolate,global *v8go.ObjectTemplate,hb HostBridge,scriptName string) {
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
	}

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
			if !hb.IsConnected(){
				throwNotConnectedException("FB_RequireUserInput")
			}
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
			if !hb.IsConnected(){
				throwNotConnectedException("FB_RequireUserInputAsync")
			}
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
				userInput:= hb.Query(str)
				value,_:=v8go.NewValue(iso,userInput)
				return value
			}
			return nil
		})); err!=nil{panic(err)}



	//// function FB_Query(info string) string
	//if err := global.Set("FB_Query",
	//	func(info *v8go.FunctionCallbackInfo) *v8go.Value {
	//		if str,ok:= hasStrIn(info,0,"FB_Query[info]"); !ok{
	//			throwException("FB_Query",str)
	//		}else {
	//
	//		}
	//	},
	//); err!=nil{fmt.Println(err)}

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