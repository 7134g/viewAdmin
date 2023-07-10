package table

import (
	"fmt"
	"github.com/7134g/viewAdmin/internel/serve"
	"github.com/7134g/viewAdmin/internel/view"
	"strings"
)

type ViewTable struct {
	cfg *view.Config
}

func NewViewTableLogic(c *view.Config) ViewTable {
	return ViewTable{cfg: c}
}

func (h *ViewTable) ViewTable(_ *serve.BaseContext) (interface{}, error) {
	tablesName := make([]string, 0)
	err := h.cfg.DB.Raw(h.cfg.SqlTableScript).Scan(&tablesName).Error
	if err != nil {
		return nil, err
	}

	table := make(map[string]map[string]interface{})
	for _, name := range tablesName {
		var m map[string]interface{}
		sqlScript := fmt.Sprintf(`SHOW CREATE TABLE %s`, name)
		err := h.cfg.DB.Raw(sqlScript).Scan(&m).Error
		if err != nil {
			return nil, err
		}
		m = parseTable(m)
		table[name] = m
	}

	return table, err
}

func parseTable(m map[string]interface{}) map[string]interface{} {
	createSql, ok := m["Create Table"]
	if !ok {
		return nil
	}
	v, ok := createSql.(string)
	if !ok {
		return nil
	}

	result := make(map[string]interface{})
	lines := strings.Split(v, "\n")
	for i := 1; i < len(lines)-1; i++ {
		line := lines[i]

		if skip(line) {
			continue
		}

		line = strings.TrimLeft(line, " ")
		line = strings.TrimRight(line, " ")
		fields := strings.Split(line, " ")
		key := strings.ReplaceAll(fields[0], "`", "")
		value := fields[1]
		result[key] = value
	}
	return result
}

func skip(line string) bool {
	switch {
	case strings.Contains(line, "UNIQUE"),
		strings.Contains(line, "PRIMARY"),
		strings.Contains(line, "CONSTRAINT"),
		strings.Contains(line, "KEY"),
		strings.Contains(line, "UNIQUE"),
		strings.Contains(line, "CONSTRAINT"):
		return true
	default:
		return false
	}
}
