package luaComponent

import "time"

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

//封装游戏内包
func NewGamePackage(packageId string, bodyString string) ConnectPackage {
	//
	//to do 将packageid与bodystring组装成json进入Body
	//

	gamePackage := ConnectPackage{
		TypeName: "GamePackage",
		Body:     bodyString,
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
