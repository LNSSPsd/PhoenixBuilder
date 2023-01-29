package blockNBT_depends

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

type EnchList []EnchSingle

type EnchSingle struct {
	Id    int16
	Level int16
}

func CheskEnchItem(input map[string]interface{}) bool {
	_, ok := input["ench"]
	return ok
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

func enchantRequest(Environment *environment.PBEnvironment, input []interface{}) error {
	got, err := parseEnchList(input)
	if err != nil {
		return fmt.Errorf("enchantRequest: %v", err)
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
			return fmt.Errorf("enchantRequest: %v", err)
		}
	}
	_, err = sender.SendWSCommandWithResponce(fmt.Sprintf("enchant @s %v %v", got[len(got)-1].Id, got[len(got)-1].Level))
	if err != nil {
		return fmt.Errorf("enchantRequest: %v", err)
	}
	return nil
}

// 物品始终会生成于 slot.hotbar 0
func ReplaceitemAndEnchant(
	Environment *environment.PBEnvironment,
	ItemInfo *types.ChestSlot,
) error {
	cmdsender := Environment.CommandSender.(*commands.CommandSender)
	_, err := cmdsender.SendWSCommandWithResponce("replaceitem entity @s slot.hotbar 0 air")
	if err != nil {
		return fmt.Errorf("ReplaceitemAndEnchant: %v", err)
	}
	Environment.Connection.(*minecraft.Conn).WritePacket(&packet.MobEquipment{
		EntityRuntimeID: Environment.Connection.(*minecraft.Conn).GameData().EntityRuntimeID,
		NewItem: protocol.ItemInstance{
			StackNetworkID: 0,
			Stack: protocol.ItemStack{
				ItemType: protocol.ItemType{
					NetworkID:     0,
					MetadataValue: 0,
				},
				BlockRuntimeID: 0,
				Count:          0,
				NBTData:        map[string]interface{}{},
				CanBePlacedOn:  []string(nil),
				CanBreak:       []string(nil),
				HasNetworkID:   false,
			},
		},
		InventorySlot: 0,
		HotBarSlot:    0,
		WindowID:      0,
	})
	// change the slot.weapon.mainhand into slot.hotbar 0
	_, err = cmdsender.SendWSCommandWithResponce(fmt.Sprintf("replaceitem entity @s slot.weapon.mainhand 0 %v %v %v %v", ItemInfo.Name, ItemInfo.Count, ItemInfo.Damage, commands_generator.GetReplaceItemEnhancement(
		&types.Module{
			ChestData: &types.ChestData{
				types.ChestSlot{
					CanPlaceOn: ItemInfo.CanPlaceOn,
					CanDestroy: ItemInfo.CanDestroy,
					ItemNBT:    ItemInfo.ItemNBT,
				},
			},
		}, 0)))
	if err != nil {
		return fmt.Errorf("ReplaceitemAndEnchant: %v", err)
	}
	// replaceitem
	_, isEnchItem := ItemInfo.ItemNBT["ench"]
	if isEnchItem {
		got, normal := ItemInfo.ItemNBT["ench"].([]interface{})
		if normal {
			err = enchantRequest(Environment, got)
			if err != nil {
				return fmt.Errorf("ReplaceitemAndEnchant: %v", err)
			}
		}
	}
	// enchant item stack
	if CheckVersion() {
		networkID, ok := ItemRunTimeID[ItemInfo.Name]
		if ok {
			if isEnchItem {
				Environment.Connection.(*minecraft.Conn).WritePacket(&packet.MobEquipment{
					EntityRuntimeID: Environment.Connection.(*minecraft.Conn).GameData().EntityRuntimeID,
					NewItem: protocol.ItemInstance{
						StackNetworkID: 0,
						Stack: protocol.ItemStack{
							ItemType: protocol.ItemType{
								NetworkID:     int32(networkID),
								MetadataValue: uint32(ItemInfo.Damage),
							},
							BlockRuntimeID: 0,
							Count:          1,
							NBTData: map[string]interface{}{
								"ench": []interface{}{},
							},
							CanBePlacedOn: ItemInfo.CanPlaceOn,
							CanBreak:      ItemInfo.CanDestroy,
							HasNetworkID:  false,
						},
					},
					InventorySlot: 0,
					HotBarSlot:    0,
					WindowID:      0,
				})
			} else {
				Environment.Connection.(*minecraft.Conn).WritePacket(&packet.MobEquipment{
					EntityRuntimeID: Environment.Connection.(*minecraft.Conn).GameData().EntityRuntimeID,
					NewItem: protocol.ItemInstance{
						StackNetworkID: 0,
						Stack: protocol.ItemStack{
							ItemType: protocol.ItemType{
								NetworkID:     int32(networkID),
								MetadataValue: uint32(ItemInfo.Damage),
							},
							BlockRuntimeID: 0,
							Count:          1,
							CanBePlacedOn:  ItemInfo.CanPlaceOn,
							CanBreak:       ItemInfo.CanDestroy,
							HasNetworkID:   false,
						},
					},
					InventorySlot: 0,
					HotBarSlot:    0,
					WindowID:      0,
				})
			}
		}
	}
	// let other players know what happened
	return nil
}
