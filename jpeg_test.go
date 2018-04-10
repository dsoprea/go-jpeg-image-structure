package exifjpeg

import (
	"os"
	"path"
	"testing"
	"bufio"
	"bytes"
	"reflect"

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

type collectorVisitor struct {
	markerList []byte
	sofList []SofSegment
}

func (v *collectorVisitor) HandleSegment(lastMarkerId byte, lastMarkerName string, counter int, lastIsScanData bool) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	v.markerList = append(v.markerList, lastMarkerId)

	return nil
}

func (v *collectorVisitor) HandleSof(sof *SofSegment) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	v.sofList = append(v.sofList, *sof)

	return nil
}

func TestJpegSplitterSplit(t *testing.T) {

	filepath := path.Join(assetsPath, testImageRelFilepath)
	f, err := os.Open(filepath)
	log.PanicIf(err)

	stat, err := f.Stat()
	log.PanicIf(err)

	size := stat.Size()

	v := new(collectorVisitor)
	js := NewJpegSplitter(v)

	s := bufio.NewScanner(f)

	// Since each segment can be any size, our buffer must allowed to grow as
	// large as the file.
	buffer := []byte {}
	s.Buffer(buffer, int(size))

	s.Split(js.Split)

	for ; s.Scan() != false; { }

	if s.Err() != nil {
		t.Fatalf("error while scanning: %v", s.Err())
	}

	expectedMarkers := []byte { 0xd8, 0xe1, 0xe1, 0xdb, 0xc0, 0xc4, 0xda, 0x00, 0xd9 }

	if bytes.Compare(v.markerList, expectedMarkers) != 0 {
		t.Fatalf("Markers found are not correct: %v\n", v.markerList)
	}

	expectedSofList := []SofSegment {
		SofSegment{
			BitsPerSample: 8,
			Width: 3840,
			Height: 2560,
			ComponentCount: 3,
		},
		SofSegment{
			BitsPerSample: 0,
			Width: 1281,
			Height: 1,
			ComponentCount: 1,
		},
	}

	if reflect.DeepEqual(v.sofList, expectedSofList) == false {
		t.Fatalf("SOF segments not equal: %v\n", v.sofList)
	}
}

func TestIsExif_True(t *testing.T) {
	js := NewJpegSplitter(nil)

	if ok := js.isExif([]byte("Exif\000\000")); ok != true {
		t.Fatalf("expected true")
	}
}

func TestIsExif_False(t *testing.T) {
	js := NewJpegSplitter(nil)

	if ok := js.isExif([]byte("something unexpected")); ok != false {
		t.Fatalf("expected false")
	}
}

func init() {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		log.Panicf("GOPATH is empty")
	}

	assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-exif-parser-jpeg", "assets")
}
