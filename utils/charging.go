package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"time"
)

type ChargePeriod struct {
	StartPeriod int64   `json:"start_period"`
	EndPeriod   int64   `json:"end_period"`
	Start       string  `json:"start"`
	End         string  `json:"end"`
	Price       float64 `json:"price"`
}

type Charging struct {
	periods               map[int64]ChargePeriod
	chargingMode          int64
	cycleMinute           int64
	startDiffDuration     int64     // 开始时段相差 分钟
	endDiffDuration       int64     // 结束时段相差 分钟
	startDiffTime         string    // 开始差集时间
	endDiffTime           string    // 结束差集时间
	startDiffTimeDivision string    // 开始差集时段
	endDiffTimeDivision   string    // 结束差集时段
	markStartTime         time.Time // 标记计费开始时间
	markEndTime           time.Time // 标记计费结束时间
}

type periodTimes struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	Index     int       `json:"index"`
	Start     string    `json:"start"`
	End       string    `json:"end"`
	Price     float64   `json:"price"`
}

type PeriodList struct {
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`
	Index      int     `json:"index"`
	Duration   int64   `json:"duration"`
	Start      string  `json:"start"`
	End        string  `json:"end"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
}

func NewCharge(periods map[int64]ChargePeriod, chargingMode int64, cycleMinute int64) *Charging {
	return &Charging{periods: periods, chargingMode: chargingMode, cycleMinute: cycleMinute}
}

// 计费
func (l *Charging) Outlay(startDate, endDate string) (float64, map[int64]float64, []*PeriodList) {
	var price float64
	periodList := make([]*PeriodList, 0)
	periods := make(map[int64]float64) // 时段费用详情
	startDateTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	endDateTime, _ := time.ParseInLocation("2006-01-02 15:04", endDate, time.Local)
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)

	startDateYMD := startDateTime.Format("2006-01-02")
	endDateYMD := endDateTime.Format("2006-01-02")
	k := 0
	times := make([]periodTimes, 0)
	count := len(l.periods)
	for {
		if currentTime.Unix() >= endDateTime.Unix() {
			break
		}
		var periodStartTime time.Time
		var periodEndTime time.Time
		for i := 0; i < count; i++ {
			if period, ok := l.periods[int64(i)]; ok {
				// 隔日跨日
				if period.Start > period.End && startDateYMD != endDateYMD {
					// 当日跨日
					periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.Start), time.Local)
					if currentTime.Unix() < periodStartTime.Unix() {
						periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", time.Unix(currentTime.Unix()-86400, 0).Format("2006-01-02"), period.Start), time.Local)
						periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
					} else {
						periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.Start), time.Local)
						periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), period.End), time.Local)
					}

				}
				// 当日
				if period.Start > period.End && startDateYMD == endDateYMD && k == 0 {
					if startDateTime.Format("15:04") > period.Start {
						periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.Start), time.Local)
						periodEndTime = endDateTime
					} else {
						periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", time.Unix(currentTime.Unix()-86400, 0).Format("2006-01-02"), period.Start), time.Local)
						periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
					}
					k++
				}

				// 判断 时段结束时间大于时段开始时间
				if period.Start < period.End {
					periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.Start), time.Local)
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
					//fmt.Println("....", periodStartTime.Format(time.DateTime), periodEndTime.Format(time.DateTime))
				}

				if currentTime.Unix() >= periodStartTime.Unix() && currentTime.Unix() <= periodEndTime.Unix() {
					if currentTime.Unix() < periodEndTime.Unix() && periodEndTime.Unix() < endDateTime.Unix() {
						times = append(times, periodTimes{StartTime: currentTime, EndTime: periodEndTime, Index: i, Start: period.Start, End: period.End, Price: period.Price})
					}
					// 判断时段结束时间 大于 传入的结束时间 则赋值且跳出
					if endDateTime.Unix() <= periodEndTime.Unix() {
						times = append(times, periodTimes{StartTime: currentTime, EndTime: endDateTime, Index: i, Start: period.Start, End: period.End, Price: period.Price})
						currentTime = endDateTime
						break
					}
					end := 1
					if currentTime.Unix() == periodEndTime.Unix() && periodEndTime.Unix() < endDateTime.Unix() {
						if periodV, ok := l.periods[int64(i+1)]; ok {
							nextPeriodStartTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), periodV.Start), time.Local)
							if nextPeriodStartTime.Unix() >= periodEndTime.Unix() && nextPeriodStartTime.Unix() < endDateTime.Unix() {
								end = 0
							}
						} else if periodV, ok := l.periods[0]; ok {
							nextPeriodStartTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), periodV.Start), time.Local)
							if nextPeriodStartTime.Unix() >= periodEndTime.Unix() && nextPeriodStartTime.Unix() < endDateTime.Unix() {
								end = 0
							}
						}
					}

					if currentTime.Unix() == periodEndTime.Unix() && periodEndTime.Unix() < endDateTime.Unix() {
						if end == 1 {
							times = append(times, periodTimes{StartTime: currentTime, EndTime: endDateTime, Index: i, Start: period.Start, End: period.End, Price: period.Price})
							currentTime = endDateTime
						} else if currentTime.Unix() != periodEndTime.Unix() {
							times = append(times, periodTimes{StartTime: currentTime, EndTime: periodEndTime, Index: i, Start: period.Start, End: period.End, Price: period.Price})
							currentTime = periodEndTime
						}
						continue
					}
					currentTime = periodEndTime
				}
			}
		}
	}
	// 1按分钟计费
	if l.chargingMode == 1 {
		price, periods, periodList = l.computeMinute(times)
	}
	// 2按半小时计费
	if l.chargingMode == 2 {
		price, periods, periodList = l.computeHalfHour(times)
	}
	// 3按小时计费
	if l.chargingMode == 3 {
		price, periods, periodList = l.computeHour(times)
	}
	// 4按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
	if l.chargingMode == 4 {
		price, periods, periodList = l.computeHalfHourOrMinute(times)
	}
	// 5按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
	if l.chargingMode == 5 {
		price, periods, periodList = l.computeHourOrMinute(times)
	}
	// 6以10分钟作为一个收费周期，第1分钟后开始计费
	if l.chargingMode == 6 {
		price, periods, periodList = l.computeCycleAndMinute(times)
	}
	return price, periods, periodList
}

