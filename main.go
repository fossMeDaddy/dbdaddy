package main

import (
	"dbdaddy/cmd"

	_ "github.com/jackc/pgx"
)

func main() {
	cmd.Execute()
}
