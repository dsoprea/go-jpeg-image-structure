package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"testing"

	"encoding/json"
	"os/exec"

	"github.com/dsoprea/go-logging"
)

var (
	assetsPath  = ""
	appFilepath = ""
)

type JsonResultExifTag struct {
	MarkerId   byte   `json:"marker_id"`
	MarkerName string `json:"market_name"`
	Offset     int    `json:"offset"`
	Data       []byte `json:"data"`
}

func TestMain_Plain_Exif(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "NDM_8901.jpg")

	cmd := exec.Command(
		"go", "run", appFilepath,
		"--filepath", imageFilepath)

	b := new(bytes.Buffer)
	cmd.Stdout = b
	cmd.Stderr = b

	err := cmd.Run()
	actual := b.String()

	if err != nil {
		fmt.Printf(actual)
		panic(err)
	}

	expected :=
		` 0: IFD-PATH=[IFD] ID=(0x10f) NAME=[Make] TYPE=(2):[ASCII] VALUE=[Canon]
 1: IFD-PATH=[IFD] ID=(0x110) NAME=[Model] TYPE=(2):[ASCII] VALUE=[Canon EOS 5D Mark III]
 2: IFD-PATH=[IFD] ID=(0x112) NAME=[Orientation] TYPE=(3):[SHORT] VALUE=[[1]]
 3: IFD-PATH=[IFD] ID=(0x11a) NAME=[XResolution] TYPE=(5):[RATIONAL] VALUE=[[{72 1}]]
 4: IFD-PATH=[IFD] ID=(0x11b) NAME=[YResolution] TYPE=(5):[RATIONAL] VALUE=[[{72 1}]]
 5: IFD-PATH=[IFD] ID=(0x128) NAME=[ResolutionUnit] TYPE=(3):[SHORT] VALUE=[[2]]
 6: IFD-PATH=[IFD] ID=(0x132) NAME=[DateTime] TYPE=(2):[ASCII] VALUE=[2017:12:02 08:18:50]
 7: IFD-PATH=[IFD] ID=(0x13b) NAME=[Artist] TYPE=(2):[ASCII] VALUE=[]
 8: IFD-PATH=[IFD] ID=(0x213) NAME=[YCbCrPositioning] TYPE=(3):[SHORT] VALUE=[[2]]
 9: IFD-PATH=[IFD] ID=(0x8298) NAME=[Copyright] TYPE=(2):[ASCII] VALUE=[]
10: IFD-PATH=[IFD] ID=(0x8769) NAME=[ExifTag] TYPE=(4):[LONG] VALUE=[[360]] CHILD-IFD-PATH=[IFD/Exif]
11: IFD-PATH=[IFD] ID=(0x8825) NAME=[GPSTag] TYPE=(4):[LONG] VALUE=[[9554]] CHILD-IFD-PATH=[IFD/GPSInfo]
12: IFD-PATH=[IFD/Exif] ID=(0x829a) NAME=[ExposureTime] TYPE=(5):[RATIONAL] VALUE=[[{1 640}]]
13: IFD-PATH=[IFD/Exif] ID=(0x829d) NAME=[FNumber] TYPE=(5):[RATIONAL] VALUE=[[{4 1}]]
14: IFD-PATH=[IFD/Exif] ID=(0x8822) NAME=[ExposureProgram] TYPE=(3):[SHORT] VALUE=[[4]]
15: IFD-PATH=[IFD/Exif] ID=(0x8827) NAME=[ISOSpeedRatings] TYPE=(3):[SHORT] VALUE=[[1600]]
16: IFD-PATH=[IFD/Exif] ID=(0x8830) NAME=[SensitivityType] TYPE=(3):[SHORT] VALUE=[[2]]
17: IFD-PATH=[IFD/Exif] ID=(0x8832) NAME=[RecommendedExposureIndex] TYPE=(4):[LONG] VALUE=[[1600]]
18: IFD-PATH=[IFD/Exif] ID=(0x9000) NAME=[ExifVersion] TYPE=(7):[UNDEFINED] VALUE=[0230]
19: IFD-PATH=[IFD/Exif] ID=(0x9003) NAME=[DateTimeOriginal] TYPE=(2):[ASCII] VALUE=[2017:12:02 08:18:50]
20: IFD-PATH=[IFD/Exif] ID=(0x9004) NAME=[DateTimeDigitized] TYPE=(2):[ASCII] VALUE=[2017:12:02 08:18:50]
21: IFD-PATH=[IFD/Exif] ID=(0x9101) NAME=[ComponentsConfiguration] TYPE=(7):[UNDEFINED] VALUE=[ComponentsConfiguration<ID=[YCBCR] BYTES=[1 2 3 0]>]
22: IFD-PATH=[IFD/Exif] ID=(0x9201) NAME=[ShutterSpeedValue] TYPE=(10):[SRATIONAL] VALUE=[[{614400 65536}]]
23: IFD-PATH=[IFD/Exif] ID=(0x9202) NAME=[ApertureValue] TYPE=(5):[RATIONAL] VALUE=[[{262144 65536}]]
24: IFD-PATH=[IFD/Exif] ID=(0x9204) NAME=[ExposureBiasValue] TYPE=(10):[SRATIONAL] VALUE=[[{0 1}]]
25: IFD-PATH=[IFD/Exif] ID=(0x9207) NAME=[MeteringMode] TYPE=(3):[SHORT] VALUE=[[5]]
26: IFD-PATH=[IFD/Exif] ID=(0x9209) NAME=[Flash] TYPE=(3):[SHORT] VALUE=[[16]]
27: IFD-PATH=[IFD/Exif] ID=(0x920a) NAME=[FocalLength] TYPE=(5):[RATIONAL] VALUE=[[{16 1}]]
28: IFD-PATH=[IFD/Exif] ID=(0x927c) NAME=[MakerNote] TYPE=(7):[UNDEFINED] VALUE=[MakerNote<TYPE-ID=[28 00 01 00 03 00 31 00 00 00 74 05 00 00 02 00 03 00 04 00]>]
29: IFD-PATH=[IFD/Exif] ID=(0x9286) NAME=[UserComment] TYPE=(7):[UNDEFINED] VALUE=[UserComment<SIZE=(256) ENCODING=[UNDEFINED] V=[0 0 0 0 0 0 0 0]... LEN=(256)>]
30: IFD-PATH=[IFD/Exif] ID=(0x9290) NAME=[SubSecTime] TYPE=(2):[ASCII] VALUE=[00]
31: IFD-PATH=[IFD/Exif] ID=(0x9291) NAME=[SubSecTimeOriginal] TYPE=(2):[ASCII] VALUE=[00]
32: IFD-PATH=[IFD/Exif] ID=(0x9292) NAME=[SubSecTimeDigitized] TYPE=(2):[ASCII] VALUE=[00]
33: IFD-PATH=[IFD/Exif] ID=(0xa000) NAME=[FlashpixVersion] TYPE=(7):[UNDEFINED] VALUE=[0100]
34: IFD-PATH=[IFD/Exif] ID=(0xa001) NAME=[ColorSpace] TYPE=(3):[SHORT] VALUE=[[1]]
35: IFD-PATH=[IFD/Exif] ID=(0xa002) NAME=[PixelXDimension] TYPE=(3):[SHORT] VALUE=[[3840]]
36: IFD-PATH=[IFD/Exif] ID=(0xa003) NAME=[PixelYDimension] TYPE=(3):[SHORT] VALUE=[[2560]]
37: IFD-PATH=[IFD/Exif] ID=(0xa005) NAME=[InteroperabilityTag] TYPE=(4):[LONG] VALUE=[[9326]] CHILD-IFD-PATH=[IFD/Exif/Iop]
38: IFD-PATH=[IFD/Exif] ID=(0xa20e) NAME=[FocalPlaneXResolution] TYPE=(5):[RATIONAL] VALUE=[[{3840000 1461}]]
39: IFD-PATH=[IFD/Exif] ID=(0xa20f) NAME=[FocalPlaneYResolution] TYPE=(5):[RATIONAL] VALUE=[[{2560000 972}]]
40: IFD-PATH=[IFD/Exif] ID=(0xa210) NAME=[FocalPlaneResolutionUnit] TYPE=(3):[SHORT] VALUE=[[2]]
41: IFD-PATH=[IFD/Exif] ID=(0xa401) NAME=[CustomRendered] TYPE=(3):[SHORT] VALUE=[[0]]
42: IFD-PATH=[IFD/Exif] ID=(0xa402) NAME=[ExposureMode] TYPE=(3):[SHORT] VALUE=[[0]]
43: IFD-PATH=[IFD/Exif] ID=(0xa403) NAME=[WhiteBalance] TYPE=(3):[SHORT] VALUE=[[0]]
44: IFD-PATH=[IFD/Exif] ID=(0xa406) NAME=[SceneCaptureType] TYPE=(3):[SHORT] VALUE=[[0]]
45: IFD-PATH=[IFD/Exif] ID=(0xa430) NAME=[CameraOwnerName] TYPE=(2):[ASCII] VALUE=[]
46: IFD-PATH=[IFD/Exif] ID=(0xa431) NAME=[BodySerialNumber] TYPE=(2):[ASCII] VALUE=[063024020097]
47: IFD-PATH=[IFD/Exif] ID=(0xa432) NAME=[LensSpecification] TYPE=(5):[RATIONAL] VALUE=[[{16 1} {35 1} {0 1} {0 1}]]
48: IFD-PATH=[IFD/Exif] ID=(0xa434) NAME=[LensModel] TYPE=(2):[ASCII] VALUE=[EF16-35mm f/4L IS USM]
49: IFD-PATH=[IFD/Exif] ID=(0xa435) NAME=[LensSerialNumber] TYPE=(2):[ASCII] VALUE=[2400001068]
50: IFD-PATH=[IFD/GPSInfo] ID=(0x00) NAME=[GPSVersionID] TYPE=(1):[BYTE] VALUE=[[2 3 0 0]]
51: IFD-PATH=[IFD] ID=(0x103) NAME=[Compression] TYPE=(3):[SHORT] VALUE=[[6]]
52: IFD-PATH=[IFD] ID=(0x11a) NAME=[XResolution] TYPE=(5):[RATIONAL] VALUE=[[{72 1}]]
53: IFD-PATH=[IFD] ID=(0x11b) NAME=[YResolution] TYPE=(5):[RATIONAL] VALUE=[[{72 1}]]
54: IFD-PATH=[IFD] ID=(0x128) NAME=[ResolutionUnit] TYPE=(3):[SHORT] VALUE=[[2]]
55: IFD-PATH=[IFD/Exif/Iop] ID=(0x01) NAME=[InteroperabilityIndex] TYPE=(2):[ASCII] VALUE=[R98]
56: IFD-PATH=[IFD/Exif/Iop] ID=(0x02) NAME=[InteroperabilityVersion] TYPE=(7):[UNDEFINED] VALUE=[0100]
`

	if actual != expected {
		fmt.Printf("ACTUAL:\n%s\n", actual)
		fmt.Printf("EXPECTED:\n%s\n", expected)

		t.Fatalf("Output not expected.")
	}
}

func TestMain_Json_Exif(t *testing.T) {
	imageFilepath := path.Join(assetsPath, "NDM_8901.jpg")

	cmd := exec.Command(
		"go", "run", appFilepath,
		"--json",
		"--filepath", imageFilepath)

	b := new(bytes.Buffer)
	cmd.Stdout = b
	cmd.Stderr = b

	err := cmd.Run()
	raw := b.Bytes()

	if err != nil {
		fmt.Printf(string(raw))
		panic(err)
	}

	result := make([]JsonResultExifTag, 0)

	err = json.Unmarshal(raw, &result)
	log.PanicIf(err)

	// TODO(dustin): !! Store the expected JSON in a file.

	if len(result) != 57 {
		t.Fatalf("Exif tag-count not correct: (%d)", len(result))
	}
}

func init() {
	goPath := os.Getenv("GOPATH")

	assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-jpeg-image-structure", "assets")
	appFilepath = path.Join(goPath, "src", "github.com", "dsoprea", "go-jpeg-image-structure", "command", "js_exif", "main.go")
}
