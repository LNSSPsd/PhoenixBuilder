package TranslateNBTInerface

import (
	"fmt"
	"strconv"
	"strings"
)

// 判断 nbt 中 value 的数据类型
func GetData(input interface{}) (string, error) {
	value1, result := input.(byte)
	if result {
		return fmt.Sprintf("%vb", int(value1)), nil
	}
	// byte
	value2, result := input.(int16)
	if result {
		return fmt.Sprintf("%vs", value2), nil
	}
	// short
	value3, result := input.(int32)
	if result {
		return fmt.Sprintf("%v", value3), nil
	}
	// int
	value4, result := input.(int64)
	if result {
		return fmt.Sprintf("%vl", value4), nil
	}
	// long
	value5, result := input.(float32)
	if result {
		return fmt.Sprintf("%vf", strconv.FormatFloat(float64(value5), 'f', 16, 32)), nil
	}
	// float
	value6, result := input.(float64)
	if result {
		return fmt.Sprintf("%vd", strconv.FormatFloat(float64(value6), 'f', 16, 64)), nil
	}
	// double
	value7, result := input.([]byte)
	if result {
		ans := []string{}
		for _, i := range value7 {
			ans = append(ans, fmt.Sprintf("%vb", int(i)))
		}
		return fmt.Sprintf("[B; %v]", strings.Join(ans, ", ")), nil
	}
	// byte_array
	value8, result := input.(string)
	if result {
		return fmt.Sprintf("\"%v\"", value8), nil
	}
	// string
	value9, result := input.([]interface{})
	if result {
		list, err := List(value9)
		if err != nil {
			return "", fmt.Errorf("GetData: Failed in %#v", value9)
		}
		return list, nil
	}
	// list
	value10, result := input.(map[string]interface{})
	if result {
		compound, err := Compound(value10, false)
		if err != nil {
			return "", fmt.Errorf("GetData: Failed in %#v", value10)
		}
		return compound, nil
	}
	// compound
	value11, result := input.([]int32)
	if result {
		ans := []string{}
		for _, i := range value11 {
			ans = append(ans, fmt.Sprintf("%v", i))
		}
		return fmt.Sprintf("[I; %v]", strings.Join(ans, ", ")), nil
	}
	// int_array
	value12, result := input.([]int64)
	if result {
		ans := []string{}
		for _, i := range value12 {
			ans = append(ans, fmt.Sprintf("%v", i))
		}
		return fmt.Sprintf("[L; %v]", strings.Join(ans, ", ")), nil
	}
	// long_array
	return "", fmt.Errorf("GetData: Failed because of unknown type of the target data, occured in %#v", input)
}

func Compound(input map[string]interface{}, outputBlockStatesMode bool) (string, error) {
	ans := make([]string, 0)
	for key, value := range input {
		if value == nil {
			return "", fmt.Errorf("Compound: Crashed in input[\"%v\"]; errorLogs = value is nil; input = %#v", key, input)
		}
		got, err := GetData(value)
		if err != nil {
			return "", fmt.Errorf("Compound: Crashed in input[\"%v\"]; errorLogs = %v; input = %#v", key, err, input)
		} else {
			if got[len(got)-1] == "b"[0] && outputBlockStatesMode {
				if got == "0b" {
					got = "false"
				} else if got == "1b" {
					got = "true"
				} else {
					return "", fmt.Errorf("Compound: Crashed in input[\"%v\"]; errorLogs = outputBlockStatesModeError; input = %#v", key, input)
				}
			}
			ans = append(ans, fmt.Sprintf("\"%v\": %v", key, got))
		}
	}
	if outputBlockStatesMode {
		return fmt.Sprintf("[%v]", strings.Join(ans, ", ")), nil
	}
	return fmt.Sprintf("{%v}", strings.Join(ans, ", ")), nil
}

func List(input []interface{}) (string, error) {
	ans := make([]string, 0)
	for key, value := range input {
		if value == nil {
			return "", fmt.Errorf("List: Crashed in input[\"%v\"]; errorLogs = value is nil; input = %#v", key, input)
		}
		got, err := GetData(value)
		if err != nil {
			return "", fmt.Errorf("List: Crashed in input[\"%v\"]; errorLogs = %v; input = %#v", key, err, input)
		}
		ans = append(ans, got)
	}
	return fmt.Sprintf("[%v]", strings.Join(ans, ", ")), nil
}
