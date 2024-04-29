package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestTrimTime(t *testing.T) {

	var duration int64 = 60
	periodStartHour := "20:00"
	periodEndHour := "22:30"
	// 隔日 2
	//trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, "2024-04-23 05:02", "2024-04-23 07:00")

	// 隔日  0000000000
	trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, "2024-04-29 16:22", "2024-04-29 17:22")

	//隔日 开始时段 12:00 - 04:00  预约时间2024-04-06 23:59 - 2024-04-07 02:30  2024-04-06 23:59 - 2024-04-07 02:30
	//trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, "2024-04-22 23:00", "2024-04-23 02:55")

	// 隔日 开始时段 12:00 - 02:00  预约时间2024-04-06 23:59 - 2024-04-07 01:30  时长180
	//trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, "2024-04-22 23:59", "2024-04-23 01:30")

	// 今日 开始时段 02:00 - 12:00  预约时间2024-04-22 06:20 - 2024-04-22 12:10  时长180   返回 111111111
	//trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, "2024-04-22 06:20", "2024-04-22 12:10")

	// 今日 开始时段 02:00 - 12:00  预约时间2024-04-22 06:20 - 2024-04-22 09:10  时长180   返回 22222222222
	//trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, "2024-04-22 06:20", "2024-04-22 09:10")

	// 今日 开始时段 02:00 - 12:00  预约时间2024-04-22 06:20 - 2024-04-22 09:50  时长180   返回 3333333333
	//trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, "2024-04-22 06:20", "2024-04-22 09:50")

	//// 今日 开始时段 02:00 - 12:00  预约时间2024-04-22 00:20 - 2024-04-22 02:00  时长180   返回 0000000000 不可使用
	//trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, "2024-04-22 00:20", "2024-04-22 02:00")

	// 今日 开始时段 02:00 - 12:00  预约时间2024-04-22 00:20 - 2024-04-22 05:00  时长180   返回 0000000000 不可使用
	//trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, "2024-04-22 02:00", "2024-04-22 18:35")

	//coupon := trimTime.CardPeriod()
	coupon := trimTime.CouponPeriod(utils.OptimalNo)
	fmt.Println(fmt.Sprintf("%+v", coupon))

}
