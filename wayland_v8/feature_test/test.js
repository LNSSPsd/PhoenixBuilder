// // console.log("1","hello!\n")
// //
// // fetch('https://storage.fastbuilder.pro').then(function(r) {
// //     r.text().then(function(d) {
// //         FB_Println(r.statusText)
// //         for (var k in r.headers._headers) {
// //             FB_Println(k + ':', r.headers.get(k))
// //         }
// //         FB_Println(d)
// //     });
// // })
//
// FB_SetName("测试脚本")
// success=FB_RequireFilePermission("data","需要访问这个文件夹来保存数据")
// if(!success){
//     FB_ScriptCrash("必须这个文件夹的权限才能工作")
// }else{
//     FB_Println("成功获得了权限!")
// }
//
// data=FB_ReadFile("data1/test.txt")
// FB_Println(data)
// isSuccess=FB_SaveFile("data/test.txt","hello!")
// if (isSuccess){
//     FB_Println("文件保存成功！")
//     data=FB_ReadFile("data/test.txt")
//     FB_Println("文件内容是:"+data)
// }else{
//     FB_Println("文件保存失败！")
// }
//
//
// function FB_SendMCCmdAndGetResultAsync(mcCmd,cb) {
//     r=_FB_SendMCCmdAndGetResultAsync(mcCmd,function (strPk) {
//         FB_Println("_FB_SendMCCmdAndGetResultAsync")
//             cb(JSON.parse(strPk))
//     })
// }
//
// var counter=0
// var deRegFn=null
// deRegFn=FB_RegPackCallBack("IDText",function (pk) {
//     FB_Println("RegPackCallBack get packet "+JSON.stringify(pk))
//     FB_Println("counter "+counter)
//     counter++
//     if(counter===3){
//         deRegFn()
//     }
// })
//
// FB_RegChat(function (name,msg) {
//     FB_Println("Recv Chat String: "+name+": "+msg)
// })
//
// // function FB_WaitConnect() None
// FB_WaitConnect()
// FB_Println("Connected to FB!")
//
// // fuunction FB_WaitConnectAsync(cb func()) None
// FB_WaitConnectAsync(function () {
//     FB_Println("Async Connected to FB!")
//
// })
//
// var user_name=FB_Query("user_name")
// var sha_token=FB_Query("sha_token")
// FB_Println("user_name: "+user_name+" sha_token:"+sha_token)
//
// FB_GeneralCmd(".tp @s ~~~")
//
// userInput=FB_RequireUserInput("your name")
// FB_Println("name is "+userInput)
//
// FB_RequireUserInputAsync("your name",function (name) {
//     FB_Println("Async name is "+name)
// })
//
//
// FB_SendMCCmd("list")
//
// cmd_list_result=FB_SendMCCmdAndGetResult("list")
// FB_Println(JSON.stringify(cmd_list_result))
//
// FB_SendMCCmdAndGetResultAsync("list",function (result) {
//     FB_Println("Async Get Result: "+JSON.stringify(result))
// })
//
//
//
// // function FB_ScriptCrash(string reason) None
// // 让脚本崩溃
// // FB_ScriptCrash("crashed here!")
//

// 当接收到新消息时，这个函数会被调用
function onNewMessage(msgType,newMessage) {
    FB_Println("type: "+msgType+" message: "+newMessage)
}

// 连接到 ws://localhost:8888/ws_test 上
sendFn=FB_WebSocketConnectV2("ws://localhost:8888/ws_test",onNewMessage)

// 使用返回的发送函数向服务器发送消息
// 1 为 Msg 类型，（TextMessage）
sendFn(1,"hello ws!")