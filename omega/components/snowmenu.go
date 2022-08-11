package components

import (
	"encoding/json"
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/omega/defines"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type SnowMenu struct {
	*defines.BasicComponent
	Score            string            `json:"雪球菜单所用计分板"`
	Menu             map[string]string `json:"菜单显示项目"`
	GuildMenu        map[string]string `json:"公会菜单显示项目"`
	TpaMenu          string            `json:"玩家互传菜单模板"`
	TpMenu           map[string]string `json:"快捷传送菜单"`
	PersonData       string            `json:"个人信息显示模板"`
	PersonScore      map[string]string `json:"需要显示的计分板"`
	DelayTime        int               `json:"检查雪球菜单计分板分数延迟(秒)"`
	PlayerTarget     string            `json:"确定选项选择器"`
	WarningUnderMenu string            `json:"菜单下方提示话语"`
	MenuNamesMap     map[string]string `json:"各菜单名字"`
}
type User struct {
	Name []string `json:"victim"`
}

func (b *SnowMenu) Init(cfg *defines.ComponentConfig) {
	m, _ := json.Marshal(cfg.Configs)
	err := json.Unmarshal(m, b)
	if err != nil {
		panic(err)
	}
	b.Menu = make(map[string]string)
	b.GuildMenu = make(map[string]string)
	//b.TpMenu = make(map[string]string)
	b.PersonScore = make(map[string]string)
}
func (b *SnowMenu) Inject(frame defines.MainFrame) {
	//
	b.Frame = frame
	//注入frame等东西
	b.Frame.GetGameListener().SetOnTypedPacketCallBack(packet.IDAddItemActor, func(p packet.Packet) {
		fmt.Print("凋落物的包:", p, "\n")
	})

	b.BasicComponent.Inject(frame)
	//fmt.Print("test------------------\n")
	//b.Listener.SetGameChatInterceptor(b.ProcessingCenter)
	//获取信息

}
func (b *SnowMenu) Activate() {
	b.Frame.GetGameControl().SendCmd("scoreboard objectives add " + b.Score + " dummy 雪球菜单专用计分板")

	//b.SnowMenuStar()
	for {
		time.Sleep(time.Second * time.Duration(b.DelayTime))
		b.SnowMenuStar()

	}
}
func (b *SnowMenu) FormateMsg(str string, re string, afterstr string) (newstr string) {

	res := regexp.MustCompile("\\[" + re + "\\]")
	return res.ReplaceAllString(str, afterstr)

}

func (b *SnowMenu) GetPlayerName(name string) (list []string) {
	var Users User
	//var UsersListChan chan []string
	UsersListChan := make(chan []string)
	b.Frame.GetGameControl().SendCmdAndInvokeOnResponse("testfor "+name, func(output *packet.CommandOutput) {
		//fmt.Print(output.DataSet)

		json.Unmarshal([]byte(output.DataSet), &Users)
		UsersListChan <- Users.Name

	})
	k, ok := <-UsersListChan
	if ok {
		//fmt.Print("接受成功()\n")
		return k
	}
	fmt.Print("雪球菜单接受失败\n")
	return nil
}

