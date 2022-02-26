// +build !ios ci

package network_popup

import "net/http"

func PopupNetwork(){
	go func() {
		// popup a network permission dialog
		http.Get("1.1.1.1")
	}()
}
