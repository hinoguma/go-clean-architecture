package utils

import "time"

type UnixTimestamp int64

func (v UnixTimestamp) Int64() int64 {
	return int64(v)
}

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

func GetNowUnixTime() UnixTimestamp {
	return UnixTimestamp(time.Now().Unix())
}
