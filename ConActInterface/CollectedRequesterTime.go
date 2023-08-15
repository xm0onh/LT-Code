package ConActInterface

import (
	"fmt"
	"sort"
)

func (c *ConActor) CollectRequestersTime(Duration int64) {
	c.RequesterTimeDurationSlice = append(c.RequesterTimeDurationSlice, Duration)
	if len(c.RequesterTimeDurationSlice) == len(c.NodeIdToDialConnMapRequestors) {
		sort.Slice(c.RequesterTimeDurationSlice, func(i, j int) bool {
			return c.RequesterTimeDurationSlice[i] < c.RequesterTimeDurationSlice[j]

		})
		fmt.Println("Sorted Requester Timer Slice is", c.RequesterTimeDurationSlice)
	}

}
