package blockNBT_API

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"
	"sync"

	"github.com/google/uuid"
)

// 向租赁服发送 WS 命令且获取返回值
func (g *GlobalAPI) SendWSCommandWithResponce(command string) (packet.CommandOutput, error) {
	uniqueId, err := uuid.NewUUID()
	if err != nil || uniqueId == uuid.Nil {
		resp, err := g.SendWSCommandWithResponce(command)
		if err != nil {
			return packet.CommandOutput{}, fmt.Errorf("SendWSCommandWithResponce: %v", err)
		}
		return resp, nil
	}
	uniqueIdString := uniqueId.String()
	// 初始化
	g.PacketHandleResult.commandRequestMapLockDown.Lock()
	g.PacketHandleResult.commandRequest[uniqueIdString] = &CommandRequest{
		LockDown: sync.Mutex{},
		Responce: packet.CommandOutput{},
	}
	g.PacketHandleResult.commandRequest[uniqueIdString].LockDown.Lock()
	tmp := &g.PacketHandleResult.commandRequest[uniqueIdString].LockDown
	g.PacketHandleResult.commandRequestMapLockDown.Unlock()
	// 写入请求到等待队列
	err = g.SendWSCommand(command, uniqueId)
	if err != nil {
		return packet.CommandOutput{}, fmt.Errorf("SendWSCommandWithResponce: %v", err)
	}
	// 发送命令
	tmp.Lock()
	tmp.Unlock()
	// 等待租赁服返回结果
	ans := g.PacketHandleResult.commandRequest[uniqueIdString].Responce
	// 取得返回值
	g.PacketHandleResult.commandRequestMapLockDown.Lock()

	delete(g.PacketHandleResult.commandRequest, uniqueIdString)
	newMap := map[string]*CommandRequest{}
	for key, value := range g.PacketHandleResult.commandRequest {
		newMap[key] = value
	}
	g.PacketHandleResult.commandRequest = newMap

	g.PacketHandleResult.commandRequestMapLockDown.Unlock()
	// 从请求列表移除当前请求
	return ans, nil
	// 返回值
}
