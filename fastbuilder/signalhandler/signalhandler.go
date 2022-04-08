package signalhandler

import (
	"fmt"
	"os"
	"os/signal"
	"phoenixbuilder/fastbuilder/i18n"
	"phoenixbuilder/minecraft"
	"phoenixbuilder/fastbuilder/readline"
	"phoenixbuilder/fastbuilder/args"
	"syscall"
)

func Install(conn *minecraft.Conn) {
	go func() {
		if(args.NoReadline()) {
			return
		}
		readline.SelfTermination=make(chan bool)
		<-readline.SelfTermination
		readline.HardInterrupt()
		conn.Close()
		fmt.Printf("%s.\n",I18n.T(I18n.QuitCorrectly))
		os.Exit(0)
	} ()
	go func() {
		if(args.NoReadline()) {
			return
		}
		for {
			sigintchannel:=make(chan os.Signal)
			signal.Notify(sigintchannel, os.Interrupt) // ^C
			<-sigintchannel
			readline.Interrupt()
		}
	} ()
	go func() {
		signalchannel:=make(chan os.Signal)
		signal.Notify(signalchannel, syscall.SIGTERM)
		signal.Notify(signalchannel, syscall.SIGQUIT) // ^\
		if(args.NoReadline()) {
			signal.Notify(signalchannel, os.Interrupt)
		}
		<-signalchannel
		readline.HardInterrupt()
		conn.Close()
		fmt.Printf("%s.\n",I18n.T(I18n.QuitCorrectly))
		os.Exit(0)
	} ()
}