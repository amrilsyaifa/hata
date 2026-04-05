package cmd

import (
	"context"
	"fmt"

	"github.com/amrilsyaifa/hata/internal/auth"
	"github.com/amrilsyaifa/hata/internal/i18n"
	"github.com/amrilsyaifa/hata/internal/sheet"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Generate per-language JSON files from Google Sheet",
	Long:  "Reads all translations from the sheet and generates one JSON file per language under the configured output directory.",
	RunE:  runPull,
}

func runPull(_ *cobra.Command, _ []string) error {
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

	_, rows, err := client.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read sheet: %w", err)
	}

	if len(rows) == 0 {
		fmt.Println("Sheet is empty. Nothing to pull.")
		return nil
	}

	fileKey := func(lang string) string {
		if alias, ok := cfg.Aliases[lang]; ok && alias != "" {
			return alias
		}
		return lang
	}

	langData := make(map[string]map[string]string, len(cfg.Languages))
	for _, lang := range cfg.Languages {
		langData[fileKey(lang)] = make(map[string]string)
	}

	warnings := 0
	for _, row := range rows {
		for _, lang := range cfg.Languages {
			val := row.Translations[lang]
			if val == "" {
				fmt.Printf("  Warning: missing translation for key %q in language %q\n", row.Key, lang)
				warnings++
				continue
			}
			langData[fileKey(lang)][row.Key] = val
		}
	}

	if err := i18n.GenerateLocaleFiles(cfg.Paths.Output, langData, cfg.Options.SortKeys, cfg.Options.NestedJSON); err != nil {
		return err
	}

	fmt.Printf("\nPull complete. %d warning(s).\n", warnings)
	return nil
}
