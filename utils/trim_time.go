package utils

import (
	"fmt"
	"math"
	"time"
)

var (
	TrimTimeAvailable                int64 = 0
	TrimTimeAvailableExceedPeriod    int64 = 1
	TrimTimeNotAvailableExceedPeriod int64 = 2
	TrimTimeNotAvailableRange        int64 = 3
	OptimalYes                       int64 = 1
	OptimalNo                        int64 = 2
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
	OffsetDuration           int64  // 抵消时长
	AvailableStartPeriodHour string // 可用开始时段
	AvailableEndPeriodHour   string // 可用结束时段
	Duration                 int64  // 时长
	MinOffsetDuration        int64  // 最小抵消时长
	Type                     int64  // 类型 1隔日 2当日
	IsExceed                 int64  // 是否超出 0完全抵扣（可用） 1抵扣部分（可用）  2超出可用范围  3不在使用时段范围内
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

func (l *TrimTime) minute(cst, cet time.Time) int64 {
	second := cet.Unix() - cst.Unix()
	if second < 60 {
		return 0
	}
	c := math.Floor(float64(second) / 60)
	return int64(c)
}

// 获取可用时段 extractPeriod
func (l *TrimTime) extractPeriod(optimal, isType int64) (error, *TimeSpan) {
	var periodStartTime time.Time
	var periodEndTime time.Time
	var lastTime time.Time
	currentTime := l.tempSubscribeStartTime
	lists := make([]TimeSpan, 0)
	var dayi int64
	for {
		timeSpan := TimeSpan{}
		timeSpan.AvailableStartPeriodHour = l.PeriodStartHour // 可用开始时段
		timeSpan.AvailableEndPeriodHour = l.PeriodEndHour     // 可用结束时段
		timeSpan.Duration = l.Duration                        // 时长
		// 优惠券
		if isType == 2 {
			timeSpan.MinOffsetDuration = l.Neutron
		}
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
			timeSpan.Type = 1

			// 判断用户预约开始时间 如果大于等于 实际可抵消结束时间 则不允许使用
			if l.StartTime.Format("2006-01-02") == l.EndTime.Format("2006-01-02") && l.EndTime.Unix() <= periodStartTime.Unix() {
				timeSpan.IsExceed = TrimTimeNotAvailableRange
				//fmt.Println("不在使用范围内0", timeSpan)
			}
			// 判断用户预约结束时间 如果小于 实际可抵消开始时间 则不允许使用
			if l.EndTime.Unix() < periodStartTime.Unix() {
				timeSpan.IsExceed = TrimTimeNotAvailableRange
				//break 如果不需要响应则可用去除注释
			}
		}
		// 当日
		if l.PeriodEndHour > l.PeriodStartHour {
			//fmt.Println("当日")
			periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodStartHour), time.Local)
			periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodEndHour), time.Local)
			timeSpan.Type = 2
		}
		// 判断用户预约开始时间 如果大于等于 实际可抵消结束时间 则不允许使用
		if l.StartTime.Unix() >= periodEndTime.Unix() {
			timeSpan.IsExceed = TrimTimeNotAvailableRange
			//fmt.Println("不在使用范围内1", timeSpan)
		}
		// 判断用户预约开始时间 如果小于等于 实际可抵消结束时间 则允许使用
		if l.StartTime.Unix() < periodEndTime.Unix() {
			timeSpan.OffsetStartDate = periodStartTime.Format("2006-01-02 15:04")               // 抵消开始时间
			reckonOffsetEndDate := periodStartTime.Add(time.Minute * time.Duration(l.Duration)) // 预估抵消结束时间
			// 判断用户预约开始时间是否大于 实际可抵消开始时间
			if l.StartTime.Unix() > periodStartTime.Unix() {
				reckonOffsetEndDate = l.tempSubscribeStartTime.Add(time.Minute * time.Duration(l.Duration)) // 预估抵消结束时间
				timeSpan.OffsetStartDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")              // 抵消开始时间
			}

			// 判断预估抵消结束时间 是否大于用户预约时间结束时间
			if reckonOffsetEndDate.Unix() > l.EndTime.Unix() {
				reckonOffsetEndDate = l.EndTime
			}

			// 判断预估抵消结束时间 是否大于 实际可抵消结束时间
			if reckonOffsetEndDate.Unix() > periodEndTime.Unix() {
				reckonOffsetEndDate = periodEndTime
			}
			//fmt.Println(reckonOffsetEndDate.Format("2006-01-02 15:04"), "------")
			//fmt.Println("====", reckonOffsetEndDate.Format("2006-01-02 15:04"), timeSpan.OffsetStartDate, timeSpan.OffsetEndDate)
			// 预估抵消结束时间 大于实际可抵消结束时间 并且 大于 实际可抵消结束时间
			if reckonOffsetEndDate.Unix() > periodStartTime.Unix() && reckonOffsetEndDate.Unix() >= periodEndTime.Unix() {
				timeSpan.OffsetEndDate = reckonOffsetEndDate.Format("2006-01-02 15:04") // 抵消结束时间
				cst, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetStartDate, time.Local)
				cet, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetEndDate, time.Local)
				timeSpan.OffsetDuration = l.minute(cst, cet)      // 抵消时长
				timeSpan.IsExceed = TrimTimeAvailableExceedPeriod // 1抵扣部分
				if timeSpan.OffsetDuration == l.Duration {
					timeSpan.IsExceed = TrimTimeAvailable // 0全部抵扣完成
				}

				// 优惠券
				if timeSpan.OffsetDuration < l.Neutron && isType == 2 {
					timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod // 2 不在使用范围 小于最小使用时长
				}

				//fmt.Println("-------||||", reckonOffsetEndDate.Format("2006-01-02 15:04"))
				//fmt.Println("验证1", timeSpan)
				//break
			}

			// 预估抵消结束时间 大于等于实际可抵消开始时间  并且 小于 实际可抵消结束时间
			if reckonOffsetEndDate.Unix() > periodStartTime.Unix() && reckonOffsetEndDate.Unix() < periodEndTime.Unix() {
				timeSpan.OffsetEndDate = reckonOffsetEndDate.Format("2006-01-02 15:04") // 抵消结束时间
				cst, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetStartDate, time.Local)
				cet, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetEndDate, time.Local)
				timeSpan.OffsetDuration = l.minute(cst, cet) // 抵消时长
				//fmt.Println("预估抵消结束时间 大于等于实际可抵消开始时间  并且 小于 实际可抵消结束时间")
				timeSpan.IsExceed = TrimTimeAvailableExceedPeriod // 1抵扣部分
				if timeSpan.OffsetDuration == l.Duration {
					timeSpan.IsExceed = TrimTimeAvailable // 0全部抵扣完成
				}
				// 优惠券
				if timeSpan.OffsetDuration < l.Neutron && isType == 2 {
					timeSpan.IsExceed = TrimTimeNotAvailableExceedPeriod // 2 不在使用范围 小于最小使用时长
				}
				//fmt.Println("验证2", timeSpan)
			}
			//fmt.Println(reckonOffsetEndDate.Format("2006-01-02 15:04"), "--------", timeSpan)
		}
		lists = append(lists, timeSpan)
		if timeSpan.IsExceed == TrimTimeAvailable {
			break
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
					timeSpan.Type = 2
					timeSpan.IsExceed = TrimTimeNotAvailableRange
					lists = append(lists, timeSpan)
					//fmt.Println("超出结束时间，且不在使用时段内")
					break
				}
			}
		}
		//time.Sleep(time.Second * 1)
	}

	//fmt.Println(lists, "--------", optimal)
	var timeSpan *TimeSpan
	// 获取最优选项
	if optimal == OptimalYes {
		for _, value := range lists {
			val := value
			if value.IsExceed == TrimTimeAvailable {
				return nil, &val
			}
		}
	}
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

// 单店卡
func (l *TrimTime) CardPeriod(optimal int64) *TimeSpan {
	err, timeSpan := l.extractPeriod(optimal, 1)
	if err == nil && timeSpan != nil {
		return timeSpan
	}
	return timeSpan
}

// 优惠券
func (l *TrimTime) CouponPeriod(optimal int64) *TimeSpan {
	err, timeSpan := l.extractPeriod(optimal, 2)
	if err == nil && timeSpan != nil {
		return timeSpan
	}
	return timeSpan
}
