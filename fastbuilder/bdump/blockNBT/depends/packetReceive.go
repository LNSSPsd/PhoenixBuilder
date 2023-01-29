package blockNBT_depends

import (
	"fmt"
	"phoenixbuilder/fastbuilder/environment"
	"phoenixbuilder/io/commands"
	"phoenixbuilder/minecraft/protocol/packet"
)

var Islistening bool
var ReceivePacket []packet.Packet

func PacketProcessor(Environment *environment.PBEnvironment, NeedWaiting bool, RequestPacketId uint32) ([]packet.Packet, error) {
	if NeedWaiting {
		cmdsender := Environment.CommandSender.(*commands.CommandSender)
		cmdsender.SendWSCommandWithResponce("list")
	}
	// waiting for the packet
	ans := []packet.Packet{}
	for _, value := range ReceivePacket {
		// fmt.Printf("%#v\n", value)
		if value.ID() == RequestPacketId {
			ans = append(ans, value)
		}
	}
	// get datas
	closeProcessor()
	// close processor
	if len(ans) == 0 {
		return []packet.Packet{}, fmt.Errorf("PacketProcessor: packet which numbered %v have been not found", RequestPacketId)
	}
	return ans, nil
}

func InitProcessor() {
	Islistening = true
	ReceivePacket = []packet.Packet{}
}

func closeProcessor() {
	Islistening = false
	ReceivePacket = []packet.Packet{}
}
