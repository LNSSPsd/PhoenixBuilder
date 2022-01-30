package plugin_beta

type mainFunction interface{} // Return true if it blocks packets.
// e.g. Main(logger ) 

type RegisteredPlugin struct {
	Priority uint64
	MainFunction mainFunction
	Logger *pluginLogger

}

type PluginSystem struct {	
	Plugins []RegisteredPlugin

}

func (pl *PluginSystem) LoadPlugin()

func (pl *PluginSystem) LoadPluginFrom()

func (pl *PluginSystem) RegisterFunction()

func (pl *PluginSystem) UnregisterFunction()
