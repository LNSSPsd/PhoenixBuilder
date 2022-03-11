package host

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/wayland_v8/host/built_in"
	"strings"

	"phoenixbuilder/fastbuilder/script"

	"github.com/gorilla/websocket"
	"go.kuoruan.net/v8go-polyfills/base64"
	"go.kuoruan.net/v8go-polyfills/fetch"
	"go.kuoruan.net/v8go-polyfills/timers"
	"go.kuoruan.net/v8go-polyfills/url"
	"rogchap.com/v8go"
)

const JSVERSION = "v8.gamma.4"

func AllowPath(path string) bool {
	if strings.Contains(path, "fbtoken") {
		return false
	}
	if strings.Contains(path, "fb_script_permission") {
		return false
	}
	return true
}

func LoadPermission(hb script.HostBridge, identifyStr string) map[string]bool {
	permission := map[string]bool{}
	fullPermission := map[string]map[string]bool{}
	file, err := hb.LoadFile("fb_script_permission.json")
	if err != nil {
		return permission
	}
	err = json.Unmarshal([]byte(file), &fullPermission)
	if err != nil {
		return permission
	}
	if savedPermission, ok := fullPermission[identifyStr]; ok {
		return savedPermission
	}
	return permission
}

func SavePermission(hb script.HostBridge, identifyStr string, permission map[string]bool) {
	fullPermission := map[string]map[string]bool{}
	file, err := hb.LoadFile("fb_script_permission.json")
	dataToSave := []byte{}
	if err == nil {
		json.Unmarshal([]byte(file), &fullPermission)
	}
	fullPermission[identifyStr] = permission
	dataToSave, _ = json.Marshal(fullPermission)
	hb.SaveFile("fb_script_permission.json", string(dataToSave))
}

