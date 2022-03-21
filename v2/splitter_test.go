package jpegstructure

import (
	"bufio"
	"bytes"
	"os"
	"reflect"
	"testing"

	log "github.com/dsoprea/go-logging"
)

type collectorVisitor struct {
	markerList []byte
	sofList    []SofSegment
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

func Test_JpegSplitter_Split(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
			t.Fatalf("Test failure.")
		}
	}()

	filepath := GetTestImageFilepath()

	f, err := os.Open(filepath)
	log.PanicIf(err)

	defer f.Close()

	stat, err := f.Stat()
	log.PanicIf(err)

	size := stat.Size()

	v := new(collectorVisitor)
	js := NewJpegSplitter(v)

	s := bufio.NewScanner(f)

	// Since each segment can be any size, our buffer must allowed to grow as
	// large as the file.
	buffer := []byte{}
	s.Buffer(buffer, int(size))

	s.Split(js.Split)

	for s.Scan() != false {
	}

	if s.Err() != nil {
		log.PrintError(s.Err())
		t.Fatalf("error while scanning: %v", s.Err())
	}

	expectedMarkers := []byte{0xd8, 0xe1, 0xe1, 0xe1, 0xdb, 0xc0, 0xc4, 0xda, 0x00, 0xd9}

	if bytes.Compare(v.markerList, expectedMarkers) != 0 {
		t.Fatalf("Markers found are not correct: %v\n", DumpBytesToString(v.markerList))
	}

	expectedSofList := []SofSegment{
		{
			BitsPerSample:  8,
			Width:          3840,
			Height:         2560,
			ComponentCount: 3,
		},
		{
			BitsPerSample:  0,
			Width:          1281,
			Height:         1,
			ComponentCount: 1,
		},
	}

	if reflect.DeepEqual(v.sofList, expectedSofList) == false {
		t.Fatalf("SOF segments not equal: %v\n", v.sofList)
	}
}
