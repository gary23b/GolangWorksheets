package main

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueryParam(t *testing.T) {
	r := httptest.NewRequest("GET", "/?q=value&b=true", nil)

	parser := New(r).AllowUnknownFields()

	// Clean, readable one-liners
	limit := parser.Int("limit").Min(1).Max(100).Default(20).Get()
	query := parser.String("q").Required().Get()
	debug := parser.Bool("debug").Default(false).Get()
	since := parser.Time("since", "2006-01-02").Default(time.Now().AddDate(-1, 0, 0)).Get()

	// Handle all processing issues in one block
	err := parser.Err()
	require.NoError(t, err)

	require.Equal(t, 20, limit)
	require.Equal(t, "value", query)
	require.Equal(t, false, debug)
	require.WithinDuration(t, time.Now().AddDate(-1, 0, 0), since, time.Second)

	// Use your safe values...
	_, _, _, _ = limit, query, debug, since
}
