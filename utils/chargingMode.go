package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"sort"
	"time"
)

// 收费时段设置
type ChargePeriodAssembly struct {
	Week    *ChargePeriodSetWeek    `json:"week"`    // 星期
	Holiday *ChargePeriodSetHoliday `json:"holiday"` // 节假日
	Hour    *ChargePeriodSetHour    `json:"hour"`    // 小时收费时段
}

// 星期
type ChargePeriodSetWeek struct {
	Week       []int64                 `json:"week"`        // 星期
	Hour       []CPHour                `json:"hour"`        // 收费时段
	HourPeak   []CPHour                `json:"hour_peak"`   // 时段封顶
	HourStairs []HourStairs            `json:"hour_stairs"` // 时段阶梯
	MinConsume *ChargePeriodMinConsume `json:"min_consume"`
}

// 小时收费时段
type ChargePeriodSetHour struct {
	Hour       []CPHour                `json:"hour"`        // 收费时段
	HourPeak   []CPHour                `json:"hour_peak"`   // 时段封顶
	HourStairs []HourStairs            `json:"hour_stairs"` // 时段阶梯
	MinConsume *ChargePeriodMinConsume `json:"min_consume"`
}

type ChargePeriodMinConsume struct {
	StartNoFeeTime int64  `json:"start_no_fee_time"`
	Idle           string `json:"idle"`
	Member         string `json:"member"`
}
type CPHour struct {
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	IdlePrice   string `json:"idle_price"`
	MemberPrice string `json:"member_price"`
}

type HourStairs struct {
	Index       int64  `json:"index"`
	IdlePrice   string `json:"idle_price"`
	MemberPrice string `json:"member_price"`
}

// 节假日
type ChargePeriodSetHoliday struct {
	Date       [][]string               `json:"date"`        // 星期
	Hour       [][]CPHour               `json:"hour"`        // 收费时段
	HourPeak   [][]CPHour               `json:"hour_peak"`   // 时段封顶
	HourStairs [][]HourStairs           `json:"hour_stairs"` // 时段阶梯
	MinConsume []ChargePeriodMinConsume `json:"min_consume"`
}

type ChargePeriod struct {
	StartPeriod int64   `json:"start_period"`
	EndPeriod   int64   `json:"end_period"`
	Start       string  `json:"start"`
	End         string  `json:"end"`
	Price       float64 `json:"price"`
}

