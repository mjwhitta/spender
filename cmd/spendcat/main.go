package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mjwhitta/cli"
	"github.com/mjwhitta/log"
	"github.com/mjwhitta/pathname"
	"github.com/mjwhitta/spender"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r)
			}

			switch r := r.(type) {
			case error:
				log.ErrX(Exception, r.Error())
			case string:
				log.ErrX(Exception, r)
			}
		}
	}()

	var cost float64
	var e error
	var f *os.File
	var lines [][]string
	var r *csv.Reader
	var s *spender.Spender = spender.New()

	validate()

	if e = readGroups(s); e != nil {
		panic(e)
	}

	for _, file := range cli.Args() {
		// Open file for read
		if f, e = os.Open(pathname.ExpandPath(file)); e != nil {
			panic(e)
		}
		defer func() {
			if e := f.Close(); e != nil {
				panic(e)
			}
		}()

		// Create CSV reader
		r = csv.NewReader(f)

		// Determine if CSV or TSV
		switch filepath.Ext(file) {
		case ".csv":
			r.Comma = ','
		case ".tsv":
			r.Comma = '\t'
		}

		// Get just filename for use as label
		file = filepath.Base(file)
		file = strings.TrimSuffix(file, filepath.Ext(file))

		// Read all lines
		if lines, e = r.ReadAll(); e != nil {
			panic(fmt.Errorf("%s: %w", file, e))
		}

		for _, ln := range lines {
			//nolint:mnd // expecting 2 fields
			if len(ln) != 2 {
				log.ErrX(Exception, "invalid file")
			}

			// Normalize data
			ln[0] = strings.TrimSpace(ln[0])
			ln[1] = strings.ReplaceAll(ln[1], "$", "")
			ln[1] = strings.ReplaceAll(ln[1], ",", "")
			ln[1] = strings.TrimSpace(ln[1])

			if cost, e = strconv.ParseFloat(ln[1], 64); e != nil {
				log.ErrXf(Exception, "invalid file: %s", e)
			}

			s.Purchase(file, ln[0], cost)
		}
	}

	s.Exclude = flags.exclude
	s.Expand = flags.expand
	s.Include = flags.include

	fmt.Printf("\n%s\n\n", s)
}

func readGroups(s *spender.Spender) error {
	var b []byte
	var e error

	if flags.groups != "" {
		if b, e = os.ReadFile(flags.groups); e != nil {
			e = fmt.Errorf("failed to read %s: %w", flags.groups, e)
			return e
		}

		if e = s.CreateGroups(b); e != nil {
			//nolint:wrapcheck // I'm not wrapping my own functions
			return e
		}
	}

	return nil
}
