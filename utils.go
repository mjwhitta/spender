package spender

import (
	"strconv"
	"strings"

	hl "github.com/mjwhitta/hilighter"
)

func colorize(
	s string,
	x int,
	y int,
	maxX int,
	maxY int,
	expanded bool,
) string {
	// Totals in red, unless expanded
	if (x == maxX) || (y == maxY) {
		if expanded {
			return hl.Yellow(s)
		}

		return hl.Red(s)
	}

	// Header is blue
	if x == 0 {
		return hl.Blue(s)
	}

	// All other cells are cyan, unless expanded
	if expanded {
		return hl.Yellow(s)
	}

	return hl.Cyan(s)
}

func floatStr(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func less(a []string) func(i int, j int) bool {
	return func(i int, j int) bool {
		return strless(a[i], a[j])
	}
}

func normalize(s string) string {
	s = strings.ReplaceAll(s, "|", "")
	s = trimBegin.ReplaceAllString(s, "")
	s = trimEnd.ReplaceAllString(s, "")

	return strings.TrimSpace(s)
}

func strless(l string, r string) bool {
	return strings.ToLower(l) < strings.ToLower(r)
}