// 分钟计费
func (l *Charging) computeMinute(times []periodTimes) (float64, map[int64]float64, []*PeriodList) {
	periods := make(map[int64]float64) // 时段费用详情
	var price float64
	periodList := make([]*PeriodList, 0)
	for _, timeV := range times {
		duration := timeV.EndTime.Sub(timeV.StartTime)
		minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		totalPrice, _ := decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		price += totalPrice
		periods[int64(timeV.Index)] += totalPrice
		periodList = append(periodList, &PeriodList{
			StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
			EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
			Index:      timeV.Index,
			Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
			Start:      timeV.Start,
			End:        timeV.End,
			Price:      timeV.Price,
			TotalPrice: totalPrice,
		})
	}

	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, periods, periodList
}

// 半小时计费 (半小时为节点，以半小时起始点计算价格，如果时间段内计算时长小于半小时，则累计至下个时段，
// 并在下个时段内扣除上个时段少于的时间[因为以收半个时段价格]，然后在计算下周期时段价格，依次类推)
func (l *Charging) computeHalfHour(times []periodTimes) (float64, map[int64]float64, []*PeriodList) {
	periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	periodList := make([]*PeriodList, 0)
	for _, timeV := range times {
		if lastDiffMinute > 0 {
			timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		}
		var totalPrice float64
		duration := timeV.EndTime.Sub(timeV.StartTime)
		if duration.Minutes() > 0 {
			haltHour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(30)).Ceil().IntPart()
			maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
			lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()

			halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
			totalPrice, _ = decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour)).Float64()

			price += totalPrice
			periods[int64(timeV.Index)] += totalPrice
		} else {
			periods[int64(timeV.Index)] += 0
		}
		periodList = append(periodList, &PeriodList{
			StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
			EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
			Index:      timeV.Index,
			Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
			Start:      timeV.Start,
			End:        timeV.End,
			Price:      timeV.Price,
			TotalPrice: totalPrice,
		})

	}
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, periods, periodList
}

