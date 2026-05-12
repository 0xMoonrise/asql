package kernel

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/0xMoonrise/asql/internal/lexer"
	"github.com/0xMoonrise/asql/internal/parser"
)

type DataType string

type KernelError struct {
	Code int
	Msg  string
}

const (
	TypeNumeric DataType = "NUMERIC"
	TypeChar    DataType = "CHAR"
	TypeDate    DataType = "DATE"
)

const (
	CodeDuplicateTable      = 301
	CodeDuplicateColumn     = 302
	CodeDuplicateConstraint = 303
	CodeTableNotFound       = 304
	CodeColumnNotFound      = 305
)

type Column struct {
	Name     string   `json:"name"`
	DataType DataType `json:"data_type"`
	Nullable bool     `json:"nullable"`
}

type Constraint struct {
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	Columns    []string `json:"columns"`
	RefTable   string   `json:"ref_table,omitempty"`
	RefColumns []string `json:"ref_columns,omitempty"`
	Expression string   `json:"expression,omitempty"`
}

type Table struct {
	Name        string           `json:"name"`
	Columns     []Column         `json:"columns"`
	Constraints []Constraint     `json:"constraints"`
	Rows        []map[string]any `json:"rows"`
}

type Catalog struct {
	Tables map[string]Table `json:"tables"`
}

type Kernel struct {
	catalog  Catalog
	jsonPath string
	parser   *parser.StackParser
}

func NewKernel(jsonPath string) *Kernel {
	k := &Kernel{
		catalog:  Catalog{Tables: make(map[string]Table)},
		jsonPath: jsonPath,
	}
	if data, err := os.ReadFile(jsonPath); err == nil {
		json.Unmarshal(data, &k.catalog)
	}
	return k
}

