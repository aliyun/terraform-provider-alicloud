package main

import (
	"fmt"
	"math/rand"
	"time"
)

var incrementalWaitJitter = rand.Float64

func incrementalWait(firstDuration time.Duration, increaseDuration time.Duration) func() {
	retryCount := 0
	return func() {
		if retryCount > 40 {
			retryCount = 40
		}
		// Exponential backoff with full jitter (AWS style):
		// sleep = random(0, first + (2^n - 1)*increase).
		backoff := firstDuration + time.Duration(1<<retryCount-1)*increaseDuration
		waitTime := time.Duration(float64(backoff) * incrementalWaitJitter())
		time.Sleep(waitTime)
		fmt.Println(waitTime.String())
		retryCount++
	}
}
func main() {
	//wait := incrementalWait(1*time.Second, 1*time.Second)
	//for _ = range 10 {
	//	wait()
	//}
	for _ = range 10 {
		fmt.Printf("%f\n", incrementalWaitJitter())
	}
}
