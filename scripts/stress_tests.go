package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"math/rand"
)

func main() {
	const (
		proxyURL    = "http://localhost:8080/books"
		totalUsers  = 30 // 30 people clicking at the same time
		reqPerUser  = 5  // Each person clicks 5 times
	)

	var wg sync.WaitGroup
	start := time.Now()

	fmt.Printf("🚀 Starting Load Test: %d users, %d requests each\n", totalUsers, reqPerUser)

	for i := 1; i <= totalUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			for j := 1; j <= reqPerUser; j++ {
				// Simulate random "thinking time" between clicks
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

				fmt.Printf("👤 User %d sending request %d...\n", userID, j)

				resp, err := http.Get(proxyURL)
				if err != nil {
					fmt.Printf("❌ User %d error: %v\n", userID, err)
					continue
				}
				if resp.StatusCode != http.StatusOK {
    				fmt.Printf("⚠️ User %d: Proxy returned ERROR %d (That's why it was fast!)\n", userID, resp.StatusCode)
				} else {
    				fmt.Printf("✅ User %d: Success 200 (This should have taken 3s)\n", userID)
				}
				resp.Body.Close()
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("\n✅ Finished in %v\n", time.Since(start))
}