// FB提供的函数中,以下四个函数为同步函数
// FB_WaitConnect()
// FB_RequireUserInput(hint)
// FB_SendMCCmdAndGetResult(mcCmd)
// FB_RequireFilePermission(dir)
// 所谓同步函数，就是脚本会完全停止，直到获得结果

// 其中，三个函数有异步版本，所谓异步，即脚本不会停止
// 当获得结果时，函数会被回调


afterGettedUserInput=function (userInput){
    FB_Println("成功获得了用户输入！"+userInput)
}

afterGettedCmdResult=function (result){
    FB_Println("成功获得了指令结果！"+result)

    // 当获得用户输入后，afterGettedUserInput会被回调
    FB_RequireUserInputAsync("随便输入一点什么",afterGettedUserInput)
}

afterConnected=function () {
    FB_Println("成功连接到MC了！")

    // 当获得指令结果后，afterGettedCmdResult会被回调
    FB_SendMCCmdAndGetResultAsync("list",afterGettedCmdResult)
}

// 当连接到MC后，afterConnected会被回调
FB_WaitConnectAsync(afterConnected)

FB_Println("和FB_WaitConnect不同，即使没有连接到FB，我也会执行")
