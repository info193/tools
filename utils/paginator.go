package utils

import (
	"math"
)

type Paginator struct {
	count     int64
	page      int64
	prepage   int64
	totalPage int64
}

func NewPaginator(page, prepage, count int64) *Paginator {
	var totalPage int64
	totalPage = 0
	if prepage > 0 {
		totalPage = int64(math.Ceil(float64(count) / float64(prepage)))
	}
	return &Paginator{
		page:      page,
		prepage:   prepage,
		count:     count,
		totalPage: totalPage,
	}
}
func (this *Paginator) Page() int64 {
	return this.page
}
func (this *Paginator) Prepage() int64 {
	return this.prepage
}
func (this *Paginator) TotalPage() int64 {
	return this.totalPage
}
func (this *Paginator) Count() int64 {
	return this.count
}