type Charging struct {
	periods      map[int64]ChargePeriod
	chargingMode int64 // 计费方式
	cycleMinute  int64 // 计费周期
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

func NewChargeMode(periods map[int64]ChargePeriod, chargingMode int64, cycleMinute int64) *Charging {
	return &Charging{periods: periods, chargingMode: chargingMode, cycleMinute: cycleMinute}
}

// 填充时间
func (l *Charging) fillTime(startDateTime, endDateTime time.Time) time.Time {
	// 半小时计费，如果时间小于半小时则初始设置半个小时
	if ContainsSliceInt64([]int64{2, 3, 6, 7}, l.chargingMode) {
		tempMinute := decimal.NewFromFloat(endDateTime.Sub(startDateTime).Minutes()).IntPart()
		if tempMinute < 30 {
			endDateTime = endDateTime.Add(time.Duration(30-tempMinute) * time.Minute)
		}
	}
	// 一小时计费，如果时间小于一小时则初始设置一个小时
	if ContainsSliceInt64([]int64{4, 5, 8, 9}, l.chargingMode) {
		tempMinute := decimal.NewFromFloat(endDateTime.Sub(startDateTime).Minutes()).IntPart()
		if tempMinute < 60 {
			endDateTime = endDateTime.Add(time.Duration(60-tempMinute) * time.Minute)
		}
	}
	// 以10分钟作为一个收费周期，第1分钟后开始计费(自定义计费方式)
	if ContainsSliceInt64([]int64{10, 11}, l.chargingMode) {
		tempMinute := decimal.NewFromFloat(endDateTime.Sub(startDateTime).Minutes()).IntPart()
		cycleNum := decimal.NewFromInt(tempMinute).Div(decimal.NewFromInt(l.cycleMinute)).Ceil().IntPart()
		maxMinute := decimal.NewFromInt(cycleNum).Mul(decimal.NewFromInt(l.cycleMinute)).IntPart()
		endDateTime = startDateTime.Add(time.Duration(maxMinute) * time.Minute)
	}
	return endDateTime
}

// 计费
func (l *Charging) Outlay(startDate, endDate string) (float64, map[int64]float64, []*PeriodList) {
	var price float64
	periodList := make([]*PeriodList, 0)
	periods := make(map[int64]float64) // 时段费用详情
	startDateTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	endDateTime, _ := time.ParseInLocation("2006-01-02 15:04", endDate, time.Local)
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	endDateTime = l.fillTime(startDateTime, endDateTime) // 填充时间

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
	// 3按半小时计费(跨时段)
	if l.chargingMode == 3 {
		price, periods, periodList = l.computeHalfHourModeTwo(times)
	}
	// 4按小时计费
	if l.chargingMode == 4 {
		price, periods, periodList = l.computeHour(times)
	}
	// 5按小时计费(跨时段)
	if l.chargingMode == 5 {
		price, periods, periodList = l.computeHourModeTwo(times)
	}
	// 6按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
	if l.chargingMode == 6 {
		price, periods, periodList = l.computeHalfHourOrMinute(times)
	}
	// 7按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费(跨时段)
	if l.chargingMode == 7 {
		price, periods, periodList = l.computeHalfHourOrMinuteModeTwo(times)
	}
	// 8按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
	if l.chargingMode == 8 {
		price, periods, periodList = l.computeHourOrMinute(times)
	}
	// 9按小时计费开台不足1小时按小时计费，超过1小时按分钟计费(跨时段)
	if l.chargingMode == 9 {
		price, periods, periodList = l.computeHourOrMinuteModeTwo(times)
	}
	// 10 以10分钟作为一个收费周期，第1分钟后开始计费
	if l.chargingMode == 10 {
		price, periods, periodList = l.computeCycleAndMinute(times)
	}
	// 11 以10分钟作为一个收费周期，第1分钟后开始计费(跨时段)
	if l.chargingMode == 11 {
		price, periods, periodList = l.computeCycleAndMinuteModeTwo(times)
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

			//price += totalPrice
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

// 以10分钟作为一个收费周期，第1分钟后开始计费，如果时间段内计算时长小于收费周期，累计至下个时段
// 并在下个时段内扣除上个时段少于的时间[因为以收一个时段价格]，然后在计算下周期时段价格，依次类推)
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

// 半小时计费 (半小时为节点，以半小时起始点计算价格，如果时间段内计算时长小于半小时，则累计至下个时段，
// 在下个时段内扣除上个时段少于的时间[如果跨时段按跨时段分钟计费,否则按半小时算]，然后在计算下周期时段价格，依次类推)
func (l *Charging) computeHalfHourModeTwo(times []periodTimes) (float64, map[int64]float64, []*PeriodList) {
	periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	periodList := make([]*PeriodList, 0)
	for _, timeV := range times {
		var totalPrice float64
		// 差值按分钟计算计费
		if lastDiffMinute > 0 {
			lastStartTime := timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
			if lastStartTime.Unix() >= timeV.EndTime.Unix() {
				lastStartTime = timeV.EndTime
			}
			duration := lastStartTime.Sub(timeV.StartTime)
			totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
			//fmt.Println(".......-----", lastDiffMinute, duration, totalMinute)
			if totalMinute > 0 {
				minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
				minuteTotalPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
				timeV.StartTime = lastStartTime
				//fmt.Println("跨时段按分钟", lastDiffMinute, duration, totalMinute, totalPrice, timeV.StartTime.Format(time.DateTime))
			}
		}
		duration := timeV.EndTime.Sub(timeV.StartTime)
		totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
		if totalMinute > 0 {
			haltHour := decimal.NewFromInt(totalMinute).Div(decimal.NewFromInt(30)).Ceil().IntPart()
			maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
			var tempEndTime time.Time
			if timeV.Start > timeV.End && timeV.StartTime.Format("15:04") >= timeV.Start {
				tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", timeV.EndTime.Add(time.Duration(86400)*time.Second).Format("2006-01-02"), timeV.End), time.Local)
				//} else if timeV.Start > timeV.End && timeV.StartTime.Format("15:04") >= "00:00" {
				//	tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", timeV.EndTime.Format("2006-01-02"), timeV.End), time.Local)
			} else {
				tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", timeV.EndTime.Format("2006-01-02"), timeV.End), time.Local)
			}

			//fmt.Println("..;;;;;;;;;;", timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format(time.DateTime), tempEndTime.Format(time.DateTime))

			if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Unix() >= tempEndTime.Unix() {
				lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(totalMinute)).IntPart()
				//fmt.Println(lastDiffMinute, "....;;;;;;;;;")
			} else {
				lastDiffMinute = 0
			}
			//fmt.Println(lastDiffMinute, "lastDiffMinute====", timeV.StartTime.Format(time.DateTime), timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format(time.DateTime), timeV.End, ".....", tempEndTime.Format(time.DateTime))
			//fmt.Println(lastDiffMinute, "======222")
			halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
			// 开始时间大于结束时间 跨天
			if timeV.Start > timeV.End {
				//fmt.Println("1111")
				// 跨时段
				if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format("15:04") >= timeV.End && lastDiffMinute != 0 {
					//minute := decimal.NewFromInt(30).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
					if totalMinute < 30 && timeV.EndTime.Format("15:04") < timeV.End {
						//fmt.Println("22222")
						totalHourPrice, _ := decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour)).Float64()
						totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
						//fmt.Println("不足半小时按半小时算", totalPrice)
					} else {
						//fmt.Println("33333")
						minute := decimal.NewFromInt(30).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
						minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
						minuteTotalPrice, _ := decimal.NewFromInt(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()

						totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
						totalHourPrice, _ := decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour - 1)).Float64()
						totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
					}
					//	fmt.Println("跨时段", minuteTotalPrice, totalPrice, lastDiffMinute)
				} else {
					//fmt.Println("444444")
					totalHourPrice, _ := decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour)).Float64()
					//fmt.Println("未跨时段", totalPrice)
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
				}
			} else {
				//fmt.Println("55555")
				//fmt.Println(haltHour, halfHourPrice, lastDiffMinute, 30-lastDiffMinute)
				//fmt.Println("-----000", timeV)
				if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format("15:04") >= timeV.End && lastDiffMinute != 0 {
					//fmt.Println("666666")
					minute := decimal.NewFromInt(30).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
					minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
					minuteTotalPrice, _ := decimal.NewFromInt(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//fmt.Println(minute, minutePrice, minuteTotalPrice, ".........1111=====")
					totalHourPrice, _ := decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour - 1)).Float64()
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
				} else {
					//fmt.Println("77777777")
					totalHourPrice, _ := decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour)).Float64()
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
					//fmt.Println("未跨时段 当日", halfHourPrice, haltHour, totalPrice)
				}
			}
			//fmt.Println("===================================")
			price += totalPrice
			periods[int64(timeV.Index)] += totalPrice
		} else {
			//fmt.Println("---11")
			periods[int64(timeV.Index)] += 0
		}
		//fmt.Println(price, "......////")
		periodList = append(periodList, &PeriodList{
			StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
			EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
			Index:      timeV.Index,
			Duration:   totalMinute,
			Start:      timeV.Start,
			End:        timeV.End,
			Price:      timeV.Price,
			TotalPrice: totalPrice,
		})

	}
	//fmt.Println(price, "..----,,,")
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, periods, periodList
}

