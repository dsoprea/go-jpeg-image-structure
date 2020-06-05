package jpegstructure

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"io/ioutil"

	"github.com/dsoprea/go-exif/v2"
	"github.com/dsoprea/go-exif/v2/common"
	"github.com/dsoprea/go-exif/v2/undefined"
	"github.com/dsoprea/go-logging"
)

func TestSegmentList_Write(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
			t.Fatalf("Test failure.")
		}
	}()

	filepath := GetTestImageFilepath()

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	br := bytes.NewReader(data)

	jmp := NewJpegMediaParser()

	intfc, err := jmp.Parse(br, len(data))
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	b := new(bytes.Buffer)

	err = sl.Write(b)
	log.PanicIf(err)

	actual := b.Bytes()

	if bytes.Compare(actual, data) != 0 {
		t.Fatalf("output bytes do not equal input bytes")
	}
}

// func TestSegmentList_WriteReconstitutedExif(t *testing.T) {
//     defer func() {
//         if state := recover(); state != nil {
//             err := log.Wrap(state.(error))
//             log.PrintErrorf(err, "Test failure.")
//             t.Fatalf("Test failure.")
//         }
//     }()

//     filepath := GetTestImageFilepath()

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

func TestSegmentList_SetExif_FromScratch(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
			t.Fatalf("Test failure.")
		}
	}()

	// Parse the image.

	filepath := GetTestImageFilepath()

	jmp := NewJpegMediaParser()

	intfc, err := jmp.ParseFile(filepath)
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	// Make sure we don't start out with EXIF data.

	wasDropped, err := sl.DropExif()
	log.PanicIf(err)

	if wasDropped != true {
		t.Fatalf("Expected the EXIF segment to be dropped, but it wasn't.")
	}

	// Set the ProcessingSoftware tag.

	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	err = exif.LoadStandardTags(ti)
	log.PanicIf(err)

	rootIb := exif.NewIfdBuilder(im, ti, exifcommon.IfdStandardIfdIdentity, exifcommon.EncodeDefaultByteOrder)

	err = rootIb.AddStandardWithName("ProcessingSoftware", "some software")
	log.PanicIf(err)

	err = sl.SetExif(rootIb)
	log.PanicIf(err)

	b := new(bytes.Buffer)

	err = sl.Write(b)
	log.PanicIf(err)

	recoveredBytes := b.Bytes()

	// Parse the re-encoded JPEG data and validate.

	recoveredIntfc, err := jmp.ParseBytes(recoveredBytes)
	log.PanicIf(err)

	recoveredSl := recoveredIntfc.(*SegmentList)

	rootIfd, _, err := recoveredSl.Exif()
	log.PanicIf(err)

	results, err := rootIfd.FindTagWithName("ProcessingSoftware")
	log.PanicIf(err)

	ucIte := results[0]

	if ucIte.TagId() != 0x000b {
		t.Fatalf("tag-ID not correct")
	}

	recoveredValueRaw, err := ucIte.Value()
	log.PanicIf(err)

	recoveredValue := recoveredValueRaw.(string)
	if recoveredValue != "some software" {
		t.Fatalf("Value of tag not correct: [%s]", recoveredValue)
	}
}

func TestSegmentList_SetExif(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
			t.Fatalf("Test failure.")
		}
	}()

	initialSegments := []*Segment{
		&Segment{MarkerId: 0},
		&Segment{MarkerId: 0},
	}

	sl := NewSegmentList(initialSegments)

	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	ib := exif.NewIfdBuilder(im, ti, exifcommon.IfdStandardIfdIdentity, exifcommon.TestDefaultByteOrder)
	ib.AddStandardWithName("ProcessingSoftware", "some software")

	err := sl.SetExif(ib)
	log.PanicIf(err)

	exifSegment := sl.Segments()[1]

	if exifSegment.MarkerId != MARKER_APP1 {
		t.Fatalf("New segment is not correct.")
	} else if len(exifSegment.Data) == 0 {
		t.Fatalf("New segment does not have data.")
	}

	originalSegment := exifSegment
	originalData := exifSegment.Data

	sl.Add(&Segment{MarkerId: 0})
	sl.Add(&Segment{MarkerId: 0})

	ib = exif.NewIfdBuilder(im, ti, exifcommon.IfdStandardIfdIdentity, exifcommon.TestDefaultByteOrder)
	ib.AddStandardWithName("ProcessingSoftware", "some software2")

	err = sl.SetExif(ib)
	log.PanicIf(err)

	exifSegment = sl.Segments()[1]

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

