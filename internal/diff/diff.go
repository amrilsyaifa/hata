package diff

import (
	"fmt"
	"sort"

	"github.com/amrilsyaifa/hata/internal/sheet"
)

type Result struct {
	MissingInSheet []string
	UnusedInBase   []string
}

func Compare(base map[string]string, rows []sheet.Row) Result {
	sheetKeys := make(map[string]bool, len(rows))
	for _, row := range rows {
		sheetKeys[row.Key] = true
	}

	var result Result
	for k := range base {
		if !sheetKeys[k] {
			result.MissingInSheet = append(result.MissingInSheet, k)
		}
	}
	for k := range sheetKeys {
		if _, ok := base[k]; !ok {
			result.UnusedInBase = append(result.UnusedInBase, k)
		}
	}

	sort.Strings(result.MissingInSheet)
	sort.Strings(result.UnusedInBase)
	return result
}

func Print(r Result) {
	if len(r.MissingInSheet) == 0 && len(r.UnusedInBase) == 0 {
		fmt.Println("Everything is in sync!")
		return
	}

	if len(r.MissingInSheet) > 0 {
		fmt.Println("Missing in sheet:")
		for _, k := range r.MissingInSheet {
			fmt.Printf("  - %s\n", k)
		}
	}

	if len(r.UnusedInBase) > 0 {
		if len(r.MissingInSheet) > 0 {
			fmt.Println()
		}
		fmt.Println("Unused in base:")
		for _, k := range r.UnusedInBase {
			fmt.Printf("  - %s\n", k)
		}
	}
}