// 一小时计费 (一小时为节点，以一小时起始点计算价格，如果时间段内计算时长小于一小时，则累计至下个时段，
// 在下个时段内扣除上个时段少于的时间[如果跨时段按跨时段分钟计费,否则按一小时算]，然后在计算下周期时段价格，依次类推)
func (l *Charging) computeHourModeTwo(times []periodTimes) (float64, map[int64]float64, []*PeriodList) {
	periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	periodList := make([]*PeriodList, 0)
	for _, timeV := range times {
		var totalPrice float64
		// 差值按分钟计算计费
		if lastDiffMinute > 0 {
			lastStartTime := timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
			if lastStartTime.Unix() >= timeV.EndTime.Unix() {
				lastStartTime = timeV.EndTime
			}
			duration := lastStartTime.Sub(timeV.StartTime)
			totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
			//fmt.Println(".......-----", lastDiffMinute, duration, totalMinute)
			if totalMinute > 0 {
				minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
				minuteTotalPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
				timeV.StartTime = lastStartTime
				//fmt.Println("跨时段按分钟", lastDiffMinute, duration, totalMinute, totalPrice, timeV.StartTime.Format(time.DateTime))
			}
		}
		duration := timeV.EndTime.Sub(timeV.StartTime)
		totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
		if totalMinute > 0 {
			hour := decimal.NewFromInt(totalMinute).Div(decimal.NewFromInt(60)).Ceil().IntPart()
			maxMinute := decimal.NewFromInt(hour).Mul(decimal.NewFromInt(60)).IntPart()
			var tempEndTime time.Time
			if timeV.Start > timeV.End && timeV.StartTime.Format("15:04") >= timeV.Start {
				tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", timeV.EndTime.Add(time.Duration(86400)*time.Second).Format("2006-01-02"), timeV.End), time.Local)
			} else {
				tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", timeV.EndTime.Format("2006-01-02"), timeV.End), time.Local)
			}

			//fmt.Println("..;;;;;;;;;;", timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format(time.DateTime), tempEndTime.Format(time.DateTime))

			if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Unix() >= tempEndTime.Unix() {
				lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(totalMinute)).IntPart()
				//fmt.Println(lastDiffMinute, "....;;;;;;;;;")
			} else {
				lastDiffMinute = 0
			}
			//fmt.Println(lastDiffMinute, "lastDiffMinute====", timeV.StartTime.Format(time.DateTime), timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format(time.DateTime), timeV.End, ".....", tempEndTime.Format(time.DateTime))
			//fmt.Println(lastDiffMinute, "======222")
			//halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
			// 开始时间大于结束时间 跨天
			if timeV.Start > timeV.End {
				//fmt.Println("1111")
				// 跨时段
				if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format("15:04") >= timeV.End && lastDiffMinute != 0 {
					if totalMinute < 60 && timeV.EndTime.Format("15:04") < timeV.End {
						//fmt.Println("22222")
						totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour)).Float64()
						totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
						//fmt.Println("不足一小时按一小时算", totalPrice)
					} else {
						//fmt.Println("33333")
						minute := decimal.NewFromInt(60).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
						minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
						minuteTotalPrice, _ := decimal.NewFromInt(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()

						totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
						totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour - 1)).Float64()
						totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
					}
					//	fmt.Println("跨时段", minuteTotalPrice, totalPrice, lastDiffMinute)
				} else {
					//fmt.Println("444444")
					totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour)).Float64()
					//fmt.Println("未跨时段", totalPrice)
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
				}
			} else {
				//fmt.Println("55555")
				//fmt.Println(haltHour, halfHourPrice, lastDiffMinute, 30-lastDiffMinute)
				//fmt.Println("-----000", timeV)
				if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format("15:04") >= timeV.End && lastDiffMinute != 0 {
					//fmt.Println("666666")
					minute := decimal.NewFromInt(60).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
					minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
					minuteTotalPrice, _ := decimal.NewFromInt(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//fmt.Println(minute, minutePrice, minuteTotalPrice, ".........1111=====")
					totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour - 1)).Float64()
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
				} else {
					//fmt.Println("77777777")
					totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour)).Float64()
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
					//fmt.Println("未跨时段 当日", halfHourPrice, haltHour, totalPrice)
				}
			}
			//fmt.Println("===================================")
			price += totalPrice
			periods[int64(timeV.Index)] += totalPrice
		} else {
			//fmt.Println("---11")
			periods[int64(timeV.Index)] += 0
		}
		//fmt.Println(price, "......////")
		periodList = append(periodList, &PeriodList{
			StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
			EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
			Index:      timeV.Index,
			Duration:   totalMinute,
			Start:      timeV.Start,
			End:        timeV.End,
			Price:      timeV.Price,
			TotalPrice: totalPrice,
		})

	}
	//fmt.Println(price, "..----,,,")
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, periods, periodList
}

// 按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费。[如果跨时段按跨时段分钟计费]
func (l *Charging) computeHalfHourOrMinuteModeTwo(times []periodTimes) (float64, map[int64]float64, []*PeriodList) {
	periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	var firstHalfHour int64
	firstHalfHour = 30
	periodList := make([]*PeriodList, 0)
	for _, timeV := range times {
		var totalPrice float64
		// 差值按分钟计算计费
		if lastDiffMinute > 0 {
			lastStartTime := timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
			if lastStartTime.Unix() >= timeV.EndTime.Unix() {
				lastStartTime = timeV.EndTime
			}
			duration := lastStartTime.Sub(timeV.StartTime)
			totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
			if totalMinute > 0 {
				minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
				minuteTotalPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
				timeV.StartTime = lastStartTime
				lastDiffMinute = 0
				//fmt.Println("跨时段按分钟", lastDiffMinute, duration, totalMinute, totalPrice, timeV.StartTime.Format(time.DateTime))
			}
		}

		duration := timeV.EndTime.Sub(timeV.StartTime)
		totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
		var surplusMinute int64
		if totalMinute > 0 {
			if firstHalfHour != 0 {
				haltHour := decimal.NewFromInt(totalMinute).Div(decimal.NewFromInt(30)).Ceil().IntPart()
				maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
				halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
				if haltHour >= 1 && totalMinute >= 30 {
					totalHourPrice, _ := decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(1)).Float64()
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
					surplusMinute = totalMinute - 30
				}
				if haltHour == 1 && totalMinute < 30 {
					lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(totalMinute)).IntPart()
					surplusMinute = decimal.NewFromInt(30).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
					//fmt.Println("检查下个时段是否有数据，没有则在当前时段计费下个时段", index, "===", times[index+1])
				}
				if surplusMinute >= 1 {
					minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
					minuteTotalPrice, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
				}
				//price += totalPrice
				periods[int64(timeV.Index)] += totalPrice
				firstHalfHour = 0 // 首个半小时计费，计费完后置为0
			} else {
				minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
				minuteTotalPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
				//price += totalPrice
				periods[int64(timeV.Index)] += totalPrice
				lastDiffMinute = 0 // 分钟计费，差值分钟置为0
			}
		} else {
			//price += totalPrice
			periods[int64(timeV.Index)] += 0
		}
		price += totalPrice
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

