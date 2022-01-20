package cqchat

import (
	"encoding/json"
	"os"
)

type MetaPost struct {
	Time          int64  `json:"time"`
	PostType      string `json:"post_type"`   // 得到消息类型, 进一步解析
	SelfID        int    `json:"self_id"`
	MetaEventType string `json:"meta_event_type"`
	Echo 		  string
}

func ParseMetaPost(data []byte) (MetaPost, error) {
	post := MetaPost{}
	
	err := json.Unmarshal(data, &post)//处理并返回
	return post, err
}

func PathExist(fp string) bool {
	_, err := os.Stat(fp)
	return !(err != nil)
}

// IsFilteredUser 检查用户命令是否合法是否被过滤.
func IsFilteredUser(user int64) bool {
	for _, u := range Setting.FilteredUserID {
		if u == user {
			return true
		}
	}
	return false
}
