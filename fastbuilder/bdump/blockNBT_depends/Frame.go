package blockNBT_depends

import (
	"fmt"
	"phoenixbuilder/fastbuilder/commands_generator"
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/io/commands"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
)

type FrameData struct {
	ItemRotation float32
	Item         types.ChestData
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
	var normal bool = false
	var itemRotation float32 = 0.0
	FRAME := *Frame
	// prepare
	_, ok := FRAME["ItemRotation"]
	if ok {
		itemRotation, normal = FRAME["ItemRotation"].(float32)
		if !normal {
			return &FrameData{}, fmt.Errorf("parseFrameData: Could not parse Frame[\"ItemRotation\"]; Frame = %#v", FRAME)
		}
	}
	// ItemRotation
	return &FrameData{
		ItemRotation: itemRotation,
		Item:         got,
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
	// parse frame data
	cmdsender := Environment.CommandSender.(*commands.CommandSender)
	err = cmdsender.SendDimensionalCommand(fmt.Sprintf("setblock %v %v %v minecraft:air", BlockInfo.Point.X, BlockInfo.Point.Y, BlockInfo.Point.Z))
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	request := commands_generator.SetBlockRequest(BlockInfo, Mainsettings)
	_, err = cmdsender.SendWSCommandWithResponce(request)
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	// place frame block
	if len(FrameData.Item) <= 0 {
		return nil
	}
	// if their is nothing in the frame block, then return
	_, err = cmdsender.SendWSCommandWithResponce("replaceitem entity @s slot.hotbar 0 air")
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	Environment.Connection.(*minecraft.Conn).WritePacket(&packet.MobEquipment{
		EntityRuntimeID: Environment.Connection.(*minecraft.Conn).GameData().EntityRuntimeID,
		NewItem: protocol.ItemInstance{
			StackNetworkID: 0,
			Stack: protocol.ItemStack{
				ItemType: protocol.ItemType{
					NetworkID:     0,
					MetadataValue: 0,
				},
				BlockRuntimeID: 0,
				Count:          0,
				NBTData:        map[string]interface{}{},
				CanBePlacedOn:  []string(nil),
				CanBreak:       []string(nil),
				HasNetworkID:   false,
			},
		},
		InventorySlot: 0,
		HotBarSlot:    0,
		WindowID:      0,
	})
	// change the slot.weapon.mainhand into slot.hotbar 0
	_, err = cmdsender.SendWSCommandWithResponce(fmt.Sprintf("replaceitem entity @s slot.weapon.mainhand 0 %v 1 %v %v", FrameData.Item[0].Name, FrameData.Item[0].Damage, commands_generator.GetReplaceItemEnhancement(
		&types.Module{
			ChestData: &types.ChestData{
				types.ChestSlot{
					ItemLock:    FrameData.Item[0].ItemLock,
					KeepOnDeath: FrameData.Item[0].KeepOnDeath,
					CanPlaceOn:  FrameData.Item[0].CanPlaceOn,
					CanDestroy:  FrameData.Item[0].CanDestroy,
				},
			},
		}, 0)))
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	// replaceitem
	err = EnchantRequest(Environment, FrameData.Item[0].EnchList)
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	// enchant item stack
	blockStates, err := ParseStringBlockStates(&BlockInfo.Block.BlockStates)
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	BLOCKSTATES := *blockStates
	_, ok := BLOCKSTATES["facing_direction"]
	if !ok {
		return fmt.Errorf("placeFrame: Could not find BlockInfo.Block.BlockStates[\"facing_direction\"]; BlockInfo.Block.BlockStates = %#v", BlockInfo.Block.BlockStates)
	}
	facing_direction, normal := BLOCKSTATES["facing_direction"].(int32)
	if !normal {
		return fmt.Errorf("placeFrame: Could not parse BlockInfo.Block.BlockStates[\"facing_direction\"]; BlockInfo.Block.BlockStates = %#v", BlockInfo.Block.BlockStates)
	}
	if facing_direction > 5 {
		return fmt.Errorf("placeFrame: Unexpected BlockInfo.Block.BlockStates[\"facing_direction\"]; BlockInfo.Block.BlockStates = %#v", BlockInfo.Block.BlockStates)
	}
	// get facing_direction
	if protocol.CurrentProtocol == 504 {
		networkID, ok := ItemRunTimeID[FrameData.Item[0].Name]
		if ok {
			var clickNum int = 1
			if FrameData.ItemRotation == 45.0 {
				clickNum = 2
			}
			if FrameData.ItemRotation == 90.0 {
				clickNum = 3
			}
			if FrameData.ItemRotation == 135.0 {
				clickNum = 4
			}
			if FrameData.ItemRotation == 180.0 {
				clickNum = 5
			}
			if FrameData.ItemRotation == 225.0 {
				clickNum = 6
			}
			if FrameData.ItemRotation == 270.0 {
				clickNum = 7
			}
			if FrameData.ItemRotation == 315.0 {
				clickNum = 8
			}
			// the rotation Angle of the item
			var blockRuntimeID uint32 = 0
			if *BlockInfo.Block.Name == "glow_frame" {
				blockRuntimeID = 179
			}
			if *BlockInfo.Block.Name == "frame" {
				blockRuntimeID = 4211
			}
			blockRuntimeID = blockRuntimeID + uint32(facing_direction)
			// get run time id of frame block
			if FrameData.Item[0].EnchList != nil {
				Environment.Connection.(*minecraft.Conn).WritePacket(&packet.MobEquipment{
					EntityRuntimeID: Environment.Connection.(*minecraft.Conn).GameData().EntityRuntimeID,
					NewItem: protocol.ItemInstance{
						StackNetworkID: 0,
						Stack: protocol.ItemStack{
							ItemType: protocol.ItemType{
								NetworkID:     int32(networkID),
								MetadataValue: uint32(FrameData.Item[0].Damage),
							},
							BlockRuntimeID: 0,
							Count:          1,
							NBTData: map[string]interface{}{
								"ench": FrameData.Item[0].EnchList,
							},
							CanBePlacedOn: FrameData.Item[0].CanPlaceOn,
							CanBreak:      FrameData.Item[0].CanDestroy,
							HasNetworkID:  false,
						},
					},
					InventorySlot: 0,
					HotBarSlot:    0,
					WindowID:      0,
				})
			} else {
				Environment.Connection.(*minecraft.Conn).WritePacket(&packet.MobEquipment{
					EntityRuntimeID: Environment.Connection.(*minecraft.Conn).GameData().EntityRuntimeID,
					NewItem: protocol.ItemInstance{
						StackNetworkID: 0,
						Stack: protocol.ItemStack{
							ItemType: protocol.ItemType{
								NetworkID:     int32(networkID),
								MetadataValue: uint32(FrameData.Item[0].Damage),
							},
							BlockRuntimeID: 0,
							Count:          1,
							CanBePlacedOn:  FrameData.Item[0].CanPlaceOn,
							CanBreak:       FrameData.Item[0].CanDestroy,
							HasNetworkID:   false,
						},
					},
					InventorySlot: 0,
					HotBarSlot:    0,
					WindowID:      0,
				})
			}
			// let other players know what happened
			// it is not necessary to do
			for i := 0; i < clickNum; i++ {
				Environment.Connection.(*minecraft.Conn).WritePacket(&packet.InventoryTransaction{
					LegacyRequestID:    0,
					LegacySetItemSlots: []protocol.LegacySetItemSlot(nil),
					Actions:            []protocol.InventoryAction{},
					TransactionData: &protocol.UseItemTransactionData{
						LegacyRequestID:    0,
						LegacySetItemSlots: []protocol.LegacySetItemSlot(nil),
						Actions:            []protocol.InventoryAction(nil),
						ActionType:         0,
						BlockPosition:      protocol.BlockPos{int32(BlockInfo.Point.X), int32(BlockInfo.Point.Y), int32(BlockInfo.Point.Z)},
						//BlockFace:          blockFace,
						HotBarSlot: 0,
						HeldItem: protocol.ItemInstance{
							StackNetworkID: 0,
							Stack: protocol.ItemStack{
								ItemType: protocol.ItemType{
									NetworkID:     int32(networkID),
									MetadataValue: uint32(FrameData.Item[0].Damage),
								},
								BlockRuntimeID: 0,
								Count:          1,
								//NBTData:        map[string]interface{}{},
								CanBePlacedOn: FrameData.Item[0].CanPlaceOn,
								CanBreak:      FrameData.Item[0].CanDestroy,
								HasNetworkID:  false,
							},
						},
						//Position:        mgl32.Vec3{float32(posx), float32(posy), float32(posz)},
						//ClickedPosition: mgl32.Vec3{0, 0, 0},
						BlockRuntimeID: blockRuntimeID,
					},
				})
			}
		}
	}
	// put item into the frame
	return nil
}
