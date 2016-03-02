// format
package web

import (
	"fmt"
	"time"
)

type TimeToStringConvert struct {
	Layout string
}

func (x TimeToStringConvert) WriteExt(interface{}) []byte {
	panic("unsupported")
}
func (x TimeToStringConvert) ReadExt(interface{}, []byte) {
	panic("unsupported")
}
func (x TimeToStringConvert) ConvertExt(v interface{}) interface{} {
	switch v2 := v.(type) {
	case time.Time:
		return v2.Format(x.Layout)
	case *time.Time:
		return v2.Format(x.Layout)
	default:
		panic(fmt.Sprintf("unsupported format for time conversion: expecting time.Time; got %T", v))
	}
}
func (x TimeToStringConvert) UpdateExt(dest interface{}, v interface{}) {
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
