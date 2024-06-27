# Spender

[![Yum](https://img.shields.io/badge/-Buy%20me%20a%20cookie-blue?labelColor=grey&logo=cookiecutter&style=for-the-badge)](https://www.buymeacoffee.com/mjwhitta)

[![Go Report Card](https://goreportcard.com/badge/github.com/mjwhitta/spender?style=for-the-badge)](https://goreportcard.com/report/github.com/mjwhitta/spender)
![License](https://img.shields.io/github/license/mjwhitta/spender?style=for-the-badge)

## Description

Spender can help track your expenses. The included `spendcat` tool
expects a list of CSV or TSV files for input. Costs are grouped by
merchant. You can optionally provide a JSON file to group multiple
merchants into categories. The CSV/TSV files should have only two
columns: merchant and cost.

The groups file file should have the form:

{
  "group1": ["merchant1", "merchant2"],
  "group2": ["merchant3", "merchant4"]
}

The merchants in the groups file can be strings or regular
expressions. Regular expressions begin and end with `/` and have an
optional trailing `i` for case-insensitive matching (e.g. `/regex/` or
`/regex/i`). It's important to note that this is a JSON file, so
backslashes should be escaped inside regular expressions. Group names
are sorted case-insensitively and the group name `ignore` is special
and will be ignored for all calculations.

## How to install

Open a terminal and run the following:

```
$ go install github.com/mjwhitta/spender/cmd/spendcat@latest
```

## Usage

Run `spendcat -h` to see the full usage.
