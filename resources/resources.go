package resources

import (
	"embed"
	"io/fs"
)

// Nuxt build files
var Nuxt fs.FS

//go:embed *
var resources embed.FS

func init() {
	var err error

	Nuxt, err = fs.Sub(resources, "dist")
	if err != nil {
		panic(err)
	}
}
