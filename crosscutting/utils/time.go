package utils

import "time"

type UnixTimestamp int64

func (v UnixTimestamp) Int64() int64 {
	return int64(v)
}

type UnixTimestampMillis int64

func (v UnixTimestampMillis) Int64() int64 {
	return int64(v)
}

func (v UnixTimestamp) Ymd_Hms() Ymd_Hms {
	return Ymd_Hms(
		time.Unix(v.Int64(), 0).Format("2006-01-02 15:04:05"),
	)
}

func (v UnixTimestampMillis) Ymd_Hms_ms() Ymd_Hms_ms {
	sec := v.Int64() / 1000
	nsec := v.Int64() % 1000 * 1_000_000
	return Ymd_Hms_ms(
		time.Unix(sec, nsec).Format("2006-01-02 15:04:05.000"),
	)
}

type Ymd_Hms string

func (v Ymd_Hms) String() string {
	return string(v)
}

type YyyyMmDd string

func (v YyyyMmDd) String() string {
	return string(v)
}

// layout: "2006-01-02 15:04:05.000"
type Ymd_Hms_ms string

func (v Ymd_Hms_ms) String() string {
	return string(v)
}
