package plugin_beta

import (
	 
	"os"
	"path/filepath"
	"fmt"
)



func StartPluginSystem () {
	logger := pluginLogger{

	}
	plugins:=loadConfigPath()
	files, _ := ioutil.ReadDir(plugins)
	pluginbridge := plugin_structs.PluginBridge(&PluginBridgeImpl {
		sessionConnection: conn,
	})
	for _, file := range files {
		path:=fmt.Sprintf("%s/%s",plugins,file.Name())
		if filepath.Ext(path)!=".so" {
			continue
		}
		go func() {
			RunPlugin(conn,path,pluginbridge)
		} ()
	}
}

func loadConfigPath() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("[PLUGIN] WARNING - Failed to obtain the user's home directory. made homedir=\".\";\n")
		homedir="."
	}
	fbconfigdir := filepath.Join(homedir, ".config/fastbuilder/plugins")
	os.MkdirAll(fbconfigdir, 0755)
	return fbconfigdir
}