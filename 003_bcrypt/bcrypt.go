/*
Master the built int package: x/crypto/bcrypt
https://pkg.go.dev/golang.org/x/crypto/bcrypt
*/
package main

import (
	"golang.org/x/crypto/bcrypt"
)

// To hash passwords securely, we can use the bcrypt package.
// It provides functions to generate and compare password hashes.

// A good way to use bcrypt is to create an interface around it
// so that it can be easily swapped out or mocked for testing.

type PasswordHasher interface {
	HashPassword(password []byte) ([]byte, error)
	HashAndPasswordMatch(hash, password []byte) bool
}

// create a struct that implements the PasswordHasher interface
type BcryptHasher struct{}

var _ PasswordHasher = BcryptHasher{}

// now we implement the PasswordHasher interface using bcrypt
// HashPassword hashes the given password using bcrypt.
func (s BcryptHasher) HashPassword(password []byte) ([]byte, error) {
	// !->
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	// <-!
}

// HashAndPasswordMatch compares a bcrypt hashed password with its possible plaintext equivalent.
func (s BcryptHasher) HashAndPasswordMatch(hash, password []byte) bool {
	// !->
	err := bcrypt.CompareHashAndPassword(hash, password)
	return err == nil
	// <-!
}
