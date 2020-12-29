package utils

import "time"

// 当前时间戳
func NowTimestamp() int64 {
	return time.Now().Unix()
}

// 时间戳转换为DateTime
func FormatDateTime(sec int64, layout string) string {
	t := time.Unix(sec, 0)
	return t.Format(layout)
}

// DateTime转换为时间戳
func FormatTimestamp(date, layout string) int64 {
	t, _ := time.ParseInLocation(layout, date, time.Local)
	return t.Unix()
}

// 返回起止时间列表
func SubUnixTimeDays(start, end int64, layout string) []string {
	st := time.Unix(start, 0)
	et := time.Unix(end, 0)
	dt := et.Sub(st)
	ds := int64(dt.Hours())/24 + 1

	days := make([]string, 0, ds)
	for i := int64(0); i < ds; i++ {
		t := time.Unix(start+i*86400, 0)
		days = append(days, t.Format(layout))
	}

	return days
}
