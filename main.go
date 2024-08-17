package main

import (
	"dbdaddy/cmd"

	_ "github.com/jackc/pgx/v5"
)

func main() {
	cmd.Execute()
}
