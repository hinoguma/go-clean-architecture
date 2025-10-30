package infra

import (
	"app/threelayeredarchitecture/crosscuttinglayer"
	"app/threelayeredarchitecture/crosscuttinglayer/adapter"
	"time"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

type basedOnPlaceTimeGenerator struct {
	tz time.Location
}

func (t basedOnPlaceTimeGenerator) GetTz() time.Location {
	return t.tz
}

func (t basedOnPlaceTimeGenerator) NowTs() crosscuttinglayer.UnixTimestamp {
	return crosscuttinglayer.UnixTimestamp(time.Now().In(&t.tz).Unix())
}

func (t basedOnPlaceTimeGenerator) NowTsMills() crosscuttinglayer.UnixTimestampMillis {
	//TODO implement me
	panic("implement me")
}

func NewTimeGenerator(tz time.Location) adapter.TimeGenerator {
	return &basedOnPlaceTimeGenerator{tz: tz}
}

func NewJstTimeGenerator() adapter.TimeGenerator {
	return &basedOnPlaceTimeGenerator{tz: *jst}
}
