package test

import (
	"testing"
)

func TestOffsetTimeCharging(t *testing.T) {

	//   测试，已废弃
	//day24HourPeriod := map[int64]string{0: "00:00", 1: "01:00", 2: "02:00", 3: "03:00", 4: "04:00", 5: "05:00", 6: "06:00", 7: "07:00", 8: "08:00", 9: "09:00", 10: "10:00", 11: "11:00", 12: "12:00", 13: "13:00", 14: "14:00", 15: "15:00", 16: "16:00", 17: "17:00", 18: "18:00", 19: "19:00", 20: "20:00", 21: "21:00", 22: "22:00", 23: "23:00", 24: "24:00"}
	//periods := make(map[int64]utils.ChargePeriod, 0)
	//periods[0] = utils.ChargePeriod{EndPeriod: 10, StartPeriod: 5, End: "10:00", Start: "05:00", Price: 5.00}
	//periods[1] = utils.ChargePeriod{EndPeriod: 15, StartPeriod: 10, End: "15:00", Start: "10:00", Price: 8.00}
	//periods[2] = utils.ChargePeriod{EndPeriod: 22, StartPeriod: 15, End: "22:00", Start: "15:00", Price: 6.00}
	//periods[3] = utils.ChargePeriod{EndPeriod: 05, StartPeriod: 22, End: "05:00", Start: "22:00", Price: 12.00}
	//price, periodAll, periodList := utils.NewChargeMode(periods, 2, 0).Outlay("2024-10-30 05:20", "2024-10-30 15:10")
	//fmt.Println(price, periodAll, periodList)
	//for _, value := range periodList {
	//	fmt.Println(value.StartTime, value.EndTime, value.Price, value.TotalPrice)
	//}
	//fmt.Println("1元", decimal.NewFromInt(1).Div(decimal.NewFromInt(60)), decimal.NewFromInt(1).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("2元", decimal.NewFromInt(2).Div(decimal.NewFromInt(60)), decimal.NewFromInt(2).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("3元", decimal.NewFromInt(3).Div(decimal.NewFromInt(60)), decimal.NewFromInt(3).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("4元", decimal.NewFromInt(4).Div(decimal.NewFromInt(60)), decimal.NewFromInt(4).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("5元", decimal.NewFromInt(5).Div(decimal.NewFromInt(60)), decimal.NewFromInt(5).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("6元", decimal.NewFromInt(6).Div(decimal.NewFromInt(60)), decimal.NewFromInt(6).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("7元", decimal.NewFromInt(7).Div(decimal.NewFromInt(60)), decimal.NewFromInt(7).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("8元", decimal.NewFromInt(8).Div(decimal.NewFromInt(60)), decimal.NewFromInt(8).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("9元", decimal.NewFromInt(9).Div(decimal.NewFromInt(60)), decimal.NewFromInt(9).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("10元", decimal.NewFromInt(10).Div(decimal.NewFromInt(60)), decimal.NewFromInt(10).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("11元", decimal.NewFromInt(11).Div(decimal.NewFromInt(60)), decimal.NewFromInt(11).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("12元", decimal.NewFromInt(12).Div(decimal.NewFromInt(60)), decimal.NewFromInt(12).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("13元", decimal.NewFromInt(13).Div(decimal.NewFromInt(60)), decimal.NewFromInt(13).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("14元", decimal.NewFromInt(14).Div(decimal.NewFromInt(60)), decimal.NewFromInt(14).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("15元", decimal.NewFromInt(15).Div(decimal.NewFromInt(60)), decimal.NewFromInt(15).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("16元", decimal.NewFromInt(16).Div(decimal.NewFromInt(60)), decimal.NewFromInt(16).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("17元", decimal.NewFromInt(17).Div(decimal.NewFromInt(60)), decimal.NewFromInt(17).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("18元", decimal.NewFromInt(18).Div(decimal.NewFromInt(60)), decimal.NewFromInt(18).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("19元", decimal.NewFromInt(19).Div(decimal.NewFromInt(60)), decimal.NewFromInt(19).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("20元", decimal.NewFromInt(20).Div(decimal.NewFromInt(60)), decimal.NewFromInt(20).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("21元", decimal.NewFromInt(21).Div(decimal.NewFromInt(60)), decimal.NewFromInt(21).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("22元", decimal.NewFromInt(22).Div(decimal.NewFromInt(60)), decimal.NewFromInt(22).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("23元", decimal.NewFromInt(23).Div(decimal.NewFromInt(60)), decimal.NewFromInt(23).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("24元", decimal.NewFromInt(24).Div(decimal.NewFromInt(60)), decimal.NewFromInt(24).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("25元", decimal.NewFromInt(25).Div(decimal.NewFromInt(60)), decimal.NewFromInt(25).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("26元", decimal.NewFromInt(26).Div(decimal.NewFromInt(60)), decimal.NewFromInt(26).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("27元", decimal.NewFromInt(27).Div(decimal.NewFromInt(60)), decimal.NewFromInt(27).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("28元", decimal.NewFromInt(28).Div(decimal.NewFromInt(60)), decimal.NewFromInt(28).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("29元", decimal.NewFromInt(29).Div(decimal.NewFromInt(60)), decimal.NewFromInt(29).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("30元", decimal.NewFromInt(30).Div(decimal.NewFromInt(60)), decimal.NewFromInt(30).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("31元", decimal.NewFromInt(31).Div(decimal.NewFromInt(60)), decimal.NewFromInt(31).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("32元", decimal.NewFromInt(32).Div(decimal.NewFromInt(60)), decimal.NewFromInt(32).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("33元", decimal.NewFromInt(33).Div(decimal.NewFromInt(60)), decimal.NewFromInt(33).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("34元", decimal.NewFromInt(34).Div(decimal.NewFromInt(60)), decimal.NewFromInt(34).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("35元", decimal.NewFromInt(35).Div(decimal.NewFromInt(60)), decimal.NewFromInt(35).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("36元", decimal.NewFromInt(36).Div(decimal.NewFromInt(60)), decimal.NewFromInt(36).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("37元", decimal.NewFromInt(37).Div(decimal.NewFromInt(60)), decimal.NewFromInt(37).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("38元", decimal.NewFromInt(38).Div(decimal.NewFromInt(60)), decimal.NewFromInt(38).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("39元", decimal.NewFromInt(39).Div(decimal.NewFromInt(60)), decimal.NewFromInt(39).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("40元", decimal.NewFromInt(40).Div(decimal.NewFromInt(60)), decimal.NewFromInt(40).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("41元", decimal.NewFromInt(41).Div(decimal.NewFromInt(60)), decimal.NewFromInt(41).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("42元", decimal.NewFromInt(42).Div(decimal.NewFromInt(60)), decimal.NewFromInt(42).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("43元", decimal.NewFromInt(43).Div(decimal.NewFromInt(60)), decimal.NewFromInt(43).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("44元", decimal.NewFromInt(44).Div(decimal.NewFromInt(60)), decimal.NewFromInt(44).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("45元", decimal.NewFromInt(45).Div(decimal.NewFromInt(60)), decimal.NewFromInt(45).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))
	//fmt.Println("46元", decimal.NewFromInt(46).Div(decimal.NewFromInt(60)), decimal.NewFromInt(46).Div(decimal.NewFromInt(60)).Mul(decimal.NewFromInt(30)).Round(2))

	//minutePrice, _ := decimal.NewFromFloat(2).Div(decimal.NewFromInt(60)).Float64()
	//totalPrice, _ := decimal.NewFromInt(30).Mul(decimal.NewFromFloat(minutePrice)).Float64()
	//fmt.Println("======", minutePrice, totalPrice)
	//fmt.Println(decimal.NewFromFloat(12.025).Round(2).Float64())
	//for _, value := range periodList {
	//	fmt.Println(value.StartTime, value.EndTime, value.Price, value.TotalPrice)
	//}
	//pricec, outlaySpecificss := utils.NewCharge(periods, day24HourPeriod).OutlaySpecifics("2024-04-22 20:40", "2024-04-23 01:59")
	//fmt.Println(pricec, outlaySpecificss, "+++++++")05:01
	//
	//startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferMinuteTime(60.00, "2024-05-10 05:00")
	//startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferHalfHourTime(150.00, "2024-05-10 22:59")
	//startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferHourTime(20.00, "2024-05-10 22:59")
	//startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferCycleAndMinuteTime(20.00, "2024-05-10 22:59")
	//startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferHourOrMinuteTime(20.00, "2024-05-10 00:29")
	//startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferHalfHourOrMinuteTime(20.00, "2024-05-10 00:29")

	//fmt.Println(startDate, endDate)
	// 4按分钟计费不足1小时按小时计费（仅首个1小时内）
	// 计费模式 1按分钟计费 2按半小时计费 3按小时计费  4按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费  5按小时计费开台不足1小时按小时计费，超过1小时按分钟计费  6自定义计费
	// 以10分钟作为一个收费周期，第1分钟后开始计费
	// custom_charge_cycle
	// start_charge_minute
}
