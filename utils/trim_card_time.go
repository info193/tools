package utils

import (
	"fmt"
	"math"
	"time"
)

var (
	TrimCardTimeAvailable                int64 = 0
	TrimCardTimeAvailableExceedPeriod    int64 = 1
	TrimCardTimeNotAvailableExceedPeriod int64 = 2
	TrimCardTimeNotAvailableRange        int64 = 3
)

type TrimCardTime struct {
	Duration               int64  // 卡券时长
	Neutron                int64  // 边界时间
	PeriodStartHour        string // 开始时段 00:00
	PeriodEndHour          string // 结束时段 05:00
	StartTime              time.Time
	EndTime                time.Time
	tempSubscribeStartTime time.Time
}

type TimeCardSpan struct {
	OffsetStartDate          string // 抵消开始时间
	OffsetEndDate            string // 抵消结束时间
	OffsetDuration           int64  // 抵消时长
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

func NewTrimCardTime(duration int64, periodStartHour, periodEndHour, subscribeStartDate, subscribeEndDate string) TrimCardTime {
	subscribeStartTime, _ := time.ParseInLocation("2006-01-02 15:04", subscribeStartDate, time.Local)
	subscribeEndTime, _ := time.ParseInLocation("2006-01-02 15:04", subscribeEndDate, time.Local)
	trimTime := TrimCardTime{
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
func (l *TrimCardTime) setNeutron() {
	l.Neutron = 1
}

func (l *TrimCardTime) minute(cst, cet time.Time) int64 {
	second := cet.Unix() - cst.Unix()
	if second < 60 {
		return 0
	}
	c := math.Floor(float64(second) / 60)
	return int64(c)
}

// 获取可用时段 extractPeriod
func (l *TrimCardTime) extractPeriod() (error, *TimeCardSpan) {
	var periodStartTime time.Time
	var periodEndTime time.Time
	//var resultEndTime time.Time
	var lastTime time.Time
	currentTime := l.tempSubscribeStartTime
	lists := make([]TimeCardSpan, 0)
	var dayi int64
	for {
		timeSpan := TimeCardSpan{}
		timeSpan.AvailableStartPeriodHour = l.PeriodStartHour // 可用开始时段
		timeSpan.AvailableEndPeriodHour = l.PeriodEndHour     // 可用结束时段
		lastTime = l.tempSubscribeStartTime
		if l.tempSubscribeStartTime.Format("2006-01-02") > l.EndTime.Format("2006-01-02") {
			//fmt.Println("跳出", l.tempSubscribeStartTime.Format("2006-01-02 15:04"))
			break
		}
		// 隔日
		if l.PeriodStartHour > l.PeriodEndHour {
			fmt.Println("隔日")
			periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodStartHour), time.Local)
			periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Add(time.Second*time.Duration(86400)).Format("2006-01-02"), l.PeriodEndHour), time.Local)
			timeSpan.Type = 1

			// 判断用户预约开始时间 如果大于等于 实际可抵消结束时间 则不允许使用
			if l.StartTime.Format("2006-01-02") == l.EndTime.Format("2006-01-02") && l.EndTime.Unix() <= periodStartTime.Unix() {
				timeSpan.IsExceed = TrimCardTimeNotAvailableRange
				fmt.Println("不在使用范围内0", timeSpan)
			}
		}
		// 当日
		if l.PeriodEndHour > l.PeriodStartHour {
			fmt.Println("当日")
			periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodStartHour), time.Local)
			periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodEndHour), time.Local)
			timeSpan.Type = 4
		}
		// 判断用户预约开始时间 如果大于等于 实际可抵消结束时间 则不允许使用
		if l.StartTime.Unix() >= periodEndTime.Unix() {
			timeSpan.IsExceed = TrimCardTimeNotAvailableRange
			fmt.Println("不在使用范围内1", timeSpan)
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
				timeSpan.OffsetDuration = l.minute(cst, cet)          // 抵消时长
				timeSpan.IsExceed = TrimCardTimeAvailableExceedPeriod // 1抵扣部分
				if timeSpan.OffsetDuration == l.Duration {
					timeSpan.IsExceed = TrimCardTimeAvailable // 0全部抵扣完成
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
				timeSpan.IsExceed = TrimCardTimeAvailableExceedPeriod // 1抵扣部分
				if timeSpan.OffsetDuration == l.Duration {
					timeSpan.IsExceed = TrimCardTimeAvailable // 0全部抵扣完成
				}
			}

			//fmt.Println(reckonOffsetEndDate.Format("2006-01-02 15:04"), "--------", timeSpan)
		}
		lists = append(lists, timeSpan)

		//// 隔日
		//if l.PeriodStartHour > l.PeriodEndHour {
		//	//fmt.Println("隔日")
		//	periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodStartHour), time.Local)
		//	periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Add(time.Second*time.Duration(86400)).Format("2006-01-02"), l.PeriodEndHour), time.Local)
		//	timeSpan.Type = 1
		//
		//	fmt.Println(periodStartTime.Format("2006-01-02 15:04"), periodEndTime.Format("2006-01-02 15:04"), "-**-*-*-**-")
		//	// 判断用户预约开始时间 如果大于等于 实际可抵消结束时间 则不允许使用
		//	if l.StartTime.Format("2006-01-02") == l.EndTime.Format("2006-01-02") && l.EndTime.Unix() <= periodStartTime.Unix() {
		//		timeSpan.IsExceed = TrimCardTimeNotAvailableRange
		//		fmt.Println("不在使用范围内0", timeSpan)
		//	}
		//	if l.StartTime.Unix() >= periodEndTime.Unix() {
		//		timeSpan.IsExceed = TrimCardTimeNotAvailableRange
		//		fmt.Println("不在使用范围内1", timeSpan)
		//	}
		//	// 判断用户预约开始时间 如果小于等于 实际可抵消结束时间 则允许使用
		//	if l.StartTime.Unix() < periodEndTime.Unix() {
		//		timeSpan.OffsetStartDate = periodStartTime.Format("2006-01-02 15:04")               // 抵消开始时间
		//		reckonOffsetEndDate := periodStartTime.Add(time.Minute * time.Duration(l.Duration)) // 预估抵消结束时间
		//		// 判断用户预约开始时间是否大于 实际可抵消开始时间
		//		if l.StartTime.Unix() > periodStartTime.Unix() {
		//			reckonOffsetEndDate = l.tempSubscribeStartTime.Add(time.Minute * time.Duration(l.Duration)) // 预估抵消结束时间
		//			timeSpan.OffsetStartDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")              // 抵消开始时间
		//		}
		//
		//		// 判断预估抵消结束时间 是否大于用户预约时间结束时间
		//		if reckonOffsetEndDate.Unix() > l.EndTime.Unix() {
		//			reckonOffsetEndDate = l.EndTime
		//		}
		//
		//		// 判断预估抵消结束时间 是否大于 实际可抵消结束时间
		//		if reckonOffsetEndDate.Unix() > periodEndTime.Unix() {
		//			reckonOffsetEndDate = periodEndTime
		//		}
		//		//fmt.Println(reckonOffsetEndDate.Format("2006-01-02 15:04"), "------")
		//		//fmt.Println("====", reckonOffsetEndDate.Format("2006-01-02 15:04"), timeSpan.OffsetStartDate, "===", timeSpan.OffsetEndDate)
		//		// 预估抵消结束时间 大于实际可抵消结束时间 并且 大于 实际可抵消结束时间
		//		if reckonOffsetEndDate.Unix() > periodStartTime.Unix() && reckonOffsetEndDate.Unix() >= periodEndTime.Unix() {
		//			timeSpan.OffsetEndDate = reckonOffsetEndDate.Format("2006-01-02 15:04") // 抵消结束时间
		//			cst, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetStartDate, time.Local)
		//			cet, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetEndDate, time.Local)
		//			timeSpan.OffsetDuration = l.minute(cst, cet)          // 抵消时长
		//			timeSpan.IsExceed = TrimCardTimeAvailableExceedPeriod // 1抵扣部分
		//			if timeSpan.OffsetDuration == l.Duration {
		//				timeSpan.IsExceed = TrimCardTimeAvailable // 0全部抵扣完成
		//			}
		//			//fmt.Println("-------||||", reckonOffsetEndDate.Format("2006-01-02 15:04"))
		//			fmt.Println("验证1", timeSpan)
		//			//break
		//		}
		//
		//		// 预估抵消结束时间 大于等于实际可抵消开始时间  并且 小于 实际可抵消结束时间
		//		if reckonOffsetEndDate.Unix() > periodStartTime.Unix() && reckonOffsetEndDate.Unix() < periodEndTime.Unix() {
		//			timeSpan.OffsetEndDate = reckonOffsetEndDate.Format("2006-01-02 15:04") // 抵消结束时间
		//			cst, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetStartDate, time.Local)
		//			cet, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetEndDate, time.Local)
		//			timeSpan.OffsetDuration = l.minute(cst, cet) // 抵消时长
		//			fmt.Println("预估抵消结束时间 大于等于实际可抵消开始时间  并且 小于 实际可抵消结束时间")
		//			timeSpan.IsExceed = TrimCardTimeAvailableExceedPeriod // 1抵扣部分
		//			if timeSpan.OffsetDuration == l.Duration {
		//				timeSpan.IsExceed = TrimCardTimeAvailable // 0全部抵扣完成
		//			}
		//		}
		//		lists = append(lists, timeSpan)
		//		//fmt.Println(reckonOffsetEndDate.Format("2006-01-02 15:04"), "--------jjjjjjjjjjjj", timeSpan)
		//	}
		//
		//}
		//// 当日
		//if l.PeriodEndHour > l.PeriodStartHour {
		//	fmt.Println("当日")
		//	periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodStartHour), time.Local)
		//	periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", l.tempSubscribeStartTime.Format("2006-01-02"), l.PeriodEndHour), time.Local)
		//	timeSpan.Type = 4
		//	// 判断用户预约开始时间 如果大于等于 实际可抵消结束时间 则不允许使用
		//	if l.StartTime.Unix() >= periodEndTime.Unix() {
		//		timeSpan.IsExceed = TrimCardTimeNotAvailableRange
		//		fmt.Println("不在使用范围内1", timeSpan)
		//	}
		//	// 判断用户预约开始时间 如果小于等于 实际可抵消结束时间 则允许使用
		//	if l.StartTime.Unix() < periodEndTime.Unix() {
		//		timeSpan.OffsetStartDate = periodStartTime.Format("2006-01-02 15:04")               // 抵消开始时间
		//		reckonOffsetEndDate := periodStartTime.Add(time.Minute * time.Duration(l.Duration)) // 预估抵消结束时间
		//		// 判断用户预约开始时间是否大于 实际可抵消开始时间
		//		if l.StartTime.Unix() > periodStartTime.Unix() {
		//			reckonOffsetEndDate = l.tempSubscribeStartTime.Add(time.Minute * time.Duration(l.Duration)) // 预估抵消结束时间
		//			timeSpan.OffsetStartDate = l.tempSubscribeStartTime.Format("2006-01-02 15:04")              // 抵消开始时间
		//		}
		//
		//		// 判断预估抵消结束时间 是否大于用户预约时间结束时间
		//		if reckonOffsetEndDate.Unix() > l.EndTime.Unix() {
		//			reckonOffsetEndDate = l.EndTime
		//		}
		//
		//		// 判断预估抵消结束时间 是否大于 实际可抵消结束时间
		//		if reckonOffsetEndDate.Unix() > periodEndTime.Unix() {
		//			reckonOffsetEndDate = periodEndTime
		//		}
		//		//fmt.Println(reckonOffsetEndDate.Format("2006-01-02 15:04"), "------")
		//		//fmt.Println("====", reckonOffsetEndDate.Format("2006-01-02 15:04"), timeSpan.OffsetStartDate, timeSpan.OffsetEndDate)
		//		// 预估抵消结束时间 大于实际可抵消结束时间 并且 大于 实际可抵消结束时间
		//		if reckonOffsetEndDate.Unix() > periodStartTime.Unix() && reckonOffsetEndDate.Unix() >= periodEndTime.Unix() {
		//			timeSpan.OffsetEndDate = reckonOffsetEndDate.Format("2006-01-02 15:04") // 抵消结束时间
		//			cst, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetStartDate, time.Local)
		//			cet, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetEndDate, time.Local)
		//			timeSpan.OffsetDuration = l.minute(cst, cet)          // 抵消时长
		//			timeSpan.IsExceed = TrimCardTimeAvailableExceedPeriod // 1抵扣部分
		//			if timeSpan.OffsetDuration == l.Duration {
		//				timeSpan.IsExceed = TrimCardTimeAvailable // 0全部抵扣完成
		//			}
		//			//fmt.Println("-------||||", reckonOffsetEndDate.Format("2006-01-02 15:04"))
		//			//fmt.Println("验证1", timeSpan)
		//			//break
		//		}
		//
		//		// 预估抵消结束时间 大于等于实际可抵消开始时间  并且 小于 实际可抵消结束时间
		//		if reckonOffsetEndDate.Unix() > periodStartTime.Unix() && reckonOffsetEndDate.Unix() < periodEndTime.Unix() {
		//			timeSpan.OffsetEndDate = reckonOffsetEndDate.Format("2006-01-02 15:04") // 抵消结束时间
		//			cst, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetStartDate, time.Local)
		//			cet, _ := time.ParseInLocation("2006-01-02 15:04", timeSpan.OffsetEndDate, time.Local)
		//			timeSpan.OffsetDuration = l.minute(cst, cet) // 抵消时长
		//			//fmt.Println("预估抵消结束时间 大于等于实际可抵消开始时间  并且 小于 实际可抵消结束时间")
		//			timeSpan.IsExceed = TrimCardTimeAvailableExceedPeriod // 1抵扣部分
		//			if timeSpan.OffsetDuration == l.Duration {
		//				timeSpan.IsExceed = TrimCardTimeAvailable // 0全部抵扣完成
		//			}
		//		}
		//
		//		//fmt.Println(reckonOffsetEndDate.Format("2006-01-02 15:04"), "--------", timeSpan)
		//	}
		//	lists = append(lists, timeSpan)
		//}

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
					timeSpan = TimeCardSpan{}
					timeSpan.Type = 4
					timeSpan.IsExceed = TrimCardTimeNotAvailableRange
					lists = append(lists, timeSpan)
					//fmt.Println("超出结束时间，且不在使用时段内")
					break
				}
			}
		}
		//time.Sleep(time.Second * 1)
	}

	fmt.Println(lists, "--------")
	var timeSpan *TimeCardSpan
	for _, value := range lists {
		val := value
		if value.IsExceed == TrimCardTimeAvailable || value.IsExceed == TrimCardTimeAvailableExceedPeriod {
			return nil, &val
		}
		if timeSpan == nil && (value.IsExceed == TrimCardTimeNotAvailableExceedPeriod || value.IsExceed == TrimCardTimeNotAvailableRange) {
			timeSpan = &val
		}
	}
	return nil, timeSpan
}

func (l *TrimCardTime) Period() *TimeCardSpan {
	//err, timeSpan := l.spanTimeTwo()
	//if err != nil {
	err, timeSpan := l.extractPeriod()
	if err == nil && timeSpan != nil {
		return timeSpan
	}
	//}
	return timeSpan
}
