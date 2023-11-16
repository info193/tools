package calculate

import "github.com/shopspring/decimal"

// 百分比 0.1333 转换成 13.33
func Percentage(discount float32) float32 {
	dDiscount := decimal.NewFromFloat32(discount)
	dBase := decimal.NewFromFloat32(100)
	discountF := dDiscount.Mul(dBase).InexactFloat64()
	return float32(discountF)
}

// 百分比率 13.33 转换成 0.1333
func PercentageRate(discount float32) float64 {
	dDiscount := decimal.NewFromFloat32(discount)
	dBase := decimal.NewFromFloat32(100)
	return dDiscount.Div(dBase).InexactFloat64()
}

// 百分比率 90.50 转换成 9.05折
func PercentageDiscount(discount float32) string {
	dDiscount := decimal.NewFromFloat32(discount)
	dBase := decimal.NewFromFloat32(100)
	sBase := decimal.NewFromFloat32(10)
	return dDiscount.Div(dBase).Mul(sBase).String()
}

// 折扣率 0.1333 转换成 1.333 string
func DiscountRateString(discount float64) string {
	dDiscount := decimal.NewFromFloat(discount)
	dBase := decimal.NewFromFloat(10)
	return dDiscount.Mul(dBase).String()
}

// 折扣率 0.1333 转换成 1.333 float64
func DiscountRateFloat(discount float64) float64 {
	dDiscount := decimal.NewFromFloat(discount)
	dBase := decimal.NewFromFloat(10)
	return dDiscount.Mul(dBase).InexactFloat64()
}

// 分成金额
func DivideIntoFloat(price float64, discount float32) float64 {
	rate := decimal.NewFromFloat(PercentageRate(discount))
	priceFloat := decimal.NewFromFloat(price)
	return priceFloat.Mul(rate).InexactFloat64()
}
