package main

import (
	"errors"
	"strings"
)

func tokenize(line string) ([]string, error) {
	var tokens []string
	var cur strings.Builder
	hasToken := false

	const (
		normal = iota
		inSingle
		inDouble
	)
	state := normal

	flush := func() {
		if hasToken {
			tokens = append(tokens, cur.String())
			cur.Reset()
			hasToken = false
		}
	}

	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		switch state {
		case normal:
			switch {
			case r == '\'':
				state = inSingle
				hasToken = true
			case r == '"':
				state = inDouble
				hasToken = true
			case r == ' ' || r == '\t':
				flush()
			default:
				cur.WriteRune(r)
				hasToken = true
			}
		case inSingle:
			if r == '\'' {
				state = normal
			} else {
				cur.WriteRune(r)
			}
		case inDouble:
			if r == '"' {
				state = normal
			} else if r == '\\' && i+1 < len(runes) && (runes[i+1] == '"' || runes[i+1] == '\\') {
				i++
				cur.WriteRune(runes[i])
			} else {
				cur.WriteRune(r)
			}
		}
	}

	if state != normal {
		return nil, errors.New("unterminated quote")
	}
	flush()

	return tokens, nil
}
