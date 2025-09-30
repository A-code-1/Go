package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"myuniq/uniq"
)

func main() {
	count := flag.Bool("c", false, "подсчитать количество повторов")
	dup := flag.Bool("d", false, "только повторяющиеся строки")
	uniqFlag := flag.Bool("u", false, "только уникальные строки")
	ignoreCase := flag.Bool("i", false, "игнорировать регистр")
	skipFields := flag.Int("f", 0, "пропустить _ полей")
	skipChars := flag.Int("s", 0, "пропустить _ символов")
	flag.Parse()

	// проверка на конфликт флагов
	modeFlags := 0
	for _, f := range []bool{*count, *dup, *uniqFlag} {
		if f {
			modeFlags++
		}
	}
	if modeFlags > 1 {
		fmt.Fprintln(os.Stderr, "флаги -c, -d и -u взаимоисключающие")
		os.Exit(1)
	}

	opts := uniq.Options{
		Count:      *count,
		OnlyDup:    *dup,
		OnlyUnique: *uniqFlag,
		IgnoreCase: *ignoreCase,
		SkipFields: *skipFields,
		SkipChars:  *skipChars,
	}

	args := flag.Args()
	var reader io.Reader = os.Stdin
	if len(args) > 0 {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "uniq:", err)
			os.Exit(1)
		}
		defer file.Close()
		reader = file
	}

	var writer io.Writer = os.Stdout
	if len(args) > 1 {
		file, err := os.Create(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "uniq:", err)
			os.Exit(1)
		}
		defer file.Close()
		writer = file
	}

	data, err := io.ReadAll(reader)
	if err != nil {
		fmt.Fprintln(os.Stderr, "uniq:", err)
		os.Exit(1)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	result := uniq.UniqueLines(lines, opts)
	_, _ = fmt.Fprintln(writer, strings.Join(result, "\n"))
}
