## Rate Limiter 

## Search Log & References

I researched different rate limiting algorithms and concurrency patterns in Go to ensure a thread-safe and accurate implementation.

*   **Search 1:** "golang rate limiter algorithms token bucket vs sliding window"
    *   *Goal:* Choose the best algorithm for "5 requests per 60 seconds".
    *   *Result:* "Token Bucket" is good for allowing bursts, but "Sliding Window Log" is more accurate for hard limits like "exactly X in the last Y seconds".
    *   *Decision:* I chose **Sliding Window Log** because it's precise and easy to implement for small-scale use cases.
    *   *URL:* [Rate Limiting Part 1](https://hechao.li/2018/06/25/Rate-Limiter-Part1/) (General concept).

*   **Search 2:** "golang mutex vs rwmutex"
    *   *Goal:* Optimize locking performance.
    *   *Result:* `sync.RWMutex` allows multiple readers (checking stats) at once, blocking only when writing (adding a new user). This is better for high-read scenarios.
    *   *URL:* [pkg.go.dev/sync#RWMutex](https://pkg.go.dev/sync#RWMutex)

*   **Search 3:** "golang time since vs now sub"
    *   *Goal:* Correctly calculate the time window.
    *   *Result:* `time.Now().Add(-window)` gives the "cutoff" time. Any timestamp before this is expired.
    *   *URL:* [pkg.go.dev/time#Time.Add](https://pkg.go.dev/time#Time.Add)



## Thought Process

### Why this approach?
I needed a way to track *exactly* when requests happened to enforce the limit accurately.
1.  **Data Structure**: I needed a map to store data for each user (`map[string]*UserLimit`).
2.  **Concurrency**: Since web servers handle many requests at once, I *had* to use locks.
    *   *Global Lock*: Protects the map itself (adding/finding users).
    *   *User Lock*: Protects the specific user's data. This is a huge optimization! It means User A doesn't block User B.

### Alternatives Considered
*   **Fixed Window Counter**: Just resetting a counter every minute (e.g., at 12:00, 12:01).
    *   *Flaw:* If I make 5 requests at 12:00:59 and 5 more at 12:01:01, I've made 10 requests in 2 seconds! This violates the spirit of the rate limit.
    *   *Decision:* Rejected in favor of Sliding Window.
*   **Token Bucket**: Adding tokens at a fixed rate.
    *   *Flaw:* Slightly more complex to implement "per 60 seconds" exactly without background timers.
    *   *Decision:* Rejected for this specific "N reqs / M time" requirement.



## Step-by-Step Solution

1.  **Struct Definition**:
    *   `UserLimit`: Holds a slice of `time.Time`. This is our "log" of requests.
    *   `Limiter`: Holds the map and configuration (`maxReqs`, `window`).
2.  **Initialization**:
    *   `NewLimiter` function to set up the map and rules.
3.  **The `Allow` Logic**:
    *   **Step 1 (Global Lock)**: `l.mu.Lock()`. Check if user exists. If not, create them. `l.mu.Unlock()`.
    *   **Step 2 (User Lock)**: `user.mu.Lock()`. Now we only block this specific user.
    *   **Step 3 (Cleanup)**: Calculate `cutoff = time.Now() - 60s`. Loop through `user.requests` and keep only the ones *after* the cutoff. This automatically removes old data, preventing memory leaks!
    *   **Step 4 (Check)**: If `len(requests) < maxReqs`, we're good. Append `time.Now()` and return `true`.
    *   **Step 5 (Block)**: If `len(requests) >= maxReqs`, return `false`.



## Why this solution is best

*   **Accuracy**: It enforces the limit precisely. No "bursts" at the turn of the minute.
*   **Thread Safety**: It uses fine-grained locking (per-user mutexes). This is critical for performance in a real server.
*   **Self-Cleaning**: We don't need a background "garbage collector" goroutine to clean up old timestamps. We do it lazily every time we check a user.
*   **Simplicity**: It achieves complex behavior (sliding window) with standard Go slices and time math.
