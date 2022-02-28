package main

import (
	"RPC/channel"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type MajorConfig struct {
}

type Major struct {
	config         *MajorConfig
	PageServes     map[string]func(w http.ResponseWriter, r *http.Request)
	gettableStatus map[string]string
}

func InitMajorFunction(config *MajorConfig) *Major {
	fn := &Major{
		config:         config,
		PageServes:     map[string]func(w http.ResponseWriter, r *http.Request){},
		gettableStatus: map[string]string{},
	}
	fn.PageServes["/status/"] = func(w http.ResponseWriter, r *http.Request) {
		p := strings.TrimPrefix(r.URL.Path, "/status/")
		s, ok := fn.gettableStatus[p]
		if ok {
			io.WriteString(w, s)
		} else {
			io.WriteString(w, "NoSuchStatus")
		}
	}
	return fn
}

func (mf *Major) Connect(connMux *channel.Mux, closeConn func()) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		log.Println(r)
	}()
	channel3 := connMux.GetSubChannel(3)
	err := channel3.Send([]byte("hello"))
	if err != nil {
		panic(err)
	}
	msgGet := channel3.Get()
	fmt.Println(string(msgGet))
}
