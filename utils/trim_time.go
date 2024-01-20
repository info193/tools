package utils

import (
	"errors"
	"fmt"
	"time"
)

type TrimTime struct {
	Duration               int64  // 卡券时长
	Neutron                int64  // 边界时间
	PeriodStartHour        string // 开始时段 00:00
	PeriodEndHour          string // 结束时段 05:00
	StartTime              time.Time
	EndTime                time.Time
	tempSubscribeStartTime time.Time
}

type TimeSpan struct {
	OffsetStartDate          string // 抵消开始时间
	OffsetEndDate            string // 抵消结束时间
	ExceedStartDate          string // 超出开始时间
	ExceedEndDate            string // 超出结束时间
	AvailableStartPeriodHour string // 可用开始时段
	AvailableEndPeriodHour   string // 可用结束时段
	ExceedDuration           int64  // 超出时长
	MaxLimitDuration         int64  // 最大限制超出时长
	Type                     int64  // 类型 1隔日（正常可用范围） 2隔日（在预约开始结束时间及时段结束时间范围内） 3隔日（预约开始时间大于时段结束时间范围内）4当日
	IsExceed                 int64  // 是否超出 0 否 1是  2超出可用范围  3不在使用时段范围内
}

func NewTrimTime(duration int64, periodStartHour, periodEndHour, subscribeStartDate, subscribeEndDate string) TrimTime {
	subscribeStartTime, _ := time.ParseInLocation("2006-01-02 15:04", subscribeStartDate, time.Local)
	subscribeEndTime, _ := time.ParseInLocation("2006-01-02 15:04", subscribeEndDate, time.Local)
	trimTime := TrimTime{
		Duration:               duration,
		PeriodStartHour:        periodStartHour,
		PeriodEndHour:          periodEndHour,
		StartTime:              subscribeStartTime,
		EndTime:                subscribeEndTime,
		tempSubscribeStartTime: subscribeStartTime}
	trimTime.setNeutron()
	return trimTime
}

// 设置
func (l *TrimTime) setNeutron() {
	switch {
	case l.Duration <= 10:
		l.Neutron = 1
	case l.Duration > 10 && l.Duration < 30:
		l.Neutron = 5
	case l.Duration >= 30 && l.Duration < 60:
		l.Neutron = 10
	case l.Duration >= 60 && l.Duration < 120:
		l.Neutron = 30
	case l.Duration >= 120 && l.Duration < 180:
		l.Neutron = 60
	case l.Duration >= 180 && l.Duration < 240:
		l.Neutron = 60
	case l.Duration >= 240 && l.Duration < 300:
		l.Neutron = 120
	case l.Duration >= 300 && l.Duration < 360:
		l.Neutron = 120
	case l.Duration >= 360 && l.Duration < 420:
		l.Neutron = 180
	case l.Duration >= 420 && l.Duration < 480:
		l.Neutron = 180
	case l.Duration >= 480 && l.Duration < 600:
		l.Neutron = 240
	case l.Duration >= 600:
		l.Neutron = 300
	}
}

