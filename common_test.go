package jpegstructure

import (
	"os"
	"path"

	"github.com/dsoprea/go-logging"
)

var (
	assetsPath = ""
)

func GetModuleRootPath() string {
	moduleRootPath := os.Getenv("JPEG_MODULE_ROOT_PATH")
	if moduleRootPath != "" {
		return moduleRootPath
	}

	currentWd, err := os.Getwd()
	log.PanicIf(err)

	currentPath := currentWd
	visited := make([]string, 0)

	for {
		tryStampFilepath := path.Join(currentPath, ".MODULE_ROOT")

		_, err := os.Stat(tryStampFilepath)
		if err != nil && os.IsNotExist(err) != true {
			log.Panic(err)
		} else if err == nil {
			break
		}

		visited = append(visited, tryStampFilepath)

		currentPath = path.Dir(currentPath)
		if currentPath == "/" {
			log.Panicf("could not find module-root: %v", visited)
		}
	}

	return currentPath
}

func init() {
	moduleRootPath := GetModuleRootPath()
	assetsPath = path.Join(moduleRootPath, "assets")
}
