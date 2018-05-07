package jpegstructure

import (
	"os"
	"path"
	"testing"
	"bufio"
	"bytes"
	"reflect"
	"fmt"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-exif"
)

var (
	assetsPath           = ""
	testImageRelFilepath = "NDM_8901.jpg"
)

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

func Test_JpegSplitter_Split(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()

	filepath := path.Join(assetsPath, testImageRelFilepath)
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
	buffer := []byte {}
	s.Buffer(buffer, int(size))

	s.Split(js.Split)

	for ; s.Scan() != false; { }

	if s.Err() != nil {
		log.PrintError(s.Err())
		t.Fatalf("error while scanning: %v", s.Err())
	}

	expectedMarkers := []byte { 0xd8, 0xe1, 0xe1, 0xdb, 0xc0, 0xc4, 0xda, 0x00, 0xd9 }

	if bytes.Compare(v.markerList, expectedMarkers) != 0 {
		t.Fatalf("Markers found are not correct: %v\n", DumpBytesToString(v.markerList))
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

func Test_SegmentList__Update_Exif(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()

    filepath := path.Join(assetsPath, testImageRelFilepath)

// TODO(dustin): !! Also test writing EXIF created from-scratch.
// TODO(dustin): !! Test adding a new EXIF (drop the existing).
// TODO(dustin): !! Test (sl).SetExif .

    sl, err := ParseFileStructure(filepath)
    log.PanicIf(err)


// TODO(dustin): !! Confirm tags before.
	_, s, tags, err := sl.DumpExif()
	log.PanicIf(err)

	s = s

	fmt.Printf("\n")
	fmt.Printf("BEFORE:\n")
	fmt.Printf("\n")

	for i, tag := range tags {
		fmt.Printf("%02d: %v\n", i, tag)
	}

	_, s, rootIb, err := sl.ConstructExifBuilder()
	log.PanicIf(err)

	// fmt.Printf("\n")
	// fmt.Printf("IB:\n")
	// fmt.Printf("\n")

	// fmt.Printf("%v\n", rootIb)

	i, err := rootIb.Find(exif.IfdExifId)
	log.PanicIf(err)

	exifBt := rootIb.Tags()[i]
	exifIb := exifBt.Value().Ib()


	uc := exif.TagUnknownType_9298_UserComment{
	    EncodingType: exif.TagUnknownType_9298_UserComment_Encoding_ASCII,
	    EncodingBytes: []byte("TEST EXIF CHANGE"),
	}

	err = exifIb.SetFromConfigWithName("UserComment", uc)
	log.PanicIf(err)



// TODO(dustin): !! Might want to test a reconstruction without actually modifying anything. This is also useful. Everything will still be reallocated and this will help us determine if we're having parsing/encoding problems versions problems with an individual tag's value.

// TODO(dustin): !! The output doesn't have the thumbnail(s).

// TODO(dustin): !! We think we're writing the original IFD *and* our updated IFD, one after the other. *Tee IFD1 data is totally different.*

	fmt.Printf("\n")
	fmt.Printf("IB TO WRITE:\n")
	fmt.Printf("\n")

	rootIb.Dump()


	err = s.SetExif(rootIb)
	log.PanicIf(err)


// TODO(dustin): !! Confirm tags after.
 	_, s, tags, err = sl.DumpExif()
 	log.PanicIf(err)

 	s = s

	fmt.Printf("\n")
	fmt.Printf("AFTER:\n")
	fmt.Printf("\n")

	for i, tag := range tags {
		fmt.Printf("%02d: %s\n", i, tag)
	}

// TODO(dustin): !! Use native/third-party EXIF support to test?

	f, err := os.Create("/tmp/updated_exif.jpg")
	log.PanicIf(err)

	defer f.Close()

	err = sl.Write(f)
	log.PanicIf(err)
}

func init() {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		log.Panicf("GOPATH is empty")
	}

	assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-jpeg-image-structure", "assets")
}
