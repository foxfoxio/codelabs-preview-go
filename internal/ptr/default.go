package ptr

import "time"

func String(v string) *string {
	return &v
}

func Bool(v bool) *bool {
	return &v
}

func Int(v int) *int {
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func Int32(v int32) *int32 {
	return &v
}

func Time(t time.Time) *time.Time {
	return &t
}

func Float64(v float64) *float64 {
	return &v
}

func Float32(v float32) *float32 {
	return &v
}
