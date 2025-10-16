package infra

import (
	"app/experimentarchitecture/crosscutting/timegenerator/infrainterface"
	"app/experimentarchitecture/crosscutting/utils"
	"time"
)

var jst = time.FixedZone("Asia/Tokyo", 9*60*60)

type basedOnPlaceTimeGenerator struct {
	tz time.Location
}

func NewTimeGenerator(tz time.Location) infrainterface.TimeGenerator {
	return &basedOnPlaceTimeGenerator{tz: tz}
}

func NewJstTimeGenerator() infrainterface.TimeGenerator {
	return &basedOnPlaceTimeGenerator{tz: *jst}
}

func (t basedOnPlaceTimeGenerator) GetTz() time.Location {
	return t.tz
}

func (t basedOnPlaceTimeGenerator) NowTs() utils.UnixTimestamp {
	return utils.UnixTimestamp(time.Now().In(&t.tz).Unix())
}
