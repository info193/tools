package utils

import (
	"github.com/shopspring/decimal"
	"time"
)

type ChargePeriod struct {
	StartPeriod int64
	EndPeriod   int64
	Start       string
	End         string
	Price       float64
}

type Charging struct {
	periods               map[int64]ChargePeriod
	day24HourPeriod       map[int64]string
	startDiffDuration     int64     // 开始时段相差 分钟
	endDiffDuration       int64     // 结束时段相差 分钟
	startDiffTime         string    // 开始差集时间
	endDiffTime           string    // 结束差集时间
	startDiffTimeDivision string    // 开始差集时段
	endDiffTimeDivision   string    // 结束差集时段
	markStartTime         time.Time // 标记计费开始时间
	markEndTime           time.Time // 标记计费结束时间
}

func NewCharge(periods map[int64]ChargePeriod, day24HourPeriod map[int64]string) *Charging {
	return &Charging{periods: periods, day24HourPeriod: day24HourPeriod}
}

// 计费
func (l *Charging) Outlay(startDate, endDate string) (float64, map[int64]float64) {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	l.markStartTime = startTime
	if startTime.Format("04") != "00" {
		tempStartTime := startTime.Add(time.Hour)
		subscriptStartTime, _ := time.ParseInLocation("2006-01-02 15:04", tempStartTime.Format("2006-01-02 15:00"), time.Local)
		l.startDiffTime = startTime.Format("200601021504")
		l.startDiffTimeDivision = startTime.Format("15:00")
		l.startDiffDuration = int64(subscriptStartTime.Sub(startTime).Minutes())
		l.markStartTime = subscriptStartTime
	}

	endTime, _ := time.ParseInLocation("2006-01-02 15:04", endDate, time.Local)
	l.markEndTime = endTime
	if endTime.Format("04") != "00" {
		subscriptEndTime, _ := time.ParseInLocation("2006-01-02 15:04", endTime.Format("2006-01-02 15:00"), time.Local)
		l.markEndTime = subscriptEndTime
		l.endDiffTime = subscriptEndTime.Format("200601021504")
		l.endDiffTimeDivision = subscriptEndTime.Format("15:04")
		l.endDiffDuration = int64(endTime.Sub(subscriptEndTime).Minutes())
	}

	allDay24HourPeriod := make(map[string]int64)
	for index, value := range l.day24HourPeriod {
		allDay24HourPeriod[value] = index
	}

	periodArr := make(map[string]float64)   // 每个时段数据
	periodDetail := make(map[int64]float64) // 时段费用详情
	// 判断 如果标记计费 开始及结束时间在同一时间则跳过
	if l.markStartTime.Format("200601021504") != l.markEndTime.Format("200601021504") {
		for {
			index := allDay24HourPeriod[l.markStartTime.Format("15:04")]
			for key, value := range l.periods {
				if value.StartPeriod <= index && value.EndPeriod > index {
					tempPrice := value.Price
					periodArr[l.markStartTime.Format("200601021504")] = tempPrice
					if value, ok := periodDetail[key]; ok {
						periodDetail[key], _ = decimal.NewFromFloat(value).Add(decimal.NewFromFloat(tempPrice)).RoundFloor(2).Float64()
					} else {
						periodDetail[key] = tempPrice
					}
				}
			}
			l.markStartTime = l.markStartTime.Add(time.Hour * 1)
			if l.markStartTime.Unix() >= l.markEndTime.Unix() {
				break
			}
		}
	}

	// 计算相差数
	for key, value := range l.periods {
		if l.startDiffTime != "" && l.endDiffTime != "" {
			startDiffIndex := allDay24HourPeriod[l.startDiffTimeDivision]
			endDiffIndex := allDay24HourPeriod[l.endDiffTimeDivision]
			// 同一时段
			if (value.StartPeriod <= startDiffIndex && value.EndPeriod > startDiffIndex) && (value.StartPeriod <= endDiffIndex && value.EndPeriod > endDiffIndex) {
				var finalPrice float64
				tempPrice := value.Price
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64() // 计算每分钟单价
				diff := decimal.NewFromInt(l.startDiffDuration).Add(decimal.NewFromInt(l.endDiffDuration)).IntPart()
				if diff >= 60 {
					residueMinute := decimal.NewFromInt(diff).Sub(decimal.NewFromInt(60)).IntPart()
					residueMinutePrice, _ := decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(residueMinute)).RoundFloor(2).Float64()
					finalPrice, _ = decimal.NewFromFloat(residueMinutePrice).Add(decimal.NewFromFloat(tempPrice)).Float64()
				} else {
					finalPrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(diff)).RoundFloor(2).Float64()
				}

				if value, ok := periodDetail[key]; ok {
					periodDetail[key], _ = decimal.NewFromFloat(value).Add(decimal.NewFromFloat(finalPrice)).RoundFloor(2).Float64()
				} else {
					periodDetail[key] = finalPrice
				}
				periodArr[l.startDiffTime] = finalPrice
				continue
			}
		}

		if l.startDiffTime != "" {
			startDiffIndex := allDay24HourPeriod[l.startDiffTimeDivision]
			if value.StartPeriod <= startDiffIndex && value.EndPeriod > startDiffIndex {
				tempPrice := value.Price
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()                            // 计算每分钟单价
				price, _ := decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.startDiffDuration)).RoundFloor(2).Float64() // 计算价格
				if value, ok := periodDetail[key]; ok {
					periodDetail[key], _ = decimal.NewFromFloat(value).Add(decimal.NewFromFloat(price)).RoundFloor(2).Float64()
				} else {
					periodDetail[key] = price
				}
				periodArr[l.startDiffTime] = price
			}
		}

		if l.endDiffTime != "" {
			endDiffIndex := allDay24HourPeriod[l.endDiffTimeDivision]
			if value.StartPeriod <= endDiffIndex && value.EndPeriod > endDiffIndex {
				tempPrice := value.Price
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()                          // 计算每分钟单价
				price, _ := decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.endDiffDuration)).RoundFloor(2).Float64() // 计算价格
				if value, ok := periodDetail[key]; ok {
					periodDetail[key], _ = decimal.NewFromFloat(value).Add(decimal.NewFromFloat(price)).RoundFloor(2).Float64()
				} else {
					periodDetail[key] = price
				}
				periodArr[l.endDiffTime] = price
			}
		}

	}
	var resultPrice float64
	for _, value := range periodArr {
		resultPrice, _ = decimal.NewFromFloat(resultPrice).Add(decimal.NewFromFloat(value)).Float64()
	}
	return resultPrice, periodDetail
}

