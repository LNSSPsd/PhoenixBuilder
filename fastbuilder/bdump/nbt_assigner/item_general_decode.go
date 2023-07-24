package NBTAssigner

import (
	"encoding/json"
	"fmt"
	"phoenixbuilder/mirror/chunk"
	"strings"
)

// 从 singleItem 解码单个物品的基本数据
func DecodeItemBasicData(singleItem ItemOrigin) (ItemBasicData, error) {
	var count uint8
	var itemData uint16
	var name string
	var slot uint8
	// 初始化
	{
		count_origin, ok := singleItem["Count"]
		if !ok {
			return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: singleItem["Count"] does not exist; singleItem = %#v`, singleItem)
		}
		count_got, normal := count_origin.(byte)
		if !normal {
			return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: Can not convert count_origin into byte(uint8); singleItem = %#v`, singleItem)
		}
		count = count_got
	}
	// 物品数量
	{
		name_origin, ok := singleItem["Name"]
		if !ok {
			return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: singleItem["Name"] does not exist; singleItem = %#v`, singleItem)
		}
		name_got, normal := name_origin.(string)
		if !normal {
			return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: Can not convert name_origin into string; singleItem = %#v`, singleItem)
		}
		name = strings.Replace(strings.ToLower(name_got), "minecraft:", "", 1)
	}
	// 物品的英文 ID (已去除命名空间)
	if slot_origin, ok := singleItem["Slot"]; ok {
		slot_got, normal := slot_origin.(byte)
		if !normal {
			return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: Can not convert slot_origin into byte(uint8); singleItem = %#v`, singleItem)
		}
		slot = slot_got
	}
	// 物品所在的槽位(对于唱片机等单槽位方块来说，此数据不存在)
	{
		{
			damage_origin, ok := singleItem["Damage"]
			if !ok {
				return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: singleItem["Damage"] does not exist; singleItem = %#v`, singleItem)
			}
			damage_got, normal := damage_origin.(int16)
			if !normal {
				return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: Can not convert damage_origin into int16; singleItem = %#v`, singleItem)
			}
			itemData = uint16(damage_got)
		}
		// Damage
		for i := 0; i < 1; i++ {
			tag_origin, ok := singleItem["tag"]
			if !ok {
				break
			}
			tag_got, normal := tag_origin.(map[string]interface{})
			if !normal {
				return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: Can not convert tag_origin into map[string]interface{}; singleItem = %#v`, singleItem)
			}
			damage_origin, ok := tag_got["Damage"]
			if !ok {
				break
			}
			damage_got, normal := damage_origin.(int32)
			if !normal {
				return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: Can not convert damage_origin into int32; singleItem = %#v`, singleItem)
			}
			itemData = uint16(damage_got)
		}
		// tag["Damage"]
		for i := 0; i < 1; i++ {
			block_origin, ok := singleItem["Block"]
			if !ok {
				break
			}
			block_got, normal := block_origin.(map[string]interface{})
			if !normal {
				return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: Can not convert block_got into map[string]interface{}; singleItem = %#v`, singleItem)
			}
			val_origin, ok := block_got["val"]
			if ok {
				val_got, normal := val_origin.(int16)
				if !normal {
					return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: Can not convert val_origin into int16; singleItem = %#v`, singleItem)
				}
				itemData = uint16(val_got)
			} else {
				states_origin, ok := block_got["states"]
				if !ok {
					break
				}
				states_got, normal := states_origin.(map[string]interface{})
				if !normal {
					return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: Can not convert states_origin into map[string]interface{}; singleItem = %#v`, singleItem)
				}
				runtimeId, found := chunk.StateToRuntimeID(name, states_got)
				if !found {
					return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: Could not convert legacy block to standard runtime id; singleItem = %#v`, singleItem)
				}
				legacyBlock, found := chunk.RuntimeIDToLegacyBlock(runtimeId)
				if !found {
					return ItemBasicData{}, fmt.Errorf(`DecodeItemBasicData: Could not convert standard runtime id to block states; singleItem = %#v`, singleItem)
				}
				itemData = legacyBlock.Val
			}
		}
		// Block["val"] or Block["states"]
	}
	/*
		物品数据值(附加值)

		以上三个方法都在拿物品数据值(附加值)，而数据值的获取优先级如下。
		Damage < tag["Damage"] < Block["val"]

		顶层复合标签下的 Damage 数据一定存在，但不一定代表物品真实的物品数据值。

		如果当前物品是武器或者工具，其 tag["Damage"] 处会说明其耐久值，
		而这个耐久值才是真正的物品数据值；
		如果当前物品是一个方块，其 block["val"] 处可能会声明其方块数据值，
		而这个方块数据值才是真正的物品数据值。
		当然，如果这个 BDX 文件是从国际版制作的，那么 block["val"] 处将不存在数据，
		此时将需要从 Block["states"] 处得到此物品所对应方块的方块数据值。

		NOTE: 不保证目前已经提供的这三个方法涵盖了所有情况，一切还需要进一步的研究
	*/
	return ItemBasicData{
		Name:     name,
		Count:    count,
		MetaData: itemData,
		Slot:     slot,
	}, nil
	// 返回值
}

