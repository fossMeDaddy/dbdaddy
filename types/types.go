package types

import "github.com/spf13/cobra"

type CobraCmdFn = func(cmd *cobra.Command, args []string)

type MiddlewareFunc = func(fn CobraCmdFn) CobraCmdFn

type DbDumpFilesMap = map[string][]string

type DbRow struct {
	DataType string
	StrValue string
}

type DbRows = map[string][]DbRow