func (l *Charging) OutlaySpecifics(startDate, endDate string) (float64, map[string]float64) {
	startTime, _ := time.ParseInLocation("2006-01-02 15:04", startDate, time.Local)
	l.markStartTime = startTime
	if startTime.Format("04") != "00" {
		tempStartTime := startTime.Add(time.Hour)
		subscriptStartTime, _ := time.ParseInLocation("2006-01-02 15:04", tempStartTime.Format("2006-01-02 15:00"), time.Local)
		l.startDiffTime = startTime.Format("200601021504")
		l.startDiffTimeDivision = startTime.Format("15:00")
		l.startDiffDuration = int64(subscriptStartTime.Sub(startTime).Minutes())
		l.markStartTime = subscriptStartTime
	}

	endTime, _ := time.ParseInLocation("2006-01-02 15:04", endDate, time.Local)
	l.markEndTime = endTime
	if endTime.Format("04") != "00" {
		subscriptEndTime, _ := time.ParseInLocation("2006-01-02 15:04", endTime.Format("2006-01-02 15:00"), time.Local)
		l.markEndTime = subscriptEndTime
		l.endDiffTime = subscriptEndTime.Format("200601021504")
		l.endDiffTimeDivision = subscriptEndTime.Format("15:04")
		l.endDiffDuration = int64(endTime.Sub(subscriptEndTime).Minutes())
	}

	allDay24HourPeriod := make(map[string]int64)
	for index, value := range l.day24HourPeriod {
		allDay24HourPeriod[value] = index
	}

	periodArr := make(map[string]float64)   // 每个时段数据
	periodDetail := make(map[int64]float64) // 时段费用详情
	// 判断 如果标记计费 开始及结束时间在同一时间则跳过
	if l.markStartTime.Format("200601021504") != l.markEndTime.Format("200601021504") {
		for {
			index := allDay24HourPeriod[l.markStartTime.Format("15:04")]
			for key, value := range l.periods {
				if value.StartPeriod <= index && value.EndPeriod > index {
					tempPrice := value.Price
					periodArr[l.markStartTime.Format("200601021504")] = tempPrice
					if value, ok := periodDetail[key]; ok {
						periodDetail[key], _ = decimal.NewFromFloat(value).Add(decimal.NewFromFloat(tempPrice)).RoundFloor(2).Float64()
					} else {
						periodDetail[key] = tempPrice
					}
				}
			}
			l.markStartTime = l.markStartTime.Add(time.Hour * 1)
			if l.markStartTime.Unix() >= l.markEndTime.Unix() {
				break
			}
		}
	}

	// 计算相差数
	for key, value := range l.periods {
		if l.startDiffTime != "" && l.endDiffTime != "" {
			startDiffIndex := allDay24HourPeriod[l.startDiffTimeDivision]
			endDiffIndex := allDay24HourPeriod[l.endDiffTimeDivision]
			// 同一时段
			if (value.StartPeriod <= startDiffIndex && value.EndPeriod > startDiffIndex) && (value.StartPeriod <= endDiffIndex && value.EndPeriod > endDiffIndex) {
				var finalPrice float64
				tempPrice := value.Price
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64() // 计算每分钟单价
				diff := decimal.NewFromInt(l.startDiffDuration).Add(decimal.NewFromInt(l.endDiffDuration)).IntPart()
				if diff >= 60 {
					residueMinute := decimal.NewFromInt(diff).Sub(decimal.NewFromInt(60)).IntPart()
					residueMinutePrice, _ := decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(residueMinute)).RoundFloor(2).Float64()
					finalPrice, _ = decimal.NewFromFloat(residueMinutePrice).Add(decimal.NewFromFloat(tempPrice)).Float64()
				} else {
					finalPrice, _ = decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(diff)).RoundFloor(2).Float64()
				}

				if value, ok := periodDetail[key]; ok {
					periodDetail[key], _ = decimal.NewFromFloat(value).Add(decimal.NewFromFloat(finalPrice)).RoundFloor(2).Float64()
				} else {
					periodDetail[key] = finalPrice
				}
				periodArr[l.startDiffTime] = finalPrice
				continue
			}
		}

		if l.startDiffTime != "" {
			startDiffIndex := allDay24HourPeriod[l.startDiffTimeDivision]
			if value.StartPeriod <= startDiffIndex && value.EndPeriod > startDiffIndex {
				tempPrice := value.Price
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()                            // 计算每分钟单价
				price, _ := decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.startDiffDuration)).RoundFloor(2).Float64() // 计算价格
				if value, ok := periodDetail[key]; ok {
					periodDetail[key], _ = decimal.NewFromFloat(value).Add(decimal.NewFromFloat(price)).RoundFloor(2).Float64()
				} else {
					periodDetail[key] = price
				}
				periodArr[l.startDiffTime] = price
			}
		}

		if l.endDiffTime != "" {
			endDiffIndex := allDay24HourPeriod[l.endDiffTimeDivision]
			if value.StartPeriod <= endDiffIndex && value.EndPeriod > endDiffIndex {
				tempPrice := value.Price
				minutePrice, _ := decimal.NewFromFloat(tempPrice).Div(decimal.NewFromInt(60)).Float64()                          // 计算每分钟单价
				price, _ := decimal.NewFromFloat(minutePrice).Mul(decimal.NewFromInt(l.endDiffDuration)).RoundFloor(2).Float64() // 计算价格
				if value, ok := periodDetail[key]; ok {
					periodDetail[key], _ = decimal.NewFromFloat(value).Add(decimal.NewFromFloat(price)).RoundFloor(2).Float64()
				} else {
					periodDetail[key] = price
				}
				periodArr[l.endDiffTime] = price
			}
		}

	}

	var resultPrice float64
	for _, value := range periodArr {
		resultPrice, _ = decimal.NewFromFloat(resultPrice).Add(decimal.NewFromFloat(value)).RoundFloor(2).Float64()
	}
	return resultPrice, periodArr
}
