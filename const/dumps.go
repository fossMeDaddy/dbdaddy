package constants

import (
	"fmt"
	"time"
)

func GetDumpFileName(dbname string) string {
	return fmt.Sprintf("%s__%s", time.Now().Local().Format("2006-01-02_15:04:05"), dbname)
}
