package NBTAssigner

// 此表描述了现阶段已经支持了的方块实体。
// 键代表方块名，而值代表这种方块应该归属的类型
var SupportBlocksPool map[string]string = map[string]string{
	"command_block":           "CommandBlock",
	"chain_command_block":     "CommandBlock",
	"repeating_command_block": "CommandBlock",
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
	"hopper":             "Container",
	"dispenser":          "Container",
	"dropper":            "Container",
	"jukebox":            "Container",
	"brewing_stand":      "Container",
	"undyed_shulker_box": "Container",
	"shulker_box":        "Container",
	"lectern":            "Container",
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

// 此表描述了现阶段已经支持了的特殊物品，如烟花等物品。
// 键代表物品名，而值代表这种物品应该归属的类型
var SupportItemsPool map[string]string = map[string]string{}

// 此表描述了可被 replaceitem 生效的容器。
// key 代表容器的方块名，而 value 则代表此容器放置物品所使用的复合标签或列表
var SupportContainerPool map[string]string = map[string]string{
	"blast_furnace":      "Items",
	"lit_blast_furnace":  "Items",
	"smoker":             "Items",
	"lit_smoker":         "Items",
	"furnace":            "Items",
	"lit_furnace":        "Items",
	"chest":              "Items", // *
	"barrel":             "Items",
	"trapped_chest":      "Items", // *
	"lectern":            "book",  // *
	"hopper":             "Items",
	"dispenser":          "Items",
	"dropper":            "Items",
	"jukebox":            "RecordItem", // *
	"brewing_stand":      "Items",
	"undyed_shulker_box": "Items", // *
	"shulker_box":        "Items", // *
}
