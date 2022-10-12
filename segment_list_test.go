package jpegstructure

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"io/ioutil"

	"github.com/dsoprea/go-exif/v2"
	exifcommon "github.com/dsoprea/go-exif/v2/common"
	exifundefined "github.com/dsoprea/go-exif/v2/undefined"
	log "github.com/dsoprea/go-logging"
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
		{MarkerId: 0},
		{MarkerId: 0},
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

	// Update the DateTime tag.

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

	// Output:
}

func TestSegmentList_ConstructExifBuilder(t *testing.T) {
	filepath := GetTestImageFilepath()

	// Parse the image.

	jmp := NewJpegMediaParser()

	intfc, err := jmp.ParseFile(filepath)
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	_, err = sl.ConstructExifBuilder()
	log.PanicIf(err)

	// As long as ConstructExifBuilder() returns an `exif.IfdBuilder`, this test is satisfied.
}

func TestSegmentList_DumpExif(t *testing.T) {
	filepath := GetTestImageFilepath()

	// Parse the image.

	jmp := NewJpegMediaParser()

	intfc, err := jmp.ParseFile(filepath)
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	_, _, exifTags, err := sl.DumpExif()
	log.PanicIf(err)

	flattened := make([]string, len(exifTags))
	for i, et := range exifTags {
		flattened[i] = fmt.Sprintf("%s", et)
	}

	expected := []string{
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x10f) TAG-NAME=[Make] TAG-TYPE=[ASCII] VALUE=[Canon] VALUE-BYTES=(6) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x110) TAG-NAME=[Model] TAG-TYPE=[ASCII] VALUE=[Canon EOS 5D Mark III] VALUE-BYTES=(22) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x112) TAG-NAME=[Orientation] TAG-TYPE=[SHORT] VALUE=[1] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x11a) TAG-NAME=[XResolution] TAG-TYPE=[RATIONAL] VALUE=[72/1] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x11b) TAG-NAME=[YResolution] TAG-TYPE=[RATIONAL] VALUE=[72/1] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x128) TAG-NAME=[ResolutionUnit] TAG-TYPE=[SHORT] VALUE=[2] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x132) TAG-NAME=[DateTime] TAG-TYPE=[ASCII] VALUE=[2017:12:02 08:18:50] VALUE-BYTES=(20) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x13b) TAG-NAME=[Artist] TAG-TYPE=[ASCII] VALUE=[] VALUE-BYTES=(1) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x213) TAG-NAME=[YCbCrPositioning] TAG-TYPE=[SHORT] VALUE=[2] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x8298) TAG-NAME=[Copyright] TAG-TYPE=[ASCII] VALUE=[] VALUE-BYTES=(1) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x8769) TAG-NAME=[ExifTag] TAG-TYPE=[LONG] VALUE=[360] VALUE-BYTES=(4) CHILD-IFD-PATH=[IFD/Exif]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x829a) TAG-NAME=[ExposureTime] TAG-TYPE=[RATIONAL] VALUE=[1/640] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x829d) TAG-NAME=[FNumber] TAG-TYPE=[RATIONAL] VALUE=[4/1] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x8822) TAG-NAME=[ExposureProgram] TAG-TYPE=[SHORT] VALUE=[4] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x8827) TAG-NAME=[ISOSpeedRatings] TAG-TYPE=[SHORT] VALUE=[1600] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x8830) TAG-NAME=[SensitivityType] TAG-TYPE=[SHORT] VALUE=[2] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x8832) TAG-NAME=[RecommendedExposureIndex] TAG-TYPE=[LONG] VALUE=[1600] VALUE-BYTES=(4) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9000) TAG-NAME=[ExifVersion] TAG-TYPE=[UNDEFINED] VALUE=[0230] VALUE-BYTES=(4) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9003) TAG-NAME=[DateTimeOriginal] TAG-TYPE=[ASCII] VALUE=[2017:12:02 08:18:50] VALUE-BYTES=(20) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9004) TAG-NAME=[DateTimeDigitized] TAG-TYPE=[ASCII] VALUE=[2017:12:02 08:18:50] VALUE-BYTES=(20) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9101) TAG-NAME=[ComponentsConfiguration] TAG-TYPE=[UNDEFINED] VALUE=[Exif9101ComponentsConfiguration<ID=[YCBCR] BYTES=[1 2 3 0]>] VALUE-BYTES=(4) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9201) TAG-NAME=[ShutterSpeedValue] TAG-TYPE=[SRATIONAL] VALUE=[614400/65536] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9202) TAG-NAME=[ApertureValue] TAG-TYPE=[RATIONAL] VALUE=[262144/65536] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9204) TAG-NAME=[ExposureBiasValue] TAG-TYPE=[SRATIONAL] VALUE=[0/1] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9207) TAG-NAME=[MeteringMode] TAG-TYPE=[SHORT] VALUE=[5] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9209) TAG-NAME=[Flash] TAG-TYPE=[SHORT] VALUE=[16] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x920a) TAG-NAME=[FocalLength] TAG-TYPE=[RATIONAL] VALUE=[16/1] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x927c) TAG-NAME=[MakerNote] TAG-TYPE=[UNDEFINED] VALUE=[MakerNote<TYPE-ID=[28 00 01 00 03 00 31 00 00 00 74 05 00 00 02 00 03 00 04 00] LEN=(8152) SHA1=[d4154aa7df5474efe7ab38de2595919b9b4cc29f]>] VALUE-BYTES=(8152) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9286) TAG-NAME=[UserComment] TAG-TYPE=[UNDEFINED] VALUE=[UserComment<SIZE=(256) ENCODING=[UNDEFINED] V=[0 0 0 0 0 0 0 0]... LEN=(256)>] VALUE-BYTES=(264) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9290) TAG-NAME=[SubSecTime] TAG-TYPE=[ASCII] VALUE=[00] VALUE-BYTES=(3) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9291) TAG-NAME=[SubSecTimeOriginal] TAG-TYPE=[ASCII] VALUE=[00] VALUE-BYTES=(3) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0x9292) TAG-NAME=[SubSecTimeDigitized] TAG-TYPE=[ASCII] VALUE=[00] VALUE-BYTES=(3) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa000) TAG-NAME=[FlashpixVersion] TAG-TYPE=[UNDEFINED] VALUE=[0100] VALUE-BYTES=(4) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa001) TAG-NAME=[ColorSpace] TAG-TYPE=[SHORT] VALUE=[1] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa002) TAG-NAME=[PixelXDimension] TAG-TYPE=[SHORT] VALUE=[3840] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa003) TAG-NAME=[PixelYDimension] TAG-TYPE=[SHORT] VALUE=[2560] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa005) TAG-NAME=[InteroperabilityTag] TAG-TYPE=[LONG] VALUE=[9326] VALUE-BYTES=(4) CHILD-IFD-PATH=[IFD/Exif/Iop]",
		"ExifTag<IFD-PATH=[IFD/Exif/Iop] TAG-ID=(0x01) TAG-NAME=[InteroperabilityIndex] TAG-TYPE=[ASCII] VALUE=[R98] VALUE-BYTES=(4) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif/Iop] TAG-ID=(0x02) TAG-NAME=[InteroperabilityVersion] TAG-TYPE=[UNDEFINED] VALUE=[0100] VALUE-BYTES=(4) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa20e) TAG-NAME=[FocalPlaneXResolution] TAG-TYPE=[RATIONAL] VALUE=[3840000/1461] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa20f) TAG-NAME=[FocalPlaneYResolution] TAG-TYPE=[RATIONAL] VALUE=[2560000/972] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa210) TAG-NAME=[FocalPlaneResolutionUnit] TAG-TYPE=[SHORT] VALUE=[2] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa401) TAG-NAME=[CustomRendered] TAG-TYPE=[SHORT] VALUE=[0] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa402) TAG-NAME=[ExposureMode] TAG-TYPE=[SHORT] VALUE=[0] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa403) TAG-NAME=[WhiteBalance] TAG-TYPE=[SHORT] VALUE=[0] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa406) TAG-NAME=[SceneCaptureType] TAG-TYPE=[SHORT] VALUE=[0] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa430) TAG-NAME=[CameraOwnerName] TAG-TYPE=[ASCII] VALUE=[] VALUE-BYTES=(1) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa431) TAG-NAME=[BodySerialNumber] TAG-TYPE=[ASCII] VALUE=[063024020097] VALUE-BYTES=(13) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa432) TAG-NAME=[LensSpecification] TAG-TYPE=[RATIONAL] VALUE=[16/1...] VALUE-BYTES=(32) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa434) TAG-NAME=[LensModel] TAG-TYPE=[ASCII] VALUE=[EF16-35mm f/4L IS USM] VALUE-BYTES=(22) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD/Exif] TAG-ID=(0xa435) TAG-NAME=[LensSerialNumber] TAG-TYPE=[ASCII] VALUE=[2400001068] VALUE-BYTES=(11) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD] TAG-ID=(0x8825) TAG-NAME=[GPSTag] TAG-TYPE=[LONG] VALUE=[9554] VALUE-BYTES=(4) CHILD-IFD-PATH=[IFD/GPSInfo]",
		"ExifTag<IFD-PATH=[IFD/GPSInfo] TAG-ID=(0x00) TAG-NAME=[GPSVersionID] TAG-TYPE=[BYTE] VALUE=[02 03 00 00] VALUE-BYTES=(4) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD1] TAG-ID=(0x103) TAG-NAME=[Compression] TAG-TYPE=[SHORT] VALUE=[6] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD1] TAG-ID=(0x11a) TAG-NAME=[XResolution] TAG-TYPE=[RATIONAL] VALUE=[72/1] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD1] TAG-ID=(0x11b) TAG-NAME=[YResolution] TAG-TYPE=[RATIONAL] VALUE=[72/1] VALUE-BYTES=(8) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD1] TAG-ID=(0x128) TAG-NAME=[ResolutionUnit] TAG-TYPE=[SHORT] VALUE=[2] VALUE-BYTES=(2) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD1] TAG-ID=(0x201) TAG-NAME=[JPEGInterchangeFormat] TAG-TYPE=[LONG] VALUE=[11444] VALUE-BYTES=(4) CHILD-IFD-PATH=[]",
		"ExifTag<IFD-PATH=[IFD1] TAG-ID=(0x202) TAG-NAME=[JPEGInterchangeFormatLength] TAG-TYPE=[LONG] VALUE=[21491] VALUE-BYTES=(4) CHILD-IFD-PATH=[]",
	}

	if reflect.DeepEqual(flattened, expected) != true {
		for _, line := range flattened {
			fmt.Printf("ACTUAL: \"%s\",\n", line)
		}

		for _, line := range expected {
			fmt.Printf("EXPECTED: \"%s\",\n", line)
		}

		t.Fatalf("Tags are not correct.")
	}
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
		t.Fatalf("EXIF data not correct")
	}
}

func TestSegmentList_FindXmp(t *testing.T) {
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

	segmentNumber, s, err := sl.FindXmp()
	log.PanicIf(err)

	if segmentNumber != 3 {
		t.Fatalf("XMP not found in right position: (%d)", segmentNumber)
	}

	actualData := string(s.Data[len(xmpPrefix):])

	// We can't display it in raw form since it frequently has oddball
	// formatting/whitespacing spread over many lines, and is usually
	// reformatted by the text editor.
	actualData, err = FormatXml(actualData)
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

func TestSegmentList_Validate(t *testing.T) {
	filepath := GetTestImageFilepath()

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	segments := []*Segment{
		{
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
