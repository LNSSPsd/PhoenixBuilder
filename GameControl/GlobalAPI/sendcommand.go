package GlobalAPI

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"

	"github.com/google/uuid"
)

// 向租赁服发送 Sizukana 命令且无视返回值
func (g *GlobalAPI) SendSettingsCommand(command string, sendDimensionalCmd bool) error {
	if sendDimensionalCmd {
		command = fmt.Sprintf(`execute @a[name="%v"] ~ ~ ~ %v`, g.BotInfo.BotName, command)
	}
	err := g.WritePacket(&packet.SettingsCommand{
		CommandLine:    command,
		SuppressOutput: true,
	})
	if err != nil {
		return fmt.Errorf("SendSettingsCommand: %v", err)
	}
	return nil
}

// 向租赁服发送 WS 命令且无视返回值
func (g *GlobalAPI) SendWSCommand(command string, uniqueId uuid.UUID) error {
	requestId, _ := uuid.Parse("96045347-a6a3-4114-94c0-1bc4cc561694")
	err := g.WritePacket(&packet.CommandRequest{
		CommandLine: command,
		CommandOrigin: protocol.CommandOrigin{
			Origin:    protocol.CommandOriginAutomationPlayer,
			UUID:      uniqueId,
			RequestID: requestId.String(),
		},
		Internal:  false,
		UnLimited: false,
	})
	if err != nil {
		return fmt.Errorf("SendWSCommand: %v", err)
	}
	return nil
}

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
	ans, err := g.Resources.Command.LoadResponceAndDelete(uniqueId)
	if err != nil {
		return packet.CommandOutput{}, fmt.Errorf("SendWSCommandWithResponce: %v", err)
	}
	// 等待租赁服响应命令请求并取得命令请求的返回值
	return ans, nil
	// 返回值
}
