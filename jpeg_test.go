package exifjpeg

import (
	"fmt"
	"os"
	"path"
	"testing"
	"bufio"

	"github.com/dsoprea/go-logging"
)

var (
	assetsPath           = ""
	testImageRelFilepath = "NDM_8901.jpg"
)

// func TestVerifyIsJpeg(t *testing.T) {
// 	filepath := path.Join(assetsPath, testImageRelFilepath)
// 	jn := NewJpegNavigator(filepath)
// 	defer jn.Close()
// }

type Visitor struct {

}

func (v *Visitor) HandleSegment(lastMarkerId byte, lastMarkerName string, counter int, lastIsScanData bool) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	fmt.Printf("VISIT (%02X)\n", lastMarkerId)

	return nil
}

func TestJpegSplitterSplit(t *testing.T) {

	filepath := path.Join(assetsPath, testImageRelFilepath)
	f, err := os.Open(filepath)
	log.PanicIf(err)

	stat, err := f.Stat()
	log.PanicIf(err)

	size := stat.Size()

	v := new(Visitor)
	js := NewJpegSplitter(v)

	s := bufio.NewScanner(f)

	// Since each segment can be any size, our buffer must allowed to grow as
	// large as the file.
	buffer := []byte {}
	s.Buffer(buffer, int(size))

	s.Split(js.Split)

	// more := s.Scan()
	// if more != true || s.Err() != nil {
	// 	t.Fatalf("more tokens expected (1): %v", s.Err())
	// }

	// fmt.Printf("MARKER1: %02X\n", js.MarkerId())

	// more = s.Scan()
	// if more != true || s.Err() != nil {
	// 	t.Fatalf("more tokens expected (2): %v", s.Err())
	// }

	// fmt.Printf("MARKER2: %02X\n", js.MarkerId())

	// more = s.Scan()
	// if more != true || s.Err() != nil {
	// 	t.Fatalf("more tokens expected (3): %v", s.Err())
	// }

	// fmt.Printf("MARKER3: %02X\n", js.MarkerId())

	// more = s.Scan()
	// if more != true || s.Err() != nil {
	// 	t.Fatalf("more tokens expected (4): %v", s.Err())
	// }

	// fmt.Printf("MARKER4: %02X\n", js.MarkerId())

	for ; s.Scan() != false; {
		fmt.Printf("Marker-ID: %02X\n", js.MarkerId())
	}

	fmt.Printf("Scan finished.\n")

	log.PanicIf(s.Err())

	fmt.Printf("No errors.\n")
}

func init() {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		log.Panicf("GOPATH is empty")
	}

	assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-exif-parser-jpeg", "assets")
}
