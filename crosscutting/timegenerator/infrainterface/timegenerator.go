package infrainterface

import (
	"app/crosscutting/utils"
	"time"
)

type TimeGenerator interface {
	GetTz() time.Location
	NowTs() utils.UnixTimestamp
	NowTsMills() utils.UnixTimestampMillis
}
