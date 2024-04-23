package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestTimeSpan(t *testing.T) {
	//current := time.Now()
	//c, _ := time.ParseInLocation("2006-01-02 15:04", "2024-04-06 10:20", time.Local)
	//if current.After(c) {
	//	fmt.Println("当前时间大于结束时间")
	//}
	//
	//if c.Before(current) {
	//	fmt.Println("结束时间大于当前时间")
	//}

	var duration int64 = 60
	periodStartHour := "20:00"
	periodEndHour := "22:30"
	//subscribeStartDate := "2024-04-09 23:11"
	//subscribeEndDate := "2024-04-10 01:00"
	//trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, subscribeStartDate, subscribeEndDate)
	//ts := trimTime.Period()
	//fmt.Println(fmt.Sprintf("%+v", ts))

	//trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, "2024-04-23 18:25", "2024-04-23 22:25")
	//coupon := trimTime.Period()
	//fmt.Println(coupon)

	trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, "2024-04-23 18:57", "2024-04-23 20:57")
	coupon := trimTime.Period()
	fmt.Println(coupon)

	//fmt.Println(trimTime.Period())
	//boundaryDuration := trimTime.Neutron * 60 // 边界时间
	//if trimTime.EndTime.Unix()-trimTime.StartTime.Unix() < boundaryDuration {
	//	fmt.Println(fmt.Sprintf("设定开始时间及结束时间小于%v分钟", trimTime.Neutron))
	//	return
	//}
	//fmt.Println(trimTime.Period(), "======")
	//给定优惠券时段限制有两种情况：
	//第一种情况跨天时段：开始时段12:00、结束时段05:00
	//第二种情况每日时段：开始时段08:00、结束时段19:00
	//请帮我写一个需求，使用golang语言，
	//给定优惠券时段限制开始时段08:00、结束时段19:00
	//检测用户预约开始时间及结束时间，并且用户选择的优惠券时长1小时，请判断优惠券是否可用，并算出优惠券在用户预约开始时间及结束时间那个时间点使用，在可用或不可用的情况下都给出可用的时间点。
	//用户预约时间可用预约好几天，需检测预约时间段内优惠券是否可用且算出优惠券在最近的某一天预约的时间内使用，并且在可用或不可用的情况下都给出可用的时间范围。
}
