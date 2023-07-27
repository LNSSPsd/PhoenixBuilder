package NBTAssigner

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
	GameInterface "phoenixbuilder/game_control/game_interface"
	ResourcesControl "phoenixbuilder/game_control/resources_control"
	"phoenixbuilder/minecraft/protocol"
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
	item := d.ItemPackage.Item
	api := d.ItemPackage.Interface.(*GameInterface.GameInterface)
	// 初始化
	err := api.ReplaceItemInInventory(
		GameInterface.TargetMySelf,
		GameInterface.ItemGenerateLocation{
			Path: "slot.hotbar",
			Slot: d.ItemPackage.AdditionalData.HotBarSlot,
		},
		types.ChestSlot{
			Name:   item.Basic.Name,
			Count:  item.Basic.Count,
			Damage: item.Basic.MetaData,
		},
		MarshalItemComponents(item.Enhancement.ItemComponents),
	)
	if err != nil {
		return fmt.Errorf("WriteData: %v", err)
	}
	// 获取物品到物品栏，并附加物品组件数据
	if item.Enhancement != nil && item.Enhancement.Enchantments != nil {
		for _, value := range *item.Enhancement.Enchantments {
			err = api.SendSettingsCommand(
				fmt.Sprintf(
					"enchant @s %d %d",
					value.ID,
					value.Level,
				),
				true,
			)
			if err != nil {
				return fmt.Errorf("WriteData: %v", err)
			}
		}
		resp := api.SendWSCommandWithResponse("list")
		if resp.Error != nil && resp.ErrorType != ResourcesControl.ErrCommandRequestTimeOut {
			return fmt.Errorf("WriteData: %v", err)
		}
	}
	// 附加附魔属性
	if item.Enhancement != nil && item.Enhancement.ItemComponents != nil && len(item.Enhancement.ItemComponents.ItemLock) != 0 {
		return nil
	}
	// 如果该物品存在 item_lock 物品组件，
	// 则后续 NBT 无需附加，
	// 因为带有该物品组件的物品不能跨容器移动
	if item.Enhancement != nil && len(item.Enhancement.DisplayName) != 0 {
		resp, err := api.RenameItemByAnvil(
			d.ItemPackage.AdditionalData.Position,
			`["direction": 0, "damage": "undamaged"]`,
			5,
			[]GameInterface.ItemRenamingRequest{
				{
					Slot: d.ItemPackage.AdditionalData.HotBarSlot,
					Name: item.Enhancement.DisplayName,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("WriteData: %v", err)
		}
		if resp[0].Destination == nil {
			return fmt.Errorf("WriteData: Inventory was full")
		}
		// 利用铁砧修改物品名称
		if resp[0].Destination.Slot != d.ItemPackage.AdditionalData.HotBarSlot {
			itemData, err := api.Resources.Inventory.GetItemStackInfo(0, resp[0].Destination.Slot)
			if err != nil {
				return fmt.Errorf("WriteData: %v", err)
			}
			// 获取已被铁砧操作后的物品数据
			err = api.ReplaceItemInInventory(
				GameInterface.TargetMySelf,
				GameInterface.ItemGenerateLocation{
					Path: "slot.hotbar",
					Slot: d.ItemPackage.AdditionalData.HotBarSlot,
				},
				types.ChestSlot{
					Name:   "air",
					Count:  1,
					Damage: 0,
				},
				"",
			)
			if err != nil {
				return fmt.Errorf("WriteData: %v", err)
			}
			// 将原有物品栏替换为空气以解除它的占用态
			res, err := api.MoveItem(
				GameInterface.ItemLocation{
					WindowID:    0,
					ContainerID: 0xc,
					Slot:        resp[0].Destination.Slot,
				},
				GameInterface.ItemLocation{
					WindowID:    0,
					ContainerID: 0xc,
					Slot:        d.ItemPackage.AdditionalData.HotBarSlot,
				},
				GameInterface.ItemChangingDetails{
					Details: map[ResourcesControl.ContainerID]ResourcesControl.StackRequestContainerInfo{
						0xc: {
							WindowID: 0,
							ChangeResult: map[uint8]protocol.ItemInstance{
								resp[0].Destination.Slot:                GameInterface.AirItem,
								d.ItemPackage.AdditionalData.HotBarSlot: itemData,
							},
						},
					},
				},
				uint8(itemData.Stack.Count),
			)
			if err != nil {
				return fmt.Errorf("WriteData: %v", err)
			}
			if res[0].Status != protocol.ItemStackResponseStatusOK {
				return fmt.Errorf("WriteData: Failed to restore the item to its original position")
			}
			// 尝试将物品恢复到原始位置
		}
	}
	// 附加物品的显示名称
	return nil
	// 返回值
}
