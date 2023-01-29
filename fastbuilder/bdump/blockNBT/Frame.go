package blockNBT

import (
	"fmt"
	"math"
	blockNBT_depends "phoenixbuilder/fastbuilder/bdump/blockNBT/depends"
	"phoenixbuilder/fastbuilder/commands_generator"
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/fastbuilder/mcstructure"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/io/commands"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
)

const glowFrameBlockRunTimeId = 179
const frameBlockRunTimeId = 4211

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
	if len(FrameData.Item) <= 0 || !blockNBT_depends.CheckVersion() {
		return nil
	}
	// if it is nothing in the frame block or the current version are not support, then return nil
	var writtenBookMark bool
	if FrameData.Item[0].Name == "written_book" {
		FrameData.Item[0].Name = "writable_book"
		writtenBookMark = true
	}
	err = blockNBT_depends.ReplaceitemAndEnchant(Environment, &types.ChestSlot{
		Name:       FrameData.Item[0].Name,
		Count:      FrameData.Item[0].Count,
		Damage:     FrameData.Item[0].Damage,
		CanPlaceOn: FrameData.Item[0].CanPlaceOn,
		CanDestroy: FrameData.Item[0].CanDestroy,
		ItemNBT:    FrameData.Item[0].ItemNBT,
	})
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	if FrameData.Item[0].Name == "writable_book" {
		err = blockNBT_depends.WriteTextToBook(Environment, &FrameData.Item[0])
	}
	if writtenBookMark {
		FrameData.Item[0].Name = "written_book"
	}
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	// replaceitem and enchant item stack
	blockStates, err := mcstructure.ParseStringNBT(BlockInfo.Block.BlockStates, true)
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	BLOCKSTATES, normal := blockStates.(map[string]interface{})
	if !normal {
		return fmt.Errorf("placeFrame: Failed to converse BlockStates to map[string]interface{}, BlockStates = %#v", BlockInfo.Block.BlockStates)
	}
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
	_, err = cmdsender.SendWSCommandWithResponce(fmt.Sprintf("tp %v %v %v", BlockInfo.Point.X, BlockInfo.Point.Y, BlockInfo.Point.Z))
	if err != nil {
		return fmt.Errorf("placeFrame: %v", err)
	}
	// teleport
	networkID, ok := blockNBT_depends.ItemRunTimeID[FrameData.Item[0].Name]
	if ok {
		clickNum := int(math.Floor(float64(FrameData.ItemRotation)/45)) + 1
		// get the click counts
		// 0.0 ≤ FrameData.ItemRotation ≤ 315.0
		var blockRuntimeID uint32 = 0
		if *BlockInfo.Block.Name == "glow_frame" {
			blockRuntimeID = glowFrameBlockRunTimeId
		}
		if *BlockInfo.Block.Name == "frame" {
			blockRuntimeID = frameBlockRunTimeId
		}
		blockRuntimeID = blockRuntimeID + uint32(facing_direction)
		// get run time id of frame block
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
	// put item into the frame
	return nil
}
