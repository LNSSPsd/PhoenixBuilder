

function onMessage(msgType,msg,sendFn,closeFn){
    FB_Println("recv Msg: "+msgType+": "+msg)
    sendFn(msgType,"server_echo: "+msg)
}

function onConnect(sendFn,closeFn){
    FB_Println("New Connection!")
    sendFn(1,"Hello Client!")
    return function (msgType,msg) {onMessage(msgType,msg,sendFn,closeFn)}
}


websocket.serve(":8888","/ws_test",onConnect)