package Timer

import (
	"go.dedis.ch/kyber"
)

type TimerStruct struct {
	Duration            int64
	IsRequesterDuration bool
}

func (T TimerStruct) Verify(IdTOPbKeyMap map[string]kyber.Point) bool {
	return true
}
