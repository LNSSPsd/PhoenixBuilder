package blockNBT

import (
	"fmt"
	"phoenixbuilder/fastbuilder/commands_generator"
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/fastbuilder/types"
	"phoenixbuilder/io/commands"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
)

type SignData struct {
	TextOwner                   string
	IgnoreLighting              byte
	SignTextColor               int32
	TextIgnoreLegacyBugResolved byte
	Text                        string
}

type SignInput struct {
	Environment  *environment.PBEnvironment
	Mainsettings *types.MainConfig
	BlockInfo    *types.Module
	Sign         *map[string]interface{}
}

func Sign(input *SignInput) error {
	err := placeSign(input.Environment, input.Mainsettings, input.BlockInfo, *input.Sign)
	if err != nil {
		return fmt.Errorf("Sign: Failed to place the entity block named %v at (%v,%v,%v), and the error log is %v", *input.BlockInfo.Block.Name, input.BlockInfo.Point.X, input.BlockInfo.Point.Y, input.BlockInfo.Point.Z, err)
	}
	return nil
}

func parseSignData(Sign *map[string]interface{}) (*SignData, error) {
	var ok bool = false
	var normal bool = false
	var textOwner string = ""
	var ignoreLighting byte = byte(0)
	var signTextColor int32 = 0
	var textIgnoreLegacyBugResolved byte = byte(0)
	var text string = ""
	// prepare
	SIGN := *Sign
	// prepare
	_, ok = SIGN["TextOwner"]
	if !ok {
		return &SignData{}, fmt.Errorf("parseSignData: Could not find Sign[\"TextOwner\"]; Sign = %#v", SIGN)
	}
	textOwner, normal = SIGN["TextOwner"].(string)
	if !normal {
		return &SignData{}, fmt.Errorf("parseSignData: Could not parse Sign[\"TextOwner\"]; Sign = %#v", SIGN)
	}
	// TextOwner
	_, ok = SIGN["IgnoreLighting"]
	if !ok {
		return &SignData{}, fmt.Errorf("parseSignData: Could not find Sign[\"IgnoreLighting\"]; Sign = %#v", SIGN)
	}
	ignoreLighting, normal = SIGN["IgnoreLighting"].(byte)
	if !normal {
		return &SignData{}, fmt.Errorf("parseSignData: Could not parse Sign[\"IgnoreLighting\"]; Sign = %#v", SIGN)
	}
	// IgnoreLighting
	_, ok = SIGN["SignTextColor"]
	if !ok {
		return &SignData{}, fmt.Errorf("parseSignData: Could not find Sign[\"SignTextColor\"]; Sign = %#v", SIGN)
	}
	signTextColor, normal = SIGN["SignTextColor"].(int32)
	if !normal {
		return &SignData{}, fmt.Errorf("parseSignData: Could not parse Sign[\"SignTextColor\"]; Sign = %#v", SIGN)
	}
	// SignTextColor
	_, ok = SIGN["TextIgnoreLegacyBugResolved"]
	if !ok {
		return &SignData{}, fmt.Errorf("parseSignData: Could not find Sign[\"TextIgnoreLegacyBugResolved\"]; Sign = %#v", SIGN)
	}
	textIgnoreLegacyBugResolved, normal = SIGN["TextIgnoreLegacyBugResolved"].(byte)
	if !normal {
		return &SignData{}, fmt.Errorf("parseSignData: Could not parse Sign[\"TextIgnoreLegacyBugResolved\"]; Sign = %#v", SIGN)
	}
	// TextIgnoreLegacyBugResolved
	_, ok = SIGN["Text"]
	if !ok {
		return &SignData{}, fmt.Errorf("parseSignData: Could not find Sign[\"Text\"]; Sign = %#v", SIGN)
	}
	text, normal = SIGN["Text"].(string)
	if !normal {
		return &SignData{}, fmt.Errorf("parseSignData: Could not parse Sign[\"Text\"]; Sign = %#v", SIGN)
	}
	// Text
	return &SignData{
		TextOwner:                   textOwner,
		IgnoreLighting:              ignoreLighting,
		SignTextColor:               signTextColor,
		TextIgnoreLegacyBugResolved: textIgnoreLegacyBugResolved,
		Text:                        text,
	}, nil
}

func placeSign(
	Environment *environment.PBEnvironment,
	Mainsettings *types.MainConfig,
	BlockInfo *types.Module,
	Sign map[string]interface{},
) error {
	SignData, err := parseSignData(&Sign)
	if err != nil {
		return fmt.Errorf("placeSign: %v", err)
	}
	// parse sign data
	cmdsender := Environment.CommandSender.(*commands.CommandSender)
	request := commands_generator.SetBlockRequest(BlockInfo, Mainsettings)
	_, err = cmdsender.SendWSCommandWithResponce(request)
	if err != nil {
		return fmt.Errorf("placeSign: %v", err)
	}
	// place sign block
	Environment.Connection.(*minecraft.Conn).WritePacket(&packet.BlockActorData{
		Position: protocol.BlockPos{int32(BlockInfo.Point.X), int32(BlockInfo.Point.Y), int32(BlockInfo.Point.Z)},
		NBTData: map[string]interface{}{
			"TextOwner":                   SignData.TextOwner,
			"IgnoreLighting":              SignData.IgnoreLighting,
			"SignTextColor":               SignData.SignTextColor,
			"TextIgnoreLegacyBugResolved": SignData.TextIgnoreLegacyBugResolved,
			"Text":                        SignData.Text,
		},
	})
	// write text
	return nil
}