// 按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
func (l *Charging) computeHalfHourOrMinute(times []periodTimes) (float64, map[int64]float64, []*PeriodList) {
	periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	var firstHalfHour int64
	firstHalfHour = 30
	periodList := make([]*PeriodList, 0)
	for _, timeV := range times {
		if lastDiffMinute != 0 {
			timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		}

		var totalPrice float64
		duration := timeV.EndTime.Sub(timeV.StartTime)
		if duration.Minutes() > 0 {
			if firstHalfHour != 0 {
				haltHour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(30)).Ceil().IntPart()
				maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
				lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
				halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
				totalPrice, _ = decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour)).Float64()
				price += totalPrice
				periods[int64(timeV.Index)] += totalPrice
				firstHalfHour = 0 // 首个半小时计费，计费完后置为0
			} else {
				minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
				totalPrice, _ = decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				price += totalPrice
				periods[int64(timeV.Index)] += totalPrice
				lastDiffMinute = 0 // 分钟计费，差值分钟置为0
			}
		} else {
			periods[int64(timeV.Index)] += 0
		}
		periodList = append(periodList, &PeriodList{
			StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
			EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
			Index:      timeV.Index,
			Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
			Start:      timeV.Start,
			End:        timeV.End,
			Price:      timeV.Price,
			TotalPrice: totalPrice,
		})

	}
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, periods, periodList
}

// 一小时计费 (一小时为节点，以一小时起始点计算价格，如果时间段内计算时长小于一小时，则累计至下个时段，
// 并在下个时段内扣除上个时段少于的时间[因为以收一个时段价格]，然后在计算下周期时段价格，依次类推)
func (l *Charging) computeHour(times []periodTimes) (float64, map[int64]float64, []*PeriodList) {
	periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	periodList := make([]*PeriodList, 0)
	for _, timeV := range times {
		if lastDiffMinute > 0 {
			timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		}
		var totalPrice float64
		duration := timeV.EndTime.Sub(timeV.StartTime)
		if duration.Minutes() > 0 {
			hour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(60)).Ceil().IntPart()
			maxMinute := decimal.NewFromInt(hour).Mul(decimal.NewFromInt(60)).IntPart()
			lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
			totalPrice, _ = decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour)).Float64()
			price += totalPrice
			periods[int64(timeV.Index)] += totalPrice
		} else {
			periods[int64(timeV.Index)] += 0
		}
		periodList = append(periodList, &PeriodList{
			StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
			EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
			Index:      timeV.Index,
			Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
			Start:      timeV.Start,
			End:        timeV.End,
			Price:      timeV.Price,
			TotalPrice: totalPrice,
		})
	}
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, periods, periodList
}

// 按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
func (l *Charging) computeHourOrMinute(times []periodTimes) (float64, map[int64]float64, []*PeriodList) {
	periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	var firstHour int64
	firstHour = 60
	periodList := make([]*PeriodList, 0)
	for _, timeV := range times {
		if lastDiffMinute != 0 {
			timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		}

		var totalPrice float64
		duration := timeV.EndTime.Sub(timeV.StartTime)
		if duration.Minutes() > 0 {
			if firstHour != 0 {
				hour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(60)).Ceil().IntPart()
				maxMinute := decimal.NewFromInt(hour).Mul(decimal.NewFromInt(60)).IntPart()
				lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
				totalPrice, _ = decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour)).Float64()
				price += totalPrice
				periods[int64(timeV.Index)] += totalPrice
				firstHour = 0 // 首个一小时计费，计费完后置为0
			} else {
				minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
				totalPrice, _ = decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				price += totalPrice
				periods[int64(timeV.Index)] += totalPrice
				lastDiffMinute = 0 // 分钟计费，差值分钟置为0
			}
		} else {
			periods[int64(timeV.Index)] += 0
		}
		periodList = append(periodList, &PeriodList{
			StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
			EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
			Index:      timeV.Index,
			Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
			Start:      timeV.Start,
			End:        timeV.End,
			Price:      timeV.Price,
			TotalPrice: totalPrice,
		})

	}
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, periods, periodList
}

