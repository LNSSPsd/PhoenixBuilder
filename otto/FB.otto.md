
// 等待直到完成FB初始化
FB_waitConnect() None

// 等效于用户在fb中输入一条指令
FB_generalUserInput(fbCmd string) None

// 发送一条MC指令，不等待其返回值
FB_sendMCCmd(mcCMD string) None

// 发送一条MC指令，并等待其返回结果，结果以json字符串形式表示
// 警告，一些指令没有返回值，这种指令会导致程序卡死
FB_sendMCCmdAndGetResult(mcCMD string) string

// 订阅一种特定类型的数据包，当收到该种数据包时，指定的函数将会被调用，数据包以json形式传入
// 警告，不合理的利用该函数可能导致性能低下
FB_regPackCallBack(packetType int,callbackFn func(string)) None

// 取消订阅一种数据包
FB_deRegPack(packetType int) None

// 订阅聊天信息
// 实际上可以通过 FB_regPackCallBack 实现，但是毕竟这种信息相当常用
FB_regChat(callBackFunc(name string, msg string))

