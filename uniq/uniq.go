package uniq

import (
	"fmt"
	"strings"
)

type Options struct {
	Count      bool // -c
	OnlyDup    bool // -d
	OnlyUnique bool // -u
	IgnoreCase bool // -i
	SkipFields int  // -f
	SkipChars  int  // -s
}

func emit(orig string, count int, opts Options, result []string) []string {
	if orig == "" {
		return result
	}
	if opts.Count {
		result = append(result, fmt.Sprintf("%d %s", count, orig))
	} else if opts.OnlyDup {
		if count > 1 {
			result = append(result, orig)
		}
	} else if opts.OnlyUnique {
		if count == 1 {
			result = append(result, orig)
		}
	} else {
		result = append(result, orig)
	}
	return result
}

func UniqueLines(lines []string, opts Options) []string {
	var result []string
	if len(lines) == 0 {
		return result
	}

	prevKey := ""
	prevOrig := ""
	count := 0

	for _, line := range lines {
		key := transformKey(line, opts)
		if key == prevKey {
			count++
		} else {
			emit(prevOrig, count, opts, result)
			prevKey = key
			prevOrig = line
			count = 1
		}
	}
	emit(prevOrig, count, opts, result)
	return result
}

func transformKey(s string, opts Options) string {
	// игнорируем поля
	fields := strings.Fields(s)
	if opts.SkipFields > 0 && len(fields) > opts.SkipFields {
		s = strings.Join(fields[opts.SkipFields:], " ")
	}
	// игнорируем символы
	if opts.SkipChars > 0 && len(s) > opts.SkipChars {
		s = s[opts.SkipChars:]
	}
	// игнорируем регистр
	if opts.IgnoreCase {
		s = strings.ToLower(s)
	}
	return s
}
