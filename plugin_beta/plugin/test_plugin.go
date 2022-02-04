package test_plugin

import (
	"phoenixbuilder/minecraft/protocol/packet"
	conn "phoenixbuilder/plugin_beta"
)

type SingleQABot struct {
	message *packet.Text
	user string
}

func (bot SingleQABot) Rule(pk packet.Packet) bool {
	switch pk.(type) {
	case  *packet.Text:
		return true
	default: 
		return false
	}
}

func (bot *SingleQABot) Init(conn *conn.PluginManager) {
	conn.RegisterPlugin(bot, true, true, 5, "SingleQABot")
}

func (bot *SingleQABot) Handler(conn *conn.PluginManager, pk packet.Packet) {
	bot.message = pk.(*packet.Text)
	
	if bot.message.Message != "留言" {
		return
	}
	bot.user = bot.message.SourceName
	conn.Method.Tellraw("您的留言内容?")
	for {
		bot.message = conn.ReadPacketFor(bot).(*packet.Text)
		if bot.message.SourceName != bot.user {
			continue
		} else {
			conn.Logger.Println(bot.message)
			return
		}
	}
}



var Plugin SingleQABot