// 以10分钟作为一个收费周期，第1分钟后开始计费
func (l *Charging) computeCycleAndMinute(times []periodTimes) (float64, map[int64]float64, []*PeriodList) {
	periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	periodList := make([]*PeriodList, 0)
	for _, timeV := range times {
		if lastDiffMinute != 0 {
			timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		}

		var totalPrice float64
		duration := timeV.EndTime.Sub(timeV.StartTime)
		if duration.Minutes() > 0 {
			// 半小时
			if l.cycleMinute == 30 {
				haltHour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(30)).Ceil().IntPart()
				maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
				lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
				halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
				totalPrice, _ = decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour)).Float64()
			}

			// 一小时
			if l.cycleMinute == 60 {
				hour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(60)).Ceil().IntPart()
				maxMinute := decimal.NewFromInt(hour).Mul(decimal.NewFromInt(60)).IntPart()
				lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
				totalPrice, _ = decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour)).Float64()
			}
			// 周期计费
			if l.cycleMinute != 30 && l.cycleMinute != 60 {
				minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
				cyclePrice, _ := decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
				cycle := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(l.cycleMinute)).Ceil().IntPart()
				maxMinute := decimal.NewFromInt(cycle).Mul(decimal.NewFromInt(l.cycleMinute)).IntPart()
				lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
				totalPrice, _ = decimal.NewFromFloat(cyclePrice).Mul(decimal.NewFromInt(cycle)).Float64()
			}
			price += totalPrice
			periods[int64(timeV.Index)] += totalPrice
		} else {
			periods[int64(timeV.Index)] += 0
		}
		periodList = append(periodList, &PeriodList{
			StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
			EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
			Index:      timeV.Index,
			Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
			Start:      timeV.Start,
			End:        timeV.End,
			Price:      timeV.Price,
			TotalPrice: totalPrice,
		})

	}
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, periods, periodList
}

// 金额算出会有误差，具体以结束计算为准(分钟)
func (l *Charging) moneyTransferMinuteTime(money float64, startDate string) (string, string) {

	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	totalMoney := money
	for {
		if totalMoney <= 0 {
			//fmt.Println("跳出循环")
			break
		}
		var periodEndTime time.Time
		var duration time.Duration

		for _, period := range l.periods {
			var minutePrice float64
			var totalMinute int64
			minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()

			// 跨日
			if period.Start > period.End {
				if currentTime.Format("15:04") >= period.Start {
					//fmt.Println("大于", period.Start, period.Price)
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < period.End {
					//fmt.Println("小于", period.End, period.Price)
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}

				if totalMinute < int64(duration.Minutes()) {
					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
					lastEndTime = currentTime
					totalMoney = 0
				} else {
					currentTime = periodEndTime
					lastEndTime = currentTime
					surplusMoney, _ := decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					totalMoney -= surplusMoney
				}
				continue
			}
			// 当日
			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End {
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
				duration = periodEndTime.Sub(currentTime)

				if totalMinute < int64(duration.Minutes()) {
					//fmt.Println("===|||今日", totalMinute, "...", int64(duration.Minutes()))
					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
					lastEndTime = currentTime
					totalMoney = 0
				} else {
					currentTime = periodEndTime
					lastEndTime = currentTime
					surplusMoney, _ := decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//fmt.Println("-----kkkk今日----", surplusMoney, totalMinute, minutePrice)
					totalMoney -= surplusMoney
				}
			}
		}
	}
	//fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime))
	return startDate, lastEndTime.Format("2006-01-02 15:04")
}

// 金额算出会有误差，具体以结束计算为准(半小时)
func (l *Charging) moneyTransferHalfHourTime(money float64, startDate string) (string, string) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	totalMoney := money
	for {
		if totalMoney <= 0 {
			//fmt.Println("跳出循环")
			break
		}
		for _, period := range l.periods {
			var halfHourPrice float64
			halfHourPrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()

			// 跨日
			if period.Start > period.End && totalMoney > 0 {
				if currentTime.Format("15:04") >= period.Start || currentTime.Format("15:04") < period.End {
					if totalMoney < halfHourPrice {
						minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
						totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
						currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
						lastEndTime = currentTime
						totalMoney = 0
						//fmt.Println(totalMoney, halfHourPrice, "隔日====下与半小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
					} else {
						currentTime = currentTime.Add(30 * time.Minute)
						lastEndTime = currentTime
						totalMoney -= halfHourPrice
						//fmt.Println(totalMoney, halfHourPrice, "隔日====满半小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
					}
				}
				continue
			}
			// 当日 6
			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End && totalMoney > 0 {
				if totalMoney < halfHourPrice {
					//fmt.Println(totalMoney, halfHourPrice, ".......")
					minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
					totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
					lastEndTime = currentTime
					totalMoney = 0
					//fmt.Println(totalMinute, halfHourPrice, totalMinute, "当日====小于半小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
				} else {
					currentTime = currentTime.Add(time.Duration(30) * time.Minute)
					lastEndTime = currentTime
					totalMoney -= halfHourPrice
					//fmt.Println(totalMoney, halfHourPrice, "当日====满足半小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
				}
			}
		}
	}
	//fmt.Println("结果：", startDate, lastEndTime.Format("2006-01-02 15:04"))
	return startDate, lastEndTime.Format("2006-01-02 15:04")
}

