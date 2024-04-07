package spender

import (
	"regexp"
	"strings"

	"github.com/mjwhitta/errors"
)

// Merchant is a struct to track spending per merchant.
type Merchant struct {
	Costs map[string]map[string]float64
	Name  string

	patterns []any
}

// Group will add a list of regex or string literal patterns to group
// merchant names into one Merchant.
func (m *Merchant) Group(patterns []string) error {
	var e error
	var ends bool
	var isRegex bool
	var r *regexp.Regexp
	var starts bool

	for _, pattern := range patterns {
		if pattern == "" {
			continue
		}

		starts = strings.HasPrefix(pattern, "/")
		ends = strings.HasSuffix(pattern, "/")
		ends = ends || strings.HasSuffix(pattern, "/i")
		isRegex = starts && ends

		if !isRegex {
			m.patterns = append(m.patterns, pattern)
			continue
		}

		if strings.HasSuffix(pattern, "/") {
			pattern = strings.TrimSuffix(pattern[1:], "/")
		} else if strings.HasSuffix(pattern, "/i") {
			pattern = "(?i)" + strings.TrimSuffix(pattern[1:], "/i")
		}

		if r, e = regexp.Compile(pattern); e != nil {
			return errors.Newf("invalid regex: %w", e)
		}

		m.patterns = append(m.patterns, r)
	}

	return nil
}

func (m *Merchant) init() {
	if m.Costs == nil {
		m.Costs = map[string]map[string]float64{}
	}
}

// Matches will return whether or not the provided merchant matches
// any of the group's patterns.
func (m *Merchant) Matches(name string) bool {
	for _, pattern := range m.patterns {
		switch pattern := pattern.(type) {
		case *regexp.Regexp:
			if pattern.MatchString(name) {
				return true
			}
		case string:
			if pattern == name {
				return true
			}
		}
	}

	return m.Name == name
}

// Purchase will track a new purchase for the Merchant.
func (m *Merchant) Purchase(label string, name string, cost float64) {
	m.init()

	if _, ok := m.Costs[name]; !ok {
		m.Costs[name] = map[string]float64{}
	}

	m.Costs[name][label] += cost
}

// Subtotals will return the total spent per submerchant, per label.
func (m *Merchant) Subtotals() map[string]map[string]float64 {
	if len(m.patterns) == 0 {
		return nil
	}

	return m.Costs
}

// Totals will return the total spent per label.
func (m *Merchant) Totals() map[string]float64 {
	var totals map[string]float64 = map[string]float64{}

	for _, costs := range m.Costs {
		for label, cost := range costs {
			totals[label] += cost
		}
	}

	return totals
}
