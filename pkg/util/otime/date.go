package otime

import "time"

// GetFirstDateOfWeek 获取传入时间所在周第一天零点时间
func GetFirstDateOfWeek(timeNow time.Time) (weekMonday time.Time) {
	offset := int(time.Monday - timeNow.Weekday())
	if offset > 0 {
		offset = -6
	}
	weekMonday = GetZeroTime(timeNow).AddDate(0, 0, offset)

	return
}

// GetZeroTime 获取传入时间当天零点时间
func GetZeroTime(timeNow time.Time) (zeroTime time.Time) {
	zeroTime = time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, time.Local)

	return
}
