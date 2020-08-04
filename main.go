//go:generate go run ./ops/builder
//go:generate sh -c "[ -d ./node_modules ] && npm run fmt || exit 0"
package main

import (
	"github.com/fuzzingbits/hub/pkg/cmd"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	cmd.Run()
}
