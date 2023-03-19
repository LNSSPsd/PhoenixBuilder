package blockNBT_API

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"sync"
)

// 向租赁服发送 ItemStackReuqest 并获取返回值。
// 考虑到物品操作请求在被批准后租赁服不会返回其他数据包用以描述对应槽位的最终结果，
// 因此在解决此问题前，此函数将暂时作为私有实现
func (g *GlobalAPI) sendItemStackRequestWithResponce(request *packet.ItemStackRequest) ([]protocol.ItemStackResponse, error) {
	requestIDList := []int32{}
	ans := []protocol.ItemStackResponse{}
	// 初始化
	for range request.Requests {
		requestIDList = append(requestIDList, g.PacketHandleResult.ItemStackOperation.GetNewRequestID())
	}
	for key := range request.Requests {
		requestID := requestIDList[key]
		request.Requests[key].RequestID = requestID
		g.PacketHandleResult.ItemStackOperation.WriteRequest(requestID)
	}
	// 重新设定每个请求的请求 ID 并写入请求到等待队列
	err := g.WritePacket(request)
	if err != nil {
		return nil, fmt.Errorf("sendItemStackRequestWithResponce: %v", err)
	}
	// 发送物品操作请求
	for _, value := range requestIDList {
		g.PacketHandleResult.ItemStackOperation.AwaitResponce(value)
		got, err := g.PacketHandleResult.ItemStackOperation.LoadResponceAndDelete(value)
		if err != nil {
			return nil, fmt.Errorf("sendItemStackRequestWithResponce: %v", err)
		}
		ans = append(ans, got)
	}
	// 等待租赁服回应所有物品操作请求。同时，每当一个请求被响应，就把对应的结果保存下来
	return ans, nil
	// 返回值
}

// 向租赁服发送 ItemStackReuqest 并无视返回值，这代表着此函数在被执行后不会发生阻塞。
// 此函数的其中一部分将以协程运行，因此返回值的第二项代表此协程的执行是否存在错误。
// 除此外，在协程执行完成后，返回值第一项 *sync.Mutex 所指向的互斥锁将被释放，此时
// 返回值的第二项 *error 将会被赋值。
// 考虑到物品操作请求在被批准后租赁服不会返回其他数据包用以描述对应槽位的最终结果，
// 因此在解决此问题前，此函数将暂时作为私有实现
func (g *GlobalAPI) sendItemStackRequest(request *packet.ItemStackRequest) (*sync.Mutex, *error, error) {
	requestIDList := []int32{}
	goRotuineReturnInfo := struct {
		IsExecuting *sync.Mutex
		Return      *error
	}{
		IsExecuting: &sync.Mutex{},
		Return:      nil,
	}
	// 初始化
	for range request.Requests {
		requestIDList = append(requestIDList, g.PacketHandleResult.ItemStackOperation.GetNewRequestID())
	}
	for key := range request.Requests {
		requestID := requestIDList[key]
		request.Requests[key].RequestID = requestID
		g.PacketHandleResult.ItemStackOperation.WriteRequest(requestID)
	}
	// 重新设定每个请求的请求 ID 并写入请求到等待队列
	err := g.WritePacket(request)
	if err != nil {
		return nil, nil, fmt.Errorf("sendItemStackRequest: %v", err)
	}
	// 发送物品操作请求
	goRotuineReturnInfo.IsExecuting.Lock()
	go func() {
		defer func() {
			goRotuineReturnInfo.IsExecuting.TryLock()
			goRotuineReturnInfo.IsExecuting.Unlock()
		}()
		// while exit
		for _, value := range requestIDList {
			g.PacketHandleResult.ItemStackOperation.AwaitResponce(value)
			_, err := g.PacketHandleResult.ItemStackOperation.LoadResponceAndDelete(value)
			if err != nil {
				goRotuineReturnInfo.Return = &err
				return
			}
		}
		// await changes and return error if their is anything working wrong
	}()
	// 等待租赁服回应所有物品操作请求。这里以协程运行的目的在于无视返回值。
	return goRotuineReturnInfo.IsExecuting, goRotuineReturnInfo.Return, nil
	// 返回值
}
