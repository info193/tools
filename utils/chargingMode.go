package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"sort"
	"strconv"
	"time"
)

// 收费时段设置
type ChargePeriodAssembly struct {
	Week       *ChargePeriodSetWeek    `json:"week"`        // 星期
	Holiday    *ChargePeriodSetHoliday `json:"holiday"`     // 节假日
	Hour       *ChargePeriodSetHour    `json:"hour"`        // 小时收费时段
	HourStairs []HourStairs            `json:"hour_stairs"` // 阶梯计费
	MinConsume *ChargePeriodMinConsume `json:"min_consume"`
}

// 星期
type ChargePeriodSetWeek struct {
	Week     []int64  `json:"week"`      // 星期
	Hour     []CPHour `json:"hour"`      // 收费时段
	HourPeak []CPHour `json:"hour_peak"` // 时段封顶
	//MinConsume *ChargePeriodMinConsume `json:"min_consume"`
}

// 小时收费时段
type ChargePeriodSetHour struct {
	Hour     []CPHour `json:"hour"`      // 收费时段
	HourPeak []CPHour `json:"hour_peak"` // 时段封顶
	//MinConsume *ChargePeriodMinConsume `json:"min_consume"`
}

// 节假日
type ChargePeriodSetHoliday struct {
	Date     [][]string `json:"date"`      // 星期
	Hour     [][]CPHour `json:"hour"`      // 收费时段
	HourPeak [][]CPHour `json:"hour_peak"` // 时段封顶
	//MinConsume []ChargePeriodMinConsume `json:"min_consume"`
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
	Start       int64  `json:"start"`
	End         int64  `json:"end"`
	IdlePrice   string `json:"idle_price"`
	MemberPrice string `json:"member_price"`
}

type ChargePeriod struct {
	StartPeriod int64   `json:"start_period"`
	EndPeriod   int64   `json:"end_period"`
	Start       string  `json:"start"`
	End         string  `json:"end"`
	Price       float64 `json:"price"`
	//HourPeak    *[]CPHour               `json:"hour_peak"` // 时段封顶
	//MinConsume  *ChargePeriodMinConsume `json:"min_consume"`
}

type Charging struct {
	//periods map[int64]ChargePeriod
	periods      ChargePeriodAssembly
	member       int64 // 是否是会员 0否  1是
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

	//HourPeak *[]CPHour `json:"hour_peak"` // 时段封顶
	////HourStairs *[]HourStairs           `json:"hour_stairs"` // 时段阶梯
	//MinConsume *ChargePeriodMinConsume `json:"min_consume"`
}

type HourPeriodList struct {
	StartDate     string  `json:"start_date"`
	EndDate       string  `json:"end_date"`
	Duration      int64   `json:"duration"`
	HourPrice     float64 `json:"hour_price"`
	Price         float64 `json:"price"`
	HourPeak      int64   `json:"hour_peak"`
	HourPeakPrice float64 `json:"hour_peak_price"`
}

func NewChargeMode(periods ChargePeriodAssembly, member int64, chargingMode int64, cycleMinute int64) *Charging {
	return &Charging{periods: periods, member: member, chargingMode: chargingMode, cycleMinute: cycleMinute}
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
		if tempMinute < l.cycleMinute {
			endDateTime = endDateTime.Add(time.Duration(l.cycleMinute-tempMinute) * time.Minute)
		}
	}
	return endDateTime
}

// 需注意：不管是否是跨日还是当日，开始时间不能为24，已0开始。
// 当日可以有24点,跨日不能出现24点，已0点开始。
// 计费
func (l *Charging) Outlay(startDate, endDate string) (float64, map[string]HourPeriodList) {
	var price float64
	//periodList := make([]*PeriodList, 0)
	hourPeriodList := make(map[string]HourPeriodList, 0)
	//periods := make(map[int64]float64) // 时段费用详情
	startDateTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	endDateTime, _ := time.ParseInLocation("2006-01-02 15:04", endDate, time.Local)
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	tempCurrentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	endDateTime = l.fillTime(startDateTime, endDateTime) // 填充时间

	//l.getCharge(startDateTime, endDateTime)
	startDateYMD := startDateTime.Format("2006-01-02")
	endDateYMD := endDateTime.Format("2006-01-02")
	k := 0
	//hourStairsArr := make(map[string][]HourStairs, 0)
	// 时段封顶
	if len(l.periods.HourStairs) == 0 {
		price, hourPeriodList = l.stairs(startDateTime, endDateTime, l.periods.HourStairs)
		return price, hourPeriodList
	} else {
		times := make([]periodTimes, 0)
		for {
			if tempCurrentTime.Format("2006-01-02 15:04") > endDateTime.Format("2006-01-02 15:04") {
				break
			}
			//isHourStairs := 0
			isOk := 0
			hours := make(map[int]ChargePeriod, 0)
			var hourPrice float64
			// 节假日
			if l.periods.Holiday != nil && isOk == 0 {
				startDate := currentTime.Format("01-02")
				for _, value := range l.periods.Holiday.Date {
					if ContainsSliceString(value, startDate) {
						for index, val := range value {
							if val == startDate {
								for ind, hour := range l.periods.Holiday.Hour[index] {
									hourPrice, _ = strconv.ParseFloat(hour.IdlePrice, 10)
									if l.member == 1 {
										hourPrice, _ = strconv.ParseFloat(hour.MemberPrice, 10)
									}
									var start string
									var end string
									if hour.Start < 10 {
										start = fmt.Sprintf("0%v:00", hour.Start)
									} else {
										start = fmt.Sprintf("%v:00", hour.Start)
									}
									if hour.End < 10 {
										end = fmt.Sprintf("0%v:00", hour.End)
									} else {
										end = fmt.Sprintf("%v:00", hour.End)
									}
									//if hour.End == 24 {
									//	end = "23:59"
									//}
									//var hourPeak *[]CPHour
									//if len(l.periods.Holiday.HourPeak[index]) >= 1 {
									//	hourPeak = &l.periods.Holiday.HourPeak[index]
									//}
									//fmt.Println("=====", hourPeak)
									//var hourStairs *[]HourStairs
									//if len(l.periods.Holiday.HourStairs[index]) >= 1 {
									//	isHourStairs = 1
									//	//hourStairs = &l.periods.Holiday.HourStairs[index]
									//	hourStairsArr[currentTime.Format("2006-01-02")] = l.periods.Holiday.HourStairs[index]
									//}
									//var minConsume *ChargePeriodMinConsume
									//if index >= 0 && index < len(l.periods.Holiday.MinConsume) {
									//	minConsume = &l.periods.Holiday.MinConsume[index]
									//}
									///////////////////////////////////////////////
									//fmt.Println(l.periods.Holiday.HourPeak[index], "------")
									//fmt.Println(l.periods.Holiday.HourStairs[index], "======")
									//fmt.Println(l.periods.Holiday.MinConsume[index], "???????")
									hours[ind] = ChargePeriod{StartPeriod: hour.Start, EndPeriod: hour.End, Start: start, End: end, Price: hourPrice}
								}
								//fmt.Println("节假日在范围内", hours)
								continue
							}
						}
						isOk = 1
						continue
					}
				}
			}

			// 星期
			if l.periods.Week != nil && isOk == 0 {
				startDay := int64(currentTime.Weekday())
				if startDay == 0 {
					startDay = 7
				}
				if ContainsSliceInt64(l.periods.Week.Week, startDay) {
					isOk = 1
					for ind, hour := range l.periods.Week.Hour {
						hourPrice, _ = strconv.ParseFloat(hour.IdlePrice, 10)
						if l.member == 1 {
							hourPrice, _ = strconv.ParseFloat(hour.MemberPrice, 10)
						}
						var start string
						var end string
						if hour.Start < 10 {
							start = fmt.Sprintf("0%v:00", hour.Start)
						} else {
							start = fmt.Sprintf("%v:00", hour.Start)
						}
						if hour.End < 10 {
							end = fmt.Sprintf("0%v:00", hour.End)
						} else {
							end = fmt.Sprintf("%v:00", hour.End)
						}
						//if hour.End == 24 {
						//	end = "23:59"
						//}
						//var hourPeak *[]CPHour
						//if len(l.periods.Week.HourPeak) >= 1 {
						//	hourPeak = &l.periods.Week.HourPeak
						//}
						//var hourStairs *[]HourStairs
						//if len(l.periods.Week.HourStairs) >= 1 {
						//	isHourStairs = 1
						//	hourStairsArr[currentTime.Format("2006-01-02")] = l.periods.Week.HourStairs
						//	hourStairs = &l.periods.Week.HourStairs
						//}
						//var minConsume *ChargePeriodMinConsume
						//if l.periods.Week.MinConsume != nil {
						//	minConsume = l.periods.Week.MinConsume
						//}
						hours[ind] = ChargePeriod{StartPeriod: hour.Start, EndPeriod: hour.End, Start: start, End: end, Price: hourPrice}
					}

					//fmt.Println(tempStartTime.Format(time.DateTime), "开始时间存在星期")
					//fmt.Println("星期在范围内", tempCurrentTime.Format(time.DateTime))
				}
			}

			// 小时
			if tempCurrentTime.Unix() <= endDateTime.Unix() && l.periods.Hour != nil && isOk == 0 {
				isOk = 1
				for ind, hour := range l.periods.Hour.Hour {
					hourPrice, _ = strconv.ParseFloat(hour.IdlePrice, 10)
					if l.member == 1 {
						hourPrice, _ = strconv.ParseFloat(hour.MemberPrice, 10)
					}
					var start string
					var end string
					if hour.Start < 10 {
						start = fmt.Sprintf("0%v:00", hour.Start)
					} else {
						start = fmt.Sprintf("%v:00", hour.Start)
					}
					if hour.End < 10 {
						end = fmt.Sprintf("0%v:00", hour.End)
					} else {
						end = fmt.Sprintf("%v:00", hour.End)
					}
					if hour.End == 24 {
						end = "23:59"
					}
					//var hourPeak *[]CPHour
					//if len(l.periods.Hour.HourPeak) >= 1 {
					//	hourPeak = &l.periods.Hour.HourPeak
					//}
					//var hourStairs *[]HourStairs
					//if len(l.periods.Hour.HourStairs) >= 1 {
					//	isHourStairs = 1
					//	hourStairsArr[currentTime.Format("2006-01-02")] = l.periods.Hour.HourStairs
					//	hourStairs = &l.periods.Hour.HourStairs
					//}
					//var minConsume *ChargePeriodMinConsume
					//if l.periods.Hour.MinConsume != nil {
					//	minConsume = l.periods.Hour.MinConsume
					//}
					hours[ind] = ChargePeriod{StartPeriod: hour.Start, EndPeriod: hour.End, Start: start, End: end, Price: hourPrice}
				}
				//fmt.Println("小时在范围内")
			}
			//fmt.Println(hours, "....|||||||||...")
			//fmt.Println(tempCurrentTime.Format(time.DateTime), count)

			//fmt.Println(hours, "====")
			var periodStartTime time.Time
			var periodEndTime time.Time
			//if isHourStairs == 1 {
			//	//tempStartTime := tempCurrentTime
			//	tempAddCurrentTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v 00:00", tempCurrentTime.Add(time.Second*86400).Format("2006-01-02")), time.Local)
			//	if tempAddCurrentTime.Unix() > endDateTime.Unix() {
			//		//fmt.Println("1----", tempStartTime.Format(time.DateTime), endDateTime.Format(time.DateTime))
			//		tempCurrentTime = endDateTime.Add(120 * time.Second)
			//	} else {
			//		//fmt.Println("2----", tempStartTime.Format(time.DateTime), tempAddCurrentTime.Format(time.DateTime))
			//		tempCurrentTime = tempAddCurrentTime
			//	}
			//	currentTime = tempCurrentTime
			//	//fmt.Println(hourStairsArr)
			//	//fmt.Println(tempCurrentTime.Format(time.DateTime), "------2")
			//	//} else {
			//	//	//tempCurrentTime = tempCurrentTime.Add(86400 * time.Second)
			//	//	//fmt.Println(tempCurrentTime.Format(time.DateTime), "------1")
			//	//}
			//}
			tempCurrentTime = tempCurrentTime.Add(86400 * time.Second)
			count := len(hours)
			for i := 0; i < count; i++ {
				if period, ok := hours[i]; ok {
					//if period.HourStairs != nil {
					//	fmt.Println(tempCurrentTime.Format(time.DateTime), period.HourStairs, "=======")
					//	continue
					//}
					//if period.HourPeak == nil && period.HourStairs != nil {
					//	//fmt.Println(period.HourPeak, "??????")
					//	// 判断如果时段封顶为空，且时段阶梯存在，则按时段阶梯计费
					//	if len(*period.HourStairs) > 0 {
					//		fmt.Println(*period.HourStairs, "..............", period.Start, period.End)
					//		break
					//	}
					//}

					//fmt.Println(period, "...", period.HourPeak)
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
							if period.EndPeriod == 24 {
								periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Add(86400*time.Second).Format("2006-01-02"), "00:00"), time.Local)
							}
						}
						k++
					}

					// 判断 时段结束时间大于时段开始时间
					if period.Start < period.End {
						periodStartTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.Start), time.Local)
						periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), period.End), time.Local)
						if period.EndPeriod == 24 {
							periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Add(86400*time.Second).Format("2006-01-02"), "00:00"), time.Local)
						}
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
							if periodV, ok := hours[int(i+1)]; ok {
								nextPeriodStartTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), periodV.Start), time.Local)
								if nextPeriodStartTime.Unix() >= periodEndTime.Unix() && nextPeriodStartTime.Unix() < endDateTime.Unix() {
									end = 0
								}
							} else if periodV, ok := hours[0]; ok {
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
			_, hourPeriodList = l.computeMinute(times)
			price, hourPeriodList = l.computePeak(hourPeriodList)
			// 时段封顶
			// 时段阶梯
			//fmt.Println(hourPeriodList, "....===")
		}
		// 2按半小时计费
		if l.chargingMode == 2 {
			_, hourPeriodList = l.computeHalfHour(times)
			price, hourPeriodList = l.computePeak(hourPeriodList)
		}
		// 3按半小时计费(跨时段)
		if l.chargingMode == 3 {
			_, hourPeriodList = l.computeHalfHourModeTwo(times)
			price, hourPeriodList = l.computePeak(hourPeriodList)
		}
		// 4按小时计费
		if l.chargingMode == 4 {
			_, hourPeriodList = l.computeHour(times)
			price, hourPeriodList = l.computePeak(hourPeriodList)
		}
		// 5按小时计费(跨时段)
		if l.chargingMode == 5 {
			_, hourPeriodList = l.computeHourModeTwo(times)
			price, hourPeriodList = l.computePeak(hourPeriodList)
		}
		// 6按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
		if l.chargingMode == 6 {
			_, hourPeriodList = l.computeHalfHourOrMinute(times)
			price, hourPeriodList = l.computePeak(hourPeriodList)
		}
		// 7按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费(跨时段)
		if l.chargingMode == 7 {
			_, hourPeriodList = l.computeHalfHourOrMinuteModeTwo(times)
			price, hourPeriodList = l.computePeak(hourPeriodList)
		}
		// 8按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
		if l.chargingMode == 8 {
			_, hourPeriodList = l.computeHourOrMinute(times)
			price, hourPeriodList = l.computePeak(hourPeriodList)
		}
		// 9按小时计费开台不足1小时按小时计费，超过1小时按分钟计费(跨时段)
		if l.chargingMode == 9 {
			_, hourPeriodList = l.computeHourOrMinuteModeTwo(times)
			price, hourPeriodList = l.computePeak(hourPeriodList)
		}
		// 10 以10分钟作为一个收费周期，第1分钟后开始计费
		if l.chargingMode == 10 {
			_, hourPeriodList = l.computeCycleAndMinute(times)
			price, hourPeriodList = l.computePeak(hourPeriodList)
		}
		// 11 以10分钟作为一个收费周期，第1分钟后开始计费(跨时段)
		if l.chargingMode == 11 {
			_, hourPeriodList = l.computeCycleAndMinuteModeTwo(times)
			price, hourPeriodList = l.computePeak(hourPeriodList)
		}
	}

	//时段封顶,算出时段所有价格,减掉该时段内金额，然后再加上封顶时段价格
	//阶段计费，按时长算
	return price, hourPeriodList
}

