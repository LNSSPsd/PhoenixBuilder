package blockNBT_API

import (
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
	"sync"
)

// 存放数据包的处理结果；理论上，这些结果应该由此结构体的 HandlePacket 方法实时更新
type PacketHandleResult struct {
	commandRequest            map[string]*CommandRequest // 存放命令请求及结果；此参数不对外公开
	commandRequestMapLockDown sync.RWMutex               // 对 commandRequest 参数进行并发读写时需要上锁
	// -----
	InventoryDatas            InventoryContents // 当机器人打开容器，此时候租赁服会发送相关的包用以描述容器内的物品数据；当容器关闭，将会自动移除本 InventoryContent 变量中对应的 value；不过容器也可能会被破坏，因此建议在使用完一个容器后，仍去尝试删除一次对应的 value 以及清空 ContainerOpenDatas
	InventoryDatasMapLockDown sync.RWMutex      // 对 InventoryDatas 参数进行并发读写时需要上锁
	// -----
	ItemStackReuqestWithResult            map[int32]*ItemStackReuqestWithAns // 在客户端发送 ItemStackRequest 后，租赁服会发送对应的 ItemStackResponce；当 ItemStackResponce 的 Stutas 字段为 0 时，视为对应的 ItemStackRequest 被租赁服通过，此时 SuccessStates 会被处理为 true，否则被处理为 false；当请求被回复后，请手动删除相应的键以释放数据
	ItemStackReuqestWithResultMapLockDown sync.RWMutex                       // 对 ItemStackReuqestWithResult 参数进行并发读写时需要上锁
	ItemStackRequestID                    int32                              // 客户端在发送 ItemStackRequest 时需要发送一个 RequestID；经过观察，这个值会随着请求发送的次数递减，且呈现为公差为 -2，首项为 -1 的递减型等差数列；特别地，如果你尝试在 RequestID 字段填写非负数或者偶数，那么客户端会被强制断开连接；因此，为了安全性，请在每次发送 ItemStackRequest 将本数值自减 2；尽管始终为 ItemStackRequest 的 RequestID 字段发送 -1 并不会造成太大的问题，但我仍然建议您使用这个变量充当 RequestID ；在修改这个值时，请务必使用以原子操作执行，例如 atomic.AddInt32 函数
	// -----
	ContainerOpenDatas  ContainerOpen  // 存放容器的打开状态及相应的数据
	ContainerCloseDatas ContainerClose // 存放容器的关闭结果及相应的数据
}

// 存放命令请求及结果
type CommandRequest struct {
	LockDown sync.Mutex           // 如果命令请求已被发出，此参数将被锁定且对应的 go 协程进入阻塞态。当命令请求得到反馈时，由对应的函数解除此参数的锁定，相应的 go 协程从阻塞变为非阻塞，然后返回值
	Responce packet.CommandOutput // 当命令请求得到反馈时此参数将被赋值
}

// uint32 代表打开的窗口 ID ，即 WindowID；
// uint8 代表槽位；
// 最内层的 protocol.ItemInstance 代表对应的 WindowID 库存中相应槽位的物品数据。
type InventoryContents map[uint32]map[uint8]protocol.ItemInstance

// 存放物品更改请求及结果
type ItemStackReuqestWithAns struct {
	LockDown      sync.Mutex // 在客户端发送请求 ItemStackRequest 前，此参数应当被锁定。当租赁服返回了对应的 ItemStackResponce 后，此参数将被其他函数解锁。如果在发送 ItemStackRequest 前未锁定此参数，将会造成程序惊慌
	SuccessStates bool       // 记录相应 ItemStackRequest 的操作结果，为真时代表成功
	ErrorCode     uint8      // ItemStackResponce 的 status 字段，一般情况下不会有任何帮助
}

// 存放容器的打开状态及打开数据
type ContainerOpen struct {
	LockDown sync.Mutex           // 如果需要打开一个容器，请为此上锁。当容器成功被打开后，其他的函数会解锁它。如果在打开容器未锁定此参数，将会造成程序惊慌
	Datas    packet.ContainerOpen // 当客户端成功打开容器时，租赁服会发送一个 ContainerOpen 数据包以说明打开的容器的相关信息，但此数据包不会包含物品信息。当租赁服以此数据包回应容器打开请求时，此参会会被赋值。如果容器被关闭，那么此参数会被恢复为 packet.ContainerOpen{} 。请不要尝试修改此参数，否则可能会造成错误
}

// 存放容器的关闭结果及相应的数据
type ContainerClose struct {
	LockDown sync.Mutex            // 如果需要关闭一个容器，请为此上锁。当容器成功被关闭后，其他的函数会解锁它。如果容器不是由租赁服强制关闭，且在关闭容器前未锁定此参数，那么将会造成程序惊慌
	Datas    packet.ContainerClose // 当客户端关闭容器时，需要发送一个 ContainerClose 数据包。当成功关闭时，租赁服会以相同的数据包回应。当租赁服以此数据包回应容器关闭请求时，此参数会被赋值。如果容器被关闭，那么此参数会被恢复为 packet.ContainerClose{} 。请不要尝试修改此参数，否则可能会造成错误
}
