package timegenerator

import (
	"app/crosscutting/timegenerator/infrainterface"
	"app/crosscutting/utils"
	"time"
)

var globalTimeGenerator infrainterface.TimeGenerator = NewStdTimeGenerator()

func SetGlobalTimeGenerator(tg infrainterface.TimeGenerator) {
	globalTimeGenerator = tg
}


func NowTs() utils.UnixTimestamp {
	return globalTimeGenerator.NowTs()
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
