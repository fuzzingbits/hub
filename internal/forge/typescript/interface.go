package typescript

import (
	"fmt"
	"strings"
)

// Interface represents a TypeScript Interface
type Interface struct {
	Name   string
	Fields []Field
}

func (i Interface) String() string {
	if len(i.Fields) < 1 {
		return fmt.Sprintf(
			"export interface %s {}",
			i.Name,
		)
	}

	fieldStrings := []string{}
	for _, field := range i.Fields {
		fieldStrings = append(fieldStrings, field.String())
	}

	return fmt.Sprintf(
		"export interface %s {\n\t%s\n}",
		i.Name,
		strings.Join(fieldStrings, "\n\t"),
	)
}
