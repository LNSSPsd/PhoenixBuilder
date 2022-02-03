// 这写插件框架有点折磨啊
// 最终版!一定是最终版!不能重写重新构思了!

// 感谢2401PT佬!!让我充满糨糊的大脑又有了一点生机!!
// 于是我试图仿照2401PT佬的想法, 却不幸理解不来,所以产生了如此屎山的东西(
// 感谢2401PT, awaqwqa, xX7912Xx, CAIMEOX, LNSSPsd... 所有帮助过咱的人!
package plugin_beta

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol/packet"
	"plugin"
	"sync"
	"time"
)

// PacketSender: plugin -> main
// PacketReceiver: main -> plugin

type Plugin struct {
	isInstantiate bool // 机翻, 是否已经有成员
	singleton bool
	block bool // true if it blocks packets.
	priority int
	name string
	
	regMu sync.RWMutex // locked when PluginManager Register and do nothing when PluginManager notify(packets)
	packetReceivers []chan packet.Packet  // plugins get it and push Packet to main process.

	packetSender chan packet.Packet
	rule func(pk packet.Packet) bool // assert pk and return ok. It should be simplified.
}


func(pl *Plugin) appendReceiver() {

}

type IPlugin interface {
	Start(manager PluginManager)
}


type PluginManager struct {	
	conn *minecraft.Conn
	Logger *log.Logger
	
	regMu sync.RWMutex
	plugins map[*Plugin]IPlugin

}
	
// routine it
func (plm *PluginManager) notify(pk packet.Packet) {
	for plugin, iplugin := range plm.plugins {
		// filter first
		if !plugin.rule(pk) {
			continue
		}
		if !plugin.Singleton && plugin.isInstantiate {
			iplugin.Start(*plm)
		}
		plugin.packetSender <- pk
	}
}


// type A, B, C, D, E... struct{} (All kinds of packet.Packet)
// func fu(i packet.Packet, j packet.Packet) {

//     if type(reflect.ValueOf(i)) == type(reflect.ValueOf(j)) {
//         ...
//     }
// }


// from /phoenixbuilder/fastbuilder/plugin/plugin.go
func (plm *PluginManager) loadPlugin() error {
	defer func ()  {
		if err := recover(); err != nil {
			plm.Logger.Printf("[WARNING] Failed to load plugins completely: %s", err)
		}
	}()
	plugindir, err := loadPluginDir()
	
	err = os.MkdirAll(plugindir, 0755)
	if err != nil { plm.Logger.Panicln("Failed to mkdir"); return err }
	
	plugins, err := ioutil.ReadDir(plugindir)
	if err != nil { plm.Logger.Panicln("Failed to read direction."); return err}
	
	for _, plugindir := range plugins {
		path:=fmt.Sprintf("%s/%s",plugins, plugindir.Name())
		if filepath.Ext(path)!=".so" {
			continue
		}

		err := plm.initPlugin(path)
		if err != nil {
			plm.Logger.Printf("Failed to load plugin: %s", path) 
			continue
		}
	}	
	return nil
}


func (plm *PluginManager) initPlugin(path string) error {
	pl, err := plugin.Open( path )
	if err != nil { return err }
	plug, err := pl.Lookup("Baka")
	if err != nil { return err }
	plugin := *plug.(*IPlugin)
	// plugin.Start(plm, )
}


// 选择Lookup一个结构体实例的理由是, 使得插件的handle(就是一般的回调函数)之间有更简单的互通渠道(不过写得加w锁).
// hanle共享其所属实例的字段.
// 插件可以选择单例, 这样当一个handle return前, 不会有新handle产生.
// 当同一插件的不同handle间想要通信时, 应使用指针方法.
// 牺牲了插件编写的简洁性, 换来了一个莫名其妙但是或许有时候会派上用场的特性.

// 后来经过2401PT讲解才知道,这原来是一个很平常不过的方法啊
func (plm *PluginManager) RegisterPlugin(singleton bool, block bool, priority int, name string, pksender chan packet.Packet, rule func(pk packet.Packet) bool) {
	
	pl := Plugin{
		isInstantiate: false,
		singleton: singleton,
		block: block,
		priority: priority,
		name: name,
		
		rule: rule,
	}
	plm.plugins = append(plm.plugins, pl)
}

// Channels that are registered can be losesd. Plugins need to care if the states of receivement from channels is true.
// e.g. value, ok := plm.GetPacket()

// It returns a memory address.

func (plm *PluginManager) RegisterChan(regipl IPlugin, lifetime int) <-chan packet.Packet {
	// regipl: 注册方
	
	receiver := make(chan packet.Packet)
	for pl, ipl := range plm.plugins {
		if ipl == regipl {
			pl.regMu.RLock()
			defer pl.regMu.RUnlock()
			pl.packetReceivers = append(pl.packetReceivers, receiver)
			return receiver
		}
	}
	return nil // it shouldn't be returned
}

func (plm *PluginManager) SendPacket(pk packet.Packet) {
	plm.conn.WritePacket(pk)
}

// // 装饰IPlugin的Start, 用于记录插件有多少个正在运行的routine
// func (pl *Plugin) WaitGroupDecorator( fn func(manager PluginManager)) func(PluginManager){
// 	return func (manager PluginManager)  {
		
// 		result := fn(manager)
// 	}
// }