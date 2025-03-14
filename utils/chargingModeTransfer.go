package utils

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"time"
)

type ChargingTransfer struct {
	//periods map[int64]ChargePeriod
	periods      ChargePeriodAssembly
	member       int64 // 是否是会员 0否  1是
	chargingMode int64 // 计费方式
	cycleMinute  int64 // 计费周期
}

type PeriodList struct {
	StartTime      string  `json:"start_time"`
	EndTime        string  `json:"end_time"`
	Index          int     `json:"index"`
	Duration       int64   `json:"duration"`
	Start          string  `json:"start"`
	End            string  `json:"end"`
	Price          float64 `json:"price"`
	TotalPrice     float64 `json:"total_price"`
	HourTotalMoney float64 `json:"hour_total_price"`
}
type PeriodChild struct {
	Hour     []CPHour
	HourPeak []CPHour
}

func NewChargeModeTransfer(periods ChargePeriodAssembly, member int64, chargingMode int64, cycleMinute int64) *ChargingTransfer {
	return &ChargingTransfer{periods: periods, member: member, chargingMode: chargingMode, cycleMinute: cycleMinute}
}

//func (l *ChargingTransfer) computePeak(hourPeriodList map[string]HourPeriodList) (float64, map[string]HourPeriodList) {
//	dateSlice := make([]string, 0)
//	dateSlicePrice := make(map[string]float64, 0)
//	//dateSliceDate := make(map[string]int64, 0)
//	for date, val := range hourPeriodList {
//		dateSlice = append(dateSlice, date)
//		dateSlicePrice[date] = val.Price
//	}
//	sort.Strings(dateSlice)
//	dateAccrual := make(map[string]map[int]float64, 0)
//	for _, date := range dateSlice {
//		isOk := 0
//		currentTime, _ := time.ParseInLocation("2006-01-02 15", date, time.Local)
//		period := int64(currentTime.Hour())
//		dateKey := currentTime.Format("2006-01-02")
//		dateTimeKey := currentTime.Format("2006-01-02 15")
//		if _, ok := dateAccrual[dateKey]; !ok {
//			dateAccrual[dateKey] = make(map[int]float64)
//		}
//		// 节假日
//		if l.periods.Holiday != nil && isOk == 0 {
//			startDate := currentTime.Format("01-02")
//			for _, value := range l.periods.Holiday.Date {
//				if ContainsSliceString(value, startDate) {
//					for index, val := range value {
//						if val == startDate {
//							if len(l.periods.Holiday.HourPeak[index]) >= 1 {
//								for _, v := range l.periods.Holiday.HourPeak[index] {
//									var peakPrice float64
//									peakPrice, _ = strconv.ParseFloat(v.IdlePrice, 10)
//									if l.member == 1 {
//										peakPrice, _ = strconv.ParseFloat(v.MemberPrice, 10)
//									}
//									// 跨日
//									if v.Start > v.End {
//										if v.Start <= period && v.End < period || v.Start > period && v.End > period {
//											//dateAccrual[dateKey][index] += dateSlicePrice[date]
//											//if _, ok := dateDeduct[dateKey][index]; !ok {
//											//	dateDeduct[dateKey][index] = peakPrice
//											//}
//											if _, ok := dateAccrual[dateKey][index]; !ok {
//												if dateSlicePrice[dateTimeKey] >= peakPrice {
//													dateSlicePrice[dateTimeKey] = peakPrice
//													dateAccrual[dateKey][index] += peakPrice
//												} else {
//													dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
//												}
//											} else if daVal, ok := dateAccrual[dateKey][index]; ok {
//												if daVal >= peakPrice {
//													dateSlicePrice[dateTimeKey] = 0
//												} else {
//													surplusPeak := (peakPrice - daVal)
//													dateAccrual[dateKey][index] += surplusPeak
//													if dateSlicePrice[date] > surplusPeak {
//														dateSlicePrice[date] = surplusPeak
//													}
//												}
//											}
//											//if dateDeduct[dateKey][index] < dateAccrual[dateKey][index] {
//											//	dateSlicePrice[dateTimeKey] = 0
//											//	//fmt.Println(dateTimeKey, "至为0")
//											//}
//											isOk = 1
//										}
//									}
//									// 当日
//									if v.Start < v.End && v.Start <= period && v.End > period {
//										//dateAccrual[dateKey][index] += dateSlicePrice[date]
//										//if _, ok := dateDeduct[dateKey][index]; !ok {
//										//	dateDeduct[dateKey][index] = peakPrice
//										//}
//										if _, ok := dateAccrual[dateKey][index]; !ok {
//											if dateSlicePrice[dateTimeKey] >= peakPrice {
//												dateSlicePrice[dateTimeKey] = peakPrice
//												dateAccrual[dateKey][index] += peakPrice
//											} else {
//												dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
//											}
//										} else if daVal, ok := dateAccrual[dateKey][index]; ok {
//											if daVal >= peakPrice {
//												dateSlicePrice[dateTimeKey] = 0
//											} else {
//												surplusPeak := (peakPrice - daVal)
//												dateAccrual[dateKey][index] += surplusPeak
//												if dateSlicePrice[date] > surplusPeak {
//													dateSlicePrice[date] = surplusPeak
//												}
//											}
//										}
//										//if dateDeduct[dateKey][index] < dateAccrual[dateKey][index] {
//										//	//fmt.Println(dateTimeKey, "至为0")
//										//	dateSlicePrice[dateTimeKey] = 0
//										//}
//										isOk = 1
//									}
//								}
//							}
//						}
//					}
//					continue
//				}
//			}
//		}
//		//星期
//		if l.periods.Week != nil && isOk == 0 {
//			startDay := int64(currentTime.Weekday())
//			if startDay == 0 {
//				startDay = 7
//			}
//			if ContainsSliceInt64(l.periods.Week.Week, startDay) {
//				if len(l.periods.Week.HourPeak) >= 1 {
//					for index, v := range l.periods.Week.HourPeak {
//						var peakPrice float64
//						peakPrice, _ = strconv.ParseFloat(v.IdlePrice, 10)
//						if l.member == 1 {
//							peakPrice, _ = strconv.ParseFloat(v.MemberPrice, 10)
//						}
//						// 跨日
//						if v.Start > v.End {
//							if v.Start <= period && v.End < period || v.Start > period && v.End > period {
//								//dateAccrual[dateKey][index] += dateSlicePrice[date]
//								//if _, ok := dateDeduct[dateKey][index]; !ok {
//								//	dateDeduct[dateKey][index] = peakPrice
//								//}
//
//								if _, ok := dateAccrual[dateKey][index]; !ok {
//									if dateSlicePrice[dateTimeKey] >= peakPrice {
//										dateSlicePrice[dateTimeKey] = peakPrice
//										dateAccrual[dateKey][index] += peakPrice
//									} else {
//										dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
//									}
//								} else if daVal, ok := dateAccrual[dateKey][index]; ok {
//									if daVal >= peakPrice {
//										dateSlicePrice[dateTimeKey] = 0
//									} else {
//										surplusPeak := (peakPrice - daVal)
//										dateAccrual[dateKey][index] += surplusPeak
//										if dateSlicePrice[date] > surplusPeak {
//											dateSlicePrice[date] = surplusPeak
//										}
//									}
//								}
//
//								isOk = 1
//							}
//						}
//						// 当日
//						if v.Start < v.End && v.Start <= period && v.End > period {
//							//dateAccrual[dateKey][index] += dateSlicePrice[date]
//							//if _, ok := dateDeduct[dateKey][index]; !ok {
//							//	dateDeduct[dateKey][index] = peakPrice
//							//}
//
//							if _, ok := dateAccrual[dateKey][index]; !ok {
//								if dateSlicePrice[dateTimeKey] >= peakPrice {
//									dateSlicePrice[dateTimeKey] = peakPrice
//									dateAccrual[dateKey][index] += peakPrice
//								} else {
//									dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
//								}
//							} else if daVal, ok := dateAccrual[dateKey][index]; ok {
//								if daVal >= peakPrice {
//									dateSlicePrice[dateTimeKey] = 0
//								} else {
//									surplusPeak := (peakPrice - daVal)
//									dateAccrual[dateKey][index] += surplusPeak
//									if dateSlicePrice[date] > surplusPeak {
//										dateSlicePrice[date] = surplusPeak
//									}
//								}
//							}
//							isOk = 1
//						}
//					}
//				}
//			}
//			//fmt.Println(tempStartTime.Format(time.DateTime), "开始时间存在星期")
//			//fmt.Println("星期在范围内", tempCurrentTime.Format(time.DateTime))
//		}
//		// 小时
//		if l.periods.Hour != nil && isOk == 0 {
//			startDay := int64(currentTime.Weekday())
//			if startDay == 0 {
//				startDay = 7
//			}
//			if len(l.periods.Hour.HourPeak) >= 1 {
//				for index, v := range l.periods.Hour.HourPeak {
//
//					var peakPrice float64
//					peakPrice, _ = strconv.ParseFloat(v.IdlePrice, 10)
//					if l.member == 1 {
//						peakPrice, _ = strconv.ParseFloat(v.MemberPrice, 10)
//					}
//					// 跨日
//					if v.Start > v.End {
//						if v.Start <= period && v.End < period || v.Start > period && v.End > period {
//							if _, ok := dateAccrual[dateKey][index]; !ok {
//								if dateSlicePrice[dateTimeKey] >= peakPrice {
//									dateSlicePrice[dateTimeKey] = peakPrice
//									dateAccrual[dateKey][index] += peakPrice
//								} else {
//									dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
//								}
//							} else if daVal, ok := dateAccrual[dateKey][index]; ok {
//								if daVal >= peakPrice {
//									dateSlicePrice[dateTimeKey] = 0
//								} else {
//									surplusPeak := (peakPrice - daVal)
//									dateAccrual[dateKey][index] += surplusPeak
//									if dateSlicePrice[date] > surplusPeak {
//										dateSlicePrice[date] = surplusPeak
//									}
//								}
//							}
//							//if dateDeduct[dateKey][index] < dateAccrual[dateKey][index] {
//							//	dateSlicePrice[dateTimeKey] = 0
//							//	//fmt.Println(dateTimeKey, "至为0")
//							//}
//						}
//					}
//					// 当日
//					if v.Start < v.End && v.Start <= period && v.End > period {
//						//dateAccrual[dateKey][index] += dateSlicePrice[date]
//						//if _, ok := dateDeduct[dateKey][index]; !ok {
//						//	dateDeduct[dateKey][index] = peakPrice
//						//}
//						if _, ok := dateAccrual[dateKey][index]; !ok {
//							if dateSlicePrice[dateTimeKey] >= peakPrice {
//								dateSlicePrice[dateTimeKey] = peakPrice
//								dateAccrual[dateKey][index] += peakPrice
//							} else {
//								dateAccrual[dateKey][index] += dateSlicePrice[dateTimeKey]
//							}
//						} else if daVal, ok := dateAccrual[dateKey][index]; ok {
//							if daVal >= peakPrice {
//								dateSlicePrice[dateTimeKey] = 0
//							} else {
//								surplusPeak := (peakPrice - daVal)
//								dateAccrual[dateKey][index] += surplusPeak
//								if dateSlicePrice[date] > surplusPeak {
//									dateSlicePrice[date] = surplusPeak
//								}
//							}
//						}
//						//if dateDeduct[dateKey][index] < dateAccrual[dateKey][index] {
//						//	dateSlicePrice[dateTimeKey] = 0
//						//	//fmt.Println(dateTimeKey, "至为0")
//						//}
//					}
//				}
//			}
//		}
//	}
//
//	for date, val := range hourPeriodList {
//		tempVal := val
//		if valPrice, ok := dateSlicePrice[date]; ok {
//			if valPrice != tempVal.Price {
//				if vm, oks := hourPeriodList[date]; oks {
//					vm.HourPeak = 1
//					vm.HourPeakPrice = valPrice
//					hourPeriodList[date] = vm
//				}
//			}
//		}
//	}
//	var price float64
//	for _, value := range dateSlicePrice {
//		tempValuePrice := value
//		price += tempValuePrice
//	}
//	//系统提示词：按什么方法什么方式做什么事情
//	return price, hourPeriodList
//}

