package utils

import (
	"database/sql/driver"
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	DaySeconds  = 24 * 60 * 60
	DayMilli    = 24 * 60 * 60 * 1000
	HourMilli   = 60 * 60 * 1000
	MinuteMilli = 60 * 1000
	SecMilli    = 1000
)

type LocalTime time.Time

func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	tlt := time.Time(t)
	if tlt.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return tlt, nil
}

func (t *LocalTime) Scan(v interface{}) error {
	if value, ok := v.(time.Time); ok {
		*t = LocalTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t *LocalTime) String() string {
	// 如果时间 null 那么我们需要把返回的值进行修改
	if t == nil || t.IsZero() {
		return ""
	}
	return fmt.Sprintf("%s", time.Time(*t).Format("2006-01-02 15:04:05"))
}

func (t *LocalTime) IsZero() bool {
	return time.Time(*t).IsZero()
}

func (t *LocalTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	var err error
	//前端接收的时间字符串
	str := string(data)
	//去除接收的str收尾多余的"
	timeStr := strings.Trim(str, "\"")
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", timeStr, time.Local)
	*t = LocalTime(t1)
	return err
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	tTime := time.Time(t)
	// 如果时间值是空或者0值 返回为null 如果写空字符串会报错
	if &t == nil || t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", tTime.Format("2006-01-02 15:04:05"))), nil
}

func GetLastMonthStartEnd() (lastMonthStart, lastMonthEnd time.Time) {
	now := time.Now()
	lastMonthFirstDay := now.AddDate(0, -1, -now.Day()+1)
	lastMonthStart = time.Date(lastMonthFirstDay.Year(), lastMonthFirstDay.Month(), lastMonthFirstDay.Day(), 0, 0, 0, 0, now.Location())
	lastMonthEndDay := lastMonthFirstDay.AddDate(0, 1, -1)
	lastMonthEnd = time.Date(lastMonthEndDay.Year(), lastMonthEndDay.Month(), lastMonthEndDay.Day(), 23, 59, 59, 1e9-1, now.Location())

	return
}

func DiffDays(t1, t2 int64) uint32 {
	return uint32(math.Abs(float64((t1+28800)/DaySeconds - (t2+28800)/DaySeconds)))
}

func FormatResidue(t int64) string {
	var str = ""
	if t >= DayMilli {
		var day = t / DayMilli
		str = fmt.Sprintf("%d天", day)
	}
	if t >= HourMilli {
		var hour = t / HourMilli % 24
		str = fmt.Sprintf("%s%d时", str, hour)
	}
	if t >= MinuteMilli {
		var minu = t / HourMilli % 60
		str = fmt.Sprintf("%s%d分", str, minu)
	}
	if t >= SecMilli {
		var sec = t / SecMilli % 60
		str = fmt.Sprintf("%s%d秒", str, sec)
	}
	return str
}

// GetTimes 获取时间
//
//	@return startTime 时间对应的0点
//	@return endTime 时间对应的23:59:59
//	@return weekAt 时间对应的周日0点
//	@return monthAt 时间对应的当月1号0点
func GetTimes(t time.Time) (startTime, endTime, weekAt, monthAt time.Time) {
	startTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	endTime = time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 1e9-1, time.Local)

	weekAt = startTime.AddDate(0, 0, -int(startTime.Weekday()))
	weekAt = time.Date(weekAt.Year(), weekAt.Month(), weekAt.Day(), 0, 0, 0, 0, time.Local)
	monthAt = time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.Local)
	return
}
