package jpegstructure

import (
	"path"

	"go/build"

	"github.com/dsoprea/go-logging"
)

var (
	assetsPath = ""
)

// GetModuleRootPath returns our source-path when running from source during
// tests.
func GetModuleRootPath() string {
	p, err := build.Default.Import(
		"github.com/dsoprea/go-jpeg-image-structure",
		build.Default.GOPATH,
		build.FindOnly)

	log.PanicIf(err)

	packagePath := p.Dir
	return packagePath
}

func init() {
	moduleRootPath := GetModuleRootPath()
	assetsPath = path.Join(moduleRootPath, "assets")
}