func (l *ChargingTransfer) sortPeriods(currentDate time.Time) *PeriodChild {
	// 节假日
	if l.periods.Holiday != nil {
		for index, val := range l.periods.Holiday.Date {
			if ContainsSliceString(val, currentDate.Format("15:04")) {
				return &PeriodChild{Hour: l.periods.Holiday.Hour[index], HourPeak: l.periods.Holiday.HourPeak[index]}
			}
		}
	}
	// 星期
	if l.periods.Week != nil {
		startDay := int64(currentDate.Weekday())
		if startDay == 0 {
			startDay = 7
		}
		if ContainsSliceInt64(l.periods.Week.Week, startDay) {
			return &PeriodChild{Hour: l.periods.Week.Hour, HourPeak: l.periods.Week.HourPeak}
		}
	}
	// 小时
	if l.periods.Hour != nil {
		return &PeriodChild{Hour: l.periods.Hour.Hour, HourPeak: l.periods.Hour.HourPeak}
	}
	return nil
}

// 金额算出会有误差，具体以结束计算为准(分钟)
func (l *ChargingTransfer) moneyTransferMinuteTime(money float64, startDate string) (string, string, []*PeriodList) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	totalMoney := money
	periodList := make([]*PeriodList, 0)
	lastEndTime = currentTime
	// 将 map 的键放入切片
	for {
		if totalMoney <= 0 {
			//fmt.Println("跳出循环")
			break
		}
		periods := l.sortPeriods(lastEndTime)

		var periodEndTime time.Time
		var duration time.Duration
		for index, period := range periods.Hour {
			var hourTotalMoney float64
			endStr := fmt.Sprintf("%v:00", period.End)
			if period.End < 10 {
				endStr = fmt.Sprintf("0%v:00", period.End)
			}
			startStr := fmt.Sprintf("%v:00", period.Start)
			if period.Start < 10 {
				startStr = fmt.Sprintf("0%v:00", period.Start)
			}
			var tempPrice float64
			tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
			if l.member == 1 {
				tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
			}
			//for _, index := range periods {
			//	period := l.periods[int64(index)]
			//}
			//for id, period := range l.periods {
			//	fmt.Println("~~~~~~~~~~~~~~", id, period.Start, period.End)
			var minutePrice float64
			var totalMinute int64
			if totalMoney <= 0 {
				continue
			}
			if currentTime.Format("15:04") == endStr {
				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
				continue
			}
			minutePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()

			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
			// 判断 如果小于1分钟价格则不取整
			if subTotalMinute < 0.001 {
				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			}

			// 跨日
			if startStr > endStr {
				tempCurrentTime := currentTime
				//fmt.Println("===========1111", currentTime.Format(time.DateTime), currentTime.Format("15:04"), period.End, period.Start)
				if currentTime.Format("15:04") >= startStr {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), endStr), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < endStr {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < endStr && endStr == "24:00" {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
				if totalMinute < durationMinute {
					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
					lastEndTime = currentTime
					hourTotalMoney += totalMoney
					totalMoney = 0
				} else {
					currentTime = periodEndTime
					lastEndTime = currentTime
					surplusMoney, _ := decimal.NewFromInt(durationMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					totalMoney -= surplusMoney
					hourTotalMoney += surplusMoney
				}
				Tduration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
					Index:          index,
					Duration:       Tduration,
					Start:          startStr,
					End:            endStr,
					Price:          tempPrice,
					TotalPrice:     totalMoney,
					HourTotalMoney: hourTotalMoney,
				})
				continue
			}
			// 当日
			if startStr < endStr && currentTime.Format("15:04") >= startStr && currentTime.Format("15:04") < endStr {
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
				if endStr == "24:00" {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(time.Second*86400).Format("2006-01-02")), time.Local)
				}
				duration = periodEndTime.Sub(currentTime)

				durationMinute := decimal.NewFromFloat(duration.Minutes()).IntPart()
				tempCurrentTime := currentTime
				if totalMinute < durationMinute {
					currentTime = currentTime.Add(time.Duration(totalMinute) * time.Minute)
					lastEndTime = currentTime
					hourTotalMoney += totalMoney
					totalMoney = 0
				} else {
					currentTime = periodEndTime
					lastEndTime = currentTime
					surplusMoney, _ := decimal.NewFromInt(durationMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					totalMoney -= surplusMoney
					hourTotalMoney += surplusMoney
				}
				//fmt.Println(lastEndTime.Format(time.DateTime))
				Tduration := decimal.NewFromFloat(lastEndTime.Sub(tempCurrentTime).Minutes()).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
					Index:          index,
					Duration:       Tduration,
					Start:          startStr,
					End:            endStr,
					Price:          tempPrice,
					TotalPrice:     totalMoney,
					HourTotalMoney: hourTotalMoney,
				})
				//fmt.Println("当日")
			}
		}
	}
	//tempStartDate,_:= time.ParseInLocation("2006-01-02 15:04",startDate,time.Local)
	//for {
	//	if tempStartDate.Unix()>=lastEndTime.Unix() {
	//		break
	//	}
	//
	//
	//}

	//fmt.Println("结果：", startDate, lastEndTime.Format(time.DateTime), hourPeriodLists)
	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
}

