
// 当收到新消息时，这个函数会被调用
function onMessage(msgType,msg,sendFn,closeFn){
    FB_Println("recv Msg: "+msgType+": "+msg)
    sendFn(msgType,"server_echo: "+msg)
}

// 当有新连接时 这个函数会被调用
function onConnect(sendFn,closeFn){
    FB_Println("New Connection!")

    // 通过这个函数可以发送数据
    sendFn(1,"Hello Client!")
    return function (msgType,msg) {onMessage(msgType,msg,sendFn,closeFn)}
}

// 可以通过 ws://localhost:8888/ws_test 连接
// 即，与例6相同
FB_WebSocketServeV2(":8888","/ws_test",onConnect)

// script wayland_v8/example/example13_ws_server.js
// script wayland_v8/example/example06_websocket.js