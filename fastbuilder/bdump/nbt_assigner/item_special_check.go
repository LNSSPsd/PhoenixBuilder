package NBTAssigner

// 检查名为 itemName 且数据值为 itemMetaData 的物品的 Tag 标签。
// 如果发现该 NBT 是默认状态下的 NBT ，
// 则认为该物品可以直接可以通过 replaceitem 生成，
// 而无需进行特殊处理，此时返回假。
// 否则，认为该物品需要进行特殊处理，
// 例如烟花这种物品需要生成工作台并合成，因此此时返回真
func ItemSpecialCheck(
	itemName string,
	itemType string,
	itemMetaData uint16,
	itemTag map[string]interface{},
) bool {
	switch itemType {
	case "fireworks":
	}
	// check
	return false
	// return
}
