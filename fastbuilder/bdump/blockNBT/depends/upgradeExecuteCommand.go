package blockNBT_depends

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pterm/pterm"
)

type command struct {
	context string
	pointer int
}

func GetNewVersionOfExecuteCommand(cmd string, blockPos [3]int) string {
	got := upgradeExecuteCommandRun(cmd)
	outputMessages(struct {
		oldCommand    string
		newCommand    string
		blockPos      [3]int
		resultType    uint8
		errorLog      error
		errorLocation int
	}{
		oldCommand:    cmd,
		newCommand:    got.upgradeResult,
		blockPos:      blockPos,
		resultType:    got.resultType,
		errorLog:      got.errorLog,
		errorLocation: got.errorLocation,
	})
	return got.upgradeResult
}

func upgradeExecuteCommandRun(cmd string) struct {
	upgradeResult string
	successStates bool
	resultType    uint8 // success=1, error=0, warning=2, noOperation=3
	errorLog      error
	errorLocation int
} {
	r := command{cmd, -1}
	ans, successStates, err, errLocation := r.upgradeExecuteCommand()
	single := struct {
		upgradeResult string
		successStates bool
		resultType    uint8
		errorLog      error
		errorLocation int
	}{
		upgradeResult: cmd,
		successStates: false,
		resultType:    3,
		errorLog:      nil,
	}
	if fmt.Sprint(err) == "upgradeExecuteCommand: searchForExecute: Try using a new version of the execute command" {
		single.resultType = 2
		single.errorLocation = errLocation
		return single
	}
	if err != nil {
		single.resultType = 0
		single.errorLog = fmt.Errorf("upgradeExecuteCommandRun: %v", err)
		single.errorLocation = errLocation
		return single
	}
	if successStates {
		single.resultType = 1
		single.upgradeResult = ans
		return single
	} else {
		single.resultType = 3
		single.upgradeResult = ans
		return single
	}
}

func outputMessages(input struct {
	oldCommand    string
	newCommand    string
	blockPos      [3]int
	resultType    uint8 // success=1, error=0, warning=2, noOperation=3
	errorLog      error
	errorLocation int
}) {
	errPart := command{input.oldCommand, input.errorLocation - 5}
	if input.resultType == 0 {
		pterm.Error.Printf("Failed to upgrade the command \"%v\", because of syntax error \"%v\", and the error maybe occurred in >>>%v<<<, and the command block is placed at (%v,%v,%v)\n", input.oldCommand, input.errorLog, errPart.getPartOfString(10), input.blockPos[0], input.blockPos[1], input.blockPos[2])
	}
	if input.resultType == 1 {
		pterm.Success.Printf("Success to upgrade command \"%v\" to \"%v\"\n", input.oldCommand, input.newCommand)
	}
	if input.resultType == 2 {
		pterm.Warning.Printf("Try using a new version of the execute command \"%v\", and the new execute command maybe is in >>>%v<<<, and the command block is placed at (%v,%v,%v)\n", input.oldCommand, errPart.getPartOfString(10), input.blockPos[0], input.blockPos[1], input.blockPos[2])
	}
}

func (cmd *command) getPartOfString(stringLength int) string {
	if cmd.pointer < 0 {
		cmd.pointer = 0
	}
	end := cmd.pointer + stringLength
	if end > len(cmd.context)-1 {
		return cmd.context[cmd.pointer:]
	} else {
		return cmd.context[cmd.pointer:end]
	}
}

func (cmd *command) jumpSpace() error {
	for {
		if cmd.pointer > len(cmd.context)-1 {
			return fmt.Errorf("jumpSpace: %v out of length(%v)", cmd.pointer, len(cmd.context))
		} else if cmd.getPartOfString(1) == " " || cmd.getPartOfString(1) == "/" {
			cmd.pointer++
		} else {
			return nil
		}
	}
}

func (cmd *command) index(searchingFor string) (int, error) {
	if cmd.pointer > len(cmd.context)-1 {
		return 0, fmt.Errorf("index: %v out of length(%v)", cmd.pointer, len(cmd.context))
	}
	find := strings.Index(cmd.context[cmd.pointer:], searchingFor)
	if find == -1 {
		return 0, fmt.Errorf("index: %v not found", searchingFor)
	} else {
		return find + cmd.pointer, nil
	}
}

func (cmd *command) highSearching(input []string) (struct {
	begin int
	end   int
}, error) {
	ansSave := []struct {
		begin int
		end   int
	}{}
	for _, value := range input {
		got, err := cmd.index(value)
		if err == nil {
			ansSave = append(ansSave, struct {
				begin int
				end   int
			}{
				begin: got,
				end:   got + len(value),
			})
		}
	}
	minRecord := struct {
		begin int
		end   int
	}{
		begin: 2147483647,
	}
	success := false
	for _, value := range ansSave {
		if value.begin < minRecord.begin {
			minRecord = value
			success = true
		}
	}
	if !success {
		return struct {
			begin int
			end   int
		}{}, fmt.Errorf("highSearching: Nothing found")
	}
	return minRecord, nil
}

