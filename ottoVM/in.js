// 订阅数据包
FB_RegPackCallBack("IDPlayerList",function (pk) {
    console.log("JS: 收到订阅数据包"+pk)
    // console.log("收到数据包"+JSON.stringify(pk))
})

// 订阅聊天信息
FB_RegChat(function (player,message) {
    console.log("JS: player: "+player+" message: "+message)
})

// 等待连接到 MC
FB_WaitConnect()



// js: 计算时间
nowTime=new Date()
year=nowTime.getFullYear()
month=nowTime.getMonth()
day=nowTime.getDay()
hour=nowTime.getHours()
minute=nowTime.getMinutes()

FB_SendMCCmd("scoreboard players set "+scoreBoardName+" year "+year)
FB_SendMCCmd("scoreboard players set "+scoreBoardName+" month "+month)
FB_SendMCCmd("scoreboard players set "+scoreBoardName+" day "+day)
FB_SendMCCmd("scoreboard players set "+scoreBoardName+" hour "+hour)
FB_SendMCCmd("scoreboard players set "+scoreBoardName+" minute "+minute)

// 向用户发送提示信息
FB_Println("时间记分板校准完成！")

// // 发送FB指令
// FB_GeneralCmd(".say 特权用户 "+superUser+" 启动了菜单系统 ")

FB_GeneralCmd(".tp @s ~~~")

listResult=FB_SendMCCmdAndGetResult("list")
currentUsers=listResult["OutputMessages"][1]["Parameters"]
console.log("当前服务器中的用户 :",currentUsers)

