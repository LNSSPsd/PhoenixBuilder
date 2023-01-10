package blockNBT_depends

import (
	"fmt"
	"phoenixbuilder/fastbuilder/commands_generator"
	environment "phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/io/commands"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"time"

	"github.com/google/uuid"
)

type CommandBlockInput struct {
	Cb           *map[string]interface{}
	BlockName    *string
	Environment  *environment.PBEnvironment
	Mainsettings *types.MainConfig
	IsFastMode   bool
	BlockInfo    *types.Module
}

func CommandBlock(input *CommandBlockInput) error {
	if input.Cb != nil {
		Cbdata, err := parseCommandBlockData(*input.Cb, *input.BlockInfo.Block.Name)
		if err != nil {
			return fmt.Errorf("CommandBlock: Failed to place the entity block named %v at (%v,%v,%v), and the error log is %v", *input.BlockInfo.Block.Name, input.BlockInfo.Point.X, input.BlockInfo.Point.Y, input.BlockInfo.Point.Z, err)
		}
		input.BlockInfo.CommandBlockData = &Cbdata
	}
	err := placeCommandBlock(input.Environment, input.Mainsettings, input.IsFastMode, input.BlockInfo)
	if err != nil {
		return fmt.Errorf("CommandBlock: Failed to place the entity block named %v at (%v,%v,%v), and the error log is %v", *input.BlockInfo.Block.Name, input.BlockInfo.Point.X, input.BlockInfo.Point.Y, input.BlockInfo.Point.Z, err)
	}
	return nil
}

func parseCommandBlockData(Cb map[string]interface{}, BlockName string) (types.CommandBlockData, error) {
	var normal bool = false
	var command string = ""
	var customName string = ""
	var lastOutput string = ""
	var mode int = 0
	var tickDelay int32 = int32(0)
	var executeOnFirstTick bool = true
	var trackOutput bool = true
	var conditionalMode bool = false
	var needRedstone bool = true
	// 初始化
	_, ok := Cb["Command"]
	if ok {
		command, normal = Cb["Command"].(string)
		if !normal {
			return types.CommandBlockData{}, fmt.Errorf("CommandBlockParseData: Crashed in Cb[\"Command\"]; Cb = %#v", Cb)
		}
	}
	// Command
	_, ok = Cb["CustomName"]
	if ok {
		customName, normal = Cb["CustomName"].(string)
		if !normal {
			return types.CommandBlockData{}, fmt.Errorf("CommandBlockParseData: Crashed in Cb[\"CustomName\"]; Cb = %#v", Cb)
		}
	}
	// CustomName
	_, ok = Cb["LastOutput"]
	if ok {
		lastOutput, normal = Cb["LastOutput"].(string)
		if !normal {
			return types.CommandBlockData{}, fmt.Errorf("CommandBlockParseData: Crashed in Cb[\"LastOutput\"]; Cb = %#v", Cb)
		}
	}
	// LastOutput
	if BlockName == "command_block" {
		mode = 0
	} else if BlockName == "repeating_command_block" {
		mode = 1
	} else if BlockName == "chain_command_block" {
		mode = 2
	} else {
		return types.CommandBlockData{}, fmt.Errorf("CommandBlockParseData: Not a command block; Cb = %#v", Cb)
	}
	// mode
	_, ok = Cb["TickDelay"]
	if ok {
		tickDelay, normal = Cb["TickDelay"].(int32)
		if !normal {
			return types.CommandBlockData{}, fmt.Errorf("CommandBlockParseData: Crashed in Cb[\"TickDelay\"]; Cb = %#v", Cb)
		}
	}
	// TickDelay
	_, ok = Cb["ExecuteOnFirstTick"]
	if ok {
		got, normal := Cb["ExecuteOnFirstTick"].(byte)
		if !normal {
			return types.CommandBlockData{}, fmt.Errorf("CommandBlockParseData: Crashed in Cb[\"ExecuteOnFirstTick\"]; Cb = %#v", Cb)
		}
		if got == byte(0) {
			executeOnFirstTick = false
		} else {
			executeOnFirstTick = true
		}
	}
	// ExecuteOnFirstTick
	_, ok = Cb["TrackOutput"]
	if ok {
		got, normal := Cb["TrackOutput"].(byte)
		if !normal {
			return types.CommandBlockData{}, fmt.Errorf("CommandBlockParseData: Crashed in Cb[\"TrackOutput\"]; Cb = %#v", Cb)
		}
		if got == byte(0) {
			trackOutput = false
		} else {
			trackOutput = true
		}
	}
	// TrackOutput
	_, ok = Cb["conditionalMode"]
	if ok {
		got, normal := Cb["conditionalMode"].(byte)
		if !normal {
			return types.CommandBlockData{}, fmt.Errorf("CommandBlockParseData: Crashed in Cb[\"conditionalMode\"]; Cb = %#v", Cb)
		}
		if got == byte(0) {
			conditionalMode = false
		} else {
			conditionalMode = true
		}
	}
	// conditionalMode
	_, ok = Cb["auto"]
	if ok {
		got, normal := Cb["auto"].(byte)
		if !normal {
			return types.CommandBlockData{}, fmt.Errorf("CommandBlockParseData: Crashed in Cb[\"auto\"]; Cb = %#v", Cb)
		}
		if got == byte(0) {
			needRedstone = true
		} else {
			needRedstone = false
		}
	}
	// auto
	return types.CommandBlockData{
		Mode:               uint32(mode),
		Command:            command,
		CustomName:         customName,
		LastOutput:         lastOutput,
		TickDelay:          tickDelay,
		ExecuteOnFirstTick: executeOnFirstTick,
		TrackOutput:        trackOutput,
		Conditional:        conditionalMode,
		NeedRedstone:       needRedstone,
	}, nil
}

