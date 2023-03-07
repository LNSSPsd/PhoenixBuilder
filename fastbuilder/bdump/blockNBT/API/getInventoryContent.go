package blockNBT_API

import (
	"phoenixbuilder/minecraft/protocol"
)

// 用于获取库存数据；当 WindowID 可以被找到时，返回的布尔值为真，否则为假
func (g *GlobalAPI) GetInventoryCotent(WindowID uint32) (map[uint8]protocol.ItemInstance, bool) {
	g.PacketHandleResult.InventoryDatasMapLockDown.RLock()
	defer g.PacketHandleResult.InventoryDatasMapLockDown.RUnlock()
	// init
	ans, success := g.PacketHandleResult.InventoryDatas[WindowID]
	if !success {
		return map[uint8]protocol.ItemInstance{}, false
	}
	// get datas
	return ans, true
	// return
}
