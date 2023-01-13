package blockNBT_depends

import (
	"fmt"
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/io/commands"
)

type EnchList []EnchSingle

type EnchSingle struct {
	Id    int16
	Level int16
}

func parseEnchList(Ench []interface{}) (EnchList, error) {
	ans := make([]EnchSingle, 0)
	for key, value := range Ench {
		single, normal := value.(map[string]interface{})
		if !normal {
			return EnchList{}, fmt.Errorf("parseEnchList: Could not parse ench[%v]; ench = %#v", key, Ench)
		}
		_, ok := single["id"]
		if !ok {
			return EnchList{}, fmt.Errorf("parseEnchList: Could not find ench[%v][\"id\"]; ench = %#v", key, Ench)
		}
		id, normal := single["id"].(int16)
		if !normal {
			return EnchList{}, fmt.Errorf("parseEnchList: Could not parse ench[%v][\"id\"]; ench = %#v", key, Ench)
		}
		_, ok = single["lvl"]
		if !ok {
			return EnchList{}, fmt.Errorf("parseEnchList: Could not find ench[%v][\"lvl\"]; ench = %#v", key, Ench)
		}
		lvl, normal := single["lvl"].(int16)
		if !normal {
			return EnchList{}, fmt.Errorf("parseEnchList: Could not parse ench[%v][\"lvl\"]; ench = %#v", key, Ench)
		}
		ans = append(ans, EnchSingle{
			Id:    id,
			Level: lvl,
		})
	}
	return ans, nil
}

func EnchantRequest(Environment *environment.PBEnvironment, input []interface{}) error {
	got, err := parseEnchList(input)
	if err != nil {
		return fmt.Errorf("EnchantRequest: %v", err)
	}
	if len(got) <= 0 {
		return nil
	}
	sender := Environment.CommandSender.(*commands.CommandSender)
	for key, value := range got {
		if key == len(got)-1 {
			break
		}
		err := sender.SendDimensionalCommand(fmt.Sprintf("enchant @s %v %v", value.Id, value.Level))
		if err != nil {
			return fmt.Errorf("EnchantRequest: %v", err)
		}
	}
	_, err = sender.SendWSCommandWithResponce(fmt.Sprintf("enchant @s %v %v", got[len(got)-1].Id, got[len(got)-1].Level))
	if err != nil {
		return fmt.Errorf("EnchantRequest: %v", err)
	}
	return nil
}