// 按小时计费开台不足1小时按小时计费，超过1小时按分钟计费。[如果跨时段按跨时段分钟计费]
func (l *Charging) computeHourOrMinuteModeTwo(times []periodTimes) (float64, map[int64]float64, []*PeriodList) {
	periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	var firstHalfHour int64
	firstHalfHour = 60
	periodList := make([]*PeriodList, 0)
	for _, timeV := range times {
		var totalPrice float64
		// 差值按分钟计算计费
		if lastDiffMinute > 0 {
			lastStartTime := timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
			if lastStartTime.Unix() >= timeV.EndTime.Unix() {
				lastStartTime = timeV.EndTime
			}
			duration := lastStartTime.Sub(timeV.StartTime)
			totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
			if totalMinute > 0 {
				minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
				minuteTotalPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
				timeV.StartTime = lastStartTime
				lastDiffMinute = 0
				//fmt.Println("跨时段按分钟", lastDiffMinute, duration, totalMinute, totalPrice, timeV.StartTime.Format(time.DateTime))
			}
		}

		duration := timeV.EndTime.Sub(timeV.StartTime)
		totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
		var surplusMinute int64
		if totalMinute > 0 {
			if firstHalfHour != 0 {
				hour := decimal.NewFromInt(totalMinute).Div(decimal.NewFromInt(60)).Ceil().IntPart()
				maxMinute := decimal.NewFromInt(hour).Mul(decimal.NewFromInt(60)).IntPart()
				if hour >= 1 && totalMinute >= 60 {
					totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(1)).Float64()
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
					surplusMinute = totalMinute - 60
				}
				if hour == 1 && totalMinute < 60 {
					lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(totalMinute)).IntPart()
					surplusMinute = decimal.NewFromInt(60).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
					//fmt.Println("检查下个时段是否有数据，没有则在当前时段计费下个时段", index, "===", times[index+1])
				}
				if surplusMinute >= 1 {
					minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
					minuteTotalPrice, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
				}
				//price += totalPrice
				periods[int64(timeV.Index)] += totalPrice
				firstHalfHour = 0 // 首个半小时计费，计费完后置为0
			} else {
				minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
				minuteTotalPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
				//price += totalPrice
				periods[int64(timeV.Index)] += totalPrice
				lastDiffMinute = 0 // 分钟计费，差值分钟置为0
			}
		} else {
			//price += totalPrice
			periods[int64(timeV.Index)] += 0
		}
		price += totalPrice
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

// 以10分钟作为一个收费周期，第1分钟后开始计费，如果时间段内计算时长小于收费周期，累计至下个时段
// 并在下个时段内扣除上个时段少于的时间[如果跨时段按跨时段分钟计费]，然后在计算下周期时段价格，依次类推)
func (l *Charging) computeCycleAndMinuteModeTwo(times []periodTimes) (float64, map[int64]float64, []*PeriodList) {
	periods := make(map[int64]float64) // 时段费用详情
	var price float64
	periodList := make([]*PeriodList, 0)
	for _, timeV := range times {
		var totalPrice float64

		duration := timeV.EndTime.Sub(timeV.StartTime)
		totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
		var surplusMinute int64
		var cyclePrice float64
		if totalMinute > 0 {
			minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
			cycle := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(l.cycleMinute)).Ceil().IntPart()

			// 半小时
			if l.cycleMinute == 30 {
				cyclePrice, _ = decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
			}
			// 一小时
			if l.cycleMinute == 60 {
				cyclePrice = timeV.Price
			}
			// 周期计费
			if l.cycleMinute != 30 && l.cycleMinute != 60 {
				cyclePrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
			}

			if cycle >= 1 && totalMinute >= l.cycleMinute {
				totalHourPrice, _ := decimal.NewFromFloat(cyclePrice).Mul(decimal.NewFromInt(cycle - 1)).Float64()
				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
				surplusMinute = totalMinute - ((cycle - 1) * l.cycleMinute)
			}
			if cycle == 1 && totalMinute < l.cycleMinute {
				totalHourPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
				//fmt.Println("检查下个时段是否有数据，没有则在当前时段计费下个时段", index, "===", times[index+1])
				surplusMinute = 0
			}
			if surplusMinute >= 1 {
				minuteTotalPrice, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
			}
			periods[int64(timeV.Index)] += totalPrice
		} else {
			periods[int64(timeV.Index)] += 0
		}
		price += totalPrice
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

func (l *Charging) sortPeriods() []int {
	// 将 map 的键放入切片
	keys := make([]int, 0, len(l.periods))
	for k := range l.periods {
		keys = append(keys, int(k))
	}
	// 对键切片进行排序
	sort.Ints(keys)
	return keys
}

// 金额算出会有误差，具体以结束计算为准(分钟)
func (l *Charging) moneyTransferMinuteTime(money float64, startDate string) (string, string, []*PeriodList) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	totalMoney := money
	periodList := make([]*PeriodList, 0)

	// 将 map 的键放入切片
	periods := l.sortPeriods()

	for {
		if totalMoney <= 0 {
			//fmt.Println("跳出循环")
			break
		}
		var periodEndTime time.Time
		var duration time.Duration

		for _, index := range periods {
			period := l.periods[int64(index)]
			//}
			//for id, period := range l.periods {
			//	fmt.Println("~~~~~~~~~~~~~~", id, period.Start, period.End)
			var minutePrice float64
			var totalMinute int64
			if totalMoney <= 0 {
				continue
			}
			if currentTime.Format("15:04") == period.End {
				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
				continue
			}
			minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()

			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
			// 判断 如果小于1分钟价格则不取整
			if subTotalMinute < 0.001 {
				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			}

			// 跨日
			if period.Start > period.End {
				tempCurrentTime := currentTime
				//fmt.Println("===========1111", currentTime.Format(time.DateTime), currentTime.Format("15:04"), period.End, period.Start)
				if currentTime.Format("15:04") >= period.Start {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < period.End {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
				if totalMinute < durationMinute {
					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
					lastEndTime = currentTime
					totalMoney = 0
				} else {
					currentTime = periodEndTime
					lastEndTime = currentTime
					surplusMoney, _ := decimal.NewFromInt(durationMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					totalMoney -= surplusMoney
				}
				Tduration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
					Index:      index,
					Duration:   Tduration,
					Start:      period.Start,
					End:        period.End,
					Price:      period.Price,
					TotalPrice: totalMoney,
				})
				continue
			}
			// 当日
			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End {
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
				duration = periodEndTime.Sub(currentTime)
				durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
				tempCurrentTime := currentTime
				if totalMinute < durationMinute {
					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
					lastEndTime = currentTime
					totalMoney = 0
				} else {
					currentTime = periodEndTime
					lastEndTime = currentTime
					surplusMoney, _ := decimal.NewFromInt(durationMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					totalMoney -= surplusMoney
				}

				Tduration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
					Index:      index,
					Duration:   Tduration,
					Start:      period.Start,
					End:        period.End,
					Price:      period.Price,
					TotalPrice: totalMoney,
				})
				//fmt.Println("当日")
			}
		}
	}
	//fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime))
	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
}

