package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/amrilsyaifa/hata/internal/auth"
	"github.com/amrilsyaifa/hata/internal/sheet"
	"github.com/spf13/cobra"
)

var (
	importFile string
	importLang string
)

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import an existing locale JSON file into Google Sheet",
	Long: `Reads a locale JSON file (flat or nested) and writes its values into
the matching language column in the sheet. Only rows whose keys already
exist in the sheet are updated — new keys are never created by this command
(use 'push' first to create them).`,
	RunE: runImport,
}

func init() {
	importCmd.Flags().StringVarP(&importFile, "file", "f", "", "path to the locale JSON file (required)")
	importCmd.Flags().StringVarP(&importLang, "lang", "l", "", "locale code as it appears in the sheet header, e.g. id-ID (required)")
	_ = importCmd.MarkFlagRequired("file")
	_ = importCmd.MarkFlagRequired("lang")
	rootCmd.AddCommand(importCmd)
}

func runImport(_ *cobra.Command, _ []string) error {
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Read and flatten the locale JSON file.
	data, err := os.ReadFile(importFile)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", importFile, err)
	}

	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("invalid JSON in %s: %w", importFile, err)
	}

	flat := make(map[string]string)
	flattenAny("", raw, flat)

	if len(flat) == 0 {
		fmt.Println("File is empty. Nothing to import.")
		return nil
	}

	fmt.Printf("Loaded %d key(s) from %s\n", len(flat), importFile)

	ctx := context.Background()
	opt, err := auth.ClientOption(ctx, cfg)
	if err != nil {
		return fmt.Errorf("auth failed: %w", err)
	}

	client, err := sheet.New(ctx, cfg, opt)
	if err != nil {
		return err
	}

	written, missing, err := client.UpdateTranslations(importLang, flat)
	if err != nil {
		return fmt.Errorf("failed to update sheet: %w", err)
	}

	fmt.Printf("Updated %d cell(s) in column %q.\n", written, importLang)

	if len(missing) > 0 {
		sort.Strings(missing)
		fmt.Printf("\n%d key(s) not found in sheet (run 'hata push' first to add them):\n", len(missing))
		for _, k := range missing {
			fmt.Printf("  - %s\n", k)
		}
	}

	return nil
}

// flattenAny recursively flattens a JSON value into dot-separated keys.
func flattenAny(prefix string, v interface{}, out map[string]string) {
	switch val := v.(type) {
	case map[string]interface{}:
		for k, child := range val {
			key := k
			if prefix != "" {
				key = prefix + "." + k
			}
			flattenAny(key, child, out)
		}
	case string:
		out[prefix] = val
	default:
		if val != nil {
			out[prefix] = strings.TrimSpace(fmt.Sprintf("%v", val))
		}
	}
}
