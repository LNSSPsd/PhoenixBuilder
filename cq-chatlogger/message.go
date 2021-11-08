package cqchat

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var CQCodeTypes = map[string]string{
	"face":    "表情",
	"record":  "语音",
	"at":      "@某人",
	"share":   "链接分享",
	"music":   "音乐分享",
	"image":   "图片",
	"reply":   "回复",
	"redbag":  "红包",
	"forward": "合并转发",
	"xml":     "XML消息",
	"json":    "json消息",
}

type User struct {
	Nickname string `json:"nickname"`
}

type UniversalMessage struct {
	Message     string `json:"message"`
	GameRawText string
	MessageType string `json:"message_type"`
}

func (msg UniversalMessage) GetMessage() string {
	return msg.Message
}

func (msg UniversalMessage) GetUser() int64 {
	return -1
}

func (msg PrivateMessage) GetUser() int64 {
	return msg.UserId
}


type PrivateMessage struct {
	UniversalMessage
	MetaPost
	UserId int64 `json:"user_id"`
	Sender User  `json:"sender"`
}

type GroupMessage struct {
	PrivateMessage
	GroupID int64 `json:"group_id"`
}

type IMessage interface {
	FormatCQMessage() string
	GetSource() string
	IsCommand() bool
	GetUser() int64
	GetMessage() string
}

func GetMessageData(data []byte) (IMessage, error) {
	msg := map[string]interface{}{}
	err := json.Unmarshal(data, &msg)
	if err != nil{
		return nil, err
	}
	msgType := msg["message_type"].(string)
	switch msgType {
	case "private":
		return PrivateMessage{}.Unmarshal(data)
	case "group":
		return GroupMessage{}.Unmarshal(data)
	default:
		return UniversalMessage{}.Unmarshal(data)
	}
}

func (msg UniversalMessage) Unmarshal(data []byte) (IMessage, error) {
	err := json.Unmarshal(data, &msg)
	return msg, err
}

// FormatCQMessage 按配置文件格式化消息.
func (msg UniversalMessage) FormatCQMessage() string {
	raw := Setting.GameMessageFormat
	raw = strings.ReplaceAll(raw, "source", msg.GetSource())
	raw = strings.ReplaceAll(raw, "message", GetRawTextFromCQMessage(msg.Message))
	raw = strings.ReplaceAll(raw, "type", strings.ToUpper(msg.MessageType))
	return raw
}


func (msg PrivateMessage) FormatCQMessage() string {
	raw := msg.UniversalMessage.FormatCQMessage()
	raw = strings.ReplaceAll(raw, "time", time.Unix(msg.Time, 0).Format("15:04:05"))
	raw = strings.ReplaceAll(raw, "user", msg.Sender.Nickname)
	return raw
}


// GetSource 返回当前信息的来源. source为在group_id_list中定义的群昵称. 如果没有定义 则以群号代替. 若为私聊消息, 则为空值.
func (msg UniversalMessage) GetSource() string {
	return ""
}

func (msg GroupMessage) GetSource() string {
	for id, title := range Setting.GroupNickname {
		if msg.GroupID == id {
			return title
		}
	}
	return strconv.FormatInt(msg.GroupID, 10)
}

// GetRawTextFromCQMessage 将图片等CQ码转为文字.
func GetRawTextFromCQMessage(msg string) string {
	for k, v := range CQCodeTypes {
		format := fmt.Sprintf(`\[CQ:%s.*?\]`, k)
		rule := regexp.MustCompile(format)
		msg = rule.ReplaceAllString(msg, fmt.Sprintf("[%s]", v))
	}
	return msg
}

// IsCommand 判断消息是否为游戏内命令
func (msg UniversalMessage) IsCommand() bool {
	if !strings.HasPrefix(msg.Message, Setting.CommandPrefix) || Setting.CommandPrefix == "" {
		return false
	}
	return true
}

// TellrawCommand 将消息转为tellraw命令
func TellrawCommand(msg string) string {
	tag := Setting.FilteredScbTitle
	msg = strings.ReplaceAll(msg, `\`, `\\`)
	msg = strings.ReplaceAll(msg, `"`, `\"`)
	cmd :=fmt.Sprintf(`tellraw @a[tag=!%s] {"rawtext":[{"text": "%s"}]}`, tag, msg)
	return cmd
}