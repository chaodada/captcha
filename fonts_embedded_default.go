package captcha

import "embed"

//go:embed fonts/*.ttf
//go:embed fonts/*.ttc
// defaultEmbeddedFontsFS Built-in font storage.
var defaultEmbeddedFontsFS embed.FS

var DefaultEmbeddedFonts = NewEmbeddedFontsStorage(defaultEmbeddedFontsFS)