func InitHostFns(iso *v8go.Isolate, global *v8go.ObjectTemplate, hb script.HostBridge, _scriptName string, identifyStr string, scriptPath string) func() {
	scriptName := _scriptName
	permission := LoadPermission(hb, identifyStr)
	updatePermission := func() {
		SavePermission(hb, identifyStr, permission)
	}

	throwException := func(funcName string, str string) *v8go.Value {
		errS := "Script triggered an exception at [" + funcName + "] due to " + str
		value, _ := v8go.NewValue(iso, errS)
		fmt.Println(errS)
		iso.ThrowException(value)
		return nil
	}
	printException := func(funcName string, str string) *v8go.Value {
		fmt.Println("Script triggered an exception at [" + funcName + "] due to " + str)
		return nil
	}
	throwNotConnectedException := func(funcName string) *v8go.Value {
		return throwException(funcName, "connection to MC not established")
	}
	hasStrIn := func(info *v8go.FunctionCallbackInfo, pos int, argName string) (string, bool) {
		if len(info.Args()) < pos+1 {
			return fmt.Sprintf("no arg %v provided in pos %v", argName, pos), false
		}
		if !info.Args()[pos].IsString() {
			return fmt.Sprintf("arg %v in pos %v is not a string (you set: %v)", argName, pos, info.Args()[pos].String()), false
		}
		return info.Args()[pos].String(), true
	}
	hasFuncIn := func(info *v8go.FunctionCallbackInfo, pos int, argName string) (string, *v8go.Function) {
		if len(info.Args()) < pos+1 {
			return fmt.Sprintf("no arg %v provided in pos %v", argName, pos), nil
		}
		function, err := info.Args()[pos].AsFunction()
		if err != nil {
			return fmt.Sprintf("arg %v in pos %v is not a function (you set: %v)", argName, pos, info.Args()[pos].String()), nil
		}
		return "", function
	}
	t := script.NewTerminator()
	t.TerminateHook = append(t.TerminateHook, func() {
		iso.TerminateExecution()
	})
	if err := global.Set("FB_SetName",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str, ok := hasStrIn(info, 0, "FB_SetName[scriptName]"); !ok {
				throwException("FB_SetName: No arguments assigned", str)
			} else {
				hb.Println("Script \""+scriptName+"\" is naming itself as \""+str+"\"", t, scriptName)
				scriptName = str
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	if err := global.Set("FB_WaitConnect",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			hb.WaitConnect(t)
			return nil
		}),
	); err != nil {
		panic(err)
	}

	if err := global.Set("FB_WaitConnectAsync",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			_args := info.Args()
			if len(_args) == 0 {
				throwException("FB_WaitConnectAsync(cb)", " No arguments assigned")
			}
			first_arg := _args[0]
			if !first_arg.IsFunction() {
				throwException("FB_WaitConnectAsync(cb)", " Callback should be a function")
			}
			f, e := first_arg.AsFunction()
			if e != nil {
				throwException("FB_WaitConnectAsync(cb)", " Callback should be a function, but got function.")
			}
			go func() {
				hb.WaitConnect(t)
				f.Call(info.This())
			}()
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_Println(msg string) None
	if err := global.Set("FB_Println",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str, ok := hasStrIn(info, 0, "FB_Println[msg]"); !ok {
				throwException("FB_Println", str)
			} else {
				hb.Println(str, t, scriptName)
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_GeneralCmd(fbCmd string) None
	if err := global.Set("FB_GeneralCmd",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected() {
				throwNotConnectedException("FB_GeneralCmd")
			}
			if str, ok := hasStrIn(info, 0, "FB_GeneralCmd[fbCmd]"); !ok {
				throwException("FB_GeneralCmd", str)
			} else {
				hb.FBCmd(str, t)
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_SendMCCmd(mcCmd string) None
	if err := global.Set("FB_SendMCCmd",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected() {
				throwNotConnectedException("FB_SendMCCmd")
			}
			if str, ok := hasStrIn(info, 0, "FB_SendMCCmd[mcCmd]"); !ok {
				throwException("FB_SendMCCmd", str)
			} else {
				hb.MCCmd(str, t, false)
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_SendMCCmdAndGetResultAsync(mcCmd string) jsObject
	if err := global.Set("FB_SendMCCmdAndGetResultAsync",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected() {
				throwNotConnectedException("FB_SendMCCmdAndGetResultAsync")
			}
			if str, ok := hasStrIn(info, 0, "FB_SendMCCmdAndGetResultAsync[mcCmd]"); !ok {
				throwException("FB_SendMCCmdAndGetResultAsync", str)
			} else {
				pk := hb.MCCmd(str, t, true)
				strPk, err := json.Marshal(pk)
				if err != nil {
					return throwException("FB_SendMCCmdAndGetResultAsync", "Cannot convert host packet to Json Str: "+str)
				}
				value, err := v8go.JSONParse(info.Context(), string(strPk))
				if err != nil {
					return throwException("FB_SendMCCmdAndGetResultAsync", str)
				} else {
					return value
				}
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_SendMCCmdAndGetResult(mcCmd string, onResult function(jsObject)) None
	// jsObject=null, if cannot get result in callback
	if err := global.Set("FB_SendMCCmdAndGetResult",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if !hb.IsConnected() {
				throwNotConnectedException("FB_SendMCCmdAndGetResult")
			}
			if str, ok := hasStrIn(info, 0, "FB_SendMCCmdAndGetResult[mcCmd]"); !ok {
				throwException("FB_SendMCCmdAndGetResult", str)
			} else {
				if _, cbFn := hasFuncIn(info, 1, "FB_SendMCCmdAndGetResult[onResult]"); cbFn == nil {
					hb.MCCmd(str, t, false)
					return nil
				} else {
					go func() {
						pk := hb.MCCmd(str, t, true)
						strPk, err := json.Marshal(pk)
						if err != nil {
							printException("FB_SendMCCmdAndGetResult", "Cannot convert host packet to Json Str: "+str)
							cbFn.Call(info.This(), v8go.Null(iso))
							return
						}
						val, err := v8go.JSONParse(info.Context(), string(strPk))
						if err != nil {
							printException("FB_SendMCCmdAndGetResult", "Cannot Parse Json Packet in Host: "+str)
							cbFn.Call(info.This(), v8go.Null(iso))
							return
						} else {
							cbFn.Call(info.This(), val)
						}
					}()
				}
				return nil
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	if err := global.Set("FB_GetBotPos",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			ot := v8go.NewObjectTemplate(iso)
			x, y, z := hb.GetBotPos()
			jsX, _ := v8go.NewValue(iso, int32(x))
			jsY, _ := v8go.NewValue(iso, int32(y))
			jsZ, _ := v8go.NewValue(iso, int32(z))
			ot.Set("x", jsX)
			ot.Set("y", jsY)
			ot.Set("z", jsZ)
			jsPos, _ := ot.NewInstance(info.Context())
			return jsPos.Value
		})); err != nil {
		panic(err)
	}

	// function FB_RequireUserInput(hint string) string
	if err := global.Set("FB_RequireUserInput",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str, ok := hasStrIn(info, 0, "FB_RequireUserInput[hint]"); !ok {
				throwException("FB_RequireUserInput", str)
			} else {
				userInput := hb.GetInput(str, t, scriptName)
				value, _ := v8go.NewValue(iso, userInput)
				return value
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_RequireUserInputAsync(hint,cb) None
	if err := global.Set("FB_RequireUserInputAsync",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str, ok := hasStrIn(info, 0, "FB_RequireUserInputAsync[hint,cb]"); !ok {
				throwException("FB_RequireUserInputAsync", str)
			} else {
				if errStr, cbFn := hasFuncIn(info, 1, "FB_RequireUserInputAsync[hint,cb]"); cbFn == nil {
					throwException("FB_RequireUserInputAsync", errStr)
				} else {
					go func() {
						userInput := hb.GetInput(str, t, scriptName)
						value, _ := v8go.NewValue(iso, userInput)
						cbFn.Call(info.This(), value)
					}()
				}

			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_RegPackCallBack(packetType,onPacketCb) deRegFn
	// when deRegFn is called, onPacketCb function will no longer be called
	if err := global.Set("FB_RegPackCallBack",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str, ok := hasStrIn(info, 0, "FB_RegPackCallBack[packetType]"); !ok {
				throwException("FB_RegPackCallBack", str)
			} else {
				if errStr, cbFn := hasFuncIn(info, 1, "game.listenPacket[onPacketCb]"); cbFn == nil {
					throwException("FB_RegPackCallBack", errStr)
				} else {
					deRegFn, err := hb.RegPacketCallBack(str, func(pk packet.Packet) {
						strPk, err := json.Marshal(pk)
						if err != nil {
							printException("FB_RegPackCallBack", "Cannot convert host packet to Json Str: "+err.Error())
							cbFn.Call(info.This(), v8go.Null(iso))
						} else {
							val, err := v8go.JSONParse(info.Context(), string(strPk))
							if err != nil {
								printException("FB_RegPackCallBack", "Cannot Parse Json Packet in Host: "+str)
								cbFn.Call(info.This(), v8go.Null(iso))
								return
							} else {
								cbFn.Call(info.This(), val)
							}
						}
					}, t)
					if err != nil {
						return throwException("FB_RegPackCallBack", err.Error())
					}
					jsCbFn := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
						deRegFn()
						return nil
					})
					return jsCbFn.GetFunction(info.Context()).Value
				}
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_RegChat(onMsg function(name,msg)) deRegFn
	if err := global.Set("FB_RegChat",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if errStr, cbFn := hasFuncIn(info, 0, "FB_RegChat[onMsg]"); cbFn == nil {
				throwException("FB_RegChat", errStr)
			} else {
				ctx := info.Context()
				deRegFn, err := hb.RegPacketCallBack("IDText", func(pk packet.Packet) {
					p := pk.(*packet.Text)
					SourceName, err := v8go.NewValue(iso, p.SourceName)
					if err != nil {
						printException("FB_RegChat", err.Error())
						cbFn.Call(info.This(), v8go.Null(iso), v8go.Null(iso))
						return
					}
					Message, err := v8go.NewValue(iso, p.Message)
					if err != nil {
						printException("FB_RegChat", err.Error())
						cbFn.Call(info.This(), v8go.Null(iso), v8go.Null(iso))
						return
					}
					cbFn.Call(info.This(), SourceName, Message)
				}, t)
				if err != nil {
					return throwException("FB_RegChat", err.Error())
				}
				jsCbFn := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
					deRegFn()
					return nil
				})
				t.TerminateHook = append(t.TerminateHook, deRegFn)
				return jsCbFn.GetFunction(ctx).Value
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_Query(info string) string
	if err := global.Set("FB_Query",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str, ok := hasStrIn(info, 0, "FB_Query[info]"); !ok {
				throwException("FB_Query", str)
			} else {
				if str == "script_sha256" {
					value, _ := v8go.NewValue(iso, identifyStr)
					return value
				} else if str == "script_path" {
					value, _ := v8go.NewValue(iso, scriptPath)
					return value
				} else if str == "engine_version" {
					value, _ := v8go.NewValue(iso, JSVERSION)
					return value
				}
				if ret, ok := hb.GetQueries()[str]; ok {
					value, _ := v8go.NewValue(iso, ret)
					return value
				}
				return nil
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_GetAbsPath(path string) string
	// I think we should allow the script to tell user where a file is
	if err := global.Set("FB_GetAbsPath",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str, ok := hasStrIn(info, 0, "FB_GetAbsPath[path]"); !ok {
				throwException("FB_GetAbsPath", str)
			} else {
				absPath := hb.GetAbsPath(str)
				value, _ := v8go.NewValue(iso, absPath)
				return value
			}
			return nil
		})); err != nil {
		panic(err)
	}

	// function FB_RequireFilePermission(hint,path) isSuccess
	if err := global.Set("FB_RequireFilePermission",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if hint, ok := hasStrIn(info, 1, "FB_RequireFilePermission[hint]"); !ok {
				throwException("FB_RequireFilePermission", hint)
				return nil
			} else {
				if dir, ok := hasStrIn(info, 0, "FB_RequireFilePermission[hint]"); !ok {
					throwException("FB_RequireFilePermission", dir)
					return nil
				} else {
					dir = hb.GetAbsPath(dir) + string(os.PathSeparator)
					if !AllowPath(dir) {
						throwException("FB_RequireFilePermission", "The script is breaking out sandbox, aborting.")
						t.Terminate()
						return nil
					}
					permissionKey := "VisitDir:" + dir
					if hasPermission, ok := permission[permissionKey]; ok && hasPermission {
						value, _ := v8go.NewValue(iso, true)
						return value
					} else {
						for {
							warning := "Script[" + scriptName + "][" + _scriptName + "]wants to access the contents of directory " + dir + ".\n" +
								"Reason " + hint + "\n" +
								"(Warning: The script will gain the ability of REMOVING, MODIFYING, CREATING any file in this directory.)\n" +
								"Allow the access? Give an answer[y/N]:"
							choose := hb.GetInput(warning, t, scriptName)
							if choose == "Y" || choose == "y" {
								value, _ := v8go.NewValue(iso, true)
								permission[permissionKey] = true
								updatePermission()
								return value
							} else {
								value, _ := v8go.NewValue(iso, false)
								return value
							}
							//hb.Println("无效输入，请输入[是/否/Y/y/N/n]其中之一",t,scriptName)
						}
					}
				}
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_ReadFile(path string) string
	// if permission is not granted or read fail, "" is returned
	if err := global.Set("FB_ReadFile",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if str, ok := hasStrIn(info, 0, "FB_ReadFile[path]"); !ok {
				throwException("FB_ReadFile", str)
			} else {
				p := hb.GetAbsPath(str)
				hasPermission := false
				for permissionName, _ := range permission {
					if strings.HasPrefix(permissionName, "VisitDir:") {
						if strings.HasPrefix(p, permissionName[len("VisitDir:"):]) {
							hasPermission = true
							break
						}
					}
				}
				if !hasPermission {
					throwException("FB_ReadFile", "The script is trying to access an external path (without permission), aborting.")
					t.Terminate()
					return nil
				}
				if !AllowPath(p) {
					throwException("FB_ReadFile", "The script is trying to access an external path (without permission), aborting.")
					t.Terminate()
					return nil
				}
				data, err := hb.LoadFile(p)
				if err != nil {
					value, _ := v8go.NewValue(iso, "")
					return value
				}
				value, _ := v8go.NewValue(iso, data)
				return value
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	// function FB_SaveFile(path string,data string) isSuccess
	if err := global.Set("FB_SaveFile",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if p, ok := hasStrIn(info, 0, "FB_SaveFile[path]"); !ok {
				throwException("FB_SaveFile", p)
			} else {
				if data, ok := hasStrIn(info, 1, "FB_SaveFile[data]"); !ok {
					throwException("FB_SaveFile", data)
				} else {
					p := hb.GetAbsPath(p)
					hasPermission := false
					for permissionName, _ := range permission {
						if strings.HasPrefix(permissionName, "VisitDir:") {
							if strings.HasPrefix(p, permissionName[len("VisitDir:"):]) {
								hasPermission = true
								break
							}
						}
					}
					if !hasPermission {
						throwException("FB_SaveFile", "The script is trying to access an external path (without permission), aborting.")
						t.Terminate()
						return nil
					}
					if !AllowPath(p) {
						throwException("FB_SaveFile", "The script is trying to access an external path (without permission), aborting.")
						t.Terminate()
						return nil
					}
					err := hb.SaveFile(p, data)
					if err != nil {
						value, _ := v8go.NewValue(iso, false)
						return value
					}
					value, _ := v8go.NewValue(iso, true)
					return value
				}
			}
			return nil
		})); err != nil {
		panic(err)
	}

	// function FB_CrashScript(string reason) None
	if err := global.Set("FB_CrashScript",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {

			if str, ok := hasStrIn(info, 0, "FB_CrashScript[reason]"); !ok {
				throwException("FB_CrashScript", str)
			} else {
				throwException("FB_CrashScript", str)
				t.Terminate()
			}
			return nil
		})); err != nil {
		panic(err)
	}

	// function FB_WaitConnect() None
	// 这里做了指数退避，几次重连失败就会放缓到1小时重连一次，见 main.go 170行
	if err := global.Set("FB_AutoRestart",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			hb.RequireAutoRestart()
			return nil
		})); err != nil {
		panic(err)
	}

	// FB_WebSocketConnectV2(address string,onNewMessage func(msgType int,data string)) func SendMsg(msgType int, data string)
	// 一般情况下，MessageType 为1(Text Messsage),即字符串类型，或者 0 byteArray (也被以字符串的方式传递)
	// onNewMessage 在连接关闭时会读取到两个null值
	if err := global.Set("FB_WebSocketConnectV2",
		v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
			if address, ok := hasStrIn(info, 0, "FB_WebSocketConnectV2[address]"); !ok {
				throwException("FB_WebSocketConnectV2", address)
			} else {
				if errStr, cbFn := hasFuncIn(info, 1, "FB_WebSocketConnectV2[onNewMessage]"); cbFn == nil {
					throwException("FB_WebSocketConnectV2", errStr)
				} else {
					conn, _, err := websocket.DefaultDialer.Dial(address, nil)
					if err != nil {
						return throwException("FB_WebSocketConnectV2", err.Error())
					}
					jsWriteFn := v8go.NewFunctionTemplate(iso, func(writeInfo *v8go.FunctionCallbackInfo) *v8go.Value {
						if t.Terminated() {
							return nil
						}
						if len(writeInfo.Args()) < 2 {
							throwException("SendMsg returned by FB_websocketConnectV2", "not enough arguments")
							return nil
						}
						if !writeInfo.Args()[1].IsString() {
							throwException("SendMsg returned by FB_websocketConnectV2", "SendMsg[data] should be string")
						}
						msgType := int(writeInfo.Args()[0].Number())
						err := conn.WriteMessage(msgType, []byte(writeInfo.Args()[1].String()))
						if err != nil {
							return throwException("SendMsg returned by FB_websocketConnectV2", "write fail")
						}
						return nil
					})
					go func() {
						for {
							msgType, data, err := conn.ReadMessage()
							if t.Terminated() {
								return
							}
							if err != nil {
								cbFn.Call(info.This(), v8go.Null(iso), v8go.Null(iso))
								return
							}
							jsMsgType, err := v8go.NewValue(iso, int32(msgType))
							jsMsgData, err := v8go.NewValue(iso, string(data))
							cbFn.Call(info.This(), jsMsgType, jsMsgData)
						}
					}()
					return jsWriteFn.GetFunction(info.Context()).Value
				}
			}
			return nil
		}),
	); err != nil {
		panic(err)
	}

	FB_WebSocketServeV2 := func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		var address, pattern, errStr string
		var jsOnConnect *v8go.Function
		var ok bool
		if address, ok = hasStrIn(info, 0, "FB_WebSocketConnectV2[address]"); !ok {
			throwException("FB_WebSocketConnectV2", address)
			return nil
		}
		if pattern, ok = hasStrIn(info, 1, "FB_WebSocketConnectV2[pattern]"); !ok {
			throwException("FB_WebSocketConnectV2", pattern)
			return nil
		}
		if errStr, jsOnConnect = hasFuncIn(info, 2, "FB_WebSocketConnectV2[onConnect]"); jsOnConnect == nil {
			throwException("FB_WebSocketConnectV2", errStr)
			return nil
		}
		http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
			wsConn, err := (&websocket.Upgrader{
				ReadBufferSize:  1024 * 10,
				WriteBufferSize: 1024 * 10,
			}).Upgrade(w, r, nil)
			if err != nil {
				throwException("FB_WebSocketServeV2", err.Error())
				return
			}

			jsSendFn := v8go.NewFunctionTemplate(iso, func(writeInfo *v8go.FunctionCallbackInfo) *v8go.Value {
				if t.Terminated() {
					return nil
				}
				if len(writeInfo.Args()) < 2 {
					throwException("FB_WebSocketServeV2.onConnect.sendMsg", "not enough arguments")
					return nil
				}
				if !writeInfo.Args()[1].IsString() {
					throwException("FB_WebSocketServeV2.onConnect.sendMsg", "SendMsg[data] should be string")
				}
				msgType := int(writeInfo.Args()[0].Number())
				err := wsConn.WriteMessage(msgType, []byte(writeInfo.Args()[1].String()))
				if err != nil {
					return throwException("FB_WebSocketServeV2.onConnect.sendMsg", "write fail")
				}
				return nil
			}).GetFunction(info.Context())
			jsCloseFn := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
				wsConn.Close()
				return nil
			}).GetFunction(info.Context())

			onMsgFn, err := jsOnConnect.Call(info.This(), jsSendFn, jsCloseFn)
			if err != nil {
				throwException("FB_WebSocketServeV2.onConnect.onMsgFn", err.Error())
				t.Terminate()
				return
			}
			jsOnMsgFn, err := onMsgFn.AsFunction()
			if err != nil {
				throwException("FB_WebSocketServeV2.onConnect.onMsgFn", err.Error())
				t.Terminate()
				return
			}
			go func() {
				for {
					msgType, data, err := wsConn.ReadMessage()
					//fmt.Println("line 739")
					if err != nil {
						jsOnMsgFn.Call(info.This(), v8go.Null(iso), v8go.Null(iso))
					} else {
						jsMsgType, _ := v8go.NewValue(iso, int32(msgType))
						jsMsgData, _ := v8go.NewValue(iso, string(data))
						jsOnMsgFn.Call(info.This(), jsMsgType, jsMsgData)
					}
					//if err!=nil{
					//	println(err)
					//	return
					//}
					//fmt.Printf("%v %v\n",msgType,data)
					//err = wsConn.WriteMessage(msgType, data)
					//if err!=nil{
					//	println(err)
					//	return
					//}
				}
			}()
		})
		go http.ListenAndServe(address, nil)

		//err := startWsServer(address, pattern, t, func(sendFn sendMessageFn, closeFn closeFn) onMessageCb {
		//	defaultOnMsgFn := func(int, string) {}
		//	jsSendFn := v8go.NewFunctionTemplate(iso, func(writeInfo *v8go.FunctionCallbackInfo) *v8go.Value {
		//		if t.Terminated() {
		//			return nil
		//		}
		//		if len(writeInfo.Args()) < 2 {
		//			throwException("FB_WebSocketServeV2.onConnect.sendMsg", "not enough arguments")
		//			return nil
		//		}
		//		if !writeInfo.Args()[1].IsString() {
		//			throwException("FB_WebSocketServeV2.onConnect.sendMsg", "SendMsg[data] should be string")
		//		}
		//		msgType := int(writeInfo.Args()[0].Number())
		//		err := sendFn(msgType, writeInfo.Args()[1].String())
		//		if err != nil {
		//			return throwException("FB_WebSocketServeV2.onConnect.sendMsg", "write fail")
		//		}
		//		return nil
		//	}).GetFunction(info.Context())
		//	jsCloseFn := v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		//		closeFn()
		//		return nil
		//	}).GetFunction(info.Context())
		//	jsOnMsgFn, err := jsOnConnect.Call(info.Context().Global(), jsSendFn, jsCloseFn)
		//	if err != nil {
		//		throwException("FB_WebSocketServeV2.onConnect.onMsgFn", err.Error())
		//		t.Terminate()
		//		return defaultOnMsgFn
		//	}
		//	onMsgFn, err := jsOnMsgFn.AsFunction()
		//	if err != nil {
		//		throwException("FB_WebSocketServeV2.onConnect.onMsgFn", err.Error())
		//		t.Terminate()
		//		return defaultOnMsgFn
		//	}
		//	wrappedOnMsgFn := func(msgType int, msg string) {
		//		if msgType == -1 {
		//			onMsgFn.Call(info.Context().Global(), v8go.Null(iso), v8go.Null(iso))
		//		} else {
		//			jsMsgType, _ := v8go.NewValue(iso, int32(msgType))
		//			jsMsgData, _ := v8go.NewValue(iso, string(msg))
		//			onMsgFn.Call(info.Context().Global(), jsMsgType, jsMsgData)
		//		}
		//	}
		//	return wrappedOnMsgFn
		//})

		return nil
	}

	// FB_WebSocketServeV2(address string,onNewMessage func(msgType int,data string)) func SendMsg(msgType int, data string)
	// 一般情况下，MessageType 为1(Text Messsage),即字符串类型，或者 0 byteArray (也被以字符串的方式传递)
	// onNewMessage 在连接关闭时会读取到两个null值
	if err := global.Set("FB_WebSocketServeV2", v8go.NewFunctionTemplate(iso, FB_WebSocketServeV2)); err != nil {
		panic(err)
	}

	// fetch
	if err := fetch.InjectTo(iso, global); err != nil {
		panic(err)
	}
	// setTimeout, clearTimeout, setInterval and clearInterval
	if err := timers.InjectTo(iso, global); err != nil {
		panic(err)
	}
	//  atob and btoa
	if err := base64.InjectTo(iso, global); err != nil {
		panic(err)
	}

	// now we use built-in js, see built_in folder
	// encryption encryption.aesEncrypt(text, key)
	//encryption:=v8go.NewObjectTemplate(iso)
	//global.Set("encryption", encryption)
	//if err := encryption.Set("aesEncrypt",
	//	v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
	//		if text, ok := hasStrIn(info, 0, "encryption.aesEncrypt[text]"); !ok {
	//			throwException("encryption.aesEncrypt", text)
	//		} else {
	//			if key, ok := hasStrIn(info, 1, "encryption.aesEncrypt[key]"); !ok {
	//				throwException("encryption.aesEncrypt", key)
	//			} else {
	//				encryptOut,iv,err := aesEncrypt(text,key)
	//				if err!=nil{
	//					throwException("encryption.aesEncrypt",err.Error())
	//					return nil
	//				}else{
	//					result:=v8go.NewObjectTemplate(iso)
	//					jsEncryptOut, _ := v8go.NewValue(iso, encryptOut)
	//					jsIV, _ := v8go.NewValue(iso, iv)
	//					result.Set("cipherText",jsEncryptOut)
	//					result.Set("iv",jsIV)
	//					obj,_:=result.NewInstance(info.Context())
	//					return obj.Value
	//				}
	//			}
	//		}
	//		return nil
	//	}),
	//); err != nil {
	//	panic(err)
	//}
	//if err := encryption.Set("aesDecrypt",
	//	v8go.NewFunctionTemplate(iso, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
	//		if text, ok := hasStrIn(info, 0, "encryption.aesDecrypt[text]"); !ok {
	//			throwException("encryption.aesDecrypt", text)
	//		} else {
	//			if key, ok := hasStrIn(info, 1, "encryption.aesDecrypt[key]"); !ok {
	//				throwException("encryption.aesDecrypt", key)
	//			} else {
	//				if iv, ok := hasStrIn(info, 2, "encryption.aesDecrypt[iv]"); !ok {
	//					throwException("encryption.aesDecrypt", key)
	//				} else{
	//					decryptOut,err := aesDecrypt(text,key,iv)
	//					if err!=nil{
	//						throwException("encryption.aesDecrypt",err.Error())
	//						return nil
	//					}else{
	//						value, _ := v8go.NewValue(iso, decryptOut)
	//						return value
	//					}
	//				}
	//			}
	//		}
	//		return nil
	//	}),
	//); err != nil {
	//	panic(err)
	//}

	return func() {
		t.Terminate()
	}
}

func CtxFunctionInject(ctx *v8go.Context) {
	// URL and URLSearchParams
	if err := url.InjectTo(ctx); err != nil {
		panic(err)
	}
	_, err := ctx.RunScript(built_in.GetbuiltIn(), "built_in")
	if err != nil {
		panic(err)
	}
}
