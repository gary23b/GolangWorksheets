# Golang Worksheets

A collection of focused worksheets demonstrating how to use Go's built-in and standard library packages effectively.

The idea is that each worksheet can be finished pretty quickly, especially, if you know what you are doing.
But you should come back to them often so that you can quickly recall and use these important concepts.

## How To Reset Worksheets

Run `go run .` in the top level directory to reset the solutions.

## Project Structure

Each exercise directory contains:
- `*.go` - Implementation file with key patterns and examples
- `*_test.go` - Test file demonstrating usage and validating behavior

## Exercises

### 001_MathRand
Master the `math/rand` package for generating random numbers.

### 002_Sha256
Master the `crypto/sha256` package for secure hashing.
- Direct byte slice hashing with `Sum256()`
- Streaming hash computation from readers

### 003_bcrypt
Master the `golang.org/x/crypto/bcrypt` package for secure password hashing.
- Password hashing with `GenerateFromPassword()`
- Password verification with `CompareHashAndPassword()`
- Interface-driven design for testability and flexibility

## Learning Approach

Each worksheet focuses on:
1. **Essential use cases** - Only the core patterns most developers need
2. **Production-ready design** - Following Go idioms and best practices
3. **Clear examples** - Commented code showing exactly how to use each feature
4. **Comprehensive tests** - Validating behavior and demonstrating usage

## Running Tests

```bash
go test ./...
```

## Notes

- Exercises are self-contained; each can be studied independently
- Tests use `testify/require` for clear, readable assertions
- Code includes fill-in-the-blank sections (marked with `// !->` and `// <-!`) for learning purposes