// 金额算出会有误差，具体以结束计算为准(一小时)
func (l *Charging) moneyTransferHourTime(money float64, startDate string) (string, string) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	totalMoney := money
	for {
		if totalMoney <= 0 {
			//fmt.Println("跳出循环")
			break
		}
		for _, period := range l.periods {
			//var hourPrice float64
			hourPrice := period.Price

			// 跨日
			if period.Start > period.End && totalMoney > 0 {
				if currentTime.Format("15:04") >= period.Start || currentTime.Format("15:04") < period.End {
					if totalMoney < hourPrice {
						minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
						totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
						currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
						lastEndTime = currentTime
						totalMoney = 0
						//fmt.Println(totalMoney, hourPrice, "隔日====下与一小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
					} else {
						currentTime = currentTime.Add(60 * time.Minute)
						lastEndTime = currentTime
						totalMoney -= hourPrice
						//fmt.Println(totalMoney, hourPrice, "隔日====满一小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
					}
				}
				continue
			}
			// 当日 6
			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End && totalMoney > 0 {
				if totalMoney < hourPrice {
					minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
					totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
					lastEndTime = currentTime
					totalMoney = 0
					//fmt.Println(totalMinute, hourPrice, totalMinute, "当日====小于一小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
				} else {
					currentTime = currentTime.Add(time.Duration(30) * time.Minute)
					lastEndTime = currentTime
					totalMoney -= hourPrice
					//fmt.Println(totalMoney, hourPrice, "当日====满足一小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
				}
			}
		}
	}
	//fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime))
	return startDate, lastEndTime.Format("2006-01-02 15:04")
}

