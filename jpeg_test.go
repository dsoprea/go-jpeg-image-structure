package jpegstructure

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-exif"
	"github.com/dsoprea/go-logging"
)

var (
	testImageRelFilepath = "NDM_8901.jpg"
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
	buffer := []byte{}
	s.Buffer(buffer, int(size))

	s.Split(js.Split)

	for s.Scan() != false {
	}

	if s.Err() != nil {
		log.PrintError(s.Err())
		t.Fatalf("error while scanning: %v", s.Err())
	}

	expectedMarkers := []byte{0xd8, 0xe1, 0xe1, 0xdb, 0xc0, 0xc4, 0xda, 0x00, 0xd9}

	if bytes.Compare(v.markerList, expectedMarkers) != 0 {
		t.Fatalf("Markers found are not correct: %v\n", DumpBytesToString(v.markerList))
	}

	expectedSofList := []SofSegment{
		SofSegment{
			BitsPerSample:  8,
			Width:          3840,
			Height:         2560,
			ComponentCount: 3,
		},
		SofSegment{
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

	jmp := NewJpegMediaParser()

	sl, err := jmp.Parse(r, len(data))
	log.PanicIf(err)

	b := new(bytes.Buffer)

	err = sl.Write(b)
	log.PanicIf(err)

	actual := b.Bytes()

	if bytes.Compare(actual, data) != 0 {
		t.Fatalf("output bytes do not equal input bytes")
	}
}

// func Test_SegmentList_WriteReconstitutedExif(t *testing.T) {
//     defer func() {
//         if state := recover(); state != nil {
//             err := log.Wrap(state.(error))
//             log.PrintErrorf(err, "Test failure.")
//         }
//     }()

//     filepath := path.Join(assetsPath, testImageRelFilepath)

//     jmp := NewJpegMediaParser()

//     sl, err := ParseFileStructure(filepath)
//     log.PanicIf(err)

// 	_, s, rootIb, err := sl.ConstructExifBuilder()
// 	log.PanicIf(err)

// 	err = s.SetExif(rootIb)
// 	log.PanicIf(err)

// 	f, err := os.Create("/tmp/no_change_exif.jpg")
// 	log.PanicIf(err)

// 	defer f.Close()

// 	err = sl.Write(f)
// 	log.PanicIf(err)
// }

func Test_Segment_SetExif(t *testing.T) {
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

	jmp := NewJpegMediaParser()

	sl, err := jmp.ParseFile(filepath)
	log.PanicIf(err)

	// Update the UserComment tag.

	rootIb, err := sl.ConstructExifBuilder()
	log.PanicIf(err)

	i, err := rootIb.Find(exif.IfdExifId)
	log.PanicIf(err)

	exifBt := rootIb.Tags()[i]
	exifIb := exifBt.Value().Ib()

	uc := exif.TagUnknownType_9298_UserComment{
		EncodingType:  exif.TagUnknownType_9298_UserComment_Encoding_ASCII,
		EncodingBytes: []byte("TEST COMMENT"),
	}

	err = exifIb.SetStandardWithName("UserComment", uc)
	log.PanicIf(err)

	// Update the exif segment.

	_, s, err := sl.FindExif()
	log.PanicIf(err)

	err = s.SetExif(rootIb)
	log.PanicIf(err)

	b := new(bytes.Buffer)

	err = sl.Write(b)
	log.PanicIf(err)

	recoveredBytes := b.Bytes()

	// Parse the re-encoded JPEG data and validate.

	recoveredSl, err := jmp.ParseBytes(recoveredBytes)
	log.PanicIf(err)

	rootIfd, _, err := recoveredSl.Exif()
	log.PanicIf(err)

	exifIfd, err := rootIfd.ChildWithIfdPath(exif.IfdPathStandardExif)
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

	expectedValueBytes = append(expectedValueBytes, []byte{'A', 'S', 'C', 'I', 'I', 0, 0, 0}...)
	expectedValueBytes = append(expectedValueBytes, []byte("TEST COMMENT")...)

	if bytes.Compare(recoveredValueBytes, expectedValueBytes) != 0 {
		t.Fatalf("Recovered UserComment does not have the right value: %v != %v", recoveredValueBytes, expectedValueBytes)
	}
}

func Test_SegmentList_SetExif(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
		}
	}()

	initialSegments := []*Segment{
		&Segment{MarkerId: 0},
		&Segment{MarkerId: 0},
	}

	sl := NewSegmentList(initialSegments)

	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	ib := exif.NewIfdBuilder(im, ti, exif.IfdPathStandard, exif.TestDefaultByteOrder)
	ib.AddStandardWithName("ProcessingSoftware", "some software")

	err := sl.SetExif(ib)
	log.PanicIf(err)

	exifSegment := sl.Segments()[2]

	if exifSegment.MarkerId != MARKER_APP1 {
		t.Fatalf("New segment is not correct.")
	} else if len(exifSegment.Data) == 0 {
		t.Fatalf("New segment does not have data.")
	}

	originalSegment := exifSegment
	originalData := exifSegment.Data

	sl.Add(&Segment{MarkerId: 0})
	sl.Add(&Segment{MarkerId: 0})

	ib = exif.NewIfdBuilder(im, ti, exif.IfdPathStandard, exif.TestDefaultByteOrder)
	ib.AddStandardWithName("ProcessingSoftware", "some software2")

	err = sl.SetExif(ib)
	log.PanicIf(err)

	exifSegment = sl.Segments()[2]

	if len(sl.Segments()) != 5 {
		t.Fatalf("Segment count not correct.")
	} else if exifSegment != originalSegment {
		// The data should change, not the segment itself.

		t.Fatalf("EXIF segment has been changed.")
	} else if exifSegment.MarkerId != MARKER_APP1 {
		t.Fatalf("EXIF segment is not correct.")
	} else if bytes.Compare(exifSegment.Data, originalData) == 0 {
		t.Fatalf("EXIF segment has not changed.")
	}
}