func (k *Kernel) save() error {
	data, err := json.MarshalIndent(k.catalog, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(k.jsonPath, data, 0644)
}

func (k *Kernel) ProcessSQL(sql string) error {
	tokensRaw := lexer.Tokenize(sql)
	lines := [][]string{tokensRaw}
	tokensStream, errs := lexer.Lexer(lines)
	if len(errs) > 0 {
		return fmt.Errorf("lexer errors: %v", errs)
	}

	p := parser.NewParser(tokensStream)
	if err := p.Parse(); err != nil {
		return fmt.Errorf("parse error: %v", err)
	}

	meta := p.Metadata

	if len(meta["ct_tables"]) > 0 {
		if err := k.handleCreateTable(meta); err != nil {
			return err
		}
	}

	if len(meta["insert_tables"]) > 0 {
		if err := k.handleInsert(meta); err != nil {
			return err
		}
	}

	if len(meta["select_tables"]) > 0 {
		if err := k.handleSelect(meta); err != nil {
			return err
		}
	}

	return nil
}

func (k *Kernel) handleCreateTable(meta map[string][]string) error {

	for _, tableName := range meta["ct_tables"] {
		if _, exists := k.catalog.Tables[tableName]; exists {
			return &KernelError{CodeDuplicateTable, fmt.Sprintf("Table '%s' already exists", tableName)}
		}

		colKey := "ct_columns:" + tableName
		colNames := meta[colKey]
		conKey := "ct_constraints:" + tableName
		constraintNames := meta[conKey]
		colTypes := meta["ct_types:"+tableName]

		newTable := Table{
			Name:        tableName,
			Columns:     []Column{},
			Constraints: []Constraint{},
			Rows:        []map[string]any{},
		}

		colSet := make(map[string]bool)
		for i, col := range colNames {
			dt := DataType(colTypes[i])
			if colSet[col] {
				return &KernelError{CodeDuplicateColumn, fmt.Sprintf("Duplicate column '%s' in table '%s'", col, tableName)}
			}

			colSet[col] = true

			newTable.Columns = append(newTable.Columns, Column{
				Name:     col,
				DataType: dt,
				Nullable: true,
			})
		}

		conSet := make(map[string]bool)
		for _, con := range constraintNames {
			if conSet[con] {
				return &KernelError{CodeDuplicateConstraint, fmt.Sprintf("Duplicate constraint '%s' in table '%s'", con, tableName)}
			}
			conSet[con] = true
			newTable.Constraints = append(newTable.Constraints, Constraint{Name: con, Type: "Default"})
		}

		k.catalog.Tables[tableName] = newTable
	}
	return k.save()
}

func (k *Kernel) handleInsert(meta map[string][]string) error {
	tableName := meta["insert_tables"][0]
	values := meta["insert_values"]

	table, exists := k.catalog.Tables[tableName]
	if !exists {
		return &KernelError{CodeTableNotFound, fmt.Sprintf("Table '%s' not found", tableName)}
	}

	expected := len(table.Columns)
	if len(values) != expected {
		return errors.New("number of values does not match number of columns")
	}

	row := make(map[string]any)
	for i, col := range table.Columns {
		raw := values[i]
		isString := strings.HasPrefix(raw, "'") && strings.HasSuffix(raw, "'")

		switch col.DataType {
		case TypeNumeric:
			if isString {
				return &KernelError{
					CodeColumnNotFound,
					fmt.Sprintf("column '%s' expects NUMERIC, got string '%s'", col.Name, raw),
				}
			}
		case TypeChar, TypeDate:
			if !isString {
				return &KernelError{
					CodeColumnNotFound,
					fmt.Sprintf("column '%s' expects %s, got NUMERIC '%s'", col.Name, col.DataType, raw),
				}
			}
		}

		if strings.HasPrefix(raw, "'") && strings.HasSuffix(raw, "'") {
			raw = raw[1 : len(raw)-1]
			row[col.Name] = raw
		} else {
			var num float64
			if _, err := fmt.Sscan(raw, &num); err == nil {
				row[col.Name] = num
			} else {
				row[col.Name] = raw
			}
		}
	}

	table.Rows = append(table.Rows, row)
	k.catalog.Tables[tableName] = table

	return k.save()
}

func (k *Kernel) handleSelect(meta map[string][]string) error {
	tables := meta["select_tables"]
	columns := meta["select_columns"]

	for _, t := range tables {
		if _, exists := k.catalog.Tables[t]; !exists {
			return &KernelError{CodeTableNotFound, fmt.Sprintf("Table '%s' not found", t)}
		}
	}

	type SelectedColumn struct {
		TableName  string
		ColumnName string
		Alias      string
	}
	var selected []SelectedColumn

	if len(columns) == 1 && columns[0] == "*" {
		for _, t := range tables {
			table := k.catalog.Tables[t]
			for _, col := range table.Columns {
				selected = append(selected, SelectedColumn{TableName: t, ColumnName: col.Name})
			}
		}
	} else {

		for _, colExpr := range columns {
			parts := strings.Split(colExpr, ".")
			var tableName, colName string
			if len(parts) == 2 {
				tableName = parts[0]
				colName = parts[1]
			} else {

				if len(tables) != 1 {
					return errors.New("column '" + colExpr + "' without table prefix but multiple tables in FROM")
				}
				tableName = tables[0]
				colName = parts[0]
			}

			table, ok := k.catalog.Tables[tableName]
			if !ok {
				return &KernelError{CodeTableNotFound, fmt.Sprintf("Table '%s' not found", tableName)}
			}

			found := false
			for _, col := range table.Columns {
				if col.Name == colName {
					found = true
					break
				}
			}

			if !found {
				return &KernelError{CodeColumnNotFound, fmt.Sprintf("Column '%s' not found in table '%s'", colName, tableName)}
			}

			selected = append(selected, SelectedColumn{TableName: tableName, ColumnName: colName})
		}
	}

	type RowMap map[string]any

	var result []RowMap

	if len(tables) == 1 {
		table := k.catalog.Tables[tables[0]]
		for _, row := range table.Rows {
			newRow := make(RowMap)
			for _, sel := range selected {
				if sel.TableName != tables[0] {
					continue
				}
				if val, ok := row[sel.ColumnName]; ok {
					newRow[sel.TableName+"."+sel.ColumnName] = val
				} else {
					newRow[sel.TableName+"."+sel.ColumnName] = nil
				}
			}
			result = append(result, newRow)
		}

	} else {

		tableRows := make([][]map[string]any, len(tables))
		for i, t := range tables {
			tableRows[i] = k.catalog.Tables[t].Rows
		}

		var combine func(depth int, current RowMap)
		combine = func(depth int, current RowMap) {
			if depth == len(tables) {
				rowCopy := make(RowMap)
				maps.Copy(rowCopy, current)

				result = append(result, rowCopy)
				return
			}

			tName := tables[depth]

			for _, row := range tableRows[depth] {
				for colName, val := range row {
					current[tName+"."+colName] = val
				}

				combine(depth+1, current)

				for colName := range row {
					delete(current, tName+"."+colName)
				}
			}
		}
		combine(0, make(RowMap))
	}

	if len(meta["where_col"]) > 0 {
		col := meta["where_col"][0]
		op := meta["where_op"][0]
		val := meta["where_val"][0]

		filtered := []RowMap{}
		for _, row := range result {
			match, err := k.matchRow(row, col, op, val)
			if err != nil {
				return err
			}
			if match {
				filtered = append(filtered, row)
			}
		}
		result = filtered
	}

	if len(result) == 0 {
		fmt.Println("(0 rows)")
		return nil
	}

	headers := make([]string, len(selected))

	for i, sel := range selected {
		headers[i] = sel.TableName + "." + sel.ColumnName
	}

	fmt.Println(strings.Join(headers, "\t"))

	for _, row := range result {
		var vals []string
		for _, h := range headers {
			if v, ok := row[h]; ok && v != nil {
				vals = append(vals, fmt.Sprintf("%v", v))
			} else {
				vals = append(vals, "NULL")
			}
		}
		fmt.Println(strings.Join(vals, "\t"))
	}

	fmt.Printf("(%d rows)\n", len(result))

	return nil
}

func (k *Kernel) matchRow(row map[string]any, col string, op string, val string) (bool, error) {
	isString := strings.HasPrefix(val, "'") && strings.HasSuffix(val, "'")
	if isString && op != "=" {
		return false, fmt.Errorf("string values only support '=' operator, got '%s'", op)
	}
	cleanVal := val
	if isString {
		cleanVal = val[1 : len(val)-1]
	}

	var rowVal any
	var found bool
	for k, v := range row {
		parts := strings.Split(k, ".")
		if parts[len(parts)-1] == col || k == col {
			rowVal = v
			found = true
			break
		}
	}
	if !found {
		return false, nil
	}

	rowStr := fmt.Sprintf("%v", rowVal)
	switch op {
	case "=":
		return rowStr == cleanVal, nil
	case ">":
		return compareNumeric(rowStr, cleanVal) > 0, nil
	case "<":
		return compareNumeric(rowStr, cleanVal) < 0, nil
	case ">=":
		return compareNumeric(rowStr, cleanVal) >= 0, nil
	case "<=":
		return compareNumeric(rowStr, cleanVal) <= 0, nil
	case "<>":
		return rowStr != cleanVal, nil
	}
	return false, nil
}

func compareNumeric(a, b string) float64 {
	var na, nb float64
	fmt.Sscan(a, &na)
	fmt.Sscan(b, &nb)
	return na - nb
}

func (e *KernelError) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Msg)
}
