package commands_generator

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
	"strconv"
)

func SetBlockRequest(module *types.Module, config *types.MainConfig) string {
	Block := module.Block
	Point := module.Point
	Method := config.Method
	if Block != nil {
		if Block.BlockStates == "" {
			return fmt.Sprintf("setblock %d %d %d %s %d %s", Point.X, Point.Y, Point.Z, *Block.Name, Block.Data, Method)
		} else if Block.BlockStates == "[]" {
			return fmt.Sprintf("setblock %d %d %d %s %d %s", Point.X, Point.Y, Point.Z, *Block.Name, Block.Data, Method)

		} else {
			if IsNum(Block.BlockStates) {
				return fmt.Sprintf("setblock %d %d %d %s %s %s", Point.X, Point.Y, Point.Z, *Block.Name, Block.BlockStates, Method)

			} else {
				return fmt.Sprintf("setblock %d %d %d %s %d %s", Point.X, Point.Y, Point.Z, *Block.Name, Block.Data, Method)
			}

		}
	} else {
		return fmt.Sprintf("setblock %d %d %d %s %d %s", Point.X, Point.Y, Point.Z, config.Block.Name, config.Block.Data, Method)
	}

}

// 这随便写的一个临时函数,随便换个位置都行,偷个懒
// Art 最喜欢摸鱼了
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
