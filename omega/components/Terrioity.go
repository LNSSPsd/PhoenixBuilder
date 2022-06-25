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
	KeyWord        map[string]string       `json:"领地关键字"`
	Price          int                     `json:"领地价格"`
	Range          []int                   `json:"领地范围"`
	TerritoryNum   int                     `json:"领地数量上限"`
	TerritoryScore string                  `json:"领地价格积分榜"`
	IsCheckMod     bool                    `json:"是否限制生存模式才能购买"`
	IsTpProtect    bool                    `json:"是否开启传送地皮保护"`
	Datas          map[string]*InTerritory //存储用户信息 的 sting是用户名字 interritory是相关信息

}

//太晚了就先写出个预备功能吧
//暂时办到购买和保护就是了
type InTerritory struct {
	//date   []string
	Pos    []int
	Range  []int
	Posx   string
	Posz   string
	RangeX string
	RangeZ string
	Member []string
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

	b.Frame.GetGameControl().SayTo("@a", "地皮插件已开启")
	//保护区域
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

//保护地皮
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
				//判断是否在白名单中
				if b.CheckArr(t.Member, p.Username) {
					//在里面不做处理
				} else {

					//fmt.Print("保护中\n对象名字为", p.Username)
					msg := "gamemode 2 @a[name=\"" + p.Username + "\",x=" + t.Posx + ",z=" + t.Posz + ",y=0,m=0,dy=255,dx=" + t.RangeX + ",dz=" + t.RangeZ + "]"
					b.Frame.GetGameControl().SayTo("@a[name=\""+p.Username+"\",x="+t.Posx+",z="+t.Posz+",y=0,m=0,dy=255,dx="+t.RangeX+",dz="+t.RangeZ+"]", "[地皮助手]你处于地皮范围\n主人:"+k+"\n地皮起始坐标为:"+t.Posx+","+t.Posz+"\n范围为 x轴延展:"+t.RangeX+"z轴延展:"+t.RangeZ)
					//fmt.Print(msg, "\n")
					b.Frame.GetGameControl().SendCmd(msg)

					//是否传送保护
					if b.IsTpProtect {
						tpx := t.Pos[0] + t.Range[0] + 1
						tpy := t.Pos[1] + t.Range[1] + 1
						_tpx := strconv.Itoa(tpx)
						_tpy := strconv.Itoa(tpy)
						b.Frame.GetGameControl().SendCmd("tp @a[name=\"" + p.Username + "\",x=" + t.Posx + ",z=" + t.Posz + ",y=0,m=0,dy=255,dx=" + t.RangeX + ",dz=" + t.RangeZ + "]" + _tpx + " " + _tpy)

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
func (b *Territory) WriteUser(Uname string, Upos []int, URange []int) {
	//写入用户信息
	fmt.Print("成功写入信息")
	//b.Frame.WriteJsonData("地皮信息.json",
	//fmt.Print("upos:", Upos, "range:", URange, "uname", Uname, "\n")
	//
	b.Datas[Uname] = &InTerritory{

		Pos:   Upos,
		Range: URange,
		Posx:  strconv.Itoa(Upos[0]),
		Posz:  strconv.Itoa(Upos[2]),

		RangeX: strconv.Itoa(URange[0]),
		RangeZ: strconv.Itoa(URange[1]),
		Member: make([]string, 99),
	}

}
func (b *Territory) DelectTerritory(Uname string) {
	//删除该人地皮
	if _, ok := b.Datas[Uname]; ok {
		delete(b.Datas, Uname)
		b.Frame.GetGameControl().SayTo("@a[name=\""+Uname+"\"]", "成功删除你的地皮")
	} else {
		b.Frame.GetGameControl().SayTo("@a[name=\""+Uname+"\"]", "你没有地皮")
	}

}
func (b *Territory) GiveMemmber(Uname string, Oname string) {
	//给予对方权限
	if _, ok := b.Datas[Uname]; ok {
		//fmt.Print(Oname, "\ntest")
		b.Datas[Uname].Member = append(b.Datas[Uname].Member, Oname)
		fmt.Println(b.Datas[Uname].Member, "成员")

		b.Frame.GetGameControl().SayTo("@a[name=\""+Uname+"\"]", "成功写入名单")
	} else {
		b.Frame.GetGameControl().SayTo("@a[name=\""+Uname+"\"]", "你没有地皮")
	}

}

//购买地皮主函数
func (b *Territory) BuyTerritory(Uname string) {
	name := utils.ToPlainName(Uname)
	//显示名字一下
	fmt.Print(name, "显示名字\n")
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
					fmt.Print("扣除的计分板:", p.Parameters[2])
					if num >= b.Price {
						//b.WritePos(name,pos,range)
						go func() {
							pos := <-b.Frame.GetGameControl().GetPlayerKit(name).GetPos("@a[name=[player]]")
							if pos != nil {
								fmt.Print(pos)
							}
							//是否检查模式
							cmdRmove := "scoreboard players remove @a[name=\"" + name + "\",m=0]  " + b.TerritoryScore + " " + price
							if b.IsCheckMod {
								b.Frame.GetGameControl().SendCmdAndInvokeOnResponse(cmdRmove, func(output *packet.CommandOutput) {
									//fmt.Print(cmdRmove, " 指令返回:", output.OutputMessages)
									if output.SuccessCount > 0 {
										b.Frame.GetGameControl().SayTo("@a[name="+name+"]", "[购买成功] 本次消费"+b.TerritoryScore+price)
										//写入用户信息
										b.WriteUser(name, pos, b.Range)
									} else {
										b.Frame.GetGameControl().SayTo("@a[name="+name+"]", "[购买失败] 请在生存模式购买")
									}
								})
								//否则就直接购买
							} else {
								b.Frame.GetGameControl().SendCmd(cmdRmove)
								b.Frame.GetGameControl().SayTo("@a[name="+name+"]", "[购买成功] 本次消费"+b.TerritoryScore+price)
								b.WriteUser(name, pos, b.Range)
							}
						}()
					} else {
						b.Frame.GetGameControl().SayTo("@a[name="+name+"]", "[余额不足] "+b.TerritoryScore+"需要达到"+price)
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

		} else {
			//如果不是想要信息就返回false让下一个组件接受
			Bstop = false
		}

	}
	return Bstop
}
