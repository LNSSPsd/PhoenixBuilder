package components

import (
	"encoding/json"
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/omega/defines"
	"time"
)

type Downplayer struct {
	*BasicComponent
	Datas   map[string]*Information
	KeyWord string   `json:"操控玩家分数关键字"`
	Master  []string `json:"管理员"`
}

//对应计分板对应积分
type Information struct {
	score map[string]string
}

func (o *Downplayer) Init(cfg *defines.ComponentConfig) {
	m, _ := json.Marshal(cfg.Configs)
	err := json.Unmarshal(m, o)
	if err != nil {
		panic(err)
	}

}
func (o *Downplayer) CheckPlayer() {
	for _, p := range o.Frame.GetUQHolder().PlayersByEntityID {
		//遍历所有人
		if k, ok := o.Datas[p.Username]; ok {
			if len(k.score) >= 1 {
				for i, j := range k.score {
					msg := "scoreboard players remove @a[name=\"" + p.Username + "\"]  " + i + " " + j
					o.Frame.GetGameControl().SendCmdAndInvokeOnResponse(msg, func(output *packet.CommandOutput) {
						if output.SuccessCount > 0 {

							delete(o.Datas, p.Username)
						}
					})
				}

			}
		}
	}

}
func (o *Downplayer) Inject(frame defines.MainFrame) {
	o.Frame = frame
	o.BasicComponent.Inject(frame)
	o.Datas = make(map[string]*Information)
	o.Frame.GetJsonData("玩家下线操作.json", &o.Datas)
	o.Listener.SetGameChatInterceptor(o.ProcessingCenter)

}
func (o *Downplayer) Activate() {
	fmt.Print("已开启下线玩家操作组件")
	go func() {
		for {
			o.CheckPlayer()
			time.Sleep(5000)
		}

	}()
}
func (o *Downplayer) Stop() error {
	fmt.Print("保存玩家下线操作.json中")
	return o.Frame.WriteJsonData("玩家下线操作.json", &o.Datas)
}
func (o *Downplayer) ProcessingCenter(entry *defines.GameChat) (stop bool) {
	if len(entry.Msg) >= 1 {
		if entry.Msg[0] == o.KeyWord && len(entry.Msg) == 4 && o.CheckArr(o.Master, entry.Name) {
			name := entry.Msg[1]
			score := entry.Msg[2]
			num := entry.Msg[3]
			fmt.Print(entry.Name)
			o.WriteInfo(entry.Name, name, score, num)
			return true
		}

	}
	return false
}
func (o *Downplayer) WriteInfo(source string, name string, Uscore string, num string) {
	if _, ok := o.Datas[name]; ok {
		//o.Datas[name].score = make(map[string]string)
		o.Datas[name].score[Uscore] = num
		o.Frame.GetGameControl().SayTo("@a[name=\""+source+"\"]", "已经改变分数等待下次玩家登录执行")

	} else {
		fmt.Print("test")
		o.Datas[name] = &Information{
			score: make(map[string]string),
		}
		o.Datas[name].score[Uscore] = num
		o.Frame.GetGameControl().SayTo("@a[name=\""+source+"\"]", "已经改变分数等待下次玩家登录执行")
	}
}
func (o *Downplayer) CheckArr(arr []string, str string) (IsIn bool) {
	if len(arr) == 0 {
		fmt.Print("数组为空")
		return false
	} else {
		var set map[string]struct{}
		set = make(map[string]struct{})
		for _, value := range arr {
			set[value] = struct{}{}
		}
		// 检查元素是否在map
		if _, ok := set[str]; ok {
			return true
		} else {
			return false
		}

	}

}