// 获取金额转时长，每小时价格及时段
func (l *ChargingTransfer) moneyTransferHourPeriodList(periodList []*PeriodList) map[string]HourPeriodList {
	hourPeriodLists := make(map[string]HourPeriodList, 0)
	for _, value := range periodList {
		tempTransferStartTime, _ := time.ParseInLocation("2006-01-02 15:04", value.StartTime, time.Local)
		transferEndTime, _ := time.ParseInLocation("2006-01-02 15:04", value.EndTime, time.Local)
		var transferPrice float64
		transferPrice = value.HourTotalMoney
		for {
			if tempTransferStartTime.Unix() >= transferEndTime.Unix() {
				break
			}
			if tempTransferStartTime.Format("04") != "00" {
				tempSplitTime, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v:00", tempTransferStartTime.Add(60*time.Minute).Format("2006-01-02 15")), time.Local)
				totalMaxMinute := decimal.NewFromFloat(tempSplitTime.Sub(tempTransferStartTime).Minutes()).IntPart()
				if totalMaxMinute == value.Duration {
					hourPeriodLists[tempTransferStartTime.Format("2006-01-02 15")] = HourPeriodList{
						StartDate: tempTransferStartTime.Format("2006-01-02 15:04"),
						EndDate:   tempTransferStartTime.Add(time.Duration(value.Duration) * time.Minute).Format("2006-01-02 15:04"),
						Duration:  value.Duration,
						HourPrice: transferPrice,
						Price:     value.Price,
					}
					tempTransferStartTime = tempTransferStartTime.Add(time.Duration(value.Duration) * time.Minute)
					continue
				}
				if totalMaxMinute < 60 {
					minutePrice, _ := decimal.NewFromFloat(value.Price).Div(decimal.NewFromInt(60)).Float64()
					splitMoney, _ := decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(totalMaxMinute)).Float64()
					hourPeriodLists[tempTransferStartTime.Format("2006-01-02 15")] = HourPeriodList{
						StartDate: tempTransferStartTime.Format("2006-01-02 15:04"),
						EndDate:   tempTransferStartTime.Add(time.Duration(totalMaxMinute) * time.Minute).Format("2006-01-02 15:04"),
						Duration:  totalMaxMinute,
						HourPrice: splitMoney,
						Price:     value.Price,
					}
					transferPrice -= splitMoney
					tempTransferStartTime = tempTransferStartTime.Add(time.Duration(totalMaxMinute) * time.Minute)
					continue
				}

			}
			nextTransferStartTime := tempTransferStartTime.Add(60 * time.Minute)
			if nextTransferStartTime.Unix() > transferEndTime.Unix() {
				tempMinute := decimal.NewFromFloat(transferEndTime.Sub(tempTransferStartTime).Minutes()).IntPart()
				hourPeriodLists[tempTransferStartTime.Format("2006-01-02 15")] = HourPeriodList{
					StartDate: tempTransferStartTime.Format("2006-01-02 15:04"),
					EndDate:   transferEndTime.Format("2006-01-02 15:04"),
					Duration:  tempMinute,
					HourPrice: transferPrice,
					Price:     value.Price,
				}
				//fmt.Println("==", tempTransferStartTime.Format(time.DateTime), transferEndTime.Format(time.DateTime), transferPrice)
				transferPrice = 0
			} else {
				//fmt.Println(tempTransferStartTime.Format(time.DateTime), nextTransferStartTime.Format(time.DateTime), transferPrice, value.Price)
				tempMinute := decimal.NewFromFloat(nextTransferStartTime.Sub(tempTransferStartTime).Minutes()).IntPart()
				hourPeriodLists[tempTransferStartTime.Format("2006-01-02 15")] = HourPeriodList{
					StartDate: tempTransferStartTime.Format("2006-01-02 15:04"),
					EndDate:   nextTransferStartTime.Format("2006-01-02 15:04"),
					Duration:  tempMinute,
					HourPrice: value.Price,
					Price:     value.Price,
				}
				transferPrice -= value.Price
			}
			tempTransferStartTime = nextTransferStartTime
		}
	}
	return hourPeriodLists
}