// 开始时段大于结束时段，同一天时间
func (l *TrimTime) spanTime() (error, *TimeSpan) {
	var timeSpan TimeSpan
	if l.StartTime.Format("2006-01-02") == l.EndTime.Format("2006-01-02") && l.PeriodStartHour > l.PeriodEndHour {
		tsPeriodStartTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodStartHour), time.Local)
		tsPeriodEndTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), "23:59"), time.Local)
		tePeriodStartTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), "00:00"), time.Local)
		tePeriodEndTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodEndHour), time.Local)

		timeSpan.AvailableStartPeriodHour = l.PeriodStartHour // 可用开始时段
		timeSpan.AvailableEndPeriodHour = l.PeriodEndHour     // 可用结束时段

		if l.tempSubscribeStartTime.Unix() >= tsPeriodStartTime.Unix() && l.tempSubscribeStartTime.Unix() <= tsPeriodEndTime.Unix() {
			//fmt.Println("大于范围内")
			timeSpan.OffsetStartDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")

			l.tempSubscribeStartTime = l.tempSubscribeStartTime.Add(time.Minute * time.Duration(l.Duration))
			//fmt.Println(tempSubscribeStartTime.Format("2006-01-02 15:04"), tsPeriodEndTime.Format("2006-01-02 15:04"), subscribeEndTime.Format("2006-01-02 15:04"))

			exceedSecond := l.tempSubscribeStartTime.Unix() - tsPeriodEndTime.Unix()
			if l.tempSubscribeStartTime.After(tsPeriodEndTime) {
				timeSpan.OffsetEndDate = tsPeriodEndTime.Format("2006-01-02 15:04")
				timeSpan.ExceedStartDate = tsPeriodEndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
				timeSpan.ExceedEndDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
				timeSpan.ExceedDuration = exceedSecond
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.IsExceed = 1
			} else {
				timeSpan.OffsetEndDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.IsExceed = 0
			}

			if exceedSecond > l.Neutron*60 {
				timeSpan.IsExceed = 2
				//fmt.Println(fmt.Sprintf("错误，超出可用范围，已超出%v分钟,最大超出限制:%v分钟,值：%v", exceedSecond/60, Neutron, exceedSecond), tempSubscribeStartTime.Format("2006-01-02 15:04"), subscribeEndTime.Format("2006-01-02 15:04"))
			}
			//fmt.Println("option===", option)
			return nil, &timeSpan
		}
		if l.tempSubscribeStartTime.Unix() >= tePeriodStartTime.Unix() && l.tempSubscribeStartTime.Unix() <= tePeriodEndTime.Unix() {
			//fmt.Println("小于 范围内")
			timeSpan.OffsetStartDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
			l.tempSubscribeStartTime = l.tempSubscribeStartTime.Add(time.Minute * time.Duration(l.Duration))
			//fmt.Println(tempSubscribeStartTime.Format("2006-01-02 15:04"), tePeriodEndTime.Format("2006-01-02 15:04"), tePeriodEndTime.Format("2006-01-02 15:04"))
			exceedSecond := l.tempSubscribeStartTime.Unix() - tePeriodEndTime.Unix()
			if l.tempSubscribeStartTime.After(tePeriodEndTime) {
				timeSpan.OffsetEndDate = tePeriodEndTime.Format("2006-01-02 15:04")
				timeSpan.ExceedStartDate = tePeriodEndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
				timeSpan.ExceedEndDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
				timeSpan.ExceedDuration = exceedSecond
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.IsExceed = 1
			} else {
				timeSpan.OffsetEndDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.IsExceed = 0
			}
			if exceedSecond > l.Neutron*60 {
				timeSpan.IsExceed = 2
				//fmt.Println(fmt.Sprintf("错误，超出可用范围，已超出%v分钟,最大超出限制:%v分钟,值：%v", exceedSecond/60, Neutron, exceedSecond), tempSubscribeStartTime.Format("2006-01-02 15:04"), tePeriodEndTime.Format("2006-01-02 15:04"))
			}
			return nil, &timeSpan
		}
		if l.tempSubscribeStartTime.Unix() < tsPeriodStartTime.Unix() && l.tempSubscribeStartTime.Unix() > tePeriodEndTime.Unix() {
			//fmt.Println("不在 范围内")
			timeSpan.OffsetStartDate = tsPeriodStartTime.Format("2006-01-02 15:04")
			l.tempSubscribeStartTime = tsPeriodStartTime.Add(time.Minute * time.Duration(l.Duration))
			//fmt.Println(tempSubscribeStartTime.Format("2006-01-02 15:04"), tsPeriodStartTime.Format("2006-01-02 15:04"), tsPeriodEndTime.Format("2006-01-02 15:04"))
			//
			//fmt.Println(tempSubscribeStartTime.Format("2006-01-02 15:04"), subscribeEndTime.Format("2006-01-02 15:04"), "..................")
			// 预约开始时间小于等于预约结束时间
			if l.tempSubscribeStartTime.Unix() <= l.EndTime.Unix() {
				timeSpan.OffsetStartDate = tsPeriodStartTime.Format("2006-01-02 15:04")
				timeSpan.OffsetEndDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.IsExceed = 0
				//fmt.Println("不在范围1：", option)
				return nil, &timeSpan
			}
			// 预约开始时间大于预约开始结束时间 ，且小于时段结束时间
			if l.tempSubscribeStartTime.Unix() > l.EndTime.Unix() && l.tempSubscribeStartTime.Unix() < tsPeriodEndTime.Unix() {
				exceedSecond := l.tempSubscribeStartTime.Unix() - l.EndTime.Unix()
				timeSpan.OffsetStartDate = tsPeriodStartTime.Format("2006-01-02 15:04")
				timeSpan.OffsetEndDate = l.EndTime.Format("2006-01-02 15:04")
				timeSpan.ExceedStartDate = l.EndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
				timeSpan.ExceedEndDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
				timeSpan.ExceedDuration = exceedSecond
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.IsExceed = 0
				if exceedSecond > l.Neutron*60 {
					timeSpan.IsExceed = 2
					//fmt.Println(fmt.Sprintf("错误，超出可用范围，已超出%v分钟,最大超出限制:%v分钟,值：%v", exceedSecond/60, Neutron, exceedSecond), tempSubscribeStartTime.Format("2006-01-02 15:04"), subscribeEndTime.Format("2006-01-02 15:04"))
				}
				//fmt.Println("不在范围2：", exceedSecond, option, tsPeriodEndTime.Format("2006-01-02 15:04"))
				return nil, &timeSpan
			}

			// 预约开始时间大于预约开始结束时间， 且大于时段结束时间。
			if l.tempSubscribeStartTime.Unix() > l.EndTime.Unix() && l.tempSubscribeStartTime.Unix() > tsPeriodEndTime.Unix() {
				exceedSecond := l.tempSubscribeStartTime.Unix() - tsPeriodEndTime.Unix()
				timeSpan.OffsetStartDate = tsPeriodStartTime.Format("2006-01-02 15:04")
				timeSpan.OffsetEndDate = tsPeriodEndTime.Format("2006-01-02 15:04")
				timeSpan.ExceedStartDate = tsPeriodEndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
				timeSpan.ExceedEndDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
				timeSpan.ExceedDuration = exceedSecond
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.IsExceed = 0
				if exceedSecond > l.Neutron*60 {
					timeSpan.IsExceed = 2
					//fmt.Println(fmt.Sprintf("错误，超出可用范围，已超出%v分钟,最大超出限制:%v分钟,值：%v", exceedSecond/60, Neutron, exceedSecond), tempSubscribeStartTime.Format("2006-01-02 15:04"), subscribeEndTime.Format("2006-01-02 15:04"))
				}
				//fmt.Println("不在范围3：", option)
				return nil, &timeSpan
			}
		}

		// 默认禁止使用
		timeSpan.OffsetStartDate = ""
		timeSpan.OffsetEndDate = ""
		timeSpan.ExceedStartDate = l.StartTime.Format("2006-01-02 15:04")
		timeSpan.ExceedEndDate = l.EndTime.Format("2006-01-02 15:04")
		timeSpan.ExceedDuration = 0
		timeSpan.MaxLimitDuration = l.Neutron
		timeSpan.IsExceed = 2
		return nil, &timeSpan
		//fmt.Println("正常时间内", subscribeStartTime.Format("2006-01-02 15:04"), subscribeEndTime.Format("2006-01-02 15:04"))
	}
	return errors.New("错误"), nil
}

