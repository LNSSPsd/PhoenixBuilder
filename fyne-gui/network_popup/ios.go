// +build !ci

// +build ios

package network_popup

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework UIKit -framework UserNotifications

#include <stdlib.h>

void popupNetwork(void);
*/
import "C"
import (

	"fyne.io/fyne/v2"
)

func PopupNetwork(){
	go func() {
		// popup a network permission dialog
		C.popupNetwork()
	}()
}

func defaultVariant() fyne.ThemeVariant {
	return 0
}
