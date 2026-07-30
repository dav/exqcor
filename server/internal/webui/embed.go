// Package webui embeds the built Svelte SPA. The dist directory is produced
// by `npm run build` in web/ (see the Makefile); everything in it except the
// placeholder index.html is gitignored.
package webui

import (
	"embed"
	"io/fs"
)

//go:embed all:dist
var dist embed.FS

// FS returns the SPA files rooted at dist/.
func FS() fs.FS {
	sub, err := fs.Sub(dist, "dist")
	if err != nil {
		panic(err)
	}
	return sub
}
