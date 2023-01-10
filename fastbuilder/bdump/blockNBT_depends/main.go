package blockNBT_depends

import (
	"bytes"
	"fmt"
	"phoenixbuilder/fastbuilder/commands_generator"
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/io/commands"

	"github.com/Tnze/go-mc/nbt"
)

// 此结构体用于本文件中 PlaceBlockWithNBTData 函数的输入部分
type input struct {
	Environment        *environment.PBEnvironment // 运行环境（必须）
	Mainsettings       *types.MainConfig          // 一些设置
	IsFastMode         bool                       // 是否是快速模式
	BlockInfo          *types.Module              // 用于存放方块信息
	BlockNBT           *map[string]interface{}    // 用于存放方块实体数据
	TypeName           *string                    // 用于存放这种方块的类型，比如不同的告示牌都可以写成 sign
	OtherNecessaryData *interface{}               // 存放其他一些必要数据
}

// 此表用于记录现阶段支持了的方块实体
var index = map[string]string{
	"command_block":           "command_block",
	"chain_command_block":     "command_block",
	"repeating_command_block": "command_block",
	// 命令方块
	"blast_furnace":      "Container",
	"lit_blast_furnace":  "Container",
	"smoker":             "Container",
	"lit_smoker":         "Container",
	"furnace":            "Container",
	"lit_furnace":        "Container",
	"chest":              "Container",
	"barrel":             "Container",
	"trapped_chest":      "Container",
	"lectern":            "Container",
	"hopper":             "Container",
	"dispenser":          "Container",
	"dropper":            "Container",
	"cauldron":           "Container",
	"lava_cauldron":      "Container",
	"jukebox":            "Container",
	"brewing_stand":      "Container",
	"undyed_shulker_box": "Container",
	"shulker_box":        "Container",
	// 容器
	"standing_sign":          "Sign",
	"spruce_standing_sign":   "Sign",
	"birch_standing_sign":    "Sign",
	"jungle_standing_sign":   "Sign",
	"acacia_standing_sign":   "Sign",
	"darkoak_standing_sign":  "Sign",
	"mangrove_standing_sign": "Sign",
	"bamboo_standing_sign":   "Sign",
	"crimson_standing_sign":  "Sign",
	"warped_standing_sign":   "Sign",
	"wall_sign":              "Sign",
	"spruce_wall_sign":       "Sign",
	"birch_wall_sign":        "Sign",
	"jungle_wall_sign":       "Sign",
	"acacia_wall_sign":       "Sign",
	"darkoak_wall_sign":      "Sign",
	"mangrove_wall_sign":     "Sign",
	"bamboo_wall_sign":       "Sign",
	"crimson_wall_sign":      "Sign",
	"warped_wall_sign":       "Sign",
	"sign":                   "Sign",
	"spruce_sign":            "Sign",
	"birch_sign":             "Sign",
	"jungle_sign":            "Sign",
	"acacia_sign":            "Sign",
	"darkoak_sign":           "Sign",
	"mangrove_sign":          "Sign",
	"bamboo_sign":            "Sign",
	"crimson_sign":           "Sign",
	"warped_sign":            "Sign",
	"oak_hanging_sign":       "Sign",
	"spruce_hanging_sign":    "Sign",
	"birch_hanging_sign":     "Sign",
	"jungle_hanging_sign":    "Sign",
	"acacia_hanging_sign":    "Sign",
	"dark_oak_hanging_sign":  "Sign",
	"mangrove_hanging_sign":  "Sign",
	"bamboo_hanging_sign":    "Sign",
	"crimson_hanging_sign":   "Sign",
	"warped_hanging_sign":    "Sign",
	// 告示牌
}

// 检查这个方块实体是否已被支持
func checkIfIsEffectiveNBTBlock(blockName string) string {
	value, ok := index[blockName]
	if ok {
		return value
	}
	return ""
}

// 带有 NBT 数据放置方块；返回值 interface{} 字段可能在后期会用到，但目前这个字段都是返回 nil
func placeBlockWithNBTData(input *input) (interface{}, error) {
	var err error
	// prepare
	switch *input.TypeName {
	case "command_block":
		err = CommandBlock(&CommandBlockInput{
			Cb:           input.BlockNBT,
			BlockName:    input.BlockInfo.Block.Name,
			Environment:  input.Environment,
			Mainsettings: input.Mainsettings,
			IsFastMode:   input.IsFastMode,
			BlockInfo:    input.BlockInfo,
		})
		if err != nil {
			return nil, fmt.Errorf("placeBlockWithNBTData: %v", err)
		}
		// 命令方块
	case "Container":
		err = Container(&ContainerInput{
			ContainerData: input.BlockNBT,
			Environment:   input.Environment,
			Mainsettings:  input.Mainsettings,
			BlockInfo:     input.BlockInfo,
		})
		if err != nil {
			return nil, fmt.Errorf("placeBlockWithNBTData: %v", err)
		}
		// 各类可被 replaceitem 生效的容器
	case "Sign":
		err = Sign(&SignInput{
			Environment:  input.Environment,
			Mainsettings: input.Mainsettings,
			BlockInfo:    input.BlockInfo,
			Sign:         input.BlockNBT,
		})
		if err != nil {
			return nil, fmt.Errorf("placeBlockWithNBTData: %v", err)
		}
		// 告示牌
	default:
		request := commands_generator.SetBlockRequest(input.BlockInfo, input.Mainsettings)
		cmdsender := input.Environment.CommandSender.(*commands.CommandSender)
		cmdsender.SendDimensionalCommand(request)
		return nil, nil
	}
	return nil, nil
}

func PlaceBlockWithNBTDataRun(
	Environment *environment.PBEnvironment,
	Mainsettings *types.MainConfig,
	IsFastMode bool,
	BlockInfo *types.Module,
) error {
	var buf bytes.Buffer
	err := nbt.NewEncoder(&buf).Encode(nbt.StringifiedMessage(*BlockInfo.StringNBT), "")
	if err != nil {
		return fmt.Errorf("PlaceBlockWithNBTDataRun: %v", err)
	}
	var BlockNBT map[string]interface{}
	nbt.Unmarshal(buf.Bytes(), &BlockNBT)
	// get inerface NBT, saved in BlockNBT
	TYPE := checkIfIsEffectiveNBTBlock(*BlockInfo.Block.Name)
	_, err = placeBlockWithNBTData(&input{
		Environment:        Environment,
		Mainsettings:       Mainsettings,
		IsFastMode:         IsFastMode,
		BlockInfo:          BlockInfo,
		BlockNBT:           &BlockNBT,
		TypeName:           &TYPE,
		OtherNecessaryData: nil,
	})
	if err != nil {
		return fmt.Errorf("PlaceBlockWithNBTDataRun: %v", err)
	}
	return nil
}
