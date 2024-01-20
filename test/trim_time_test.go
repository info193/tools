package test

import (
	"fmt"
	"github.com/info193/tools/utils"
	"testing"
)

func TestTimeSpan(t *testing.T) {
	var duration int64 = 120
	periodStartHour := "08:00"
	periodEndHour := "00:00"
	subscribeStartDate := "2024-01-12 23:10"
	subscribeEndDate := "2024-01-13 05:30"
	trimTime := utils.NewTrimTime(duration, periodStartHour, periodEndHour, subscribeStartDate, subscribeEndDate)

	boundaryDuration := trimTime.Neutron * 60 // 边界时间
	if trimTime.EndTime.Unix()-trimTime.StartTime.Unix() < boundaryDuration {
		fmt.Println(fmt.Sprintf("设定开始时间及结束时间小于%v分钟", trimTime.Neutron))
		return
	}
	fmt.Println(trimTime.Period(), "======")
}
