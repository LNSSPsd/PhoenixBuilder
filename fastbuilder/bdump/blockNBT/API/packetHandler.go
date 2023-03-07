package blockNBT_API

import (
	"phoenixbuilder/minecraft/protocol"
	"phoenixbuilder/minecraft/protocol/packet"
)

// 解析数据包并执行相应动作，如更新记录的背包数据等
func (o *PacketHandleResult) HandlePacket(pk *packet.Packet) {
	switch p := (*pk).(type) {
	case *packet.CommandOutput:
		uniqueIdString := p.CommandOrigin.UUID.String()
		o.commandRequestMapLockDown.RLock()
		if o.commandRequest[uniqueIdString] != nil {
			o.commandRequestMapLockDown.RUnlock()
			o.commandRequestMapLockDown.Lock()

			o.commandRequest[uniqueIdString].Responce = *p
			o.commandRequest[uniqueIdString].LockDown.Unlock()

			o.commandRequestMapLockDown.Unlock()
		} else {
			o.commandRequestMapLockDown.RUnlock()
		}
		// send ws command with responce
	case *packet.InventoryContent:
		for key, value := range p.Content {
			if value.Stack.ItemType.NetworkID != -1 {
				o.InventoryDatasMapLockDown.Lock()

				if o.InventoryDatas[p.WindowID] == nil {
					o.InventoryDatas[p.WindowID] = map[uint8]protocol.ItemInstance{uint8(key): value}
				} else {
					o.InventoryDatas[p.WindowID][uint8(key)] = value
				}

				o.InventoryDatasMapLockDown.Unlock()
			}
		}
		// inventory contents(global)
	case *packet.InventoryTransaction:
		for _, value := range p.Actions {
			if value.SourceType == protocol.InventoryActionSourceCreative {
				continue
			}
			o.InventoryDatasMapLockDown.Lock()

			if o.InventoryDatas[uint32(value.WindowID)] == nil {
				o.InventoryDatas[uint32(value.WindowID)] = map[uint8]protocol.ItemInstance{
					uint8(value.InventorySlot): value.NewItem,
				}
			} else {
				o.InventoryDatas[uint32(value.WindowID)][uint8(value.InventorySlot)] = value.NewItem
			}

			o.InventoryDatasMapLockDown.Unlock()
		}
		// inventory contents(for enchant command...)
	case *packet.ItemStackResponse:
		for _, value := range p.Responses {
			//o.ItemStackRequestID = value.RequestID - 2
			o.ItemStackReuqestWithResultMapLockDown.RLock()

			if o.ItemStackReuqestWithResult[value.RequestID] == nil {
				panic("HandlePacket: Attempt to send packet ItemStackRequest without using Bdump/blockNBT API")
			}

			o.ItemStackReuqestWithResultMapLockDown.RUnlock()
			o.ItemStackReuqestWithResultMapLockDown.Lock()

			if value.Status != 0 {
				o.ItemStackReuqestWithResult[value.RequestID].SuccessStates = false
				o.ItemStackReuqestWithResult[value.RequestID].ErrorCode = value.Status
			} else {
				o.ItemStackReuqestWithResult[value.RequestID].SuccessStates = true
			}
			o.ItemStackReuqestWithResult[value.RequestID].LockDown.Unlock()

			o.ItemStackReuqestWithResultMapLockDown.Unlock()
		}
		// item stack request
	case *packet.ContainerOpen:
		o.ContainerCloseDatas.Datas = packet.ContainerClose{}
		o.ContainerOpenDatas.Datas = *p
		unsuccess := o.ContainerOpenDatas.LockDown.TryLock()
		if unsuccess {
			panic("HandlePacket: Attempt to send packet ContainerOpen without using Bdump/blockNBT API")
		}
		o.ContainerOpenDatas.LockDown.Unlock()
		// while open a container
	case *packet.ContainerClose:
		if p.WindowID != 0 && p.WindowID != 119 && p.WindowID != 120 && p.WindowID != 124 {
			o.InventoryDatasMapLockDown.Lock()

			delete(o.InventoryDatas, uint32(p.WindowID))
			newMap := map[uint32]map[uint8]protocol.ItemInstance{}
			for key, value := range o.InventoryDatas {
				newMap[key] = value
			}
			o.InventoryDatas = newMap

			o.InventoryDatasMapLockDown.Unlock()
		}
		o.ContainerOpenDatas.Datas = packet.ContainerOpen{}
		o.ContainerCloseDatas.Datas = *p
		unsuccess := o.ContainerCloseDatas.LockDown.TryLock()
		if !p.ServerSide && unsuccess {
			panic("HandlePacket: Attempt to send packet ContainerClose without using Bdump/blockNBT API")
		}
		o.ContainerCloseDatas.LockDown.Unlock()
		// while a container is closed
	}
}
