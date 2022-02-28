
// 等待直到完成FB初始化
function FB_WaitConnect() None

// 等效于用户在fb中输入一条指令
function FB_GeneralUserInput(fbCmd string) None

// 发送一条MC指令，不等待其返回值
function FB_SendMCCmd(mcCMD string) None

// 发送一条MC指令，并等待其返回结果，结果以json字符串形式表示
// 警告，一些指令没有返回值，这种指令会导致程序卡死
function FB_SendMCCmdAndGetResult(mcCMD string) string

// 订阅一种特定类型的数据包，当收到该种数据包时，指定的函数将会被调用，数据包以json形式传入
// 警告，不合理的利用该函数可能导致性能低下
function FB_RegPackCallBack(packetType string,callbackFn func(string)) None

// 取消订阅一种数据包
function FB_DeRegPack(packetType string) None

// 订阅聊天信息
// 实际上可以通过 FB_regPackCallBack 实现，但是毕竟这种信息相当常用
function FB_RegChat(callBackFunc(name string, msg string)) None

// 请求用户输入信息
function FB_UserInput(hint string) string

// 向用户显示一条信息
function FB_Println(s string) None

// 获得获取fb的某些信息，例如，用户的游戏名，无论是何种值，结果都以string形式返回
function FB_Query(info string) string 