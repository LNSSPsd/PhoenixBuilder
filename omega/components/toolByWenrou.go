package components

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"phoenixbuilder/omega/collaborate"
	"phoenixbuilder/omega/defines"
	"strings"
	"time"
)

// 获取fb用户名字(废弃)
func GetFbUserNameByToken(fbtoken string) string {

	url := "https://uc.fastbuilder.pro/api/v2/login.web"
	token := fbtoken
	fmt.Println(fbtoken)
	//json序列化
	post := "{\"token\":\"" + token +
		"\"}"
	var jsonStr = []byte(post)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("抱歉获取时间超时")
		return ""
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var fbname ToGetFbName
	fmt.Println("body:", string(body))
	json.Unmarshal(body, &fbname)
	return string(fbname.Name)
	//fmt.Println("response Body:", )

}
func GetYsCoreNameList() map[string]string {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	resp, err := client.Get("https://pans-1259150973.cos-website.ap-shanghai.myqcloud.com")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("getYsCoreName:", string(body))
	arr := strings.Split(string(body), " ")
	list := make(map[string]string)
	for _, v := range arr {
		list[v] = ""
	}
	return list
}

// 监听fb的用户名称
func ListenFbUserName(Maxnum int, b defines.MainFrame) string {
	num := 0
	for {
		if num >= Maxnum {
			break
		}
		time.Sleep(time.Second * 1)
		if username, ok := (*b.GetContext())[collaborate.INTERFACE_FB_USERNAME]; ok && username != "" {
			name := username.(collaborate.STRING_FB_USERNAME)
			//fmt.Println("fb用户名取得:", name)
			return string(name)
		}
		num++
	}
	return ""

}
func ListenFbUserNamer(Maxnum int, b defines.MainFrame, componentsName string) {
	name := ListenFbUserName(10, b)
	if name == "" {
		panic("[抱歉] " + componentsName + "组件为yscore会员组件 获取你的fb用户名字时链接超时或者你没有开启YsCore验证组件 请重复尝试如果多次尝试依旧报错请关闭所有的yscore出品组件")
	}
}