// 金额算出会有误差，具体以结束计算为准(半小时)[不足半小时或跨时段，剩余金额则按分钟计算]
func (l *ChargingTransfer) moneyTransferHalfHourTime(money float64, startDate string) (string, string, []*PeriodList) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	periodList := make([]*PeriodList, 0)
	//periods := l.sortPeriods()
	totalMoney := money
	lastEndTime = currentTime
	for {
		if totalMoney <= 0 {
			break
		}
		periods := l.sortPeriods(lastEndTime)
		var periodEndTime time.Time
		var duration time.Duration

		for index, period := range periods.Hour {
			var hourTotalMoney float64
			endStr := fmt.Sprintf("%v:00", period.End)
			if period.End < 10 {
				endStr = fmt.Sprintf("0%v:00", period.End)
			}
			startStr := fmt.Sprintf("%v:00", period.Start)
			if period.Start < 10 {
				startStr = fmt.Sprintf("0%v:00", period.Start)
			}
			var tempPrice float64
			tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
			if l.member == 1 {
				tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
			}
			var minutePrice float64
			var totalMinute int64
			if totalMoney <= 0 {
				continue
			}
			if currentTime.Format("15:04") == endStr {
				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
				continue
			}
			tempCurrentTime := currentTime
			minutePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
			halfHourPrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()

			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
			// 判断 如果小于1分钟价格则不取整
			if subTotalMinute < 0.001 {
				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			}

			// 跨日
			if startStr > endStr {
				if currentTime.Format("15:04") >= startStr {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), endStr), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < endStr {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < endStr && endStr == "24:00" {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
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
			if startStr < endStr && currentTime.Format("15:04") >= startStr && currentTime.Format("15:04") < endStr {
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
				if endStr == "24:00" {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(time.Second*86400).Format("2006-01-02")), time.Local)
				}
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
				hourTotalMoney += surplusMoney
				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
				if totalMoney < 0.001 {
					hourTotalMoney += totalMoney
					totalMoney = 0
				}

				Tduration := decimal.NewFromInt(durationMinute).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
					Index:          index,
					Duration:       Tduration,
					Start:          startStr,
					End:            endStr,
					Price:          tempPrice,
					TotalPrice:     totalMoney,
					HourTotalMoney: hourTotalMoney,
				})
			}
		}
	}

	//hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
	//fmt.Println(hourPeriodLists)
	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
}

