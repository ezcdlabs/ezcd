//go:build !dev
// +build !dev

package public

import (
	"embed"
	"io/fs"
)

//go:embed web-dist/*
var ebox embed.FS

func init() {
	IsEmbedded = true

	sub, err := fs.Sub(ebox, "web-dist")
	if err != nil {
		panic(err)
	}

	Box = sub
}
