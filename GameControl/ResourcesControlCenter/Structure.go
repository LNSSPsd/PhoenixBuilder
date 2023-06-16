package ResourcesControlCenter

import "phoenixbuilder/minecraft/protocol/packet"

// 用于在请求结构前执行。
// 便于后续调用 AwaitResponceAfterSendPacket 以阻塞程序的执行从而
// 达到等待租赁服响应结构请求的目的
func (m *mcstructure) AwaitResponceBeforeSendPacket() {
	m.awaitChanges.Lock()
}

// 用于在请求结构后执行。
// 用于等待租赁服响应结构请求。
// 在调用此函数后，会持续阻塞直到相关操作所对应的互斥锁被释放
func (m *mcstructure) AwaitResponceAfterSendPacket() {
	m.awaitChanges.Lock()
	m.awaitChanges.Unlock()
}

// 写入租赁服响应的结构请求数据，同时释放 m.awaitChanges 中关于容器操作的互斥锁并写入。
// 如果互斥锁未被锁定，程序也仍不会发生惊慌。
// 当且仅当租赁服响应客户端发送的结构请求后，此函数才会被调用。
// 属于私有实现
func (m *mcstructure) writeStructureResponce(
	resp *packet.StructureTemplateDataResponse,
) {
	m.responce.lockDown.Lock()
	defer m.responce.lockDown.Unlock()
	// init
	m.responce.datas = resp
	// write datas
	m.awaitChanges.TryLock()
	m.awaitChanges.Unlock()
	// release lock
}

// 加载租赁服返回的结构请求结果并删除它。
// 如果当前未发现租赁服返回了结果，那么返回值将为 nil
func (m *mcstructure) LoadStructureResponceDataAndDelete() *packet.StructureTemplateDataResponse {
	m.responce.lockDown.Lock()
	defer m.responce.lockDown.Unlock()
	// init
	if m.responce.datas == nil {
		return nil
	}
	// if their is no responces
	new := *m.responce.datas
	m.responce.datas = nil
	// load datas and delete it
	return &new
	// return
}
