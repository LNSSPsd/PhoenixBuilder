package blockNBT_API

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"sync"
	"sync/atomic"
)

// 向租赁服发送 ItemStackReuqest 并获取返回值
func (g *GlobalAPI) SendItemStackRequestWithResponce(request *packet.ItemStackRequest) ([]*ItemStackReuqestWithAns, error) {
	requestIDList := []int32{}
	waitingList := []*ItemStackReuqestWithAns{}
	// 初始化
	for range request.Requests {
		requestIDList = append(requestIDList, atomic.AddInt32(&g.PacketHandleResult.ItemStackRequestID, -2))
	}
	for key := range request.Requests {
		request.Requests[key].RequestID = requestIDList[key]
	}
	// 重新设定每个请求的请求 ID
	g.PacketHandleResult.ItemStackReuqestWithResultMapLockDown.Lock()

	for _, value := range requestIDList {
		g.PacketHandleResult.ItemStackReuqestWithResult[value] = &ItemStackReuqestWithAns{
			LockDown:      sync.Mutex{},
			SuccessStates: false,
			ErrorCode:     0,
		}
		g.PacketHandleResult.ItemStackReuqestWithResult[value].LockDown.Lock()
		waitingList = append(waitingList, g.PacketHandleResult.ItemStackReuqestWithResult[value])
	}

	g.PacketHandleResult.ItemStackReuqestWithResultMapLockDown.Unlock()
	// 写入请求到等待队列
	err := g.WritePacket(request)
	if err != nil {
		return nil, fmt.Errorf("SendItemStackRequestWithResponce: %v", err)
	}
	// 发送物品操作请求
	for _, value := range waitingList {
		value.LockDown.Lock()
		value.LockDown.Unlock()
	}
	// 等待租赁服回应所有物品操作请求
	g.PacketHandleResult.ItemStackReuqestWithResultMapLockDown.Lock()

	for _, value := range requestIDList {
		delete(g.PacketHandleResult.ItemStackReuqestWithResult, value)
	}
	newMap := map[int32]*ItemStackReuqestWithAns{}
	for key, value := range g.PacketHandleResult.ItemStackReuqestWithResult {
		newMap[key] = value
	}
	g.PacketHandleResult.ItemStackReuqestWithResult = newMap

	g.PacketHandleResult.ItemStackReuqestWithResultMapLockDown.Unlock()
	// 将已完成的请求释放(移除请求)
	return waitingList, nil
	// 返回值
}

/*
将背包中槽位为 inventorySlot 的物品移动到已打开容器的第 containerSlot 槽位，且只移动 moveCount 个物品

此函数将 containerSlot 处的物品当作空气处理。如果涉及到交换物品等操作，或许您需要使用其他函数
*/
func (g *GlobalAPI) PlaceItemIntoContainer(
	inventorySlot uint8,
	containerSlot uint8,
	moveCount uint8,
) error {
	datas, _ := g.GetInventoryCotent(0)
	// 取得 Window ID 为 0 的库存数据，也就是取得背包的库存数据
	got, ok := datas[inventorySlot]
	if !ok {
		return fmt.Errorf("PlaceItemIntoContainer: %v is not in inventory contents; datas = %#v", inventorySlot, datas)
	}
	// 取得“来源”的物品信息
	placeStackRequestAction := protocol.PlaceStackRequestAction{}
	if moveCount <= uint8(got.Stack.Count) {
		placeStackRequestAction.Count = moveCount
	} else {
		placeStackRequestAction.Count = uint8(got.Stack.Count)
	}
	// 得到欲移动的物品数量
	placeStackRequestAction.Source = protocol.StackRequestSlotInfo{
		ContainerID:    12,
		Slot:           inventorySlot,
		StackNetworkID: got.StackNetworkID,
	}
	placeStackRequestAction.Destination = protocol.StackRequestSlotInfo{
		ContainerID:    g.PacketHandleResult.ContainerOpenDatas.Datas.ContainerType,
		Slot:           containerSlot,
		StackNetworkID: 0,
	}
	// 前置准备
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
		return fmt.Errorf("PlaceItemIntoContainer: %v", err)
	}
	// 发送物品操作请求
	if ans[0].SuccessStates == false {
		return fmt.Errorf("PlaceItemIntoContainer: Operation is canceled, and the errorCode(status) is %v; inventorySlot = %v, containerSlot = %v, moveCount = %v", ans[0].ErrorCode, inventorySlot, containerSlot, moveCount)
	}
	// 当操作失败时
	return nil
	// 返回值
}
