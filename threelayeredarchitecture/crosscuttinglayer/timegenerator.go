package crosscuttinglayer

import (
	"app/threelayeredarchitecture/crosscuttinglayer/adapter"
	"time"
)

var globalTimeGenerator adapter.TimeGenerator = NewStdTimeGenerator()

func SetGlobalTimeGenerator(tg adapter.TimeGenerator) {
	globalTimeGenerator = tg
}

func NowTs() UnixTimestamp {
	return globalTimeGenerator.NowTs()
}

func Now() time.Time {
	return globalTimeGenerator.Now()
}

func NowTsMills() UnixTimestampMillis {
	return globalTimeGenerator.NowTsMills()
}

func NowYmd_Hms_ms() Ymd_Hms_ms {
	return NowTsMills().Ymd_Hms_ms()
}

type StdTimeGenerator struct {
	tz time.Location
}

func NewStdTimeGenerator() adapter.TimeGenerator {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return &StdTimeGenerator{tz: *jst}
}

func (t StdTimeGenerator) GetTz() time.Location {
	return t.tz
}

func (t StdTimeGenerator) NowTs() UnixTimestamp {
	return UnixTimestamp(time.Now().In(&t.tz).Unix())
}

func (t StdTimeGenerator) NowTsMills() UnixTimestampMillis {
	return UnixTimestampMillis(time.Now().In(&t.tz).UnixMilli())
}
