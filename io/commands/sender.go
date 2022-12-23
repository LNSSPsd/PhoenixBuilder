//go:build !is_tweak
// +build !is_tweak

package commands

import (
	"fmt"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"sync"

	"github.com/google/uuid"
)

func (sender *CommandSender) GetBlockUpdateSubscribeMap() *sync.Map {
	return &sender.BlockUpdateSubscribeMap
}

func (sender *CommandSender) GetUUIDMap() *sync.Map {
	return &sender.UUIDMap
}

func (sender *CommandSender) ClearUUIDMap() {
	sender.UUIDMap = sync.Map{}
}

func (sender *CommandSender) getConn() *minecraft.Conn {
	conn := sender.env.Connection.(*minecraft.Conn)
	return conn
}

func (sender *CommandSender) SendCommand(command string, UUID uuid.UUID) error {
	requestId, _ := uuid.Parse("96045347-a6a3-4114-94c0-1bc4cc561694")
	origin := protocol.CommandOrigin{
		Origin:         protocol.CommandOriginPlayer,
		UUID:           UUID,
		RequestID:      requestId.String(),
		PlayerUniqueID: 0,
	}
	commandRequest := &packet.CommandRequest{
		CommandLine:   command,
		CommandOrigin: origin,
		Internal:      false,
		UnLimited:     false,
	}
	return sender.getConn().WritePacket(commandRequest)
}

func (sender *CommandSender) SendWSCommand(command string, UUID uuid.UUID) error {
	requestId, _ := uuid.Parse("96045347-a6a3-4114-94c0-1bc4cc561694")
	origin := protocol.CommandOrigin{
		Origin:         protocol.CommandOriginAutomationPlayer,
		UUID:           UUID,
		RequestID:      requestId.String(),
		PlayerUniqueID: 0,
	}
	commandRequest := &packet.CommandRequest{
		CommandLine:   command,
		CommandOrigin: origin,
		Internal:      false,
		UnLimited:     false,
	}
	return sender.getConn().WritePacket(commandRequest)
}

func (sender *CommandSender) SendSizukanaCommand(command string) error {
	return sender.getConn().WritePacket(&packet.SettingsCommand{
		CommandLine:    command,
		SuppressOutput: true,
	})
}

func (sender *CommandSender) SendChat(content string) error {
	conn := sender.getConn()
	idd := conn.IdentityData()
	return conn.WritePacket(&packet.Text{
		TextType:         packet.TextTypeChat,
		NeedsTranslation: false,
		SourceName:       idd.DisplayName,
		Message:          content,
		XUID:             idd.XUID,
		PlayerRuntimeID:  fmt.Sprintf("%d", conn.GameData().EntityUniqueID),
	})
}

func (sender *CommandSender) GetBotName() string {
	ud, _ := uuid.NewUUID()
	chann := make(chan *packet.CommandOutput)
	sender.UUIDMap.Store(ud.String(), chann)
	sender.SendCommand("testfor @s", ud)
	resp := <-chann
	return resp.OutputMessages[0].Parameters[0]
}
