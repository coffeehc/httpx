package httpx

import (
	"fmt"
	"strconv"
	"time"
)

//ConvertError 请求参数转换错误
type ConvertError struct {
	typeName string
	msg      string
}

func (ce ConvertError) Error() string {
	return ce.msg
}

func newConvertError(typeName string) *ConvertError {
	return &ConvertError{
		typeName: typeName,
		msg:      "not convert to " + typeName,
	}
}

var (
	errorConvertInt     = newConvertError("int")
	errorConvertInt32   = newConvertError("int32")
	errorConvertInt64   = newConvertError("int64")
	errorConvertFloat32 = newConvertError("float32")
	errorConvertFloat64 = newConvertError("float64")
	errorConvertTime    = newConvertError("Time")
)

//RequestParam 统一Request 请求参数类型
type RequestParam string

//AsInt 转换为 Int 类型
func (rp RequestParam) AsInt() (int, error) {
	i, err := strconv.Atoi(rp.AsString())
	if err != nil {
		return 0, errorConvertInt
	}
	return i, nil
}

//AsInt32 to int32
func (rp RequestParam) AsInt32() (int32, error) {
	i, err := strconv.ParseInt(rp.AsString(), 10, 32)
	if err != nil {
		return 0, errorConvertInt32
	}
	return int32(i), nil
}

//AsInt64 to int64
func (rp RequestParam) AsInt64() (int64, error) {
	i, err := strconv.ParseInt(rp.AsString(), 10, 64)
	if err != nil {
		return 0, errorConvertInt64
	}
	return i, nil
}

//AsFloat32 to float32
func (rp RequestParam) AsFloat32() (float32, error) {
	f, err := strconv.ParseFloat(rp.AsString(), 32)
	if err != nil {
		return 0, errorConvertFloat32
	}
	return float32(f), nil
}

//AsFloat64 to float64
func (rp RequestParam) AsFloat64() (float64, error) {
	f, err := strconv.ParseFloat(rp.AsString(), 64)
	if err != nil {
		return 0, errorConvertFloat64
	}
	return f, nil
}

//AsString to string
func (rp RequestParam) AsString() string {
	return string(rp)
}

//AsTime to time type
func (rp RequestParam) AsTime(layout string) (time.Time, error) {
	t, err := time.Parse(layout, rp.AsString())
	if err != nil {
		return time.Now(), errorConvertTime
	}
	return t, err
}

//PathFragment 动态 Path 里面的变量片段
type PathFragment map[string]RequestParam

//EmptyParam 空参数
var EmptyParam = RequestParam("")

//Get get RequestParam form key
func (pf PathFragment) Get(key string) (RequestParam, error) {
	if p, ok := pf[key]; ok {
		return p, nil
	}
	return EmptyParam, fmt.Errorf("param [%s] is nil", key)
}

//GetAsInt get a int from key
func (pf PathFragment) GetAsInt(key string) (int, error) {
	v, err := pf.Get(key)
	if err != nil {
		return 0, err
	}
	return v.AsInt()
}

//GetAsInt32 get a int32 form key
func (pf PathFragment) GetAsInt32(key string) (int32, error) {
	v, err := pf.Get(key)
	if err != nil {
		return 0, err
	}
	return v.AsInt32()
}

//GetAsInt64 get a int64 from from key
func (pf PathFragment) GetAsInt64(key string) (int64, error) {
	v, err := pf.Get(key)
	if err != nil {
		return 0, err
	}
	return v.AsInt64()
}

//GetAsFloat32 get a float32 from key
func (pf PathFragment) GetAsFloat32(key string) (float32, error) {
	v, err := pf.Get(key)
	if err != nil {
		return 0, err
	}
	return v.AsFloat32()
}

//GetAsFloat64 get a float64 from key
func (pf PathFragment) GetAsFloat64(key string) (float64, error) {
	v, err := pf.Get(key)
	if err != nil {
		return 0, err
	}
	return v.AsFloat64()
}

//GetAsTime get a Time from key
func (pf PathFragment) GetAsTime(key string, layout string) (time.Time, error) {
	v, err := pf.Get(key)
	if err != nil {
		return time.Now(), err
	}
	return v.AsTime(layout)
}

//GetAsString get a string form key
func (pf PathFragment) GetAsString(key string) (string, error) {
	v, err := pf.Get(key)
	if err != nil {
		return "", err
	}
	return v.AsString(), nil
}
