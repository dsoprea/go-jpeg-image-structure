package main

import (
	"bytes"
	"fmt"
	"path"
	"testing"

	"encoding/json"
	"os/exec"

	log "github.com/dsoprea/go-logging"

	jpegstructure "github.com/dsoprea/go-jpeg-image-structure"
)

type JsonResultJpegSegmentListItem struct {
	MarkerId   byte   `json:"marker_id"`
	MarkerName string `json:"market_name"`
	Offset     int    `json:"offset"`
	Data       []byte `json:"data"`
}

type JsonResultJpegSegmentIndexItem struct {
	MarkerName string `json:"marker_name"`
	Offset     int    `json:"offset"`
	Data       []byte `json:"data"`
}

func TestMain_Plain(t *testing.T) {
	imageFilepath := jpegstructure.GetTestImageFilepath()
	appFilepath := getAppFilepath()

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
		`JPEG Segments:

 0: OFFSET=(0x00000000          0) ID=(0xd8) NAME=[SOI  ] SIZE=(         0) SHA1=[da39a3ee5e6b4b0d3255bfef95601890afd80709]
 1: OFFSET=(0x00000002          2) ID=(0xe1) NAME=[APP1 ] SIZE=(     32942) SHA1=[81dce16a2abe2232049b5aa430ccf4095d240071] [EXIF]
 2: OFFSET=(0x000080b4      32948) ID=(0xe1) NAME=[APP1 ] SIZE=(         0) SHA1=[da39a3ee5e6b4b0d3255bfef95601890afd80709]
 3: OFFSET=(0x000080b8      32952) ID=(0xe1) NAME=[APP1 ] SIZE=(      2558) SHA1=[b56f13aa6bc3410a7d866302ef51c8b9798113af] [XMP]
 4: OFFSET=(0x00008aba      35514) ID=(0xdb) NAME=[DQT  ] SIZE=(       130) SHA1=[40441c843ce4c8027cbd3dbdc174ac13d7555aec]
 5: OFFSET=(0x00008b40      35648) ID=(0xc0) NAME=[SOF0 ] SIZE=(        15) SHA1=[2458a7e3cf26aed68a0becb123a0a022c03d1243]
 6: OFFSET=(0x00008b53      35667) ID=(0xc4) NAME=[DHT  ] SIZE=(       416) SHA1=[41b700bdd457862ce170bec95c9dac272e415470]
 7: OFFSET=(0x00008cf7      36087) ID=(0xda) NAME=[SOS  ] SIZE=(         0) SHA1=[da39a3ee5e6b4b0d3255bfef95601890afd80709]
 8: OFFSET=(0x00008cf9      36089) ID=(0x00) NAME=[     ] SIZE=(   5554296) SHA1=[16e7465a831a075b096dbd7f2d6f2c931e509edd]
 9: OFFSET=(0x00554d71    5590385) ID=(0xd9) NAME=[EOI  ] SIZE=(         0) SHA1=[da39a3ee5e6b4b0d3255bfef95601890afd80709]
`

	if actual != expected {
		fmt.Printf("ACTUAL:\n%s\n", actual)
		fmt.Printf("EXPECTED:\n%s\n", expected)

		t.Fatalf("Output not expected.")
	}
}

