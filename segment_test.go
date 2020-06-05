package jpegstructure

import (
	"bytes"
	"fmt"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-exif/v2"
	"github.com/dsoprea/go-exif/v2/common"
	"github.com/dsoprea/go-exif/v2/undefined"
	"github.com/dsoprea/go-logging"
)

func TestSegment_SetExif_Update(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
			t.Fatalf("Test failure.")
		}
	}()

	filepath := GetTestImageFilepath()

	// TODO(dustin): !! Might want to test a reconstruction without actually modifying anything. This is also useful. Everything will still be reallocated and this will help us determine if we're having parsing/encoding problems versions problems with an individual tag's value.
	// TODO(dustin): !! Use native/third-party EXIF support to test?

	// Parse the image.

	jmp := NewJpegMediaParser()

	intfc, err := jmp.ParseFile(filepath)
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	// Update the UserComment tag.

	rootIb, err := sl.ConstructExifBuilder()
	log.PanicIf(err)

	i, err := rootIb.Find(exifcommon.IfdExifStandardIfdIdentity.TagId())
	log.PanicIf(err)

	exifBt := rootIb.Tags()[i]
	exifIb := exifBt.Value().Ib()

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

	recoveredBytes := b.Bytes()

	// Parse the re-encoded JPEG data and validate.

	recoveredIntfc, err := jmp.ParseBytes(recoveredBytes)
	log.PanicIf(err)

	recoveredSl := recoveredIntfc.(*SegmentList)

	rootIfd, _, err := recoveredSl.Exif()
	log.PanicIf(err)

	exifIfd, err := rootIfd.ChildWithIfdPath(exifcommon.IfdExifStandardIfdIdentity)
	log.PanicIf(err)

	results, err := exifIfd.FindTagWithName("UserComment")
	log.PanicIf(err)

	ucIte := results[0]

	if ucIte.TagId() != 0x9286 {
		t.Fatalf("tag-ID not correct")
	}

	recoveredValueBytes, err := ucIte.GetRawBytes()
	log.PanicIf(err)

	expectedValueBytes := make([]byte, 0)

	expectedValueBytes = append(expectedValueBytes, []byte{'A', 'S', 'C', 'I', 'I', 0, 0, 0}...)
	expectedValueBytes = append(expectedValueBytes, []byte("TEST COMMENT")...)

	if bytes.Compare(recoveredValueBytes, expectedValueBytes) != 0 {
		t.Fatalf("Recovered UserComment does not have the right value: %v != %v", recoveredValueBytes, expectedValueBytes)
	}
}

func TestSegment_SetExif_FromScratch(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
			t.Fatalf("Test failure.")
		}
	}()

	// Create the IB.

	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	err := exif.LoadStandardTags(ti)
	log.PanicIf(err)

	rootIb := exif.NewIfdBuilder(im, ti, exifcommon.IfdStandardIfdIdentity, exifcommon.EncodeDefaultByteOrder)

	err = rootIb.AddStandardWithName("ProcessingSoftware", "some software")
	log.PanicIf(err)

	// Encode.

	s := makeEmptyExifSegment()

	err = s.SetExif(rootIb)
	log.PanicIf(err)

	// Decode.

	rootIfd, _, err := s.Exif()
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

func TestSegment_Exif(t *testing.T) {
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

	_, s, err := sl.FindExif()
	log.PanicIf(err)

	rootIfd, data, err := s.Exif()
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
