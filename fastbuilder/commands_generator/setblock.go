package commands_generator

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
)

func SetBlockRequest(module *types.Module, config *types.MainConfig, BotName string) string {
	Block := module.Block
	Point := module.Point
	Method := config.Method
	if Block != nil {
		if len(Block.BlockStates) != 0 {
			return fmt.Sprintf("execute @a[name=\"%v\"] ~ ~ ~ setblock %d %d %d %s %s %s", BotName, Point.X, Point.Y, Point.Z, *Block.Name, Block.BlockStates, Method)
		} else {
			return fmt.Sprintf("execute @a[name=\"%v\"] ~ ~ ~ setblock %d %d %d %s %d %s", BotName, Point.X, Point.Y, Point.Z, *Block.Name, Block.Data, Method)
		}
	} else {
		return fmt.Sprintf("execute @a[name=\"%v\"] ~ ~ ~ setblock %d %d %d %s %d %s", BotName, Point.X, Point.Y, Point.Z, config.Block.Name, config.Block.Data, Method)
	}

}
