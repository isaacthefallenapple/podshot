package rss

import (
	"bufio"
	"io"
	"regexp"
)

type itemScanner struct {
	*bufio.Scanner
}

func newItemScanner(reader io.Reader) *itemScanner {
	scanner := itemScanner{bufio.NewScanner(reader)}
	scanner.Split(splitByItem)
	scanner.Buffer(make([]byte, 4096*4), 4096*8)
	return &scanner
}

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
