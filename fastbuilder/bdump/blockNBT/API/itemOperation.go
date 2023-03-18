package blockNBT_API

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"

	"github.com/pterm/pterm"
)

// 在描述物品移动操作时使用的结构体
type MoveItemDatas struct {
	WindowID                  int16 // 物品所在库存的窗口 ID
	ItemStackNetworkIDProvide int32 // 主动提供的 StackNetworkID
	ContainerID               uint8 // 物品所在库存的库存类型 ID
	Slot                      uint8 // 物品所在的槽位
}

/*
将库存编号为 source 所指代的物品移动到 destination 所指代的物品。
当 MoveItemDatas 结构体的 WindowID 为 -1 时将会认为对应物品所在的库存不存在窗口 ID ，
此时将转而使用此结构体中的 ItemStackNetworkIDProvide 值 。
当且仅当物品操作得到租赁服的响应后，此函数才会返回物品操作结果。
*/
func (g *GlobalAPI) MoveItem(
	source MoveItemDatas,
	destination MoveItemDatas,
	moveCount uint8,
) ([]protocol.ItemStackResponse, error) {
	placeStackRequestAction := protocol.PlaceStackRequestAction{}
	var itemOnSource protocol.ItemInstance = protocol.ItemInstance{}
	var itemOnDestination protocol.ItemInstance = protocol.ItemInstance{}
	var err error = nil
	// 初始化
	if source.WindowID != -1 {
		itemOnSource, err = g.PacketHandleResult.Inventory.GetItemStackInfo(uint32(source.WindowID), source.Slot)
		if err != nil {
			return []protocol.ItemStackResponse{}, fmt.Errorf("MoveItem: %v", err)
		}
	} else {
		itemOnSource.StackNetworkID = source.ItemStackNetworkIDProvide
	}
	if destination.WindowID != -1 {
		itemOnDestination, err = g.PacketHandleResult.Inventory.GetItemStackInfo(uint32(destination.WindowID), destination.Slot)
		if err != nil {
			return []protocol.ItemStackResponse{}, fmt.Errorf("MoveItem: %v", err)
		}
	} else {
		itemOnDestination.StackNetworkID = destination.ItemStackNetworkIDProvide
	}
	// 取得 source 和 destination 处的物品信息
	if moveCount <= uint8(itemOnSource.Stack.Count) || source.WindowID == -1 {
		placeStackRequestAction.Count = moveCount
	} else {
		placeStackRequestAction.Count = uint8(itemOnSource.Stack.Count)
	}
	// 得到欲移动的物品数量
	placeStackRequestAction.Source = protocol.StackRequestSlotInfo{
		ContainerID:    source.ContainerID,
		Slot:           source.Slot,
		StackNetworkID: itemOnSource.StackNetworkID,
	}
	placeStackRequestAction.Destination = protocol.StackRequestSlotInfo{
		ContainerID:    destination.ContainerID,
		Slot:           destination.Slot,
		StackNetworkID: itemOnDestination.StackNetworkID,
	}
	// 构造 placeStackRequestAction 结构体
	ans, err := g.SendItemStackRequestWithResponce(&packet.ItemStackRequest{
		Requests: []protocol.ItemStackRequest{
			{
				Actions: []protocol.StackRequestAction{
					&placeStackRequestAction,
				},
				FilterStrings: []string{},
			},
		},
	})
	if err != nil {
		return []protocol.ItemStackResponse{}, fmt.Errorf("MoveItem: %v", err)
	}
	// 发送物品操作请求
	return ans, nil
	// 返回值
}

// 根据铁砧操作的返回值 resp 更新背包中对应物品栏的物品数据，属于私有实现。
// 此函数仅被铁砧的改名操作所使用，因为在进行改名操作后，租赁服似乎只会返回 ItemStackResponce 包
// 来告知客户端关于物品的最终操作结果，所以我们不得不手动构造改名后的 NBT 数据，然后利用
// 租赁服返回的 ItemStackResponce 包来更新客户端已保存的背包库存数据
func (g *GlobalAPI) updateSlotInfoOnlyUseForAnvilChangeItemName(resp protocol.ItemStackResponse) error {
	var correctDatas protocol.StackResponseSlotInfo = protocol.StackResponseSlotInfo{}
	for _, value := range resp.ContainerInfo {
		if value.ContainerID == 12 {
			correctDatas = value.SlotInfo[0]
			break
		}
	}
	// 从 resp 中提取有效数据
	oldItem, err := g.PacketHandleResult.Inventory.GetItemStackInfo(0, correctDatas.Slot)
	if err != nil {
		return fmt.Errorf("updateSlotInfoOnlyUseForAnvilChangeItemName: %v", err)
	}
	nbt := oldItem.Stack.NBTData
	// 获取物品的旧数据
	_, ok := nbt["tag"]
	if !ok {
		nbt["tag"] = map[string]interface{}{}
	}
	tag, normal := nbt["tag"].(map[string]interface{})
	if !normal {
		return fmt.Errorf("updateSlotInfoOnlyUseForAnvilChangeItemName: Failed to convert nbt[\"tag\"] into map[string]interface{}; nbt = %#v", nbt)
	}
	// tag
	_, ok = tag["display"]
	if !ok {
		tag["display"] = map[string]interface{}{}
		nbt["tag"].(map[string]interface{})["display"] = map[string]interface{}{}
	}
	display, normal := tag["display"].(map[string]interface{})
	if !normal {
		return fmt.Errorf("updateSlotInfoOnlyUseForAnvilChangeItemName: Failed to convert tag[\"display\"] into map[string]interface{}; tag = %#v", tag)
	}
	// display
	_, ok = display["Name"]
	if !ok {
		display["Name"] = correctDatas.CustomName
	}
	// name
	nbt["tag"].(map[string]interface{})["display"].(map[string]interface{})["Name"] = correctDatas.CustomName
	// 更新物品名称
	oldItem.Stack.NBTData = nbt
	newItem := protocol.ItemInstance{
		StackNetworkID: correctDatas.StackNetworkID,
		Stack: protocol.ItemStack{
			ItemType:       oldItem.Stack.ItemType,
			BlockRuntimeID: oldItem.Stack.BlockRuntimeID,
			Count:          uint16(correctDatas.Count),
			NBTData:        nbt,
			CanBePlacedOn:  oldItem.Stack.CanBePlacedOn,
			CanBreak:       oldItem.Stack.CanBreak,
			HasNetworkID:   oldItem.Stack.HasNetworkID,
		},
	}
	g.PacketHandleResult.Inventory.writeItemStackInfo(0, correctDatas.Slot, newItem)
	// 更新槽位数据
	return nil
	// 返回值
}

