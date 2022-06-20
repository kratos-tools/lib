package otime

import (
	"reflect"
	"testing"
	"time"
)

func TestGetFirstDateOfWeekWeekBegin(t *testing.T) {
	type args struct {
		timeNow time.Time
	}
	timeNow, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-04-26 10:00:00", time.Local)
	tests := []struct {
		name           string
		args           args
		wantWeekMonday string
	}{
		{
			name:           "获取传入时间所在周第一天零点时间(周一)",
			args:           args{timeNow},
			wantWeekMonday: "2021-04-26",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWeekMonday := GetFirstDateOfWeek(tt.args.timeNow).Format("2006-01-02"); !reflect.DeepEqual(gotWeekMonday, tt.wantWeekMonday) {
				t.Errorf("GetFirstDateOfWeek() = %v, want %v", gotWeekMonday, tt.wantWeekMonday)
			}
		})
	}
}

func TestGetFirstDateOfWeekWeek(t *testing.T) {
	type args struct {
		timeNow time.Time
	}
	timeNow, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-05-01 00:00:00", time.Local)
	tests := []struct {
		name           string
		args           args
		wantWeekMonday string
	}{
		{
			name:           "获取传入时间所在周第一天零点时间(周中)",
			args:           args{timeNow},
			wantWeekMonday: "2021-04-26",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWeekMonday := GetFirstDateOfWeek(tt.args.timeNow).Format("2006-01-02"); !reflect.DeepEqual(gotWeekMonday, tt.wantWeekMonday) {
				t.Errorf("GetFirstDateOfWeek() = %v, want %v", gotWeekMonday, tt.wantWeekMonday)
			}
		})
	}
}

func TestGetFirstDateOfWeekWeekEnd(t *testing.T) {
	type args struct {
		timeNow time.Time
	}
	timeNow, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-05-02 00:00:00", time.Local)
	tests := []struct {
		name           string
		args           args
		wantWeekMonday string
	}{
		{
			name:           "获取传入时间所在周第一天零点时间(周日)",
			args:           args{timeNow},
			wantWeekMonday: "2021-04-26",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotWeekMonday := GetFirstDateOfWeek(tt.args.timeNow).Format("2006-01-02"); !reflect.DeepEqual(gotWeekMonday, tt.wantWeekMonday) {
				t.Errorf("GetFirstDateOfWeek() = %v, want %v", gotWeekMonday, tt.wantWeekMonday)
			}
		})
	}
}

func TestGetZeroTime(t *testing.T) {
	type args struct {
		timeNow time.Time
	}
	timeNow, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-04-26 10:00:00", time.Local)
	tests := []struct {
		name         string
		args         args
		wantZeroTime string
	}{
		{
			name:         "获取传入时间所在周第一天零点时间(周日)",
			args:         args{timeNow},
			wantZeroTime: "2021-04-26",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotZeroTime := GetZeroTime(tt.args.timeNow).Format("2006-01-02"); !reflect.DeepEqual(gotZeroTime, tt.wantZeroTime) {
				t.Errorf("GetZeroTime() = %v, want %v", gotZeroTime, tt.wantZeroTime)
			}
		})
	}
}
