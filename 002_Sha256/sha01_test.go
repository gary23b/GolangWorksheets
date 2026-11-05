package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashByteSlice(t *testing.T) {

	gotBytes := HashByteSlice([]byte("hello world\n"))
	gotStr := fmt.Sprintf("%x", gotBytes)
	want := "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447"
	require.Equal(t, want, gotStr)

	gotBytes = HashByteSlice([]byte("1234567890-13444"))
	gotStr = fmt.Sprintf("%x", gotBytes)
	want = "f6ef92e32fdc3793683dd6ecc1876859d2dc6824adf32779c3829dca604e47b2"
	require.Equal(t, want, gotStr)
}

func TestHash(t *testing.T) {
	gotBytes := HashReader(strings.NewReader("hello world\n"))
	gotStr := fmt.Sprintf("%x", gotBytes)
	want := "a948904f2f0f479b8f8197694b30184b0d2ed1c1cd2a1ec0fb85d299a192a447"
	require.Equal(t, want, gotStr)
	gotBytes = HashReader(strings.NewReader("1234567890-13444"))
	gotStr = fmt.Sprintf("%x", gotBytes)
	want = "f6ef92e32fdc3793683dd6ecc1876859d2dc6824adf32779c3829dca604e47b2"
	require.Equal(t, want, gotStr)
}
