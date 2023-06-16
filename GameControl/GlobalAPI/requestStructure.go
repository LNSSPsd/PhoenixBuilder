package GlobalAPI

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"
)

// 请求 request 代表的结构请求并获取 packet.StructureTemplateDataResponse 。
// 当且仅当租赁服响应结构请求时本函数才会返回值。
//
// 请确保在使用此函数前占用了结构资源，否则您将无法收到结构的请求结果
func (g *GlobalAPI) SendStructureRequestWithResponce(
	request *packet.StructureTemplateDataRequest,
) (*packet.StructureTemplateDataResponse, error) {
	g.Resources.Structure.AwaitResponceBeforeSendPacket()
	// prepare
	err := g.WritePacket(request)
	if err != nil {
		return nil, fmt.Errorf("SendStructureRequestWithResponce: %v", err)
	}
	// send packet
	g.Resources.Structure.AwaitResponceAfterSendPacket()
	// await changes
	return g.Resources.Structure.LoadStructureResponceDataAndDelete(), nil
	// return
}
