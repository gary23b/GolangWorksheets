package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBcryptHasher_HashPassword(t *testing.T) {
	b := BcryptHasher{}
	hash, err := b.HashPassword([]byte("abc"))
	require.NoError(t, err)
	require.NotNil(t, hash)
	require.True(t, b.HashAndPasswordMatch(hash, []byte("abc")))
	fmt.Println(string(hash))

	require.False(t, b.HashAndPasswordMatch(hash, []byte("wrongpassword")))

	prevHash := []byte("$2a$10$aUy1ofsNfAreM3FPtCU8duf6XcXdWF0MqkZsJawd99.6qqGBL8CqW")
	require.True(t, b.HashAndPasswordMatch(prevHash, []byte("abc")))
}
