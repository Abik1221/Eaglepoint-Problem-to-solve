# Eagle Point AI test for coding challenge code submision doc

**Author:** Nahom Keneni  
**Project:** Coding Challenge Solutions

## Overview
This repository contains my solutions for the Eagle Point AI entrance test. I've implemented three distinct tasks using **Golang**, focusing on clean, readable, and efficient code. Each task is isolated in its own folder with dedicated documentation explaining the logic and approach.

## Project Structure

### [Task 1: Smart Text Analyzer](./task_1)
A function that analyzes text to extract metrics like word count, average length, and frequency.
- **Code:** `task_1/main.go`
- **Documentation:** [Read the approach](./task_1/task1.md)

### [Task 2: Async Data Fetcher](./task_2)
An async data fetcher with robust retry logic to handle network failures gracefully.
- **Code:** `task_2/main.go`
- **Documentation:** [Read the approach](./task_2/task2.md)

### [Task 3: Rate Limiter](./task_3)
A custom rate limiter implementation that restricts user requests within a specific time window.
- **Code:** `task_3/main.go`
- **Documentation:** [Read the approach](./task_3/task3.md)

---

## Technology & Resources

### Language
- **Golang** (Go): Chosen for its simplicity, strong standard library, and excellent support for concurrency (goroutines/channels), which was perfect for the async and rate-limiting tasks.

### Packages Used
I relied almost exclusively on the **Go Standard Library** to keep the solutions lightweight and dependency-free:
- `fmt`: For formatting I/O.
- `strings`: For efficient text manipulation (Task 1).
- `time`: For handling timeouts, sleep durations, and ticker logic (Task 2 & 3).
- `sync`: For mutexes to ensure thread safety in the rate limiter (Task 3).
- `math/rand`: For simulating API failures and random data.
- `encoding/json`: For formatting output.

### References & Inspiration
To ensure I was following best practices, I consulted:
- **[Effective Go](https://go.dev/doc/effective_go)**: For idiomatic Go patterns.
- **[Go Documentation](https://pkg.go.dev/)**: specifically the `sync` and `time` package docs.
- **[Go by Example](https://gobyexample.com/)**: A great quick reference for syntax.

---

## How to Run

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Abik1221/Eaglepoint-Problem-to-solve.git
   cd Eaglepoint-Problem-to-solve
   ```

2. **Run a specific task:**
   Navigate to the folder and run the `main.go` file.

   **Task 1:**
   ```bash
   cd task_1
   go run main.go
   ```

   **Task 2:**
   ```bash
   cd task_2
   go run main.go
   ```

   **Task 3:**
   ```bash
   cd task_3
   go run main.go
   ```


Thank you 
Nahom Keneni
nahomkeneni4@gmail.com