//// 金额算出会有误差，具体以结束计算为准(按小时计费开台不足1小时按小时计费，超过1小时按分钟计费)
//func (l *Charging) MoneyTransferHourOrMinuteTime(money float64, startDate string) (string, string) {
//	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
//	var lastEndTime time.Time
//	totalMoney := money
//	isExceed := 0
//	for {
//		if totalMoney <= 0 {
//			fmt.Println("跳出循环")
//			break
//		}
//		var periodEndTime time.Time
//		var duration time.Duration
//		for _, period := range l.periods {
//			hourPrice := period.Price
//			// 跨日
//			if period.Start > period.End && totalMoney > 0 {
//				if currentTime.Format("15:04") >= period.Start || currentTime.Format("15:04") < period.End {
//					if totalMoney >= hourPrice && isExceed == 0 {
//						currentTime = currentTime.Add(60 * time.Minute)
//						lastEndTime = currentTime
//						totalMoney -= hourPrice
//						isExceed = 1
//						fmt.Println(totalMoney, hourPrice, "隔日====不足一小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//
//						continue
//					}
//
//					if isExceed == 1 {
//						if currentTime.Format("15:04") >= period.Start {
//							fmt.Println("大于", period.Start, period.Price)
//							periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), period.End), time.Local)
//							duration = periodEndTime.Sub(currentTime)
//						}
//						if currentTime.Format("15:04") < period.End {
//							fmt.Println("小于", period.End, period.Price)
//							periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
//							duration = periodEndTime.Sub(currentTime)
//						}
//						var minutePrice float64
//						var totalMinute int64
//						minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//						totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//
//						if totalMinute < int64(duration.Minutes()) {
//							currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//							lastEndTime = currentTime
//							totalMoney = 0
//						} else {
//							currentTime = periodEndTime
//							lastEndTime = currentTime
//							surplusMoney, _ := decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
//							totalMoney -= surplusMoney
//						}
//						fmt.Println(totalMoney, hourPrice, "隔日====不足一小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					}
//					continue
//				}
//			}
//			// 当日 6
//			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End && totalMoney > 0 {
//				if totalMoney >= hourPrice && isExceed == 0 {
//					currentTime = currentTime.Add(time.Duration(60) * time.Minute)
//					lastEndTime = currentTime
//					totalMoney -= hourPrice
//					isExceed = 1
//					fmt.Println(totalMoney, hourPrice, "当日====满足一小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//				} else {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//
//					var minutePrice float64
//					var totalMinute int64
//					minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//					totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//					if totalMinute < int64(duration.Minutes()) {
//						fmt.Println("===|||今日", totalMinute, "...", int64(duration.Minutes()))
//						currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//						lastEndTime = currentTime
//						totalMoney = 0
//						fmt.Println(totalMoney, hourPrice, "今日====不足一小时按分钟计费=====1", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					} else {
//						currentTime = periodEndTime
//						lastEndTime = currentTime
//						surplusMoney, _ := decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
//						totalMoney -= surplusMoney
//						fmt.Println(totalMoney, hourPrice, "今日====不足一小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					}
//
//				}
//			}
//		}
//	}
//	fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime))
//	return "", ""
//}
//
//// 金额算出会有误差，具体以结束计算为准(按半小时计费开台不足半小时按半时计费，超过半小时按分钟计费)
//func (l *Charging) MoneyTransferHalfHourOrMinuteTime(money float64, startDate string) (string, string) {
//	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
//	var lastEndTime time.Time
//	totalMoney := money
//	isExceed := 0
//	for {
//		if totalMoney <= 0 {
//			fmt.Println("跳出循环")
//			break
//		}
//		var periodEndTime time.Time
//		var duration time.Duration
//		for _, period := range l.periods {
//			//hourPrice := period.Price
//			halfHourPrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
//			// 跨日
//			if period.Start > period.End && totalMoney > 0 {
//				if currentTime.Format("15:04") >= period.Start || currentTime.Format("15:04") < period.End {
//					if totalMoney >= halfHourPrice && isExceed == 0 {
//						currentTime = currentTime.Add(30 * time.Minute)
//						lastEndTime = currentTime
//						totalMoney -= halfHourPrice
//						isExceed = 1
//						fmt.Println(totalMoney, halfHourPrice, "隔日====不足半小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//						continue
//					}
//
//					if isExceed == 1 {
//						if currentTime.Format("15:04") >= period.Start {
//							fmt.Println("大于", period.Start, period.Price)
//							periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), period.End), time.Local)
//							duration = periodEndTime.Sub(currentTime)
//						}
//						if currentTime.Format("15:04") < period.End {
//							fmt.Println("小于", period.End, period.Price)
//							periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
//							duration = periodEndTime.Sub(currentTime)
//						}
//						var minutePrice float64
//						var totalMinute int64
//						minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//						totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//
//						if totalMinute < int64(duration.Minutes()) {
//							currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//							lastEndTime = currentTime
//							totalMoney = 0
//						} else {
//							currentTime = periodEndTime
//							lastEndTime = currentTime
//							surplusMoney, _ := decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
//							totalMoney -= surplusMoney
//						}
//						fmt.Println(totalMoney, halfHourPrice, "隔日====不足半小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					}
//					continue
//				}
//			}
//			// 当日 6
//			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End && totalMoney > 0 {
//				if totalMoney >= halfHourPrice && isExceed == 0 {
//					currentTime = currentTime.Add(time.Duration(30) * time.Minute)
//					lastEndTime = currentTime
//					totalMoney -= halfHourPrice
//					isExceed = 1
//					fmt.Println(totalMoney, halfHourPrice, "当日====满足半小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//				} else {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//
//					var minutePrice float64
//					var totalMinute int64
//					minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//					totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//					if totalMinute < int64(duration.Minutes()) {
//						fmt.Println("===|||今日", totalMinute, "...", int64(duration.Minutes()))
//						currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//						lastEndTime = currentTime
//						totalMoney = 0
//						fmt.Println(totalMoney, halfHourPrice, "今日====不足一小时按分钟计费=====1", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					} else {
//						currentTime = periodEndTime
//						lastEndTime = currentTime
//						surplusMoney, _ := decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
//						totalMoney -= surplusMoney
//						fmt.Println(totalMoney, halfHourPrice, "今日====不足一小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					}
//
//				}
//			}
//		}
//	}
//	fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime))
//	return "", ""
//}

