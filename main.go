package main

import (
	"github.com/fuzzingbits/hub/internal/cmd"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	cmd.Run()
}