// 金额算出会有误差，具体以结束计算为准(一小时)[不足一小时或跨时段，剩余金额则按分钟计算]
func (l *ChargingTransfer) moneyTransferHourTime(money float64, startDate string) (string, string, []*PeriodList) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	periodList := make([]*PeriodList, 0)
	//periods := l.sortPeriods()
	totalMoney := money
	lastEndTime = currentTime
	for {
		if totalMoney <= 0 {
			break
		}
		periods := l.sortPeriods(lastEndTime)

		var periodEndTime time.Time
		var duration time.Duration
		for index, period := range periods.Hour {
			var hourTotalMoney float64
			endStr := fmt.Sprintf("%v:00", period.End)
			if period.End < 10 {
				endStr = fmt.Sprintf("0%v:00", period.End)
			}
			startStr := fmt.Sprintf("%v:00", period.Start)
			if period.Start < 10 {
				startStr = fmt.Sprintf("0%v:00", period.Start)
			}
			var tempPrice float64
			tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
			if l.member == 1 {
				tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
			}

			//for _, index := range periods {
			//	period := l.periods[int64(index)]
			var minutePrice float64
			var totalMinute int64
			if totalMoney <= 0 {
				continue
			}
			if currentTime.Format("15:04") == endStr {
				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
				continue
			}

			tempCurrentTime := currentTime
			minutePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()

			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
			// 判断 如果小于1分钟价格则不取整
			if subTotalMinute < 0.001 {
				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			}

			// 跨日
			if startStr > endStr {
				if currentTime.Format("15:04") >= startStr {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), endStr), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < endStr {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < endStr && endStr == "24:00" {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
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
			if startStr < endStr && currentTime.Format("15:04") >= startStr && currentTime.Format("15:04") < endStr {
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
				if endStr == "24:00" {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(time.Second*86400).Format("2006-01-02")), time.Local)
				}
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
				surplusMoney, _ := decimal.NewFromInt(haltHour - 1).Mul(decimal.NewFromFloat(tempPrice)).Float64()
				surplusMinuteMoney, _ := decimal.NewFromInt(surplusMinute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
				surplusMoney, _ = decimal.NewFromFloat(surplusMoney).Add(decimal.NewFromFloat(surplusMinuteMoney)).Float64()
				totalMoney -= surplusMoney
				hourTotalMoney += surplusMoney
				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
				if totalMoney < 0.001 {
					hourTotalMoney += totalMoney
					totalMoney = 0
				}

				Tduration := decimal.NewFromInt(durationMinute).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
					Index:          index,
					Duration:       Tduration,
					Start:          startStr,
					End:            endStr,
					Price:          tempPrice,
					TotalPrice:     totalMoney,
					HourTotalMoney: hourTotalMoney,
				})
			}
		}
	}
	//hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
	//fmt.Println(hourPeriodLists)
	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
}

// 金额算出会有误差，具体以结束计算为准(半小时\一小时)[不足半小时\一小时或跨时段，剩余金额则按分钟计算]
func (l *ChargingTransfer) moneyTransferHalfOrHourAndMinuteTime(money float64, startDate string) (string, string, []*PeriodList) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	periodList := make([]*PeriodList, 0)
	//periods := l.sortPeriods()
	totalMoney := money
	lastEndTime = currentTime
	for {
		if totalMoney <= 0 {
			break
		}
		periods := l.sortPeriods(lastEndTime)

		var periodEndTime time.Time
		var duration time.Duration

		for index, period := range periods.Hour {
			var hourTotalMoney float64
			endStr := fmt.Sprintf("%v:00", period.End)
			if period.End < 10 {
				endStr = fmt.Sprintf("0%v:00", period.End)
			}
			startStr := fmt.Sprintf("%v:00", period.Start)
			if period.Start < 10 {
				startStr = fmt.Sprintf("0%v:00", period.Start)
			}
			var tempPrice float64
			tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
			if l.member == 1 {
				tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
			}
			//for _, index := range periods {
			//	period := l.periods[int64(index)]
			var minutePrice float64
			var totalMinute int64
			var halfOrHourPrice float64
			var minute int64
			if totalMoney <= 0 {
				continue
			}
			//halfHourPrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
			if currentTime.Format("15:04") == endStr {
				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
				continue
			}
			tempCurrentTime := currentTime
			minutePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()

			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
			// 判断 如果小于1分钟价格则不取整
			if subTotalMinute < 0.001 {
				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			}

			// 不足半小时计费按小时计费，超过半小时按分钟计费
			if ContainsSliceInt64([]int64{6, 7}, l.chargingMode) {
				halfOrHourPrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
				minute = 30
			}
			// 不足一小时计费按小时计费，超过一小时按分钟计费
			if ContainsSliceInt64([]int64{8, 9}, l.chargingMode) {
				halfOrHourPrice = tempPrice
				minute = 60
			}

			// 跨日
			if startStr > endStr {
				if currentTime.Format("15:04") >= startStr {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), endStr), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < endStr {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < endStr && endStr == "24:00" {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
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
			if startStr < endStr && currentTime.Format("15:04") >= startStr && currentTime.Format("15:04") < endStr {
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
				if endStr == "24:00" {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(time.Second*86400).Format("2006-01-02")), time.Local)
				}
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
				hourTotalMoney += surplusMoney
				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
				if totalMoney < 0.001 {
					hourTotalMoney += totalMoney
					totalMoney = 0
				}

				Tduration := decimal.NewFromInt(durationMinute).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
					Index:          index,
					Duration:       Tduration,
					Start:          startStr,
					End:            endStr,
					Price:          tempPrice,
					TotalPrice:     totalMoney,
					HourTotalMoney: hourTotalMoney,
				})
			}
		}
	}
	//hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
	//fmt.Println(hourPeriodLists)
	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
}

