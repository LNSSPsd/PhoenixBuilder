package plugin_beta

import (
	"fmt"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol/packet"
	"reflect"
	"plugin"
	
	"time"
)

// function that has been registered
// type RegisteredFunction func(interface{})

// e.g. RegisterFunction {
//
//
// }

/*
e.g.

func _MainFunction(ps PluginSystem)  {
	fn := func (pk packet.Text)  {
		// ps.Conn.WritePacket()
		ps.Logger.log("hi!")
	}
	plugin := RegisteredPlugin {
		Singleton: false,
		Block: true,
		Priority: 5,
		Name: "tpa",
		Main: fn,
	}
	ps.RegisterFunc(plugin)
}
*/


type RegisteredPlugin struct {
	
	Singleton bool
	Block bool // true if it blocks packets.
	Priority int
	Name string
	Main interface{}

	 // True: Only one example can exist at the same time
	
}


// show plugins
type IPluginSystem interface {
	RegisterFunc(RegisteredPlugin)
	RegisterChan(chan interface{})
	// UnregisterFunc()
}

func (pl RegisteredPlugin) Notify()

