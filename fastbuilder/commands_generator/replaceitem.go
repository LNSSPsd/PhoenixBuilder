package commands_generator

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
)

func ReplaceItemRequest(module *types.Module, config *types.MainConfig) []string {
	if module.ChestSlot != nil {
		return []string{fmt.Sprintf("replaceitem block %d %d %d slot.container %d %s %d %d", module.Point.X, module.Point.Y, module.Point.Z, module.ChestSlot.Slot, module.ChestSlot.Name, module.ChestSlot.Count, module.ChestSlot.Damage)}
	} else {
		ans := []string{}
		for _, value := range *module.ChestData {
			ans = append(ans, fmt.Sprintf("replaceitem block %d %d %d slot.container %d %s %d %d", module.Point.X, module.Point.Y, module.Point.Z, value.Slot, value.Name, value.Count, value.Damage))
		}
		return ans
	}
}
