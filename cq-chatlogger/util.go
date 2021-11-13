package cqchat

import (
	"encoding/json"
	"os"
)

type MetaPost struct {
	Time          int64  `json:"time"`
	PostType      string `json:"post_type"`
	SelfID        int    `json:"self_id"`
	MetaEventType string `json:"meta_event_type"`
}

func ParseMetaPost(data []byte) (MetaPost, error) {
	post := MetaPost{}
	err := json.Unmarshal(data, &post)
	return post, err
}

func PathExist(fp string) bool {
	_, err := os.Stat(fp)
	return !(err != nil)
}
