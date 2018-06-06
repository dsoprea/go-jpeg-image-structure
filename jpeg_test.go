package jpegstructure

import (
	"os"
	"path"
	"testing"
	"bufio"
	"bytes"
	"reflect"
    "fmt"

	"io/ioutil"

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

func Test_SegmentList_Write(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()

    filepath := path.Join(assetsPath, testImageRelFilepath)

    data, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    r := bytes.NewBuffer(data)

    sl, err := ParseSegments(r, len(data))
    log.PanicIf(err)

    b := new(bytes.Buffer)

	err = sl.Write(b)
	log.PanicIf(err)

	actual := b.Bytes()

	if bytes.Compare(actual, data) != 0 {
		t.Fatalf("output bytes do not equal input bytes")
	}
}

func Test_SegmentList_WriteReconstitutedExif(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()

    filepath := path.Join(assetsPath, testImageRelFilepath)

    sl, err := ParseFileStructure(filepath)
    log.PanicIf(err)

	_, s, rootIb, err := sl.ConstructExifBuilder()
	log.PanicIf(err)

	err = s.SetExif(rootIb)
	log.PanicIf(err)

	f, err := os.Create("/tmp/no_change_exif.jpg")
	log.PanicIf(err)

	defer f.Close()

	err = sl.Write(f)
	log.PanicIf(err)
}

func Test_SegmentList__UpdateExif(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()

    filepath := path.Join(assetsPath, testImageRelFilepath)

// TODO(dustin): !! Also test writing EXIF created from-scratch.
// TODO(dustin): !! Test adding a new EXIF (drop the existing).
// TODO(dustin): !! Might want to test a reconstruction without actually modifying anything. This is also useful. Everything will still be reallocated and this will help us determine if we're having parsing/encoding problems versions problems with an individual tag's value.
// TODO(dustin): !! Use native/third-party EXIF support to test?


    // Parse the image.

    sl, err := ParseFileStructure(filepath)
    log.PanicIf(err)


    // Update the UserComment tag.

	_, s, rootIb, err := sl.ConstructExifBuilder()
	log.PanicIf(err)

	i, err := rootIb.Find(exif.IfdExifId)
	log.PanicIf(err)

	exifBt := rootIb.Tags()[i]
	exifIb := exifBt.Value().Ib()


	uc := exif.TagUnknownType_9298_UserComment{
	    EncodingType: exif.TagUnknownType_9298_UserComment_Encoding_ASCII,
	    EncodingBytes: []byte("TEST COMMENT"),
	}

	err = exifIb.SetStandardWithName("UserComment", uc)
	log.PanicIf(err)


    // Update the exif segment.

	err = s.SetExif(rootIb)
	log.PanicIf(err)

    b := new(bytes.Buffer)

	err = sl.Write(b)
	log.PanicIf(err)

    recoveredBytes := b.Bytes()


    // Parse the re-encoded JPEG data and validate.

    recoveredSl, err := ParseBytesStructure(recoveredBytes)
    log.PanicIf(err)

    rootIfd, _, err := recoveredSl.Exif()
    log.PanicIf(err)

    exifIfd, err := rootIfd.ChildWithIfdIdentity(exif.ExifIi)
    log.PanicIf(err)

    results, err := exifIfd.FindTagWithName("UserComment")
    log.PanicIf(err)

    ucIte := results[0]

    if ucIte.TagId != 0x9286 {
        t.Fatalf("tag-ID not correct")
    }

    recoveredValueBytes, err := exifIfd.TagValueBytes(ucIte)
    log.PanicIf(err)

    expectedValueBytes := make([]byte, 0)

    expectedValueBytes = append(expectedValueBytes, []byte{ 'A', 'S', 'C', 'I', 'I', 0, 0, 0 }...)
    expectedValueBytes = append(expectedValueBytes, []byte("TEST COMMENT")...)

    if bytes.Compare(recoveredValueBytes, expectedValueBytes) != 0 {
        t.Fatalf("Recovered UserComment does not have the right value: %v != %v", recoveredValueBytes, expectedValueBytes)
    }
}

func ExampleSegment_SetExif() {
    filepath := path.Join(assetsPath, testImageRelFilepath)

    // Parse the image.

    sl, err := ParseFileStructure(filepath)
    log.PanicIf(err)


    // Update the UserComment tag.

    _, s, rootIb, err := sl.ConstructExifBuilder()
    log.PanicIf(err)

    i, err := rootIb.Find(exif.IfdExifId)
    log.PanicIf(err)

    exifBt := rootIb.Tags()[i]
    exifIb := exifBt.Value().Ib()


    uc := exif.TagUnknownType_9298_UserComment{
        EncodingType: exif.TagUnknownType_9298_UserComment_Encoding_ASCII,
        EncodingBytes: []byte("TEST COMMENT"),
    }

    err = exifIb.SetStandardWithName("UserComment", uc)
    log.PanicIf(err)


    // Update the exif segment.

    err = s.SetExif(rootIb)
    log.PanicIf(err)

    b := new(bytes.Buffer)

    err = sl.Write(b)
    log.PanicIf(err)

    updatedImageBytes := b.Bytes()
    updatedImageBytes = updatedImageBytes
    // Output:
}

func Test_SegmentList__FindExif(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()

    imageFilepath := path.Join(assetsPath, testImageRelFilepath)

    // Parse the image.

    sl, err := ParseFileStructure(imageFilepath)
    log.PanicIf(err)

    segmentNumber, s, err := sl.FindExif()
    log.PanicIf(err)

    if segmentNumber != 1 {
        t.Fatalf("exif not found in right position: (%d)", segmentNumber)
    }

    exifFilepath := fmt.Sprintf("%s.exif", imageFilepath)

    expectedExifBytes, err := ioutil.ReadFile(exifFilepath)
    log.PanicIf(err)

    if bytes.Compare(s.Data[6:], expectedExifBytes) != 0 {
        t.Fatalf("exif data not correct")
    }
}

func Test_SegmentList__Exif(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
        }
    }()

    imageFilepath := path.Join(assetsPath, testImageRelFilepath)

    // Parse the image.

    sl, err := ParseFileStructure(imageFilepath)
    log.PanicIf(err)

    rootIfd, s, err := sl.Exif()
    log.PanicIf(err)

    if rootIfd.Ii != exif.RootIi {
        t.Fatalf("root IFD does not have correct identity")
    }

    exifFilepath := fmt.Sprintf("%s.exif", imageFilepath)

    expectedExifBytes, err := ioutil.ReadFile(exifFilepath)
    log.PanicIf(err)

    if bytes.Compare(s.Data[6:], expectedExifBytes) != 0 {
        t.Fatalf("exif data not correct")
    }
}

func init() {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		log.Panicf("GOPATH is empty")
	}

	assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-jpeg-image-structure", "assets")
}