// 金额算出会有误差，具体以结束计算为准(按半小时计费开台不足半小时按半时计费，超过半小时按分钟计费)
func (l *Charging) moneyTransferHalfOrHourAndMinuteTime(money float64, startDate string) (string, string) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	totalMoney := money
	isExceed := 0
	for {
		if totalMoney <= 0 {
			//fmt.Println("跳出循环")
			break
		}
		var periodEndTime time.Time
		var duration time.Duration
		for _, period := range l.periods {
			var halfOrHourPrice float64
			var minute int
			// 不足半小时计费按小时计费，超过半小时按分钟计费
			if l.chargingMode == 4 {
				halfOrHourPrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
				minute = 30
			}
			// 不足一小时计费按小时计费，超过一小时按分钟计费
			if l.chargingMode == 5 {
				halfOrHourPrice = period.Price
				minute = 60
			}
			// 跨日
			if period.Start > period.End && totalMoney > 0 {
				if currentTime.Format("15:04") >= period.Start || currentTime.Format("15:04") < period.End {
					if totalMoney >= halfOrHourPrice && isExceed == 0 {
						currentTime = currentTime.Add(time.Duration(minute) * time.Minute)
						lastEndTime = currentTime
						totalMoney -= halfOrHourPrice
						isExceed = 1
						//fmt.Println(totalMoney, halfOrHourPrice, "隔日====不足半小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
						continue
					}

					if isExceed == 1 {
						if currentTime.Format("15:04") >= period.Start {
							//fmt.Println("大于", period.Start, period.Price)
							periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), period.End), time.Local)
							duration = periodEndTime.Sub(currentTime)
						}
						if currentTime.Format("15:04") < period.End {
							//fmt.Println("小于", period.End, period.Price)
							periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
							duration = periodEndTime.Sub(currentTime)
						}
						var minutePrice float64
						var totalMinute int64
						minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
						totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()

						if totalMinute < int64(duration.Minutes()) {
							currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
							lastEndTime = currentTime
							totalMoney = 0
						} else {
							currentTime = periodEndTime
							lastEndTime = currentTime
							surplusMoney, _ := decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
							totalMoney -= surplusMoney
						}
						//fmt.Println(totalMoney, halfOrHourPrice, "隔日====不足半小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
					}
					continue
				}
			}
			// 当日 6
			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End && totalMoney > 0 {
				if totalMoney >= halfOrHourPrice && isExceed == 0 {
					currentTime = currentTime.Add(time.Duration(minute) * time.Minute)
					lastEndTime = currentTime
					totalMoney -= halfOrHourPrice
					isExceed = 1
					//fmt.Println(totalMoney, halfOrHourPrice, "当日====满足半小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
				} else {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)

					var minutePrice float64
					var totalMinute int64
					minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
					totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
					if totalMinute < int64(duration.Minutes()) {
						//fmt.Println("===|||今日", totalMinute, "...", int64(duration.Minutes()))
						currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
						lastEndTime = currentTime
						totalMoney = 0
						//fmt.Println(totalMoney, halfOrHourPrice, "今日====不足一小时按分钟计费=====1", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
					} else {
						currentTime = periodEndTime
						lastEndTime = currentTime
						surplusMoney, _ := decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
						totalMoney -= surplusMoney
						//fmt.Println(totalMoney, halfOrHourPrice, "今日====不足一小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
					}

				}
			}
		}
	}
	//fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime))
	return startDate, lastEndTime.Format("2006-01-02 15:04")
}

