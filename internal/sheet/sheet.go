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

// UpdateTranslations writes flat translation values for the given language column.
// It returns the number of cells written and a list of keys that were not found in the sheet.
func (c *Client) UpdateTranslations(lang string, translations map[string]string) (int, []string, error) {
	langs, rows, err := c.ReadAll()
	if err != nil {
		return 0, nil, err
	}

	// Locate the column index for the requested language (1-based for Sheets A=1).
	// Header layout: col1=key, col2=base, then each lang in order.
	langColIndex := -1
	for i, h := range langs {
		if h == lang {
			langColIndex = i + 2 // +1 for key col, +1 for 1-based
			break
		}
	}
	if langColIndex == -1 {
		return 0, nil, fmt.Errorf("language column %q not found in sheet (run 'hata push' first)", lang)
	}

	// Build a key→row lookup.
	rowByKey := make(map[string]Row, len(rows))
	for _, row := range rows {
		rowByKey[row.Key] = row
	}

	var updates []CellUpdate
	var missing []string

	for key, value := range translations {
		r, ok := rowByKey[key]
		if !ok {
			missing = append(missing, key)
			continue
		}
		updates = append(updates, CellUpdate{
			RowIndex: r.RowIndex,
			ColIndex: langColIndex,
			Value:    value,
		})
	}

	if err := c.BatchUpdateCells(updates); err != nil {
		return 0, missing, err
	}
	return len(updates), missing, nil
}

// DeleteRows removes the given 1-based row indices from the sheet.
// Rows are deleted bottom-up to avoid index shifting.
func (c *Client) DeleteRows(rowIndices []int) error {
	if len(rowIndices) == 0 {
		return nil
	}

	// Resolve the numeric sheetId for the named tab.
	spreadsheet, err := c.service.Spreadsheets.Get(c.sheetID).Do()
	if err != nil {
		return fmt.Errorf("failed to get spreadsheet metadata: %w", err)
	}
	var tabID int64
	found := false
	for _, s := range spreadsheet.Sheets {
		if s.Properties.Title == c.sheetName {
			tabID = s.Properties.SheetId
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("sheet tab %q not found", c.sheetName)
	}

	// Sort descending so deleting one row does not shift the next.
	sorted := make([]int, len(rowIndices))
	copy(sorted, rowIndices)
	for i := 0; i < len(sorted)-1; i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j] > sorted[i] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	// Build one DeleteDimension request per row.
	requests := make([]*sheets.Request, 0, len(sorted))
	for _, rowIdx := range sorted {
		zero := int64(rowIdx - 1)
		requests = append(requests, &sheets.Request{
			DeleteDimension: &sheets.DeleteDimensionRequest{
				Range: &sheets.DimensionRange{
					SheetId:    tabID,
					Dimension:  "ROWS",
					StartIndex: zero,
					EndIndex:   zero + 1,
				},
			},
		})
	}

	_, err = c.service.Spreadsheets.BatchUpdate(c.sheetID, &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}).Do()
	if err != nil {
		return fmt.Errorf("failed to delete rows: %w", err)
	}
	return nil
}