func placeCommandBlock(
	Environment *environment.PBEnvironment,
	Mainsettings *types.MainConfig,
	IsFastMode bool,
	BlockInfo *types.Module,
) error {
	if !Mainsettings.ExcludeCommands && BlockInfo.CommandBlockData != nil {
		cmdsender := Environment.CommandSender.(*commands.CommandSender)
		if BlockInfo.Block != nil {
			request := commands_generator.SetBlockRequest(BlockInfo, Mainsettings)
			//<-time.After(time.Second)
			wc := make(chan bool)
			(*cmdsender.GetBlockUpdateSubscribeMap()).Store(protocol.BlockPos{int32(BlockInfo.Point.X), int32(BlockInfo.Point.Y), int32(BlockInfo.Point.Z)}, wc)
			err := cmdsender.SendDimensionalCommand(request)
			select {
			case <-wc:
				break
			case <-time.After(time.Second * 2):
				(*cmdsender.GetBlockUpdateSubscribeMap()).Delete(protocol.BlockPos{int32(BlockInfo.Point.X), int32(BlockInfo.Point.Y), int32(BlockInfo.Point.Z)})
			}
			close(wc)
			if err != nil {
				return fmt.Errorf("placeCommandBlock: %v", err)
			}
		}
		Cbdata := BlockInfo.CommandBlockData
		if Mainsettings.InvalidateCommands {
			Cbdata.Command = "|" + Cbdata.Command
		}
		if !IsFastMode {
			UUID := uuid.New()
			w := make(chan *packet.CommandOutput)
			(*cmdsender.GetUUIDMap()).Store(UUID.String(), w)
			err := cmdsender.SendWSCommand(fmt.Sprintf("tp %d %d %d", BlockInfo.Point.X, BlockInfo.Point.Y+1, BlockInfo.Point.Z), UUID)
			select {
			case <-time.After(time.Second):
				(*cmdsender.GetUUIDMap()).Delete(UUID.String())
				break
			case <-w:
			}
			close(w)
			if err != nil {
				return fmt.Errorf("placeCommandBlock: %v", err)
			}
		}
		cmdsender.UpdateCommandBlock(int32(BlockInfo.Point.X), int32(BlockInfo.Point.Y), int32(BlockInfo.Point.Z), Cbdata)
	}
	return nil
}