// 金额算出会有误差，具体以结束计算为准(以10分钟作为一个收费周期，第1分钟后开始计费)[不足周期计费或跨时段，剩余金额则按分钟计算]
func (l *ChargingTransfer) moneyTransferCycleAndMinuteTime(money float64, startDate string) (string, string, []*PeriodList) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	periodList := make([]*PeriodList, 0)
	//periods := l.sortPeriods()
	totalMoney := money
	lastEndTime = currentTime
	for {
		if totalMoney <= 0 {
			break
		}
		periods := l.sortPeriods(lastEndTime)
		var periodEndTime time.Time
		var duration time.Duration

		for index, period := range periods.Hour {
			var hourTotalMoney float64
			endStr := fmt.Sprintf("%v:00", period.End)
			if period.End < 10 {
				endStr = fmt.Sprintf("0%v:00", period.End)
			}
			startStr := fmt.Sprintf("%v:00", period.Start)
			if period.Start < 10 {
				startStr = fmt.Sprintf("0%v:00", period.Start)
			}
			var tempPrice float64
			tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
			if l.member == 1 {
				tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
			}
			//for _, index := range periods {
			//	period := l.periods[int64(index)]
			var minutePrice float64
			var totalMinute int64
			//var halfOrHourPrice float64
			//var minute int64
			var cyclePrice float64
			if totalMoney <= 0 {
				continue
			}
			//halfHourPrice, _ := decimal.NewFromFloat(period.Price).Div(decimal.NewFromInt(2)).Float64()
			if currentTime.Format("15:04") == endStr {
				//fmt.Println("跳过本次日期", currentTime.Format(time.DateTime), period.End, period.Start)
				continue
			}
			tempCurrentTime := currentTime
			minutePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
			totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()

			tempTotalMinute, _ := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Float64()
			tempTotalMinuteWhole := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			subTotalMinute, _ := decimal.NewFromFloat(tempTotalMinute).Sub(decimal.NewFromInt(tempTotalMinuteWhole)).Float64()
			// 判断 如果小于1分钟价格则不取整
			if subTotalMinute < 0.001 {
				totalMinute = decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).IntPart()
			}

			if l.cycleMinute == 30 {
				cyclePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(2)).Float64()
			}
			if l.cycleMinute == 60 {
				cyclePrice = tempPrice
			}
			if l.cycleMinute != 30 && l.cycleMinute != 60 {
				cyclePrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.cycleMinute)).Float64()
			}

			// 跨日
			if startStr > endStr {
				if currentTime.Format("15:04") >= startStr {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", time.Unix(currentTime.Unix()+86400, 0).Format("2006-01-02"), endStr), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < endStr {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
					duration = periodEndTime.Sub(currentTime)
				}
				if currentTime.Format("15:04") < endStr && endStr == "24:00" {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
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
			if startStr < endStr && currentTime.Format("15:04") >= startStr && currentTime.Format("15:04") < endStr {
				periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", currentTime.Format("2006-01-02"), endStr), time.Local)
				if endStr == "24:00" {
					periodEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s 00:00", currentTime.Add(time.Second*86400).Format("2006-01-02")), time.Local)
				}
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
				hourTotalMoney += surplusMoney
				// 小于0.001元，强制至为0元,因为已不足按1元一小时，每分钟价格
				if totalMoney < 0.001 {
					hourTotalMoney += totalMoney
					totalMoney = 0
				}

				Tduration := decimal.NewFromInt(durationMinute).IntPart()
				periodList = append(periodList, &PeriodList{
					StartTime:      tempCurrentTime.Format("2006-01-02 15:04"),
					EndTime:        lastEndTime.Format("2006-01-02 15:04"),
					Index:          index,
					Duration:       Tduration,
					Start:          startStr,
					End:            endStr,
					Price:          tempPrice,
					TotalPrice:     totalMoney,
					HourTotalMoney: hourTotalMoney,
				})
			}
		}
	}
	//hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
	//fmt.Println(hourPeriodLists)
	return startDate, lastEndTime.Format("2006-01-02 15:04"), periodList
}

