package plugin_beta

import (
	"fmt"
	"log"
	"sync"
	"github.com/pterm/pterm"
	
)

var IsLogFile bool

type pluginLogger struct{
	logger log.Logger
	prefix string

	Mu sync.Mutex
}

// 看起来像在第二抽象层, 所以这个函数应该是必要的
// 第一层是log.logger本身噢.
// It seems to be the second abstraction level and therefore be written.
// the first level is log.Logger, I think.
func (l *pluginLogger) setPlugin (pluginName string) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	l.prefix = pterm.Yellow(fmt.Sprintf("[%s]", pluginName))
}

func (l *pluginLogger) SPrintln (plugin struct{name string}, v ...interface{}) string {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	l.setPlugin(plugin.name)
	return pterm.Sprintln(l.prefix, l.prefix)
}

func (l *pluginLogger) Println (plugin struct{name string}, v ...interface{}) {
	l.Mu.Lock()
	defer l.Mu.Unlock()
	lg := l.SPrintln(plugin, v...)
	pterm.Println(lg)
	if IsLogFile {
		// todo
	}
}


func New() *pluginLogger{
	logger := log.Logger{
	}
	// flag = 1 + 2 + 16. How it works? See https://pkg.go.dev/log#Logger.Flags.
	logger.SetFlags(19)
	 
	return &pluginLogger{logger: logger}
}