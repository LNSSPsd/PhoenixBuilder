package components

import (
	"encoding/json"
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/omega/defines"
	"phoenixbuilder/omega/utils"
	"strconv"
	"time"
)

type Territory struct {
	*BasicComponent
	KeyWord        map[string]string `json:"领地关键字"`
	Price          int               `json:"领地价格"`
	Range          []int             `json:"领地范围"`
	TerritoryNum   int               `json:"领地数量上限"`
	TerritoryScore string            `json:"领地价格积分榜"`
	IsCheckMod     bool              `json:"是否限制生存模式才能购买"`
	//IsTpProtect    bool                    `json:"是否开启传送地皮保护"`
	McTerriotyMsg string                  `json:"非法进入地皮提示信息"`
	Datas         map[string]*InTerritory //存储用户信息 的 sting是用户名字 interritory是相关信息

}

//太晚了就先写出个预备功能吧
//暂时办到购买和保护就是了
type InTerritory struct {
	//date   []string
	IsTpProtect bool //是否这个玩家开启地皮保护
	Pos         []int
	Range       []int
	Posx        string
	Posz        string
	Posy        string
	RangeX      string
	RangeZ      string
	Member      []string
}

//检测价格是否合适
func (b *Territory) Init(cfg *defines.ComponentConfig) {
	//解析并对应存储进入b
	m, _ := json.Marshal(cfg.Configs)
	err := json.Unmarshal(m, b)
	if err != nil {
		panic(err)
	}
	b.Datas = make(map[string]*InTerritory)
	///fmt.Print("读取")
	//Upos := []int{0, 0, 0}
	//Range := []int{0, 0, 0}
	//b.WriteUser("test", Upos, Range)

	//b.Frame.GetJsonData("地皮信息.json", &b.Datas)
	//fmt.Println("读取完毕")
}
func (b *Territory) Inject(frame defines.MainFrame) {
	//
	b.Frame = frame
	//注入frame等东西
	b.Frame.GetJsonData("地皮信息.json", &b.Datas)
	b.BasicComponent.Inject(frame)
	b.Listener.SetGameChatInterceptor(b.ProcessingCenter)
	//获取信息

}
func (b *Territory) Activate() {
	//分别执行两段函数

	b.Frame.GetGameControl().SayTo("@a", "地皮插件已开启\n输入地皮菜单获取相关指令信息")
	//保护区域
	//初始化分数
	b.Frame.GetGameControl().SendCmd("scoreboard players add @a " + b.TerritoryScore + " 0")
	go func() {
		fmt.Print("[地皮插件] 地皮保护已开启")
		for {
			//先遍历是否有地皮
			//获取所任玩家id然后如果名字不在遍历到的v对象里面 记住有两个 一个是name 一个是members 看是否这个名字在里面 如果不在 分别执行指令gamemode 2 @a[name=玩家名,x=,y=z=,dx=,dy=]这样就做到了保护了
			if len(b.Datas) >= 1 {
				b.preserver()
			}
		}

	}()

}
func (b *Territory) Stop() error {
	fmt.Print("正在保存[地皮信息]进入<地皮信息.json>")
	return b.Frame.WriteJsonData("地皮信息.json", &b.Datas)

}

//保护地皮函数
func (b *Territory) preserver() {
	//先遍历所有地皮开始保护
	for k, t := range b.Datas {
		for _, p := range b.Frame.GetUQHolder().PlayersByEntityID {

			//fmt.Print(p.Username)
			//判断是否在地皮名单中
			//if _, ok := b.Datas[p.Username]; ok {
			if k == p.Username {
				//不处理
			} else {
				NoMmember := "@a[name=\"" + p.Username + "\",x=" + t.Posx + ",z=" + t.Posz + ",y=0,m=0,dy=255,dx=" + t.RangeX + ",dz=" + t.RangeZ + "]"
				//判断是否在白名单中
				if b.CheckArr(t.Member, p.Username) {
					//在里面不做处理
				} else {

					//fmt.Print("保护中\n对象名字为", p.Username)
					msg := "gamemode 2 " + NoMmember
					b.Frame.GetGameControl().SayTo(NoMmember, "[地皮助手]你处于地皮范围\n主人:"+k+"\n地皮起始坐标为:"+t.Posx+","+t.Posz+"\n范围为 x轴延展:"+t.RangeX+"z轴延展:"+t.RangeZ)
					//fmt.Print(msg, "\n")
					b.Frame.GetGameControl().SayTo(NoMmember, b.McTerriotyMsg)
					b.Frame.GetGameControl().SendCmd(msg)

					//是否传送保护
					if t.IsTpProtect {
						tpx := t.Pos[0] + t.Range[0] + 1
						tpy := t.Pos[1] + t.Range[1] + 1
						_tpx := strconv.Itoa(tpx)
						_tpy := strconv.Itoa(tpy)
						b.Frame.GetGameControl().SendCmd("tp " + NoMmember + " " + _tpx + " " + _tpy)

					}
					time.Sleep(5000)
				}

			}

		}

	}
}

