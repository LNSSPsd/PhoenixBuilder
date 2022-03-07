// 本脚演示了时间记分板的校正
// 演示了 FB_SetName FB_WaitConnect，FB_RequireUserInput，FB_Println，FB_SendMCCmd 的功能
// 假设用户有一个记分板，记分板里有 year, month, day, hour, minute 四个项目
// 需要与现实时间同步

// 这个不是必须的，不设置时会以脚本文件名作为名字
FB_SetName("时间校准")

// 等待连接到 MC
FB_WaitConnect()
FB_Println("已经连接到服务器!")

// 请求用户输入信息 (时间相关记分板的名字)
scoreBoardName=FB_RequireUserInput("时间记分板的名字是?")

// js: 计算时间
nowTime=new Date()
nowYear=nowTime.getFullYear()
nowMonth=nowTime.getMonth()
nowDay=nowTime.getDay()
nowHour=nowTime.getHours()
nowMinute=nowTime.getMinutes()

// 发送指令
FB_SendMCCmd("scoreboard objectives add "+scoreBoardName+" dummy 时间记分板")
FB_SendMCCmd("scoreboard players set year "+scoreBoardName+" "+nowYear)
FB_SendMCCmd("scoreboard players set month "+scoreBoardName+" "+nowMonth)
FB_SendMCCmd("scoreboard players set day "+scoreBoardName+" "+nowDay)
FB_SendMCCmd("scoreboard players set hour "+scoreBoardName+" "+nowHour)
FB_SendMCCmd("scoreboard players set minute "+scoreBoardName+" "+nowMinute)

// 向用户发送提示信息
FB_Println("时间记分板校准完成！")