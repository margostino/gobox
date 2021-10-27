package common

import (
	"fmt"
	"time"
)

func RegisterTime(timeType string, requestId int) {
	fmt.Printf("%s #%d: %s\n", timeType, requestId, time.Now().String())
}
