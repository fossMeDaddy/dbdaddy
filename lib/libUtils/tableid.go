package libUtils

import (
	"fmt"
	"strings"
)

func GetTableId(schemaname, tablename string) string {
	if len(schemaname) > 0 && len(tablename) > 0 {
		return fmt.Sprintf("%s.%s", schemaname, tablename)
	} else if len(tablename) > 0 {
		return tablename
	}

	return ""
}

// returns (schemaname, tablename) after taking in "schemaname.tablename"
func GetTableFromId(tableid string) (string, string) {
	table := strings.Split(tableid, ".")
	if len(table) < 2 {
		return "", ""
	}
	return table[0], table[1]
}