// 金额转换成开始及结束时间
func (l *ChargingTransfer) MoneyTransfer(money float64, startDateParam string) (string, string, map[string]HourPeriodList) {
	var startData string
	var endData string
	hourPeriodList := make(map[string]HourPeriodList, 0)
	//计费模式 1按分钟计费 2按半小时计费 3按半小时计费(跨时段) 4按小时计费 5按小时计费(跨时段)
	//6按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费  7按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费(跨时段)
	//8按小时计费开台不足1小时按小时计费，超过1小时按分钟计费   9按小时计费开台不足1小时按小时计费，超过1小时按分钟计费(跨时段)
	//10自定义计费 11自定义计费(跨时段)
	startData, endData, hourPeriodList = l.moneyTransferHourPeek(money, startDateParam)
	return startData, endData, hourPeriodList
	//fmt.Println(startData, endData, hourPeriodList)

	//periodList := make([]*PeriodList, 0)

	//计费模式 1按分钟计费 2按半小时计费 3按半小时计费(跨时段) 4按小时计费 5按小时计费(跨时段)
	//6按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费  7按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费(跨时段)
	//8按小时计费开台不足1小时按小时计费，超过1小时按分钟计费   9按小时计费开台不足1小时按小时计费，超过1小时按分钟计费(跨时段)
	//10自定义计费 11自定义计费(跨时段)

	//// 分钟计费
	//if l.chargingMode == 1 {
	//	startData, endData, periodList = l.moneyTransferMinuteTime(money, startDateParam)
	//	hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
	//	fmt.Println(hourPeriodLists)
	//}
	//// 半小时计费
	//if ContainsSliceInt64([]int64{2, 3}, l.chargingMode) {
	//	startData, endData, periodList = l.moneyTransferHalfHourTime(money, startDateParam)
	//	hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
	//	fmt.Println(hourPeriodLists)
	//}
	//// 小时计费
	//if ContainsSliceInt64([]int64{4, 5}, l.chargingMode) {
	//	startData, endData, periodList = l.moneyTransferHourTime(money, startDateParam)
	//	hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
	//	fmt.Println(hourPeriodLists)
	//}
	//// 按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费
	//if ContainsSliceInt64([]int64{6, 7}, l.chargingMode) {
	//	startData, endData, periodList = l.moneyTransferHalfOrHourAndMinuteTime(money, startDateParam)
	//	hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
	//	fmt.Println(hourPeriodLists)
	//}
	//// 按小时计费开台不足1小时按小时计费，超过1小时按分钟计费
	//if ContainsSliceInt64([]int64{8, 9}, l.chargingMode) {
	//	startData, endData, periodList = l.moneyTransferHalfOrHourAndMinuteTime(money, startDateParam)
	//	hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
	//	fmt.Println(hourPeriodLists)
	//}
	//// 自定义计费
	//if ContainsSliceInt64([]int64{10, 11}, l.chargingMode) {
	//	startData, endData, periodLista := l.peek(money, startDateParam)
	//	fmt.Println(startData, endData, periodLista)
	//	//startData, endData, periodList = l.moneyTransferCycleAndMinuteTime(money, startDateParam)
	//	//hourPeriodLists := l.moneyTransferHourPeriodList(periodList)
	//	//fc, sc := l.computePeak(hourPeriodLists)
	//	//for _, pl := range periodList {
	//	//	fmt.Println(pl, "======")
	//	//}
	//	////fmt.Println(money-fc, "===", hourPeriodLists)
	//	//_, endData, periodList = l.moneyTransferCycleAndMinuteTime(money-fc, endData)
	//	//for _, pl := range periodList {
	//	//	fmt.Println(pl, "......")
	//	//}
	//	//
	//	//hourPeriodLists = l.moneyTransferHourPeriodList(periodList)
	//	//fmt.Println(".====..", hourPeriodLists)
	//	//fc, sc = l.computePeak(hourPeriodLists)
	//	//fmt.Println("结束-----", startData, endData, fc, periodList, sc)
	////}
	//return startData, endData, periodList
}
func (l *ChargingTransfer) checkHourPeek(hourPeak []CPHour, currentTime time.Time) (bool, float64, int) {
	for index, period := range hourPeak {
		//  是否是会员 0否  1是
		var tempPrice float64
		tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
		if l.member == 1 {
			tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
		}
		endStr := fmt.Sprintf("%v:00", period.End)
		if period.End < 10 {
			endStr = fmt.Sprintf("0%v:00", period.End)
		}
		startStr := fmt.Sprintf("%v:00", period.Start)
		if period.Start < 10 {
			startStr = fmt.Sprintf("0%v:00", period.Start)
		}
		// 跨日
		if startStr > endStr {
			var periodEnd time.Time
			var periodStart time.Time
			if currentTime.Format("15:04") < endStr {
				periodStart, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", currentTime.Add(-86400*time.Second).Format("2006-01-02"), startStr), time.Local)
				periodEnd, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", currentTime.Format("2006-01-02"), endStr), time.Local)
				if endStr == "24:00" {
					periodEnd, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v 00:00", currentTime.Format("2006-01-02")), time.Local)
				}
			} else {
				periodStart, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", currentTime.Format("2006-01-02"), startStr), time.Local)
				periodEnd, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", currentTime.Format("2006-01-02"), endStr), time.Local)
				if endStr == "24:00" {
					periodEnd, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v 00:00", currentTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
				}
			}
			//fmt.Println("====11", periodStart.Format(time.DateTime), periodEnd.Format(time.DateTime), currentTime.Format(time.DateTime))
			if currentTime.Unix() >= periodStart.Unix() && currentTime.Unix() < periodEnd.Unix() {
				//fmt.Println("====11 tttttttt")
				return true, tempPrice, index
			}
			//lastEndTime
			//if int64(mins) < period.Start && int64(mins) < period.End {
			//	return true
			//}
			//fmt.Println("跨日")
		}
		// 当日
		if startStr < endStr {
			//fmt.Println("====22")
			periodStart, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", currentTime.Format("2006-01-02"), startStr), time.Local)
			periodEnd, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", currentTime.Format("2006-01-02"), endStr), time.Local)
			//fmt.Println("====22", periodStart.Format(time.DateTime), periodEnd.Format(time.DateTime), currentTime.Format(time.DateTime))
			if currentTime.Unix() >= periodStart.Unix() && currentTime.Unix() < periodEnd.Unix() {
				//fmt.Println("====22 tttttttt")
				return true, tempPrice, index
			}
		}
	}
	return false, 0, -1
}
func (l *ChargingTransfer) moneyTransferHourPeek(money float64, startDate string) (string, string, map[string]HourPeriodList) {
	currentTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	var lastEndTime time.Time
	totalMoney := money
	//periodList := make([]*PeriodList, 0)
	hourPeriodList := make(map[string]HourPeriodList, 0)
	lastEndTime = currentTime
	mpPeek := make(map[string]float64, 0)
	// 将 map 的键放入切片
	for {
		if totalMoney <= 0 {
			//fmt.Println("跳出循环")
			break
		}
		periods := l.sortPeriods(lastEndTime)

		for _, period := range periods.Hour {
			endStr := fmt.Sprintf("%v:00", period.End)
			if period.End < 10 {
				endStr = fmt.Sprintf("0%v:00", period.End)
			}
			startStr := fmt.Sprintf("%v:00", period.Start)
			if period.Start < 10 {
				startStr = fmt.Sprintf("0%v:00", period.Start)
			}
			var minutePrice float64
			var tempPrice float64
			tempPrice, _ = strconv.ParseFloat(period.IdlePrice, 10)
			if l.member == 1 {
				tempPrice, _ = strconv.ParseFloat(period.MemberPrice, 10)
			}
			minutePrice, _ = decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()
			// 跨日 开始及结束时间不能出现24点的时段
			if startStr > endStr {
				var tempEndTime time.Time
				if startStr <= lastEndTime.Format("15:04") && lastEndTime.Format("15:04") > endStr {
					tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", lastEndTime.Add(86400*time.Second).Format("2006-01-02"), endStr), time.Local)
				} else {
					tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", lastEndTime.Format("2006-01-02"), endStr), time.Local)
				}
				for {
					if totalMoney <= 0 {
						break
					}
					nextHourMinute, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v:00", lastEndTime.Add(60*time.Minute).Format("2006-01-02 15")), time.Local)
					if nextHourMinute.Unix() > tempEndTime.Unix() {
						break
					}
					var tempJSPrice float64
					minute := decimal.NewFromFloat(nextHourMinute.Sub(lastEndTime).Minutes()).IntPart()
					if minute == 60 {
						tempJSPrice = tempPrice
					} else {
						tempJSPrice, _ = decimal.NewFromInt(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					}

					var hourPeak int64
					res, cappedPrice, ind := l.checkHourPeek(periods.HourPeak, lastEndTime)
					if res == true {
						hourPeak = 1
						key := fmt.Sprintf(fmt.Sprintf("%v%v", lastEndTime.Format("2006-01-02"), ind))
						if val, ok := mpPeek[key]; ok {
							if val >= cappedPrice {
								tempJSPrice = 0
							} else {
								var tsTotalCpric float64
								tsTotalCpric = tempJSPrice + val
								if tsTotalCpric > cappedPrice {
									tempJSPrice = cappedPrice - val
								}
								mpPeek[key] += tempJSPrice
							}
						} else {
							if cappedPrice < tempJSPrice {
								tempJSPrice = cappedPrice
							}
							mpPeek[key] = tempJSPrice
						}
					}

					if tempJSPrice > totalMoney {
						tempJSPrice = totalMoney
						totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
						nextHourMinute = lastEndTime.Add(time.Duration(totalMinute) * time.Minute)
					}

					hourPeriodList[lastEndTime.Format("2006-01-02 15")] = HourPeriodList{
						StartDate:     lastEndTime.Format("2006-01-02 15:04"),
						EndDate:       nextHourMinute.Format("2006-01-02 15:04"),
						Duration:      decimal.NewFromFloat(nextHourMinute.Sub(lastEndTime).Minutes()).IntPart(),
						HourPrice:     tempPrice,
						Price:         tempJSPrice,
						HourPeak:      hourPeak,
						HourPeakPrice: cappedPrice,
					}
					//fmt.Println("时长:", nextHourMinute.Sub(lastEndTime).Minutes(), "===", lastEndTime.Format(time.DateTime), nextHourMinute.Format(time.DateTime), "*****", tempJSPrice)
					totalMoney -= tempJSPrice
					lastEndTime = nextHourMinute
					continue
				}
			}
			if startStr < endStr && lastEndTime.Format("15:04") >= startStr && lastEndTime.Format("15:04") < endStr {
				var tempEndTime time.Time
				if endStr == "24:00" {
					tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v 00:00", lastEndTime.Add(86400*time.Second).Format("2006-01-02")), time.Local)
				} else {
					tempEndTime, _ = time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v %v", lastEndTime.Format("2006-01-02"), endStr), time.Local)
				}
				for {
					if totalMoney <= 0 {
						break
					}
					nextHourMinute, _ := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%v:00", lastEndTime.Add(60*time.Minute).Format("2006-01-02 15")), time.Local)
					if nextHourMinute.Unix() > tempEndTime.Unix() {
						break
					}
					var tempJSPrice float64
					minute := decimal.NewFromFloat(nextHourMinute.Sub(lastEndTime).Minutes()).IntPart()
					if minute == 60 {
						tempJSPrice = tempPrice
					} else {
						tempJSPrice, _ = decimal.NewFromInt(minute).Mul(decimal.NewFromFloat(minutePrice)).Float64()
					}
					var hourPeak int64
					res, cappedPrice, ind := l.checkHourPeek(periods.HourPeak, lastEndTime)
					if res == true {
						hourPeak = 1
						key := fmt.Sprintf(fmt.Sprintf("%v%v", lastEndTime.Format("2006-01-02"), ind))
						if val, ok := mpPeek[key]; ok {
							if val >= cappedPrice {
								tempJSPrice = 0
							} else {
								var tsTotalCpric float64
								tsTotalCpric = tempJSPrice + val
								if tsTotalCpric > cappedPrice {
									tempJSPrice = cappedPrice - val
								}
								mpPeek[key] += tempJSPrice
							}
						} else {
							if cappedPrice < tempJSPrice {
								tempJSPrice = cappedPrice
							}
							mpPeek[key] = tempJSPrice
						}
					}
					if tempJSPrice > totalMoney {
						tempJSPrice = totalMoney
						totalMinute := decimal.NewFromFloat(totalMoney).Div(decimal.NewFromFloat(minutePrice)).Ceil().IntPart()
						nextHourMinute = lastEndTime.Add(time.Duration(totalMinute) * time.Minute)
					}
					hourPeriodList[lastEndTime.Format("2006-01-02 15")] = HourPeriodList{
						StartDate:     lastEndTime.Format("2006-01-02 15:04"),
						EndDate:       nextHourMinute.Format("2006-01-02 15:04"),
						Duration:      decimal.NewFromFloat(nextHourMinute.Sub(lastEndTime).Minutes()).IntPart(),
						HourPrice:     tempPrice,
						Price:         tempJSPrice,
						HourPeak:      hourPeak,
						HourPeakPrice: cappedPrice,
					}
					//fmt.Println("时长:", nextHourMinute.Sub(lastEndTime).Minutes(), "===", lastEndTime.Format(time.DateTime), nextHourMinute.Format(time.DateTime), "*****", tempJSPrice)
					totalMoney -= tempJSPrice
					lastEndTime = nextHourMinute
				}
			}
		}

	}
	return startDate, lastEndTime.Format("2006-01-02 15:04"), hourPeriodList
}
