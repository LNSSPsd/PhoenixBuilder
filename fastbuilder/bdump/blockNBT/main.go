package blockNBT

import (
	"fmt"
	blockNBT_depends "phoenixbuilder/fastbuilder/bdump/blockNBT/depends"
	"phoenixbuilder/fastbuilder/commands_generator"
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/io/commands"
	"phoenixbuilder/minecraft/nbt"
	"sync"
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

/*
带有 NBT 数据放置方块；返回值 interface{} 字段可能在后期会用到，但目前这个字段都是返回 nil

如果你也想参与更多方块实体的支持，可以去看看这个库 https://github.com/df-mc/dragonfly

这个库也是用了 gophertunnel 的
*/
func placeBlockWithNBTData(input *input) (interface{}, error) {
	var err error
	// prepare
	switch *input.TypeName {
	case "CommandBlock":
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
	case "Frame":
		err = Frame(&FrameInput{
			Environment:  input.Environment,
			Mainsettings: input.Mainsettings,
			BlockInfo:    input.BlockInfo,
			Frame:        input.BlockNBT,
		})
		if err != nil {
			return nil, fmt.Errorf("placeBlockWithNBTData: %v", err)
		}
		// 物品展示框
	case "lectern":
		err = Lectern(&LecternInput{
			Environment:  input.Environment,
			Mainsettings: input.Mainsettings,
			BlockInfo:    input.BlockInfo,
			LecternData:  input.BlockNBT,
		})
		if err != nil {
			return nil, fmt.Errorf("placeBlockWithNBTData: %v", err)
		}
		// 讲台
	default:
		request := commands_generator.SetBlockRequest(input.BlockInfo, input.Mainsettings)
		cmdsender := input.Environment.CommandSender.(*commands.CommandSender)
		cmdsender.SendDimensionalCommand(request)
		return nil, nil
		// 其他没有支持的方块实体
	}
	return nil, nil
}

var apiIsUsing sync.Mutex

// 此函数是 package blockNBT 的主函数
func PlaceBlockWithNBTDataRun(
	Environment *environment.PBEnvironment,
	Mainsettings *types.MainConfig,
	IsFastMode bool,
	BlockInfo *types.Module,
) error {
	defer apiIsUsing.Unlock()
	apiIsUsing.Lock()
	// lock(or unlock) api
	var BlockNBT map[string]interface{}
	err := nbt.Unmarshal(BlockInfo.NBTData, &BlockNBT)
	if err != nil {
		return fmt.Errorf("PlaceBlockWithNBTDataRun: Failed to place the entity block named %v at (%v,%v,%v), and the error log is %v", *BlockInfo.Block.Name, BlockInfo.Point.X, BlockInfo.Point.Y, BlockInfo.Point.Z, err)
	}
	// get interface nbt
	TYPE := blockNBT_depends.CheckIfIsEffectiveNBTBlock(*BlockInfo.Block.Name)
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
