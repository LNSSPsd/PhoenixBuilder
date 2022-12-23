//go:build do_not_add_this_tag__not_implemented
// +build do_not_add_this_tag__not_implemented

package commands_generator

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
)

func SummonRequest(module *types.Module, config *types.MainConfig, BotName string) string {
	Entity := module.Entity
	Point := module.Point
	Method := config.Method
	if Entity != nil {
		return fmt.Sprintf("execute @a[name=%v] ~ ~ ~ summon %s %v %v %v", BotName, *Entity, Point.X, Point.Y, Point.Z)
	} else {
		return fmt.Sprintf("execute @a[name=%v] ~ ~ ~ summon %s %v %v %v", BotName, *Entity, Point.X, Point.Y, Point.Z)
	}
}
