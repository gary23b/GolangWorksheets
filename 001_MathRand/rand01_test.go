package main

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_randomInt(t *testing.T) {
	bins := make([][]int64, 8)
	for i := range bins {
		bins[i] = make([]int64, 256)
	}

	loopCount := 1000000

	for range loopCount {
		x := RandomInt()

		for binIndex := range 8 {
			y := x & 0xFF
			bins[binIndex][y] += 1
			x = x >> 8
		}
	}

	binMax := int64(float64(loopCount) / 256.0 * 1.1)
	binMin := int64(float64(loopCount) / 256.0 * 0.9)

	binMaxLastByte := int64(float64(loopCount) / 128.0 * 1.1)
	binMinLastByte := int64(float64(loopCount) / 128.0 * 0.9)

	for b1, ByteBins := range bins {
		for b2, byteBin := range ByteBins {
			if b1 != 7 {
				require.Less(t, byteBin, binMax, "Byte %d, value %d count is too high", b1, b2)
				require.Greater(t, byteBin, binMin, "Byte %d, value %d count is too low", b1, b2)
			} else {
				// The top byte has the sign bit with is always zero
				if b2 >= 128 {
					require.Equal(t, int64(0), byteBin, "the sign bit should always be zero")
				} else {
					require.Less(t, byteBin, binMaxLastByte, "Byte %d, value %d count is too high", b1, b2)
					require.Greater(t, byteBin, binMinLastByte, "Byte %d, value %d count is too low", b1, b2)
				}
			}
		}
	}
}

func TestRandomFloat64(t *testing.T) {
	loopCount := 1000000
	gotMap := make(map[float64]int, loopCount)
	sum := 0.0

	for range loopCount {
		x := RandomFloat64()
		if gotMap[x] != 0 {
			t.Error("Saw the same value twice")
		}
		gotMap[x] = 1
		sum += x
	}

	average := sum / float64(loopCount)
	require.Less(t, average, .51, "Average is too high")
	require.Greater(t, average, .49, "Average is too high")
}

func TestRandomInt32Between(t *testing.T) {

	for range 1000 {
		x := RandomInt32Between(5, 5)
		require.Equal(t, int32(5), x)
	}

	gotMap := make(map[int32]bool)
	for range 1000 {
		x := RandomInt32Between(-10, 5)
		require.GreaterOrEqual(t, x, int32(-10))
		require.LessOrEqual(t, x, int32(5))
		gotMap[x] = true
	}

	for i := -10; i <= 5; i++ {
		require.True(t, gotMap[int32(i)])
	}
}

func TestShuffleIntsInPlace(t *testing.T) {

	getStartArray := func(n int) []int {
		ret := make([]int, 0, 10)
		for i := range n {
			ret = append(ret, i)
		}
		return ret
	}

	for range 10 {
		n := rand.Intn(11) + 10
		start := getStartArray(n)
		gotMap := make(map[int]int)
		for range 100 {
			a := getStartArray(n)
			ShuffleIntsInPlace(a)

			mapKey := 0
			for i, val := range a {
				mapKey += int(math.Pow(10, float64(i))) * val
			}
			gotMap[mapKey] = gotMap[mapKey] + 1
			require.NotEqual(t, a, start, "list not changed")
			require.ElementsMatch(t, start, a, "list doesn't have the same elements as it started with")

		}
		for _, val := range gotMap {
			require.Less(t, val, 3, "Getting the same shuffle to many times")
		}
	}
}

func TestNewRandomSource(t *testing.T) {
	s := NewRandomSource()
	s.Seed(10)
	got := s.RandomBytes(10)
	want := []byte{0x83, 0x47, 0x2e, 0xda, 0x6e, 0xb4, 0x75, 0x90, 0x6a, 0xee}
	require.Equal(t, want, got)

	s.Seed(67)
	got = s.RandomBytes(5)
	want = []byte{0xa5, 0x6c, 0x88, 0x3f, 0xf4}
	require.Equal(t, want, got)

	s.Seed(100)
	require.Equal(t, -0.2850021332665101, s.Float64Normal())
	require.Equal(t, -0.19284757603646052, s.Float64Normal())
	require.Equal(t, 0.16092035055548803, s.Float64Normal())
	require.Equal(t, 0.6025740677484148, s.Float64Normal())
}