// 将已放入铁砧第一格(注意是第一格)的物品的物品名称修改为 name 并返还到背包中的 slot 处。
// 当且仅当租赁服回应操作结果后再返回值。
// resp 参数指代把物品放入铁砧第一格时租赁服返回的结果。
func (g *GlobalAPI) ChangeItemName(resp protocol.ItemStackResponse, name string, slot uint8) (bool, error) {
	var stackNetworkID int32
	var count uint8
	for _, value := range resp.ContainerInfo {
		if value.ContainerID == 0 {
			stackNetworkID = value.SlotInfo[0].StackNetworkID
			count = value.SlotInfo[0].Count
			break
		}
	}
	// 初始化
	newRequestID := g.PacketHandleResult.ItemStackOperation.GetNewRequestID()
	// 请求一个新的 RequestID 用于 ItemStackRequest
	placeStackRequestAction := protocol.PlaceStackRequestAction{}
	placeStackRequestAction.Count = count
	placeStackRequestAction.Source = protocol.StackRequestSlotInfo{
		ContainerID:    0x3c,
		Slot:           0x32,
		StackNetworkID: newRequestID,
	}
	placeStackRequestAction.Destination = protocol.StackRequestSlotInfo{
		ContainerID:    0xc,
		Slot:           slot,
		StackNetworkID: 0,
	}
	// 构造一个新的 PlaceStackRequestAction 结构体
	newItemStackRequest := packet.ItemStackRequest{
		Requests: []protocol.ItemStackRequest{
			{
				RequestID: newRequestID,
				Actions: []protocol.StackRequestAction{
					&protocol.CraftRecipeOptionalStackRequestAction{
						RecipeNetworkID:   0,
						FilterStringIndex: 0,
					},
					&protocol.ConsumeStackRequestAction{
						DestroyStackRequestAction: protocol.DestroyStackRequestAction{
							Count: count,
							Source: protocol.StackRequestSlotInfo{
								ContainerID:    0,
								Slot:           1,
								StackNetworkID: stackNetworkID,
							},
						},
					},
					&placeStackRequestAction,
				},
				FilterStrings: []string{name},
			},
		},
	}
	// 构造一个新的 ItemStackRequest 结构体
	g.PacketHandleResult.ItemStackOperation.WriteRequest(newRequestID)
	// 写入请求到等待队列
	err := g.WritePacket(&newItemStackRequest)
	if err != nil {
		return false, fmt.Errorf("ChangeItemName: %v", err)
	}
	// 发送物品操作请求
	g.PacketHandleResult.ItemStackOperation.AwaitResponce(newRequestID)
	ans, err := g.PacketHandleResult.ItemStackOperation.LoadResponceAndDelete(newRequestID)
	if err != nil {
		return false, fmt.Errorf("ChangeItemName: %v", err)
	}
	// 取得物品操作请求的结果
	if ans.Status == 0 {
		err = g.updateSlotInfoOnlyUseForAnvilChangeItemName(ans)
		if err != nil {
			return false, fmt.Errorf("ChangeItemName: %v", err)
		}
	}
	// 更新槽位数据
	if ans.Status == 9 {
		source := MoveItemDatas{
			WindowID:                  -1,
			ItemStackNetworkIDProvide: stackNetworkID,
			ContainerID:               0,
			Slot:                      1,
		}
		destination := MoveItemDatas{
			WindowID:                  -1,
			ItemStackNetworkIDProvide: 0,
			ContainerID:               12,
			Slot:                      slot,
		}
		newAns, err := g.MoveItem(source, destination, count)
		if err != nil {
			return false, fmt.Errorf("ChangeItemName: %v", err)
		}
		if newAns[0].Status != 0 {
			panic(pterm.Error.Sprintf("ChangeItemName: Could not revert operation %v because of the new operation which numbered %v have been canceled by error code %v. This maybe is a BUG, please provide this logs to the developers!\nnewAns = %#v; source = %#v; destination = %#v; moveCount = %v", ans.RequestID, newAns, source, destination, count))
		}
		return false, nil
	}
	// 如果名称未发生变化或者因为其他一些原因所导致的改名失败 (error code = 9)
	if ans.Status != 0 {
		return false, fmt.Errorf("ChangeItemName: Operation %v have been canceled by error code %v; ans = %#v", ans.RequestID, ans.Status, ans)
	}
	// 如果物品操作请求被拒绝 (error code = others)
	return true, nil
	// 返回值
}
