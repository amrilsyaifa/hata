package cmd

import (
	"context"
	"fmt"

	idiff "github.com/amrilsyaifa/hata/internal/diff"
	"github.com/amrilsyaifa/hata/internal/auth"
	"github.com/amrilsyaifa/hata/internal/i18n"
	"github.com/amrilsyaifa/hata/internal/sheet"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Show differences between base.json and Google Sheet",
	Long: `Compares keys in base.json against the sheet and reports:
  - Keys missing in the sheet (need to be pushed)
  - Keys in the sheet not in base.json (potentially stale)`,
	RunE: runDiff,
}

func runDiff(_ *cobra.Command, _ []string) error {
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

	result := idiff.Compare(base, rows)
	idiff.Print(result)
	return nil
}
