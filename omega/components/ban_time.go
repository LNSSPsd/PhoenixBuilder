package components

import (
	"encoding/json"
	"fmt"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/omega/defines"
	"phoenixbuilder/omega/utils"
	"sync"
	"time"

	"github.com/pterm/pterm"
)

type BanTime struct {
	*defines.BasicComponent
	Selector             string      `json:"选择器"`
	Duration             float64     `json:"检查周期"`
	OnOmegaTakeOverIn    interface{} `json:"Omega接管时指令"`
	KickCmdIn            interface{} `json:"踢出指令"`
	AfterOmegaTakeOverIn interface{} `json:"到达封禁时间Omega结束接管时指令"`
	ScoreboardName       string      `json:"读取封禁时间的计分板名"`
	FileName             string      `json:"文件名"`
	KickDelay            int         `json:"延迟踢出时间"`
	LoginDelay           int         `json:"登录时延迟发送"`
	OnOmegaTakeOver      []defines.Cmd
	KickCmd              []defines.Cmd
	AfterOmegaTakeOver   []defines.Cmd
	banTime              map[string]time.Time
	banTimeStr           map[string]string
	fileChange           bool
	mu                   sync.Mutex
}

func (o *BanTime) Init(cfg *defines.ComponentConfig) {
	m, _ := json.Marshal(cfg.Configs)
	var err error
	if err = json.Unmarshal(m, o); err != nil {
		panic(err)
	}
	if o.OnOmegaTakeOver, err = utils.ParseAdaptiveCmd(o.OnOmegaTakeOverIn); err != nil {
		panic(err)
	}
	if o.KickCmd, err = utils.ParseAdaptiveCmd(o.KickCmdIn); err != nil {
		panic(err)
	}
	if o.AfterOmegaTakeOver, err = utils.ParseAdaptiveCmd(o.AfterOmegaTakeOverIn); err != nil {
		panic(err)
	}
	o.mu = sync.Mutex{}
	o.banTime = make(map[string]time.Time)
	o.banTimeStr = make(map[string]string)
}

func (o *BanTime) Inject(frame defines.MainFrame) {
	o.Frame = frame
	err := frame.GetJsonData(o.FileName, &o.banTimeStr)
	if err != nil {
		panic(err)
	}
	for k, v := range o.banTimeStr {
		o.banTime[k], err = utils.StringToTimeWithLocal(v + " +0800 CST")
		if err != nil {
			panic(err)
		}
	}
	o.Frame.GetGameListener().AppendLoginInfoCallback(func(entry protocol.PlayerListEntry) {
		o.mu.Lock()
		defer o.mu.Unlock()
		if banTime, hasK := o.banTime[entry.Username]; hasK {
			if banTime.After(time.Now()) {
				go func() {
					<-time.NewTimer(time.Duration(o.KickDelay) * time.Second).C
					o.kick(entry.Username)
				}()
			} else {
				delete(o.banTime, entry.Username)
				delete(o.banTimeStr, entry.Username)
				o.fileChange = true
				go func() {
					<-time.NewTimer(time.Duration(o.LoginDelay) * time.Second).C
					utils.LaunchCmdsArray(o.Frame.GetGameControl(), o.AfterOmegaTakeOver, map[string]interface{}{
						"[player]": entry.Username,
					}, o.Frame.GetBackendDisplay())
				}()
			}
		}
	})
}

func (o *BanTime) Signal(signal int) (err error) {
	switch signal {
	case defines.SIGNAL_DATA_CHECKPOINT:
		if o.fileChange {
			err = o.Frame.WriteJsonDataWithTMP(o.FileName, ".ckpt", o.banTimeStr)
			o.fileChange = false
		}
	}
	return err
}

func (o *BanTime) Stop() error {
	fmt.Printf("正在保存 %v\n", o.FileName)
	return o.Frame.WriteJsonDataWithTMP(o.FileName, ".final", o.banTimeStr)
}

func (o *BanTime) kick(name string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if banTime, hasK := o.banTime[name]; hasK {
		duration := time.Until(banTime)
		go utils.LaunchCmdsArray(o.Frame.GetGameControl(), o.KickCmd, map[string]interface{}{
			"[player]": name,
			"[day]":    int(duration.Hours()) / 24,
			"[hour]":   int(duration.Hours()) % 24,
			"[min]":    int(duration.Minutes()) % 60,
			"[sec]":    int(duration.Seconds()) % 60,
		}, o.Frame.GetBackendDisplay())
	}
}

func (o *BanTime) takeOver(name string) {
	utils.GetPlayerScore(o.Frame.GetGameControl(), name, o.ScoreboardName, func(val int, err error) {
		if err != nil {
			pterm.Error.Printfln("无法获取封禁时间信息 %v %v %v", name, o.ScoreboardName, err)
		} else if val < 0 {
			pterm.Error.Printfln("封禁时间指令设计配置有问题，如果封禁时间小于等于 0，则不应该被选择器选中 %v %v", name, o.ScoreboardName)
		} else {
			duration := time.Second * time.Duration(val)
			banTime := time.Now().Add(duration)
			o.mu.Lock()
			o.banTime[name] = banTime
			o.banTimeStr[name] = utils.TimeToString(banTime)
			o.mu.Unlock()
			o.fileChange = true
			go func() {
				utils.LaunchCmdsArray(o.Frame.GetGameControl(), o.OnOmegaTakeOver, map[string]interface{}{
					"[player]": name,
					"[day]":    int(duration.Hours()) / 24,
					"[hour]":   int(duration.Hours()) % 24,
					"[min]":    int(duration.Minutes()) % 60,
					"[sec]":    int(duration.Seconds()) % 60,
				}, o.Frame.GetBackendDisplay())
				o.kick(name)
			}()
		}
	})
}
func (o *BanTime) Activate() {
	t := time.NewTicker(time.Second * time.Duration(o.Duration))
	for {
		<-t.C
		o.Frame.GetGameControl().SendCmdAndInvokeOnResponse(fmt.Sprintf("/testfor %v", o.Selector), func(output *packet.CommandOutput) {
			if output.SuccessCount > 0 && len(output.OutputMessages) > 0 {
				ban := &Banned{Victim: []string{}}
				err := json.Unmarshal([]byte(output.DataSet), &ban)
				if err != nil {
					o.Frame.GetBackendDisplay().Write(fmt.Sprintf("fail to get kick info " + err.Error()))
				} else {
					o.Frame.GetBackendDisplay().Write(fmt.Sprintf("try to kick %v", ban.Victim))
					for _, v := range ban.Victim {
						o.takeOver(v)
					}
				}
			}
		})
	}
}
