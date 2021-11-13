package cqchat

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Player map[string]int32

type scoreboards struct {
	Scoreboard map[string][]Player
}

var UserInfo scoreboards

var filename = "./cq-chatlogger/scoreboard.json"

func init() {
	f, err := os.Open(filename)
	if err != nil {
		f, _ = os.Create(filename)
		_, _ = f.Write([]byte("{}"))
	}
	data, _ := ioutil.ReadAll(f)
	_ = f.Close()
	err = json.Unmarshal(data, &UserInfo)
	if err != nil {
		log.Println("解析scoreboard文件出错.")
		return
	}
	return
}

func UpdateScoreboardInfo() error {
	data, err := json.Marshal(UserInfo)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0777)
	return err
}
