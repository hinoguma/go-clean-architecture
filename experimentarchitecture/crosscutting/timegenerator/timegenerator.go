package timegenerator

import (
	"app/experimentarchitecture/crosscutting/timegenerator/infrainterface"
	"app/experimentarchitecture/crosscutting/utils"
	"time"
)

var globalTimeGenerator infrainterface.TimeGenerator = NewStdTimeGenerator()

func SetGlobalTimeGenerator(tg infrainterface.TimeGenerator) {
	globalTimeGenerator = tg
}

func NowTs() utils.UnixTimestamp {
	return globalTimeGenerator.NowTs()
}

func NowTsMills() utils.UnixTimestampMillis {
	return globalTimeGenerator.NowTsMills()
}

func NowYmd_Hms_ms() utils.Ymd_Hms_ms {
	return NowTsMills().Ymd_Hms_ms()
}

type StdTimeGenerator struct {
	tz time.Location
}

func NewStdTimeGenerator() infrainterface.TimeGenerator {
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	return &StdTimeGenerator{tz: *jst}
}

func (t StdTimeGenerator) GetTz() time.Location {
	return t.tz
}

func (t StdTimeGenerator) NowTs() utils.UnixTimestamp {
	return utils.UnixTimestamp(time.Now().In(&t.tz).Unix())
}

func (t StdTimeGenerator) NowTsMills() utils.UnixTimestampMillis {
	return utils.UnixTimestampMillis(time.Now().In(&t.tz).UnixMilli())
}
