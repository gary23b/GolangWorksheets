package main

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCountUpTo(t *testing.T) {
	d1 := []int{}
	count := 0
	for i := range GenFib() {
		d1 = append(d1, i)
		count++
		if count >= 8 {
			break
		}
	}
	require.Equal(t, []int{0, 1, 1, 2, 3, 5, 8, 13}, d1)
}

func TestRampWaveform(t *testing.T) {
	d1 := []int{}
	count := 0
	for i := range RampWaveform() {
		d1 = append(d1, i)
		count++
		if count >= 16 {
			break
		}
	}
	require.Equal(t, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 0, 1, 2, 3, 4}, d1)

	// break early test
	d2 := []int{}
	count = 0
	for i := range RampWaveform() {
		d2 = append(d2, i)
		count++
		if count >= 1 {
			break
		}
	}
	require.Equal(t, []int{0}, d2)
}

func TestGenNumbers(t *testing.T) {
	got := Collect(GenNumbers(1, 5))
	require.Equal(t, []int{1, 2, 3, 4, 5}, got)

	got = slices.Collect(GenNumbers(1, 5))
	require.Equal(t, []int{1, 2, 3, 4, 5}, got)
}

func TestFilterEvenNumbers(t *testing.T) {
	got := Collect(FilterEvenNumbers(GenNumbers(1, 10)))
	require.Equal(t, []int{2, 4, 6, 8, 10}, got)
}

func TestGetLotsOfData(t *testing.T) {
	resetErr := func() {
		dbQueryErr = nil
		dbNextErr = nil
		dbScanErr = nil
	}
	fillDatabase := func(length int) []int {
		databaseData = make([]int, length)
		want := make([]int, length)
		for i := range length {
			databaseData[i] = rand.Int()
			want[i] = databaseData[i]
		}
		return want
	}

	resetErr()
	want := fillDatabase(10)

	data := []int{}
	for v, err := range GetLotsOfData() {
		require.NoError(t, err)
		data = append(data, v)
	}
	require.Equal(t, want, data)

	// run query error
	resetErr()
	dbQueryErr = fmt.Errorf("query error")
	for v, err := range GetLotsOfData() {
		require.Error(t, err)
		require.Equal(t, 0, v)
	}

	// next error
	resetErr()
	want = fillDatabase(10)
	count := 0
	for v, err := range GetLotsOfData() {
		require.NoError(t, err)
		require.Equal(t, want[count], v)
		dbNextErr = fmt.Errorf("next error")
		count++
	}
	require.Equal(t, 1, count)

	// scan error
	resetErr()
	want = fillDatabase(10)
	count = 0
	for v, err := range GetLotsOfData() {
		if count == 0 {
			require.NoError(t, err)
			require.Equal(t, want[count], v)
			dbScanErr = fmt.Errorf("scan error")
		}
		if count > 0 {
			require.Error(t, err)
		}
		count++
	}
	require.Equal(t, 2, count)

	// break early test
	resetErr()
	want = fillDatabase(10)
	for v, err := range GetLotsOfData() {
		require.NoError(t, err)
		data = append(data, v)
		require.Equal(t, want[0], v)
		break
	}
}
