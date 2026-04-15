// Copyright 2026 Phillip Cloud
// Licensed under the Apache License, Version 2.0

//go:build gui

package gui

import (
	"fmt"
	"strconv"
	"time"
)

// centsStr formats a nullable cents value as a dollar amount (e.g. "$12.50").
func centsStr(c *int64) string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("$%.2f", float64(*c)/100)
}

// timeStr formats a nullable time as "YYYY-MM-DD".
func timeStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("2006-01-02")
}

// parseDateOr parses a "YYYY-MM-DD" string, returning nil on failure.
func parseDateOr(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil
	}
	return &t
}

// parseIntOr parses s as an int, returning def on failure.
func parseIntOr(s string, def int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return v
}

// parseFloatOr parses s as a float64, returning def on failure.
func parseFloatOr(s string, def float64) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return def
	}
	return v
}
