/*
Master the built int package: math/random
https://pkg.go.dev/math/rand
*/
package main

import (
	"math/rand"
)

// To make random numbers in Go, we can use the math/rand package.
// As long as we don't need cryptographic security, math/rand is sufficient.

// Make a function that returns a random integer (non negative).
func RandomInt() int {
	// !->
	return rand.Int()
	// <-!
}

// make a function that return a random float64 [0.0,1.0)
func RandomFloat64() float64 {
	// !->
	return rand.Float64()
	// <-!
}

// make a function that return a int32 between a and b inclusive
func RandomInt32Between(a, b int32) int32 {
	// !->
	return rand.Int31n(b-a+1) + a
	// <-!
}

// make a function to shuffle a slice of integers
func ShuffleIntsInPlace(in []int) {
	// !->
	rand.Shuffle(len(in), func(i, j int) {
		in[i], in[j] = in[j], in[i]
	})
	// <-!
}

// We can also make a new random source that is independent of the default source.
// The source can only be used by one goroutine at a time.
// but that us to generate them faster without mutex contention.

// Create a new random source and that implements the interface below:
type RandomSource interface {
	Seed(in int64)
	Float64Normal() float64
	RandomBytes(len int) []byte
}

// !->
type rs struct {
	r *rand.Rand
}

func (s *rs) Seed(in int64) {
	s.r.Seed(in)
}

func (s *rs) Float64Normal() float64 {
	return s.r.NormFloat64()
}

func (s *rs) RandomBytes(len int) []byte {
	ret := make([]byte, len)
	n, err := s.r.Read(ret)
	if err != nil {
		panic(err)
	}
	if n != len {
		panic("n doesn't match len")
	}
	return ret
}

// <-!

func NewRandomSource() RandomSource {
	// !->
	ret := &rs{
		r: rand.New(rand.NewSource(99)),
	}
	return ret
	// <-!
}
