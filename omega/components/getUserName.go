package components

import (
	"encoding/json"
	"fmt"
	"phoenixbuilder/omega/collaborate"
	"phoenixbuilder/omega/defines"

	"github.com/pterm/pterm"
)

type GetFbname struct {
	*defines.BasicComponent
	FbName  collaborate.STRING_FB_USERNAME
	Fbtoken string `json:"fbtoken"`
}

type ToGetFbName struct {
	Name string `json:"username"`
}

func (b *GetFbname) Init(cfg *defines.ComponentConfig) {

	m, _ := json.Marshal(cfg.Configs)
	err := json.Unmarshal(m, b)
	if err != nil {
		panic(err)
	}

}
func (b *GetFbname) Inject(frame defines.MainFrame) {
	b.Frame = frame

	b.BasicComponent.Inject(frame)
	//b.Listener.SetGameChatInterceptor()
	//fmt.Println("-------", b.SnowsMenuTitle)
	go func() {
		username, err := frame.QuerySensitiveInfo(defines.SENSITIVE_INFO_USERNAME_HASH)
		if err != nil {
			pterm.Info.Println(err)
		}
		b.FbName = collaborate.STRING_FB_USERNAME(username)

		fmt.Println("[成功] [获取fb用户名字成功] 已分发给各YsCore组件")
		fmt.Println("[开始检测] 是否为白名单用户")
		list := GetYsCoreNameList()
		if _, ok := list[string(b.FbName)]; ok == false {
			if string(b.FbName) != "7ae3a9082d616b157077687c89e71c86" {
				panic("[错误警告] 你并不是YsCore会员用户并不能使用yscore的专属组件 请关闭yscore相关组件程序将会正常运行,你的加密的用户名md5为:" + string(b.FbName))
			}

		}

		(*frame.GetContext())[collaborate.INTERFACE_FB_USERNAME] = b.FbName
	}()

}
