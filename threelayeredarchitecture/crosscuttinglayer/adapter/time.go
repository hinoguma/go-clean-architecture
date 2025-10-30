package adapter

import (
	"app/threelayeredarchitecture/crosscuttinglayer"
	"time"
)

type TimeGenerator interface {
	GetTz() time.Location
	NowTs() crosscuttinglayer.UnixTimestamp
	NowTsMills() crosscuttinglayer.UnixTimestampMillis
}
