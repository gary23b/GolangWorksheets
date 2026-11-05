# Learning Model

This project uses a focused, practical learning model for mastering Go built-in packages.

## Philosophy

Rather than exhaustive coverage of every feature, each worksheet covers:
- **The core use cases** - What 90% of developers actually use
- **The "why"** - Why this package matters and when to reach for it
- **Best practices** - How professionals structure code using these packages

## Structure of Each Worksheet

### 1. Clear Purpose
Each file starts with a comment explaining what package is being mastered and linking to official documentation.

### 2. Focused Scope
Two to three essential functions/patterns, nothing more. Better to deeply understand the essentials than superficially cover everything.

### 3. Interface-Driven Design
When applicable, abstractions are introduced early:
- Makes code testable
- Shows professional Go patterns
- Demonstrates how to swap implementations

### 4. Comprehensive Tests
Tests serve double duty:
- Validate the implementation works
- Show how to use the package (documentation through code)

## Example: SHA256 Worksheet

The `002_Sha256` worksheet covers exactly two use cases:
1. **Hashing a byte slice** - `sha256.Sum256()`
2. **Streaming hash computation** - `sha256.New()` + `io.Copy()`

Why these two? Because they cover 99% of real-world hashing needs. Everything else (incremental writes, hex encoding, etc.) is just variation on these patterns.

## Example: bcrypt Worksheet

The `003_bcrypt` worksheet demonstrates:
1. **Interface design** - Abstracting the crypto implementation
2. **Password hashing** - Secure password storage
3. **Password verification** - Checking plaintext against stored hash

The `HashAndPasswordMatch()` helper wraps the error-based API into a boolean, showing how to adapt package APIs to your domain needs.

## Quality Indicators

A worksheet is well-designed when someone can:
1. Read it top-to-bottom in 10 minutes
2. Understand all the core patterns
3. Look at the tests and see exactly how to use it
4. Immediately apply it in real code
