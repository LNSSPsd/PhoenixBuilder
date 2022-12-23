package commands_generator

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
)

func ReplaceItemRequest(module *types.Module, config *types.MainConfig, BotName string) string {
	return fmt.Sprintf("execute @a[name=%v] ~ ~ ~ replaceitem block %d %d %d slot.container %d %s %d %d", BotName, module.Point.X, module.Point.Y, module.Point.Z, module.ChestSlot.Slot, module.ChestSlot.Name, module.ChestSlot.Count, module.ChestSlot.Damage)
}
