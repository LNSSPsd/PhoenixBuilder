cq-chatlogger的一些简单说明
========================

---
> cq-chatlogger依赖于phoenixBuilder和负责发送消息的客户端. 主要用于qq和游戏间消息通信, 所以客户端一般使用 [go-cqhttp].
>
> 可以选择自定义一个客户端. 具体方法和示例将在文末说明.

[go-cqhttp]: 请移步https://github.com/Mrs4s/go-cqhttp

### 总体的结构

| 文件名            | 作用                        |
| :------:         | ----                       |
| chatlogger.go    |用来建立通信
| config.go        | 读取config.yml
| config.yml       |对chatlogger的各种功能进行限制和细分
| event.go         |cq协议端发来的消息类型主要有event和message. event一般有加群请求、 添加删除好友的事件, 在此我们不特意关注.|
| scoreboard.go    |建立和更新mc->qq消息的玩家过滤表. <font color=red>可能将要废弃.</font>     |
| scoreboard.json  |如上所述. <font color=red>可能将要废弃.</font>        |
| uiil.go          |~~把不知道该放哪的东西都放在这里面~~            |

---

### 工作原理

+ 建立通信
    1. cq-chatlogger(以下简称server)建立websocket正向通信(作为服务端)
    2. 理所当然地客户端进行连接.
    3. 如果客户端是go-cqhttp, 建立连接时会发一个lifecycle的post_type用于说明这个时候已经建立好连接,可以开始愉快地通信了.

+ server接收消息
    1. server会先尝试获取post_type. 此时一般有三种情况.
        - event: 一般是qq的各种杂七杂八的事件.
        - message: 私聊消息或群消息.
        - meta_event: 貌似只会在刚开始建立通信时收到.
    2. 解析message
        - 按照message_type分为private和group. 其他情况则为universal.
        - 然后按照每种消息的对应规则发消息.
        - 用tellraw发到游戏里, 用游戏中的选择器过滤不允许接收消息的玩家.

+ server发送消息
    1. todo