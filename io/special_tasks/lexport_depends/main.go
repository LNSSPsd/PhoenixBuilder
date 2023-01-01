package lexport_depends

import (
	"fmt"
	"math"
	"phoenixbuilder/fastbuilder/types"
	TranslateNBTInerface "phoenixbuilder/io/special_tasks/lexport_depends/TranslateNBTInterface"
	"strconv"
	"strings"
)

// 用于描述一个区域的基本信息，也就是区域的起点位置及区域的尺寸
type Area struct {
	BeginX int
	BeginY int
	BeginZ int
	SizeX  int
	SizeY  int
	SizeZ  int
}

// 用于描述一个区域的坐标
type AreaLocation struct {
	Posx int
	Posz int
}

// 用于描述一个方块的坐标
type BlockPos struct {
	Posx int
	Posy int
	Posz int
}

/*
用于存放一个 MCBE 的结构；这里面的数据稍微作了一些处理，只保留了需要的部分

如果后期要给这个结构体添加别的东西，请参见本文件中的 GetMCStructureData 函数
*/
type Mcstructure struct {
	info                     Area                           // 用于描述这个结构的基本信息，也就是起点位置及尺寸
	blockPalette             []string                       // 用于存放调色板(方块池)中的方块名
	blockPalette_blockStates []string                       // 用于存放调色板(方块池)中的数据；这里的方块池稍作了处理，只保留了方块状态(string)，且这种方块状态正是 setblock 命令所需要的部分；需要特别说明的是，方块状态里面所有的 TAG_Byte 都被处理成了布尔值，如果有 BUG 记得提 Issue
	blockPalette_blockData   []int16                        // 用于存放调色板(方块池)中的数据；这里的方块池稍作了处理，只保留了方块数据值，也就是附加值(int)；这个东西只是为了支持容器而做的
	foreground               []int16                        // 用于描述一个方块的前景层；这里应该用 int32 的，不过 PhoenixBuilder 只能表示 int16 个方块，所以我这里就省一下内存
	background               []int16                        // 用于描述一个方块的背景层；这里应该用 int32 的，不过 PhoenixBuilder 只能表示 int16 个方块，所以我这里就省一下内存
	blockNBT                 map[int]map[string]interface{} // 用于存放方块实体数据
}

/*
用于拆分一个大区域为若干个小区域；当 useSpecialSplitWay 为真时，将蛇形拆分区域

返回值 []Area 代表一个已经排好顺序的若干个小区域

返回值 map[AreaLocation]int 代表可以通过 区域坐标(AreaLocation) 来访问 []Area 的对应项

因此，返回值 map[int]AreaLocation 是返回值 map[AreaLocation]int 的逆过程
*/
func SplitArea(
	startX int, startY int, startZ int,
	endX int, endY int, endZ int,
	splitSizeX int, splitSizeZ int,
	useSpecialSplitWay bool,
) ([]Area, map[AreaLocation]int, map[int]AreaLocation) {
	if splitSizeX < 0 {
		splitSizeX = splitSizeX * -1
	}
	if splitSizeZ < 0 {
		splitSizeZ = splitSizeZ * -1
	}
	// 考虑一些特殊的情况，此举是为了更高的兼容性
	var save int
	if endX < startX {
		save = startX
		startX = endX
		endX = save
	}
	if endY < startY {
		save = startY
		startY = endY
		endY = save
	}
	if endZ < startZ {
		save = startZ
		startZ = endZ
		endZ = save
	}
	// 考虑一些特殊的情况，此举是为了更高的兼容性
	sizeX := endX - startX + 1
	sizeY := endY - startY + 1
	sizeZ := endZ - startZ + 1
	// 取得 Area 的大小
	chunkX_length := int(math.Ceil(float64(sizeX) / float64(splitSizeX)))
	chunkZ_length := int(math.Ceil(float64(sizeZ) / float64(splitSizeZ)))
	// 取得各轴上需要拆分的区域数
	ans := make([]Area, chunkX_length*chunkZ_length) // 这个东西最终会 return 掉
	areaLoctionToInt := map[AreaLocation]int{}       // 知道了区域的坐标求区域在 []Area 的位置
	IntToareaLoction := map[int]AreaLocation{}       // 知道了区域在 []Area 的位置求区域坐标
	facing := -1                                     // 蛇形处理的时候需要用到这个
	key := -1                                        // 向 ans 插入数据的时候需要用到这个
	// 初始化
	for chunkX := 1; chunkX <= chunkX_length; chunkX++ {
		facing = facing * -1
		BeginX := splitSizeX*(chunkX-1) + startX
		xLength := splitSizeX
		if BeginX+xLength-1 > endX {
			xLength = endX - BeginX + 1
		}
		for chunkZ := 1; chunkZ <= chunkZ_length; chunkZ++ {
			key++ // p = p + 1
			currentChunkZ := chunkZ
			if useSpecialSplitWay && facing == -1 {
				currentChunkZ = chunkZ_length - currentChunkZ + 1
			}
			BeginZ := splitSizeZ*(currentChunkZ-1) + startZ
			zLength := splitSizeZ
			if BeginZ+zLength-1 > endZ {
				zLength = endZ - BeginZ + 1
			}
			ans[key] = Area{
				BeginX: BeginX,
				BeginY: startY,
				BeginZ: BeginZ,
				SizeX:  xLength,
				SizeY:  sizeY,
				SizeZ:  zLength,
			}
			areaLoctionToInt[AreaLocation{chunkX - 1, currentChunkZ - 1}] = key
			IntToareaLoction[key] = AreaLocation{chunkX - 1, currentChunkZ - 1}
		}
	}
	return ans, areaLoctionToInt, IntToareaLoction
}

