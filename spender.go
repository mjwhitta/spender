package spender

import (
	"sort"
	"strings"

	"github.com/mjwhitta/errors"
	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/jq"
)

// Spender is a struct containing Groups and Merchants, and allows for
// tracking spending on a per Merchant basis.
type Spender struct {
	Exclude []string
	Expand  bool
	Include []string

	labels    map[string]set
	merchants []*Merchant
	results   [][]string
	tallied   bool
}

// New will return a pointer to a new Spender instance.
func New() *Spender {
	return &Spender{}
}

// CreateGroups will parse a JSON blob to create a list of Groups with
// associated patterns to match against Merchant names.
func (s *Spender) CreateGroups(json string) error {
	var e error
	var j *jq.JSON
	var m *Merchant
	var patterns []string

	if j, e = jq.New(json); e != nil {
		return errors.Newf("failed to parse JSON: %w", e)
	}

	for _, group := range j.GetKeys() {
		if patterns, e = j.MustGetStringArray(group); e != nil {
			return errors.Newf("unexpected JSON: %w", e)
		}

		m = &Merchant{Name: group}
		if e = m.Group(patterns); e != nil {
			return e
		}

		s.merchants = append(s.merchants, m)
	}

	return nil
}

func (s *Spender) filtered(name string) bool {
	for _, ex := range s.Exclude {
		if strings.EqualFold(ex, name) {
			return true
		}
	}

	if len(s.Include) == 0 {
		return strings.EqualFold(name, "ignore")
	}

	for _, in := range s.Include {
		if strings.EqualFold(in, name) {
			return false
		}
	}

	return true
}

func (s *Spender) init() {
	if s.labels == nil {
		s.labels = map[string]set{}
	}
}

// Purchase will track a new purchase (cost) for a specific merchant,
// organized by label, and applying any configured grouping.
func (s *Spender) Purchase(label string, merch string, cost float64) {
	var m *Merchant

	s.init()
	s.tallied = false

	s.labels[label] = set{}
	merch = normalize(merch)

	for _, m = range s.merchants {
		if m.Matches(merch) {
			m.Purchase(label, merch, cost)
			return
		}
	}

	m = &Merchant{Name: merch}
	s.merchants = append(s.merchants, m)
	m.Purchase(label, merch, cost)
}

// String will return a string representation of the Spender.
func (s *Spender) String() string {
	var fill string
	var filler string
	var out []string
	var tmp []string
	var widths map[int]int = map[int]int{}

	if !s.tallied {
		s.Tally()
	}

	// // Get max column widths
	for _, row := range s.results {
		for i, column := range row {
			if len(column) > widths[i] {
				widths[i] = len(column)
			}
		}
	}

	// Build table output
	for i, row := range s.results {
		filler = " "
		if (i == 1) || (i == len(s.results)-2) {
			filler = hl.Green("-")
		}

		tmp = []string{}

		for j, column := range row {
			fill = strings.Repeat(filler, widths[j]-len(column))

			column = colorize(
				column,
				i,
				j,
				len(s.results)-1,
				len(row)-1,
				strings.HasPrefix(row[0], " "),
			)

			if j == 0 {
				tmp = append(tmp, column+fill)
			} else {
				tmp = append(tmp, fill+column)
			}
		}

		filler = " | "
		if (i == 1) || (i == len(s.results)-2) {
			filler = "-+-"
		}

		out = append(out, strings.Join(tmp, hl.Green(filler)))
	}

	return strings.Join(out, "\n")
}

// Tally will finalize the financials and should be called before you
// display the spending results.
func (s *Spender) Tally() {
	var labels []string
	var names []string
	var row []string
	var subtotals map[string]map[string]float64
	var total float64
	var totals map[string]float64

	for label := range s.labels {
		labels = append(labels, label)
	}

	sort.Slice(labels, less(labels))
	sort.Slice(
		s.merchants,
		func(i int, j int) bool {
			return strless(s.merchants[i].Name, s.merchants[j].Name)
		},
	)

	// Header
	s.results = [][]string{}
	row = []string{"Merchant"}
	row = append(row, labels...)
	row = append(row, "Total")
	s.results = append(s.results, row)

	// Divider
	s.results = append(s.results, make([]string, len(labels)+2))

	// Build results row by row
	for _, m := range s.merchants {
		if s.filtered(m.Name) {
			continue
		}

		row = []string{m.Name}
		total = 0
		totals = m.Totals()

		for _, label := range labels {
			if _, ok := totals[label]; !ok {
				row = append(row, "0")
			} else {
				row = append(row, floatStr(totals[label]))
				total += totals[label]
			}
		}

		row = append(row, floatStr(total))
		s.results = append(s.results, row)

		if s.Expand {
			names = []string{}
			subtotals = m.Subtotals()

			for name := range subtotals {
				names = append(names, name)
			}

			sort.Slice(names, less(names))

			for _, name := range names {
				row = []string{" \\_ " + name}
				total = 0
				totals = subtotals[name]

				for _, label := range labels {
					if _, ok := totals[label]; !ok {
						row = append(row, "0")
					} else {
						row = append(row, floatStr(totals[label]))
						total += totals[label]
					}
				}

				row = append(row, floatStr(total))
				s.results = append(s.results, row)
			}
		}
	}

	// Divider
	s.results = append(s.results, make([]string, len(labels)+2))

	// Add label totals
	row = []string{"Total"}
	total = 0
	totals = map[string]float64{}

	for _, label := range labels {
		for _, m := range s.merchants {
			if s.filtered(m.Name) {
				continue
			}

			totals[label] += m.Totals()[label]
			total += m.Totals()[label]
		}

		row = append(row, floatStr(totals[label]))
	}

	// Add total total
	row = append(row, floatStr(total))

	s.results = append(s.results, row)
	s.tallied = true
}
