## Example

```
package main

import (
	"fmt"
	"time"

	ratelimiter "github.com/dong-tran/go-utils/rate-limiter"
)

func main() {
	// Initialize a rate limiter with a limit of 1 call per second and 10 tasks allowed to run immediately
	limiter := ratelimiter.NewRateLimiter(time.Second, 10*time.Second, 10, 2)

	startTime := time.Now()
	fmt.Printf("Task started at: %02d:%02d\n", startTime.Minute(), startTime.Second())

	var caller ratelimiter.CallerFunction = func(num ratelimiter.CallerParam) ratelimiter.ReturnResult {
		fmt.Printf("Executing task %d ...\n", num)
		startTime := time.Now()
		fmt.Printf("Task started at: %02d:%02d\n", startTime.Minute(), startTime.Second())
		time.Sleep(500 * time.Millisecond) // Simulating some work
		fmt.Printf("Task %d completed.\n", num)
		return startTime
	}

	// Simulate 20 calls to Execute, which should be rate-limited after the initial window
	for i := 0; i < 20; i++ {
		res := limiter.Execute(caller, i)
		if value, ok := res.(time.Time); ok {
			fmt.Printf("Returned task started at: %02d:%02d\n.", value.Minute(), value.Second())
		}

	}

	stopTime := time.Now()
	fmt.Printf("Task finished at: %02d:%02d\n", stopTime.Minute(), stopTime.Second())
}

```