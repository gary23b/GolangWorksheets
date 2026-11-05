/*
Master the built int package: math/random
https://pkg.go.dev/math/rand
*/
package main

import (
	"crypto/sha256"
	"io"
)

// To make cryptographically secure hashes we need to use the crypto packages
// a very common one is sha256
// A hash algorithm takes bytes and produces a fixed-size hash value.

// The easiest usage is to feed it a byte slice and get the hash
func HashByteSlice(in []byte) []byte {
	// !->
	ret := sha256.Sum256(in)
	return ret[:]
	// <-!
}

// We can also stream data into it because the hash.HashReader interface implements the io.Writer interface.
// That means we can hash data as it comes in, without needing to have it all in memory at once.
func HashReader(reader io.Reader) []byte {
	// !->
	h := sha256.New()
	io.Copy(h, reader)

	return h.Sum(nil)
	// <-!
}