func (l *Charging) computePeak(hourPeriodList map[string]HourPeriodList) (float64, map[string]HourPeriodList) {
	dateSlice := make([]string, 0)
	dateSlicePrice := make(map[string]float64, 0)
	for date, val := range hourPeriodList {
		dateSlice = append(dateSlice, date)
		dateSlicePrice[date] = val.Price
	}
	sort.Strings(dateSlice)
	dateAccrual := make(map[string]map[int]float64, 0)
	for _, date := range dateSlice {
		isOk := 0
		currentTime, _ := time.ParseInLocation("2006-01-02 15", date, time.Local)
		period := int64(currentTime.Hour())
		dateKey := currentTime.Format("2006-01-02")
		dateTimeKey := currentTime.Format("2006-01-02 15")
		if _, ok := dateAccrual[dateKey]; !ok {
			dateAccrual[dateKey] = make(map[int]float64)
		}
		// 节假日
		if l.periods.Holiday != nil && isOk == 0 {
			startDate := currentTime.Format("01-02")
			for _, value := range l.periods.Holiday.Date {
				if ContainsSliceString(value, startDate) {
					for index, val := range value {
						if val == startDate {
							if len(l.periods.Holiday.HourPeak[index]) >= 1 {
								for _, v := range l.periods.Holiday.HourPeak[index] {
									var peakPrice float64
									peakPrice, _ = strconv.ParseFloat(v.IdlePrice, 10)
									if l.member == 1 {
										peakPrice, _ = strconv.ParseFloat(v.MemberPrice, 10)
									}
									// 跨日
									if v.Start > v.End {
										if v.Start <= period && v.End < period || v.Start > period && v.End > period {
											//dateAccrual[dateKey][index] += dateSlicePrice[date]
											//if _, ok := dateDeduct[dateKey][index]; !ok {
											//	dateDeduct[dateKey][index] = peakPrice
											//}
											if _, ok := dateAccrual[dateKey][index]; !ok {
												if dateSlicePrice[dateTimeKey] >= peakPrice {
													dateSlicePrice[dateTimeKey] = peakPrice
													dateAccrual[dateKey][index] += peakPrice
												} else {
													dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
												}
											} else if daVal, ok := dateAccrual[dateKey][index]; ok {
												if daVal >= peakPrice {
													dateSlicePrice[dateTimeKey] = 0
												} else {
													surplusPeak := (peakPrice - daVal)
													dateAccrual[dateKey][index] += surplusPeak
													if dateSlicePrice[date] > surplusPeak {
														dateSlicePrice[date] = surplusPeak
													}
												}
											}
											//if dateDeduct[dateKey][index] < dateAccrual[dateKey][index] {
											//	dateSlicePrice[dateTimeKey] = 0
											//	//fmt.Println(dateTimeKey, "至为0")
											//}
											isOk = 1
										}
									}
									// 当日
									if v.Start < v.End && v.Start <= period && v.End > period {
										//dateAccrual[dateKey][index] += dateSlicePrice[date]
										//if _, ok := dateDeduct[dateKey][index]; !ok {
										//	dateDeduct[dateKey][index] = peakPrice
										//}
										if _, ok := dateAccrual[dateKey][index]; !ok {
											if dateSlicePrice[dateTimeKey] >= peakPrice {
												dateSlicePrice[dateTimeKey] = peakPrice
												dateAccrual[dateKey][index] += peakPrice
											} else {
												dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
											}
										} else if daVal, ok := dateAccrual[dateKey][index]; ok {
											if daVal >= peakPrice {
												dateSlicePrice[dateTimeKey] = 0
											} else {
												surplusPeak := (peakPrice - daVal)
												dateAccrual[dateKey][index] += surplusPeak
												if dateSlicePrice[date] > surplusPeak {
													dateSlicePrice[date] = surplusPeak
												}
											}
										}
										//if dateDeduct[dateKey][index] < dateAccrual[dateKey][index] {
										//	//fmt.Println(dateTimeKey, "至为0")
										//	dateSlicePrice[dateTimeKey] = 0
										//}
										isOk = 1
									}
								}
							}
						}
					}
					continue
				}
			}
		}
		//星期
		if l.periods.Week != nil && isOk == 0 {
			startDay := int64(currentTime.Weekday())
			if startDay == 0 {
				startDay = 7
			}
			if ContainsSliceInt64(l.periods.Week.Week, startDay) {
				if len(l.periods.Week.HourPeak) >= 1 {
					for index, v := range l.periods.Week.HourPeak {
						var peakPrice float64
						peakPrice, _ = strconv.ParseFloat(v.IdlePrice, 10)
						if l.member == 1 {
							peakPrice, _ = strconv.ParseFloat(v.MemberPrice, 10)
						}
						// 跨日
						if v.Start > v.End {
							if v.Start <= period && v.End < period || v.Start > period && v.End > period {
								//dateAccrual[dateKey][index] += dateSlicePrice[date]
								//if _, ok := dateDeduct[dateKey][index]; !ok {
								//	dateDeduct[dateKey][index] = peakPrice
								//}

								if _, ok := dateAccrual[dateKey][index]; !ok {
									if dateSlicePrice[dateTimeKey] >= peakPrice {
										dateSlicePrice[dateTimeKey] = peakPrice
										dateAccrual[dateKey][index] += peakPrice
									} else {
										dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
									}
								} else if daVal, ok := dateAccrual[dateKey][index]; ok {
									if daVal >= peakPrice {
										dateSlicePrice[dateTimeKey] = 0
									} else {
										surplusPeak := (peakPrice - daVal)
										dateAccrual[dateKey][index] += surplusPeak
										if dateSlicePrice[date] > surplusPeak {
											dateSlicePrice[date] = surplusPeak
										}
									}
								}

								isOk = 1
							}
						}
						// 当日
						if v.Start < v.End && v.Start <= period && v.End > period {
							//dateAccrual[dateKey][index] += dateSlicePrice[date]
							//if _, ok := dateDeduct[dateKey][index]; !ok {
							//	dateDeduct[dateKey][index] = peakPrice
							//}

							if _, ok := dateAccrual[dateKey][index]; !ok {
								if dateSlicePrice[dateTimeKey] >= peakPrice {
									dateSlicePrice[dateTimeKey] = peakPrice
									dateAccrual[dateKey][index] += peakPrice
								} else {
									dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
								}
							} else if daVal, ok := dateAccrual[dateKey][index]; ok {
								if daVal >= peakPrice {
									dateSlicePrice[dateTimeKey] = 0
								} else {
									surplusPeak := (peakPrice - daVal)
									dateAccrual[dateKey][index] += surplusPeak
									if dateSlicePrice[date] > surplusPeak {
										dateSlicePrice[date] = surplusPeak
									}
								}
							}
							isOk = 1
						}
					}
				}
			}
			//fmt.Println(tempStartTime.Format(time.DateTime), "开始时间存在星期")
			//fmt.Println("星期在范围内", tempCurrentTime.Format(time.DateTime))
		}
		// 小时
		if l.periods.Hour != nil && isOk == 0 {
			startDay := int64(currentTime.Weekday())
			if startDay == 0 {
				startDay = 7
			}
			if len(l.periods.Hour.HourPeak) >= 1 {
				for index, v := range l.periods.Hour.HourPeak {
					var peakPrice float64
					peakPrice, _ = strconv.ParseFloat(v.IdlePrice, 10)
					if l.member == 1 {
						peakPrice, _ = strconv.ParseFloat(v.MemberPrice, 10)
					}
					// 跨日
					if v.Start > v.End {
						if v.Start <= period && v.End < period || v.Start > period && v.End > period {
							//dateAccrual[dateKey][index] += dateSlicePrice[date]
							//if _, ok := dateDeduct[dateKey][index]; !ok {
							//	dateDeduct[dateKey][index] = peakPrice
							//}
							if _, ok := dateAccrual[dateKey][index]; !ok {
								if dateSlicePrice[dateTimeKey] >= peakPrice {
									dateSlicePrice[dateTimeKey] = peakPrice
									dateAccrual[dateKey][index] += peakPrice
								} else {
									dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
								}
							} else if daVal, ok := dateAccrual[dateKey][index]; ok {
								if daVal >= peakPrice {
									dateSlicePrice[dateTimeKey] = 0
								} else {
									surplusPeak := (peakPrice - daVal)
									dateAccrual[dateKey][index] += surplusPeak
									if dateSlicePrice[date] > surplusPeak {
										dateSlicePrice[date] = surplusPeak
									}
								}
							}
							//if dateDeduct[dateKey][index] < dateAccrual[dateKey][index] {
							//	dateSlicePrice[dateTimeKey] = 0
							//	//fmt.Println(dateTimeKey, "至为0")
							//}
						}
					}
					// 当日
					if v.Start < v.End && v.Start <= period && v.End > period {
						//dateAccrual[dateKey][index] += dateSlicePrice[date]
						//if _, ok := dateDeduct[dateKey][index]; !ok {
						//	dateDeduct[dateKey][index] = peakPrice
						//}
						if _, ok := dateAccrual[dateKey][index]; !ok {
							if dateSlicePrice[dateTimeKey] >= peakPrice {
								dateSlicePrice[dateTimeKey] = peakPrice
								dateAccrual[dateKey][index] += peakPrice
							} else {
								dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
							}
						} else if daVal, ok := dateAccrual[dateKey][index]; ok {
							if daVal >= peakPrice {
								dateSlicePrice[dateTimeKey] = 0
							} else {
								surplusPeak := (peakPrice - daVal)
								dateAccrual[dateKey][index] += surplusPeak
								if dateSlicePrice[date] > surplusPeak {
									dateSlicePrice[date] = surplusPeak
								}
							}
						}
						//if dateDeduct[dateKey][index] < dateAccrual[dateKey][index] {
						//	dateSlicePrice[dateTimeKey] = 0
						//	//fmt.Println(dateTimeKey, "至为0")
						//}
					}
				}
			}
		}
	}

	for date, val := range hourPeriodList {
		tempVal := val
		if valPrice, ok := dateSlicePrice[date]; ok {
			if valPrice != tempVal.Price {
				if vm, oks := hourPeriodList[date]; oks {
					vm.HourPeak = 1
					vm.HourPeakPrice = valPrice
					hourPeriodList[date] = vm
				}
			}
		}
	}
	var price float64
	for _, value := range dateSlicePrice {
		tempValuePrice := value
		price += tempValuePrice
	}
	//系统提示词：按什么方法什么方式做什么事情
	return price, hourPeriodList
}

func (l *Charging) stairs(startDateTime, endDateTime time.Time, stairs []HourStairs) (float64, map[string]HourPeriodList) {
	var price float64
	duration := decimal.NewFromFloat(endDateTime.Sub(startDateTime).Minutes()).IntPart()
	hourPeriodLists := make(map[string]HourPeriodList, 0)
	tempCreateTime := startDateTime
	for _, value := range stairs {
		if duration < value.Start {
			break
		}
		var totalPrice float64
		tempPrice, _ := strconv.ParseFloat(value.IdlePrice, 10)
		if l.member == 1 {
			tempPrice, _ = strconv.ParseFloat(value.MemberPrice, 10)
		}

		if duration <= value.End && value.Start == 0 {
			if l.chargingMode == 1 {
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
				totalPrice, _ = decimal.NewFromInt(duration).Mul(decimal.NewFromFloat(minutePrice)).Float64()
			}
			// 半小时
			if ContainsSliceInt64([]int64{2, 3}, l.chargingMode) {
				halfPrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
				maxHalfHour := decimal.NewFromInt(duration).Div(decimal.NewFromInt(30)).Ceil().IntPart()
				totalPrice, _ = decimal.NewFromInt(maxHalfHour).Mul(decimal.NewFromFloat(halfPrice)).Float64()
			}
			// 小时
			if ContainsSliceInt64([]int64{4, 5}, l.chargingMode) {
				maxHour := decimal.NewFromInt(duration).Div(decimal.NewFromInt(60)).Ceil().IntPart()
				totalPrice, _ = decimal.NewFromInt(maxHour).Mul(decimal.NewFromFloat(tempPrice)).Float64()
			}
			// 按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
			if ContainsSliceInt64([]int64{6, 7}, l.chargingMode) {
				halfPrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
				if duration < 30 {
					totalPrice = halfPrice
				} else {
					minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
					totalPrice, _ = decimal.NewFromInt(duration - 30).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					totalPrice += halfPrice
				}
			}
			// 按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
			if ContainsSliceInt64([]int64{8, 9}, l.chargingMode) {
				if duration < 60 {
					totalPrice = tempPrice
				} else {
					minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
					totalPrice, _ = decimal.NewFromInt(duration - 60).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					totalPrice += tempPrice
				}
			}
			// 以10分钟作为一个收费周期，第1分钟后开始计费
			if ContainsSliceInt64([]int64{10, 11}, l.chargingMode) {
				var cyclePrice float64
				// 半小时
				if l.cycleMinute == 30 {
					halfPrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
					cyclePrice = halfPrice
				}
				// 小时
				if l.cycleMinute == 60 {
					cyclePrice = tempPrice
				}
				if l.cycleMinute != 30 && l.cycleMinute != 60 {
					minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
					cyclePrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
				}
				maxCycleHour := decimal.NewFromInt(duration).Div(decimal.NewFromInt(l.cycleMinute)).Ceil().IntPart()
				totalPrice, _ = decimal.NewFromInt(maxCycleHour).Mul(decimal.NewFromFloat(cyclePrice)).Float64()
			}
			tempCurrentTime := tempCreateTime.Add(time.Duration(duration) * time.Minute)
			hourPeriodLists[tempCurrentTime.Format("2006-01-02 15")] = HourPeriodList{
				StartDate: tempCreateTime.Format("2006-01-02 15:04"),
				EndDate:   tempCurrentTime.Format("2006-01-02 15:04"),
				Duration:  duration,
				HourPrice: tempPrice,
				Price:     totalPrice,
			}
			tempCreateTime = tempCurrentTime
			// 立邦、三棵树、多乐士 水性漆
			price += totalPrice
		}

		// 累加前期计费
		if duration >= value.End {
			tempValue := value.End - value.Start
			if l.chargingMode == 1 {
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
				totalPrice, _ = decimal.NewFromInt(tempValue).Mul(decimal.NewFromFloat(minutePrice)).Float64()
			}
			// 半小时
			if ContainsSliceInt64([]int64{2, 3}, l.chargingMode) {
				halfPrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
				maxHalfHour := decimal.NewFromInt(tempValue).Div(decimal.NewFromInt(30)).Ceil().IntPart()
				totalPrice, _ = decimal.NewFromInt(maxHalfHour).Mul(decimal.NewFromFloat(halfPrice)).Float64()
			}
			// 小时
			if ContainsSliceInt64([]int64{4, 5}, l.chargingMode) {
				maxHour := decimal.NewFromInt(tempValue).Div(decimal.NewFromInt(60)).Ceil().IntPart()
				totalPrice, _ = decimal.NewFromInt(maxHour).Mul(decimal.NewFromFloat(tempPrice)).Float64()
			}
			// 按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
			if ContainsSliceInt64([]int64{6, 7}, l.chargingMode) {
				halfPrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
				totalPrice, _ = decimal.NewFromInt(tempValue - 30).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				totalPrice += halfPrice
			}

			// 按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
			if ContainsSliceInt64([]int64{8, 9}, l.chargingMode) {
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
				totalPrice, _ = decimal.NewFromInt(tempValue - 60).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				totalPrice += tempPrice
			}
			// 以10分钟作为一个收费周期，第1分钟后开始计费
			if ContainsSliceInt64([]int64{10, 11}, l.chargingMode) {
				var cyclePrice float64
				// 半小时
				if l.cycleMinute == 30 {
					halfPrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
					cyclePrice = halfPrice
				}
				// 小时
				if l.cycleMinute == 60 {
					cyclePrice = tempPrice
				}
				if l.cycleMinute != 30 && l.cycleMinute != 60 {
					minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
					cyclePrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
				}
				maxCycleHour := decimal.NewFromInt(tempValue).Div(decimal.NewFromInt(l.cycleMinute)).Ceil().IntPart()
				totalPrice, _ = decimal.NewFromInt(maxCycleHour).Mul(decimal.NewFromFloat(cyclePrice)).Float64()
			}
			tempCurrentTime := tempCreateTime.Add(time.Duration(tempValue) * time.Minute)
			hourPeriodLists[tempCreateTime.Format("2006-01-02 15")] = HourPeriodList{
				StartDate: tempCreateTime.Format("2006-01-02 15:04"),
				EndDate:   tempCurrentTime.Format("2006-01-02 15:04"),
				Duration:  tempValue,
				HourPrice: tempPrice,
				Price:     totalPrice,
			}
			tempCreateTime = tempCurrentTime
			price += totalPrice
			continue
		}

		// 在最后计费段内
		if duration > value.Start && duration <= value.End && value.Start != 0 {
			tempValue := duration - value.Start
			if l.chargingMode == 1 {
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
				totalPrice, _ = decimal.NewFromInt(tempValue).Mul(decimal.NewFromFloat(minutePrice)).Float64()
			}
			// 半小时
			if ContainsSliceInt64([]int64{2, 3}, l.chargingMode) {
				halfPrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
				maxHalfHour := decimal.NewFromInt(tempValue).Div(decimal.NewFromInt(30)).Ceil().IntPart()
				totalPrice, _ = decimal.NewFromInt(maxHalfHour).Mul(decimal.NewFromFloat(halfPrice)).Float64()
			}
			// 小时
			if ContainsSliceInt64([]int64{4, 5}, l.chargingMode) {
				maxHour := decimal.NewFromInt(tempValue).Div(decimal.NewFromInt(60)).Ceil().IntPart()
				totalPrice, _ = decimal.NewFromInt(maxHour).Mul(decimal.NewFromFloat(tempPrice)).Float64()
			}
			// 按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
			if ContainsSliceInt64([]int64{6, 7}, l.chargingMode) {
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
				totalPrice, _ = decimal.NewFromInt(tempValue).Mul(decimal.NewFromFloat(minutePrice)).Float64()
			}
			// 按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
			if ContainsSliceInt64([]int64{8, 9}, l.chargingMode) {
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
				totalPrice, _ = decimal.NewFromInt(tempValue).Mul(decimal.NewFromFloat(minutePrice)).Float64()
			}
			// 以10分钟作为一个收费周期，第1分钟后开始计费
			if ContainsSliceInt64([]int64{10, 11}, l.chargingMode) {
				var cyclePrice float64
				// 半小时
				if l.cycleMinute == 30 {
					halfPrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
					cyclePrice = halfPrice
				}
				// 小时
				if l.cycleMinute == 60 {
					cyclePrice = tempPrice
				}
				if l.cycleMinute != 30 && l.cycleMinute != 60 {
					minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
					cyclePrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
				}
				maxCycleHour := decimal.NewFromInt(tempValue).Div(decimal.NewFromInt(l.cycleMinute)).Ceil().IntPart()
				totalPrice, _ = decimal.NewFromInt(maxCycleHour).Mul(decimal.NewFromFloat(cyclePrice)).Float64()
			}
			tempCurrentTime := tempCreateTime.Add(time.Duration(tempValue) * time.Minute)
			hourPeriodLists[tempCreateTime.Format("2006-01-02 15")] = HourPeriodList{
				StartDate: tempCreateTime.Format("2006-01-02 15:04"),
				EndDate:   tempCurrentTime.Format("2006-01-02 15:04"),
				Duration:  tempValue,
				HourPrice: tempPrice,
				Price:     totalPrice,
			}
			tempCreateTime = tempCurrentTime
			price += totalPrice
		}
	}
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	//fmt.Println(price, duration, hourPeriodLists)
	return price, hourPeriodLists
}

