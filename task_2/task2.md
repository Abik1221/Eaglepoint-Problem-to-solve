# Async Data Fetcher with Retry

## Search Log & References

I researched standard patterns for handling retries and timeouts in Go to ensure the solution was robust.

*   **Search 1:** "golang retry pattern best practice"
    *   *Goal:* Find the idiomatic way to retry an operation.
    *   *Result:* Found that a simple `for` loop is often preferred over complex libraries for simple cases.
    *   *URL:* [Go by Example: Timeouts](https://gobyexample.com/timeouts) (Used for understanding `time.Sleep` vs `select` with channels).

*   **Search 2:** "golang random number generator seed"
    *   *Goal:* Ensure the mock API behaves randomly each run.
    *   *Result:* Need to call `rand.Seed(time.Now().UnixNano())` once at startup, otherwise `rand.Intn` returns the same sequence.
    *   *URL:* [pkg.go.dev/math/rand](https://pkg.go.dev/math/rand)

*   **Search 3:** "golang sleep duration"
    *   *Goal:* Correct syntax for sleeping.
    *   *Result:* `time.Sleep(1 * time.Second)`. It's important to multiply by `time.Second` because the argument is a `Duration` (nanoseconds), not just an integer.
    *   *URL:* [pkg.go.dev/time#Sleep](https://pkg.go.dev/time#Sleep)



## Thought Process

### Why this approach?
I considered two approaches for the retry mechanism:
1.  **Recursive Approach**: Calling the function again inside itself.
    *   *Pros:* Elegant for some algorithms.
    *   *Cons:* Can lead to stack overflow (unlikely here, but bad practice) and harder to reason about state (like "retries left").
2.  **Iterative Approach**: A simple `for` loop.
    *   *Pros:* Extremely clear. "Do this X times". Easy to insert a `Sleep` between iterations.
    *   *Decision:* **Selected Approach 2**. It is the most readable and "Go-like" way to solve this.

### Alternatives Considered
*   **Exponential Backoff**: Instead of waiting 1s, wait 1s, then 2s, then 4s.
    *   *Decision:* The requirements specifically asked for "Waits 1 second between retries", so I stuck to the constant backoff to meet the spec exactly.



## Step-by-Step Solution

1.  **Mocking the API**:
    *   I created `callAPI` to simulate a real network request.
    *   Used `time.Sleep(500 * time.Millisecond)` to simulate latency.
    *   Used `rand.Float64() > 0.6` to create a 60% failure rate, ensuring the retry logic would actually be triggered during testing.
2.  **The Retry Loop**:
    *   I wrapped the call in a loop: `for i := 1; i <= retries; i++`.
    *   Inside the loop, I call the API.
    *   **Success Case**: If `err == nil`, I return the response immediately. No need to keep looping.
    *   **Failure Case**: If `err != nil`, I save the error to a variable `err` (so I can return it later if all attempts fail).
3.  **The Wait**:
    *   I added a check `if i < retries`. This is crucial. We don't want to sleep after the *last* attempt, because we're about to give up anyway.
    *   Used `time.Sleep(time.Second)` to pause execution.
4.  **Final Return**:
    *   If the loop finishes without returning, it means all attempts failed. I return a formatted error message: `fmt.Errorf("gave up after %d tries: %s", ...)` so the user knows exactly what happened.