// 金额算出会有误差，具体以结束计算为准(以10分钟作为一个收费周期，第1分钟后开始计费)
func (l *Charging) moneyTransferCycleAndMinuteTime(money float64, startDate string) (string, string) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	totalMoney := money

	for {
		if totalMoney <= 0 {
			//fmt.Println("跳出循环")
			break
		}
		for _, period := range l.periods {
			var cyclePrice float64
			if l.cycleMinute == 30 {
				cyclePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
			}
			if l.cycleMinute == 60 {
				cyclePrice = period.Price
			}
			if l.cycleMinute != 30 && l.cycleMinute != 60 {
				minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
				cyclePrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
			}
			// 跨日
			if period.Start > period.End && totalMoney > 0 {
				if currentTime.Format("15:04") >= period.Start || currentTime.Format("15:04") < period.End {
					if totalMoney < cyclePrice {
						minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
						totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
						currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
						lastEndTime = currentTime
						totalMoney = 0
						//fmt.Println(totalMoney, cyclePrice, "隔日====小于周期内=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
					} else {
						currentTime = currentTime.Add(time.Duration(l.cycleMinute) * time.Minute)
						lastEndTime = currentTime
						totalMoney -= cyclePrice
						//fmt.Println(totalMoney, cyclePrice, "隔日====满足周期内=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
					}
				}
				continue
			}
			// 当日 6
			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End && totalMoney > 0 {
				if totalMoney < cyclePrice {
					//fmt.Println(totalMoney, cyclePrice, ".......")
					minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
					totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
					lastEndTime = currentTime
					totalMoney = 0
					//fmt.Println(totalMinute, cyclePrice, totalMinute, "当日====小于周期内=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
				} else {
					currentTime = currentTime.Add(time.Duration(l.cycleMinute) * time.Minute)
					lastEndTime = currentTime
					totalMoney -= cyclePrice
					//fmt.Println(totalMoney, cyclePrice, "当日====满足周期内=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
				}
			}
		}
	}
	//fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime))
	return startDate, lastEndTime.Format("2006-01-02 15:04")
}

// 金额转换成开始及结束时间
func (l *Charging) MoneyTransfer(money float64, startDateParam string) (string, string) {
	var startData string
	var endData string
	// 分钟计费
	if l.chargingMode == 1 {
		startData, endData = l.moneyTransferMinuteTime(money, startDateParam)
	}
	// 半小时计费
	if l.chargingMode == 2 {
		startData, endData = l.moneyTransferHalfHourTime(money, startDateParam)
	}
	// 小时计费
	if l.chargingMode == 3 {
		startData, endData = l.moneyTransferHourTime(money, startDateParam)
	}
	// 按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
	if l.chargingMode == 4 {
		startData, endData = l.moneyTransferHalfOrHourAndMinuteTime(money, startDateParam)
	}
	// 5按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
	if l.chargingMode == 5 {
		startData, endData = l.moneyTransferHalfOrHourAndMinuteTime(money, startDateParam)
	}
	// 6自定义计费
	if l.chargingMode == 6 {
		startData, endData = l.moneyTransferCycleAndMinuteTime(money, startDateParam)
	}
	return startData, endData
}
