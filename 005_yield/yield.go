/*
Master a usage of yield
https://go.dev/blog/range-functions
https://gobyexample.com/range-over-iterators
*/
package main

import (
	"io"
	"iter"
)

// In golang version 1.23 they added support for iterators and the yield pattern.
// a basic example is generating Fibonacci numbers using yield.
func GenFib() func(yield func(int) bool) {
	return func(yield func(int) bool) {
		a, b := 0, 1
		for {
			if !yield(a) {
				return
			}
			a, b = b, a+b
		}
	}
}

// when you pass function with the signature func(yield func(T) bool) to a for range loop,
// golang now allows you to iterate over the generated values as if they were a channel.
func PrintFib() {
	loopCount := 0
	for v := range GenFib() {
		println(v)

		loopCount++
		if loopCount >= 20 {
			break
		}
	}
}

// They also created a helper type called iter.Seq and iter.Seq2 that simplifies the function signature for yielding data.
// It makes the definition shorter, but hides the full function signature and the control flow of yielding values.
/*
type Seq[V any] func(yield func(V) bool)
type Seq2[K, V any] func(yield func(K, V) bool)
*/

// To use iter.Seq we just replace
func GenFib2() iter.Seq[int] {
	return func(yield func(int) bool) {
		a, b := 0, 1
		for {
			if !yield(a) {
				return
			}
			a, b = b, a+b
		}
	}
}

// Can you make one that outputs 0,1,2,...,10,0,1,2,... called RampWaveform
// !->
func RampWaveform() iter.Seq[int] {
	return func(yield func(int) bool) {
		i := 0
		for {
			if !yield(i % 11) {
				return
			}
			i++
		}
	}
}

// <-!

// We can also rely on the iterator function to stop itself:
func GenNumbers(start, stop int) iter.Seq[int] {
	return func(yield func(int) bool) {
		for i := start; i <= stop; i++ {
			if !yield(i) {
				return
			}
		}
	}
}

// Then since we can make function that take a sequence as an input and Collect it all
func Collect[T any](seq iter.Seq[T]) []T {
	// !->
	var result []T
	for v := range seq {
		result = append(result, v)
	}
	return result
	// <-!
}

// this is the same thing as slices.Collect in the standard library

// Then since we can make function that take a sequence as an input and also output a sequence,
// we can build pipelines of data processing without ever storing all the intermediate results in memory.
func FilterEvenNumbers(seq iter.Seq[int]) iter.Seq[int] {
	// !->
	return func(yield func(int) bool) {
		for v := range seq {
			if v%2 == 0 {
				if !yield(v) {
					return
				}
			}
		}
	}
	// <-!
}

// I have started to see that when working with a database that can give you huge amounts of data,
// you need to start thinking about streaming the data. But to do that without bringing implementation details into your business logic,
// you can use a new concept called yield.
// It makes it so you can loop over the data as it becomes available, without loading everything into memory at once.

// Say you have a database query that functions sort of like this:

var databaseData = []int{}
var dbQueryErr error
var dbNextErr error
var dbScanErr error

type Rows struct {
	index int
}

func RunQuery() (*Rows, error) {
	return &Rows{
		index: -1,
	}, dbQueryErr
}

func (s *Rows) Close() error {
	// In a real database, you would release resources here.
	return nil
}

func (s *Rows) Next() bool {
	if dbNextErr != nil {
		return false
	}
	s.index++
	if s.index < 0 || s.index >= len(databaseData) {
		return false
	}
	return true
}

func (s *Rows) Scan() (int, error) {
	if s.index >= len(databaseData) {
		return 0, io.EOF
	}
	if dbScanErr != nil {
		return 0, dbScanErr
	}
	return databaseData[s.index], nil
}

// So, we can use the Rows type to iterate over query results in a memory-efficient way.
func GetLotsOfData() iter.Seq2[int, error] {
	// !->
	return func(yield func(int, error) bool) {
		rows, err := RunQuery()
		if err != nil {
			yield(0, err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			data, err := rows.Scan()
			if err != nil {
				yield(0, err)
				return
			}
			if !yield(data, nil) {
				return
			}
		}
	}
	// <-!
}
