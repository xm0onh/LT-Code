package ConActInterface

import (
	"fmt"
	"sort"
)

func (c *ConActor) CollectRespondersTime(Duration int64) {
	c.ResponderTimeDurationSlice = append(c.ResponderTimeDurationSlice, Duration)
	if len(c.ResponderTimeDurationSlice) == len(c.NodeIdToDialConnMapResponders) {
		sort.Slice(c.ResponderTimeDurationSlice, func(i, j int) bool {
			return c.ResponderTimeDurationSlice[i] < c.ResponderTimeDurationSlice[j]

		})
	}
	fmt.Println("Sorted Responder Timer Slice is", c.ResponderTimeDurationSlice)

}
