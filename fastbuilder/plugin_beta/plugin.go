package plugin_beta

import (
	"log"
	"os"
	"path/filepath"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/minecraft/protocol/packet"
)


func StartPluginSystem (conn *minecraft.Conn) {
	manager := PluginManager {
		conn: conn,
		Logger: &log.Logger{},
		PacketReceiver: make(chan packet.Packet),
	}
	manager.Logger.SetPrefix("[PLUGIN]")
	manager.loadPlugin()
	manager.notify()
}


func loadPluginDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		homedir="."
	}
	plugindir := filepath.Join(homedir, ".config/fastbuilder/plugins")
	return plugindir, err
}
