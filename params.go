package web

import (
	"strconv"
	"time"
)

type RequestParam string

func (this RequestParam) AsInt() int {
	i, _ := strconv.Atoi(this.AsString())
	return i
}

func (this RequestParam) AsInt32() int32 {
	i, _ := strconv.ParseInt(this.AsString(), 10, 32)
	return int32(i)
}

func (this RequestParam) AsInt64() int64 {
	i, _ := strconv.ParseInt(this.AsString(), 10, 64)
	return i
}

func (this RequestParam) AsFloat32() float32 {
	f, _ := strconv.ParseFloat(this.AsString(), 32)
	return float32(f)
}
func (this RequestParam) AsFloat64() float64 {
	f, _ := strconv.ParseFloat(this.AsString(), 64)
	return f
}

func (this RequestParam) AsString() string {
	return string(this)
}

func (this RequestParam) AsTime(layout string) time.Time {
	t, _ := time.Parse(layout, this.AsString())
	return t
}

type PathFragment map[string]RequestParam

var EmptyParam = RequestParam("")
