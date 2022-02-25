// +build fyne_gui

package args
import "C"

func AuthServer() string {
	return "wss://api.fastbuilder.pro:2053/"
}

func ShouldDisableHashCheck() bool {
	return false
}
