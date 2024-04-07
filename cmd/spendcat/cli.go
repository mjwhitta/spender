package main

import (
	"os"

	"github.com/mjwhitta/cli"
	"github.com/mjwhitta/errors"
	hl "github.com/mjwhitta/hilighter"
	"github.com/mjwhitta/pathname"
	"github.com/mjwhitta/spender"
)

// Exit status
const (
	Good = iota
	InvalidOption
	MissingOption
	InvalidArgument
	MissingArgument
	ExtraArgument
	Exception
)

// Flags
var flags struct {
	exclude cli.StringList
	expand  bool
	groups  string
	include cli.StringList
	nocolor bool
	verbose bool
	version bool
}

func init() {
	// Configure cli package
	cli.Align = true // Defaults to false
	cli.Authors = []string{"Miles Whittaker <mj@whitta.dev>"}
	cli.Banner = hl.Sprintf(
		"%s [OPTIONS] <file1>... [fileN]",
		os.Args[0],
	)
	cli.BugEmail = "spendcat.bugs@whitta.dev"
	cli.ExitStatus(
		"Normally the exit status is 0. In the event of an error the",
		"exit status will be one of the below:\n\n",
		hl.Sprintf("  %d: Invalid option\n", InvalidOption),
		hl.Sprintf("  %d: Missing option\n", MissingOption),
		hl.Sprintf("  %d: Invalid argument\n", InvalidArgument),
		hl.Sprintf("  %d: Missing argument\n", MissingArgument),
		hl.Sprintf("  %d: Extra argument\n", ExtraArgument),
		hl.Sprintf("  %d: Exception", Exception),
	)
	cli.Info(
		"SpendCat expects a list of CSV or TSV files for input.",
		"Costs are grouped by merchant. You can optionally provide a",
		"JSON file to group multiple merchants into categories.\n\n",
		"The CSV/TSV files should have only two columns: merchant",
		"and cost.\n\n",
		"The groups file should have the form:\n\n",
		"{\n",
		"\"group1\": [\"merchant1\", \"merchant2\"],\n",
		"\"group2\": [\"merchant3\", \"merchant4\"]\n",
		"}\n\n",
		"The merchants in the groups file can be strings or regular",
		"expressions. Regular expressions begin and end with / and",
		"have an optional trailing \"i\" for case-insensitive",
		"matching (e.g. /.*regex.*/ or /.*regex.*/i). It's important",
		"to note that this is a JSON file, so backslashes should be",
		"escaped in regular expressions. All group names are",
		"case-insensitive and the group name \"ignore\" is special",
		"and will be ignored for calculations and output.",
	)
	cli.Title = "SpendCat"

	// Parse cli flags
	cli.Flag(
		&flags.exclude,
		"e",
		"exclude",
		"Filter out provided group or merchant (can be used more",
		"than once).",
	)
	cli.Flag(
		&flags.expand,
		"x",
		"expand",
		false,
		"Expand to show subtotals for grouped merchants.",
	)
	cli.Flag(
		&flags.groups,
		"g",
		"groups",
		"",
		"Group merchants using patterns in provided JSON file.",
	)
	cli.Flag(
		&flags.include,
		"i",
		"include",
		"Filter results to show only provided group or merchant (can",
		"be used more than once).",
	)
	cli.Flag(
		&flags.nocolor,
		"no-color",
		false,
		"Disable colorized output.",
	)
	cli.Flag(
		&flags.verbose,
		"v",
		"verbose",
		false,
		"Show stacktrace, if error.",
	)
	cli.Flag(&flags.version, "V", "version", false, "Show version.")
	cli.Parse()
}

// Process cli flags and ensure no issues
func validate() {
	hl.Disable(flags.nocolor)

	// Short circuit, if version was requested
	if flags.version {
		hl.Printf("spendcat version %s\n", spender.Version)
		os.Exit(Good)
	}

	// Validate cli flags
	if cli.NArg() == 0 {
		cli.Usage(MissingArgument)
	}

	// Ensure group file exists, if provided
	if flags.groups != "" {
		if ok, e := pathname.DoesExist(flags.groups); e != nil {
			panic(e)
		} else if !ok {
			e = errors.Newf("group file %s not found", flags.groups)
			panic(e)
		}
	}

	// Ensure each CSV file exists
	for _, csv := range cli.Args() {
		if ok, e := pathname.DoesExist(csv); e != nil {
			panic(e)
		} else if !ok {
			panic(errors.Newf("csv file %s not found", csv))
		}
	}
}
