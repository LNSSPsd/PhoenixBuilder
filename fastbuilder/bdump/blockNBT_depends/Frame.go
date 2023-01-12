package blockNBT_depends

import (
	"encoding/json"
	"fmt"
	"phoenixbuilder/fastbuilder/commands_generator"
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/io/commands"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"

	"github.com/go-gl/mathgl/mgl32"
)

type FrameData struct {
	ItemRotation   float32
	ItemDropChance float32
	Item           *FrameItemData
}

type FrameItemData struct {
	Name       string
	Data       uint16
	CanDestroy []interface{}
	CanPlaceOn []interface{}
	EnchList   []interface{}
}

type FrameInput struct {
	Environment  *environment.PBEnvironment
	Mainsettings *types.MainConfig
	BlockInfo    *types.Module
	Frame        *map[string]interface{}
}

func Frame(input *FrameInput) error {
	err := placeFrame(input.Environment, input.Mainsettings, input.BlockInfo, *input.Frame)
	if err != nil {
		return fmt.Errorf("Frame: Failed to place the entity block named %v at (%v,%v,%v), and the error log is %v", *input.BlockInfo.Block.Name, input.BlockInfo.Point.X, input.BlockInfo.Point.Y, input.BlockInfo.Point.Z, err)
	}
	return nil
}

func parseFrameData(Frame *map[string]interface{}, BlockName *string) (*FrameData, error) {
	var got types.ChestData = types.ChestData{}
	var err error = nil
	got, err = getContainerDataRun(*Frame, *BlockName)
	if err != nil {
		return &FrameData{}, fmt.Errorf("parseFrameData: %v", err)
	}
	// get item info
	var ok bool = false
	var normal bool = false
	var itemRotation float32 = 0.0
	var itemDropChance float32 = 0.0
	var enchList []interface{} = []interface{}{}
	FRAME := *Frame
	// prepare
	_, ok = FRAME["ItemRotation"]
	if ok {
		itemRotation, normal = FRAME["ItemRotation"].(float32)
		if !normal {
			return &FrameData{}, fmt.Errorf("parseFrameData: Could not parse Frame[\"ItemRotation\"]; Frame = %#v", FRAME)
		}
	}
	// ItemRotation
	_, ok = FRAME["ItemDropChance"]
	if ok {
		itemDropChance, normal = FRAME["ItemDropChance"].(float32)
		if !normal {
			return &FrameData{}, fmt.Errorf("parseFrameData: Could not parse Frame[\"ItemDropChance\"]; Frame = %#v", FRAME)
		}
	}
	// ItemDropChance
	FRAME_Item := FRAME["Item"].(map[string]interface{})
	_, ok = FRAME_Item["tag"]
	if ok {
		FRAME_tag, normal := FRAME_Item["tag"].(map[string]interface{})
		if !normal {
			return &FrameData{}, fmt.Errorf("parseFrameData: Could not parse Frame[\"Item\"][\"tag\"]; Frame = %#v", FRAME)
		}
		_, ok = FRAME_tag["ench"]
		if ok {
			enchList, normal = FRAME_tag["ench"].([]interface{})
			if !normal {
				return &FrameData{}, fmt.Errorf("parseFrameData: Could not parse Frame[\"Item\"][\"tag\"][\"ench\"]; Frame = %#v", FRAME)
			}
		}
	}
	// ench
	return &FrameData{
		ItemRotation:   itemRotation,
		ItemDropChance: itemDropChance,
		Item: &FrameItemData{
			Name:     got[0].Name,
			Data:     got[0].Damage,
			EnchList: enchList,
		},
	}, nil
}

func placeFrame(
	Environment *environment.PBEnvironment,
	Mainsettings *types.MainConfig,
	BlockInfo *types.Module,
	Frame map[string]interface{},
) error {
	FrameData, err := parseFrameData(&Frame, BlockInfo.Block.Name)
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	// parse sign data
	cmdsender := Environment.CommandSender.(*commands.CommandSender)
	var got interface{}
	var position map[string]float32
	resp, err := cmdsender.SendWSCommandWithResponce("querytarget @s")
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	json.Unmarshal([]byte(resp.OutputMessages[0].Parameters[0]), &got)
	position = got.([]interface{})[0].(map[string]interface{})["position"].(map[string]float32)
	// get bot pos
	_, err = cmdsender.SendWSCommandWithResponce(fmt.Sprintf("replaceitem entity @s slot.weapon.mainhand 0 %v 1 %v", FrameData.Item.Name, FrameData.Item.Data))
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	// replaceitem
	err = SendEnchantCommand(*Environment, &FrameData.Item.EnchList)
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	// enchant item stack
	request := commands_generator.SetBlockRequest(BlockInfo, Mainsettings)
	_, err = cmdsender.SendWSCommandWithResponce(request)
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	// place frame block
	if protocol.CurrentProtocol == 504 {
		networkID, ok := ItemRunTimeID[FrameData.Item.Name]
		if ok {
			Environment.Connection.(*minecraft.Conn).WritePacket(&packet.InventoryTransaction{
				LegacyRequestID:    0,
				LegacySetItemSlots: []protocol.LegacySetItemSlot(nil),
				Actions:            []protocol.InventoryAction{},
				TransactionData: &protocol.UseItemTransactionData{
					LegacyRequestID:    0,
					LegacySetItemSlots: []protocol.LegacySetItemSlot(nil),
					Actions:            []protocol.InventoryAction(nil),
					ActionType:         0x0,
					BlockPosition:      protocol.BlockPos{int32(BlockInfo.Point.X), int32(BlockInfo.Point.Y), int32(BlockInfo.Point.Z)},
					BlockFace:          1,
					HotBarSlot:         0,
					HeldItem: protocol.ItemInstance{
						StackNetworkID: 0,
						Stack: protocol.ItemStack{
							ItemType: protocol.ItemType{
								NetworkID:     int32(networkID),
								MetadataValue: uint32(FrameData.Item.Data),
							},
							BlockRuntimeID: 0,
							Count:          1,
							NBTData:        map[string]interface{}{},
							CanBePlacedOn:  []string{},
							CanBreak:       []string{},
							HasNetworkID:   false,
						},
					},
					Position:        mgl32.Vec3{position["x"], position["y"], position["z"]},
					ClickedPosition: mgl32.Vec3{0.0, 0.0, 0.0},
					//BlockRuntimeID:  180,
				},
			})
		}
	}
	// write nbt
	return nil
}
