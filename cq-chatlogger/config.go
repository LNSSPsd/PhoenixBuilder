package cqchat

import (
	"fmt"
	"github.com/pterm/pterm"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

func writeConf() []byte {
	return []byte(`
# cq-chatlogger 配置

#
address: "127.0.0.1:5555"


# 是否过滤消息. 如果为真,则过滤 filtered_group_id 中的群消息(黑名单), 反之则仅接收群消息(白名单)
is_filtered_group: false

filtered_group_id: [
#  1098232840,
]

# 游戏内的消息将默认转发至哪个群. 如果为空, 则默认不转发, 只能通过指定别名发送指定群消息.
normal_group_id: 123456789

# 给每个群设置别名来指定群聊发送消息.
# 例如按照如下配置, 在游戏中发送此消息是合法的:
# FBP: alpha版什么时候可以插件化啊kuso!
group_nickname: {
#  1098232840: FBP,
#  961748506: MR,
}


# qq消息转发至游戏的消息格式.
# time: 消息时间.
# message: 消息主体. 其中表情、 图片等消息将转化为 [表情] [图片] 等纯文字形式.
# source: 在group_id_list中定义的群昵称. 如果没有定义 则以群号代替. 若为私聊消息, 则为空值.
# type: 消息类型. 默认有 private 和 group . 可以自定义协议端并传入其他类型. 游戏中将以大写字母呈现.
# 参数可以重复, 可以省略, 也可以加一点括号或颜色符号§之类的.
# 在游戏中仍然会受到屏蔽词影响.

# 几个示例:
# <user> message
# <香音> 你好谢谢小笼包再见

# [type] user: message
# [GROUP] 达达鸭: 破绽 烧冻鸡翅!

# [time] §r user: message (source)
# [12:33:04]  菜月昴: EMT Maji Tenshi! (FBP)
game_message_format: "[time] [user] message (source)"


# 游戏聊天转发至qq的消息格式. user为游戏ID, source为租赁服号. 不建议使用time参数(因为没啥必要).
qq_message_format: "[user] message [source]"

filtered_scb_title: "filtered"
# 定义一个计分板, 以玩家分数表示在该插件中的权限
# 玩家分数为: (如果不是使用自定义协议端, 则以下协议端默认为qq.)
# 1: 可以收到协议端消息, 但消息不会被转发至协议端.
# 2: 可以收到协议端消息或将消息转发至协议端.
# 其他(或没有值): 不能收发协议端消息


# 在qq中使用命令: 选择一个前缀,来标识它是一个命令. 不要选择空字符, 它将永远无法生效.
command_prefix: "/"

# 是否选择过滤用户. 如果为真, 则仅允许 command_user 中的用户在qq中使用命令(白名单),反之仅阻止.
is_filtered_user: false

# 哪些用户可以(或不可以)在qq中使用命令. 填入qq号. 示例配置:
# command_user: [123456789, 987654321]
filtered_user_id: []
`)
}

type ChatSettings struct {
	Port              string           `yaml:"address"`
	NormalGroupID     int64            `yaml:"normal_group_id"`
	GroupNickname     map[int64]string `yaml:"group_nickname"`
	GameMessageFormat string           `yaml:"game_message_format"`
	QQMessageFormat   string           `yaml:"qq_message_format"`
	FilteredPlayerTag string           `yaml:"filtered_player_tag"`
	CommandPrefix     string           `yaml:"command_prefix"`
	FilteredUserID    []int64          `yaml:"filtered_user_id"`
}

var Setting ChatSettings
var ErrSetting error

func ReadSettings(fp string) (ChatSettings, error) {
	f, err := ioutil.ReadFile(fp)
	out := ChatSettings{}
	err = yaml.Unmarshal(f, &out)
	return out, err
}

func init() {
	fp := "./cq-chatlogger"
	conf := fp + "config.yml"
	if !PathExist(fp) {
		os.Mkdir(fp, os.ModePerm)
	}
	if !PathExist(conf) {
		ioutil.WriteFile(conf, writeConf(), os.ModePerm)
		fmt.Println("chatlogger配置文件已创建. 配置后下次启动生效.")
	}
	Setting, ErrSetting = ReadSettings("./cq-chatlogger/config.yml")
	if ErrSetting != nil {
		pterm.Println(pterm.Red("WARNING: config.yml解析异常:", ErrSetting))
	}
	return
}
