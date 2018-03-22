package exifjpeg

import (
	// "fmt"
	"os"
	"path"
	"testing"

	"github.com/dsoprea/go-logging"
)

var (
	assetsPath           = ""
	testImageRelFilepath = "NDM_8901.jpg"
)

func TestSeekToNextSegment(t *testing.T) {
	filepath := path.Join(assetsPath, testImageRelFilepath)
	jn := NewJpegNavigator(filepath)
	defer jn.Close()

	err := jn.SeekToNextSegment()
	log.PanicIf(err)

	si, err := jn.ReadSegment()
	log.PanicIf(err)

	if si.MarkerId != 0xFB {
		t.Fatalf("Marker 0 not correct: (%X)", si.MarkerId)
	}

	err = jn.SeekToNextSegment()
	log.PanicIf(err)

	si, err = jn.ReadSegment()
	log.PanicIf(err)

	if si.MarkerId != 0x0 {
		t.Fatalf("Marker 1 not correct: (%X)", si.MarkerId)
	}
}

func TestVisitSegments(t *testing.T) {
	filepath := path.Join(assetsPath, testImageRelFilepath)
	jn := NewJpegNavigator(filepath)
	defer jn.Close()

	cb := func(markerId byte) (continue_ bool, err error) {
		// fmt.Printf("CB (%X)\n", markerId)

		return true, nil
	}

	err := jn.VisitSegments(cb)
	log.PanicIf(err)
}

func init() {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		log.Panicf("GOPATH is empty")
	}

	assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-exif-parser-jpeg", "assets")
}
