// 本脚本演示了自动将机器人移动到玩家身边，并设置全局延迟为 100
// 演示了 FB_Query，FB_SendMCCmdAndGetResult，FB_GeneralCmd 的功能

FB_SetName("FB信息和通用指令")

// 等待连接到 MC
FB_WaitConnect()

// 通用fb功能，相当于用户在fb中输入了这条指令
FB_GeneralCmd("delay set 100")

// 通过FB_Query 查询信息
userName = FB_Query("user_name")

// 查看当前玩家有哪些，只是为了演示功能才那么做，其实没必要
listResult = FB_SendMCCmdAndGetResult("list")
currentPlayers = listResult["OutputMessages"][1]["Parameters"] // "玩家1, 玩家2"

currentPlayersList = String(currentPlayers).split(", ")

FB_Println("当前的玩家有:")
currentPlayersList.forEach(function (playerName) {
    FB_Println(playerName)
    if (playerName === userName) {
        result = FB_SendMCCmdAndGetResult("tp @s " + userName)
        FB_Println("成功移动! " + JSON.stringify(result))
    }
})


// FB_Query 能查询的所有信息
// 脚本内容的哈希值
FB_Println(FB_Query("script_sha256"))
// 脚本所在路径
FB_Println(FB_Query("script_path"))
// JS解释器实现
FB_Println(FB_Query("engine_version"))
//用户名
FB_Println(FB_Query("user_name"))
//用户FB Token的哈希值
FB_Println(FB_Query("sha_token"))
//服务器代码
FB_Println(FB_Query("server_code"))
//FB 版本信息
FB_Println(FB_Query("fb_version"))
//工作路径(一般情况下就是fb所在路径)
FB_Println(FB_Query("fb_dir"))