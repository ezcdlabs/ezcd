//go:build dev
// +build dev

package public

func init() {
	IsEmbedded = false
}
