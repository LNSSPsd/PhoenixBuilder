package blockNBT_API

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"phoenixbuilder/mirror/chunk"
)

/*
让客户端点击 pos 处名为 blockName 且方块状态为 blockStates 的方块

如果 needWaiting 为真，则会等待点击完成后再返回值

你可以对容器使用这样的操作，这会使得容器被打开
*/
func (g *GlobalAPI) UseItemOnBlocks(
	hotBarSlotID uint8,
	pos [3]int32,
	blockName string,
	blockStates map[string]interface{},
	needWaiting bool,
) error {
	standardRuntimeID, found := chunk.StateToRuntimeID(blockName, blockStates)
	if !found {
		return fmt.Errorf("UseItemOnBlocks: Failed to get the runtimeID of block %v; blockStates = %#v", blockName, blockStates)
	}
	blockRuntimeID := chunk.StandardRuntimeIDToNEMCRuntimeID(standardRuntimeID)
	if blockRuntimeID == chunk.AirRID || blockRuntimeID == chunk.NEMCAirRID {
		return fmt.Errorf("UseItemOnBlocks: Failed to converse StandardRuntimeID to NEMCRuntimeID; standardRuntimeID = %#v, blockName = %#v, blockStates = %#v", standardRuntimeID, blockName, blockStates)
	}
	// get block RunTime ID
	err := g.ChangeSelectedHotbarSlot(hotBarSlotID, true)
	if err != nil {
		return fmt.Errorf("UseItemOnBlocks: %v", err)
	}
	// change selected hotbar slot
	datas, _ := g.GetInventoryCotent(0)
	// get item contents of window 0
	got, ok := datas[hotBarSlotID]
	if !ok {
		return fmt.Errorf("UseItemOnBlocks: %v is not in inventory contents; datas = %#v", hotBarSlotID, datas)
	}
	// get target item datas
	err = g.WritePacket(&packet.InventoryTransaction{
		LegacyRequestID:    0,
		LegacySetItemSlots: []protocol.LegacySetItemSlot(nil),
		Actions:            []protocol.InventoryAction{},
		TransactionData: &protocol.UseItemTransactionData{
			LegacyRequestID:    0,
			LegacySetItemSlots: []protocol.LegacySetItemSlot(nil),
			Actions:            []protocol.InventoryAction(nil),
			ActionType:         protocol.UseItemActionClickBlock,
			BlockPosition:      pos,
			BlockFace:          0,
			HotBarSlot:         int32(hotBarSlotID),
			HeldItem:           got,
			BlockRuntimeID:     blockRuntimeID,
		},
	})
	if err != nil {
		return fmt.Errorf("UseItemOnBlocks: %v", err)
	}
	// write packet
	if needWaiting {
		_, err = g.SendWSCommandWithResponce("list")
		if err != nil {
			return fmt.Errorf("UseItemOnBlocks: %v", err)
		}
	}
	// wait changes
	return nil
	// return
}
