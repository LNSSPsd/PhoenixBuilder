package blockNBT_API

import (
	"phoenixbuilder/minecraft/protocol/packet"
	"sync"
)

/*
初始化 PacketHandleResult 结构体中的各个参数

这件事已经在 PheonixBuilder 启动时做过了，非必要请不要使用此函数
*/
func (p *PacketHandleResult) InitValue() {
	p.commandRequest = make(map[string]*CommandRequest)
	p.commandRequestMapLockDown = sync.RWMutex{}
	// -----
	p.InventoryDatas = make(InventoryContents)
	p.InventoryDatasMapLockDown = sync.RWMutex{}
	// -----
	p.ItemStackReuqestWithResult = make(map[int32]*ItemStackReuqestWithAns)
	p.ItemStackReuqestWithResultMapLockDown = sync.RWMutex{}
	p.ItemStackRequestID = -1
	// -----
	p.ContainerOpenDatas = ContainerOpen{LockDown: sync.Mutex{}, Datas: packet.ContainerOpen{}}
	p.ContainerCloseDatas = ContainerClose{LockDown: sync.Mutex{}, Datas: packet.ContainerClose{}}
}
