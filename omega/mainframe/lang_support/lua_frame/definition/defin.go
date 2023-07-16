package definition

// 消息包
type PackageChan struct {
	//监听包的类型
	PackageType string
	//通道
	PackageMsgChan chan interface{}
}