func (l *Charging) computeMinute(times []periodTimes) (float64, map[string]HourPeriodList) {
	//periods := make(map[int64]float64) // 时段费用详情
	var price float64
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	hourPeriodLists := make(map[string][]HourPeriodList, 0)
	//dec := make(map[string]float64, 0)
	for _, timeV := range times {
		//duration := timeV.EndTime.Sub(timeV.StartTime)
		minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//totalPrice, _ := decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		tempStartTime := timeV.StartTime
		for {
			if timeV.EndTime.Unix() <= tempStartTime.Unix() {
				break
			}
			originStartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", tempStartTime.Format("2006-01-02 15:00:00"), time.Local)
			tempMinutes := decimal.NewFromFloat(tempStartTime.Sub(originStartTime).Minutes()).IntPart()
			if tempMinutes != 0 {
				hourTotalPrice, _ := decimal.NewFromInt(60 - tempMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				tempStartTimeD := tempStartTime.Add(time.Duration(60-tempMinutes) * time.Minute)
				//hourPeriodList[tempStartTime.Format("2006-01-02 15")] = hourTotalPrice
				//fmt.Println(hourTotalPrice, "====", tempStartTime.Format("2006-01-02 15"), "......", tempStartTime.Format(time.DateTime))
				hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
					HourPeriodList{Price: hourTotalPrice, Duration: 60 - tempMinutes, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
				tempStartTime = tempStartTimeD
			} else {
				if tempStartTime.Unix() <= timeV.EndTime.Unix() {
					//var hourTotalPrice float64
					tempEndMinutes := decimal.NewFromFloat(timeV.EndTime.Sub(tempStartTime).Minutes()).IntPart()
					if tempEndMinutes < 60 {
						hourTotalPrice, _ := decimal.NewFromInt(tempEndMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
						//fmt.Println(hourTotalPrice, ".....")
						//hourPeriodList[tempStartTime.Format("2006-01-02 15")] = hourTotalPrice
						hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
							HourPeriodList{Price: hourTotalPrice, Duration: tempEndMinutes, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTime.Add(time.Duration(tempEndMinutes) * time.Minute).Format("2006-01-02 15:04"), HourPrice: timeV.Price})
						tempStartTime = tempStartTime.Add(time.Duration(tempEndMinutes) * time.Minute)
					} else {
						//hourTotalPrice, _ = decimal.NewFromFloat(60).Mul(decimal.NewFromFloat(minutePrice)).Float64()
						//hourPeriodList[tempStartTime.Format("2006-01-02 15")] = timeV.Price
						hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
							HourPeriodList{Price: timeV.Price, Duration: 60, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTime.Add(60 * time.Minute).Format("2006-01-02 15:04"), HourPrice: timeV.Price})
						tempStartTime = tempStartTime.Add(60 * time.Minute)
					}
					//fmt.Println(hourTotalPrice, "====", tempStartTime.Format("2006-01-02 15"), "|||||||||||", tempStartTime.Format(time.DateTime))
				}
			}
		}
		//fmt.Println(totalPrice, "=======", duration.Minutes(), timeV.EndTime.Format(time.DateTime), "...", timeV.StartTime.Format(time.DateTime))
		//price += totalPrice
		//periods[int64(timeV.Index)] += totalPrice
		////fmt.Println(timeV.StartTime.Format(time.DateTime), timeV.EndTime.Format(time.DateTime), ".....")
		////fmt.Println(price, totalPrice, "======")
		//periodList = append(periodList, &PeriodList{
		//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
		//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
		//	Index:      timeV.Index,
		//	Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
		//	Start:      timeV.Start,
		//	End:        timeV.End,
		//	Price:      timeV.Price,
		//	TotalPrice: totalPrice,
		//})
	}
	//fmt.Println(".......", hourPeriodLists)
	//allHourPeriodList := make(map[string][]HourPeriodList, 0)
	//for _, value := range hourPeriodLists {
	//	for _, val := range value {
	//		tempDateTime, _ := time.ParseInLocation("2006-01-02 15:04", val.StartDate, time.Local)
	//		allHourPeriodList[tempDateTime.Format("2006-01-02 15")] = append(allHourPeriodList[tempDateTime.Format("2006-01-02 15")], HourPeriodList{StartDate: val.StartDate, EndDate: val.EndDate, Duration: val.Duration, HourPrice: val.HourPrice, Price: val.Price})
	//	}
	//}
	//fmt.Println(allHourPeriodList)
	alHourPeriodList := make(map[string]HourPeriodList, 0)

	dateSlice := make([]string, 0)
	for date, _ := range hourPeriodLists {
		dateSlice = append(dateSlice, date)
	}
	sort.Strings(dateSlice)
	for _, date := range dateSlice {
		value := hourPeriodLists[date]
		var tempPrice float64
		var tempDuration int64
		var tempHourPrice float64
		var tempStartDate string
		var tempEndDate string
		for _, val := range value {
			if tempStartDate == "" {
				tempStartDate = val.StartDate
			}
			tempEndDate = val.EndDate
			tempPrice += val.Price
			tempDuration += val.Duration
			tempHourPrice = val.HourPrice
			price += val.Price
		}

		alHourPeriodList[date] = HourPeriodList{StartDate: tempStartDate, EndDate: tempEndDate, HourPrice: tempHourPrice, Price: tempPrice, Duration: tempDuration}
	}

	//fmt.Println("价格", price, "===", alHourPeriodList)
	//fmt.Println("=======", times)
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, alHourPeriodList
}

// 半小时计费 (半小时为节点，以半小时起始点计算价格，如果时间段内计算时长小于半小时，则累计至下个时段，
// 并在下个时段内扣除上个时段少于的时间[因为以收半个时段价格]，然后在计算下周期时段价格，依次类推)
func (l *Charging) computeHalfHour(times []periodTimes) (float64, map[string]HourPeriodList) {
	//periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	hourPeriodLists := make(map[string][]HourPeriodList, 0)

	for _, timeV := range times {
		if lastDiffMinute > 0 {
			timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
			lastDiffMinute = 0
		}

		halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
		// 每小时消费金额
		tempStartTime := timeV.StartTime
		for {
			if timeV.EndTime.Unix() <= tempStartTime.Unix() {
				break
			}
			tempMinute := timeV.EndTime.Sub(tempStartTime).Minutes()
			if tempMinute < 30 {
				lastDiffMinute = 30 - int64(tempMinute)
			}
			hourTotalPrice, _ := decimal.NewFromFloat(1).Mul(decimal.NewFromFloat(halfHourPrice)).Float64()
			//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
			tempStartTimeD := tempStartTime.Add(30 * time.Minute)
			hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
				HourPeriodList{Price: hourTotalPrice, Duration: 30, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
			//fmt.Println(tempStartTimeD.Format(time.DateTime), ".............")
			tempStartTime = tempStartTimeD
		}

		//var totalPrice float64
		//duration := timeV.EndTime.Sub(timeV.StartTime)
		//if duration.Minutes() > 0 {
		//	haltHour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(30)).Ceil().IntPart()
		//	maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
		//	lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
		//
		//	//halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
		//	totalPrice, _ = decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour)).Float64()
		//
		//	price += totalPrice
		//	periods[int64(timeV.Index)] += totalPrice
		//} else {
		//	periods[int64(timeV.Index)] += 0
		//}
		//periodList = append(periodList, &PeriodList{
		//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
		//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
		//	Index:      timeV.Index,
		//	Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
		//	Start:      timeV.Start,
		//	End:        timeV.End,
		//	Price:      timeV.Price,
		//	TotalPrice: totalPrice,
		//})
	}
	//fmt.Println("-++++++", hourPeriodLists)
	allHourPeriodList := make(map[string]HourPeriodList, 0)
	//fmt.Println(allHourPeriodList, "========")
	//var cdP float64
	for date, value := range hourPeriodLists {
		var tempPrice float64
		var tempDuration int64
		var tempHourPrice float64
		var tempStartDate string
		var tempEndDate string
		for _, val := range value {
			//fmt.Println(val.StartDate, val.EndDate, "时长=", val.Duration, "----", val.HourPrice, "总价", val.Price)
			if tempStartDate == "" {
				tempStartDate = val.StartDate
			}
			tempEndDate = val.EndDate
			tempPrice += val.Price
			tempDuration += val.Duration
			tempHourPrice = val.HourPrice
			price += val.Price
		}
		allHourPeriodList[date] = HourPeriodList{StartDate: tempStartDate, EndDate: tempEndDate, HourPrice: tempHourPrice, Price: tempPrice, Duration: tempDuration}
	}
	//fmt.Println("价格", price, "===", allHourPeriodList)

	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, allHourPeriodList
}

// 一小时计费 (一小时为节点，以一小时起始点计算价格，如果时间段内计算时长小于一小时，则累计至下个时段，
// 并在下个时段内扣除上个时段少于的时间[因为以收一个时段价格]，然后在计算下周期时段价格，依次类推)
func (l *Charging) computeHour(times []periodTimes) (float64, map[string]HourPeriodList) {
	//periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	hourPeriodLists := make(map[string][]HourPeriodList, 0)

	for _, timeV := range times {
		if lastDiffMinute > 0 {
			timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
			lastDiffMinute = 0
		}

		// 每小时消费金额
		tempStartTime := timeV.StartTime
		for {
			if timeV.EndTime.Unix() <= tempStartTime.Unix() {
				break
			}
			tempMinute := timeV.EndTime.Sub(tempStartTime).Minutes()
			if tempMinute < 60 {
				lastDiffMinute = 60 - int64(tempMinute)
			}
			//hourTotalPrice, _ := decimal.NewFromFloat(1).Mul(decimal.NewFromFloat(timeV.Price)).Float64()
			//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += timeV.Price
			tempStartTimeD := tempStartTime.Add(60 * time.Minute)
			hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
				HourPeriodList{Price: timeV.Price, Duration: 60, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
			tempStartTime = tempStartTimeD
		}

		//var totalPrice float64
		//duration := timeV.EndTime.Sub(timeV.StartTime)
		//if duration.Minutes() > 0 {
		//	hour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(60)).Ceil().IntPart()
		//	maxMinute := decimal.NewFromInt(hour).Mul(decimal.NewFromInt(60)).IntPart()
		//	lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
		//	totalPrice, _ = decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour)).Float64()
		//	price += totalPrice
		//	periods[int64(timeV.Index)] += totalPrice
		//} else {
		//	periods[int64(timeV.Index)] += 0
		//}
		//periodList = append(periodList, &PeriodList{
		//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
		//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
		//	Index:      timeV.Index,
		//	Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
		//	Start:      timeV.Start,
		//	End:        timeV.End,
		//	Price:      timeV.Price,
		//	TotalPrice: totalPrice,
		//})
	}
	//fmt.Println("-++++++", hourPeriodLists)
	allHourPeriodList := make(map[string]HourPeriodList, 0)
	//fmt.Println(allHourPeriodList, "========")
	//var cdP float64
	for date, value := range hourPeriodLists {
		var tempPrice float64
		var tempDuration int64
		var tempHourPrice float64
		var tempStartDate string
		var tempEndDate string
		for _, val := range value {
			//fmt.Println(val.StartDate, val.EndDate, val.HourPrice, val.Price)
			if tempStartDate == "" {
				tempStartDate = val.StartDate
			}
			tempEndDate = val.EndDate
			tempPrice += val.Price
			tempDuration += val.Duration
			tempHourPrice = val.HourPrice
			price += val.Price
		}
		allHourPeriodList[date] = HourPeriodList{StartDate: tempStartDate, EndDate: tempEndDate, HourPrice: tempHourPrice, Price: tempPrice, Duration: tempDuration}
	}
	//fmt.Println("价格", price, "===", allHourPeriodList)

	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, allHourPeriodList
}

// 按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
func (l *Charging) computeHalfHourOrMinute(times []periodTimes) (float64, map[string]HourPeriodList) {
	//periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	//var firstHalfHour int64
	//firstHalfHour = 30
	firstHalfHourFor := 30
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	hourPeriodLists := make(map[string][]HourPeriodList, 0)
	lenIndex := len(times)
	for index, timeV := range times {

		if lastDiffMinute != 0 {
			timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		}
		halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
		minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()

		tempStartTime := timeV.StartTime
		for {
			if timeV.EndTime.Unix() <= tempStartTime.Unix() {
				break
			}
			var tempStartTimeD time.Time
			if firstHalfHourFor > 0 {
				tempStartTimeD = tempStartTime.Add(30 * time.Minute)
				//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += halfHourPrice
				hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: halfHourPrice, Duration: 30, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
				firstHalfHourFor = 0
				if timeV.EndTime.Sub(tempStartTime).Minutes() < 30 {
					lastDiffMinute = 30 - int64(timeV.EndTime.Sub(tempStartTime).Minutes())
				}
				//fmt.Println(halfHourPrice, "-------", tempStartTimeD.Format(time.DateTime))
				tempStartTime = tempStartTimeD
			} else {
				tempStartTimeD = tempStartTime.Add(1800 * time.Second)
				// 判断 相差值，必须为数组最后之前的数据，如果是数组最后一个则按半小时收费
				if tempStartTimeD.Unix() > timeV.EndTime.Unix() && lenIndex-1 > index {
					minute := timeV.EndTime.Sub(tempStartTime).Minutes()
					hourTotalPrice, _ := decimal.NewFromFloat(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					tempStartTimeD = tempStartTime.Add(time.Duration(minute) * time.Minute)
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: int64(minute), StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//lastDiffMinute = 30 - int64(minute)
					//fmt.Println(hourTotalPrice, "=|||||||===", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime), minute)
					tempStartTime = tempStartTimeD
				} else {
					var hourTotalPrice float64
					tempEndMinutes := decimal.NewFromFloat(timeV.EndTime.Sub(tempStartTime).Minutes()).IntPart()
					var tempDuration int64
					if tempEndMinutes < 30 {
						tempDuration = tempEndMinutes
						hourTotalPrice, _ = decimal.NewFromInt(tempEndMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
						tempStartTimeD = tempStartTime.Add(time.Duration(tempEndMinutes) * time.Minute)
						//fmt.Println("------1", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime), tempEndMinutes)
					} else {
						tempDuration = 30
						hourTotalPrice, _ = decimal.NewFromFloat(30).Mul(decimal.NewFromFloat(minutePrice)).Float64()
						tempStartTimeD = tempStartTime.Add(time.Duration(tempDuration) * time.Minute)
						//fmt.Println("------0000", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime))
					}
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: tempDuration, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//fmt.Println(hourTotalPrice, "====", tempStartTime.Format("2006-01-02 15"), "|||||||||||", tempStartTime.Format(time.DateTime))
					tempStartTime = tempStartTimeD
				}
			}
			//}
		}

		//var totalPrice float64
		//duration := timeV.EndTime.Sub(timeV.StartTime)
		//if duration.Minutes() > 0 {
		//	if firstHalfHour != 0 {
		//		haltHour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(30)).Ceil().IntPart()
		//		maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
		//		lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
		//		tempMinute := maxMinute - 30 - lastDiffMinute
		//		totalPrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(tempMinute)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(halfHourPrice)).Float64()
		//		price += totalPrice
		//		periods[int64(timeV.Index)] += totalPrice
		//		firstHalfHour = 0 // 首个半小时计费，计费完后置为0
		//	} else {
		//		//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//		price += totalPrice
		//		periods[int64(timeV.Index)] += totalPrice
		//		lastDiffMinute = 0 // 分钟计费，差值分钟置为0
		//	}
		//} else {
		//	periods[int64(timeV.Index)] += 0
		//}
		//periodList = append(periodList, &PeriodList{
		//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
		//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
		//	Index:      timeV.Index,
		//	Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
		//	Start:      timeV.Start,
		//	End:        timeV.End,
		//	Price:      timeV.Price,
		//	TotalPrice: totalPrice,
		//})

	}

	//fmt.Println(hourPeriodLists, "--------")
	allHourPeriodList := make(map[string]HourPeriodList, 0)
	//var cdP float64
	//fmt.Println(allHourPeriodList, "========")
	for date, value := range hourPeriodLists {
		var tempPrice float64
		var tempDuration int64
		var tempHourPrice float64
		var tempStartDate string
		var tempEndDate string
		for _, val := range value {
			//fmt.Println(val.StartDate, val.EndDate, "时长=", val.Duration, "----", val.HourPrice, "总价", val.Price)
			if tempStartDate == "" {
				tempStartDate = val.StartDate
			}
			tempEndDate = val.EndDate
			tempPrice += val.Price
			tempDuration += val.Duration
			tempHourPrice = val.HourPrice
			price += val.Price
		}
		allHourPeriodList[date] = HourPeriodList{StartDate: tempStartDate, EndDate: tempEndDate, HourPrice: tempHourPrice, Price: tempPrice, Duration: tempDuration}
	}
	//fmt.Println(price, "金额", allHourPeriodList, "..----,,,")

	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, allHourPeriodList
}

// 按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
func (l *Charging) computeHourOrMinute(times []periodTimes) (float64, map[string]HourPeriodList) {
	//periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	firstHourFor := 60
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	hourPeriodLists := make(map[string][]HourPeriodList, 0)

	lenIndex := len(times)
	for index, timeV := range times {

		if lastDiffMinute != 0 {
			timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		}
		//halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
		minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()

		//
		tempStartTime := timeV.StartTime
		for {
			if timeV.EndTime.Unix() <= tempStartTime.Unix() {
				break
			}
			var tempStartTimeD time.Time
			if firstHourFor > 0 {
				tempStartTimeD = tempStartTime.Add(60 * time.Minute)
				//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += timeV.Price
				hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: timeV.Price, Duration: 60, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
				firstHourFor = 0
				if timeV.EndTime.Sub(tempStartTime).Minutes() < 60 {
					lastDiffMinute = 60 - int64(timeV.EndTime.Sub(tempStartTime).Minutes())
				}
				//fmt.Println(halfHourPrice, "-------", tempStartTimeD.Format(time.DateTime))
				tempStartTime = tempStartTimeD
			} else {
				tempStartTimeD = tempStartTime.Add(3600 * time.Second)
				// 判断 相差值，必须为数组最后之前的数据，如果是数组最后一个则按半小时收费
				if tempStartTimeD.Unix() > timeV.EndTime.Unix() && lenIndex-1 > index {
					minute := timeV.EndTime.Sub(tempStartTime).Minutes()
					hourTotalPrice, _ := decimal.NewFromFloat(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					tempStartTimeD = tempStartTime.Add(time.Duration(minute) * time.Minute)
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: int64(minute), StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//lastDiffMinute = 30 - int64(minute)
					//fmt.Println(hourTotalPrice, "=|||||||===", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime), minute)
					tempStartTime = tempStartTimeD
				} else {
					var hourTotalPrice float64
					tempEndMinutes := decimal.NewFromFloat(timeV.EndTime.Sub(tempStartTime).Minutes()).IntPart()
					var tempDuration int64
					if tempEndMinutes < 60 {
						tempDuration = tempEndMinutes
						hourTotalPrice, _ = decimal.NewFromInt(tempEndMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
						tempStartTimeD = tempStartTime.Add(time.Duration(tempEndMinutes) * time.Minute)
						//fmt.Println("------1", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime), tempEndMinutes)
					} else {
						tempDuration = 60
						hourTotalPrice, _ = decimal.NewFromFloat(60).Mul(decimal.NewFromFloat(minutePrice)).Float64()
						tempStartTimeD = tempStartTime.Add(time.Duration(tempDuration) * time.Minute)
						//fmt.Println("------0000", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime))
					}
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: tempDuration, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//fmt.Println(hourTotalPrice, "====", tempStartTime.Format("2006-01-02 15"), "|||||||||||", tempStartTime.Format(time.DateTime))
					tempStartTime = tempStartTimeD
				}
			}
			//}
		}

		//var totalPrice float64
		//duration := timeV.EndTime.Sub(timeV.StartTime)
		//if duration.Minutes() > 0 {
		//	if firstHalfHour != 0 {
		//		haltHour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(30)).Ceil().IntPart()
		//		maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
		//		lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
		//		tempMinute := maxMinute - 30 - lastDiffMinute
		//		totalPrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(tempMinute)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(halfHourPrice)).Float64()
		//		price += totalPrice
		//		periods[int64(timeV.Index)] += totalPrice
		//		firstHalfHour = 0 // 首个半小时计费，计费完后置为0
		//	} else {
		//		//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//		price += totalPrice
		//		periods[int64(timeV.Index)] += totalPrice
		//		lastDiffMinute = 0 // 分钟计费，差值分钟置为0
		//	}
		//} else {
		//	periods[int64(timeV.Index)] += 0
		//}
		//periodList = append(periodList, &PeriodList{
		//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
		//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
		//	Index:      timeV.Index,
		//	Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
		//	Start:      timeV.Start,
		//	End:        timeV.End,
		//	Price:      timeV.Price,
		//	TotalPrice: totalPrice,
		//})

	}

	//fmt.Println(hourPeriodLists, "--------")
	allHourPeriodList := make(map[string]HourPeriodList, 0)
	//var cdP float64
	//fmt.Println(allHourPeriodList, "========")
	for date, value := range hourPeriodLists {
		var tempPrice float64
		var tempDuration int64
		var tempHourPrice float64
		var tempStartDate string
		var tempEndDate string
		for _, val := range value {
			//fmt.Println(val.StartDate, val.EndDate, "时长=", val.Duration, "----", val.HourPrice, "总价", val.Price)
			if tempStartDate == "" {
				tempStartDate = val.StartDate
			}
			tempEndDate = val.EndDate
			tempPrice += val.Price
			tempDuration += val.Duration
			tempHourPrice = val.HourPrice
			price += val.Price
		}
		allHourPeriodList[date] = HourPeriodList{StartDate: tempStartDate, EndDate: tempEndDate, HourPrice: tempHourPrice, Price: tempPrice, Duration: tempDuration}
	}
	//fmt.Println(price, "金额", allHourPeriodList, "..----,,,")

	//for _, timeV := range times {
	//	if lastDiffMinute != 0 {
	//		timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
	//	}
	//
	//	minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
	//	tempStartTime := timeV.StartTime
	//	for {
	//		if timeV.EndTime.Unix() <= tempStartTime.Unix() {
	//			break
	//		}
	//		if firstHourFor != 0 {
	//			tempStartTimeD := tempStartTime.Add(60 * time.Minute)
	//			hourPeriodList[tempStartTime.Format("2006-01-02 15")] += timeV.Price
	//			//fmt.Println(timeV.Price, "-------", tempStartTimeD.Format(time.DateTime))
	//			tempStartTime = tempStartTimeD
	//			firstHourFor = 0
	//		} else {
	//			originStartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", tempStartTime.Format("2006-01-02 15:00:00"), time.Local)
	//			tempMinutes := decimal.NewFromFloat(tempStartTime.Sub(originStartTime).Minutes()).IntPart()
	//			if tempMinutes != 0 {
	//				hourTotalPrice, _ := decimal.NewFromInt(60 - tempMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
	//				tempStartTimeD := tempStartTime.Add(time.Duration(60-tempMinutes) * time.Minute)
	//				hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
	//				//fmt.Println(hourTotalPrice, "====", tempStartTime.Format("2006-01-02 15"), "......", tempStartTime.Format(time.DateTime))
	//				tempStartTime = tempStartTimeD
	//			} else {
	//				if tempStartTime.Unix() <= timeV.EndTime.Unix() {
	//					var hourTotalPrice float64
	//					tempEndMinutes := decimal.NewFromFloat(timeV.EndTime.Sub(tempStartTime).Minutes()).IntPart()
	//					if tempEndMinutes < 60 {
	//						hourTotalPrice, _ = decimal.NewFromInt(tempEndMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
	//					} else {
	//						hourTotalPrice, _ = decimal.NewFromFloat(60).Mul(decimal.NewFromFloat(minutePrice)).Float64()
	//					}
	//					hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
	//					//fmt.Println(hourTotalPrice, "====", tempStartTime.Format("2006-01-02 15"), "|||||||||||", tempStartTime.Format(time.DateTime))
	//					tempStartTime = tempStartTime.Add(3600 * time.Second)
	//				}
	//			}
	//		}
	//
	//	}
	//
	//	//var totalPrice float64
	//	//duration := timeV.EndTime.Sub(timeV.StartTime)
	//	//if duration.Minutes() > 0 {
	//	//	if firstHour != 0 {
	//	//		hour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(60)).Ceil().IntPart()
	//	//		maxMinute := decimal.NewFromInt(hour).Mul(decimal.NewFromInt(60)).IntPart()
	//	//		lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
	//	//
	//	//		tempMinute := maxMinute - 60 - lastDiffMinute
	//	//		totalPrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(tempMinute)).Float64()
	//	//		totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(timeV.Price)).Float64()
	//	//		price += totalPrice
	//	//		periods[int64(timeV.Index)] += totalPrice
	//	//		firstHour = 0 // 首个一小时计费，计费完后置为0
	//	//	} else {
	//	//		//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
	//	//		totalPrice, _ = decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
	//	//		price += totalPrice
	//	//		periods[int64(timeV.Index)] += totalPrice
	//	//		lastDiffMinute = 0 // 分钟计费，差值分钟置为0
	//	//	}
	//	//} else {
	//	//	periods[int64(timeV.Index)] += 0
	//	//}
	//	//periodList = append(periodList, &PeriodList{
	//	//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
	//	//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
	//	//	Index:      timeV.Index,
	//	//	Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
	//	//	Start:      timeV.Start,
	//	//	End:        timeV.End,
	//	//	Price:      timeV.Price,
	//	//	TotalPrice: totalPrice,
	//	//})
	//}

	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, allHourPeriodList
}

// 以10分钟作为一个收费周期，第1分钟后开始计费，如果时间段内计算时长小于收费周期，累计至下个时段
// 并在下个时段内扣除上个时段少于的时间[因为以收一个时段价格]，然后在计算下周期时段价格，依次类推)
func (l *Charging) computeCycleAndMinute(times []periodTimes) (float64, map[string]HourPeriodList) {
	//periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var lastDiffMinute int64
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	hourPeriodLists := make(map[string][]HourPeriodList, 0)
	for _, timeV := range times {
		if lastDiffMinute != 0 {
			timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		}

		var cyclePrice float64
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
			minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
			cyclePrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
		}

		// 周期消费金额
		tempStartTime := timeV.StartTime
		for {
			if timeV.EndTime.Unix() <= tempStartTime.Unix() {
				break
			}

			tempMinute := int64(timeV.EndTime.Sub(tempStartTime).Minutes())
			if tempMinute < l.cycleMinute {
				lastDiffMinute = l.cycleMinute - int64(tempMinute)
			}
			//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += cyclePrice
			tempStartTimeD := tempStartTime.Add(time.Duration(l.cycleMinute) * time.Minute)
			hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
				HourPeriodList{Price: cyclePrice, Duration: l.cycleMinute, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
			tempStartTime = tempStartTimeD
		}

		//var totalPrice float64
		//duration := timeV.EndTime.Sub(timeV.StartTime)
		//if duration.Minutes() > 0 {
		//	// 半小时
		//	if l.cycleMinute == 30 {
		//		haltHour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(30)).Ceil().IntPart()
		//		maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
		//		lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
		//		//halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(cyclePrice).Mul(decimal.NewFromInt(haltHour)).Float64()
		//	}
		//
		//	// 一小时
		//	if l.cycleMinute == 60 {
		//		hour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(60)).Ceil().IntPart()
		//		maxMinute := decimal.NewFromInt(hour).Mul(decimal.NewFromInt(60)).IntPart()
		//		lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
		//		totalPrice, _ = decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour)).Float64()
		//	}
		//	// 周期计费
		//	if l.cycleMinute != 30 && l.cycleMinute != 60 {
		//		//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//		//cyclePrice, _ := decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
		//		cycle := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(l.cycleMinute)).Ceil().IntPart()
		//		maxMinute := decimal.NewFromInt(cycle).Mul(decimal.NewFromInt(l.cycleMinute)).IntPart()
		//		lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
		//		totalPrice, _ = decimal.NewFromFloat(cyclePrice).Mul(decimal.NewFromInt(cycle)).Float64()
		//	}
		//	price += totalPrice
		//	periods[int64(timeV.Index)] += totalPrice
		//} else {
		//	periods[int64(timeV.Index)] += 0
		//}
		//periodList = append(periodList, &PeriodList{
		//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
		//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
		//	Index:      timeV.Index,
		//	Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
		//	Start:      timeV.Start,
		//	End:        timeV.End,
		//	Price:      timeV.Price,
		//	TotalPrice: totalPrice,
		//})
	}

	//fmt.Println("-++++++", hourPeriodLists)
	allHourPeriodList := make(map[string]HourPeriodList, 0)
	//fmt.Println(allHourPeriodList, "========")
	//var cdP float64
	for date, value := range hourPeriodLists {
		var tempPrice float64
		var tempDuration int64
		var tempHourPrice float64
		var tempStartDate string
		var tempEndDate string
		for _, val := range value {
			//fmt.Println(val.StartDate, val.EndDate, val.HourPrice, val.Price)
			if tempStartDate == "" {
				tempStartDate = val.StartDate
			}
			tempEndDate = val.EndDate
			tempPrice += val.Price
			tempDuration += val.Duration
			tempHourPrice = val.HourPrice
			price += val.Price
		}
		allHourPeriodList[date] = HourPeriodList{StartDate: tempStartDate, EndDate: tempEndDate, HourPrice: tempHourPrice, Price: tempPrice, Duration: tempDuration}
	}
	//fmt.Println("价格", price, "===", allHourPeriodList)

	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, allHourPeriodList
}

// 半小时计费 (半小时为节点，以半小时起始点计算价格，如果时间段内计算时长小于半小时，则累计至下个时段，
// 在下个时段内扣除上个时段少于的时间[如果跨时段按跨时段分钟计费,否则按半小时算]，然后在计算下周期时段价格，依次类推)
func (l *Charging) computeHalfHourModeTwo(times []periodTimes) (float64, map[string]HourPeriodList) {
	//periods := make(map[int64]float64) // 时段费用详情
	var price float64
	//var lastDiffMinute int64
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	hourPeriodLists := make(map[string][]HourPeriodList, 0)

	var hourLastMinute int64
	lenIndex := len(times)
	for index, timeV := range times {
		//var totalPrice float64

		halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
		minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()

		var tempStartTime time.Time
		var tempStartTimeD time.Time
		//if hourLastMinute != 0 {
		//	tempStartTime = timeV.StartTime.Add(time.Duration(hourLastMinute) * time.Minute)
		//} else {
		tempStartTime = timeV.StartTime
		//}
		for {
			if timeV.EndTime.Unix() <= tempStartTime.Unix() {
				break
			}
			if hourLastMinute > 0 {
				// 跨时段按分钟计费
				tempFirstMinutes := int64(timeV.EndTime.Sub(tempStartTime).Minutes())
				if tempFirstMinutes < 30 {
					tempStartTimeD = tempStartTime.Add(time.Duration(tempFirstMinutes) * time.Minute)
					hourTotalPrice, _ := decimal.NewFromInt(tempFirstMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: tempFirstMinutes, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTime.Add(time.Duration(tempFirstMinutes) * time.Minute).Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					if tempFirstMinutes > hourLastMinute {
						hourLastMinute = 0
					} else {
						hourLastMinute -= tempFirstMinutes
					}

				} else if tempFirstMinutes >= 30 && hourLastMinute != 0 {
					tempStartTimeD = tempStartTime.Add(time.Duration(hourLastMinute) * time.Minute)
					hourTotalPrice, _ := decimal.NewFromInt(hourLastMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: hourLastMinute, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTime.Add(time.Duration(hourLastMinute) * time.Minute).Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					hourLastMinute = 0
				} else {
					hourTotalPrice, _ := decimal.NewFromInt(hourLastMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
						HourPeriodList{Price: hourTotalPrice, Duration: hourLastMinute, StartDate: tempStartTime.Add(time.Duration(-hourLastMinute) * time.Minute).Format("2006-01-02 15:04"), EndDate: tempStartTime.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					tempStartTimeD = tempStartTime.Add(time.Duration(hourLastMinute) * time.Minute)
					hourLastMinute = 0
				}

				//hourTotalPrice, _ := decimal.NewFromInt(hourLastMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
				//hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
				//	HourPeriodList{Price: hourTotalPrice, Duration: hourLastMinute, StartDate: tempStartTime.Add(time.Duration(-hourLastMinute) * time.Minute).Format("2006-01-02 15:04"), EndDate: tempStartTime.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
				//hourLastMinute = 0
				//fmt.Println("11111111", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime), "===", hourTotalPrice)
			} else {
				tempStartTimeD = tempStartTime.Add(1800 * time.Second)
				// 判断 相差值，必须为数组最后之前的数据，如果是数组最后一个则按半小时收费
				if tempStartTimeD.Unix() > timeV.EndTime.Unix() && lenIndex-1 > index {
					minute := timeV.EndTime.Sub(tempStartTime).Minutes()
					hourLastMinute = 30 - int64(minute)
					hourTotalPrice, _ := decimal.NewFromFloat(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					tempStartTimeD = timeV.EndTime
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: int64(minute), StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: timeV.EndTime.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//fmt.Println("22222222", tempStartTime.Format(time.DateTime), "===", hourTotalPrice, minute, timeV.EndTime.Format(time.DateTime), hourLastMinute)
				} else {
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += halfHourPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: halfHourPrice, Duration: 30, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//fmt.Println("333333", tempStartTimeD.Format(time.DateTime), "===", halfHourPrice)
				}
			}
			tempStartTime = tempStartTimeD
		}

		//// 差值按分钟计算计费
		//if lastDiffMinute > 0 {
		//	lastStartTime := timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		//	if lastStartTime.Unix() >= timeV.EndTime.Unix() {
		//		lastStartTime = timeV.EndTime
		//	}
		//	duration := lastStartTime.Sub(timeV.StartTime)
		//	totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
		//	//fmt.Println(".......-----", lastDiffMinute, duration, totalMinute)
		//	if totalMinute > 0 {
		//		//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//		minuteTotalPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
		//		timeV.StartTime = lastStartTime
		//		//fmt.Println("跨时段按分钟", lastDiffMinute, duration, totalMinute, totalPrice, timeV.StartTime.Format(time.DateTime))
		//	}
		//}
		//
		//duration := timeV.EndTime.Sub(timeV.StartTime)
		//totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
		//if totalMinute > 0 {
		//	haltHour := decimal.NewFromInt(totalMinute).Div(decimal.NewFromInt(30)).Ceil().IntPart()
		//	maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
		//	var tempEndTime time.Time
		//	if timeV.Start > timeV.End && timeV.StartTime.Format("15:04") >= timeV.Start {
		//		tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", timeV.EndTime.Add(time.Duration(86400)*time.Second).Format("2006-01-02"), timeV.End), time.Local)
		//		//} else if timeV.Start > timeV.End && timeV.StartTime.Format("15:04") >= "00:00" {
		//		//	tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", timeV.EndTime.Format("2006-01-02"), timeV.End), time.Local)
		//	} else {
		//		tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", timeV.EndTime.Format("2006-01-02"), timeV.End), time.Local)
		//	}
		//
		//	//fmt.Println("..;;;;;;;;;;", timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format(time.DateTime), tempEndTime.Format(time.DateTime))
		//
		//	if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Unix() >= tempEndTime.Unix() {
		//		lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(totalMinute)).IntPart()
		//		//fmt.Println(lastDiffMinute, "....;;;;;;;;;")
		//	} else {
		//		lastDiffMinute = 0
		//	}
		//	//fmt.Println(lastDiffMinute, "lastDiffMinute====", timeV.StartTime.Format(time.DateTime), timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format(time.DateTime), timeV.End, ".....", tempEndTime.Format(time.DateTime))
		//	//fmt.Println(lastDiffMinute, "======222")
		//	//halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
		//	// 开始时间大于结束时间 跨天
		//	if timeV.Start > timeV.End {
		//		//fmt.Println("1111")
		//		// 跨时段
		//		if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format("15:04") >= timeV.End && lastDiffMinute != 0 {
		//			//minute := decimal.NewFromInt(30).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
		//			if totalMinute < 30 && timeV.EndTime.Format("15:04") < timeV.End {
		//				//fmt.Println("22222")
		//				totalHourPrice, _ := decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour)).Float64()
		//				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//				//fmt.Println("不足半小时按半小时算", totalPrice)
		//			} else {
		//				//fmt.Println("33333")
		//				minute := decimal.NewFromInt(30).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
		//				//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//				minuteTotalPrice, _ := decimal.NewFromInt(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//
		//				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
		//				totalHourPrice, _ := decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour - 1)).Float64()
		//				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//			}
		//			//	fmt.Println("跨时段", minuteTotalPrice, totalPrice, lastDiffMinute)
		//		} else {
		//			//fmt.Println("444444")
		//			totalHourPrice, _ := decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour)).Float64()
		//			//fmt.Println("未跨时段", totalPrice)
		//			totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//		}
		//	} else {
		//		//fmt.Println("55555")
		//		//fmt.Println(haltHour, halfHourPrice, lastDiffMinute, 30-lastDiffMinute)
		//		//fmt.Println("-----000", timeV)
		//		if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format("15:04") >= timeV.End && lastDiffMinute != 0 {
		//			//fmt.Println("666666")
		//			minute := decimal.NewFromInt(30).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
		//			//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//			minuteTotalPrice, _ := decimal.NewFromInt(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//			//fmt.Println(minute, minutePrice, minuteTotalPrice, ".........1111=====")
		//			totalHourPrice, _ := decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour - 1)).Float64()
		//			totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//			totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
		//		} else {
		//			//fmt.Println("77777777")
		//			totalHourPrice, _ := decimal.NewFromFloat(halfHourPrice).Mul(decimal.NewFromInt(haltHour)).Float64()
		//			totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//			//fmt.Println("未跨时段 当日", halfHourPrice, haltHour, totalPrice)
		//		}
		//	}
		//	//fmt.Println("===================================")
		//	price += totalPrice
		//	periods[int64(timeV.Index)] += totalPrice
		//} else {
		//	//fmt.Println("---11")
		//	periods[int64(timeV.Index)] += 0
		//}
		////fmt.Println(price, "......////")
		//periodList = append(periodList, &PeriodList{
		//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
		//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
		//	Index:      timeV.Index,
		//	Duration:   totalMinute,
		//	Start:      timeV.Start,
		//	End:        timeV.End,
		//	Price:      timeV.Price,
		//	TotalPrice: totalPrice,
		//})

	}
	//fmt.Println("=========11111---", hourPeriodList)
	//
	//fmt.Println("=========222222---", hourPeriodLists)
	allHourPeriodList := make(map[string]HourPeriodList, 0)
	//var cdP float64
	//fmt.Println(allHourPeriodList, "========")
	for date, value := range hourPeriodLists {
		var tempPrice float64
		var tempDuration int64
		var tempHourPrice float64
		var tempStartDate string
		var tempEndDate string
		for _, val := range value {
			//fmt.Println(val.StartDate, val.EndDate, "时长=", val.Duration, "----", val.HourPrice, "总价", val.Price)
			if tempStartDate == "" {
				tempStartDate = val.StartDate
			}
			tempEndDate = val.EndDate
			tempPrice += val.Price
			tempDuration += val.Duration
			tempHourPrice = val.HourPrice
			price += val.Price
		}
		allHourPeriodList[date] = HourPeriodList{StartDate: tempStartDate, EndDate: tempEndDate, HourPrice: tempHourPrice, Price: tempPrice, Duration: tempDuration}
	}
	//fmt.Println(price, "金额", allHourPeriodList, "..----,,,")
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, allHourPeriodList
}

// 一小时计费 (一小时为节点，以一小时起始点计算价格，如果时间段内计算时长小于一小时，则累计至下个时段，
// 在下个时段内扣除上个时段少于的时间[如果跨时段按跨时段分钟计费,否则按一小时算]，然后在计算下周期时段价格，依次类推)
func (l *Charging) computeHourModeTwo(times []periodTimes) (float64, map[string]HourPeriodList) {
	//periods := make(map[int64]float64) // 时段费用详情
	var price float64
	//var lastDiffMinute int64
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	hourPeriodLists := make(map[string][]HourPeriodList, 0)

	var hourLastMinute int64
	lenIndex := len(times)
	for index, timeV := range times {
		//var totalPrice float64

		minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()

		var tempStartTime time.Time
		var tempStartTimeD time.Time
		tempStartTime = timeV.StartTime
		for {
			if timeV.EndTime.Unix() <= tempStartTime.Unix() {
				break
			}
			if hourLastMinute > 0 {
				// 跨时段按分钟计费
				tempFirstMinutes := int64(timeV.EndTime.Sub(tempStartTime).Minutes())
				if tempFirstMinutes < 60 {
					tempStartTimeD = tempStartTime.Add(time.Duration(tempFirstMinutes) * time.Minute)
					hourTotalPrice, _ := decimal.NewFromInt(tempFirstMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: tempFirstMinutes, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTime.Add(time.Duration(tempFirstMinutes) * time.Minute).Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					if tempFirstMinutes > hourLastMinute {
						hourLastMinute = 0
					} else {
						hourLastMinute -= tempFirstMinutes
					}

				} else if tempFirstMinutes >= 60 && hourLastMinute != 0 {
					tempStartTimeD = tempStartTime.Add(time.Duration(hourLastMinute) * time.Minute)
					hourTotalPrice, _ := decimal.NewFromInt(hourLastMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: hourLastMinute, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTime.Add(time.Duration(hourLastMinute) * time.Minute).Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					hourLastMinute = 0
				} else {
					hourTotalPrice, _ := decimal.NewFromInt(hourLastMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
						HourPeriodList{Price: hourTotalPrice, Duration: hourLastMinute, StartDate: tempStartTime.Add(time.Duration(-hourLastMinute) * time.Minute).Format("2006-01-02 15:04"), EndDate: tempStartTime.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					tempStartTimeD = tempStartTime.Add(time.Duration(hourLastMinute) * time.Minute)
					hourLastMinute = 0
				}

				//hourTotalPrice, _ := decimal.NewFromInt(hourLastMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
				//hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
				//	HourPeriodList{Price: hourTotalPrice, Duration: hourLastMinute, StartDate: tempStartTime.Add(time.Duration(-hourLastMinute) * time.Minute).Format("2006-01-02 15:04"), EndDate: tempStartTime.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
				//hourLastMinute = 0
				//fmt.Println("11111111", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime), "===", hourTotalPrice)
			} else {
				tempStartTimeD = tempStartTime.Add(3600 * time.Second)
				// 判断 相差值，必须为数组最后之前的数据，如果是数组最后一个则按小时收费
				if tempStartTimeD.Unix() > timeV.EndTime.Unix() && lenIndex-1 > index {
					minute := timeV.EndTime.Sub(tempStartTime).Minutes()
					hourLastMinute = 60 - int64(minute)
					hourTotalPrice, _ := decimal.NewFromFloat(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					tempStartTimeD = timeV.EndTime
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: int64(minute), StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: timeV.EndTime.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//fmt.Println("22222222", tempStartTime.Format(time.DateTime), "===", hourTotalPrice, minute, timeV.EndTime.Format(time.DateTime), hourLastMinute)
				} else {
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += timeV.Price
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: timeV.Price, Duration: 60, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//fmt.Println("333333", tempStartTimeD.Format(time.DateTime), "===", halfHourPrice)
				}
			}
			tempStartTime = tempStartTimeD
		}
		//for {
		//	if timeV.EndTime.Unix() <= tempStartTime.Unix() {
		//		break
		//	}
		//	if hourLastMinute > 0 {
		//		hourTotalPrice, _ := decimal.NewFromInt(hourLastMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//		hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
		//		hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
		//			HourPeriodList{Price: hourTotalPrice, Duration: hourLastMinute, StartDate: tempStartTime.Add(time.Duration(-hourLastMinute) * time.Minute).Format("2006-01-02 15:04"), EndDate: tempStartTime.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
		//		hourLastMinute = 0
		//		tempStartTimeD = tempStartTime.Add(time.Duration(hourLastMinute) * time.Minute)
		//	} else {
		//		tempStartTimeD = tempStartTime.Add(3600 * time.Second)
		//		// 判断 相差值，必须为数组最后之前的数据，如果是数组最后一个则按半小时收费
		//		if tempStartTimeD.Unix() > timeV.EndTime.Unix() && lenIndex-1 > index {
		//			minute := timeV.EndTime.Sub(tempStartTime).Minutes()
		//			hourLastMinute = 60 - int64(minute)
		//			hourTotalPrice, _ := decimal.NewFromFloat(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//			hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
		//			hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: int64(minute), StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: timeV.EndTime.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
		//		} else {
		//			hourPeriodList[tempStartTime.Format("2006-01-02 15")] += timeV.Price
		//			hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: timeV.Price, Duration: 60, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
		//		}
		//	}
		//	tempStartTime = tempStartTimeD
		//}

		//// 差值按分钟计算计费
		//if lastDiffMinute > 0 {
		//	lastStartTime := timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		//	if lastStartTime.Unix() >= timeV.EndTime.Unix() {
		//		lastStartTime = timeV.EndTime
		//	}
		//	duration := lastStartTime.Sub(timeV.StartTime)
		//	totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
		//	//fmt.Println(".......-----", lastDiffMinute, duration, totalMinute)
		//	if totalMinute > 0 {
		//		//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//		minuteTotalPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
		//		timeV.StartTime = lastStartTime
		//		//fmt.Println("跨时段按分钟", lastDiffMinute, duration, totalMinute, totalPrice, timeV.StartTime.Format(time.DateTime))
		//	}
		//}
		//duration := timeV.EndTime.Sub(timeV.StartTime)
		//totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
		//if totalMinute > 0 {
		//	hour := decimal.NewFromInt(totalMinute).Div(decimal.NewFromInt(60)).Ceil().IntPart()
		//	maxMinute := decimal.NewFromInt(hour).Mul(decimal.NewFromInt(60)).IntPart()
		//	var tempEndTime time.Time
		//	if timeV.Start > timeV.End && timeV.StartTime.Format("15:04") >= timeV.Start {
		//		tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", timeV.EndTime.Add(time.Duration(86400)*time.Second).Format("2006-01-02"), timeV.End), time.Local)
		//	} else {
		//		tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", timeV.EndTime.Format("2006-01-02"), timeV.End), time.Local)
		//	}
		//
		//	//fmt.Println("..;;;;;;;;;;", timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format(time.DateTime), tempEndTime.Format(time.DateTime))
		//
		//	if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Unix() >= tempEndTime.Unix() {
		//		lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(totalMinute)).IntPart()
		//		//fmt.Println(lastDiffMinute, "....;;;;;;;;;")
		//	} else {
		//		lastDiffMinute = 0
		//	}
		//	//fmt.Println(lastDiffMinute, "lastDiffMinute====", timeV.StartTime.Format(time.DateTime), timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format(time.DateTime), timeV.End, ".....", tempEndTime.Format(time.DateTime))
		//	//fmt.Println(lastDiffMinute, "======222")
		//	//halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
		//	// 开始时间大于结束时间 跨天
		//	if timeV.Start > timeV.End {
		//		//fmt.Println("1111")
		//		// 跨时段
		//		if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format("15:04") >= timeV.End && lastDiffMinute != 0 {
		//			if totalMinute < 60 && timeV.EndTime.Format("15:04") < timeV.End {
		//				//fmt.Println("22222")
		//				totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour)).Float64()
		//				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//				//fmt.Println("不足一小时按一小时算", totalPrice)
		//			} else {
		//				//fmt.Println("33333")
		//				minute := decimal.NewFromInt(60).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
		//				//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//				minuteTotalPrice, _ := decimal.NewFromInt(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//
		//				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
		//				totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour - 1)).Float64()
		//				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//			}
		//			//	fmt.Println("跨时段", minuteTotalPrice, totalPrice, lastDiffMinute)
		//		} else {
		//			//fmt.Println("444444")
		//			totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour)).Float64()
		//			//fmt.Println("未跨时段", totalPrice)
		//			totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//		}
		//	} else {
		//		//fmt.Println("55555")
		//		//fmt.Println(haltHour, halfHourPrice, lastDiffMinute, 30-lastDiffMinute)
		//		//fmt.Println("-----000", timeV)
		//		if timeV.StartTime.Add(time.Duration(maxMinute)*time.Minute).Format("15:04") >= timeV.End && lastDiffMinute != 0 {
		//			//fmt.Println("666666")
		//			minute := decimal.NewFromInt(60).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
		//			//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//			minuteTotalPrice, _ := decimal.NewFromInt(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//			//fmt.Println(minute, minutePrice, minuteTotalPrice, ".........1111=====")
		//			totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour - 1)).Float64()
		//			totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//			totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
		//		} else {
		//			//fmt.Println("77777777")
		//			totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(hour)).Float64()
		//			totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//			//fmt.Println("未跨时段 当日", halfHourPrice, haltHour, totalPrice)
		//		}
		//	}
		//	//fmt.Println("===================================")
		//	price += totalPrice
		//	periods[int64(timeV.Index)] += totalPrice
		//} else {
		//	//fmt.Println("---11")
		//	periods[int64(timeV.Index)] += 0
		//}
		////fmt.Println(price, "......////")
		//periodList = append(periodList, &PeriodList{
		//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
		//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
		//	Index:      timeV.Index,
		//	Duration:   totalMinute,
		//	Start:      timeV.Start,
		//	End:        timeV.End,
		//	Price:      timeV.Price,
		//	TotalPrice: totalPrice,
		//})

	}

	//fmt.Println(hourPeriodLists, "--------")
	allHourPeriodList := make(map[string]HourPeriodList, 0)
	//var cdP float64
	//fmt.Println(allHourPeriodList, "========")
	for date, value := range hourPeriodLists {
		var tempPrice float64
		var tempDuration int64
		var tempHourPrice float64
		var tempStartDate string
		var tempEndDate string
		for _, val := range value {
			//fmt.Println(val.StartDate, val.EndDate, "时长=", val.Duration, "----", val.HourPrice, "总价", val.Price)
			if tempStartDate == "" {
				tempStartDate = val.StartDate
			}
			tempEndDate = val.EndDate
			tempPrice += val.Price
			tempDuration += val.Duration
			tempHourPrice = val.HourPrice
			price += val.Price
		}
		allHourPeriodList[date] = HourPeriodList{StartDate: tempStartDate, EndDate: tempEndDate, HourPrice: tempHourPrice, Price: tempPrice, Duration: tempDuration}
	}
	//fmt.Println(price, "金额", allHourPeriodList, "..----,,,")

	//fmt.Println(price, "..----,,,")
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, allHourPeriodList
}

// 按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费。[如果跨时段按跨时段分钟计费]
func (l *Charging) computeHalfHourOrMinuteModeTwo(times []periodTimes) (float64, map[string]HourPeriodList) {
	//periods := make(map[int64]float64) // 时段费用详情
	var price float64
	//var lastDiffMinute int64
	var firstHalfHourFor int64
	firstHalfHourFor = 30
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	hourPeriodLists := make(map[string][]HourPeriodList, 0)
	lenIndex := len(times)
	for index, timeV := range times {
		//if lastDiffMinute != 0 {
		//	timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		//}
		halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
		minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()

		tempStartTime := timeV.StartTime
		for {
			if timeV.EndTime.Unix() <= tempStartTime.Unix() {
				break
			}
			var tempStartTimeD time.Time
			if firstHalfHourFor > 0 {
				tempFirstMinute := int64(timeV.EndTime.Sub(tempStartTime).Minutes())
				if tempFirstMinute < 30 {
					//lastDiffMinute = 30 - tempFirstMinute
					tempStartTimeD = tempStartTime.Add(time.Duration(tempFirstMinute) * time.Minute)
					hourTotalPrice, _ := decimal.NewFromInt(tempFirstMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: tempFirstMinute, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//firstHalfHourFor -= tempFirstMinute
				} else {
					tempStartTimeD = tempStartTime.Add(30 * time.Minute)
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += halfHourPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: halfHourPrice, Duration: 30, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//firstHalfHourFor = 0
				}
				firstHalfHourFor = 0
				//fmt.Println(halfHourPrice, "-------", tempStartTimeD.Format(time.DateTime))
				tempStartTime = tempStartTimeD
			} else {
				tempStartTimeD = tempStartTime.Add(1800 * time.Second)
				// 判断 相差值，必须为数组最后之前的数据，如果是数组最后一个则按半小时收费
				if tempStartTimeD.Unix() > timeV.EndTime.Unix() && lenIndex-1 > index {
					minute := timeV.EndTime.Sub(tempStartTime).Minutes()
					hourTotalPrice, _ := decimal.NewFromFloat(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					tempStartTimeD = tempStartTime.Add(time.Duration(minute) * time.Minute)
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: int64(minute), StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//lastDiffMinute = 30 - int64(minute)
					//fmt.Println(hourTotalPrice, "=|||||||===", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime), minute)
					tempStartTime = tempStartTimeD
				} else {
					var hourTotalPrice float64
					tempEndMinutes := decimal.NewFromFloat(timeV.EndTime.Sub(tempStartTime).Minutes()).IntPart()
					var tempDuration int64
					if tempEndMinutes < 30 {
						tempDuration = tempEndMinutes
						hourTotalPrice, _ = decimal.NewFromInt(tempEndMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
						tempStartTimeD = tempStartTime.Add(time.Duration(tempEndMinutes) * time.Minute)
						//fmt.Println("------1", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime), tempEndMinutes)
					} else {
						tempDuration = 30
						hourTotalPrice, _ = decimal.NewFromFloat(30).Mul(decimal.NewFromFloat(minutePrice)).Float64()
						tempStartTimeD = tempStartTime.Add(time.Duration(tempDuration) * time.Minute)
						//fmt.Println("------0000", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime))
					}
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: tempDuration, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//fmt.Println(hourTotalPrice, "====", tempStartTime.Format("2006-01-02 15"), "|||||||||||", tempStartTime.Format(time.DateTime))
					tempStartTime = tempStartTimeD
				}
			}
			//}
		}

		//var totalPrice float64
		//duration := timeV.EndTime.Sub(timeV.StartTime)
		//if duration.Minutes() > 0 {
		//	if firstHalfHour != 0 {
		//		haltHour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(30)).Ceil().IntPart()
		//		maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
		//		lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
		//		tempMinute := maxMinute - 30 - lastDiffMinute
		//		totalPrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(tempMinute)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(halfHourPrice)).Float64()
		//		price += totalPrice
		//		periods[int64(timeV.Index)] += totalPrice
		//		firstHalfHour = 0 // 首个半小时计费，计费完后置为0
		//	} else {
		//		//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//		price += totalPrice
		//		periods[int64(timeV.Index)] += totalPrice
		//		lastDiffMinute = 0 // 分钟计费，差值分钟置为0
		//	}
		//} else {
		//	periods[int64(timeV.Index)] += 0
		//}
		//periodList = append(periodList, &PeriodList{
		//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
		//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
		//	Index:      timeV.Index,
		//	Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
		//	Start:      timeV.Start,
		//	End:        timeV.End,
		//	Price:      timeV.Price,
		//	TotalPrice: totalPrice,
		//})

	}

	//fmt.Println(hourPeriodLists, "--------")
	allHourPeriodList := make(map[string]HourPeriodList, 0)
	//var cdP float64
	//fmt.Println(allHourPeriodList, "========")
	for date, value := range hourPeriodLists {
		var tempPrice float64
		var tempDuration int64
		var tempHourPrice float64
		var tempStartDate string
		var tempEndDate string
		for _, val := range value {
			//fmt.Println(val.StartDate, val.EndDate, "时长=", val.Duration, "----", val.HourPrice, "总价", val.Price)
			if tempStartDate == "" {
				tempStartDate = val.StartDate
			}
			tempEndDate = val.EndDate
			tempPrice += val.Price
			tempDuration += val.Duration
			tempHourPrice = val.HourPrice
			price += val.Price
		}
		allHourPeriodList[date] = HourPeriodList{StartDate: tempStartDate, EndDate: tempEndDate, HourPrice: tempHourPrice, Price: tempPrice, Duration: tempDuration}
	}
	//fmt.Println(price, "金额", allHourPeriodList, "..----,,,")

	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, allHourPeriodList
}

// 按小时计费开台不足1小时按小时计费，超过1小时按分钟计费。[如果跨时段按跨时段分钟计费]
func (l *Charging) computeHourOrMinuteModeTwo(times []periodTimes) (float64, map[string]HourPeriodList) {
	//periods := make(map[int64]float64) // 时段费用详情
	var price float64
	//var lastDiffMinute int64
	var firstHourFor int64
	firstHourFor = 60
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	hourPeriodLists := make(map[string][]HourPeriodList, 0)
	lenIndex := len(times)
	for index, timeV := range times {
		//if lastDiffMinute != 0 {
		//	timeV.StartTime = timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
		//}
		//halfHourPrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
		minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()

		//
		tempStartTime := timeV.StartTime
		for {
			if timeV.EndTime.Unix() <= tempStartTime.Unix() {
				break
			}
			var tempStartTimeD time.Time
			if firstHourFor > 0 {
				tempFirstMinute := int64(timeV.EndTime.Sub(tempStartTime).Minutes())
				if tempFirstMinute < 60 {
					tempStartTimeD = tempStartTime.Add(time.Duration(tempFirstMinute) * time.Minute)
					hourTotalPrice, _ := decimal.NewFromInt(tempFirstMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: tempFirstMinute, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
				} else {
					tempStartTimeD = tempStartTime.Add(60 * time.Minute)
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += timeV.Price
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: timeV.Price, Duration: 60, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
				}
				firstHourFor = 0
				//fmt.Println(halfHourPrice, "-------", tempStartTimeD.Format(time.DateTime))
				tempStartTime = tempStartTimeD
			} else {
				tempStartTimeD = tempStartTime.Add(3600 * time.Second)
				// 判断 相差值，必须为数组最后之前的数据，如果是数组最后一个则按半小时收费
				if tempStartTimeD.Unix() > timeV.EndTime.Unix() && lenIndex-1 > index {
					minute := timeV.EndTime.Sub(tempStartTime).Minutes()
					hourTotalPrice, _ := decimal.NewFromFloat(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					tempStartTimeD = tempStartTime.Add(time.Duration(minute) * time.Minute)
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: int64(minute), StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//lastDiffMinute = 30 - int64(minute)
					//fmt.Println(hourTotalPrice, "=|||||||===", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime), minute)
					tempStartTime = tempStartTimeD
				} else {
					var hourTotalPrice float64
					tempEndMinutes := decimal.NewFromFloat(timeV.EndTime.Sub(tempStartTime).Minutes()).IntPart()
					var tempDuration int64
					if tempEndMinutes < 60 {
						tempDuration = tempEndMinutes
						hourTotalPrice, _ = decimal.NewFromInt(tempEndMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
						tempStartTimeD = tempStartTime.Add(time.Duration(tempEndMinutes) * time.Minute)
						//fmt.Println("------1", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime), tempEndMinutes)
					} else {
						tempDuration = 60
						hourTotalPrice, _ = decimal.NewFromFloat(60).Mul(decimal.NewFromFloat(minutePrice)).Float64()
						tempStartTimeD = tempStartTime.Add(time.Duration(tempDuration) * time.Minute)
						//fmt.Println("------0000", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime))
					}
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: tempDuration, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//fmt.Println(hourTotalPrice, "====", tempStartTime.Format("2006-01-02 15"), "|||||||||||", tempStartTime.Format(time.DateTime))
					tempStartTime = tempStartTimeD
				}
			}
			//}
		}

		//var totalPrice float64
		//duration := timeV.EndTime.Sub(timeV.StartTime)
		//if duration.Minutes() > 0 {
		//	if firstHalfHour != 0 {
		//		haltHour := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(30)).Ceil().IntPart()
		//		maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
		//		lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromFloat(duration.Minutes())).IntPart()
		//		tempMinute := maxMinute - 30 - lastDiffMinute
		//		totalPrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(tempMinute)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(halfHourPrice)).Float64()
		//		price += totalPrice
		//		periods[int64(timeV.Index)] += totalPrice
		//		firstHalfHour = 0 // 首个半小时计费，计费完后置为0
		//	} else {
		//		//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(duration.Minutes()).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//		price += totalPrice
		//		periods[int64(timeV.Index)] += totalPrice
		//		lastDiffMinute = 0 // 分钟计费，差值分钟置为0
		//	}
		//} else {
		//	periods[int64(timeV.Index)] += 0
		//}
		//periodList = append(periodList, &PeriodList{
		//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
		//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
		//	Index:      timeV.Index,
		//	Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
		//	Start:      timeV.Start,
		//	End:        timeV.End,
		//	Price:      timeV.Price,
		//	TotalPrice: totalPrice,
		//})

	}

	//fmt.Println(hourPeriodLists, "--------")
	allHourPeriodList := make(map[string]HourPeriodList, 0)
	//var cdP float64
	//fmt.Println(allHourPeriodList, "========")
	for date, value := range hourPeriodLists {
		var tempPrice float64
		var tempDuration int64
		var tempHourPrice float64
		var tempStartDate string
		var tempEndDate string
		for _, val := range value {
			//fmt.Println(val.StartDate, val.EndDate, "时长=", val.Duration, "----", val.HourPrice, "总价", val.Price)
			if tempStartDate == "" {
				tempStartDate = val.StartDate
			}
			tempEndDate = val.EndDate
			tempPrice += val.Price
			tempDuration += val.Duration
			tempHourPrice = val.HourPrice
			price += val.Price
		}
		allHourPeriodList[date] = HourPeriodList{StartDate: tempStartDate, EndDate: tempEndDate, HourPrice: tempHourPrice, Price: tempPrice, Duration: tempDuration}
	}
	//fmt.Println(price, "金额", allHourPeriodList, "..----,,,")

	//var price float64
	//var lastDiffMinute int64
	//var firstHalfHour int64
	//firstHalfHour = 60
	//firstHalfHourFor := 60
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	//var hourLastMinute int64
	//for _, timeV := range times {
	//	var totalPrice float64
	//
	//	minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
	//
	//	// 每小时消费
	//	var tempStartTime time.Time
	//	var tempStartTimeD time.Time
	//	if hourLastMinute != 0 {
	//		tempStartTime = timeV.StartTime.Add(time.Duration(hourLastMinute) * time.Minute)
	//	} else {
	//		tempStartTime = timeV.StartTime
	//	}
	//	for {
	//		if timeV.EndTime.Unix() <= tempStartTime.Unix() {
	//			break
	//		}
	//		if firstHalfHourFor != 0 {
	//			tempStartTimeD = tempStartTime.Add(60 * time.Minute)
	//			hourPeriodList[tempStartTime.Format("2006-01-02 15")] += timeV.Price
	//			//fmt.Println(timeV.Price, "-------", tempStartTimeD.Format(time.DateTime))
	//			//tempStartTime = tempStartTimeD
	//			firstHalfHourFor = 0
	//		} else {
	//			originStartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", tempStartTime.Format("2006-01-02 15:00:00"), time.Local)
	//			tempMinutes := decimal.NewFromFloat(tempStartTime.Sub(originStartTime).Minutes()).IntPart()
	//			if tempMinutes != 0 {
	//				hourTotalPrice, _ := decimal.NewFromInt(60 - tempMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
	//				tempStartTimeD = tempStartTime.Add(time.Duration(60-tempMinutes) * time.Minute)
	//				hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
	//				//fmt.Println(hourTotalPrice, "====", tempStartTime.Format("2006-01-02 15"), "......", tempStartTime.Format(time.DateTime))
	//				//tempStartTime = tempStartTimeD
	//			} else {
	//				if tempStartTime.Unix() <= timeV.EndTime.Unix() {
	//					var hourTotalPrice float64
	//					tempEndMinutes := decimal.NewFromFloat(timeV.EndTime.Sub(tempStartTime).Minutes()).IntPart()
	//					if tempEndMinutes < 60 {
	//						hourTotalPrice, _ = decimal.NewFromInt(tempEndMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
	//					} else {
	//						hourTotalPrice, _ = decimal.NewFromFloat(60).Mul(decimal.NewFromFloat(minutePrice)).Float64()
	//					}
	//					hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
	//					//fmt.Println(hourTotalPrice, "====", tempStartTime.Format("2006-01-02 15"), "|||||||||||", tempStartTime.Format(time.DateTime))
	//					tempStartTimeD = tempStartTime.Add(3600 * time.Second)
	//				}
	//			}
	//		}
	//
	//		tempStartTime = tempStartTimeD
	//	}
	//
	//	// 差值按分钟计算计费
	//	if lastDiffMinute > 0 {
	//		lastStartTime := timeV.StartTime.Add(time.Duration(lastDiffMinute) * time.Minute)
	//		if lastStartTime.Unix() >= timeV.EndTime.Unix() {
	//			lastStartTime = timeV.EndTime
	//		}
	//		duration := lastStartTime.Sub(timeV.StartTime)
	//		totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
	//		if totalMinute > 0 {
	//			//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
	//			minuteTotalPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
	//			totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
	//			timeV.StartTime = lastStartTime
	//			lastDiffMinute = 0
	//			//fmt.Println("跨时段按分钟", lastDiffMinute, duration, totalMinute, totalPrice, timeV.StartTime.Format(time.DateTime))
	//		}
	//	}
	//
	//	duration := timeV.EndTime.Sub(timeV.StartTime)
	//	totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
	//	var surplusMinute int64
	//	if totalMinute > 0 {
	//		if firstHalfHour != 0 {
	//			hour := decimal.NewFromInt(totalMinute).Div(decimal.NewFromInt(60)).Ceil().IntPart()
	//			maxMinute := decimal.NewFromInt(hour).Mul(decimal.NewFromInt(60)).IntPart()
	//			if hour >= 1 && totalMinute >= 60 {
	//				totalHourPrice, _ := decimal.NewFromFloat(timeV.Price).Mul(decimal.NewFromInt(1)).Float64()
	//				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
	//				surplusMinute = totalMinute - 60
	//			}
	//			if hour == 1 && totalMinute < 60 {
	//				lastDiffMinute = decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(totalMinute)).IntPart()
	//				surplusMinute = decimal.NewFromInt(60).Sub(decimal.NewFromInt(lastDiffMinute)).IntPart()
	//				//fmt.Println("检查下个时段是否有数据，没有则在当前时段计费下个时段", index, "===", times[index+1])
	//			}
	//			if surplusMinute >= 1 {
	//				//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
	//				minuteTotalPrice, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
	//				totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
	//			}
	//			//price += totalPrice
	//			periods[int64(timeV.Index)] += totalPrice
	//			firstHalfHour = 0 // 首个半小时计费，计费完后置为0
	//		} else {
	//			//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
	//			minuteTotalPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
	//			totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
	//			//price += totalPrice
	//			periods[int64(timeV.Index)] += totalPrice
	//			lastDiffMinute = 0 // 分钟计费，差值分钟置为0
	//		}
	//	} else {
	//		//price += totalPrice
	//		periods[int64(timeV.Index)] += 0
	//	}
	//	price += totalPrice
	//	periodList = append(periodList, &PeriodList{
	//		StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
	//		EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
	//		Index:      timeV.Index,
	//		Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
	//		Start:      timeV.Start,
	//		End:        timeV.End,
	//		Price:      timeV.Price,
	//		TotalPrice: totalPrice,
	//	})
	//
	//}
	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, allHourPeriodList
}

// 以10分钟作为一个收费周期，第1分钟后开始计费，如果时间段内计算时长小于收费周期，累计至下个时段
// 并在下个时段内扣除上个时段少于的时间[如果跨时段按跨时段分钟计费]，然后在计算下周期时段价格，依次类推)
func (l *Charging) computeCycleAndMinuteModeTwo(times []periodTimes) (float64, map[string]HourPeriodList) {
	//periods := make(map[int64]float64) // 时段费用详情
	var price float64
	var hourLastMinute int64
	//periodList := make([]*PeriodList, 0)
	//hourPeriodList := make(map[string]float64, 0)
	hourPeriodLists := make(map[string][]HourPeriodList, 0)
	lenIndex := len(times)
	for index, timeV := range times {
		//var totalPrice float64
		minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		var cyclePrice float64
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
		// 周期消费金额
		// 每小时消费
		var tempStartTime time.Time
		var tempStartTimeD time.Time
		tempStartTime = timeV.StartTime

		for {
			if timeV.EndTime.Unix() <= tempStartTime.Unix() {
				break
			}
			if hourLastMinute > 0 {
				// 跨时段按分钟计费
				tempFirstMinutes := int64(timeV.EndTime.Sub(tempStartTime).Minutes())
				if tempFirstMinutes < l.cycleMinute {
					tempStartTimeD = tempStartTime.Add(time.Duration(tempFirstMinutes) * time.Minute)
					hourTotalPrice, _ := decimal.NewFromInt(tempFirstMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: tempFirstMinutes, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTime.Add(time.Duration(tempFirstMinutes) * time.Minute).Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					if tempFirstMinutes > hourLastMinute {
						hourLastMinute = 0
					} else {
						hourLastMinute -= tempFirstMinutes
					}
					//fmt.Println(tempStartTimeD.Format(time.DateTime), "------111111--------")
				} else if tempFirstMinutes >= l.cycleMinute && hourLastMinute != 0 {
					tempStartTimeD = tempStartTime.Add(time.Duration(hourLastMinute) * time.Minute)
					hourTotalPrice, _ := decimal.NewFromInt(hourLastMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: hourLastMinute, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTime.Add(time.Duration(hourLastMinute) * time.Minute).Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					hourLastMinute = 0
					//fmt.Println(tempStartTimeD.Format(time.DateTime), "------222222--------")
				} else {
					hourTotalPrice, _ := decimal.NewFromInt(hourLastMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")],
						HourPeriodList{Price: hourTotalPrice, Duration: hourLastMinute, StartDate: tempStartTime.Add(time.Duration(-hourLastMinute) * time.Minute).Format("2006-01-02 15:04"), EndDate: tempStartTime.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					tempStartTimeD = tempStartTime.Add(time.Duration(hourLastMinute) * time.Minute)
					hourLastMinute = 0
					//fmt.Println(tempStartTimeD.Format(time.DateTime), "------333333--------")
				}
			} else {
				tempStartTimeD = tempStartTime.Add(time.Duration(l.cycleMinute) * time.Minute)
				//fmt.Println(tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime), timeV.EndTime.Format(time.DateTime), "---", l.cycleMinute, "*********************", lenIndex-1, index)

				// 判断 相差值，必须为数组最后之前的数据，如果是数组最后一个则按半小时收费
				if tempStartTimeD.Unix() > timeV.EndTime.Unix() && lenIndex-1 >= index {
					//fmt.Println(tempStartTimeD.Format(time.DateTime), "-------+++++++++++++111111")
					minute := timeV.EndTime.Sub(tempStartTime).Minutes()
					hourLastMinute = l.cycleMinute - int64(minute)
					hourTotalPrice, _ := decimal.NewFromFloat(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
					tempStartTimeD = timeV.EndTime
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: int64(minute), StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: timeV.EndTime.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//fmt.Println("22222222", tempStartTime.Format(time.DateTime), "===", hourTotalPrice, minute, timeV.EndTime.Format(time.DateTime), hourLastMinute)
				} else {
					//fmt.Println(tempStartTimeD.Format(time.DateTime), "-------+++++++++++++2222")
					//hourPeriodList[tempStartTime.Format("2006-01-02 15")] += cyclePrice
					hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: cyclePrice, Duration: l.cycleMinute, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
					//fmt.Println("333333", tempStartTimeD.Format(time.DateTime), "===", halfHourPrice)
				}
			}
			tempStartTime = tempStartTimeD
		}

		//for {
		//	if timeV.EndTime.Unix() <= tempStartTime.Unix() {
		//		break
		//	}
		//
		//	originStartTime, _ := time.ParseInLocation("2006-01-02 15:04:05", tempStartTime.Format("2006-01-02 15:00:00"), time.Local)
		//	tempMinutes := decimal.NewFromFloat(tempStartTime.Sub(originStartTime).Minutes()).IntPart()
		//	if tempMinutes != 0 && tempMinutes < l.cycleMinute {
		//		diffMinute := decimal.NewFromInt(l.cycleMinute).Sub(decimal.NewFromInt(tempMinutes)).IntPart()
		//		hourTotalPrice, _ := decimal.NewFromInt(diffMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//		tempStartTimeD = tempStartTime.Add(time.Duration(diffMinute) * time.Minute)
		//		hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
		//		hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: tempMinutes, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
		//		//fmt.Println("========1", tempStartTime.Format(time.DateTime), tempStartTimeD.Format(time.DateTime))
		//	} else {
		//		if tempStartTime.Unix() <= timeV.EndTime.Unix() {
		//			var hourTotalPrice float64
		//			tempMinutes := int64(timeV.EndTime.Sub(tempStartTime).Minutes())
		//			if tempMinutes < l.cycleMinute {
		//				hourTotalPrice, _ = decimal.NewFromInt(tempMinutes).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//			} else {
		//				hourTotalPrice = cyclePrice
		//			}
		//			hourPeriodList[tempStartTime.Format("2006-01-02 15")] += hourTotalPrice
		//			tempStartTimeD = tempStartTime.Add(time.Duration(l.cycleMinute) * time.Minute)
		//			hourPeriodLists[tempStartTime.Format("2006-01-02 15")] = append(hourPeriodLists[tempStartTime.Format("2006-01-02 15")], HourPeriodList{Price: hourTotalPrice, Duration: tempMinutes, StartDate: tempStartTime.Format("2006-01-02 15:04"), EndDate: tempStartTimeD.Format("2006-01-02 15:04"), HourPrice: timeV.Price})
		//
		//		}
		//	}
		//	tempStartTime = tempStartTimeD
		//}

		//duration := timeV.EndTime.Sub(timeV.StartTime)
		//totalMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
		//var surplusMinute int64
		////var cyclePrice float64
		//if totalMinute > 0 {
		//	//minutePrice, _ := decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(60)).Float64()
		//	cycle := decimal.NewFromFloat(duration.Minutes()).Div(decimal.NewFromInt(l.cycleMinute)).Ceil().IntPart()
		//
		//	//// 半小时
		//	//if l.cycleMinute == 30 {
		//	//	cyclePrice, _ = decimal.NewFromFloat(timeV.Price).Div(decimal.NewFromInt(2)).Float64()
		//	//}
		//	//// 一小时
		//	//if l.cycleMinute == 60 {
		//	//	cyclePrice = timeV.Price
		//	//}
		//	//// 周期计费
		//	//if l.cycleMinute != 30 && l.cycleMinute != 60 {
		//	//	cyclePrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
		//	//}
		//
		//	if cycle >= 1 && totalMinute >= l.cycleMinute {
		//		totalHourPrice, _ := decimal.NewFromFloat(cyclePrice).Mul(decimal.NewFromInt(cycle - 1)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//		surplusMinute = totalMinute - ((cycle - 1) * l.cycleMinute)
		//	}
		//	if cycle == 1 && totalMinute < l.cycleMinute {
		//		totalHourPrice, _ := decimal.NewFromInt(totalMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(totalHourPrice)).Float64()
		//		//fmt.Println("检查下个时段是否有数据，没有则在当前时段计费下个时段", index, "===", times[index+1])
		//		surplusMinute = 0
		//	}
		//	if surplusMinute >= 1 {
		//		minuteTotalPrice, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
		//		totalPrice, _ = decimal.NewFromFloat(totalPrice).Add(decimal.NewFromFloat(minuteTotalPrice)).Float64()
		//	}
		//	periods[int64(timeV.Index)] += totalPrice
		//} else {
		//	periods[int64(timeV.Index)] += 0
		//}
		//price += totalPrice
		//periodList = append(periodList, &PeriodList{
		//	StartTime:  timeV.StartTime.Format("2006-01-02 15:04"),
		//	EndTime:    timeV.EndTime.Format("2006-01-02 15:04"),
		//	Index:      timeV.Index,
		//	Duration:   decimal.NewFromFloat(duration.Minutes()).IntPart(),
		//	Start:      timeV.Start,
		//	End:        timeV.End,
		//	Price:      timeV.Price,
		//	TotalPrice: totalPrice,
		//})
	}

	allHourPeriodList := make(map[string]HourPeriodList, 0)
	//var cdP float64
	for date, value := range hourPeriodLists {
		var tempPrice float64
		var tempDuration int64
		var tempHourPrice float64
		var tempStartDate string
		var tempEndDate string
		for _, val := range value {
			//fmt.Println(val.StartDate, val.EndDate, "时长=", val.Duration, "----", val.HourPrice, "总价", val.Price)
			if tempStartDate == "" {
				tempStartDate = val.StartDate
			}
			tempEndDate = val.EndDate
			tempPrice += val.Price
			tempDuration += val.Duration
			tempHourPrice = val.HourPrice
			price += val.Price
		}
		allHourPeriodList[date] = HourPeriodList{StartDate: tempStartDate, EndDate: tempEndDate, HourPrice: tempHourPrice, Price: tempPrice, Duration: tempDuration}
	}
	//fmt.Println(price, "金额", allHourPeriodList, "..----,,,")

	price, _ = decimal.NewFromFloat(price).RoundFloor(2).Float64()
	return price, allHourPeriodList
}

//func (l *Charging) sortPeriods(currentDate time.Time) *PeriodChild {
//	// 节假日
//	if l.periods.Holiday != nil {
//		for index, val := range l.periods.Holiday.Date {
//			if ContainsSliceString(val, currentDate.Format("15:04")) {
//				return &PeriodChild{Hour: l.periods.Holiday.Hour[index], HourPeak: l.periods.Holiday.HourPeak[index]}
//			}
//		}
//	}
//	// 星期
//	if l.periods.Week != nil {
//		startDay := int64(currentDate.Weekday())
//		if startDay == 0 {
//			startDay = 7
//		}
//		if ContainsSliceInt64(l.periods.Week.Week, startDay) {
//			return &PeriodChild{Hour: l.periods.Week.Hour, HourPeak: l.periods.Week.HourPeak}
//		}
//	}
//	// 小时
//	if l.periods.Hour != nil {
//		return &PeriodChild{Hour: l.periods.Hour.Hour, HourPeak: l.periods.Hour.HourPeak}
//	}
//	return nil
//}

// 金额算出会有误差，具体以结束计算为准(分钟)
//func (l *Charging) moneyTransferMinuteTime(money float64, startDate string) (string, string, []*PeriodList) {
//	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
//	var lastEndTime time.Time
//	totalMoney := money
//	periodList := make([]*PeriodList, 0)
//	lastEndTime = currentTime
//	// 将 map 的键放入切片
//	for {
//		if totalMoney <= 0 {
//			//fmt.Println("跳出循环")
//			break
//		}
//		periods := l.sortPeriods(lastEndTime)
//
//		var periodEndTime time.Time
//		var duration time.Duration
//		for index, period := range periods.Hour {
//			var hourTotalMoney float64
//			endStr := fmt.Sprintf("%v:00", period.End)
//			if period.End < 10 {
//				endStr = fmt.Sprintf("0%v:00", period.End)
//			}
//			startStr := fmt.Sprintf("%v:00", period.Start)
//			if period.Start < 10 {
//				startStr = fmt.Sprintf("0%v:00", period.Start)
//			}
//			var tempPrice float64
//			tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
//			if l.member == 1 {
//				tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
//			}
//			//for _, index := range periods {
//			//	period := l.periods[int64(index)]
//			//}
//			//for id, period := range l.periods {
//			//	fmt.Println("~~~~~~~~~~~~~~", id, period.Start, period.End)
//			var minutePrice float64
//			var totalMinute int64
//			if totalMoney <= 0 {
//				continue
//			}
//			if currentTime.Format("15:04") == endStr {
//				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
//				continue
//			}
//			minutePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
//			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//
//			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
//			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
//			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
//			// 判断 如果小于1分钟价格则不取整
//			if subTotalMinute < 0.001 {
//				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
//			}
//
//			// 跨日
//			if startStr > endStr {
//				tempCurrentTime := currentTime
//				//fmt.Println("===========1111", currentTime.Format(time.DateTime), currentTime.Format("15:04"), period.End, period.Start)
//				if currentTime.Format("15:04") >= startStr {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), endStr), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				if currentTime.Format("15:04") < endStr {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				if currentTime.Format("15:04") < endStr && endStr == "24:00" {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
//				if totalMinute < durationMinute {
//					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//					lastEndTime = currentTime
//					hourTotalMoney += totalMoney
//					totalMoney = 0
//				} else {
//					currentTime = periodEndTime
//					lastEndTime = currentTime
//					surplusMoney, _ := decimal.NewFromInt(durationMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
//					totalMoney -= surplusMoney
//					hourTotalMoney += surplusMoney
//				}
//				Tduration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
//					Index:          index,
//					Duration:       Tduration,
//					Start:          startStr,
//					End:            endStr,
//					Price:          tempPrice,
//					TotalPrice:     totalMoney,
//					HourTotalMoney: hourTotalMoney,
//				})
//				continue
//			}
//			// 当日
//			if startStr < endStr && currentTime.Format("15:04") >= startStr && currentTime.Format("15:04") < endStr {
//				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
//				if endStr == "24:00" {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(time.Second*86400).Format("2006-01-02")), time.Local)
//				}
//				duration = periodEndTime.Sub(currentTime)
//
//				durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
//				tempCurrentTime := currentTime
//				if totalMinute < durationMinute {
//					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
//					lastEndTime = currentTime
//					hourTotalMoney += totalMoney
//					totalMoney = 0
//				} else {
//					currentTime = periodEndTime
//					lastEndTime = currentTime
//					surplusMoney, _ := decimal.NewFromInt(durationMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
//					totalMoney -= surplusMoney
//					hourTotalMoney += surplusMoney
//				}
//				fmt.Println(lastEndTime.Format(time.DateTime))
//				Tduration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
//					Index:          index,
//					Duration:       Tduration,
//					Start:          startStr,
//					End:            endStr,
//					Price:          tempPrice,
//					TotalPrice:     totalMoney,
//					HourTotalMoney: hourTotalMoney,
//				})
//				//fmt.Println("当日")
//			}
//		}
//	}
//	//tempStartDate,_:= time.ParseInLocation("2006-01-02 15:04",startDate,time.Local)
//	//for {
//	//	if tempStartDate.Unix()>=lastEndTime.Unix() {
//	//		break
//	//	}
//	//
//	//
//	//}
//	hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
//
//	fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime), hourPeriodLists)
//	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
//}
//
//// 获取金额转时长，每小时价格及时段
//func (l *Charging) moneyTransferHourPeriodList(periodList []*PeriodList) map[string]HourPeriodList {
//	hourPeriodLists := make(map[string]HourPeriodList, 0)
//	for _, value := range periodList {
//		tempTransferStartTime, _ := time.ParseInLocation("2006-01-02 15:04", value.StartTime, time.Local)
//		transferEndTime, _ := time.ParseInLocation("2006-01-02 15:04", value.EndTime, time.Local)
//		var transferPrice float64
//		transferPrice = value.HourTotalMoney
//		for {
//			if tempTransferStartTime.Unix() >= transferEndTime.Unix() {
//				break
//			}
//			if tempTransferStartTime.Format("04") != "00" {
//				tempSplitTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v:00", tempTransferStartTime.Add(60*time.Minute).Format("2006-01-02 15")), time.Local)
//				totalMaxMinute := decimal.NewFromFloat(tempSplitTime.Sub(tempTransferStartTime).Minutes()).IntPart()
//				if totalMaxMinute == value.Duration {
//					hourPeriodLists[tempTransferStartTime.Format("2006-01-02 15")] = HourPeriodList{
//						StartDate: tempTransferStartTime.Format("2006-01-02 15:04"),
//						EndDate:   tempTransferStartTime.Add(time.Duration(value.Duration) * time.Minute).Format("2006-01-02 15:04"),
//						Duration:  value.Duration,
//						HourPrice: transferPrice,
//						Price:     value.Price,
//					}
//					tempTransferStartTime = tempTransferStartTime.Add(time.Duration(value.Duration) * time.Minute)
//					continue
//				}
//				if totalMaxMinute < 60 {
//					minutePrice, _ := decimal.NewFromFloat(value.Price).Div(decimal.NewFromInt(60)).Float64()
//					splitMoney, _ := decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(totalMaxMinute)).Float64()
//					hourPeriodLists[tempTransferStartTime.Format("2006-01-02 15")] = HourPeriodList{
//						StartDate: tempTransferStartTime.Format("2006-01-02 15:04"),
//						EndDate:   tempTransferStartTime.Add(time.Duration(totalMaxMinute) * time.Minute).Format("2006-01-02 15:04"),
//						Duration:  totalMaxMinute,
//						HourPrice: splitMoney,
//						Price:     value.Price,
//					}
//					transferPrice -= splitMoney
//					tempTransferStartTime = tempTransferStartTime.Add(time.Duration(totalMaxMinute) * time.Minute)
//					continue
//				}
//
//			}
//			nextTransferStartTime := tempTransferStartTime.Add(60 * time.Minute)
//			if nextTransferStartTime.Unix() > transferEndTime.Unix() {
//				tempMinute := decimal.NewFromFloat(transferEndTime.Sub(tempTransferStartTime).Minutes()).IntPart()
//				hourPeriodLists[tempTransferStartTime.Format("2006-01-02 15")] = HourPeriodList{
//					StartDate: tempTransferStartTime.Format("2006-01-02 15:04"),
//					EndDate:   transferEndTime.Format("2006-01-02 15:04"),
//					Duration:  tempMinute,
//					HourPrice: transferPrice,
//					Price:     value.Price,
//				}
//				//fmt.Println("==", tempTransferStartTime.Format(time.DateTime), transferEndTime.Format(time.DateTime), transferPrice)
//				transferPrice = 0
//			} else {
//				//fmt.Println(tempTransferStartTime.Format(time.DateTime), nextTransferStartTime.Format(time.DateTime), transferPrice, value.Price)
//				tempMinute := decimal.NewFromFloat(nextTransferStartTime.Sub(tempTransferStartTime).Minutes()).IntPart()
//				hourPeriodLists[tempTransferStartTime.Format("2006-01-02 15")] = HourPeriodList{
//					StartDate: tempTransferStartTime.Format("2006-01-02 15:04"),
//					EndDate:   nextTransferStartTime.Format("2006-01-02 15:04"),
//					Duration:  tempMinute,
//					HourPrice: value.Price,
//					Price:     value.Price,
//				}
//				transferPrice -= value.Price
//			}
//			tempTransferStartTime = nextTransferStartTime
//		}
//	}
//	return hourPeriodLists
//}
//
//// 金额算出会有误差，具体以结束计算为准(半小时)[不足半小时或跨时段，剩余金额则按分钟计算]
//func (l *Charging) moneyTransferHalfHourTime(money float64, startDate string) (string, string, []*PeriodList) {
//	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
//	var lastEndTime time.Time
//	periodList := make([]*PeriodList, 0)
//	//periods := l.sortPeriods()
//	totalMoney := money
//	lastEndTime = currentTime
//	for {
//		if totalMoney <= 0 {
//			break
//		}
//		periods := l.sortPeriods(lastEndTime)
//		var periodEndTime time.Time
//		var duration time.Duration
//
//		for index, period := range periods.Hour {
//			var hourTotalMoney float64
//			endStr := fmt.Sprintf("%v:00", period.End)
//			if period.End < 10 {
//				endStr = fmt.Sprintf("0%v:00", period.End)
//			}
//			startStr := fmt.Sprintf("%v:00", period.Start)
//			if period.Start < 10 {
//				startStr = fmt.Sprintf("0%v:00", period.Start)
//			}
//			var tempPrice float64
//			tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
//			if l.member == 1 {
//				tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
//			}
//			var minutePrice float64
//			var totalMinute int64
//			if totalMoney <= 0 {
//				continue
//			}
//			if currentTime.Format("15:04") == endStr {
//				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
//				continue
//			}
//			tempCurrentTime := currentTime
//			minutePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
//			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//			halfHourPrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
//
//			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
//			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
//			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
//			// 判断 如果小于1分钟价格则不取整
//			if subTotalMinute < 0.001 {
//				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
//			}
//
//			// 跨日
//			if startStr > endStr {
//				if currentTime.Format("15:04") >= startStr {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), endStr), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				if currentTime.Format("15:04") < endStr {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				if currentTime.Format("15:04") < endStr && endStr == "24:00" {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
//				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
//					periodEndTime = tempTotalEndTime
//					duration = periodEndTime.Sub(currentTime)
//				}
//				lastEndTime = periodEndTime
//				currentTime = periodEndTime
//			}
//			// 当日
//			if startStr < endStr && currentTime.Format("15:04") >= startStr && currentTime.Format("15:04") < endStr {
//				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
//				if endStr == "24:00" {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(time.Second*86400).Format("2006-01-02")), time.Local)
//				}
//				duration = periodEndTime.Sub(currentTime)
//				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
//				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
//					periodEndTime = tempTotalEndTime
//					duration = periodEndTime.Sub(currentTime)
//				}
//				lastEndTime = periodEndTime
//				currentTime = periodEndTime
//			}
//			durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
//			if durationMinute >= 1 {
//				haltHour := decimal.NewFromInt(durationMinute).Div(decimal.NewFromInt(30)).Ceil().IntPart()
//				maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(30)).IntPart()
//				subMinute := decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(durationMinute)).IntPart()
//				surplusMinute := decimal.NewFromInt(30).Sub(decimal.NewFromInt(subMinute)).IntPart()
//				surplusMoney, _ := decimal.NewFromInt(haltHour - 1).Mul(decimal.NewFromFloat(halfHourPrice)).Float64()
//				surplusMinuteMoney, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
//				surplusMoney, _ = decimal.NewFromFloat(surplusMoney).Add(decimal.NewFromFloat(surplusMinuteMoney)).Float64()
//				totalMoney -= surplusMoney
//				hourTotalMoney += surplusMoney
//				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
//				if totalMoney < 0.001 {
//					hourTotalMoney += totalMoney
//					totalMoney = 0
//				}
//
//				Tduration := decimal.NewFromInt(durationMinute).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
//					Index:          index,
//					Duration:       Tduration,
//					Start:          startStr,
//					End:            endStr,
//					Price:          tempPrice,
//					TotalPrice:     totalMoney,
//					HourTotalMoney: hourTotalMoney,
//				})
//			}
//		}
//	}
//
//	hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
//	fmt.Println(hourPeriodLists)
//	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
//}
//
//// 金额算出会有误差，具体以结束计算为准(一小时)[不足一小时或跨时段，剩余金额则按分钟计算]
//func (l *Charging) moneyTransferHourTime(money float64, startDate string) (string, string, []*PeriodList) {
//	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
//	var lastEndTime time.Time
//	periodList := make([]*PeriodList, 0)
//	//periods := l.sortPeriods()
//	totalMoney := money
//	lastEndTime = currentTime
//	for {
//		if totalMoney <= 0 {
//			break
//		}
//		periods := l.sortPeriods(lastEndTime)
//
//		var periodEndTime time.Time
//		var duration time.Duration
//		for index, period := range periods.Hour {
//			var hourTotalMoney float64
//			endStr := fmt.Sprintf("%v:00", period.End)
//			if period.End < 10 {
//				endStr = fmt.Sprintf("0%v:00", period.End)
//			}
//			startStr := fmt.Sprintf("%v:00", period.Start)
//			if period.Start < 10 {
//				startStr = fmt.Sprintf("0%v:00", period.Start)
//			}
//			var tempPrice float64
//			tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
//			if l.member == 1 {
//				tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
//			}
//
//			//for _, index := range periods {
//			//	period := l.periods[int64(index)]
//			var minutePrice float64
//			var totalMinute int64
//			if totalMoney <= 0 {
//				continue
//			}
//			if currentTime.Format("15:04") == endStr {
//				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
//				continue
//			}
//
//			tempCurrentTime := currentTime
//			minutePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
//			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//
//			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
//			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
//			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
//			// 判断 如果小于1分钟价格则不取整
//			if subTotalMinute < 0.001 {
//				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
//			}
//
//			// 跨日
//			if startStr > endStr {
//				if currentTime.Format("15:04") >= startStr {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), endStr), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				if currentTime.Format("15:04") < endStr {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				if currentTime.Format("15:04") < endStr && endStr == "24:00" {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//
//				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
//				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
//					periodEndTime = tempTotalEndTime
//					duration = periodEndTime.Sub(currentTime)
//				}
//				lastEndTime = periodEndTime
//				currentTime = periodEndTime
//			}
//			// 当日
//			if startStr < endStr && currentTime.Format("15:04") >= startStr && currentTime.Format("15:04") < endStr {
//				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
//				if endStr == "24:00" {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(time.Second*86400).Format("2006-01-02")), time.Local)
//				}
//				duration = periodEndTime.Sub(currentTime)
//				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
//				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
//					periodEndTime = tempTotalEndTime
//					duration = periodEndTime.Sub(currentTime)
//				}
//				lastEndTime = periodEndTime
//				currentTime = periodEndTime
//			}
//			durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
//			if durationMinute >= 1 {
//				haltHour := decimal.NewFromInt(durationMinute).Div(decimal.NewFromInt(60)).Ceil().IntPart()
//				maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(60)).IntPart()
//				subMinute := decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(durationMinute)).IntPart()
//				surplusMinute := decimal.NewFromInt(60).Sub(decimal.NewFromInt(subMinute)).IntPart()
//				surplusMoney, _ := decimal.NewFromInt(haltHour - 1).Mul(decimal.NewFromFloat(tempPrice)).Float64()
//				surplusMinuteMoney, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
//				surplusMoney, _ = decimal.NewFromFloat(surplusMoney).Add(decimal.NewFromFloat(surplusMinuteMoney)).Float64()
//				totalMoney -= surplusMoney
//				hourTotalMoney += surplusMoney
//				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
//				if totalMoney < 0.001 {
//					hourTotalMoney += totalMoney
//					totalMoney = 0
//				}
//
//				Tduration := decimal.NewFromInt(durationMinute).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
//					Index:          index,
//					Duration:       Tduration,
//					Start:          startStr,
//					End:            endStr,
//					Price:          tempPrice,
//					TotalPrice:     totalMoney,
//					HourTotalMoney: hourTotalMoney,
//				})
//			}
//		}
//	}
//	hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
//	fmt.Println(hourPeriodLists)
//	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
//}
//
//// 金额算出会有误差，具体以结束计算为准(半小时\一小时)[不足半小时\一小时或跨时段，剩余金额则按分钟计算]
//func (l *Charging) moneyTransferHalfOrHourAndMinuteTime(money float64, startDate string) (string, string, []*PeriodList) {
//	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
//	var lastEndTime time.Time
//	periodList := make([]*PeriodList, 0)
//	//periods := l.sortPeriods()
//	totalMoney := money
//	lastEndTime = currentTime
//	for {
//		if totalMoney <= 0 {
//			break
//		}
//		periods := l.sortPeriods(lastEndTime)
//
//		var periodEndTime time.Time
//		var duration time.Duration
//
//		for index, period := range periods.Hour {
//			var hourTotalMoney float64
//			endStr := fmt.Sprintf("%v:00", period.End)
//			if period.End < 10 {
//				endStr = fmt.Sprintf("0%v:00", period.End)
//			}
//			startStr := fmt.Sprintf("%v:00", period.Start)
//			if period.Start < 10 {
//				startStr = fmt.Sprintf("0%v:00", period.Start)
//			}
//			var tempPrice float64
//			tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
//			if l.member == 1 {
//				tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
//			}
//			//for _, index := range periods {
//			//	period := l.periods[int64(index)]
//			var minutePrice float64
//			var totalMinute int64
//			var halfOrHourPrice float64
//			var minute int64
//			if totalMoney <= 0 {
//				continue
//			}
//			//halfHourPrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
//			if currentTime.Format("15:04") == endStr {
//				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
//				continue
//			}
//			tempCurrentTime := currentTime
//			minutePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
//			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//
//			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
//			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
//			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
//			// 判断 如果小于1分钟价格则不取整
//			if subTotalMinute < 0.001 {
//				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
//			}
//
//			// 不足半小时计费按小时计费，超过半小时按分钟计费
//			if ContainsSliceInt64([]int64{6, 7}, l.chargingMode) {
//				halfOrHourPrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
//				minute = 30
//			}
//			// 不足一小时计费按小时计费，超过一小时按分钟计费
//			if ContainsSliceInt64([]int64{8, 9}, l.chargingMode) {
//				halfOrHourPrice = tempPrice
//				minute = 60
//			}
//
//			// 跨日
//			if startStr > endStr {
//				if currentTime.Format("15:04") >= startStr {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), endStr), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				if currentTime.Format("15:04") < endStr {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				if currentTime.Format("15:04") < endStr && endStr == "24:00" {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
//				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
//					periodEndTime = tempTotalEndTime
//					duration = periodEndTime.Sub(currentTime)
//				}
//				lastEndTime = periodEndTime
//				currentTime = periodEndTime
//			}
//			// 当日
//			if startStr < endStr && currentTime.Format("15:04") >= startStr && currentTime.Format("15:04") < endStr {
//				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
//				if endStr == "24:00" {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(time.Second*86400).Format("2006-01-02")), time.Local)
//				}
//				duration = periodEndTime.Sub(currentTime)
//				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
//				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
//					periodEndTime = tempTotalEndTime
//					duration = periodEndTime.Sub(currentTime)
//				}
//				lastEndTime = periodEndTime
//				currentTime = periodEndTime
//			}
//			durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
//			if durationMinute >= 1 {
//				haltHour := decimal.NewFromInt(durationMinute).Div(decimal.NewFromInt(minute)).Ceil().IntPart()
//				maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(minute)).IntPart()
//				subMinute := decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(durationMinute)).IntPart()
//				surplusMinute := decimal.NewFromInt(minute).Sub(decimal.NewFromInt(subMinute)).IntPart()
//				surplusMoney, _ := decimal.NewFromInt(haltHour - 1).Mul(decimal.NewFromFloat(halfOrHourPrice)).Float64()
//				surplusMinuteMoney, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
//				surplusMoney, _ = decimal.NewFromFloat(surplusMoney).Add(decimal.NewFromFloat(surplusMinuteMoney)).Float64()
//				totalMoney -= surplusMoney
//				hourTotalMoney += surplusMoney
//				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
//				if totalMoney < 0.001 {
//					hourTotalMoney += totalMoney
//					totalMoney = 0
//				}
//
//				Tduration := decimal.NewFromInt(durationMinute).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
//					Index:          index,
//					Duration:       Tduration,
//					Start:          startStr,
//					End:            endStr,
//					Price:          tempPrice,
//					TotalPrice:     totalMoney,
//					HourTotalMoney: hourTotalMoney,
//				})
//			}
//		}
//	}
//	hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
//	fmt.Println(hourPeriodLists)
//	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
//}
//
//// 金额算出会有误差，具体以结束计算为准(以10分钟作为一个收费周期，第1分钟后开始计费)[不足周期计费或跨时段，剩余金额则按分钟计算]
//func (l *Charging) moneyTransferCycleAndMinuteTime(money float64, startDate string) (string, string, []*PeriodList) {
//	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
//	var lastEndTime time.Time
//	periodList := make([]*PeriodList, 0)
//	//periods := l.sortPeriods()
//	totalMoney := money
//	lastEndTime = currentTime
//	for {
//		if totalMoney <= 0 {
//			break
//		}
//		periods := l.sortPeriods(lastEndTime)
//		var periodEndTime time.Time
//		var duration time.Duration
//
//		for index, period := range periods.Hour {
//			var hourTotalMoney float64
//			endStr := fmt.Sprintf("%v:00", period.End)
//			if period.End < 10 {
//				endStr = fmt.Sprintf("0%v:00", period.End)
//			}
//			startStr := fmt.Sprintf("%v:00", period.Start)
//			if period.Start < 10 {
//				startStr = fmt.Sprintf("0%v:00", period.Start)
//			}
//			var tempPrice float64
//			tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
//			if l.member == 1 {
//				tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
//			}
//			//for _, index := range periods {
//			//	period := l.periods[int64(index)]
//			var minutePrice float64
//			var totalMinute int64
//			//var halfOrHourPrice float64
//			//var minute int64
//			var cyclePrice float64
//			if totalMoney <= 0 {
//				continue
//			}
//			//halfHourPrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
//			if currentTime.Format("15:04") == endStr {
//				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
//				continue
//			}
//			tempCurrentTime := currentTime
//			minutePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
//			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
//
//			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
//			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
//			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
//			// 判断 如果小于1分钟价格则不取整
//			if subTotalMinute < 0.001 {
//				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
//			}
//
//			if l.cycleMinute == 30 {
//				cyclePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
//			}
//			if l.cycleMinute == 60 {
//				cyclePrice = tempPrice
//			}
//			if l.cycleMinute != 30 && l.cycleMinute != 60 {
//				cyclePrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
//			}
//
//			// 跨日
//			if startStr > endStr {
//				if currentTime.Format("15:04") >= startStr {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), endStr), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				if currentTime.Format("15:04") < endStr {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				if currentTime.Format("15:04") < endStr && endStr == "24:00" {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
//					duration = periodEndTime.Sub(currentTime)
//				}
//				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
//				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
//					periodEndTime = tempTotalEndTime
//					duration = periodEndTime.Sub(currentTime)
//				}
//				lastEndTime = periodEndTime
//				currentTime = periodEndTime
//			}
//			// 当日
//			if startStr < endStr && currentTime.Format("15:04") >= startStr && currentTime.Format("15:04") < endStr {
//				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
//				if endStr == "24:00" {
//					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(time.Second*86400).Format("2006-01-02")), time.Local)
//				}
//				duration = periodEndTime.Sub(currentTime)
//				tempTotalEndTime := currentTime.Add(time.Duration(totalMinute) * time.Minute)
//				if tempTotalEndTime.Unix() < periodEndTime.Unix() {
//					periodEndTime = tempTotalEndTime
//					duration = periodEndTime.Sub(currentTime)
//				}
//				lastEndTime = periodEndTime
//				currentTime = periodEndTime
//			}
//			durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
//			if durationMinute >= 1 {
//				haltHour := decimal.NewFromInt(durationMinute).Div(decimal.NewFromInt(l.cycleMinute)).Ceil().IntPart()
//				maxMinute := decimal.NewFromInt(haltHour).Mul(decimal.NewFromInt(l.cycleMinute)).IntPart()
//				subMinute := decimal.NewFromInt(maxMinute).Sub(decimal.NewFromInt(durationMinute)).IntPart()
//				surplusMinute := decimal.NewFromInt(l.cycleMinute).Sub(decimal.NewFromInt(subMinute)).IntPart()
//				surplusMoney, _ := decimal.NewFromInt(haltHour - 1).Mul(decimal.NewFromFloat(cyclePrice)).Float64()
//				surplusMinuteMoney, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
//				surplusMoney, _ = decimal.NewFromFloat(surplusMoney).Add(decimal.NewFromFloat(surplusMinuteMoney)).Float64()
//				totalMoney -= surplusMoney
//				hourTotalMoney += surplusMoney
//				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
//				if totalMoney < 0.001 {
//					hourTotalMoney += totalMoney
//					totalMoney = 0
//				}
//
//				Tduration := decimal.NewFromInt(durationMinute).IntPart()
//				periodList = append(periodList, &PeriodList{
//					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
//					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
//					Index:          index,
//					Duration:       Tduration,
//					Start:          startStr,
//					End:            endStr,
//					Price:          tempPrice,
//					TotalPrice:     totalMoney,
//					HourTotalMoney: hourTotalMoney,
//				})
//			}
//		}
//	}
//	hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
//	fmt.Println(hourPeriodLists)
//	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
//}
//
//// 金额转换成开始及结束时间
//func (l *Charging) MoneyTransfer(money float64, startDateParam string) (string, string, []*PeriodList) {
//	var startData string
//	var endData string
//	periodList := make([]*PeriodList, 0)
//
//	//计费模式 1按分钟计费 2按半小时计费 3按半小时计费(跨时段) 4按小时计费 5按小时计费(跨时段)
//	//6按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费  7按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费(跨时段)
//	//8按小时计费开台不足1小时按小时计费，超过1小时按分钟计费   9按小时计费开台不足1小时按小时计费，超过1小时按分钟计费(跨时段)
//	//10自定义计费 11自定义计费(跨时段)
//
//	// 分钟计费
//	if l.chargingMode == 1 {
//		startData, endData, periodList = l.moneyTransferMinuteTime(money, startDateParam)
//	}
//	// 半小时计费
//	if ContainsSliceInt64([]int64{2, 3}, l.chargingMode) {
//		startData, endData, periodList = l.moneyTransferHalfHourTime(money, startDateParam)
//		fmt.Println(startData, endData, periodList)
//	}
//	// 小时计费
//	if ContainsSliceInt64([]int64{4, 5}, l.chargingMode) {
//		startData, endData, periodList = l.moneyTransferHourTime(money, startDateParam)
//		fmt.Println(startData, endData, periodList)
//	}
//	// 按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
//	if ContainsSliceInt64([]int64{6, 7}, l.chargingMode) {
//		startData, endData, periodList = l.moneyTransferHalfOrHourAndMinuteTime(money, startDateParam)
//		fmt.Println(startData, endData, periodList)
//	}
//	// 按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
//	if ContainsSliceInt64([]int64{8, 9}, l.chargingMode) {
//		startData, endData, periodList = l.moneyTransferHalfOrHourAndMinuteTime(money, startDateParam)
//		fmt.Println(startData, endData, periodList)
//	}
//	// 自定义计费
//	if ContainsSliceInt64([]int64{10, 11}, l.chargingMode) {
//		startData, endData, periodList = l.moneyTransferCycleAndMinuteTime(money, startDateParam)
//		fmt.Println(startData, endData, periodList)
//	}
//	return startData, endData, periodList
//}
