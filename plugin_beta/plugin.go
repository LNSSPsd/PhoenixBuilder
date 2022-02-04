package plugin_beta

import (
	"log"
	"os"
	"path/filepath"
	"phoenixbuilder/minecraft"
	"sync"
)


func StartPluginSystem (conn *minecraft.Conn) *PluginManager{
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
	return &manager
}


func loadPluginDir() (string, error) {
	homedir, err := os.UserHomeDir()
	if err != nil {
		homedir="."
	}
	plugindir := filepath.Join(homedir, ".config/fastbuilder/plugins")
	return plugindir, err
}
