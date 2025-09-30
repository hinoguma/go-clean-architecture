package utils

type UnixTimestamp int64

func (v UnixTimestamp) Int64() int64 {
	return int64(v)
}
