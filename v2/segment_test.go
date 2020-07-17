package jpegstructure

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"io/ioutil"

	"github.com/dsoprea/go-exif/v3"
	"github.com/dsoprea/go-exif/v3/common"
	"github.com/dsoprea/go-exif/v3/undefined"
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

	im, err := exifcommon.NewIfdMappingWithStandard()
	log.PanicIf(err)

	ti := exif.NewTagIndex()

	err = exif.LoadStandardTags(ti)
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

func TestSegment_IsExif_Hit(t *testing.T) {
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

	if s.IsExif() != true {
		t.Fatalf("Did not return true.")
	}
}

func TestSegment_IsExif_Miss(t *testing.T) {
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

	if sl.Segments()[4].IsExif() != false {
		t.Fatalf("Did not return false.")
	}
}

func TestSegment_IsXmp_Hit(t *testing.T) {
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

	_, s, err := sl.FindXmp()
	log.PanicIf(err)

	if s.IsXmp() != true {
		t.Fatalf("Did not return true.")
	}
}

func TestSegment_IsXmp_Miss(t *testing.T) {
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

	if sl.Segments()[4].IsXmp() != false {
		t.Fatalf("Did not return false.")
	}
}

func TestSegment_FormattedXmp(t *testing.T) {
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

	_, s, err := sl.FindXmp()
	log.PanicIf(err)

	actualData, err := s.FormattedXmp()
	log.PanicIf(err)

	// Filter out the Unicode BOM character since this would add unnecessary
	// complexity to the test.
	actualData = strings.ReplaceAll(actualData, "\ufeff", "")

	// Replace Windows-style newlines to Unix.
	actualData = strings.ReplaceAll(actualData, "\r\n", "\n")

	expectedData := `<?xpacket begin='' id='W5M0MpCehiHzreSzNTczkc9d'?>
    <x:xmpmeta xmlns:x="adobe:ns:meta/">
      <rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
        <rdf:Description rdf:about="" xmlns:xmp="http://ns.adobe.com/xap/1.0/">
          <xmp:Rating>0
          </xmp:Rating>
        </rdf:Description>
      </rdf:RDF>
    </x:xmpmeta>
    <?xpacket end='w'?>`

	if actualData != expectedData {
		t.Fatalf("XMP data is not correct:\nACTUAL:\n>>>%s<<<\n\nEXPECTED:\n>>>%s<<<\n", actualData, expectedData)
	}
}