//检查某个东西是否在数组里面
func (b *Territory) CheckArr(arr []string, str string) (IsIn bool) {
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
func (b *Territory) WriteUser(Uname string, Upos []int, URange []int) bool {

	//fmt.Printf("Uname:%v,Upos:%v,Urange%v\n", Uname, Upos, URange)

	//检查是否重合
	for _, j := range b.Datas {
		if b.CheckIsoverlap(Upos, URange, j.Pos, j.Range) {
			b.Frame.GetGameControl().SayTo("@a[name=\""+Uname+"\"]", "§c[写入失败] §a你周围有地皮请远离后再购买")
			return false
		}
	}

	b.Datas[Uname] = &InTerritory{
		IsTpProtect: false,
		Pos:         Upos,
		Range:       URange,
		Posx:        strconv.Itoa(Upos[0]),
		Posz:        strconv.Itoa(Upos[2]),
		Posy:        strconv.Itoa(Upos[1]),

		RangeX: strconv.Itoa(URange[0]),
		RangeZ: strconv.Itoa(URange[1]),
		Member: make([]string, 99),
	}
	return true

}
func (b *Territory) DelectTerritory(Uname string) {
	//删除该人地皮
	if _, ok := b.Datas[Uname]; ok {
		delete(b.Datas, Uname)
		b.Frame.GetGameControl().SayTo("@a[name=\""+Uname+"\"]", "§e§l成功删除你的地皮")
	} else {
		b.Frame.GetGameControl().SayTo("@a[name=\""+Uname+"\"]", "§c你没有地皮")
	}

}
func (b *Territory) CheckIsoverlap(pos []int, Epos []int, Spos []int, SEpos []int) (IsOverlap bool) {
	//检查地皮是否重合
	x1 := pos[0]
	y1 := pos[2]
	x2 := pos[0] + Epos[0]
	y2 := pos[2] + Epos[1]
	x3 := Spos[0]
	y3 := Spos[2]
	x4 := Spos[0] + SEpos[0]
	y4 := Spos[2] + SEpos[1]

	if x1 <= x4 && x3 <= x2 && y1 <= y4 && y3 <= y2 {
		return true
	}
	return false
}
func (b *Territory) GiveMemmber(Uname string, Oname string) {
	//给予对方权限
	if _, ok := b.Datas[Uname]; ok {
		//fmt.Print(Oname, "\ntest")
		b.Datas[Uname].Member = append(b.Datas[Uname].Member, Oname)
		fmt.Println(b.Datas[Uname].Member, "成员")

		b.Frame.GetGameControl().SayTo("@a[name=\""+Uname+"\"]", "§e§l成功写入名单")
	} else {
		b.Frame.GetGameControl().SayTo("@a[name=\""+Uname+"\"]", "§c§l§o你没有地皮")
	}

}

//购买地皮主函数
func (b *Territory) BuyTerritory(Uname string) {
	name := utils.ToPlainName(Uname)
	//显示名字一下
	//fmt.Print(name, "显示名字\n")
	msgCmds := "scoreboard players list \"" + name + "\""
	price := strconv.Itoa(b.Price)
	//获取对象的积分并判断是否为生存 积分是否超过设定价格
	b.Frame.GetGameControl().SendCmdAndInvokeOnResponse(msgCmds, func(output *packet.CommandOutput) {
		if output.SuccessCount > 0 {

			for _, p := range output.OutputMessages[1:] {

				num, err := strconv.Atoi(p.Parameters[0])
				if err != nil {
					fmt.Print(err)
				}
				//fmt.Print("转化后数据", num, "\n")
				//fmt.Print("要求积分榜", b.TerritoryScore, "\n")
				//是否长度达标 是否积分为指定积分
				if len(p.Parameters) == 3 && (p.Parameters[2] == b.TerritoryScore) {
					//fmt.Print("扣除的计分板:", p.Parameters[2])
					if num >= b.Price {
						//b.WritePos(name,pos,range)
						go func() {
							pos := <-b.Frame.GetGameControl().GetPlayerKit(name).GetPos("@a[name=[player]]")
							if pos != nil {
								fmt.Print(pos)
							}
							//是否检查模式
							cmdTest := "scoreboard players test @a[name=\"" + name + "\",m=0]  " + b.TerritoryScore + " " + price + " *"
							cmdRemove := "scoreboard players remove @a[name=\"" + name + "\",m=0]  " + b.TerritoryScore + " " + price
							if b.IsCheckMod {
								b.Frame.GetGameControl().SendCmdAndInvokeOnResponse(cmdTest, func(output *packet.CommandOutput) {
									//fmt.Print(cmdRmove, " 指令返回:", output.OutputMessages)
									if output.SuccessCount > 0 {
										k := b.WriteUser(name, pos, b.Range)
										if k {
											b.Frame.GetGameControl().SendCmd(cmdRemove)
											b.Frame.GetGameControl().SayTo("@a[name=\""+name+"\"]", "§b[购买成功] §a本次消费"+b.TerritoryScore+price)
										}
										//写入用户信息

									} else {
										b.Frame.GetGameControl().SayTo("@a[name="+name+"]", "§c[购买失败] §e请在生存模式购买")
									}
								})
								//否则就直接购买
							} else {

								k := b.WriteUser(name, pos, b.Range)
								if k {
									b.Frame.GetGameControl().SendCmd("scoreboard players remove @a[name=\"" + name + "\"]  " + b.TerritoryScore + " " + price)
									b.Frame.GetGameControl().SayTo("@a[name="+name+"]", "§b[购买成功] §a本次消费"+b.TerritoryScore+price)
								}
							}
						}()
					} else {
						b.Frame.GetGameControl().SayTo("@a[name="+name+"]", "§c[余额不足] §b"+b.TerritoryScore+"需要达到"+price)
					}
					//fmt.Print(p.Parameters, "\n")
				}
			}
		}

	})

}
func (b *Territory) ProcessingCenter(entry *defines.GameChat) (stop bool) {
	Bstop := false
	//判断是否存在
	if len(entry.Msg) >= 1 {
		fmt.Println(entry.Msg)
		if entry.Msg[0] == b.KeyWord["购买地皮关键字"] {
			//如果是想要的信息就返回true
			b.BuyTerritory(entry.Name)
			Bstop = true
		} else if entry.Msg[0] == b.KeyWord["删除领地关键字"] {
			b.DelectTerritory(entry.Name)
			//如果不是想要信息就返回false让下一个组件接受
			Bstop = true
		} else if len(entry.Msg) >= 2 && entry.Msg[0] == b.KeyWord["赋予地皮白名单关键字"] && len(entry.Msg[1]) >= 1 {

			b.GiveMemmber(entry.Name, entry.Msg[1])
			Bstop = true

		} else if entry.Msg[0] == b.KeyWord["返回地皮关键字"] {
			if k, ok := b.Datas[entry.Name]; ok {

				b.Frame.GetGameControl().SendCmd("tp @a[name=\"" + entry.Name + "\"] " + k.Posx + " " + k.Posy + " " + k.Posz)
			} else {
				b.Frame.GetGameControl().SayTo("@a[name=\""+entry.Name+"\"]", "§c[返回失败] §a你没有地皮")
			}

		} else if entry.Msg[0] == b.KeyWord["查看白名单关键字"] {
			if _, InOk := b.Datas[entry.Name]; InOk {
				msg := "你的地皮成员名单:"
				for _, v := range b.Datas[entry.Name].Member {
					if v == "" {
						//如果为空则跳过
					} else {
						msg = msg + v + ","
					}

				}
				msg = msg + "\n"
				b.Frame.GetGameControl().SayTo("@a[name=\""+entry.Name+"\"]", msg)
			}

		} else if entry.Msg[0] == "地皮菜单" {
			//打印地皮指令相关菜单
			msg := "  		§b§l[地皮菜单相关快捷指令]"

			for i, j := range b.KeyWord {
				msg = msg + "\n§e§l" + i + ": §a§l" + j + "\n"
			}
			b.Frame.GetGameControl().SayTo("@a[name=\""+entry.Name+"\"]", msg)
		} else {
			//如果不是想要信息就返回false让下一个组件接受
			Bstop = false
		}
	}
	return Bstop
}
