package blockNBT_API

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
)

// 在描述物品移动操作时使用的结构体
type MoveItemDatas struct {
	WindowID                  int16 // 物品所在库存的窗口 ID
	ItemStackNetworkIDProvide int32 // 主动提供的 StackNetworkID
	ContainerID               uint8 // 物品所在库存的库存类型 ID
	Slot                      uint8 // 物品所在的槽位
}

// 将库存编号为 source 所指代的槽位中的物品移动到 destination 所指代的槽位。
// 当 MoveItemDatas 结构体的 WindowID 为 -1 时对应槽位中物品的 ItemStackNetworkID 值
// 将使用此结构体中提供的 ItemStackNetworkIDProvide 值 。
// 当且仅当物品操作得到租赁服的响应后，此函数才会返回物品操作结果。
// 考虑到物品操作请求在被批准后租赁服不会返回其他数据包用以描述对应槽位的最终结果，
// 因此在解决此问题前，此函数将暂时作为私有实现
func (g *GlobalAPI) moveItem(
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
			return []protocol.ItemStackResponse{}, fmt.Errorf("moveItem: %v", err)
		}
	} else {
		itemOnSource.StackNetworkID = source.ItemStackNetworkIDProvide
	}
	if destination.WindowID != -1 {
		itemOnDestination, err = g.PacketHandleResult.Inventory.GetItemStackInfo(uint32(destination.WindowID), destination.Slot)
		if err != nil {
			return []protocol.ItemStackResponse{}, fmt.Errorf("moveItem: %v", err)
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
	ans, err := g.sendItemStackRequestWithResponce(&packet.ItemStackRequest{
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
		return []protocol.ItemStackResponse{}, fmt.Errorf("moveItem: %v", err)
	}
	// 发送物品操作请求
	return ans, nil
	// 返回值
}

// 根据铁砧操作的返回值 resp 更新背包中对应物品栏的物品数据，属于私有实现。
// 此函数仅被铁砧的改名操作所使用，因为在进行改名操作后，租赁服似乎只会返回 ItemStackResponce 包
// 来告知客户端关于物品的最终操作结果，所以我们不得不手动构造改名后的 NBT 数据，然后利用
// 租赁服返回的 ItemStackResponce 包来更新客户端已保存的背包库存数据。
// oldItem 指被改名物品的原始信息
func (g *GlobalAPI) updateSlotInfoOnlyUseForAnvilChangeItemName(
	resp protocol.ItemStackResponse,
	oldItem protocol.ItemInstance,
) error {
	var correctDatas protocol.StackResponseSlotInfo = protocol.StackResponseSlotInfo{}
	for _, value := range resp.ContainerInfo {
		if value.ContainerID == 12 {
			correctDatas = value.SlotInfo[0]
			break
		}
	}
	// 从 resp 中提取有效数据
	nbt := oldItem.Stack.NBTData
	// 获取物品的旧 NBT 数据
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
	g.PacketHandleResult.Inventory.WriteItemStackInfo(0, correctDatas.Slot, newItem)
	// 更新槽位数据
	return nil
	// 返回值
}

// 根据铁砧操作的返回值 resp 更新背包中对应物品栏的物品数据，属于私有实现。
// 此函数仅被铁砧的改名操作所使用，因为在进行改名操作后，租赁服似乎只会返回 ItemStackResponce 包
// 来告知客户端关于物品的最终操作结果，所以我们不得不手动更新客户端已保存的背包库存数据。
// oldItem 指被改名物品的原始信息
func (g *GlobalAPI) updateSlotInfoOnlyUseForAnvilRevertOperation(
	resp protocol.ItemStackResponse,
	oldItem protocol.ItemInstance,
) error {
	var correctDatas protocol.StackResponseSlotInfo = protocol.StackResponseSlotInfo{}
	for _, value := range resp.ContainerInfo {
		if value.ContainerID == 12 {
			correctDatas = value.SlotInfo[0]
			break
		}
	}
	// 从 resp 中提取有效数据
	newItem := protocol.ItemInstance{
		StackNetworkID: correctDatas.StackNetworkID,
		Stack: protocol.ItemStack{
			ItemType:       oldItem.Stack.ItemType,
			BlockRuntimeID: oldItem.Stack.BlockRuntimeID,
			Count:          uint16(correctDatas.Count),
			NBTData:        oldItem.Stack.NBTData,
			CanBePlacedOn:  oldItem.Stack.CanBePlacedOn,
			CanBreak:       oldItem.Stack.CanBreak,
			HasNetworkID:   oldItem.Stack.HasNetworkID,
		},
	}
	g.PacketHandleResult.Inventory.WriteItemStackInfo(0, correctDatas.Slot, newItem)
	// 更新槽位数据
	return nil
	// 返回值
}

// 将已放入铁砧第一格(注意是第一格)的物品的物品名称修改为 name 并返还到背包中的 slot 处。
// oldItem 指被改名物品的原始信息，此参数将作为更新本地库存数据的依据，
// 因为租赁服返回的物品操作请求的批准数据包并不会描述一个完整的 protocol.ItemInstance 信息。
// 额外的，似乎对于任何物品操作请求，租赁服都只会用 packet.ItemStackResponce 来描述对应槽位
// 的最终变动结果，因此我们必须自行更新本地库存数据。
//
// 当且仅当租赁服回应操作结果后再返回值。
// resp 参数指代把物品放入铁砧第一格时租赁服返回的结果。
// 部分情况下此函数可能会遇见无法处理的错误，届时程序将抛出严重错误(panic)
func (g *GlobalAPI) ChangeItemName(
	resp protocol.ItemStackResponse,
	name string,
	slot uint8,
	oldItem protocol.ItemInstance,
) (bool, error) {
	var stackNetworkID int32
	var count uint8
	for _, value := range resp.ContainerInfo {
		if value.ContainerID == 0 {
			stackNetworkID = value.SlotInfo[0].StackNetworkID
			count = value.SlotInfo[0].Count
			break
		}
	}
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
		err = g.updateSlotInfoOnlyUseForAnvilChangeItemName(ans, oldItem)
		if err != nil {
			panic(fmt.Sprintf("ChangeItemName: %v", err))
		}
		// 更新槽位数据
	} else {
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
		newAns, err := g.moveItem(source, destination, count)
		if err != nil {
			panic(fmt.Sprintf("ChangeItemName: %v", err))
		}
		if newAns[0].Status != 0 {
			panic(fmt.Sprintf("ChangeItemName: Could not revert operation %v because of the new operation which numbered %v have been canceled by error code %v. This maybe is a BUG, please provide this logs to the developers!\nnewAns = %#v; source = %#v; destination = %#v; moveCount = %v", ans.RequestID, newAns[0].RequestID, newAns[0].Status, newAns, source, destination, count))
		}
		err = g.updateSlotInfoOnlyUseForAnvilRevertOperation(newAns[0], oldItem)
		if err != nil {
			panic(fmt.Sprintf("ChangeItemName: %v", err))
		}
		// 如果名称未发生变化或者经验值不足等造成的改名失败
	}
	// 改名后更新对应槽位的数据
	if ans.Status == 9 {
		return false, nil
	}
	// 如果名称未发生变化或者因为其他一些原因所导致的改名失败 (ans.Status = 9)
	if ans.Status != 0 {
		return false, fmt.Errorf("ChangeItemName: Operation %v have been canceled by error code %v; ans = %#v", ans.RequestID, ans.Status, ans)
	}
	// 如果物品操作请求被拒绝 (ans.Status = others)
	return true, nil
	// 返回值
}

// 将 source 所指代的槽位中的全部物品丢出。
// windowID 指代被丢出物品所在库存的窗口 ID
func (g *GlobalAPI) DropItemAll(
	source protocol.StackRequestSlotInfo,
	windowID uint32,
) (bool, error) {
	ans, err := g.sendItemStackRequestWithResponce(&packet.ItemStackRequest{
		Requests: []protocol.ItemStackRequest{
			{
				Actions: []protocol.StackRequestAction{
					&protocol.DropStackRequestAction{
						Count:    64,
						Source:   source,
						Randomly: false,
					},
				},
			},
		},
	})
	if err != nil {
		return false, fmt.Errorf("DropItemAll: %v", err)
	}
	if ans[0].Status != 0 {
		return false, nil
	}
	// 发送物品丢掷请求
	g.PacketHandleResult.Inventory.WriteItemStackInfo(windowID, source.Slot, protocol.ItemInstance{
		StackNetworkID: 0,
		Stack: protocol.ItemStack{
			ItemType: protocol.ItemType{
				NetworkID:     0,
				MetadataValue: 0,
			},
			BlockRuntimeID: 0,
			Count:          0,
			NBTData:        map[string]interface{}(nil),
			CanBePlacedOn:  []string(nil),
			CanBreak:       []string(nil),
			HasNetworkID:   false,
		},
	})
	// 刷新本地保存的物品数据
	return true, nil
	// 返回值
}