// 金额算出会有误差，具体以结束计算为准(半小时)[不足半小时或跨时段，剩余金额则按分钟计算]
func (l *Charging) moneyTransferHalfHourTime(money float64, startDate string) (string, string, []*PeriodList) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	periodList := make([]*PeriodList, 0)
	periods := l.sortPeriods()
	totalMoney := money
	for {
		if totalMoney <= 0 {
			break
		}
		var periodEndTime time.Time
		var duration time.Duration

		for _, index := range periods {
			period := l.periods[int64(index)]
			var minutePrice float64
			var totalMinute int64
			if totalMoney <= 0 {
				continue
			}
			if currentTime.Format("15:04") == period.End {
				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
				continue
			}
			tempCurrentTime := currentTime
			minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
			halfHourPrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()

			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
			// 判断 如果小于1分钟价格则不取整
			if subTotalMinute < 0.001 {
				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			}

			// 跨日
			if period.Start > period.End {
				if currentTime.Format("15:04") >= period.Start {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < period.End {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
					periodEndTime = tempTotalEndTime
					duration = periodEndTime.Sub(currentTime)
				}
				lastEndTime = periodEndTime
				currentTime = periodEndTime
			}
			// 当日
			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End {
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
				duration = periodEndTime.Sub(currentTime)
				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
					periodEndTime = tempTotalEndTime
					duration = periodEndTime.Sub(currentTime)
				}
				lastEndTime = periodEndTime
				currentTime = periodEndTime
			}
			durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
			if durationMinute >= 1 {
				haltHour := decimal.NewFromInt(durationMinute).Div(decimal.NewFromInt(30)).Ceil().IntPart()
				maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
				subMinute := decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(durationMinute)).IntPart()
				surplusMinute := decimal.NewFromInt(30).Sub(decimal.NewFromInt(subMinute)).IntPart()
				surplusMoney, _ := decimal.NewFromInt(haltHour - 1).Mul(decimal.NewFromFloat(halfHourPrice)).Float64()
				surplusMinuteMoney, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				surplusMoney, _ = decimal.NewFromFloat(surplusMoney).Add(decimal.NewFromFloat(surplusMinuteMoney)).Float64()
				totalMoney -= surplusMoney
				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
				if totalMoney < 0.001 {
					totalMoney = 0
				}

				Tduration := decimal.NewFromInt(durationMinute).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
					Index:      index,
					Duration:   Tduration,
					Start:      period.Start,
					End:        period.End,
					Price:      period.Price,
					TotalPrice: totalMoney,
				})
			}
		}
	}
	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
}

// 金额算出会有误差，具体以结束计算为准(一小时)[不足一小时或跨时段，剩余金额则按分钟计算]
func (l *Charging) moneyTransferHourTime(money float64, startDate string) (string, string, []*PeriodList) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	periodList := make([]*PeriodList, 0)
	periods := l.sortPeriods()
	totalMoney := money
	for {
		if totalMoney <= 0 {
			break
		}
		var periodEndTime time.Time
		var duration time.Duration

		for _, index := range periods {
			period := l.periods[int64(index)]
			var minutePrice float64
			var totalMinute int64
			if totalMoney <= 0 {
				continue
			}
			if currentTime.Format("15:04") == period.End {
				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
				continue
			}

			tempCurrentTime := currentTime
			minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()

			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
			// 判断 如果小于1分钟价格则不取整
			if subTotalMinute < 0.001 {
				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			}

			// 跨日
			if period.Start > period.End {
				if currentTime.Format("15:04") >= period.Start {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < period.End {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
					periodEndTime = tempTotalEndTime
					duration = periodEndTime.Sub(currentTime)
				}
				lastEndTime = periodEndTime
				currentTime = periodEndTime
			}
			// 当日
			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End {
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
				duration = periodEndTime.Sub(currentTime)
				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
					periodEndTime = tempTotalEndTime
					duration = periodEndTime.Sub(currentTime)
				}
				lastEndTime = periodEndTime
				currentTime = periodEndTime
			}
			durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
			if durationMinute >= 1 {
				haltHour := decimal.NewFromInt(durationMinute).Div(decimal.NewFromInt(60)).Ceil().IntPart()
				maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(60)).IntPart()
				subMinute := decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(durationMinute)).IntPart()
				surplusMinute := decimal.NewFromInt(60).Sub(decimal.NewFromInt(subMinute)).IntPart()
				surplusMoney, _ := decimal.NewFromInt(haltHour - 1).Mul(decimal.NewFromFloat(period.Price)).Float64()
				surplusMinuteMoney, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				surplusMoney, _ = decimal.NewFromFloat(surplusMoney).Add(decimal.NewFromFloat(surplusMinuteMoney)).Float64()
				totalMoney -= surplusMoney
				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
				if totalMoney < 0.001 {
					totalMoney = 0
				}

				Tduration := decimal.NewFromInt(durationMinute).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
					Index:      index,
					Duration:   Tduration,
					Start:      period.Start,
					End:        period.End,
					Price:      period.Price,
					TotalPrice: totalMoney,
				})
			}
		}
	}
	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
}

// 金额算出会有误差，具体以结束计算为准(半小时\一小时)[不足半小时\一小时或跨时段，剩余金额则按分钟计算]
func (l *Charging) moneyTransferHalfOrHourAndMinuteTime(money float64, startDate string) (string, string, []*PeriodList) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	periodList := make([]*PeriodList, 0)
	periods := l.sortPeriods()
	totalMoney := money
	for {
		if totalMoney <= 0 {
			break
		}
		var periodEndTime time.Time
		var duration time.Duration

		for _, index := range periods {
			period := l.periods[int64(index)]
			var minutePrice float64
			var totalMinute int64
			var halfOrHourPrice float64
			var minute int64
			if totalMoney <= 0 {
				continue
			}
			//halfHourPrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
			if currentTime.Format("15:04") == period.End {
				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
				continue
			}
			tempCurrentTime := currentTime
			minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()

			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
			// 判断 如果小于1分钟价格则不取整
			if subTotalMinute < 0.001 {
				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			}

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
			if period.Start > period.End {
				if currentTime.Format("15:04") >= period.Start {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < period.End {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
					periodEndTime = tempTotalEndTime
					duration = periodEndTime.Sub(currentTime)
				}
				lastEndTime = periodEndTime
				currentTime = periodEndTime
			}
			// 当日
			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End {
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
				duration = periodEndTime.Sub(currentTime)
				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
					periodEndTime = tempTotalEndTime
					duration = periodEndTime.Sub(currentTime)
				}
				lastEndTime = periodEndTime
				currentTime = periodEndTime
			}
			durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
			if durationMinute >= 1 {
				haltHour := decimal.NewFromInt(durationMinute).Div(decimal.NewFromInt(minute)).Ceil().IntPart()
				maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(minute)).IntPart()
				subMinute := decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(durationMinute)).IntPart()
				surplusMinute := decimal.NewFromInt(minute).Sub(decimal.NewFromInt(subMinute)).IntPart()
				surplusMoney, _ := decimal.NewFromInt(haltHour - 1).Mul(decimal.NewFromFloat(halfOrHourPrice)).Float64()
				surplusMinuteMoney, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				surplusMoney, _ = decimal.NewFromFloat(surplusMoney).Add(decimal.NewFromFloat(surplusMinuteMoney)).Float64()
				totalMoney -= surplusMoney
				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
				if totalMoney < 0.001 {
					totalMoney = 0
				}

				Tduration := decimal.NewFromInt(durationMinute).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
					Index:      index,
					Duration:   Tduration,
					Start:      period.Start,
					End:        period.End,
					Price:      period.Price,
					TotalPrice: totalMoney,
				})
			}
		}
	}
	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
}

