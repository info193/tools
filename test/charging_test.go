package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestCharging(t *testing.T) {
	day24HourPeriod := map[int64]string{0: "00:00", 1: "01:00", 2: "02:00", 3: "03:00", 4: "04:00", 5: "05:00", 6: "06:00", 7: "07:00", 8: "08:00", 9: "09:00", 10: "10:00", 11: "11:00", 12: "12:00", 13: "13:00", 14: "14:00", 15: "15:00", 16: "16:00", 17: "17:00", 18: "18:00", 19: "19:00", 20: "20:00", 21: "21:00", 22: "22:00", 23: "23:00", 24: "24:00"}
	periods := make(map[int64]utils.ChargePeriod, 0)
	periods[0] = utils.ChargePeriod{EndPeriod: 8, StartPeriod: 0, End: "00:00", Start: "08:00", Price: 1.99}
	periods[1] = utils.ChargePeriod{EndPeriod: 12, StartPeriod: 8, End: "08:00", Start: "12:00", Price: 9.90}
	periods[2] = utils.ChargePeriod{EndPeriod: 18, StartPeriod: 12, End: "12:00", Start: "18:00", Price: 10.00}
	periods[3] = utils.ChargePeriod{EndPeriod: 24, StartPeriod: 18, End: "18:00", Start: "24:00", Price: 60.00}
	price, periodAll := utils.NewCharge(periods, day24HourPeriod).Outlay("2024-04-22 20:40", "2024-04-23 01:59")
	fmt.Println(price, periodAll, "-----")
	pricec, outlaySpecificss := utils.NewCharge(periods, day24HourPeriod).OutlaySpecifics("2024-04-22 20:40", "2024-04-23 01:59")
	fmt.Println(pricec, outlaySpecificss, "+++++++")

	startDate, endDate := utils.NewCharge(periods, day24HourPeriod).MoneyTransferTime(60.00, "2024-05-22 22:40")
	fmt.Println(startDate, endDate)
}
