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
	"reflect"
)

// PacketSender: packet -> main
// PacketReceiver: main -> packet

type Plugin struct {
	isInstantiate bool // 机翻, 实例化
	Singleton bool
	Block bool // true if it blocks packets.
	Priority int
	Name string
	Main interface{}
	PacketSender chan packet.Packet
	rule reflect.Type
}

type IPlugin interface {
	Start(manager PluginManager, PacketReceiver chan packet.Packet) (error)
}


type PluginManager struct {	
	conn *minecraft.Conn
	plugins map[Plugin]IPlugin
	Logger *log.Logger
	PacketReceiver chan packet.Packet // plugins get it and push Packet to main process.
}
	

func (pls *PluginManager) notify(pk packet.Packet) {
	for plugin, iplugin := range pls.plugins {
		// filter first
		if reflect.TypeOf(pk) != plugin.rule {
			continue
		}
		plugin		
		plugin.PacketSender <- pk
		
		}
	}
}

// routine it.
func (pls *PluginManager) forward() {

	for {
		pk := <- pls.PacketReceiver	
		if pk == nil { continue }
		// todo check the type of pk!!
		pls.conn.WritePacket(pk)
	}
}// if pk.Unmarshal()

// from /phoenixbuilder/fastbuilder/plugin/plugin.go
func (pls *PluginManager) loadPlugin() error {
	defer func ()  {
		if err := recover(); err != nil {
			pls.Logger.Printf("[WARNING] Failed to load plugins completely: %s", err)
		}
	}()
	plugindir, err := loadPluginDir()
	
	err = os.MkdirAll(plugindir, 0755)
	if err != nil { pls.Logger.Panicln("Failed to mkdir"); return err }
	
	plugins, err := ioutil.ReadDir(plugindir)
	if err != nil { pls.Logger.Panicln("Failed to read direction."); return err}
	
	for _, plugindir := range plugins {
		path:=fmt.Sprintf("%s/%s",plugins, plugindir.Name())
		if filepath.Ext(path)!=".so" {
			continue
		}

		err := pls.initPlugin(path)
		if err != nil {
			pls.Logger.Printf("Failed to load plugin: %s", path) 
			continue
		}
	}	
	return nil
}


func (pls *PluginManager) initPlugin(path string) error {
	pl, err := plugin.Open( path )
	if err != nil { return err }
	plug, err := pl.Lookup("Baka")
	if err != nil { return err }
	plugin := *plug.(*IPlugin)
	// plugin.Start(pls, )
}


// 选择Lookup一个结构体实例的理由是, 使得插件的handle(就是一般的回调函数)之间有更简单的互通渠道(不过写得加w锁).
// hanle共享其所属实例的字段.
// 插件可以选择单例, 这样当一个handle return前, 不会有新handle产生.
// 当同一插件的不同handle间想要通信时, 应使用指针方法.
// 牺牲了插件编写的简洁性, 换来了一个莫名其妙但是或许有时候会派上用场的特性.
func (pls *PluginManager) RegisterPlugin(singleton bool, block bool, priority int, name string, pksender chan packet.Packet, rule reflect.Type) {
	pl := Plugin{
		isInstantiate: false,
		Singleton: singleton,
		Block: block,
		Priority: priority,
		Name: name,
		PacketSender: pksender,
		rule: rule,
	}
	pls.plugins = append(pls.plugins, pl)
}