package blockNBT_depends

import (
	"bytes"
	"fmt"

	"github.com/Tnze/go-mc/nbt"
)

func ParseStringNBT(stringNBT *string) (*interface{}, error) {
	var buf bytes.Buffer
	err := nbt.NewEncoder(&buf).Encode(nbt.StringifiedMessage(*stringNBT), "")
	if err != nil {
		return nil, fmt.Errorf("ParseStringNBT: %v", err)
	}
	var BlockNBT interface{}
	nbt.Unmarshal(buf.Bytes(), &BlockNBT)
	return &BlockNBT, nil
}

// 得益于那个库兼容性不太强的原因，方块状态里面的布尔值在解析后会变成字符串
func ParseStringBlockStates(stringBlockStates *string) (*map[string]interface{}, error) {
	input := *stringBlockStates
	for {
		if input[0] == " "[0] {
			input = input[1:]
		} else {
			break
		}
	}
	for {
		if input[len(input)-1] == " "[0] {
			input = input[:len(input)-1]
		} else {
			break
		}
	}
	input = fmt.Sprintf("{%v}", input[1:len(input)-1])
	got, err := ParseStringNBT(&input)
	if err != nil {
		return &map[string]interface{}{}, fmt.Errorf("ParseStringBlockStates: Could not parse this block state; errorLog = %v; stringBlockStates = %#v", err, *stringBlockStates)
	}
	GOT := *got
	ans, normal := GOT.(map[string]interface{})
	if !normal {
		return &map[string]interface{}{}, fmt.Errorf("ParseStringBlockStates: This string is a string-nbt, but not describe a block state; stringBlockStates = %#v", *stringBlockStates)
	}
	return &ans, nil
}
