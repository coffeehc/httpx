// format
package web

import (
	"fmt"
	"time"
)

type TimeToNanoFormat struct{}

func (x TimeToNanoFormat) WriteExt(interface{}) []byte { panic("unsupported") }
func (x TimeToNanoFormat) ReadExt(interface{}, []byte) { panic("unsupported") }
func (x TimeToNanoFormat) ConvertExt(v interface{}) interface{} {
	switch v2 := v.(type) {
	case time.Time:
		return v2.UnixNano()
	case *time.Time:
		return v2.UnixNano()
	default:
		panic(fmt.Sprintf("unsupported format for time conversion: expecting time.Time; got %T", v))
	}
}
func (x TimeToNanoFormat) UpdateExt(dest interface{}, v interface{}) {
	tt := dest.(*time.Time)
	switch v2 := v.(type) {
	case int64:
		*tt = time.Unix(0, v2)
	case uint64:
		*tt = time.Unix(0, int64(v2))
	default:
		panic(fmt.Sprintf("unsupported format for time conversion: expecting int64/uint64; got %T", v))
	}
}

type TimeToStringFormat struct {
	Layout string
}

func (x TimeToStringFormat) WriteExt(interface{}) []byte { panic("unsupported") }
func (x TimeToStringFormat) ReadExt(interface{}, []byte) { panic("unsupported") }
func (x TimeToStringFormat) ConvertExt(v interface{}) interface{} {
	switch v2 := v.(type) {
	case time.Time:
		return v2.Format(x.Layout)
	case *time.Time:
		return v2.Format(x.Layout)
	default:
		panic(fmt.Sprintf("unsupported format for time conversion: expecting time.Time; got %T", v))
	}
}
func (x TimeToStringFormat) UpdateExt(dest interface{}, v interface{}) {
	tt := dest.(*time.Time)
	var err error
	switch v2 := v.(type) {
	case string:
		*tt, err = time.Parse(x.Layout, v2)
		if err != nil {
			tt = nil
		}
	default:
		panic(fmt.Sprintf("unsupported format for time conversion: expecting string; got %T", v))
	}
}
