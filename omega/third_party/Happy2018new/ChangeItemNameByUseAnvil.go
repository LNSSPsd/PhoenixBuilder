package Happy2018new

import (
	"encoding/json"
	"fmt"
	"math"
	blockNBT_API "phoenixbuilder/fastbuilder/bdump/blockNBT/API"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/omega/defines"

	"github.com/pterm/pterm"
)

type ChangeItemNameByUseAnvil struct {
	*defines.BasicComponent
	apis     blockNBT_API.GlobalAPI
	Triggers []string `json:"菜单触发词"`
	Usage    string   `json:"菜单项描述"`
	FilePath string   `json:"从何处提取物品的新名称(填写路径)"`
}

func (o *ChangeItemNameByUseAnvil) Init(settings *defines.ComponentConfig, storage defines.StorageAndLogProvider) {
	marshal, _ := json.Marshal(settings.Configs)
	if err := json.Unmarshal(marshal, o); err != nil {
		panic(err)
	}
}

func (o *ChangeItemNameByUseAnvil) Inject(frame defines.MainFrame) {
	o.Frame = frame
	o.apis = blockNBT_API.GlobalAPI{
		WritePacket: func(p packet.Packet) error {
			o.Frame.GetGameControl().SendMCPacket(p)
			return nil
		},
		BotName:            o.Frame.GetUQHolder().GetBotName(),
		BotIdentity:        "",
		BotUniqueID:        o.Frame.GetUQHolder().BotUniqueID,
		BotRunTimeID:       o.Frame.GetUQHolder().BotRuntimeID,
		PacketHandleResult: o.Frame.GetNewUQHolder(),
	}
	o.Frame.GetGameListener().SetGameMenuEntry(&defines.GameMenuEntry{
		MenuEntry: defines.MenuEntry{
			Triggers:     o.Triggers,
			FinalTrigger: false,
			Usage:        o.Usage,
		},
		OptionalOnTriggerFn: o.ChangeItemName,
	})
}

func (o *ChangeItemNameByUseAnvil) ChangeItemName(chat *defines.GameChat) bool {
	go func() {
		o.apis.BotName = o.Frame.GetUQHolder().GetBotName()
		// 初始化
		datas, err := o.Frame.GetFileData(o.FilePath)
		if err != nil {
			o.Frame.GetGameControl().SayTo(chat.Name, fmt.Sprintf("§c无法打开 §bomega_storage/data/%v §c处的文件\n详细日志已发送到控制台", o.FilePath))
			pterm.Error.Printf("修改物品名称: %v\n", err)
			return
		}
		if len(datas) <= 0 {
			o.Frame.GetGameControl().SayTo(chat.Name, fmt.Sprintf("§bomega_storage/data/%v §c处的文件没有填写物品名称§f，§c可能这个文件是个空文件§f，§c也可能是文件本身不存在", o.FilePath))
			return
		}
		itemName := string(datas)
		// 获取物品的新名称
		itemDatas, err := o.apis.PacketHandleResult.Inventory.GetItemStackInfo(0, 0)
		if err != nil {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c在读取快捷栏 §b0 §c时发送了错误\n详细日志已发送到控制台")
			pterm.Error.Printf("修改物品名称: %v\n", err)
			return
		}
		if itemDatas.Stack.ItemType.NetworkID == 0 {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c请确保机器人在快捷栏 §b0 §c有一个物品\n详细日志已发送到控制台")
			pterm.Warning.Printf("修改物品名称: itemDatas = %#v\n", itemDatas)
			return
		}
		// 确定被改名物品存在
		cmdResp, err := o.apis.SendWSCommandWithResponce("querytarget @s")
		if err != nil {
			panic(pterm.Error.Sprintf("修改物品名称: %v", err))
		}
		parseAns, err := o.apis.ParseQuerytargetInfo(cmdResp)
		if err != nil {
			panic(pterm.Error.Sprintf("修改物品名称: %v", err))
		}
		if len(parseAns) <= 0 {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c机器人可能没有 §bOP §c权限")
			return
		}
		pos := [3]int32{
			int32(math.Floor(float64(parseAns[0].Position[0]))),
			int32(math.Floor(float64(parseAns[0].Position[1]))),
			int32(math.Floor(float64(parseAns[0].Position[2]))),
		}
		// 取得机器人当前的坐标
		successStates, err := o.apis.ChangeItemNameByUsingAnvil(
			pos,
			`["direction": 0, "damage": "undamaged"]`,
			[]blockNBT_API.AnvilChangeItemName{
				{
					Slot: 0,
					Name: itemName,
				},
			},
			true,
		)
		if err != nil {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c物品名称修改失败\n详细日志已发送到控制台")
			pterm.Error.Printf("修改物品名称: %v\n", err)
			return
		}
		if successStates[0] == false {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c物品名称修改失败§f，§c请检查新的名称是否与原始名称相同")
			return
		}
		// 修改物品名称
		newItemDatas, err := o.apis.PacketHandleResult.Inventory.GetItemStackInfo(0, 0)
		if err != nil {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c在读取快捷栏 §b0 §c时发送了错误\n详细日志已发送到控制台")
			pterm.Error.Printf("修改物品名称: %v\n", err)
			return
		}
		// 读取新物品的数据
		dropResp, err := o.apis.DropItemAll(
			protocol.StackRequestSlotInfo{
				ContainerID:    28,
				Slot:           0,
				StackNetworkID: newItemDatas.StackNetworkID,
			},
			0,
		)
		if err != nil {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c尝试丢出新物品时失败\n详细日志已发送到控制台")
			pterm.Error.Printf("修改物品名称: %v\n", err)
			return
		}
		if !dropResp {
			o.Frame.GetGameControl().SayTo(chat.Name, "§c尝试丢出新物品时失败\n详细日志已发送到控制台")
			pterm.Error.Printf("修改物品名称: dropResp = %#v\n", dropResp)
			return
		}
		// 丢出新物品
		o.Frame.GetGameControl().SayTo(chat.Name, "§a已成功修改物品名称")
		return
		// 返回值
	}()
	return true
}