// 金额算出会有误差，具体以结束计算为准(以10分钟作为一个收费周期，第1分钟后开始计费)[不足周期计费或跨时段，剩余金额则按分钟计算]
func (l *Charging) moneyTransferCycleAndMinuteTime(money float64, startDate string) (string, string, []*PeriodList) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	periodList := make([]*PeriodList, 0)
	periods := l.sortPeriods()
	totalMoney := money
	for {
		if totalMoney <= 0 {
			break
		}
		var periodEndTime time.Time
		var duration time.Duration

		for _, index := range periods {
			period := l.periods[int64(index)]
			var minutePrice float64
			var totalMinute int64
			//var halfOrHourPrice float64
			//var minute int64
			var cyclePrice float64
			if totalMoney <= 0 {
				continue
			}
			//halfHourPrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
			if currentTime.Format("15:04") == period.End {
				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
				continue
			}
			tempCurrentTime := currentTime
			minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()

			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
			// 判断 如果小于1分钟价格则不取整
			if subTotalMinute < 0.001 {
				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			}

			if l.cycleMinute == 30 {
				cyclePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
			}
			if l.cycleMinute == 60 {
				cyclePrice = period.Price
			}
			if l.cycleMinute != 30 && l.cycleMinute != 60 {
				cyclePrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
			}

			// 跨日
			if period.Start > period.End {
				if currentTime.Format("15:04") >= period.Start {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < period.End {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
					periodEndTime = tempTotalEndTime
					duration = periodEndTime.Sub(currentTime)
				}
				lastEndTime = periodEndTime
				currentTime = periodEndTime
			}
			// 当日
			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End {
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
				duration = periodEndTime.Sub(currentTime)
				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
					periodEndTime = tempTotalEndTime
					duration = periodEndTime.Sub(currentTime)
				}
				lastEndTime = periodEndTime
				currentTime = periodEndTime
			}
			durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
			if durationMinute >= 1 {
				haltHour := decimal.NewFromInt(durationMinute).Div(decimal.NewFromInt(l.cycleMinute)).Ceil().IntPart()
				maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(l.cycleMinute)).IntPart()
				subMinute := decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(durationMinute)).IntPart()
				surplusMinute := decimal.NewFromInt(l.cycleMinute).Sub(decimal.NewFromInt(subMinute)).IntPart()
				surplusMoney, _ := decimal.NewFromInt(haltHour - 1).Mul(decimal.NewFromFloat(cyclePrice)).Float64()
				surplusMinuteMoney, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				surplusMoney, _ = decimal.NewFromFloat(surplusMoney).Add(decimal.NewFromFloat(surplusMinuteMoney)).Float64()
				totalMoney -= surplusMoney
				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
				if totalMoney < 0.001 {
					totalMoney = 0
				}

				Tduration := decimal.NewFromInt(durationMinute).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
					Index:      index,
					Duration:   Tduration,
					Start:      period.Start,
					End:        period.End,
					Price:      period.Price,
					TotalPrice: totalMoney,
				})
			}
		}
	}
	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
}

// 金额转换成开始及结束时间
func (l *Charging) MoneyTransfer(money float64, startDateParam string) (string, string, []*PeriodList) {
	var startData string
	var endData string
	periodList := make([]*PeriodList, 0)

	//计费模式 1按分钟计费 2按半小时计费 3按半小时计费(跨时段) 4按小时计费 5按小时计费(跨时段)
	//6按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费  7按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费(跨时段)
	//8按小时计费开台不足1小时按小时计费，超过1小时按分钟计费   9按小时计费开台不足1小时按小时计费，超过1小时按分钟计费(跨时段)
	//10自定义计费 11自定义计费(跨时段)

	// 分钟计费
	if l.chargingMode == 1 {
		startData, endData, periodList = l.moneyTransferMinuteTime(money, startDateParam)
	}
	// 半小时计费
	if ContainsSliceInt64([]int64{2, 3}, l.chargingMode) {
		startData, endData, periodList = l.moneyTransferHalfHourTime(money, startDateParam)
	}
	// 小时计费
	if ContainsSliceInt64([]int64{4, 5}, l.chargingMode) {
		startData, endData, periodList = l.moneyTransferHourTime(money, startDateParam)
	}
	// 按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
	if ContainsSliceInt64([]int64{6, 7}, l.chargingMode) {
		startData, endData, periodList = l.moneyTransferHalfOrHourAndMinuteTime(money, startDateParam)
	}
	// 按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
	if ContainsSliceInt64([]int64{8, 9}, l.chargingMode) {
		startData, endData, periodList = l.moneyTransferHalfOrHourAndMinuteTime(money, startDateParam)
	}
	// 自定义计费
	if ContainsSliceInt64([]int64{10, 11}, l.chargingMode) {
		startData, endData, periodList = l.moneyTransferCycleAndMinuteTime(money, startDateParam)
	}
	return startData, endData, periodList
}

