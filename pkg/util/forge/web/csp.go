package web

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"golang.org/x/net/html"
)

// CSPEntries stores various values to be converted into the CSP Header value
type CSPEntries struct {
	Default []string
	Script  []string
	Image   []string
	Style   []string
}

func (c CSPEntries) String() string {
	cspParts := []string{}

	if len(c.Default) > 0 {
		cspParts = append(cspParts, fmt.Sprintf("default-src %s", strings.Join(c.Default, " ")))
	}

	if len(c.Script) > 0 {
		cspParts = append(cspParts, fmt.Sprintf("script-src %s", strings.Join(c.Script, " ")))
	}

	if len(c.Style) > 0 {
		cspParts = append(cspParts, fmt.Sprintf("style-src %s", strings.Join(c.Style, " ")))
	}
	if len(c.Image) > 0 {
		cspParts = append(cspParts, fmt.Sprintf("img-src %s", strings.Join(c.Image, " ")))
	}

	return strings.Join(cspParts, "; ")
}

// GenerateContentSecurityPolicy for the given html file contents
func GenerateContentSecurityPolicy(fileContents []byte, cspEntries CSPEntries) string {
	doc, _ := html.Parse(bytes.NewReader(fileContents))

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.FirstChild != nil {
			switch n.Data {
			case "script":
				cspEntries.Script = append(cspEntries.Script, generateNodeHash(n.FirstChild))
			case "style":
				cspEntries.Style = append(cspEntries.Style, generateNodeHash(n.FirstChild))
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	return cspEntries.String()
}

func generateNodeHash(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)

	_ = html.Render(w, n)

	str := html.UnescapeString(buf.String())
	sum := sha256.Sum256([]byte(str))
	sha1Hash := base64.StdEncoding.EncodeToString(sum[:])

	return fmt.Sprintf("'sha256-%s'", sha1Hash)
}
