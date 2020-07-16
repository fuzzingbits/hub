package typescript

import "strings"

type jsonFieldTag struct {
	NameOverride string
	TypeOverride string
	Omitempty    bool
	Ignored      bool
}

func parseJSONFieldTag(tagString string) jsonFieldTag {
	if tagString == "-" {
		return jsonFieldTag{
			Ignored: true,
		}
	}

	parts := strings.Split(tagString, ",")
	if len(parts) == 1 {
		return jsonFieldTag{
			NameOverride: parts[0],
		}
	}

	tag := jsonFieldTag{}
	for i, part := range parts {
		if i == 0 {
			tag.NameOverride = part
			continue
		}

		if strings.Contains(part, "omitempty") {
			tag.Omitempty = true
		} else {
			tag.TypeOverride = part
		}
	}

	return tag
}
