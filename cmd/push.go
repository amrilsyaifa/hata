package cmd

import (
	"context"
	"fmt"

	"github.com/amrilsyaifa/hata/internal/auth"
	"github.com/amrilsyaifa/hata/internal/i18n"
	"github.com/amrilsyaifa/hata/internal/sheet"
	"github.com/spf13/cobra"
)

// baseColIndex is the 1-based column index of the "base" column in the sheet.
// Sheet layout: A=key, B=base, C=lang1, D=lang2, ...
const baseColIndex = 2

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Sync translation keys from base.json to Google Sheet",
	Long: `Reads base.json and syncs it to the sheet.

The sheet has a dedicated "base" column (column B) that is always populated
and updated from base.json by push.

Language columns (en-US, id-ID, ...) are NEVER modified by push — translators
fill those directly in the sheet. Use 'pull' to download them to local files.`,
	RunE: runPush,
}

func runPush(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	ctx := context.Background()
	opt, err := auth.ClientOption(ctx, cfg)
	if err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	client, err := sheet.New(ctx, cfg, opt)
	if err != nil {
		return err
	}

	base, err := i18n.ReadBase(cfg.Paths.Base)
	if err != nil {
		return fmt.Errorf("failed to read base file %s: %w", cfg.Paths.Base, err)
	}

	if err := client.EnsureHeaders(cfg.Languages); err != nil {
		return fmt.Errorf("failed to ensure sheet headers: %w", err)
	}

	_, rows, err := client.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read sheet: %w", err)
	}

	// Index existing sheet rows by key.
	existingRows := make(map[string]sheet.Row, len(rows))
	for _, row := range rows {
		existingRows[row.Key] = row
	}

	var newRows [][]interface{}
	var cellUpdates []sheet.CellUpdate

	for _, key := range i18n.SortedKeys(base) {
		baseVal := base[key]
		if existing, ok := existingRows[key]; ok {
			// Key already exists — update the "base" column only if value changed.
			// Language columns (en-US, id-ID, ...) are NEVER touched by push.
			if existing.Translations["base"] != baseVal {
				cellUpdates = append(cellUpdates, sheet.CellUpdate{
					RowIndex: existing.RowIndex,
					ColIndex: baseColIndex, // column B = "base"
					Value:    baseVal,
				})
			}
		} else {
			// New key — append a row with: key | baseValue | (empty lang cols)
			row := make([]interface{}, len(cfg.Languages)+2)
			row[0] = key
			row[1] = baseVal
			// columns 2..N are left empty for translators
			newRows = append(newRows, row)
		}
	}

	if len(newRows) == 0 && len(cellUpdates) == 0 {
		fmt.Println("Sheet is already up to date. Nothing to push.")
		return nil
	}

	if len(cellUpdates) > 0 {
		if err := client.BatchUpdateCells(cellUpdates); err != nil {
			return fmt.Errorf("failed to update base column: %w", err)
		}
		fmt.Printf("Updated base value for %d existing key(s).\n", len(cellUpdates))
	}

	if len(newRows) > 0 {
		if err := client.AppendRows(newRows); err != nil {
			return fmt.Errorf("failed to append rows: %w", err)
		}
		fmt.Printf("Pushed %d new key(s) to the sheet.\n", len(newRows))
	}

	return nil
}
