package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/amrilsyaifa/hata/internal/config"
	"github.com/amrilsyaifa/hata/internal/locale"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize hata configuration interactively",
	RunE:  runInit,
}

func runInit(_ *cobra.Command, _ []string) error {
	reader := bufio.NewReader(os.Stdin)

	printBanner()
	fmt.Println("Let's set up your i18n configuration.")
	fmt.Println()

	projectID := prompt(reader, "Project ID", "my-i18n")
	sheetID := prompt(reader, "Google Sheet ID", "")
	sheetName := prompt(reader, "Sheet tab name", "Translations")

	fmt.Println("\nAuthentication method:")
	fmt.Println("  1. Service Account (recommended for CI/CD)")
	fmt.Println("  2. OAuth (for interactive / multi-user environments)")
	authChoice := prompt(reader, "Choice [1/2]", "1")

	var authType string
	switch authChoice {
	case "2":
		authType = "oauth"
	default:
		authType = "service_account"
	}

	credPath := prompt(reader, "Credentials file path", ".i18n/credentials.json")
	tokenPath := ".i18n/token.json"
	if authType == "oauth" {
		tokenPath = prompt(reader, "Token cache path", ".i18n/token.json")
	}

	// Interactive locale multi-select
	fmt.Println()
	langs, err := locale.Select()
	if err != nil {
		return err
	}
	if len(langs) == 0 {
		return fmt.Errorf("no languages selected — at least one is required")
	}

	basePath := prompt(reader, "\nBase file path", "./base.json")
	outputPath := prompt(reader, "Output directory", "./locales")

	newCfg := &config.Config{
		ProjectID: projectID,
		Sheet: config.SheetConf{
			ID:   sheetID,
			Name: sheetName,
		},
		Auth: config.AuthConf{
			Type:            authType,
			CredentialsPath: credPath,
			TokenPath:       tokenPath,
		},
		Languages: langs,
		Paths: config.PathConf{
			Base:   basePath,
			Output: outputPath,
		},
		Options: config.Options{
			NestedJSON: true,
			SortKeys:   true,
			KeepUnused: true,
		},
	}

	if err := config.Save(cfgFile, newCfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("\nConfig saved to %s\n", cfgFile)
	fmt.Printf("Selected languages: %s\n", strings.Join(langs, ", "))
	return nil
}

func prompt(reader *bufio.Reader, label, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", label, defaultVal)
	} else {
		fmt.Printf("%s: ", label)
	}
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultVal
	}
	return line
}