/////////////////////////////////
//
//// 金额算出会有误差，具体以结束计算为准(半小时)
//func (l *Charging) moneyTransferHalfHourTimeModeTwo(money float64, startDate string) (string, string, []*PeriodList) {
//	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
//	var lastEndTime time.Time
//	periodList := make([]*PeriodList, 0)
//	periods := l.sortPeriods()
//
//	totalMoney := money
//	for {
//		if totalMoney <= 0 {
//			//fmt.Println("跳出循环")
//			break
//		}
//		for _, index := range periods {
//			period := l.periods[int64(index)]
//			var halfHourPrice float64
//			halfHourPrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
//
//			if currentTime.Format("15:04") == period.End {
//				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
//				continue
//			}
//			// 跨日
//			if period.Start > period.End && totalMoney > 0 {
//				tempCurrentTime := currentTime
//				if currentTime.Format("15:04") >= period.Start || currentTime.Format("15:04") < period.End {
//					if totalMoney < halfHourPrice {
//						minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//						totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//						currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//						lastEndTime = currentTime
//						totalMoney = 0
//						//fmt.Println(totalMoney, halfHourPrice, "隔日====下与半小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					} else {
//						currentTime = currentTime.Add(30 * time.Minute)
//						lastEndTime = currentTime
//						totalMoney -= halfHourPrice
//						//fmt.Println(totalMoney, halfHourPrice, "隔日====满半小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					}
//					duration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
//					periodList = append(periodList, &PeriodList{
//						StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
//						EndTime:    lastEndTime.Format("2006-01-02 15:04"),
//						Index:      index,
//						Duration:   duration,
//						Start:      period.Start,
//						End:        period.End,
//						Price:      period.Price,
//						TotalPrice: totalMoney,
//					})
//				}
//				continue
//			}
//			// 当日 6
//			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End && totalMoney > 0 {
//				tempCurrentTime := currentTime
//				if totalMoney < halfHourPrice {
//					//fmt.Println(totalMoney, halfHourPrice, ".......")
//					minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//					totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//					lastEndTime = currentTime
//					totalMoney = 0
//					//fmt.Println(totalMinute, halfHourPrice, totalMinute, "当日====小于半小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//				} else {
//					currentTime = currentTime.Add(time.Duration(30) * time.Minute)
//					lastEndTime = currentTime
//					totalMoney -= halfHourPrice
//					//fmt.Println(totalMoney, halfHourPrice, "当日====满足半小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//				}
//				duration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
//					Index:      index,
//					Duration:   duration,
//					Start:      period.Start,
//					End:        period.End,
//					Price:      period.Price,
//					TotalPrice: totalMoney,
//				})
//			}
//		}
//	}
//	//fmt.Println("结果：", startDate, lastEndTime.Format("2006-01-02 15:04"))
//	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
//}
//
//// 金额算出会有误差，具体以结束计算为准(一小时)
//func (l *Charging) moneyTransferHourTimeModeTwo(money float64, startDate string) (string, string, []*PeriodList) {
//	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
//	var lastEndTime time.Time
//	periodList := make([]*PeriodList, 0)
//	periods := l.sortPeriods()
//
//	totalMoney := money
//	for {
//		if totalMoney <= 0 {
//			//fmt.Println("跳出循环")
//			break
//		}
//		for _, index := range periods {
//			period := l.periods[int64(index)]
//			//var hourPrice float64
//			hourPrice := period.Price
//			if currentTime.Format("15:04") == period.End {
//				fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
//				continue
//			}
//			// 跨日
//			if period.Start > period.End && totalMoney > 0 {
//				tempCurrentTime := currentTime
//				if currentTime.Format("15:04") >= period.Start || currentTime.Format("15:04") < period.End {
//					if totalMoney < hourPrice {
//						minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//						totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//						currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//						lastEndTime = currentTime
//						totalMoney = 0
//						//fmt.Println(totalMoney, hourPrice, "隔日====下与一小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					} else {
//						currentTime = currentTime.Add(60 * time.Minute)
//						lastEndTime = currentTime
//						totalMoney -= hourPrice
//						//fmt.Println(totalMoney, hourPrice, "隔日====满一小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					}
//				}
//				duration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
//					Index:      index,
//					Duration:   duration,
//					Start:      period.Start,
//					End:        period.End,
//					Price:      period.Price,
//					TotalPrice: totalMoney,
//				})
//				continue
//			}
//			// 当日 6
//			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End && totalMoney > 0 {
//				tempCurrentTime := currentTime
//				if totalMoney < hourPrice {
//					minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//					totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//					lastEndTime = currentTime
//					totalMoney = 0
//					//fmt.Println(totalMinute, hourPrice, totalMinute, "当日====小于一小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//				} else {
//					currentTime = currentTime.Add(60 * time.Minute)
//					lastEndTime = currentTime
//					totalMoney -= hourPrice
//
//					//currentTime = currentTime.Add(time.Duration(30) * time.Minute)
//					//lastEndTime = currentTime
//					//totalMoney -= hourPrice
//					//fmt.Println(totalMoney, hourPrice, "当日====满足一小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//				}
//				duration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
//					Index:      index,
//					Duration:   duration,
//					Start:      period.Start,
//					End:        period.End,
//					Price:      period.Price,
//					TotalPrice: totalMoney,
//				})
//			}
//		}
//	}
//	//fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime))
//	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
//}
//
//// 金额算出会有误差，具体以结束计算为准(按半小时计费开台不足半小时按半时计费，超过半小时按分钟计费)
//func (l *Charging) moneyTransferHalfOrHourAndMinuteTimeModeTwo(money float64, startDate string) (string, string, []*PeriodList) {
//	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
//	var lastEndTime time.Time
//	periodList := make([]*PeriodList, 0)
//	periods := l.sortPeriods()
//
//	totalMoney := money
//	isExceed := 0
//	for {
//		if totalMoney <= 0 {
//			//fmt.Println("跳出循环")
//			break
//		}
//		var periodEndTime time.Time
//		var duration time.Duration
//		for _, index := range periods {
//			period := l.periods[int64(index)]
//			var halfOrHourPrice float64
//			var minute int
//			//fmt.Println("----====", currentTime.Format("15:04"), period.Start, period.End)
//			if currentTime.Format("15:04") == period.End {
//				fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
//				continue
//			}
//			// 不足半小时计费按小时计费，超过半小时按分钟计费
//			if l.chargingMode == 4 {
//				halfOrHourPrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
//				minute = 30
//			}
//			// 不足一小时计费按小时计费，超过一小时按分钟计费
//			if l.chargingMode == 5 {
//				halfOrHourPrice = period.Price
//				minute = 60
//			}
//			// 跨日
//			if period.Start > period.End && totalMoney > 0 {
//				if currentTime.Format("15:04") >= period.Start || currentTime.Format("15:04") < period.End {
//					tempCurrentTime := currentTime
//					// 总金额不足
//					if totalMoney < halfOrHourPrice && isExceed == 0 {
//						isExceed = 1
//					}
//					if totalMoney >= halfOrHourPrice && isExceed == 0 {
//						currentTime = currentTime.Add(time.Duration(minute) * time.Minute)
//						lastEndTime = currentTime
//						totalMoney -= halfOrHourPrice
//						isExceed = 1
//
//						Tduration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
//						periodList = append(periodList, &PeriodList{
//							StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
//							EndTime:    lastEndTime.Format("2006-01-02 15:04"),
//							Index:      index,
//							Duration:   Tduration,
//							Start:      period.Start,
//							End:        period.End,
//							Price:      period.Price,
//							TotalPrice: totalMoney,
//						})
//						//fmt.Println(totalMoney, halfOrHourPrice, "隔日====不足半小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//						continue
//					}
//					fmt.Println(isExceed, "isExceedisExceedisExceedisExceed")
//
//					if isExceed == 1 {
//						if currentTime.Format("15:04") >= period.Start {
//							//fmt.Println("大于", period.Start, period.Price)
//							periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), period.End), time.Local)
//							duration = periodEndTime.Sub(currentTime)
//						}
//						if currentTime.Format("15:04") < period.End {
//							//fmt.Println("小于", period.End, period.Price)
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
//						//fmt.Println(totalMoney, halfOrHourPrice, "隔日====不足半小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					}
//
//					Tduration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
//					periodList = append(periodList, &PeriodList{
//						StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
//						EndTime:    lastEndTime.Format("2006-01-02 15:04"),
//						Index:      index,
//						Duration:   Tduration,
//						Start:      period.Start,
//						End:        period.End,
//						Price:      period.Price,
//						TotalPrice: totalMoney,
//					})
//					continue
//				}
//			}
//			// 当日 6
//			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End && totalMoney > 0 {
//				tempCurrentTime := currentTime
//				if totalMoney >= halfOrHourPrice && isExceed == 0 {
//					currentTime = currentTime.Add(time.Duration(minute) * time.Minute)
//					lastEndTime = currentTime
//					totalMoney -= halfOrHourPrice
//					isExceed = 1
//					//fmt.Println(totalMoney, halfOrHourPrice, "当日====满足半小时=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//				} else {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//
//					var minutePrice float64
//					var totalMinute int64
//					minutePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//					totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//					if totalMinute < int64(duration.Minutes()) {
//						//fmt.Println("===|||今日", totalMinute, "...", int64(duration.Minutes()))
//						currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//						lastEndTime = currentTime
//						totalMoney = 0
//						//fmt.Println(totalMoney, halfOrHourPrice, "今日====不足一小时按分钟计费=====1", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					} else {
//						currentTime = periodEndTime
//						lastEndTime = currentTime
//						surplusMoney, _ := decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
//						totalMoney -= surplusMoney
//						//fmt.Println(totalMoney, halfOrHourPrice, "今日====不足一小时按分钟计费=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					}
//
//				}
//				Tduration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
//					Index:      index,
//					Duration:   Tduration,
//					Start:      period.Start,
//					End:        period.End,
//					Price:      period.Price,
//					TotalPrice: totalMoney,
//				})
//			}
//		}
//	}
//	//fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime))
//	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
//}
//
//// 金额算出会有误差，具体以结束计算为准(以10分钟作为一个收费周期，第1分钟后开始计费)
//func (l *Charging) moneyTransferCycleAndMinuteTimeModeTwo(money float64, startDate string) (string, string, []*PeriodList) {
//	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
//	var lastEndTime time.Time
//	totalMoney := money
//	periodList := make([]*PeriodList, 0)
//	periods := l.sortPeriods()
//
//	for {
//		if totalMoney <= 0 {
//			//fmt.Println("跳出循环")
//			break
//		}
//		for _, index := range periods {
//			var cyclePrice float64
//			period := l.periods[int64(index)]
//			if currentTime.Format("15:04") == period.End {
//				fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
//				continue
//			}
//
//			if l.cycleMinute == 30 {
//				cyclePrice, _ = decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
//			}
//			if l.cycleMinute == 60 {
//				cyclePrice = period.Price
//			}
//			if l.cycleMinute != 30 && l.cycleMinute != 60 {
//				minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//				cyclePrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
//			}
//			// 跨日
//			if period.Start > period.End && totalMoney > 0 {
//				tempCurrentTime := currentTime
//				if currentTime.Format("15:04") >= period.Start || currentTime.Format("15:04") < period.End {
//					if totalMoney < cyclePrice {
//						minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//						totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//						currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//						lastEndTime = currentTime
//						totalMoney = 0
//						//fmt.Println(totalMoney, cyclePrice, "隔日====小于周期内=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					} else {
//						currentTime = currentTime.Add(time.Duration(l.cycleMinute) * time.Minute)
//						lastEndTime = currentTime
//						totalMoney -= cyclePrice
//						//fmt.Println(totalMoney, cyclePrice, "隔日====满足周期内=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//					}
//				}
//				Tduration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
//					Index:      index,
//					Duration:   Tduration,
//					Start:      period.Start,
//					End:        period.End,
//					Price:      period.Price,
//					TotalPrice: totalMoney,
//				})
//				continue
//			}
//			// 当日 6
//			if period.Start < period.End && currentTime.Format("15:04") >= period.Start && currentTime.Format("15:04") < period.End && totalMoney > 0 {
//				tempCurrentTime := currentTime
//				if totalMoney < cyclePrice {
//					//fmt.Println(totalMoney, cyclePrice, ".......")
//					minutePrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(60)).Float64()
//					totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//					lastEndTime = currentTime
//					totalMoney = 0
//					//fmt.Println(totalMinute, cyclePrice, totalMinute, "当日====小于周期内=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//				} else {
//					currentTime = currentTime.Add(time.Duration(l.cycleMinute) * time.Minute)
//					lastEndTime = currentTime
//					totalMoney -= cyclePrice
//					//fmt.Println(totalMoney, cyclePrice, "当日====满足周期内=====", currentTime.Format(time.DateTime), lastEndTime.Format(time.DateTime))
//				}
//				Tduration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:  tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:    lastEndTime.Format("2006-01-02 15:04"),
//					Index:      index,
//					Duration:   Tduration,
//					Start:      period.Start,
//					End:        period.End,
//					Price:      period.Price,
//					TotalPrice: totalMoney,
//				})
//			}
//		}
//	}
//	//fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime))
//	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
//}
