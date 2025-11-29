package main

import (
	"fmt"
	"sync"
	"time"
)

// this struct tracks the request history for a single user
type UserLimit struct {
	requests []time.Time // a list of timestamps when users are made requests
	mu       sync.Mutex  // caches this specific user's data safe from concurrent access
}

// limiter manages rate limits for all users
type Limiter struct {
	users   map[string]*UserLimit // map of user_ID
	mu      sync.RWMutex          // protects the users map itself
	maxReqs int                   // how many requests are allowed
	window  time.Duration         // in what time frame does it made
}

// new limiter creates a fresh rate limiter
func NewLimiter(maxReqs int, window time.Duration) *Limiter {
	return &Limiter{
		users:   make(map[string]*UserLimit),
		maxReqs: maxReqs,
		window:  window,
	}
}

// allow checks if a user can make a request right now
func (l *Limiter) Allow(userID string) bool {
	// first, we need to get the user's record. then
	// we use a read/write lock here because we might need to write to the map if the user is new.
	l.mu.Lock()
	if l.users[userID] == nil {
		l.users[userID] = &UserLimit{
			requests: []time.Time{},
		}
	}
	user := l.users[userID]
	l.mu.Unlock()

	// now we lock just this user's data to check their limits.
	// this is great because checking User A doesn't block User B.
	user.mu.Lock()
	defer user.mu.Unlock()

	now := time.Now()
	// calculate the "cutoff" time. Any request before this time is too old to matter.
	cutoff := now.Add(-l.window)

	// filter out the old requests. We only keep the ones inside the current window.
	valid := []time.Time{}
	for _, t := range user.requests {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	user.requests = valid

	// check if they've hit their limit
	if len(user.requests) >= l.maxReqs {
		return false // Too many requests! block them.
	}

	// if they're good, record this request and let them .
	user.requests = append(user.requests, now)
	return true
}

// getStats is a helper to see how many requests a user has made
func (l *Limiter) GetStats(userID string) (current int, limit int) {
	l.mu.RLock()
	user := l.users[userID]
	l.mu.RUnlock()

	if user == nil {
		return 0, l.maxReqs
	}

	user.mu.Lock()
	defer user.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-l.window)

	count := 0
	for _, t := range user.requests {
		if t.After(cutoff) {
			count++
		}
	}

	return count, l.maxReqs
}

func main() {
	// create a limiter: 5 requests every 60 seconds
	limiter := NewLimiter(5, 60*time.Second)

	fmt.Println("Rate Limiter Demo")
	fmt.Println("Limit: 5 requests per 60 seconds")

	testUser := "user_123"

	fmt.Println("Simulating rapid requests from user_123...")

	// Try to make 8 requests in a row (only 5 should work)
	for i := 1; i <= 8; i++ {
		allowed := limiter.Allow(testUser)
		current, limit := limiter.GetStats(testUser)

		if allowed {
			fmt.Printf("Request %d: Allowed (%d/%d)\n", i, current, limit)
		} else {
			fmt.Printf("Request %d: Blocked (%d/%d)\n", i, current, limit)
		}

		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("\nWaiting 3 seconds...")
	time.Sleep(3 * time.Second)

	fmt.Println("Testing with a different user (user_456)...")

	// new user should have a fresh limit
	testUser2 := "user_456"
	for i := 1; i <= 3; i++ {
		allowed := limiter.Allow(testUser2)
		current, limit := limiter.GetStats(testUser2)

		if allowed {
			fmt.Printf("Request %d:Allowed (%d/%d)\n", i, current, limit)
		} else {
			fmt.Printf("Request %d:Blocked (%d/%d)\n", i, current, limit)
		}
	}

	fmt.Println("\nChecking user_123 stats after waiting...")
	current, limit := limiter.GetStats(testUser)
	fmt.Printf("user_123: %d/%d requests in window\n", current, limit)
}
