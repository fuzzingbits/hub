package typescript

import (
	"reflect"
	"regexp"
	"strings"
)

// Generate a list of TypeScript Interfaces
func Generate(goStructs []interface{}) []Interface {
	tsInterfaces := []Interface{}
	for _, goStruct := range goStructs {
		rv := reflect.ValueOf(goStruct)

		tsInterface := Interface{
			Name: rv.Type().Name(),
		}
		for i := 0; i < rv.NumField(); i++ {
			valueField := rv.Field(i)
			if !valueField.CanInterface() {
				continue
			}

			typeField := rv.Type().Field(i)
			tag := parseJSONFieldTag(typeField.Tag.Get("json"))
			if tag.Ignored {
				continue
			}

			tsInterface.Fields = append(
				tsInterface.Fields,
				genTsField(
					typeField.Name,
					typeField.Type.String(),
					tag,
				),
			)
		}

		tsInterfaces = append(tsInterfaces, tsInterface)
	}

	return tsInterfaces
}

func genTsField(goFieldName string, goFieldType string, tag jsonFieldTag) Field {
	tsField := Field{}
	if strings.HasPrefix(goFieldType, "[]") {
		tsField.Array = true
		goFieldType = strings.TrimPrefix(goFieldType, "[]")
	}

	if strings.HasPrefix(goFieldType, "*") {
		tsField.Null = true
		goFieldType = strings.TrimPrefix(goFieldType, "*")
	}

	tsField.Name = getTsFieldName(goFieldName, tag)
	tsField.Type = getTsFieldType(goFieldType, tag)
	if tag.Omitempty {
		tsField.Optional = true
	}

	return tsField
}

func getTsFieldName(goFieldName string, tag jsonFieldTag) string {
	if tag.NameOverride != "" {
		return tag.NameOverride
	}

	return goFieldName
}

var fieldTypeReplacer = regexp.MustCompile(`(^[\w]+\.)`)

func getTsFieldType(goFieldType string, tag jsonFieldTag) string {
	if tag.TypeOverride != "" {
		goFieldType = tag.TypeOverride
	}

	switch goFieldType {
	case "interface {}":
		return "any"
	case "int":
		return "number"
	case "uint":
		return "number"
	case "bool":
		return "boolean"
	case "time.Time":
		return "string"
	}

	goFieldType = fieldTypeReplacer.ReplaceAllString(goFieldType, "")

	return goFieldType
}