func ExampleSegmentList_SetExif_unknowntype() {
	filepath := GetTestImageFilepath()

	// Parse the image.

	jmp := NewJpegMediaParser()

	intfc, err := jmp.ParseFile(filepath)
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	// Update the UserComment tag.

	rootIb, err := sl.ConstructExifBuilder()
	log.PanicIf(err)

	ifdPath := "IFD/Exif"

	exifIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPath)
	log.PanicIf(err)

	uc := exifundefined.Tag9286UserComment{
		EncodingType:  exifundefined.TagUndefinedType_9286_UserComment_Encoding_ASCII,
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

// ExampleSegmentList_SetExif shows how to construct a chain of
// `IfdBuilder` structs for the existing IFDs, identify the builder for the IFD
// that we know hosts the tag we want to change, and how to change it.
func ExampleSegmentList_SetExif() {
	filepath := GetTestImageFilepath()

	// Parse the image.

	jmp := NewJpegMediaParser()

	intfc, err := jmp.ParseFile(filepath)
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	// Update the UserComment tag.

	rootIb, err := sl.ConstructExifBuilder()
	log.PanicIf(err)

	ifdPath := "IFD0"

	ifdIb, err := exif.GetOrCreateIbFromRootIb(rootIb, ifdPath)
	log.PanicIf(err)

	now := time.Now().UTC()
	updatedTimestampPhrase := exif.ExifFullTimestampString(now)

	err = ifdIb.SetStandardWithName("DateTime", updatedTimestampPhrase)
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

func TestSegmentList_FindExif(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
			t.Fatalf("Test failure.")
		}
	}()

	imageFilepath := GetTestImageFilepath()

	// Parse the image.

	jmp := NewJpegMediaParser()

	intfc, err := jmp.ParseFile(imageFilepath)
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	segmentNumber, s, err := sl.FindExif()
	log.PanicIf(err)

	if segmentNumber != 1 {
		t.Fatalf("exif not found in right position: (%d)", segmentNumber)
	}

	exifFilepath := fmt.Sprintf("%s.just_exif", imageFilepath)

	expectedExifBytes, err := ioutil.ReadFile(exifFilepath)
	log.PanicIf(err)

	if bytes.Compare(s.Data[6:], expectedExifBytes) != 0 {
		t.Fatalf("exif data not correct")
	}
}

func TestSegmentList_Exif(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
			t.Fatalf("Test failure.")
		}
	}()

	imageFilepath := GetTestImageFilepath()

	// Parse the image.

	jmp := NewJpegMediaParser()

	intfc, err := jmp.ParseFile(imageFilepath)
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	rootIfd, data, err := sl.Exif()
	log.PanicIf(err)

	if rootIfd.IfdIdentity().Equals(exifcommon.IfdStandardIfdIdentity) != true {
		t.Fatalf("root IFD does not have correct identity")
	}

	exifFilepath := fmt.Sprintf("%s.just_exif", imageFilepath)

	expectedExifBytes, err := ioutil.ReadFile(exifFilepath)
	log.PanicIf(err)

	if bytes.Compare(data, expectedExifBytes) != 0 {
		t.Fatalf("exif data not correct")
	}
}

func TestSegmentList_Validate(t *testing.T) {
	filepath := GetTestImageFilepath()

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
