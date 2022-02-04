package plugin_beta

import (
	"log"
	"os"
	"path/filepath"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol/packet"
	"sync"
)


func StartPluginSystem (conn *minecraft.Conn) chan packet.Packet{
	receiver := make(chan packet.Packet)
	manager := PluginManager {
		conn: conn,
		Logger: &log.Logger{},
		regMu: sync.RWMutex{},
		pluginPriority: []IPlugin{},
		plugins: map[IPlugin]*Plugin{},
	}
	manager.Logger.SetPrefix("[PLUGIN]")
	err := manager.loadPlugins()
	if err != nil {
		manager.Logger.Println("Plugin system crashed")
	}
	go func ()  {
		manager.Notify(<-receiver)
	}()
	return receiver

}


func loadPluginDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		homedir="."
	}
	plugindir := filepath.Join(homedir, ".config/fastbuilder/plugins")
	return plugindir, err
}
