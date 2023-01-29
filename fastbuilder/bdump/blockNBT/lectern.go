package blockNBT

import (
	"fmt"
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

const lecternBlockRunTimeId = 6857

type LecternInput struct {
	Environment  *environment.PBEnvironment
	Mainsettings *types.MainConfig
	BlockInfo    *types.Module
	LecternData  *map[string]interface{}
}

func Lectern(input *LecternInput) error {
	err := placeLectern(input.Environment, input.Mainsettings, input.BlockInfo, *input.LecternData)
	if err != nil {
		return fmt.Errorf("Frame: Failed to place the entity block named %v at (%v,%v,%v), and the error log is %v", *input.BlockInfo.Block.Name, input.BlockInfo.Point.X, input.BlockInfo.Point.Y, input.BlockInfo.Point.Z, err)
	}
	return nil
}

func parseLecternData(input *map[string]interface{}, BlockName *string) (*types.ChestData, error) {
	got, err := getContainerDataRun(*input, *BlockName)
	if err != nil {
		return &types.ChestData{}, fmt.Errorf("parseLecternData: %v", err)
	}
	// get book info
	return &got, nil
}

func placeLectern(
	Environment *environment.PBEnvironment,
	Mainsettings *types.MainConfig,
	BlockInfo *types.Module,
	LecternData map[string]interface{},
) error {
	Lectern, err := parseLecternData(&LecternData, BlockInfo.Block.Name)
	if err != nil {
		return fmt.Errorf("placeLectern: %v", err)
	}
	LecternItemData := *Lectern
	// parse lectern data
	cmdsender := Environment.CommandSender.(*commands.CommandSender)
	err = cmdsender.SendDimensionalCommand(fmt.Sprintf("setblock %v %v %v minecraft:air", BlockInfo.Point.X, BlockInfo.Point.Y, BlockInfo.Point.Z))
	if err != nil {
		return fmt.Errorf("placeLectern: %v", err)
	}
	// refresh current block
	blockStates, err := mcstructure.ParseStringNBT(BlockInfo.Block.BlockStates, true)
	if err != nil {
		return fmt.Errorf("placeLectern: %v", err)
	}
	BLOCKSTATES, normal := blockStates.(map[string]interface{})
	if !normal {
		return fmt.Errorf("placeLectern: Failed to converse BlockStates to map[string]interface{}, BlockStates = %#v", BlockInfo.Block.BlockStates)
	}
	_, ok := BLOCKSTATES["direction"]
	if !ok {
		return fmt.Errorf("placeLectern: Could not find BlockInfo.Block.BlockStates[\"direction\"]; BlockInfo.Block.BlockStates = %#v", BlockInfo.Block.BlockStates)
	}
	direction, normal := BLOCKSTATES["direction"].(int32)
	if !normal {
		return fmt.Errorf("placeLectern: Could not parse BlockInfo.Block.BlockStates[\"direction\"]; BlockInfo.Block.BlockStates = %#v", BlockInfo.Block.BlockStates)
	}
	if direction > 3 {
		return fmt.Errorf("placeLectern: Unexpected BlockInfo.Block.BlockStates[\"direction\"]; BlockInfo.Block.BlockStates = %#v", BlockInfo.Block.BlockStates)
	}
	// get direction
	BlockInfo.Block.BlockStates = fmt.Sprintf(`["direction": %v]`, direction)
	request := commands_generator.SetBlockRequest(BlockInfo, Mainsettings)
	_, err = cmdsender.SendWSCommandWithResponce(request)
	if err != nil {
		return fmt.Errorf("placeLectern: %v", err)
	}
	// place lectern block
	if len(LecternItemData) <= 0 || !blockNBT_depends.CheckVersion() {
		return nil
	}
	// if it is nothing in the lectern block or the current version are not support, then return nil
	err = blockNBT_depends.ReplaceitemAndEnchant(Environment, &types.ChestSlot{
		Name:       "writable_book",
		Count:      LecternItemData[0].Count,
		Damage:     LecternItemData[0].Damage,
		CanPlaceOn: LecternItemData[0].CanPlaceOn,
		CanDestroy: LecternItemData[0].CanDestroy,
		ItemNBT:    LecternItemData[0].ItemNBT,
	})
	if err != nil {
		return fmt.Errorf("placeLectern: %v", err)
	}
	err = blockNBT_depends.WriteTextToBook(Environment, &LecternItemData[0])
	if err != nil {
		return fmt.Errorf("placeLectern: %v", err)
	}
	// replaceitem and write text to the book
	_, err = cmdsender.SendWSCommandWithResponce(fmt.Sprintf("tp %v %v %v", BlockInfo.Point.X, BlockInfo.Point.Y, BlockInfo.Point.Z))
	if err != nil {
		return fmt.Errorf("placeLectern: %v", err)
	}
	// teleport
	networkID, ok := blockNBT_depends.ItemRunTimeID[LecternItemData[0].Name]
	if ok {
		blockRuntimeID := lecternBlockRunTimeId + uint32(direction)
		// get run time id of lectern block
		Environment.Connection.(*minecraft.Conn).WritePacket(&packet.InventoryTransaction{
			TransactionData: &protocol.UseItemTransactionData{
				BlockPosition: protocol.BlockPos{int32(BlockInfo.Point.X), int32(BlockInfo.Point.Y), int32(BlockInfo.Point.Z)},
				HotBarSlot:    0,
				HeldItem: protocol.ItemInstance{
					Stack: protocol.ItemStack{
						ItemType: protocol.ItemType{
							NetworkID: int32(networkID),
						},
						Count:         1,
						CanBePlacedOn: LecternItemData[0].CanPlaceOn,
						CanBreak:      LecternItemData[0].CanDestroy,
					},
				},
				BlockRuntimeID: blockRuntimeID,
			},
		})
	}
	// put item into the lectern
	return nil
}
