package sheet

import (
	"context"
	"fmt"

	"github.com/amrilsyaifa/hata/internal/config"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Client struct {
	service   *sheets.Service
	sheetID   string
	sheetName string
}

type Row struct {
	Key          string
	Translations map[string]string
	RowIndex     int // 1-based sheet row index (row 1 = header)
}

// CellUpdate describes a single cell to overwrite.
type CellUpdate struct {
	RowIndex int // 1-based
	ColIndex int // 1-based (A=1, B=2, ...)
	Value    string
}

// columnLetter converts a 1-based column index to a spreadsheet letter (1→A, 2→B, …).
func columnLetter(col int) string {
	result := ""
	for col > 0 {
		col--
		result = string(rune('A'+col%26)) + result
		col /= 26
	}
	return result
}

func New(ctx context.Context, cfg *config.Config, opt option.ClientOption) (*Client, error) {
	svc, err := sheets.NewService(ctx, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create Sheets service: %w", err)
	}
	return &Client{
		service:   svc,
		sheetID:   cfg.Sheet.ID,
		sheetName: cfg.Sheet.Name,
	}, nil
}

func (c *Client) rangeStr(r string) string {
	return fmt.Sprintf("'%s'!%s", c.sheetName, r)
}

func (c *Client) ReadAll() ([]string, []Row, error) {
	resp, err := c.service.Spreadsheets.Values.Get(c.sheetID, c.sheetName).Do()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read sheet: %w", err)
	}

	if len(resp.Values) == 0 {
		return nil, nil, nil
	}

	headerRow := resp.Values[0]
	langs := make([]string, 0, len(headerRow)-1)
	for _, h := range headerRow[1:] {
		langs = append(langs, fmt.Sprintf("%v", h))
	}

	var rows []Row
	seen := make(map[string]bool)
	for rowIdx, rawRow := range resp.Values[1:] {
		if len(rawRow) == 0 {
			continue
		}
		key := fmt.Sprintf("%v", rawRow[0])
		if key == "" {
			continue
		}
		if seen[key] {
			return nil, nil, fmt.Errorf("duplicate key found in sheet: %q", key)
		}
		seen[key] = true

		row := Row{
			Key:          key,
			Translations: make(map[string]string, len(langs)),
			RowIndex:     rowIdx + 2, // +1 for header, +1 for 1-based indexing
		}
		for i, lang := range langs {
			if i+1 < len(rawRow) {
				row.Translations[lang] = fmt.Sprintf("%v", rawRow[i+1])
			}
		}
		rows = append(rows, row)
	}

	return langs, rows, nil
}

func (c *Client) EnsureHeaders(languages []string) error {
	resp, err := c.service.Spreadsheets.Values.Get(c.sheetID, c.rangeStr("A1:Z1")).Do()
	if err != nil {
		return fmt.Errorf("failed to read sheet headers: %w", err)
	}

	if len(resp.Values) > 0 && len(resp.Values[0]) > 0 {
		return nil
	}

	// Sheet columns: key | base | lang1 | lang2 | ...
	// "base" holds the source-language value from base.json (managed by push).
	// Language columns are filled by translators directly in the sheet.
	headers := make([]interface{}, len(languages)+2)
	headers[0] = "key"
	headers[1] = "base"
	for i, lang := range languages {
		headers[i+2] = lang
	}

	_, err = c.service.Spreadsheets.Values.Update(
		c.sheetID,
		c.rangeStr("A1"),
		&sheets.ValueRange{Values: [][]interface{}{headers}},
	).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("failed to write headers: %w", err)
	}
	return nil
}

func (c *Client) AppendRows(rows [][]interface{}) error {
	if len(rows) == 0 {
		return nil
	}
	_, err := c.service.Spreadsheets.Values.Append(
		c.sheetID,
		c.sheetName,
		&sheets.ValueRange{Values: rows},
	).ValueInputOption("RAW").Do()
	if err != nil {
		return fmt.Errorf("failed to append rows: %w", err)
	}
	return nil
}

// BatchUpdateCells overwrites specific cells without touching any other cells.
func (c *Client) BatchUpdateCells(updates []CellUpdate) error {
	if len(updates) == 0 {
		return nil
	}
	var valueRanges []*sheets.ValueRange
	for _, u := range updates {
		rangeStr := fmt.Sprintf("'%s'!%s%d", c.sheetName, columnLetter(u.ColIndex), u.RowIndex)
		valueRanges = append(valueRanges, &sheets.ValueRange{
			Range:  rangeStr,
			Values: [][]interface{}{{u.Value}},
		})
	}
	_, err := c.service.Spreadsheets.Values.BatchUpdate(c.sheetID, &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "RAW",
		Data:             valueRanges,
	}).Do()
	if err != nil {
		return fmt.Errorf("failed to batch update cells: %w", err)
	}
	return nil
}
