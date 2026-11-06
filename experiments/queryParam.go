/*
Look at a way to get query parameters without manually parsing or an external library
https://pkg.go.dev/net/http#Request.URL.Query
*/
package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type parser struct {
	values             url.Values
	errors             []error
	allowUnknownFields bool
	usedValues         map[string]bool
}

func New(r *http.Request) *parser {
	return &parser{
		values:     r.URL.Query(),
		usedValues: make(map[string]bool),
	}
}

func (p *parser) AllowUnknownFields() *parser {
	p.allowUnknownFields = true
	return p
}

func (p *parser) Err() error {
	if !p.allowUnknownFields {
		for key := range p.values {
			if _, ok := p.usedValues[key]; !ok {
				p.errors = append(p.errors, fmt.Errorf("unknown query parameter: %s", key))
			}
		}
	}

	return errors.Join(p.errors...)
}

// ============================================================================
// INT BUILDER
// ============================================================================

type intBuilder struct {
	p          *parser
	key        string
	min        int
	hasMin     bool
	max        int
	hasMax     bool
	required   bool
	defaultVal int
	hasDefault bool
}

func (p *parser) Int(key string) *intBuilder { return &intBuilder{p: p, key: key} }

func (b *intBuilder) Min(m int) *intBuilder       { b.min = m; b.hasMin = true; return b }
func (b *intBuilder) Max(m int) *intBuilder       { b.max = m; b.hasMax = true; return b }
func (b *intBuilder) Required() *intBuilder       { b.required = true; return b }
func (b *intBuilder) Default(def int) *intBuilder { b.defaultVal = def; b.hasDefault = true; return b }

func (b *intBuilder) Get() int {
	// Guard rails against developer configuration error
	if b.required && b.hasDefault {
		panic(fmt.Sprintf("key '%s': cannot set both required and default", b.key))
	}
	if b.hasMin && b.hasMax && b.min > b.max {
		panic(fmt.Sprintf("key '%s': min (%d) cannot be greater than max (%d)", b.key, b.min, b.max))
	}
	if b.hasDefault {
		if b.hasMin && b.defaultVal < b.min {
			panic(fmt.Sprintf("key '%s': default (%d) cannot be less than min (%d)", b.key, b.defaultVal, b.min))
		}
		if b.hasMax && b.defaultVal > b.max {
			panic(fmt.Sprintf("key '%s': default (%d) cannot be greater than max (%d)", b.key, b.defaultVal, b.max))
		}
	}

	valStr := b.p.values[b.key]
	if len(valStr) == 0 {
		if b.required {
			b.p.errors = append(b.p.errors, fmt.Errorf("%s is required", b.key))
		}
		return b.defaultVal
	}
	if len(valStr) > 1 {
		b.p.errors = append(b.p.errors, fmt.Errorf("%s must have only one value", b.key))
		return b.defaultVal
	}
	b.p.usedValues[b.key] = true

	i, err := strconv.Atoi(valStr[0])
	if err != nil {
		b.p.errors = append(b.p.errors, fmt.Errorf("%s must be a valid integer", b.key))
		return b.defaultVal
	}

	if b.hasMin && i < b.min {
		b.p.errors = append(b.p.errors, fmt.Errorf("%s must be at least %d", b.key, b.min))
		return b.defaultVal
	}
	if b.hasMax && i > b.max {
		b.p.errors = append(b.p.errors, fmt.Errorf("%s must be at most %d", b.key, b.max))
		return b.defaultVal
	}
	return i
}

// ============================================================================
// STRING BUILDER
// ============================================================================

type stringBuilder struct {
	p          *parser
	key        string
	required   bool
	defaultVal string
	hasDefault bool
}

func (p *parser) String(key string) *stringBuilder { return &stringBuilder{p: p, key: key} }
func (b *stringBuilder) Required() *stringBuilder  { b.required = true; return b }
func (b *stringBuilder) Default(def string) *stringBuilder {
	b.defaultVal = def
	b.hasDefault = true
	return b
}

func (b *stringBuilder) Get() string {
	if b.required && b.hasDefault {
		panic(fmt.Sprintf("key '%s': cannot set both required and default", b.key))
	}

	valStr := b.p.values[b.key]
	if len(valStr) == 0 {
		if b.required {
			b.p.errors = append(b.p.errors, fmt.Errorf("%s is required", b.key))
		}
		return b.defaultVal
	}
	if len(valStr) > 1 {
		b.p.errors = append(b.p.errors, fmt.Errorf("%s must have only one value", b.key))
		return b.defaultVal
	}
	b.p.usedValues[b.key] = true

	return valStr[0]
}

// ============================================================================
// BOOL BUILDER
// ============================================================================

type boolBuilder struct {
	p          *parser
	key        string
	defaultVal bool
}

func (p *parser) Bool(key string) *boolBuilder       { return &boolBuilder{p: p, key: key} }
func (b *boolBuilder) Default(def bool) *boolBuilder { b.defaultVal = def; return b }

func (b *boolBuilder) Get() bool {
	valStr := b.p.values[b.key]
	if len(valStr) == 0 {
		return b.defaultVal
	}
	if len(valStr) > 1 {
		b.p.errors = append(b.p.errors, fmt.Errorf("%s must have only one value", b.key))
		return b.defaultVal
	}
	b.p.usedValues[b.key] = true

	// Accepts true, false, 1, 0, t, f
	val, err := strconv.ParseBool(valStr[0])
	if err != nil {
		b.p.errors = append(b.p.errors, fmt.Errorf("%s must be a valid boolean", b.key))
		return b.defaultVal
	}
	return val
}

// ============================================================================
// TIME BUILDER
// ============================================================================

type timeBuilder struct {
	p          *parser
	key        string
	layout     string
	required   bool
	defaultVal time.Time
	hasDefault bool
}

// Time accepts a key and a Go layout string (e.g., time.RFC3339 or "2006-01-02")
func (p *parser) Time(key string, layout string) *timeBuilder {
	return &timeBuilder{p: p, key: key, layout: layout}
}
func (b *timeBuilder) Required() *timeBuilder { b.required = true; return b }
func (b *timeBuilder) Default(def time.Time) *timeBuilder {
	b.defaultVal = def
	b.hasDefault = true
	return b
}

func (b *timeBuilder) Get() time.Time {
	if b.required && b.hasDefault {
		panic(fmt.Sprintf("key '%s': cannot set both required and default", b.key))
	}

	valStr := b.p.values[b.key]
	if len(valStr) == 0 {
		if b.required {
			b.p.errors = append(b.p.errors, fmt.Errorf("%s is required", b.key))
		}
		return b.defaultVal
	}

	if len(valStr) > 1 {
		b.p.errors = append(b.p.errors, fmt.Errorf("%s must have only one value", b.key))
		return b.defaultVal
	}
	b.p.usedValues[b.key] = true

	t, err := time.Parse(b.layout, valStr[0])
	if err != nil {
		b.p.errors = append(b.p.errors, fmt.Errorf("%s must match format %s", b.key, b.layout))
		return b.defaultVal
	}
	return t
}
