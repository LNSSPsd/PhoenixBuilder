package NBTAssigner

import (
	"encoding/json"
	"fmt"
	"phoenixbuilder/fastbuilder/mcstructure"
	"phoenixbuilder/fastbuilder/types"
	"strings"
)

// 从 SupportBlocksPool 检查这个方块实体是否已被支持。
// 如果尚未被支持，则返回空字符串，否则返回这种方块的类型。
// 以告示牌为例，所有的告示牌都可以写作为 Sign
func IsNBTBlockSupported(blockName string) string {
	value, ok := SupportBlocksPool[blockName]
	if ok {
		return value
	}
	return ""
}

// 从 SupportItemsPool 检查这个 NBT 物品是否已被支持。
// 如果尚未被支持，则返回空字符串，否则返回这种物品的类型。
// 以告示牌为例，所有的告示牌都可以写作为 Sign
func IsNBTItemSupported(itemName string) string {
	value, ok := SupportItemsPool[itemName]
	if ok {
		return value
	}
	return ""
}

// 将 itemComponents 编码为游戏支持的 JSON 格式。
// 如果传入的 itemComponents 为空指针，则返回空字符串
func MarshalItemComponents(itemComponents *ItemComponents) string {
	type can_place_on_or_can_destroy struct {
		Blocks []string `json:"blocks"`
	}
	type item_lock struct {
		Mode string `json:"mode"`
	}
	res := map[string]interface{}{}
	// 初始化
	if itemComponents == nil {
		return ""
	}
	// 如果物品组件不存在，那么应该返回空字符串而非 {}
	if len(itemComponents.CanPlaceOn) > 0 {
		res["can_place_on"] = can_place_on_or_can_destroy{Blocks: itemComponents.CanPlaceOn}
	}
	if len(itemComponents.CanDestroy) > 0 {
		res["can_destroy"] = can_place_on_or_can_destroy{Blocks: itemComponents.CanDestroy}
	}
	if itemComponents.KeepOnDeath {
		res["keep_on_death"] = struct{}{}
	}
	if len(itemComponents.ItemLock) != 0 {
		res["item_lock"] = item_lock{Mode: itemComponents.ItemLock}
	}
	// 赋值
	bytes, _ := json.Marshal(res)
	return string(bytes)
	// 返回值
}

// 将 types.Module 解析为 GeneralBlock
func ParseBlockModule(singleBlock *types.Module) (GeneralBlock, error) {
	got, err := mcstructure.ParseStringNBT(singleBlock.Block.BlockStates, true)
	if err != nil {
		return GeneralBlock{}, fmt.Errorf("ParseBlockModule: Could not parse block states; singleBlock.Block.BlockStates = %#v", singleBlock.Block.BlockStates)
	}
	blockStates, normal := got.(map[string]interface{})
	if !normal {
		return GeneralBlock{}, fmt.Errorf("ParseBlockModule: The target block states is not map[string]interface{}; got = %#v", got)
	}
	// get block states
	return GeneralBlock{
		Name:   strings.Replace(strings.ToLower(strings.ReplaceAll(*singleBlock.Block.Name, " ", "")), "minecraft:", "", 1),
		States: blockStates,
		NBT:    singleBlock.NBTMap,
	}, nil
	// return
}

/*
将 singleItem 解析为 GeneralItem 。

特别地，如果此物品存在 item_lock 物品组件，
则只会解析物品组件的相关数据，
因为存在 item_lock 的物品无法跨容器移动；

如果此物品是一个 NBT 方块，
则附魔属性将被丢弃，因为无法为方块附魔
*/
func ParseItemFromNBT(
	singleItem ItemOrigin,
	supportBlocksPool map[string]string,
) (GeneralItem, error) {
	itemBasicData, err := DecodeItemBasicData(singleItem)
	if err != nil {
		return GeneralItem{}, fmt.Errorf("ParseItemFromNBT: %v", err)
	}
	// basic
	itemAdditionalData, err := DecodeItemEnhancementData(singleItem)
	if err != nil {
		return GeneralItem{}, fmt.Errorf("ParseItemFromNBT: %v", err)
	}
	// additional
	if itemAdditionalData != nil && itemAdditionalData.ItemComponents != nil && len(itemAdditionalData.ItemComponents.ItemLock) != 0 {
		return GeneralItem{
			Basic:       itemBasicData,
			Enhancement: itemAdditionalData,
			Custom:      nil,
		}, nil
	}
	// 如果此物品使用了物品组件 item_lock ，
	// 则后续数据将不被解析。
	// 因为存在 item_lock 的物品无法跨容器移动
	itemCustomData, err := DecodeItemCustomData(itemBasicData, singleItem)
	if err != nil {
		return GeneralItem{}, fmt.Errorf("ParseItemFromNBT: %v", err)
	}
	// custom
	if itemCustomData != nil && itemCustomData.SubBlockData != nil && itemAdditionalData != nil {
		itemAdditionalData.Enchantments = nil
	}
	// 如果此物品是一个 NBT 方块，
	// 则附魔属性将被丢弃，因为无法为方块附魔
	return GeneralItem{
		Basic:       itemBasicData,
		Enhancement: itemAdditionalData,
		Custom:      itemCustomData,
	}, nil
	// return
}
