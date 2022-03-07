// 本脚演示了玩家无法进入租赁服时，使用脚本紧急执行指令
// 演示了 engine.waitConnectSync 的一个 **重要特性**
// 该脚本必须以启动脚本的方式运行，即:
// ./fastbuilder -S 脚本路径/example02_紧急救援.js

engine.setName("紧急救援")

// 询问用户是否要停用所有命令方块
let choose = engine.questionSync("是否要紧急停止所有命令方块? 输入y:")
if (choose === "y") {
    // 当作为启动脚本运行时，FB将暂停连接到MC，直到 FB_WaitConnect 被调用
    engine.waitConnectionSync()
    // 在连接到MC后，立刻发送指令
    game.sendCommand("gamerule commandblocksenabled false")
    // 向用户发送提示信息
    engine.message("时间记分板校准完成！")
} else {
    engine.message("好的吧")
}