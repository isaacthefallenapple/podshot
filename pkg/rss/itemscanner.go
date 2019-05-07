// Package rss provides a Scanner-type to parse a podcast RSS feed item by item.
package rss

import (
	"bufio"
	"io"
	"regexp"
)

// itemScanner is a generic wrapper for a *bufio.Scanner
type itemScanner struct { //TODO get rid of this, the constructor will do.
	*bufio.Scanner
}

// newItemScanner returns a new itemScanner
func newItemScanner(reader io.Reader) *itemScanner {
	scanner := itemScanner{bufio.NewScanner(reader)}
	scanner.Split(splitByItem)
	scanner.Buffer(make([]byte, 4096*4), 4096*8)
	return &scanner
}

// splitByItem is the SplitFunc for the itemScanner type. It parses the RSS feed from <item> to </item>.
func splitByItem(data []byte, atEOF bool) (advance int, token []byte, err error) {

	if atEOF && len(data) == 0 {
		return
	}

	var startIndex, endIndex int

	startRegex := regexp.MustCompile(`(?m)<item>`)
	endRegex := regexp.MustCompile(`(?m)</item>`)

	if startMatch := startRegex.FindIndex(data); len(startMatch) != 0 {
		startIndex = startMatch[0]
	} else {
		return
	}

	if endMatch := endRegex.FindIndex(data); len(endMatch) != 0 {
		endIndex = endMatch[1]
	} else {
		return
	}

	advance, token = endIndex, data[startIndex:endIndex]

	return
}
