package main

import (
	"github.com/fuzzingbits/hub/internal/cmd"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	cmd.Run()
}
