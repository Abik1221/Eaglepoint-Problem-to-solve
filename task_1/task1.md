#  Smart Text Analyzer 

## Search Log & References

Before starting, I researched the most efficient ways to handle string manipulation in Go to ensure performance and readability.

*   **Search 1:** "golang split string by whitespace vs split"
    *   *Goal:* Find the best way to handle multiple spaces and newlines.
    *   *Result:* Found `strings.Fields` in the official docs. It's better than `strings.Split` because it ignores empty results from consecutive spaces.
    *   *URL:* [pkg.go.dev/strings#Fields](https://pkg.go.dev/strings#Fields)

*   **Search 2:** "golang lowercase string efficient"
    *   *Goal:* Ensure case-insensitive comparison without allocating too much memory.
    *   *Result:* `strings.ToLower` is the standard approach.
    *   *URL:* [pkg.go.dev/strings#ToLower](https://pkg.go.dev/strings#ToLower)

*   **Search 3:** "golang map iteration order"
    *   *Goal:* Check if I needed to sort keys for the output.
    *   *Result:* Go map iteration is random. Since the requirement didn't specify sorted output, I left it as is, but noted that `json.Marshal` sorts keys alphabetically by default, which is a nice bonus.


## Thought Process

### Why this approach?
I considered two main approaches:
1.  **Regex Approach**: Using `regexp` to find words.
    *   *Pros:* Very flexible for complex patterns.
    *   *Cons:* Slower and harder to read. Overkill for simple text.
2.  **Standard Library Approach**: Using `strings.Fields`.
    *   *Pros:* O(N) complexity, built-in, handles all whitespace (tabs, newlines) automatically.
    *   *Decision:* **Selected Approach 2**. It's the "Go way" - simple, readable, and efficient.

### Alternatives Considered
I initially thought about calculating the "Average Word Length" by iterating through the map of frequencies.
*   *Correction:* That would be inaccurate because it would count unique words, not total words. I realized I needed to sum the length of *every* word as I processed the input list, not the unique map.


## Step-by-Step Solution

1.  **Setup**: I created the `TextAnalysis` struct with JSON tags. This ensures the output matches the requested format exactly (snake_case keys).
2.  **Input Cleaning**:
    *   Called `strings.Fields(text)` to get a clean slice of words.
    *   Added a check `if len(words) == 0` to handle empty input immediately.
3.  **The Main Loop**:
    *   I iterated through the slice of words.
    *   **Normalization**: Converted each word to lowercase immediately so "The" and "the" count as the same word.
    *   **Frequency**: Used `analysis.WordFrequency[lowerWord]++`. Go maps handle zero-values automatically, so I didn't need to check if the key existed first.
    *   **Length Tracking**: Added the length of the current word to `totalLetters`.
4.  **Longest Word Logic**:
    *   If `len(word) > maxWordLength`: Found a new longest word. I cleared the `LongestWords` slice and started fresh.
    *   If `len(word) == maxWordLength`: It's a tie. I checked if the word was already in the list (to avoid duplicates like `["fox", "fox"]`) and appended it.
5.  **Final Calculation**:
    *   Calculated average: `totalLetters / wordCount`.
    *   Rounding: Used `int(val * 100) / 100.0` to truncate to 2 decimal places as requested.


check for the flow chart image in this folder 