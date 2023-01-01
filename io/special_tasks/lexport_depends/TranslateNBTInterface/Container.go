package TranslateNBTInerface

import (
	"fmt"
	"phoenixbuilder/fastbuilder/types"
)

// 检查一个方块是否是有效的容器；这里的有效指的是可以被 replaceitem 命令生效的容器
func CheckIfIsEffectiveContainer(name string) (string, error) {
	index := map[string]string{
		"blast_furnace":      "Items",
		"lit_blast_furnace":  "Items",
		"smoker":             "Items",
		"lit_smoker":         "Items",
		"furnace":            "Items",
		"lit_furnace":        "Items",
		"chest":              "Items",
		"barrel":             "Items",
		"trapped_chest":      "Items",
		"lectern":            "book",
		"hopper":             "Items",
		"dispenser":          "Items",
		"dropper":            "Items",
		"cauldron":           "Items",
		"lava_cauldron":      "Items",
		"jukebox":            "RecordItem",
		"brewing_stand":      "Items",
		"undyed_shulker_box": "Items",
		"shulker_box":        "Items",
	}
	value, ok := index[name]
	if ok {
		return value, nil
	} else {
		return "", fmt.Errorf("CheckIfIsEffectiveContainer: \"%v\" not found", name)
	}
}

// 将 Interface NBT 转换为 types.ChestData
func GetContainerData(container interface{}) (types.ChestData, error) {
	var correct []interface{} = make([]interface{}, 0)
	// 初始化
	got, normal := container.([]interface{})
	if !normal {
		got, normal := container.(map[string]interface{})
		if normal {
			correct = append(correct, got)
		} else {
			return types.ChestData{}, fmt.Errorf("Crashed in input")
		}
	} else {
		correct = got
	}
	// 把物品丢入 correct 里面
	ans := make(types.ChestData, 0)
	for key, value := range correct {
		var count uint8 = uint8(0)
		var itemData uint16 = uint16(0)
		var name string = ""
		var slot uint8 = uint8(0)
		// 初始化
		containerData, normal := value.(map[string]interface{})
		if normal {
			_, ok := containerData["Count"]
			if ok {
				got, normal := containerData["Count"].(byte)
				if normal {
					count = uint8(got)
				} else {
					return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Count\"]", key)
				}
			} else {
				return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Count\"]", key)
			}
			// 物品数量
			_, ok = containerData["Damage"]
			if ok {
				got, normal := containerData["Damage"].(int16)
				if normal {
					itemData = uint16(got)
				} else {
					return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Damage\"]", key)
				}
			} else {
				return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Damage\"]", key)
			}
			_, ok = containerData["tag"]
			if ok {
				tag, normal := containerData["tag"].(map[string]interface{})
				if normal {
					_, ok = tag["Damage"]
					if ok {
						got, normal := tag["Damage"].(int32)
						if normal {
							itemData = uint16(got)
						} else {
							return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"tag\"]", key)
						}
					}
				} else {
					return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"tag\"]", key)
				}
			}
			_, ok = containerData["Block"]
			if ok {
				Block, normal := containerData["Block"].(map[string]interface{})
				if normal {
					_, ok = Block["Block"]
					if ok {
						got, normal := Block["Block"].(map[string]interface{})
						if normal {
							_, ok = got["states"]
							if ok {
								states, normal := got["states"].(map[string]interface{})
								if normal {
									_, ok = states["val"]
									if ok {
										got, normal := states["val"].(int16)
										if normal {
											itemData = uint16(got)
										} else {
											return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Block\"][\"states\"][\"val\"]", key)
										}
									} else {
										return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Block\"][\"states\"][\"val\"]", key)
									}
								} else {
									return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Block\"][\"states\"]", key)
								}
							} else {
								return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Block\"][\"states\"]", key)
							}
						} else {
							return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Block\"]", key)
						}
					}
				} else {
					return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Block\"]", key)
				}
			}
			// 物品数据值(附加值)
			// 需要说明的是，数据值的获取优先级是这样的
			// Damage < tag["Damage"] < Block["states"]["val"]
			// 我目前还没有找到一种妥善的办法可以解决数据值的问题，所以我希望能有人帮帮我！
			_, ok = containerData["Name"]
			if ok {
				got, normal := containerData["Name"].(string)
				if normal {
					name = got
				} else {
					return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Name\"]", key)
				}
			} else {
				return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Name\"]", key)
			}
			// 物品名称
			_, ok = containerData["Slot"]
			if ok {
				got, normal := containerData["Slot"].(byte)
				if normal {
					slot = uint8(got)
				} else {
					return types.ChestData{}, fmt.Errorf("Crashed in input[%v][\"Slot\"]", key)
				}
			}
			// 物品所在的栏位
			ans = append(ans, types.ChestSlot{
				Name:   name,
				Count:  count,
				Damage: itemData,
				Slot:   slot,
			})
		} else {
			return types.ChestData{}, fmt.Errorf("Crashed in input[%v]", key)
		}
	}
	return ans, nil
}

// 主函数
func GetContainerDataRun(blockNBT map[string]interface{}, blockName string) (types.ChestData, error) {
	key, err := CheckIfIsEffectiveContainer(blockName)
	if err != nil {
		return types.ChestData{}, fmt.Errorf("GetContainerDataRun: Not a container")
	}
	got, ok := blockNBT[key]
	if ok {
		ans, err := GetContainerData(got)
		if err != nil {
			return types.ChestData{}, fmt.Errorf("GetContainerData(Started by GetContainerDataRun): %v", err)
		}
		return ans, nil
	} else {
		return types.ChestData{}, nil
	}
}