func (cmd *command) getRightBarrier() (int, error) {
	var quotationMark int = 0
	var barrier int = 0
	var err1 error = nil
	var err2 error = nil
	for {
		quotationMark, err1 = cmd.index("\"")
		barrier, err2 = cmd.index("]")
		if err2 != nil {
			return 0, fmt.Errorf("getRightBarrier: Right barrier not found")
		}
		if err1 != nil {
			return barrier, nil
		} else if quotationMark < barrier {
			cmd.pointer = quotationMark + 1
			for {
				if cmd.pointer > len(cmd.context)-1 {
					return 0, fmt.Errorf("getRightBarrier: Right barrier not found")
				}
				if cmd.getPartOfString(2) == "\\\"" {
					cmd.pointer = cmd.pointer + 2
				} else if cmd.getPartOfString(1) == "\"" {
					cmd.pointer++
					break
				} else {
					cmd.pointer++
				}
			}
		} else {
			return barrier, nil
		}
	}
}

func (cmd *command) searchForExecute() (bool, error) {
	err := cmd.jumpSpace()
	if err != nil {
		return false, nil
	}
	commandHeader := strings.ToLower(cmd.getPartOfString(7))
	if commandHeader == "execute" {
		cmd.pointer = cmd.pointer + 7
		return true, nil
	}
	list := []string{"align", "anchored", "as", "at", "facing", "in", "positioned", "rotated", "if", "unless", "run"}
	for _, value := range list {
		if strings.ToLower(cmd.getPartOfString(len(value))) == value {
			return false, fmt.Errorf("searchForExecute: Try using a new version of the execute command")
		}
	}
	return false, nil
}

func (cmd *command) getSelector() (string, error) {
	err := cmd.jumpSpace()
	if err != nil {
		return "", fmt.Errorf("getSelector: Incomplete parameter")
	}
	if cmd.getPartOfString(1) == "@" {
		ans, err := cmd.highSearching([]string{"@s", "@a", "@p", "@e", "@r", "@initiator", "@c", "@v"})
		if err != nil {
			return "", fmt.Errorf("getSelector: Incomplete selector prefix")
		}
		selector := cmd.context[ans.begin:ans.end]
		cmd.pointer = ans.end
		err = cmd.jumpSpace()
		if err != nil {
			return "", fmt.Errorf("getSelector: Incomplete selector parameter")
		}
		if cmd.getPartOfString(1) != "[" {
			cmd.pointer = ans.end
			return selector, nil
		} else {
			save := cmd.pointer
			transit, err := cmd.getRightBarrier()
			if err != nil {
				return "", fmt.Errorf("getSelector: Incomplete selector parameter")
			}
			cmd.pointer = transit + 1
			return fmt.Sprintf("%v%v", selector, cmd.context[save:transit+1]), nil
		}
	} else if cmd.getPartOfString(1) == "\"" {
		save := cmd.pointer
		cmd.pointer++
		for {
			if cmd.pointer > len(cmd.context)-1 {
				return "", fmt.Errorf("getSelector: Incomplete selector")
			}
			if cmd.getPartOfString(2) == "\\\"" {
				cmd.pointer = cmd.pointer + 2
			} else if cmd.getPartOfString(1) == "\"" {
				cmd.pointer++
				break
			} else {
				cmd.pointer++
			}
		}
		return cmd.context[save:cmd.pointer], nil
	} else {
		transit, err := cmd.highSearching([]string{" ", "^", "~"})
		if err != nil {
			return "", fmt.Errorf("getSelector: Invalid selector")
		}
		save := cmd.pointer
		cmd.pointer = transit.end - 1
		return cmd.context[save : transit.end-1], nil
	}
}

