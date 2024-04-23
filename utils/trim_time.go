package utils

import (
	"errors"
	"fmt"
	"math"
	"time"
)

var (
	TrimTimeAvailable                int64 = 0
	TrimTimeAvailableExceedPeriod    int64 = 1
	TrimTimeNotAvailableExceedPeriod int64 = 2
	TrimTimeNotAvailableRange        int64 = 3
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
	DeductionDuration        int64  // 抵扣时长 isExceed  =1 的时候才会有值
	ExceedDuration           int64  // 超出时长
	MaxLimitDuration         int64  // 最大限制超出时长
	Type                     int64  // 类型 1隔日（正常可用范围） 2隔日（在预约开始结束时间及时段结束时间范围内） 3隔日（预约开始时间大于时段结束时间范围内）4当日
	IsExceed                 int64  // 是否超出 0未超出（可用） 1超出预约时段范围（可用）  2超出可用范围  3不在使用时段范围内
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
	// 同一天 优惠券不在时段开始及结束使用范围内
	//if l.StartTime.Format("2006-01-02") == l.EndTime.Format("2006-01-02") {
	//	tsPeriodStartTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.StartTime.Format("2006-01-02"), l.PeriodStartHour), time.Local)
	//	tsPeriodEndTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.StartTime.Format("2006-01-02"), l.PeriodEndHour), time.Local)
	//	//fmt.Println(tsPeriodStartTime.Format("2006-01-02 15:04"), tsPeriodEndTime.Format("2006-01-02 15:04"))
	//	if tsPeriodStartTime.Unix() < tsPeriodEndTime.Unix() && l.EndTime.Unix() <= tsPeriodStartTime.Unix() {
	//		timeSpan.ExceedStartDate = ""
	//		timeSpan.ExceedEndDate = ""
	//		timeSpan.ExceedDuration = 0
	//		timeSpan.MaxLimitDuration = l.Neutron
	//		timeSpan.Type = 4
	//		timeSpan.IsExceed = 3
	//		return nil, &timeSpan
	//	}
	//}
	//fmt.Println("spanTime........")
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
				timeSpan.IsExceed = TrimTimeAvailableExceedPeriod
			} else {
				timeSpan.OffsetEndDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.IsExceed = TrimTimeAvailable
			}

			if exceedSecond > l.Neutron*60 {
				timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod
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
				timeSpan.IsExceed = TrimTimeAvailableExceedPeriod
			} else {
				timeSpan.OffsetEndDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.IsExceed = TrimTimeAvailable
			}
			if exceedSecond > l.Neutron*60 {
				timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod
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
				timeSpan.IsExceed = TrimTimeAvailable
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
				timeSpan.IsExceed = TrimTimeAvailable
				if exceedSecond > l.Neutron*60 {
					timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod
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
				timeSpan.IsExceed = TrimTimeAvailable
				if exceedSecond > l.Neutron*60 {
					timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod
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
		timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod
		return nil, &timeSpan
		//fmt.Println("正常时间内", subscribeStartTime.Format("2006-01-02 15:04"), subscribeEndTime.Format("2006-01-02 15:04"))
	}
	return errors.New("错误"), nil
}

func (l *TrimTime) minute(cst, cet time.Time) int64 {
	second := cet.Unix() - cst.Unix()
	if second < 60 {
		return 0
	}
	c := math.Floor(float64(second) / 60)
	return int64(c)
}

// 获取可用时段 extractPeriod
func (l *TrimTime) extractPeriod() (error, *TimeSpan) {
	var periodStartTime time.Time
	var periodEndTime time.Time
	var resultEndTime time.Time
	var lastTime time.Time
	currentTime := l.tempSubscribeStartTime
	lists := make([]TimeSpan, 0)
	// 判断隔日 预约同一天的
	if l.PeriodStartHour > l.PeriodEndHour {
		tperiodEndTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodEndHour), time.Local)
		//fmt.Println(periodStartTime.Format("2006-01-02 15:04"), periodEndTime.Format("2006-01-02 15:04"), "======")
		if l.tempSubscribeStartTime.Unix() < tperiodEndTime.Unix() {

			timeSpan := TimeSpan{}
			timeSpan.AvailableStartPeriodHour = l.PeriodStartHour
			timeSpan.AvailableEndPeriodHour = l.PeriodEndHour
			timeSpan.OffsetStartDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
			tResultEndTime := l.tempSubscribeStartTime.Add(time.Minute * time.Duration(l.Duration))
			// 超出部分
			if tResultEndTime.Unix() > tperiodEndTime.Unix() {
				timeSpan.OffsetEndDate = tperiodEndTime.Format("2006-01-02 15:04")
				timeSpan.ExceedStartDate = tperiodEndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
				timeSpan.ExceedEndDate = tResultEndTime.Format("2006-01-02 15:04")
				exceedSecond := tResultEndTime.Unix() - tperiodEndTime.Unix()
				timeSpan.ExceedDuration = exceedSecond
				timeSpan.MaxLimitDuration = l.Neutron
				timeSpan.Type = 2
				timeSpan.IsExceed = TrimTimeAvailable
				if exceedSecond > l.Neutron*60 {
					timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod
					lists = append(lists, timeSpan)
				} else {
					cst, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetStartDate, time.Local)
					cet, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetEndDate, time.Local)
					timeSpan.DeductionDuration = l.minute(cst, cet)
					timeSpan.IsExceed = TrimTimeAvailableExceedPeriod
					lists = append(lists, timeSpan)
				}
				//fmt.Println("1111111")
				// 未超出
			} else {
				//fmt.Println("未超出")
				if l.tempSubscribeStartTime.Unix() < tperiodEndTime.Unix() && tResultEndTime.Unix() <= tperiodEndTime.Unix() && tResultEndTime.Unix() < l.EndTime.Unix() {
					//fmt.Println(tResultEndTime.Format("2006-01-02 15:04"), "ssssssss11111")
					timeSpan.OffsetEndDate = tResultEndTime.Format("2006-01-02 15:04")
					timeSpan.ExceedStartDate = ""
					timeSpan.ExceedEndDate = ""
					timeSpan.ExceedDuration = 0
					timeSpan.MaxLimitDuration = l.Neutron
					timeSpan.Type = 2
					timeSpan.IsExceed = TrimTimeAvailable
					lists = append(lists, timeSpan)
				}
				if l.tempSubscribeStartTime.Unix() < tperiodEndTime.Unix() && tResultEndTime.Unix() <= tperiodEndTime.Unix() && tResultEndTime.Unix() > l.EndTime.Unix() {
					timeSpan.OffsetEndDate = l.EndTime.Format("2006-01-02 15:04")
					timeSpan.ExceedStartDate = l.EndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
					timeSpan.ExceedEndDate = tResultEndTime.Format("2006-01-02 15:04")
					timeSpan.ExceedDuration = 0
					timeSpan.MaxLimitDuration = l.Neutron
					timeSpan.Type = 2
					timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod
					exceedSecond := tResultEndTime.Unix() - l.EndTime.Unix()
					if exceedSecond > l.Neutron*60 {
						timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod
						lists = append(lists, timeSpan)
					} else {
						cst, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetStartDate, time.Local)
						cet, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetEndDate, time.Local)
						timeSpan.DeductionDuration = l.minute(cst, cet)
						timeSpan.IsExceed = TrimTimeAvailableExceedPeriod
						lists = append(lists, timeSpan)
					}
				}
			}
		}
	}
	if len(lists) == 0 {
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
				//fmt.Println("隔日")
				periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodStartHour), time.Local)
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Add(time.Second*time.Duration(86400)).Format("2006-01-02"), l.PeriodEndHour), time.Local)
				//fmt.Println(periodStartTime.Format("2006-01-02 15:04"), periodEndTime.Format("2006-01-02 15:04"), "======")
				timeSpan.OffsetStartDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")
				if !l.tempSubscribeStartTime.After(periodStartTime) {
					l.tempSubscribeStartTime = periodStartTime
					timeSpan.OffsetStartDate = periodStartTime.Format("2006-01-02 15:04")
				}
				resultEndTime = l.tempSubscribeStartTime.Add(time.Minute * time.Duration(l.Duration))
				timeSpan.OffsetEndDate = resultEndTime.Format("2006-01-02 15:04")

				// 预约开始时间+优惠时段 大于 预约结束时间 开始时段 12:00 - 04:00  预约时间2024-04-06 23:59 - 2024-04-07 02:30  时长180
				if resultEndTime.Unix() > l.EndTime.Unix() && resultEndTime.Unix() < periodEndTime.Unix() {
					timeSpan.OffsetEndDate = l.EndTime.Format("2006-01-02 15:04")
					timeSpan.ExceedStartDate = l.EndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
					timeSpan.ExceedEndDate = resultEndTime.Format("2006-01-02 15:04")
					exceedSecond := resultEndTime.Unix() - l.EndTime.Unix()
					timeSpan.ExceedDuration = exceedSecond
					timeSpan.MaxLimitDuration = l.Neutron
					timeSpan.Type = 2
					timeSpan.IsExceed = TrimTimeAvailable
					//fmt.Println("resultEndTime>EndTime 222", exceedSecond)
					if exceedSecond > l.Neutron*60 {
						timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod
						lists = append(lists, timeSpan)
						//fmt.Println(fmt.Sprintf("错误1，超出可用范围，已超出%v分钟,最大超出限制:%v分钟,值：%v", exceedSecond/60, Neutron, exceedSecond), resultEndTime.Format("2006-01-02 15:04"), periodEndTime.Format("2006-01-02 15:04"))
						//fmt.Println("错误1", timeSpan)
						//return
						//break
					} else {
						cst, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetStartDate, time.Local)
						cet, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetEndDate, time.Local)
						timeSpan.DeductionDuration = l.minute(cst, cet)
						timeSpan.IsExceed = TrimTimeAvailableExceedPeriod
						lists = append(lists, timeSpan)
						//// 正常使用
						//break
					}

					//fmt.Println("超出，但是可以使用，在可用范围内1", timeSpan)
				}

				// 预约开始时间+优惠时段  大于 时段结束时间 开始时段 12:00 - 02:00  预约时间2024-04-06 23:59 - 2024-04-07 01:30  时长180
				if resultEndTime.Unix() >= periodEndTime.Unix() {
					timeSpan.OffsetEndDate = periodEndTime.Format("2006-01-02 15:04")
					timeSpan.ExceedStartDate = periodEndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
					timeSpan.ExceedEndDate = resultEndTime.Format("2006-01-02 15:04")
					exceedSecond := resultEndTime.Unix() - periodEndTime.Unix()
					timeSpan.ExceedDuration = exceedSecond
					timeSpan.MaxLimitDuration = l.Neutron
					timeSpan.Type = 3
					timeSpan.IsExceed = TrimTimeAvailable
					//fmt.Println("resultEndTime>periodEndTime 111", exceedSecond)
					if exceedSecond > l.Neutron*60 {
						timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod
						lists = append(lists, timeSpan)
						//fmt.Println(fmt.Sprintf("错误2，超出可用范围，已超出%v分钟,最大超出限制:%v分钟,值：%v", exceedSecond/60, Neutron, exceedSecond), resultEndTime.Format("2006-01-02 15:04"), periodEndTime.Format("2006-01-02 15:04"))
						//fmt.Println("错误2", timeSpan)
						//return
						//break
					} else {
						cst, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetStartDate, time.Local)
						cet, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetEndDate, time.Local)
						timeSpan.DeductionDuration = l.minute(cst, cet)
						timeSpan.IsExceed = TrimTimeAvailableExceedPeriod
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
					timeSpan.IsExceed = TrimTimeAvailable
					lists = append(lists, timeSpan)
					//fmt.Println("resultEndTime>periodEndTime 33333")
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

				if l.StartTime.Format("2006-01-02") == l.EndTime.Format("2006-01-02") {
					// 时段开始时间 大于等于 预约结束时间
					if l.EndTime.Unix() <= periodStartTime.Unix() {
						timeSpan.ExceedStartDate = ""
						timeSpan.ExceedEndDate = ""
						timeSpan.ExceedDuration = 0
						timeSpan.MaxLimitDuration = l.Neutron
						timeSpan.Type = 4
						timeSpan.IsExceed = TrimTimeNotAvailableRange
						lists = append(lists, timeSpan)
						break
					}
					// 可用时段开始时间 大于等于 预约结束时间
					if resultEndTime.Unix() >= l.EndTime.Unix() {
						exceedSecond := resultEndTime.Unix() - l.EndTime.Unix()
						if l.Duration*60-exceedSecond < l.Neutron*60 {
							timeSpan.OffsetStartDate = periodStartTime.Format("2006-01-02 15:04")
							timeSpan.OffsetEndDate = l.EndTime.Format("2006-01-02 15:04")
							timeSpan.ExceedStartDate = ""
							timeSpan.ExceedEndDate = ""
							timeSpan.ExceedDuration = 0
							timeSpan.MaxLimitDuration = l.Neutron
							timeSpan.Type = 4
							timeSpan.IsExceed = TrimTimeNotAvailableRange
							lists = append(lists, timeSpan)
							break
						}
					}
				}

				if resultEndTime.Unix() >= periodEndTime.Unix() {
					exceedSecond := resultEndTime.Unix() - periodEndTime.Unix()
					timeSpan.OffsetEndDate = periodEndTime.Format("2006-01-02 15:04")
					timeSpan.ExceedStartDate = periodEndTime.Add(time.Minute * 1).Format("2006-01-02 15:04")
					timeSpan.ExceedEndDate = resultEndTime.Format("2006-01-02 15:04")
					timeSpan.ExceedDuration = exceedSecond
					timeSpan.MaxLimitDuration = l.Neutron
					timeSpan.Type = 4
					if exceedSecond > l.Neutron*60 {
						timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod
						//fmt.Println(fmt.Sprintf("错误，超出可用范围，已超出%v分钟,最大超出限制:%v分钟,值：%v", exceedSecond/60, Neutron, exceedSecond), resultEndTime.Format("2006-01-02 15:04"), periodEndTime.Format("2006-01-02 15:04"))
						lists = append(lists, timeSpan)
					} else {
						cst, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetStartDate, time.Local)
						cet, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetEndDate, time.Local)
						timeSpan.DeductionDuration = l.minute(cst, cet)
						timeSpan.IsExceed = TrimTimeAvailableExceedPeriod
						lists = append(lists, timeSpan)
						//// 正常使用
						//break
					}
					//fmt.Println("大于时间段时间", resultEndTime.Format("2006-01-02 15:04"), periodEndTime.Format("2006-01-02 15:04"), exceedSecond)
				} else {
					timeSpan.OffsetStartDate = periodStartTime.Format("2006-01-02 15:04")
					if l.StartTime.Unix() > periodStartTime.Unix() {
						timeSpan.OffsetStartDate = l.StartTime.Format("2006-01-02 15:04")
					}
					if l.EndTime.Unix() < periodEndTime.Unix() && resultEndTime.Unix() < periodEndTime.Unix() && l.EndTime.Unix() < resultEndTime.Unix() {
						timeSpan.OffsetEndDate = l.EndTime.Format("2006-01-02 15:04")
					}
					if l.EndTime.Unix() < periodEndTime.Unix() && resultEndTime.Unix() < periodEndTime.Unix() && l.EndTime.Unix() > resultEndTime.Unix() {
						timeSpan.OffsetEndDate = resultEndTime.Format("2006-01-02 15:04")
					}
					// 正常使用时段
					timeSpan.ExceedStartDate = ""
					timeSpan.ExceedEndDate = ""
					timeSpan.ExceedDuration = 0
					timeSpan.MaxLimitDuration = l.Neutron
					timeSpan.Type = 4
					timeSpan.IsExceed = TrimTimeAvailable
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
						timeSpan.IsExceed = TrimTimeNotAvailableRange
						lists = append(lists, timeSpan)
						//fmt.Println("超出结束时间，且不在使用时段内")
						break
					}
				}
			}
			//time.Sleep(time.Second * 1)
		}
	}
	var timeSpan *TimeSpan
	for _, value := range lists {
		val := value
		if value.IsExceed == TrimTimeAvailable || value.IsExceed == TrimTimeAvailableExceedPeriod {
			return nil, &val
		}
		if timeSpan == nil && (value.IsExceed == TrimTimeNotAvailableExceedPeriod || value.IsExceed == TrimTimeNotAvailableRange) {
			timeSpan = &val
		}
	}
	return nil, timeSpan
}

func (l *TrimTime) Period() *TimeSpan {
	//err, timeSpan := l.spanTimeTwo()
	//if err != nil {
	err, timeSpan := l.extractPeriod()
	if err == nil && timeSpan != nil {
		return timeSpan
	}
	//}
	return timeSpan
}
