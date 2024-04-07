package main

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/mjwhitta/cli"
	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/log"
	"github.com/mjwhitta/spender"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if flags.verbose {
				panic(r.(error).Error())
			}
			log.ErrX(Exception, r.(error).Error())
		}
	}()

	var b []byte
	var cost float64
	var e error
	var f *os.File
	var lines [][]string
	var r *csv.Reader
	var s *spender.Spender = spender.New()

	validate()

	if flags.groups != "" {
		if b, e = os.ReadFile(flags.groups); e != nil {
			panic(e)
		}

		if e = s.CreateGroups(string(b)); e != nil {
			panic(e)
		}
	}

	for _, file := range cli.Args() {
		// Open file for read
		if f, e = os.Open(file); e != nil {
			panic(e)
		}
		defer f.Close()

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
			panic(hl.Errorf("%s: %w", file, e))
		}

		for _, ln := range lines {
			if len(ln) != 2 {
				log.ErrX(Exception, "invalid file")
			}

			// Normalize data
			ln[0] = strings.TrimSpace(ln[0])
			ln[1] = strings.ReplaceAll(ln[1], "$", "")
			ln[1] = strings.ReplaceAll(ln[1], ",", "")
			ln[1] = strings.TrimSpace(ln[1])

			if cost, e = strconv.ParseFloat(ln[1], 64); e != nil {
				log.ErrfX(Exception, "invalid file: %s", e)
			}

			s.Purchase(file, ln[0], cost)
		}
	}

	s.Exclude = flags.exclude
	s.Expand = flags.expand
	s.Include = flags.include

	hl.Printf("\n%s\n\n", s)
}