// 用于提取得到的 MCBE 结构文件中的一些数据，具体拿了什么数据，你可以看返回值字段
func GetMCStructureData(area Area, structure map[string]interface{}) (Mcstructure, error) {
	var value_default map[string]interface{} = map[string]interface{}{}
	var ok bool = false
	var normal = false

	var value_structure map[string]interface{} = map[string]interface{}{}

	var blockPalette = []string{}
	var blockPalette_blockStates []string = []string{}
	var blockPalette_blockData []int16 = []int16{}
	var blockNBT map[int]map[string]interface{} = map[int]map[string]interface{}{}
	var foreground []int16 = []int16{}
	var background []int16 = []int16{}
	// 初始化
	_, ok = structure["structure"]
	if ok {
		value_structure, normal = structure["structure"].(map[string]interface{})
		if normal {
			_, ok = value_structure["palette"]
			if ok {
				value_palette, normal := value_structure["palette"].(map[string]interface{})
				if normal {
					_, ok = value_palette["default"]
					if ok {
						value_default, normal = value_palette["default"].(map[string]interface{})
						if normal {
							_, ok = value_default["block_palette"]
							if ok {
								value_block_palette, normal := value_default["block_palette"].([]interface{})
								if normal {
									for key, value := range value_block_palette {
										got, normal := value.(map[string]interface{})
										if normal {
											_, ok = got["name"]
											if ok {
												value_name, normal := got["name"].(string)
												if normal {
													blockPalette = append(blockPalette, value_name)
												} else {
													return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"][\"block_palette\"][%v][\"name\"]", key)
												}
											} else {
												return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"][\"block_palette\"][%v][\"name\"]", key)
											}
											// get block name
											_, ok = got["states"]
											if ok {
												value_states, normal := got["states"].(map[string]interface{})
												if normal {
													blockStates, err := TranslateNBTInerface.Compound(value_states, true)
													if err != nil {
														return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"][\"block_palette\"][%v][\"states\"]", key)
													} else {
														blockPalette_blockStates = append(blockPalette_blockStates, blockStates)
													}
												} else {
													return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"][\"block_palette\"][%v][\"states\"]", key)
												}
											} else {
												return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"][\"block_palette\"][%v]", key)
											}
											// get block states
											_, ok = got["val"]
											if ok {
												val, normal := got["val"].(int16)
												if normal {
													blockPalette_blockData = append(blockPalette_blockData, val)
												} else {
													return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"][\"block_palette\"][%v][\"val\"]", key)
												}
											} else {
												return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"][\"block_palette\"][%v][\"val\"]", key)
											}
											// get block data
										} else {
											return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"][\"block_palette\"][%v]", key)
										}
									}
								} else {
									return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"][\"block_palette\"]")
								}
							} else {
								return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"][\"block_palette\"]")
							}
						} else {
							return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"]")
						}
					} else {
						return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"][\"default\"]")
					}
				} else {
					return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"palette\"]")
				}
			}
		} else {
			return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"]")
		}
	} else {
		return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"]")
	}
	// 先把方块状态得到，放于列表 blockPalette 中
	// 很抱歉，我写出了金字塔屎山（
	_, ok = value_default["block_position_data"]
	if ok {
		value_block_position_data, normal := value_default["block_position_data"].(map[string]interface{})
		if normal {
			for key, value := range value_block_position_data {
				block_position_data, ok := value.(map[string]interface{})
				if ok {
					location_of_block_position_data, err := strconv.ParseInt(key, 10, 64)
					if err != nil {
						return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"default\"][\"block_position_data\"][%v]", key)
					}
					if blockNBT[int(location_of_block_position_data)] == nil {
						blockNBT[int(location_of_block_position_data)] = make(map[string]interface{})
					}
					blockNBT[int(location_of_block_position_data)] = map[string]interface{}{"block_position_data": block_position_data}
				} else {
					return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"default\"][\"block_position_data\"][%v]", key)
				}
			}
		} else {
			return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"default\"][\"block_position_data\"]")
		}
	}
	// 众所不一定周知，这个方块实体数据可能是不存在的(当然这个我没测试过)
	// 然后找到所有的方块实体数据，放于 map(blockNBT) 中
	_, ok = value_structure["block_indices"]
	if ok {
		value_block_indices, normal := value_structure["block_indices"].([]interface{})
		if normal {
			if len(value_block_indices) == 2 {
				value_block_indices_0, normal := value_block_indices[0].([]interface{})
				if normal {
					for blockLocation_key, blockLocation := range value_block_indices_0 {
						got, normal := blockLocation.(int32)
						if normal {
							foreground = append(foreground, int16(got))
						} else {
							return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"block_indices\"][0][%v]", blockLocation_key)
						}
					}
				} else {
					return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"block_indices\"][0]")
				}
				value_block_indices_1, normal := value_block_indices[1].([]interface{})
				if normal {
					for blockLocation_key, blockLocation := range value_block_indices_1 {
						got, normal := blockLocation.(int32)
						if normal {
							background = append(background, int16(got))
						} else {
							return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"block_indices\"][1][%v]", blockLocation_key)
						}
					}
				} else {
					return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"block_indices\"][1]")
				}
			} else {
				return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"block_indices\"]")
			}
		} else {
			return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"block_indices\"]")
		}
	} else {
		return Mcstructure{}, fmt.Errorf("GetMCStructureData: Crashed in input[\"structure\"][\"block_indices\"]")
	}
	// 然后分别拿到方块池的前景层和背景层
	return Mcstructure{
		info:                     area,
		blockPalette:             blockPalette,
		blockPalette_blockStates: blockPalette_blockStates,
		blockPalette_blockData:   blockPalette_blockData,
		foreground:               foreground,
		background:               background,
		blockNBT:                 blockNBT,
	}, nil
}

