package commands_generator

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
	"strings"
)

func ReplaceItemRequest(module *types.Module, config *types.MainConfig) []string {
	if module.ChestSlot != nil {
		return []string{fmt.Sprintf("replaceitem block %d %d %d slot.container %d %s %d %d", module.Point.X, module.Point.Y, module.Point.Z, module.ChestSlot.Slot, module.ChestSlot.Name, module.ChestSlot.Count, module.ChestSlot.Damage)}
	} else {
		ans := []string{}
		for key, value := range *module.ChestData {
			ans = append(ans, fmt.Sprintf("replaceitem block %d %d %d slot.container %d %s %d %d %v", module.Point.X, module.Point.Y, module.Point.Z, value.Slot, value.Name, value.Count, value.Damage, GetReplaceItemEnhancement(module, key)))
		}
		return ans
	}
}

func GetReplaceItemEnhancement(module *types.Module, location int) string {
	ans := make([]string, 0)
	i := *module.ChestData
	value := i[location]
	// prepare
	single := make([]string, 0)
	for _, VALUE := range value.CanPlaceOn {
		single = append(single, fmt.Sprintf(`"%v"`, VALUE))
	}
	if len(single) > 0 {
		ans = append(ans, fmt.Sprintf(`"can_place_on": {"blocks": [%v]}`, strings.Join(single, ", ")))
	}
	// can_place_on
	single = []string{}
	for _, VALUE := range value.CanDestroy {
		single = append(single, fmt.Sprintf(`"%v"`, VALUE))
	}
	if len(single) > 0 {
		ans = append(ans, fmt.Sprintf(`"can_destroy": {"blocks": [%v]}`, strings.Join(single, ", ")))
	}
	// can_destroy
	if value.ItemLock == 1 {
		ans = append(ans, `"item_lock": {"mode": "lock_in_slot"}`)
	}
	if value.ItemLock == 2 {
		ans = append(ans, `"item_lock": {"mode": "lock_in_inventory"}`)
	}
	// item_lock
	if value.KeepOnDeath == 1 {
		ans = append(ans, `"keep_on_death": {}`)
	}
	// keep_on_death
	if len(ans) > 0 {
		return fmt.Sprintf("{%v}", strings.Join(ans, ", "))
	}
	return ""
}
