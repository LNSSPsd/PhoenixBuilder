package NBTAssigner

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
	GameInterface "phoenixbuilder/game_control/game_interface"
)

// 从容器的 NBT 数据中提取物品数据。
// 返回的列表代表提取到的每个物品
func (c *Container) getContainerContents() ([]ItemOrigin, error) {
	key := SupportContainerPool[c.BlockEntity.Block.Name]
	if len(key) == 0 {
		return nil, ErrNotASupportedContainer
	}
	// 确定目标容器是否已被支持
	itemContentsGot, ok := c.BlockEntity.Block.NBT[key]
	// 从 containerOriginNBT 获取物品的数据
	if !ok {
		return []ItemOrigin{}, nil
	}
	// 对于唱片机和讲台这种容器，如果它们没有被放物品的话，
	// 那么对应的 key 是找不到的，但是这并非是错误
	switch itemContents := itemContentsGot.(type) {
	case map[string]interface{}:
		return []ItemOrigin{itemContents}, nil
		// 如果这个物品是一个唱片机或者讲台，
		// 那么传入的 itemContents 是一个复合标签而非列表。
		// 因此，为了统一数据格式，
		// 我们将复合标签处理成通常情况下的列表
	case []interface{}:
		res := []ItemOrigin{}
		for key, value := range itemContents {
			singleItem, success := value.(map[string]interface{})
			if !success {
				return nil, fmt.Errorf(`getContainerContents: Crashed on itemContents[%d]; itemContents = %#v`, key, itemContents)
			}
			res = append(res, singleItem)
		}
		return res, nil
		// 常规型物品的(多个)物品数据存放在一张表中，
		// 而每个物品都用一个复合标签来描述
	default:
		return nil, fmt.Errorf(`getContainerContents: Unexpected data type of itemContentsGot; itemContentsGot = %#v`, itemContentsGot)
	}
	// 处理方块实体数据并返回值
}

// 从 c.Package.Block.NBT 提取物品数据并保存在 c.Contents 中
func (c *Container) Decode() error {
	// 初始化
	itemContents, err := c.getContainerContents()
	if err != nil {
		return fmt.Errorf("Decode: %v", err)
	}
	// 获取容器内的物品数据
	for _, value := range itemContents {
		got, err := ParseItemFromNBT(value, SupportBlocksPool)
		if err != nil {
			return fmt.Errorf("Decode: %v", err)
		}
		c.Contents = append(c.Contents, got)
	}
	// 解码
	return nil
	// 返回值
}

// 放置一个容器并填充物品
func (c *Container) WriteData() error {
	err := c.BlockEntity.Interface.SetBlock(c.BlockEntity.AdditionalData.Position, c.BlockEntity.Block.Name, c.BlockEntity.AdditionalData.BlockStates)
	if err != nil {
		return fmt.Errorf("WriteData: %v", err)
	}
	// 放置容器
	for _, value := range c.Contents {
		err := c.BlockEntity.Interface.(*GameInterface.GameInterface).ReplaceItemInContainer(
			c.BlockEntity.AdditionalData.Position,
			types.ChestSlot{
				Name:   value.Basic.Name,
				Count:  value.Basic.Count,
				Damage: value.Basic.MetaData,
				Slot:   value.Basic.Slot,
			},
			"",
		)
		if err != nil {
			return fmt.Errorf("WriteData: %v", err)
		}
	}
	// 向容器内填充物品
	return nil
	// 返回值
}
