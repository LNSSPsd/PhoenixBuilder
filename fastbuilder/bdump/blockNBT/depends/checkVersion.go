package blockNBT_depends

import "phoenixbuilder/minecraft/protocol"

func CheckVersion() bool {
	return protocol.CurrentProtocol == 504
}
