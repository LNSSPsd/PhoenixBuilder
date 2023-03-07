package blockNBT_API

import "phoenixbuilder/minecraft/protocol/packet"

// 以下罗列了每次放置方块实体时需要外部实现赋值的 API
type GlobalAPI struct {
	WritePacket        func(packet.Packet) error // 向租赁服发送数据包的函数
	BotName            string                    // 机器人的游戏昵称
	BotIdentity        string                    // 机器人的唯一标识符 [当前还未使用]
	BotUniqueID        int64                     // 机器人的唯一 ID [当前还未使用]
	BotRunTimeID       uint64                    // 机器人的运行时 ID [当前还未使用]
	PacketHandleResult *PacketHandleResult       // 保存包处理结果；由外部实现实时更新
}
