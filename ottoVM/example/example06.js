// example06.js
// 本脚本演示了FB_WaitConnect,FB_SendMCCmdAndGetResult,FB_RequireUserInput
// 的异步版本 FB_WaitConnectAsync,FB_SendMCCmdAndGetResultAsync,FB_RequireUserInputAsync

FB_WaitConnectAsync(function () {
    FB_Println("成功连接到服务器了!")
    FB_SendMCCmdAndGetResultAsync("tp @s @r",function (pk) {
        FB_Println("成功收到指令结果了!")
        FB_Println(JSON.stringify(pk))
        FB_RequireUserInputAsync("随便输入一点什么",function (userInput) {
            FB_Println("成功接收到用户输入了！"+userInput)
        })
    })
})