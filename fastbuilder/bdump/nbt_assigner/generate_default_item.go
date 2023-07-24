package NBTAssigner

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
	GameInterface "phoenixbuilder/game_control/game_interface"
)

// DefaultItem 结构体用于描述一个完整的 NBT 物品的数据。
// 任何未被支持的 NBT 物品都会被重定向为此结构体
type DefaultItem struct {
	ItemPackage *ItemPackage // 该 NBT 物品的详细数据
}

// 这只是为了保证接口一致而设
func (d *DefaultItem) Decode() error {
	return nil
}

// 生成目标物品到快捷栏但不写入 NBT 数据
func (d *DefaultItem) WriteData() error {
	err := d.ItemPackage.Interface.(*GameInterface.GameInterface).ReplaceItemInInventory(
		GameInterface.TargetMySelf,
		GameInterface.ItemGenerateLocation{
			Path: "slot.hotbar",
			Slot: d.ItemPackage.AdditionalData.HotBarSlot,
		},
		types.ChestSlot{
			Name:   d.ItemPackage.Item.Basic.Name,
			Count:  d.ItemPackage.Item.Basic.Count,
			Damage: d.ItemPackage.Item.Basic.MetaData,
		},
		MarshalItemComponents(d.ItemPackage.Item.Enhancement.ItemComponents),
	)
	if err != nil {
		return fmt.Errorf("WriteData: %v", err)
	}
	return nil
}
