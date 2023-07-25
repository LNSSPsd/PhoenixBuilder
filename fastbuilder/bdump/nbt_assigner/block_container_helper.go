package NBTAssigner

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
	GameInterface "phoenixbuilder/game_control/game_interface"
)

// 获取一个潜藏盒到快捷栏 5 。
// 此函数仅应当在放置潜藏盒时被使用
func (c *Container) getShulkerBox() error {
	var blockMetaData uint16
	api := c.BlockEntity.Interface.(*GameInterface.GameInterface)
	// 初始化
	blockMetaData, _ = get_block_data_from_states(
		c.BlockEntity.Block.Name,
		c.BlockEntity.Block.States,
	)
	// 取得潜影盒的方块数据值(附加值)
	err := api.ReplaceItemInInventory(
		GameInterface.TargetMySelf,
		GameInterface.ItemGenerateLocation{
			Path: "slot.hotbar",
			Slot: 5,
		},
		types.ChestSlot{
			Name:   c.BlockEntity.Block.Name,
			Count:  1,
			Damage: blockMetaData,
		},
		"",
	)
	if err != nil {
		return fmt.Errorf("GetShulkerBox: %v", err)
	}
	// 将潜影盒替换至快捷栏 5
	return nil
	// 返回值
}

/*
func (c *Container) PlaceContainer() {
	api := c.BlockEntity.Interface.(*GameInterface.GameInterface)
	// 初始化
	if strings.Contains(c.BlockEntity.Block.Name, "shulker_box") {

	}
}
*/
