package web

import (
	"fmt"
	"strconv"
	"time"
)

type ConvertError struct {
	typeName string
	msg      string
}

func (this ConvertError) Error() string {
	return this.msg
}

func newConvertError(typeName string) *ConvertError {
	return &ConvertError{
		typeName: typeName,
		msg:      "not convert to " + typeName,
	}
}

var (
	ERROR_CONVERT_INT     = newConvertError("int")
	ERROR_CONVERT_INT32   = newConvertError("int32")
	ERROR_CONVERT_INT64   = newConvertError("int64")
	ERROR_CONVERT_FLOAT32 = newConvertError("float32")
	ERROR_CONVERT_FLOAT64 = newConvertError("float64")
	ERROR_CONVERT_TIME    = newConvertError("Time")
)

type RequestParam string

func (this RequestParam) AsInt() (int, error) {
	i, err := strconv.Atoi(this.AsString())
	if err != nil {
		return 0, ERROR_CONVERT_INT
	}
	return i, nil
}

func (this RequestParam) AsInt32() (int32, error) {
	i, err := strconv.ParseInt(this.AsString(), 10, 32)
	if err != nil {
		return 0, ERROR_CONVERT_INT32
	}
	return int32(i), nil
}

func (this RequestParam) AsInt64() (int64, error) {
	i, err := strconv.ParseInt(this.AsString(), 10, 64)
	if err != nil {
		return 0, ERROR_CONVERT_INT64
	}
	return i, nil
}

func (this RequestParam) AsFloat32() (float32, error) {
	f, err := strconv.ParseFloat(this.AsString(), 32)
	if err != nil {
		return 0, ERROR_CONVERT_FLOAT32
	}
	return float32(f), nil
}
func (this RequestParam) AsFloat64() (float64, error) {
	f, err := strconv.ParseFloat(this.AsString(), 64)
	if err != nil {
		return 0, ERROR_CONVERT_FLOAT64
	}
	return f, nil
}

func (this RequestParam) AsString() string {
	return string(this)
}

func (this RequestParam) AsTime(layout string) (time.Time, error) {
	t, err := time.Parse(layout, this.AsString())
	if err != nil {
		return time.Now(), ERROR_CONVERT_TIME
	}
	return t, err
}

type PathFragment map[string]RequestParam

var EmptyParam = RequestParam("")

func (this PathFragment) Get(key string) (RequestParam, error) {
	if p, ok := this[key]; ok {
		return p, nil
	}
	return EmptyParam, fmt.Errorf("param [%s] is nil", key)
}

func (this PathFragment) GetAsInt(key string) (int, error) {
	v, err := this.Get(key)
	if err != nil {
		return 0, err
	}
	return v.AsInt()
}

func (this PathFragment) GetAsInt32(key string) (int32, error) {
	v, err := this.Get(key)
	if err != nil {
		return 0, err
	}
	return v.AsInt32()
}

func (this PathFragment) GetAsInt64(key string) (int64, error) {
	v, err := this.Get(key)
	if err != nil {
		return 0, err
	}
	return v.AsInt64()
}

func (this PathFragment) GetAsFloat32(key string) (float32, error) {
	v, err := this.Get(key)
	if err != nil {
		return 0, err
	}
	return v.AsFloat32()
}

func (this PathFragment) GetAsFloat64(key string) (float64, error) {
	v, err := this.Get(key)
	if err != nil {
		return 0, err
	}
	return v.AsFloat64()
}

func (this PathFragment) GetAsTime(key string, layout string) (time.Time, error) {
	v, err := this.Get(key)
	if err != nil {
		return time.Now(), err
	}
	return v.AsTime(layout)
}

func (this PathFragment) GetAsString(key string) (string, error) {
	v, err := this.Get(key)
	if err != nil {
		return "", err
	}
	return v.AsString(), nil
}
