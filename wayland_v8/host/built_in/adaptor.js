class engine{
    static setName(name){
        return FB_SetName(name)
    }
    static waitConnectionSync(){
        return FB_WaitConnect()
    }
    static waitConnection(onConnect){
        return FB_WaitConnectAsync(onConnect)
    }
    static message(msg){
        return FB_Println(msg)
    }
    static question(hint,onResult){
        return FB_RequireUserInputAsync(hint,onResult)
    }
    static query(info){
        return FB_Query(info)
    }
    static crash(reason){
        return FB_CrashScript(reason)
    }
    static autoRestart(){
        return FB_AutoRestart()
    }
}
class game{
    static eval(fbCmd){
        return FB_GeneralCmd(fbCmd)
    }
    static oneShotCommand(mcCmd){
        return FB_SendMCCmd(mcCmd)
    }
    static sendCommandSync(mcCmd){
        return FB_SendMCCmdAndGetResultAsync(mcCmd)
    }
    static sendCommand(mcCmd,onResult){
        return FB_SendMCCmdAndGetResultAsync(mcCmd,onResult)
    }
    static botPos(){
        return FB_GetBotPos()
    }
    static subscribePacket(packetType,onPacket){
        return FB_RegPackCallBack(packetType,onPacket)
    }
    static listenChat(onMsg){
        return FB_RegChat(onMsg)
    }
}

class storage{
    static requestFilePermission(hint,path){
        return FB_RequireFilePermission(hint,path)
    }
    static readFile(path){
        return FB_ReadFile(path)
    }
    static writeFile(path,string){
        return FB_SaveFile(path,string)
    }
}

class websocket{
    static connect(address,onNewMessage){
        return FB_WebSocketConnectV2(address,onNewMessage)
    }
    static serve(address,pattern,onNewConnection){
        return FB_WebSocketServeV2(address,pattern,onNewConnection)
    }
}

// FB_SetName
// FB_WaitConnect
// FB_WaitConnectAsync
// FB_RequireUserInput
// FB_RequireUserInputAsync
// FB_Println
// FB_Query
// FB_GeneralCmd
// FB_SendMCCmd
// FB_SendMCCmdAndGetResultAsync
// FB_SendMCCmdAndGetResult
// FB_RegPackCallBack
// FB_RegChat
// FB_GetBotPos
// FB_GetAbsPath
// FB_RequireFilePermission
// FB_ReadFile
// FB_SaveFile
// FB_CrashScript
// FB_AutoRestart
// FB_WebSocketConnectV2
// FB_WebSocketServeV2