// 根据 mcstructure 的起点和尺寸，以及提供的方块坐标，寻找这个方块在 mcstructure 中的角标
func SearchForBlock(structureInfo Area, pos BlockPos) (int, error) {
	pos.Posx = pos.Posx - structureInfo.BeginX
	pos.Posy = pos.Posy - structureInfo.BeginY
	pos.Posz = pos.Posz - structureInfo.BeginZ
	// 将方块的绝对坐标转换为相对坐标(相对于 mcstructure)
	blockCount := structureInfo.SizeX * structureInfo.SizeY * structureInfo.SizeZ
	// 计算结构的尺寸
	angleMark := 0
	angleMark = angleMark + structureInfo.SizeY*structureInfo.SizeZ*pos.Posx
	angleMark = angleMark + structureInfo.SizeZ*pos.Posy
	angleMark = angleMark + pos.Posz
	// 计算方块相对于 mcstructure 的角标
	if angleMark > blockCount-1 {
		return -1, fmt.Errorf("Index out of the list, occured in input[%v]", angleMark)
	}
	return angleMark, nil
}

func ExportBaseOnChunk(
	allAreas []Mcstructure,
	allAreasFindUse map[AreaLocation]int,
	currentExport Area,
) ([]*types.Module, error) {
	ans := make([]*types.Module, 0)
	// 这个东西最后会 return 掉
	allChunks, _, allChunksFindUse := SplitArea(
		currentExport.BeginX, currentExport.BeginY, currentExport.BeginZ,
		currentExport.BeginX+currentExport.SizeX-1,
		currentExport.BeginY+currentExport.SizeY-1,
		currentExport.BeginZ+currentExport.SizeZ-1,
		16, 16, true,
	)
	// 将所有待导出区域按 16*16 的大小拆分为区块，且蛇形拆分
	// 然后按照得到的结果重排处理
	for key, value := range allChunks {
		chunkPos := allChunksFindUse[key]
		chunkPos.Posx = int(math.Floor(float64(chunkPos.Posx) / 4))
		chunkPos.Posz = int(math.Floor(float64(chunkPos.Posz) / 4))
		// 取得当前遍历的区块的坐标
		// 这里已经把坐标变换到 allAreas 下的坐标系中
		targetAreaPos := allAreasFindUse[chunkPos]
		targetArea := allAreas[targetAreaPos]
		// 取得被遍历区块对应的 mcstructure
		i, _, _ := SplitArea(
			value.BeginX, value.BeginY, value.BeginZ,
			value.BeginX+value.SizeX-1,
			value.BeginY+value.SizeY-1,
			value.BeginZ+value.SizeZ-1,
			1, 1, true,
		)
		allBlocksInCurrentChunk := make([]int32, 0)
		for _, VALUE := range i {
			got, err := SearchForBlock(targetArea.info, BlockPos{
				Posx: VALUE.BeginX,
				Posy: VALUE.BeginY,
				Posz: VALUE.BeginZ,
			})
			if err != nil {
				return []*types.Module{}, fmt.Errorf("SearchForBlock(Started by ExportBaseOnChunk): %v", err)
			}
			allBlocksInCurrentChunk = append(allBlocksInCurrentChunk, int32(got))
		}
		// 枚举出被遍历区块中所有方块的坐标(只枚举其中一层)
		for KEY, VALUE := range allBlocksInCurrentChunk {
			VALUE = VALUE - int32(targetArea.info.SizeZ)
			// 这个前置处理方法可能不太优雅
			// 凑合着用吧
			for j := 0; j < targetArea.info.SizeY; j++ {
				VALUE = VALUE + int32(targetArea.info.SizeZ)
				// 前往下一层
				foreground_blockName := "undefined"
				background_blockName := "undefined"
				foreground_blockStates := "undefined"
				background_blockStates := "undefined"
				foreground_blockData := int16(-1)
				// 初始化
				fgId := targetArea.foreground[VALUE] // 前景层方块在调色板中的id
				bgId := targetArea.background[VALUE] // 背景层方块在调色板中的id
				if fgId != -1 {
					foreground_blockName = strings.Replace(targetArea.blockPalette[fgId], "minecraft:", "", 1) // 前景层方块的名称
					foreground_blockStates = targetArea.blockPalette_blockStates[fgId]                         // 前景层方块的方块状态
					foreground_blockData = targetArea.blockPalette_blockData[fgId]                             // 前景层方块的方块数据值(附加值)
				}
				if bgId != -1 {
					background_blockName = strings.Replace(targetArea.blockPalette[bgId], "minecraft:", "", 1) // 背景层方块的名称
					background_blockStates = targetArea.blockPalette_blockStates[bgId]                         // 背景层方块的方块状态
				}
				// 获得基本信息
				var hasNBT bool = false
				var containerDataMark bool = false
				var containerData types.ChestData = types.ChestData{}
				var commandBlockDataMark bool = false
				var commandBlockData types.CommandBlockData = types.CommandBlockData{}
				var string_nbt string = ""
				var err error = fmt.Errorf("")

				got, ok := targetArea.blockNBT[int(VALUE)]
				if ok {
					_, ok := got["block_position_data"]
					if ok {
						block_position_data, normal := got["block_position_data"].(map[string]interface{})
						if normal {
							_, ok := block_position_data["block_entity_data"]
							if ok {
								block_entity_data, normal := block_position_data["block_entity_data"].(map[string]interface{})
								if normal {
									containerData, err = TranslateNBTInerface.GetContainerDataRun(block_entity_data, foreground_blockName)
									if fmt.Sprintf("%v", err) != "GetContainerDataRun: Not a container" && err != nil {
										return []*types.Module{}, fmt.Errorf("%v", err)
									} else if err == nil {
										if foreground_blockName == "chest" {
											useOfChest := "trapped_chest"
											ans = append(ans, &types.Module{
												Block: &types.Block{
													Name: &useOfChest,
													Data: 0,
												},
												Point: types.Position{
													X: i[KEY].BeginX - currentExport.BeginX,
													Y: i[KEY].BeginY + j - currentExport.BeginY,
													Z: i[KEY].BeginZ - currentExport.BeginZ,
												},
											})
										}
										// 这么处理是为了解决箱子间的连接问题，让所有的箱子都不再连接
										if foreground_blockName == "trapped_chest" {
											useOfChest := "chest"
											ans = append(ans, &types.Module{
												Block: &types.Block{
													Name: &useOfChest,
													Data: 0,
												},
												Point: types.Position{
													X: i[KEY].BeginX - currentExport.BeginX,
													Y: i[KEY].BeginY + j - currentExport.BeginY,
													Z: i[KEY].BeginZ - currentExport.BeginZ,
												},
											})
										}
										// 这么处理是为了解决箱子间的连接问题，让所有的箱子都不再连接
										containerDataMark = true
									}
									// 容器
									if foreground_blockName == "command_block" || foreground_blockName == "repeating_command_block" || foreground_blockName == "chain_command_block" {
										commandBlockData, err = TranslateNBTInerface.GetCommandBlockData(block_entity_data, foreground_blockName)
										if err != nil {
											return []*types.Module{}, fmt.Errorf("GetCommandBlockData(Started by ExportBaseOnChunk): %v", err)
										}
										commandBlockDataMark = true
									}
									// 命令方块
									hasNBT = true
									string_nbt, err = TranslateNBTInerface.Compound(block_entity_data, false)
									string_nbt = fmt.Sprintf("{\"block_entity_data\": %v}", string_nbt)
									if err != nil {
										return []*types.Module{}, fmt.Errorf("%v", err)
									}
									// 取得 snbt
								} else {
									return []*types.Module{}, fmt.Errorf("ExportBaseOnChunk: Crashed by invalid \"block_entity_data\"")
								}
							} else {
								//return []*types.Module{}, fmt.Errorf("ExportBaseOnChunk: Crashed by could not found \"block_entity_data\"")
							}
						} else {
							return []*types.Module{}, fmt.Errorf("ExportBaseOnChunk: Crashed by invalid \"block_position_data\"")
						}
					} else {
						return []*types.Module{}, fmt.Errorf("ExportBaseOnChunk: Crashed by could not found \"block_position_data\"")
					}
				}
				// 取得方块实体数据
				if foreground_blockName != "" && (background_blockName == "water" || background_blockName == "flowing_water") {
					ans = append(ans, &types.Module{
						Block: &types.Block{
							Name:        &background_blockName,
							BlockStates: &background_blockStates,
						},
						Point: types.Position{
							X: i[KEY].BeginX - currentExport.BeginX,
							Y: i[KEY].BeginY + j - currentExport.BeginY,
							Z: i[KEY].BeginZ - currentExport.BeginZ,
						},
					})
				}
				// 含水类方块
				if foreground_blockName != "" && foreground_blockName != "air" {
					if hasNBT && commandBlockDataMark {
						ans = append(ans, &types.Module{
							Block: &types.Block{
								Name: &foreground_blockName,
								Data: uint16(foreground_blockData),
							},
							CommandBlockData: &commandBlockData,
							NBTData:          []byte(string_nbt),
							Point: types.Position{
								X: i[KEY].BeginX - currentExport.BeginX,
								Y: i[KEY].BeginY + j - currentExport.BeginY,
								Z: i[KEY].BeginZ - currentExport.BeginZ,
							},
						})
					} else if hasNBT && containerDataMark {
						ans = append(ans, &types.Module{
							Block: &types.Block{
								Name: &foreground_blockName,
								Data: uint16(foreground_blockData),
							},
							NBTData:   []byte(string_nbt),
							ChestData: &containerData,
							Point: types.Position{
								X: i[KEY].BeginX - currentExport.BeginX,
								Y: i[KEY].BeginY + j - currentExport.BeginY,
								Z: i[KEY].BeginZ - currentExport.BeginZ,
							},
						})
					} else if hasNBT {
						ans = append(ans, &types.Module{
							Block: &types.Block{
								Name:        &foreground_blockName,
								BlockStates: &foreground_blockStates,
							},
							NBTData: []byte(string_nbt),
							Point: types.Position{
								X: i[KEY].BeginX - currentExport.BeginX,
								Y: i[KEY].BeginY + j - currentExport.BeginY,
								Z: i[KEY].BeginZ - currentExport.BeginZ,
							},
						})
					} else {
						ans = append(ans, &types.Module{
							Block: &types.Block{
								Name:        &foreground_blockName,
								BlockStates: &foreground_blockStates,
							},
							Point: types.Position{
								X: i[KEY].BeginX - currentExport.BeginX,
								Y: i[KEY].BeginY + j - currentExport.BeginY,
								Z: i[KEY].BeginZ - currentExport.BeginZ,
							},
						})
					}
				}
				// 放置前景层的方块
			}
		}
	}
	return ans, nil
}
