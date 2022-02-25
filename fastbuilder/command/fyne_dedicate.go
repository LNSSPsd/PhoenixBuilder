// +build fyne_gui

package command

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	bridge_fmt "phoenixbuilder/dedicate/fyne/bridge"
	"strings"
	"sync"
	"time"
)


var UUIDMap sync.Map //= make(map[string]func(*minecraft.Conn,*[]protocol.CommandOutputMessage))
var BlockUpdateSubscribeMap sync.Map
var AdditionalChatCb func(string)
var AdditionalTitleCb func(s string)


func init() {
	AdditionalChatCb = func(s string) {}
	AdditionalTitleCb = func(s string) {}
}


func ReplaceItemRequest(module *types.Module, config *types.MainConfig) string {
	//C.replaceItemRequestInternal(unsafe.Pointer(buf), C.int(module.Point.X), C.int(module.Point.Y), C.int(module.Point.Z), C.uchar(module.ChestSlot.Slot),C.CString(module.ChestSlot.Name),C.uchar(module.ChestSlot.Count), C.ushort(module.ChestSlot.Damage))
	return fmt.Sprintf("replaceitem block %d %d %d slot.container %d %s %d %d", module.Point.X, module.Point.Y, module.Point.Z, module.ChestSlot.Slot, module.ChestSlot.Name, module.ChestSlot.Count, module.ChestSlot.Damage)
}


func ClearUUIDMap() {
	UUIDMap = sync.Map{}
}

func SendCommand(command string, UUID uuid.UUID, conn *minecraft.Conn) error {
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
	return conn.WritePacket(commandRequest)
}

func SendWSCommand(command string, UUID uuid.UUID, conn *minecraft.Conn) error {
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
	return conn.WritePacket(commandRequest)
}

func SendSizukanaCommand(command string, conn *minecraft.Conn) error {
	return conn.WritePacket(&packet.SettingsCommand{
		CommandLine:    command,
		SuppressOutput: true,
	})
}

func SendChat(content string, conn *minecraft.Conn) error {
	AdditionalChatCb(content)
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

func SetBlockRequest(module *types.Module, config *types.MainConfig) string {
	Block := module.Block
	Point := module.Point
	Method := config.Method
	if Block != nil {
		return fmt.Sprintf("setblock %v %v %v %v %v %v", Point.X, Point.Y, Point.Z, *Block.Name, Block.Data, Method)
	} else {
		return fmt.Sprintf("setblock %v %v %v %v %v %v", Point.X, Point.Y, Point.Z, config.Block.Name, config.Block.Data, Method)
	}
}


type TellrawItem struct {
	Text string `json:"text"`
}

type TellrawStruct struct {
	RawText []TellrawItem `json:"rawtext"`
}

func TellRawRequest(target types.Target, lines ...string) string {
	now := time.Now().Format("§6{15:04:05}§b")
	var items []TellrawItem
	for _, text := range lines {
		msg := fmt.Sprintf("%v %v", now, strings.Replace(text, "schematic", "sc***atic", -1))
		items = append(items, TellrawItem{Text: msg})
	}
	final := &TellrawStruct{
		RawText: items,
	}
	content, _ := json.Marshal(final)
	bridge_fmt.Printf("%s\n", content)
	cmd := fmt.Sprintf("tellraw %v %s", target, content)
	return cmd
}

func Tellraw(conn *minecraft.Conn, lines ...string) error {
	//uuid1, _ := uuid.NewUUID()
	bridge_fmt.Printf("%s\n", lines[0])
	//return nil
	msg := strings.Replace(lines[0], "schematic", "sc***atic", -1)
	msg = strings.Replace(msg, ".", "．", -1)
	// Netease set .bdx, .schematic, .mcacblock, etc as blocked words
	// So we should replace half-width points w/ full-width points to avoid being
	// blocked
	//return SendChat(fmt.Sprintf("§b%s",msg), conn)
	return SendSizukanaCommand(TellRawRequest(types.AllPlayers, lines...), conn)
}

func RawTellRawRequest(target types.Target, line string) string {
	var items []TellrawItem
	msg := strings.Replace(line, "schematic", "sc***atic", -1)
	items = append(items, TellrawItem{Text: msg})
	final := &TellrawStruct{
		RawText: items,
	}
	content, _ := json.Marshal(final)
	cmd := fmt.Sprintf("tellraw %v %s", target, content)
	return cmd
}

func WorldChatTellraw(conn *minecraft.Conn, sender string, content string) error {
	bridge_fmt.Printf("W <%s> %s\n", sender, content)
	str := fmt.Sprintf("§eW §r<%s> %s", sender, content)
	return SendSizukanaCommand(RawTellRawRequest(types.AllPlayers, str), conn)
}

func TitleRequest(target types.Target, lines ...string) string {
	var items []TellrawItem
	for _, text := range lines {
		items = append(items, TellrawItem{Text: strings.Replace(text, "schematic", "sc***atic", -1)})
	}
	final := &TellrawStruct{
		RawText: items,
	}
	content, _ := json.Marshal(final)
	AdditionalTitleCb(string(content))
	cmd := fmt.Sprintf("titleraw %v actionbar %s", target, content)
	return cmd
}

func Title(conn *minecraft.Conn, lines ...string) error {
	l := TitleRequest(types.AllPlayers, lines...)
	return SendSizukanaCommand(l, conn)
}