func (cmd *command) getPos() (string, error) {
	ans := []string{}
	for i := 0; i < 3; i++ {
		err := cmd.jumpSpace()
		if err != nil {
			return "", fmt.Errorf("getPos: Incomplete parameter")
		}
		save := cmd.pointer
		switch cmd.getPartOfString(1) {
		case "~":
			cmd.pointer++
		case "^":
			cmd.pointer++
		}
		for {
			if cmd.pointer > len(cmd.context)-1 {
				return "", fmt.Errorf("getPos: EOF not found")
			}
			header := cmd.getPartOfString(1)
			success := true
			for _, value := range []string{"+", "-", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "."} {
				if header == value {
					success = false
					break
				}
			}
			if success {
				ans = append(ans, cmd.context[save:cmd.pointer])
				break
			}
			cmd.pointer++
		}
	}
	for i, value := range ans {
		if value[0] == "^"[0] || value[0] == "~"[0] {
			value = value[1:]
			if value != "" {
				j, err := strconv.ParseFloat(value, 32)
				if err != nil {
					return "", fmt.Errorf("getPos: Invalid position occurred in %v", ans[i])
				}
				value = strconv.FormatFloat(j, 'f', -1, 32)
				if value == "0" || value == "-0" {
					value = ""
				}
			}
			if ans[i][0] == "^"[0] {
				ans[i] = "^" + value
			} else {
				ans[i] = "~" + value
			}
		} else {
			if strings.Contains(value, ".") {
				j, err := strconv.ParseFloat(value, 32)
				if err != nil {
					return "", fmt.Errorf("getPos: Invalid position occurred in %v", ans[i])
				}
				value = strconv.FormatFloat(j, 'f', -1, 32)
				if value == "-0" {
					value = "0"
				}
				if !strings.Contains(value, ".") {
					ans[i] = value + ".0"
				}
			} else {
				j, err := strconv.ParseFloat(value, 32)
				if err != nil {
					return "", fmt.Errorf("getPos: Invalid position occurred in %v", ans[i])
				}
				value = strconv.FormatFloat(j, 'f', -1, 32)
				if value == "-0" {
					value = "0"
				}
				ans[i] = value
			}
		}
	}
	if ans[0][0] == "^"[0] || ans[1][0] == "^"[0] || ans[2][0] == "^"[0] {
		if ans[0][0] != "^"[0] || ans[1][0] != "^"[0] || ans[2][0] != "^"[0] {
			return "", fmt.Errorf("getPos: Incorrect position")
		}
	}
	return fmt.Sprintf("%v %v %v", ans[0], ans[1], ans[2]), nil
}

func (cmd *command) detectBlock() (string, error) {
	err := cmd.jumpSpace()
	if err != nil {
		return "", fmt.Errorf("detectBlock: Incomplete parameter")
	}
	if strings.ToLower(cmd.getPartOfString(6)) == "detect" {
		cmd.pointer = cmd.pointer + 6
		pos, err := cmd.getPos()
		if err != nil {
			return "", fmt.Errorf("detectBlock: Failed to get the position, and the error log is %v", err)
		}
		err = cmd.jumpSpace()
		if err != nil {
			return "", fmt.Errorf("detectBlock: Incomplete parameter")
		}
		startLocation := cmd.pointer
		endLocation, err := cmd.index(" ")
		if err != nil {
			return "", fmt.Errorf("detectBlock: Incomplete parameter")
		}
		cmd.pointer = endLocation + 1
		err = cmd.jumpSpace()
		if err != nil {
			return "", fmt.Errorf("detectBlock: Incomplete parameter")
		}
		spaceLocation, err := cmd.index(" ")
		if err != nil {
			return "", fmt.Errorf("detectBlock: Incomplete parameter")
		}
		cmd.pointer = spaceLocation
		return fmt.Sprintf(" if block %v %v %v", pos, cmd.context[startLocation:endLocation], cmd.context[endLocation+1:spaceLocation]), nil
	} else {
		return "", nil
	}
}

func (reader *command) upgradeExecuteCommand() (string, bool, error, int) {
	ans := []string{}
	for {
		reader.pointer++
		found, err := reader.searchForExecute()
		if err != nil {
			return "", false, fmt.Errorf("upgradeExecuteCommand: %v", err), reader.pointer
		}
		if found {
			selector, err := reader.getSelector()
			if err != nil {
				return "", false, fmt.Errorf("upgradeExecuteCommand: %v", err), reader.pointer
			}
			position, err := reader.getPos()
			if err != nil {
				return "", false, fmt.Errorf("upgradeExecuteCommand: %v", err), reader.pointer
			}
			detect, err := reader.detectBlock()
			if err != nil {
				return "", false, fmt.Errorf("upgradeExecuteCommand: %v", err), reader.pointer
			}
			reader.pointer--
			if position == "~ ~ ~" || position == "^ ^ ^" {
				ans = append(ans, fmt.Sprintf("as %v at @s%v ", selector, detect))
			} else {
				ans = append(ans, fmt.Sprintf("as %v at @s positioned %v%v ", selector, position, detect))
			}
		} else {
			ans = append(ans, reader.context[reader.pointer:])
			break
		}
	}
	if len(ans) <= 1 {
		return strings.Join(ans, ""), false, nil, 0
	} else {
		ans[len(ans)-1] = fmt.Sprintf("run %v", ans[len(ans)-1])
		return fmt.Sprintf("execute %v", strings.Join(ans, "")), true, nil, 0
	}
}
