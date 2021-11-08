package test

import (
	"phoenixbuilder/cq-chatlogger"
	"reflect"
	"testing"
)

var want1 = cqchat.ChatSettings{
	IsFilteredGroup: false,
	FilteredGroupID: []int64{961748506},
	NormalGroupID:   123456789,
	GroupNickname: map[int64]string{},
	//GameMessageFormat: GameMessageFormat{Format: "time-user-flag-message-source", Flag: ""},
	//QQMessageFormat:   QQMessageFormat{Format: "flag-user-message-source", Flag: ""},
}

func TestReadSettings(t *testing.T) {
	type args struct {
		fp string
	}
	tests := []struct {
		name    string
		args    args
		want    cqchat.ChatSettings
		wantErr bool
	}{
		{"模板测试", args{fp: "./settings.yml"}, want1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cqchat.ReadSettings(tt.args.fp)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadSettings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadSettings() got = %v, want %v", got, tt.want)
			}
		})
	}
}