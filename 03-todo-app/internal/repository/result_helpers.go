package repository

import (
	"github.com/chameleon-db/chameleondb/chameleon/pkg/engine"
)

func rowToMap(row engine.Row) map[string]interface{} {
	return map[string]interface{}(row)
}

func rowsToMaps(rows []engine.Row) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		out = append(out, rowToMap(row))
	}
	return out
}
