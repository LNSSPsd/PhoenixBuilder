package test

import (
	"phoenixbuilder/cq-chatlogger"
	"testing"
)



func TestTellrawCommand(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		//{"t0", args{"hello world"}, "hello world"},
		{"t0", args{`"Emi\\lia"`}, `tellraw @a[tag=!msg_listener] {"rawtext":[{"text": "\"Emi\\\\lia\""}]}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cqchat.TellrawCommand(tt.args.msg); got != tt.want {
				t.Errorf("TellrawCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}


func Test_getRawTextFromCQMessage(t *testing.T) {
	type args struct {
		msg string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"", args{"[CQ:face,id=39]"}, "[表情]"},
		{"", args{"[CQ:record,resource:balabala]一段录音"},"[语音]一段录音"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cqchat.GetRawTextFromCQMessage(tt.args.msg); got != tt.want {
				t.Errorf("getRawTextFromCQMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}