func TestMain_Json_NoData(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
		}
	}()

	imageFilepath := jpegstructure.GetTestImageFilepath()
	appFilepath := getAppFilepath()

	cmd := exec.Command(
		"go", "run", appFilepath,
		"--json-list",
		"--filepath", imageFilepath)

	b := new(bytes.Buffer)
	cmd.Stdout = b
	cmd.Stderr = b

	err := cmd.Run()
	actual := b.String()

	if err != nil {
		fmt.Println(actual)
		panic(err)
	}

	expected := `[
  {
    "marker_id": 216,
    "marker_name": "SOI",
    "offset": 0,
    "data": null,
    "length": 0
  },
  {
    "marker_id": 225,
    "marker_name": "APP1",
    "offset": 2,
    "data": null,
    "length": 32942
  },
  {
    "marker_id": 225,
    "marker_name": "APP1",
    "offset": 32948,
    "data": null,
    "length": 0
  },
  {
    "marker_id": 225,
    "marker_name": "APP1",
    "offset": 32952,
    "data": null,
    "length": 2558
  },
  {
    "marker_id": 219,
    "marker_name": "DQT",
    "offset": 35514,
    "data": null,
    "length": 130
  },
  {
    "marker_id": 192,
    "marker_name": "SOF0",
    "offset": 35648,
    "data": null,
    "length": 15
  },
  {
    "marker_id": 196,
    "marker_name": "DHT",
    "offset": 35667,
    "data": null,
    "length": 416
  },
  {
    "marker_id": 218,
    "marker_name": "SOS",
    "offset": 36087,
    "data": null,
    "length": 0
  },
  {
    "marker_id": 0,
    "marker_name": "!SCANDATA",
    "offset": 36089,
    "data": null,
    "length": 5554296
  },
  {
    "marker_id": 217,
    "marker_name": "EOI",
    "offset": 5590385,
    "data": null,
    "length": 0
  }
]
`

	if actual != expected {
		fmt.Printf("ACTUAL:\n%s\n\nEXPECTED:\n%s\n", actual, expected)

		t.Fatalf("output not expected.")
	}
}

func TestMain_Json_NoData_SegmentIndex(t *testing.T) {
	imageFilepath := jpegstructure.GetTestImageFilepath()
	appFilepath := getAppFilepath()

	cmd := exec.Command(
		"go", "run", appFilepath,
		"--json-object",
		"--filepath", imageFilepath)

	b := new(bytes.Buffer)
	cmd.Stdout = b
	cmd.Stderr = b

	err := cmd.Run()
	actual := b.String()

	if err != nil {
		fmt.Println(actual)
		panic(err)
	}

	expected := `{
  "!SCANDATA": [
    {
      "offset": 36089,
      "data": null,
      "length": 5554296
    }
  ],
  "APP1": [
    {
      "offset": 2,
      "data": null,
      "length": 32942
    },
    {
      "offset": 32948,
      "data": null,
      "length": 0
    },
    {
      "offset": 32952,
      "data": null,
      "length": 2558
    }
  ],
  "DHT": [
    {
      "offset": 35667,
      "data": null,
      "length": 416
    }
  ],
  "DQT": [
    {
      "offset": 35514,
      "data": null,
      "length": 130
    }
  ],
  "EOI": [
    {
      "offset": 5590385,
      "data": null,
      "length": 0
    }
  ],
  "SOF0": [
    {
      "offset": 35648,
      "data": null,
      "length": 15
    }
  ],
  "SOI": [
    {
      "offset": 0,
      "data": null,
      "length": 0
    }
  ],
  "SOS": [
    {
      "offset": 36087,
      "data": null,
      "length": 0
    }
  ]
}
`

	if actual != expected {
		fmt.Printf("ACTUAL:\n%s\n\nEXPECTED:\n%s\n", actual, expected)

		t.Fatalf("output not expected.")
	}
}

func TestMain_Json_Data(t *testing.T) {
	imageFilepath := jpegstructure.GetTestImageFilepath()
	appFilepath := getAppFilepath()

	cmd := exec.Command(
		"go", "run", appFilepath,
		"--json-list",
		"--data",
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

	result := make([]JsonResultJpegSegmentListItem, 0)

	err = json.Unmarshal(raw, &result)
	log.PanicIf(err)

	if len(result) != 10 {
		t.Fatalf("JPEG segment count not correct: (%d)", len(result))
	}

	hasData := false
	for _, s := range result {
		if s.Data != nil {
			hasData = true
			break
		}
	}

	if hasData != true {
		t.Fatalf("No segments have data but were expected to.")
	}
}

func getAppFilepath() string {
	moduleRootPath := jpegstructure.GetModuleRootPath()
	return path.Join(moduleRootPath, "command", "js_dump", "main.go")
}
