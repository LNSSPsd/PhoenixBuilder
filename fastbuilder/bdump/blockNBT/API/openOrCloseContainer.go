package blockNBT_API

import (
	"fmt"
	"phoenixbuilder/minecraft/protocol/packet"
)

// 打开 pos 处名为 blockName 且方块状态为 blockStates 的容器；只有当打开完成后才会返回值
func (g *GlobalAPI) OpenContainer(
	pos [3]int32,
	blockName string,
	blockStates map[string]interface{},
) error {
	g.PacketHandleResult.ContainerOpenDatas.LockDown.Lock()
	defer g.PacketHandleResult.ContainerOpenDatas.LockDown.Unlock()
	// init
	err := g.UseItemOnBlocks(0, pos, blockName, blockStates, false)
	if err != nil {
		return fmt.Errorf("OpenContainer: %v", err)
	}
	// open container
	g.PacketHandleResult.ContainerOpenDatas.LockDown.Lock()
	g.PacketHandleResult.ContainerOpenDatas.LockDown.Unlock()
	// wait changes
	return nil
	// return
}

// 关闭已经打开的容器；如果容器已被关闭，则返回的布尔值为 false；只有当容器被关闭后才会返回值
func (g *GlobalAPI) CloseContainer() (bool, error) {
	g.PacketHandleResult.ContainerCloseDatas.LockDown.Lock()
	defer g.PacketHandleResult.ContainerCloseDatas.LockDown.Unlock()
	// init
	err := g.WritePacket(&packet.ContainerClose{
		WindowID:   g.PacketHandleResult.ContainerOpenDatas.Datas.WindowID,
		ServerSide: false,
	})
	if err != nil {
		return false, fmt.Errorf("CloseContainer: %v", err)
	}
	// close container
	g.PacketHandleResult.ContainerCloseDatas.LockDown.Lock()
	g.PacketHandleResult.ContainerCloseDatas.LockDown.Unlock()
	// wait changes
	return true, nil
	// return
}
