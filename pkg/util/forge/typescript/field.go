package typescript

import (
	"fmt"
	"strings"
)

// Field represents a TypeScript Interface field
type Field struct {
	Name     string
	Type     string
	Optional bool
	Null     bool
}

func (f Field) String() string {
	name := f.Name

	separator := ":"
	if f.Optional {
		separator = "?:"
	}

	typeString := f.Type
	if strings.HasSuffix(typeString, "[]") {
		typeString += "|null"
	} else if f.Null {
		typeString += "|null"
	}

	return fmt.Sprintf(
		"%s%s %s;",
		name,
		separator,
		typeString,
	)
}