// 获取可用时段 extractPeriod
func (l *TrimTime) extractPeriod() (error, *TimeSpan) {
	var periodStartTime time.Time
	var periodEndTime time.Time
	var resultEndTime time.Time
	var lastTime time.Time
	currentTime := l.tempSubscribeStartTime
	lists := make([]TimeSpan, 0)
	var dayi int64
	for {
		timeSpan := TimeSpan{}
		timeSpan.AvailableStartPeriodHour = l.PeriodStartHour // 可用开始时段
		timeSpan.AvailableEndPeriodHour = l.PeriodEndHour     // 可用结束时段
		lastTime = l.tempSubscribeStartTime
		if l.tempSubscribeStartTime.Format("2006-01-02") > l.EndTime.Format("2006-01-02") {
			//fmt.Println("跳出", l.tempSubscribeStartTime.Format("2006-01-02 15:04"))
			break
		}
		// 隔日
		if l.PeriodStartHour > l.PeriodEndHour {
			periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodStartHour), time.Local)
			periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Add(time.Second*time.Duration(86400)).Format("2006-01-02"), l.PeriodEndHour), time.Local)
			timeSpan.OffsetStartDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
			if !l.tempSubscribeStartTime.After(periodStartTime) {
				l.tempSubscribeStartTime = periodStartTime
				timeSpan.OffsetStartDate = periodStartTime.Format("2006-01-02 15:04")
			}
			resultEndTime = l.tempSubscribeStartTime.Add(time.Minute * time.Duration(l.Duration))
			timeSpan.OffsetEndDate = resultEndTime.Format("2006-01-02 15:04")

			// 预约开始时间+优惠时段 大于 预约结束时间
			if resultEndTime.Unix() > l.EndTime.Unix() && resultEndTime.Unix() < periodEndTime.Unix() {
				timeSpan.ExceedStartDate = l.EndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
				timeSpan.ExceedEndDate = resultEndTime.Format("2006-01-02 15:04")
				exceedSecond := resultEndTime.Unix() - l.EndTime.Unix()
				timeSpan.ExceedDuration = exceedSecond
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.Type = 2
				timeSpan.IsExceed = 0
				if exceedSecond > l.Neutron*60 {
					timeSpan.IsExceed = 2
					lists = append(lists, timeSpan)
					//fmt.Println(fmt.Sprintf("错误1，超出可用范围，已超出%v分钟,最大超出限制:%v分钟,值：%v", exceedSecond/60, Neutron, exceedSecond), resultEndTime.Format("2006-01-02 15:04"), periodEndTime.Format("2006-01-02 15:04"))
					//fmt.Println("错误1", timeSpan)
					//return
				} else {
					timeSpan.IsExceed = 1
					lists = append(lists, timeSpan)
					//// 正常使用
					//break
				}
				//fmt.Println("超出，但是可以使用，在可用范围内1", timeSpan)
			}

			// 预约开始时间+优惠时段  大于 时段结束时间
			if resultEndTime.Unix() >= periodEndTime.Unix() {
				timeSpan.OffsetEndDate = periodEndTime.Format("2006-01-02 15:04")
				timeSpan.ExceedStartDate = periodEndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
				timeSpan.ExceedEndDate = resultEndTime.Format("2006-01-02 15:04")
				exceedSecond := resultEndTime.Unix() - periodEndTime.Unix()
				timeSpan.ExceedDuration = exceedSecond
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.Type = 3
				timeSpan.IsExceed = 0
				if exceedSecond > l.Neutron*60 {
					timeSpan.IsExceed = 2
					lists = append(lists, timeSpan)
					//fmt.Println(fmt.Sprintf("错误2，超出可用范围，已超出%v分钟,最大超出限制:%v分钟,值：%v", exceedSecond/60, Neutron, exceedSecond), resultEndTime.Format("2006-01-02 15:04"), periodEndTime.Format("2006-01-02 15:04"))
					//fmt.Println("错误2", timeSpan)
					//return
				} else {
					timeSpan.IsExceed = 1
					lists = append(lists, timeSpan)
					//// 正常使用
					//break
				}
				//fmt.Println("超出，但是可以使用，在可用范围内2", timeSpan)
			}

			if resultEndTime.Unix() < l.EndTime.Unix() && resultEndTime.Unix() < periodEndTime.Unix() {
				timeSpan.ExceedStartDate = ""
				timeSpan.ExceedEndDate = ""
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.Type = 1
				timeSpan.IsExceed = 0
				lists = append(lists, timeSpan)
				// 正常使用
				break
				//fmt.Println("正常为超出在可用范围内", timeSpan)
			}
		}
		// 当日
		if l.PeriodEndHour > l.PeriodStartHour {
			//fmt.Println("当日")
			periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodStartHour), time.Local)
			periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodEndHour), time.Local)
			timeSpan.OffsetStartDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
			resultEndTime = l.tempSubscribeStartTime.Add(time.Minute * time.Duration(l.Duration))
			timeSpan.OffsetEndDate = resultEndTime.Format("2006-01-02 15:04")
			// 判断开始时间是否小于 时段开始时间，则使用开始时间
			if l.tempSubscribeStartTime.Unix() < periodStartTime.Unix() {
				resultEndTime = periodStartTime.Add(time.Minute * time.Duration(l.Duration))
				timeSpan.OffsetStartDate = periodStartTime.Format("2006-01-02 15:04")
				timeSpan.OffsetEndDate = resultEndTime.Format("2006-01-02 15:04")
				//fmt.Println("小于", resultEndTime.Format("2006-01-02 15:04"))
			}

			if resultEndTime.Unix() >= periodEndTime.Unix() {
				exceedSecond := resultEndTime.Unix() - periodEndTime.Unix()
				timeSpan.ExceedStartDate = periodEndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
				timeSpan.ExceedEndDate = resultEndTime.Format("2006-01-02 15:04")
				timeSpan.ExceedDuration = exceedSecond
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.Type = 4
				if exceedSecond > l.Neutron*60 {
					timeSpan.IsExceed = 2
					//fmt.Println(fmt.Sprintf("错误，超出可用范围，已超出%v分钟,最大超出限制:%v分钟,值：%v", exceedSecond/60, Neutron, exceedSecond), resultEndTime.Format("2006-01-02 15:04"), periodEndTime.Format("2006-01-02 15:04"))
					lists = append(lists, timeSpan)
				} else {
					timeSpan.IsExceed = 1
					lists = append(lists, timeSpan)
					//// 正常使用
					//break
				}
				//fmt.Println("大于时间段时间", resultEndTime.Format("2006-01-02 15:04"), periodEndTime.Format("2006-01-02 15:04"), exceedSecond)
			} else {
				// 正常使用时段
				timeSpan.ExceedStartDate = ""
				timeSpan.ExceedEndDate = ""
				timeSpan.ExceedDuration = 0
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.Type = 4
				timeSpan.IsExceed = 0
				lists = append(lists, timeSpan)
				// 正常使用
				break
			}
		}
		dayi++
		l.tempSubscribeStartTime = l.StartTime.Add(time.Hour * time.Duration(24*dayi))
		currentTime = l.tempSubscribeStartTime
		if currentTime.After(lastTime) {
			l.tempSubscribeStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodStartHour), time.Local)
			//fmt.Println("---当", currentTime.Format("2006-01-02 15:04"), lastTime.Format("2006-01-02 15:04"), tempSubscribeStartTime.Format("2006-01-02 15:04"))
		}

		// 结束时间大于 开始时间
		if l.PeriodEndHour > l.PeriodStartHour {
			if l.tempSubscribeStartTime.Format("2006-01-02") == l.EndTime.Format("2006-01-02") && l.tempSubscribeStartTime.After(l.EndTime) {
				if l.EndTime.Format("15:04") < l.PeriodStartHour {
					timeSpan = TimeSpan{}
					timeSpan.Type = 4
					timeSpan.IsExceed = 3
					lists = append(lists, timeSpan)
					//fmt.Println("超出结束时间，且不在使用时段内")
					break
				}
			}
		}
		//time.Sleep(time.Second * 1)
	}
	var timeSpan *TimeSpan
	for _, value := range lists {
		val := value
		if value.IsExceed == 0 || value.IsExceed == 1 {
			return nil, &val
		}
		if timeSpan == nil && (value.IsExceed == 2 || value.IsExceed == 3) {
			timeSpan = &val
		}
	}
	return nil, timeSpan
}

func (l *TrimTime) Period() *TimeSpan {
	err, timeSpan := l.spanTime()
	if err != nil {
		err, timeSpan := l.extractPeriod()
		if err == nil && timeSpan != nil {
			return timeSpan
		}
	}
	return timeSpan
}
