package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestCharging(t *testing.T) {
	//totalAmount := 60.0 // 总金额60元
	//unitCost := 0.01    // 每分钟计费0.033333元
	//minutes := 0.0      // 初始化时间（分钟）为0
	//
	//// 循环直到总金额达到
	//for totalAmount > 0 {
	//	totalAmount -= unitCost
	//	minutes += 1.0
	//	//time.Sleep(time.Minute) // 等待一分钟
	//}

	//// 输出使用的时间
	//fmt.Printf("所需时间: %.2f 分钟\n", minutes)
	//
	//return 18.5+1.5

	//day24HourPeriod := map[int64]string{0: "00:00", 1: "01:00", 2: "02:00", 3: "03:00", 4: "04:00", 5: "05:00", 6: "06:00", 7: "07:00", 8: "08:00", 9: "09:00", 10: "10:00", 11: "11:00", 12: "12:00", 13: "13:00", 14: "14:00", 15: "15:00", 16: "16:00", 17: "17:00", 18: "18:00", 19: "19:00", 20: "20:00", 21: "21:00", 22: "22:00", 23: "23:00", 24: "24:00"}
	periods := make(map[int64]utils.ChargePeriod, 0)
	periods[0] = utils.ChargePeriod{EndPeriod: 10, StartPeriod: 5, End: "10:00", Start: "05:00", Price: 5.00}
	periods[1] = utils.ChargePeriod{EndPeriod: 15, StartPeriod: 10, End: "15:00", Start: "10:00", Price: 8.00}
	periods[2] = utils.ChargePeriod{EndPeriod: 22, StartPeriod: 15, End: "22:00", Start: "15:00", Price: 6.00}
	periods[3] = utils.ChargePeriod{EndPeriod: 05, StartPeriod: 22, End: "05:00", Start: "22:00", Price: 12.00}
	//price, periodAll, periodList := utils.NewCharge(periods, 6, 10).Outlay("2024-10-30 04:56", "2024-10-30 06:10") 29+24
	//fmt.Println(price, periodAll, "-----", periodList)
	//pricec, outlaySpecificss := utils.NewCharge(periods, day24HourPeriod).OutlaySpecifics("2024-04-22 20:40", "2024-04-23 01:59")
	//fmt.Println(pricec, outlaySpecificss, "+++++++")05:01
	//
	//startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferMinuteTime(60.00, "2024-05-10 05:00")
	//startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferHalfHourTime(150.00, "2024-05-10 22:59")
	//startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferHourTime(20.00, "2024-05-10 22:59")
	//startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferCycleAndMinuteTime(20.00, "2024-05-10 22:59")
	//startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferHourOrMinuteTime(20.00, "2024-05-10 00:29")
	startDate, endDate := utils.NewCharge(periods, 6, 10).MoneyTransferHalfHourOrMinuteTime(20.00, "2024-05-10 00:29")

	fmt.Println(startDate, endDate)
	// 4按分钟计费不足1小时按小时计费（仅首个1小时内）
	// 计费模式 1按分钟计费 2按半小时计费 3按小时计费  4按半小时计费开台不足半小时按半小时计费，超过半小时按分钟计费  5按小时计费开台不足1小时按小时计费，超过1小时按分钟计费  6自定义计费
	// 以10分钟作为一个收费周期，第1分钟后开始计费
	// custom_charge_cycle
	// start_charge_minute
}
