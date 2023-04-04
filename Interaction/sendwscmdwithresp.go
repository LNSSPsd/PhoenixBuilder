package GlobalAPI

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"

	"github.com/google/uuid"
)

// 向租赁服发送 WS 命令且获取返回值
func (g *GlobalAPI) SendWSCommandWithResponce(command string) (packet.CommandOutput, error) {
	uniqueId, err := uuid.NewUUID()
	if err != nil || uniqueId == uuid.Nil {
		return g.SendWSCommandWithResponce(command)
	}
	err = g.Resources.Command.WriteRequest(uniqueId)
	if err != nil {
		return packet.CommandOutput{}, fmt.Errorf("SendWSCommandWithResponce: %v", err)
	}
	// 写入请求到等待队列
	err = g.SendWSCommand(command, uniqueId)
	if err != nil {
		return packet.CommandOutput{}, fmt.Errorf("SendWSCommandWithResponce: %v", err)
	}
	// 发送命令
	g.Resources.Command.AwaitResponce(uniqueId)
	// 等待租赁服响应命令请求
	ans, err := g.Resources.Command.LoadResponceAndDelete(uniqueId)
	if err != nil {
		return packet.CommandOutput{}, fmt.Errorf("SendWSCommandWithResponce: %v", err)
	}
	// 取得命令请求的返回值
	return ans, nil
	// 返回值
}