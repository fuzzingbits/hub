//go:generate go run ./ops/builder
package main

import (
	"github.com/fuzzingbits/hub/pkg/cmd"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	cmd.Run()
}
