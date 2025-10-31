package adapter

import (
	"app/threelayeredarchitecture/crosscuttinglayer"
	"time"
)

type TimeGenerator interface {
	GetTz() time.Location
	Now() time.Time
	NowTs() crosscuttinglayer.UnixTimestamp
	NowTsMills() crosscuttinglayer.UnixTimestampMillis
}
