package NBTAssigner

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
	GameInterface "phoenixbuilder/game_control/game_interface"
	"strings"
)

// 获取一个潜影盒到快捷栏 5 。
// 此函数仅应当在放置潜影盒时被使用
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

// 从 c.BlockEntity.Block.NBT 获取潜影盒的朝向。
// 此函数仅应当在放置潜影盒时被使用。
// 如果朝向不存在，则重定向为 1(朝上)
func (c *Container) getFacingOfShulkerBox() (uint8, error) {
	if facing_origin, ok := c.BlockEntity.Block.NBT["facing"]; ok {
		facing_got, success := facing_origin.(byte)
		if !success {
			return 0, fmt.Errorf(`getFacingOfShulkerBox: Can not convert facing_origin into byte(uint8); c.BlockEntity.Block.NBT = %#v`, c.BlockEntity.Block.NBT)
		}
		return facing_got, nil
	}
	return 1, nil
}

// 放置 c.BlockEntity 所代表的容器。
// 此函数侧重于对潜影盒的专门化处理，
// 以保证放置出的潜影盒能拥有正确的朝向
func (c *Container) PlaceContainer() error {
	api := c.BlockEntity.Interface.(*GameInterface.GameInterface)
	// 初始化
	if strings.Contains(c.BlockEntity.Block.Name, "shulker_box") {
		facing, err := c.getFacingOfShulkerBox()
		if err != nil {
			return fmt.Errorf("PlaceContainer: %v", err)
		}
		// 获取潜影盒的朝向
		err = api.SendSettingsCommand(
			fmt.Sprintf(
				"tp %d %d %d",
				c.BlockEntity.AdditionalData.Position[0],
				c.BlockEntity.AdditionalData.Position[1],
				c.BlockEntity.AdditionalData.Position[2],
			),
			true,
		)
		if err != nil {
			return fmt.Errorf("PlaceContainer: %v", err)
		}
		// 将机器人传送到潜影盒处
		err = c.getShulkerBox()
		if err != nil {
			return fmt.Errorf("PlaceContainer: %v", err)
		}
		// 获取一个潜影盒到快捷栏 5
		err = api.PlaceShulkerBox(c.BlockEntity.AdditionalData.Position, 5, facing)
		if err != nil {
			return fmt.Errorf("PlaceContainer: %v", err)
		}
		// 生成潜影盒
	} else {
		err := api.SetBlock(
			c.BlockEntity.AdditionalData.Position,
			c.BlockEntity.Block.Name,
			c.BlockEntity.AdditionalData.BlockStates,
		)
		if err != nil {
			return fmt.Errorf("PlaceContainer: %v", err)
		}
	}
	// 放置容器
	return nil
	// 返回值
}

/*
打开已放置的容器，因此该函数应当后于 PlaceContainer 执行。

如果该容器不可被打开，则返回假，否则返回真。

此函数不会等待租赁服响应更改，它不是阻塞式的实现。
您应当在上层实现中占用容器资源并测定容器是否被打开
*/
func (c *Container) OpenContainer() (bool, error) {
	api := c.BlockEntity.Interface.(*GameInterface.GameInterface)
	backupBlockPos := c.BlockEntity.AdditionalData.Position
	// 初始化
	if c.BlockEntity.Block.Name == "lectern" || c.BlockEntity.Block.Name == "jukebox" {
		return false, nil
	}
	// 如果这个容器不能打开
	if strings.Contains(c.BlockEntity.Block.Name, "shulker_box") || strings.Contains(c.BlockEntity.Block.Name, "chest") {
		if strings.Contains(c.BlockEntity.Block.Name, "shulker_box") {
			facing, err := c.getFacingOfShulkerBox()
			if err != nil {
				return false, fmt.Errorf("OpenContainer: %v", err)
			}
			switch facing {
			case 0:
				backupBlockPos[1] = backupBlockPos[1] - 1
			case 1:
				backupBlockPos[1] = backupBlockPos[1] + 1
			case 2:
				backupBlockPos[2] = backupBlockPos[2] - 1
			case 3:
				backupBlockPos[2] = backupBlockPos[2] + 1
			case 4:
				backupBlockPos[0] = backupBlockPos[0] - 1
			case 5:
				backupBlockPos[0] = backupBlockPos[0] + 1
			}
		} else {
			backupBlockPos[1] = backupBlockPos[1] + 1
		}
		// 确定容器开启方向上前一格方块的位置
		uniqueId, err := api.BackupStructure(GameInterface.MCStructure{
			BeginX: backupBlockPos[0],
			BeginY: backupBlockPos[1],
			BeginZ: backupBlockPos[2],
			SizeX:  1,
			SizeY:  1,
			SizeZ:  1,
		})
		if err != nil {
			return false, fmt.Errorf("OpenContainer: %v", err)
		}
		defer func() {
			api.RevertStructure(uniqueId, backupBlockPos)
		}()
		err = api.SendSettingsCommand(
			fmt.Sprintf(
				"kill @e[x=%d,y=%d,z=%d,dx=0]",
				backupBlockPos[0],
				backupBlockPos[1],
				backupBlockPos[2],
			),
			true,
		)
		if err != nil {
			return false, fmt.Errorf("OpenContainer: %v", err)
		}
		err = api.SetBlockAsync(backupBlockPos, "air", "[]")
		if err != nil {
			return false, fmt.Errorf("OpenContainer: %v", err)
		}
		/*
			我们需要保证潜影盒开启方向上的方块为空气且没有生物，
			否则潜影盒将无法正常开启。
			然而，对这个方块进行操作并杀死该处的生物不是预期的行为，
			所以需要确定其坐标并发起一次备份，
			然后强行将其变更为空气并执行一次 kill 命令
		*/
	}
	// 对潜影盒或者箱子的特殊化处理
	err := api.ChangeSelectedHotbarSlot(5)
	if err != nil {
		return false, fmt.Errorf("OpenContainer: %v", err)
	}
	err = api.ClickBlock(GameInterface.UseItemOnBlocks{
		HotbarSlotID: 5,
		BlockPos:     c.BlockEntity.AdditionalData.Position,
		BlockName:    c.BlockEntity.Block.Name,
		BlockStates:  c.BlockEntity.Block.States,
	})
	if err != nil {
		return false, fmt.Errorf("OpenContainer: %v", err)
	}
	// 将快捷栏切换至 5 号槽位，
	// 然后使用该槽位的物品点击容器，
	// 以达到开启容器的目的
	return true, nil
	// 返回值
}