// 从 singleItem 解码单个物品的增强数据，
// 其中包含物品组件、显示名称和附魔属性。
// 特别地，如果此物品存在 item_lock 物品组件，
// 则只会解析物品组件的相关数据，
// 因为存在 item_lock 的物品无法跨容器移动
func DecodeItemEnhancementData(
	singleItem ItemOrigin,
) (*ItemEnhancementData, error) {
	var displayName string
	var enchantments *[]Enchantment
	var itemComponents *ItemComponents
	var nbt_tag_got map[string]interface{}
	var normal bool
	// 初始化
	nbt_tag_origin, ok := singleItem["tag"]
	if ok {
		nbt_tag_got, normal = nbt_tag_origin.(map[string]interface{})
		if !normal {
			return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert nbt_tag_origin into map[string]interface{}; singleItem = %#v`, singleItem)
		}
	}
	// 获取当前物品的 tag 数据
	{
		if can_place_on_origin, ok := singleItem["CanPlaceOn"]; ok {
			can_place_on_got, normal := can_place_on_origin.([]interface{})
			if !normal {
				return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert can_place_on_origin into []interface{}; singleItem = %#v`, singleItem)
			}
			if itemComponents == nil {
				itemComponents = &ItemComponents{}
			}
			for key, value := range can_place_on_got {
				blockName, normal := value.(string)
				if !normal {
					return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert can_place_on_got[%d] into string; singleItem = %#v`, key, singleItem)
				}
				itemComponents.CanPlaceOn = append(itemComponents.CanPlaceOn, blockName)
			}
		}
		// can_place_on
		if can_destroy_origin, ok := singleItem["CanDestroy"]; ok {
			can_destroy_got, normal := can_destroy_origin.([]interface{})
			if !normal {
				return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert can_destroy_origin into []interface{}; singleItem = %#v`, singleItem)
			}
			if itemComponents == nil {
				itemComponents = &ItemComponents{}
			}
			for key, value := range can_destroy_got {
				blockName, normal := value.(string)
				if !normal {
					return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert can_destroy_got[%d] into string; singleItem = %#v`, key, singleItem)
				}
				itemComponents.CanDestroy = append(itemComponents.CanDestroy, blockName)
			}
		}
		// can_destroy
		if nbt_tag_got != nil {
			if item_lock_origin, ok := nbt_tag_got["minecraft:item_lock"]; ok {
				item_lock_got, normal := item_lock_origin.(byte)
				if !normal {
					return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert item_lock_origin into byte(uint8); singleItem = %#v`, singleItem)
				}
				if itemComponents == nil {
					itemComponents = &ItemComponents{}
				}
				switch item_lock_got {
				case 1:
					itemComponents.ItemLock = "lock_in_slot"
				case 2:
					itemComponents.ItemLock = "lock_in_inventory"
				default:
					return nil, fmt.Errorf(`DecodeItemAdditionalData: Unknown value(%d) of item_lock; singleItem = %#v`, item_lock_got, singleItem)
				}
			}
			// item_lock
			if keep_on_death_origin, ok := nbt_tag_got["minecraft:keep_on_death"]; ok {
				keep_on_death_got, normal := keep_on_death_origin.(byte)
				if !normal {
					return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert keep_on_death_origin into byte(uint8); singleItem = %#v`, singleItem)
				}
				if keep_on_death_got != 1 {
					return nil, fmt.Errorf(`DecodeItemAdditionalData: Unknown value(%d) of kepp_on_death; singleItem = %#v`, keep_on_death_got, singleItem)
				}
				if itemComponents == nil {
					itemComponents = &ItemComponents{}
				}
				itemComponents.KeepOnDeath = true
			}
			// keep_on_death
		}
		// item_lock and keep_on_death
	}
	// 物品组件
	if itemComponents != nil && len(itemComponents.ItemLock) != 0 {
		return &ItemEnhancementData{
			DisplayName:    "",
			Enchantments:   nil,
			ItemComponents: itemComponents,
		}, nil
	}
	// 如果当前物品已经使用了 item_lock 物品组件，
	// 则无需再解析后续的数据，
	// 因为存在 item_lock 的物品无法跨容器移动
	for i := 0; i < 1; i++ {
		display_origin, ok := nbt_tag_got["display"]
		if !ok {
			break
		}
		display_got, normal := display_origin.(map[string]interface{})
		if !normal {
			return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert display_origin into map[string]interface{}; singleItem = %#v`, singleItem)
		}
		name_origin, ok := display_got["Name"]
		if !ok {
			break
		}
		name_got, normal := name_origin.(string)
		if !normal {
			return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert name_origin into string; singleItem = %#v`, singleItem)
		}
		displayName = name_got
	}
	// 物品的显示名称
	if ench_origin, ok := nbt_tag_got["ench"]; ok {
		ench_got, normal := ench_origin.([]interface{})
		if !normal {
			return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert ench_origin into []interface{}; singleItem = %#v`, singleItem)
		}
		if len(ench_got) > 0 {
			enchantments = &[]Enchantment{}
		}
		for key, value := range ench_got {
			value_got, normal := value.(map[string]interface{})
			if !normal {
				return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert ench_got[%d] into map[string]interface{}; singleItem = %#v`, key, singleItem)
			}
			id_origin, ok := value_got["id"]
			if !ok {
				return nil, fmt.Errorf(`DecodeItemAdditionalData: ench_got[%d]["id"] does not exist; singleItem = %#v`, key, singleItem)
			}
			id_got, normal := id_origin.(int16)
			if !normal {
				return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert ench_got[%d]["id"] into int16; singleItem = %#v`, key, singleItem)
			}
			lvl_origin, ok := value_got["lvl"]
			if !ok {
				return nil, fmt.Errorf(`DecodeItemAdditionalData: ench_got[%d]["lvl"] does not exist; singleItem = %#v`, key, singleItem)
			}
			lvl_got, normal := lvl_origin.(int16)
			if !normal {
				return nil, fmt.Errorf(`DecodeItemAdditionalData: Can not convert ench_got[%d]["lvl"] into int16; singleItem = %#v`, key, singleItem)
			}
			*enchantments = append(*enchantments, Enchantment{ID: uint8(id_got), Level: lvl_got})
		}
	}
	// 物品的附魔属性
	if len(displayName) != 0 || enchantments != nil || itemComponents != nil {
		return &ItemEnhancementData{
			DisplayName:    displayName,
			Enchantments:   enchantments,
			ItemComponents: itemComponents,
		}, nil
	}
	return nil, nil
	// 返回值
}

// 从 singleItem 解码单个物品的自定义 NBT 数据
func DecodeItemCustomData(
	itemBasicData ItemBasicData,
	singleItem ItemOrigin,
) (*ItemCustomData, error) {
	nbt_tag_origin, ok := singleItem["tag"]
	if !ok {
		return nil, nil
	}
	nbt_tag_got, normal := nbt_tag_origin.(map[string]interface{})
	if !normal {
		return nil, fmt.Errorf(`DecodeItemCustomData: Can not convert nbt_tag_origin into map[string]interface{}; singleItem = %#v`, singleItem)
	}
	// 获取当前物品的 tag 数据
	/*
		此代码块还在 WIP 阶段
		blockType,ok := supportBlocksPool[itemBasicData.Name]
		if ok {
			switch blockType {
				case
			}
		}
	*/
	json.Marshal(nbt_tag_got)
	return nil, nil
}