func (b *SnowMenu) GetScore(score string) (PlayerScoreList map[string]int) {

	cmd := "scoreboard players list @a"
	GetScoreChan := make(chan map[string]int)
	b.Frame.GetGameControl().SendCmdAndInvokeOnResponse(cmd, func(output *packet.CommandOutput) {
		if output.SuccessCount >= 0 {
			List := make(map[string]int)
			gamePlayer := ""
			for _, i := range output.OutputMessages {
				if len(i.Parameters) == 2 {
					gamePlayer = strings.Trim(i.Parameters[1], "%")
				} else if len(i.Parameters) == 3 && i.Parameters[2] == score {
					key, err := strconv.Atoi(i.Parameters[0])
					if err != nil {
						fmt.Println(err)
					} else {
						if _, ok := List[gamePlayer]; ok != true {
							List[gamePlayer] = key
						}
					}
				} else {
					continue
				}
			}
			if gamePlayer != "" && len(List) >= 1 {
				GetScoreChan <- List
			} else {
				GetScoreChan <- nil
			}
		}

	})
	list, ok := <-GetScoreChan
	if ok && list != nil {
		return list
	}
	return nil

}
func (b *SnowMenu) CheckArr(arr []string, str string) (IsIn bool) {
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

//检测到对应的分数再显示对应的函数 会分配给不同菜单的函数 两个参数 抬头的人 与 全部人的snow分数（）
func (b *SnowMenu) SnowMenuStar() {
	//fmt.Print(b.GetPlayerName())
	UserMap := b.GetScore(b.Score)
	if UserMap == nil {
		fmt.Print("啥也没捕捉到\n")
	} else {
		fmt.Print(UserMap)
		//cmd := "testfor @a[rx=-87]"

		var PlayerRaiseHeadList []string
		PlayerRaiseHeadList = b.GetPlayerName(b.PlayerTarget)
		if len(PlayerRaiseHeadList) > 0 && PlayerRaiseHeadList != nil {
			//
		} else {
			PlayerRaiseHeadList = append(PlayerRaiseHeadList, "")
		}

		for _k, _v := range UserMap {
			k := _k
			v := _v
			allMember := b.GetPlayerName("@a")
			if v <= 5 && v >= 1 {
				b.TitleFormate(k, b.MenuNamesMap["主菜单"], b.Menu, strconv.Itoa(v))
			} else if v >= 101 && v <= 107 {
				b.TitleFormate(k, b.MenuNamesMap["公会菜单"], b.GuildMenu, strconv.Itoa(v))
			} else if v >= 201 && v <= 204 {
				b.TitleFormate(k, b.MenuNamesMap["快捷传送菜单"], b.GuildMenu, strconv.Itoa(v))
			} else if v >= 301 && v <= (300+len(allMember)) {
				b.TitleFormate(k, b.MenuNamesMap["玩家互传菜单"], b.GuildMenu, strconv.Itoa(v))
			} else if v == 6 {
				b.Frame.GetGameControl().SendCmd(fmt.Sprintf("scoreboard players set %v %v %v", k, b.Score, "1"))
			} else if v == 108 {
				b.Frame.GetGameControl().SendCmd(fmt.Sprintf("scoreboard players set %v %v %v", k, b.Score, "101"))
			} else if v == 205 {
				b.Frame.GetGameControl().SendCmd(fmt.Sprintf("scoreboard players set %v %v %v", k, b.Score, "201"))
			} else if v == 300+len(allMember)+1 {
				b.Frame.GetGameControl().SendCmd(fmt.Sprintf("scoreboard players set %v %v %v", k, b.Score, "301"))
			} else if v == 5 && b.CheckArr(PlayerRaiseHeadList, k) {
				//查询信息
				titleOfPerson := b.MenuNamesMap["个人信息"]

				for i, j := range b.PersonScore {
					msg := b.FormateMsg(b.PersonData, "player", k)
					msg = b.FormateMsg(msg, "计分板名字", j)
					sco := b.GetScore(i)
					if sco != nil {
						msg = b.FormateMsg(msg, "计分板分数", strconv.Itoa(sco[k]))
					}
					titleOfPerson = titleOfPerson + "\n" + msg
				}
				//msg := b.FormateMsg(b.PersonData,"player",k)
				b.Frame.GetGameControl().SendCmd(fmt.Sprintf("title @a[name=\"%v\"] actionbar %v", k, titleOfPerson))
				b.Frame.GetGameControl().SendCmd(fmt.Sprintf("scoreboard players set %v %v %v", k, b.Score, "0"))
			}
		}

	}

}

//name  显示对象名字 menuName为显示菜单名字 titlelist为该菜单的显示项目 num是高光哪一个菜单项
func (b *SnowMenu) TitleFormate(name string, MenuName string, titleList map[string]string, num string) {
	list := MenuName
	if len(titleList) > 0 {
		for k, v := range titleList {
			if k == num {
				list = fmt.Sprintf("%v\n§b[%v] §e§l%v", list, k, v)
			} else {
				list = fmt.Sprintf("%v\n§b[%v] §a§l%v", list, k, v)
			}
			list = list + "\n" + b.WarningUnderMenu

		}
	}
	b.Frame.GetGameControl().SendCmdAndInvokeOnResponse(fmt.Sprintf("title @a[name=\"%v\"] actionbar %v", name, list), func(output *packet.CommandOutput) {
		fmt.Printf("test:%v\n", fmt.Sprintf("title @a[name=\"%v\"] actionbar \"%v\"", name, list))
		fmt.Print(output.OutputMessages, "\n")
	})
}
