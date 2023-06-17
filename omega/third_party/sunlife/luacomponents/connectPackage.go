package luaComponent

import "time"

const (
	//主要是描述执行命令 比如重载插件之类的
	COMMANDTYPE = "commmandpackage"
	//绑定包的时候用
	RESOURCETYPE = "resourcepackage"
	//行为包 玩家行为之类的
	BEHAVIORTYPE = "behaviorpackage"
	//游戏指令包 主要是给游戏发送指令
	GAMETYPE = "gamecommandpackage"
)

//lua程序与交互器之间的包规范
type ConnectPackage struct {
	//包名
	TypeName string
	//发包主内容 内容主要是json
	Body string
	//发包日期
	Date int64
}

//包的主体
type PackageBody struct {
}
type GameBody struct {
}

func NewGamePackage(packageId string, bodystring string) ConnectPackage {
	//
	//to do 将packageid与bodystring组装成json进入Body
	//

	gamePackage := ConnectPackage{
		TypeName: "GamePackage",
		Body:     bodystring,
		Date:     time.Now().Unix(),
	}
	return gamePackage
}

//封装命令包
func NewCommandPackage(commandType string, commandBody string) ConnectPackage {
	//
	//to do
	//
	commandPackage := ConnectPackage{
		TypeName: "CommandPackage",
		Body:     commandBody,
		Date:     time.Now().Unix(),
	}
	return commandPackage
}
