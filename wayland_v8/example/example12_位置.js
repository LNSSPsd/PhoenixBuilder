FB_WaitConnect()

// 获得机器人位置 pos.x pos.y pos.z
pos = FB_GetBotPos()
FB_Println(JSON.stringify(pos))

// 移动
FB_SendMCCmdAndGetResult("tp @s 100 200 300")
pos = FB_GetBotPos()
FB_Println(JSON.stringify(pos))