func ExampleSegmentList_SetExif() {
	filepath := path.Join(assetsPath, testImageRelFilepath)

	// Parse the image.

	jmp := NewJpegMediaParser()

	sl, err := jmp.ParseFile(filepath)
	log.PanicIf(err)

	// Update the UserComment tag.

	rootIb, err := sl.ConstructExifBuilder()
	log.PanicIf(err)

	i, err := rootIb.Find(exif.IfdExifId)
	log.PanicIf(err)

	exifBt := rootIb.Tags()[i]
	exifIb := exifBt.Value().Ib()

	uc := exif.TagUnknownType_9298_UserComment{
		EncodingType:  exif.TagUnknownType_9298_UserComment_Encoding_ASCII,
		EncodingBytes: []byte("TEST COMMENT"),
	}

	err = exifIb.SetStandardWithName("UserComment", uc)
	log.PanicIf(err)

	// Update the exif segment.

	err = sl.SetExif(rootIb)
	log.PanicIf(err)

	b := new(bytes.Buffer)

	err = sl.Write(b)
	log.PanicIf(err)

	updatedImageBytes := b.Bytes()
	updatedImageBytes = updatedImageBytes
	// Output:
}

func Test_SegmentList_FindExif(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
		}
	}()

	imageFilepath := path.Join(assetsPath, testImageRelFilepath)

	// Parse the image.

	jmp := NewJpegMediaParser()

	sl, err := jmp.ParseFile(imageFilepath)
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

func Test_SegmentList_Exif(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
		}
	}()

	imageFilepath := path.Join(assetsPath, testImageRelFilepath)

	// Parse the image.

	jmp := NewJpegMediaParser()

	sl, err := jmp.ParseFile(imageFilepath)
	log.PanicIf(err)

	rootIfd, data, err := sl.Exif()
	log.PanicIf(err)

	if rootIfd.IfdPath != exif.IfdPathStandard {
		t.Fatalf("root IFD does not have correct identity")
	}

	exifFilepath := fmt.Sprintf("%s.exif", imageFilepath)

	expectedExifBytes, err := ioutil.ReadFile(exifFilepath)
	log.PanicIf(err)

	if bytes.Compare(data, expectedExifBytes) != 0 {
		t.Fatalf("exif data not correct")
	}
}

func TestSegmentList_Validate(t *testing.T) {
	filepath := path.Join(assetsPath, testImageRelFilepath)

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	segments := []*Segment{
		&Segment{
			MarkerId: 0x0,
			Offset:   0x0,
		},
	}

	sl := NewSegmentList(segments)

	err = sl.Validate(data)
	if err == nil {
		t.Fatalf("Expected error about missing minimum segments.")
	} else if err.Error() != "minimum segments not found" {
		log.Panic(err)
	}
}
