package cmd

import (
	"context"
	"fmt"
	"sort"

	"github.com/AlecAivazis/survey/v2"
	"github.com/amrilsyaifa/hata/internal/auth"
	"github.com/amrilsyaifa/hata/internal/i18n"
	"github.com/amrilsyaifa/hata/internal/sheet"
	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Interactively remove stale keys from Google Sheet",
	Long: `Finds keys that exist in the sheet but are no longer in base.json,
shows them in an interactive selector, and deletes the confirmed ones.`,
	RunE: runClean,
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

func runClean(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	base, err := i18n.ReadBase(cfg.Paths.Base)
	if err != nil {
		return fmt.Errorf("failed to read base file %s: %w", cfg.Paths.Base, err)
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

	_, rows, err := client.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read sheet: %w", err)
	}

	// Find keys in sheet that are NOT in base.json.
	var stale []sheet.Row
	for _, row := range rows {
		if _, ok := base[row.Key]; !ok {
			stale = append(stale, row)
		}
	}

	if len(stale) == 0 {
		fmt.Println("Sheet is clean. No stale keys found.")
		return nil
	}

	sort.Slice(stale, func(i, j int) bool {
		return stale[i].Key < stale[j].Key
	})

	options := make([]string, len(stale))
	for i, row := range stale {
		options[i] = row.Key
	}

	fmt.Printf("Found %d stale key(s) in the sheet (not in base.json):\n\n", len(stale))

	var selected []string
	selectPrompt := &survey.MultiSelect{
		Message:  "Select keys to DELETE from the sheet (Space to toggle, Enter to confirm):",
		Options:  options,
		PageSize: 20,
	}
	if err := survey.AskOne(selectPrompt, &selected); err != nil {
		return fmt.Errorf("selection cancelled: %w", err)
	}

	if len(selected) == 0 {
		fmt.Println("Nothing selected. No changes made.")
		return nil
	}

	var confirm bool
	confirmPrompt := &survey.Confirm{
		Message: fmt.Sprintf("Permanently delete %d key(s) from the sheet?", len(selected)),
		Default: false,
	}
	if err := survey.AskOne(confirmPrompt, &confirm); err != nil || !confirm {
		fmt.Println("Aborted. No changes made.")
		return nil
	}

	// Map selected keys back to row indices.
	selectedSet := make(map[string]bool, len(selected))
	for _, k := range selected {
		selectedSet[k] = true
	}
	var rowIndices []int
	for _, row := range stale {
		if selectedSet[row.Key] {
			rowIndices = append(rowIndices, row.RowIndex)
		}
	}

	if err := client.DeleteRows(rowIndices); err != nil {
		return fmt.Errorf("failed to delete rows: %w", err)
	}

	fmt.Printf("\nDeleted %d key(s) from the sheet.\n", len(rowIndices))
	return nil
}
