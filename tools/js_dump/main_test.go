package main

import (
    "testing"
    "os"
    "path"
    "bytes"
    "fmt"

    "os/exec"
    "encoding/json"

    "github.com/dsoprea/go-logging"
)

var (
    assetsPath = ""
    appFilepath = ""
)

type JsonResultJpegSegmentListItem struct {
    MarkerId byte `json:"marker_id"`
    MarkerName string `json:"market_name"`
    Offset int `json:"offset"`
    Data []byte `json:"data"`
}

type JsonResultJpegSegmentIndexItem struct {
    MarkerName string `json:"marker_name"`
    Offset int `json:"offset"`
    Data []byte `json:"data"`
}

func TestMain_Plain(t *testing.T) {
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
`JPEG Segments:

 0: OFFSET=(0x00000000          0) ID=(0x000000d8) NAME=[SOI ] SIZE=(         0) SHA1=[da39a3ee5e6b4b0d3255bfef95601890afd80709]
 1: OFFSET=(0x00000002          2) ID=(0x000000e1) NAME=[APP1] SIZE=(     32942) SHA1=[81dce16a2abe2232049b5aa430ccf4095d240071]
 2: OFFSET=(0x000080b4      32948) ID=(0x000000e1) NAME=[APP1] SIZE=(      2558) SHA1=[b56f13aa6bc3410a7d866302ef51c8b9798113af]
 3: OFFSET=(0x00008ab6      35510) ID=(0x000000db) NAME=[DQT ] SIZE=(       130) SHA1=[40441c843ce4c8027cbd3dbdc174ac13d7555aec]
 4: OFFSET=(0x00008b3c      35644) ID=(0x000000c0) NAME=[SOF0] SIZE=(        15) SHA1=[2458a7e3cf26aed68a0becb123a0a022c03d1243]
 5: OFFSET=(0x00008b4f      35663) ID=(0x000000c4) NAME=[DHT ] SIZE=(       416) SHA1=[41b700bdd457862ce170bec95c9dac272e415470]
 6: OFFSET=(0x00008cf3      36083) ID=(0x000000da) NAME=[SOS ] SIZE=(         0) SHA1=[da39a3ee5e6b4b0d3255bfef95601890afd80709]
 7: OFFSET=(0x00008cf5      36085) ID=(0x00000000) NAME=[    ] SIZE=(   5554296) SHA1=[16e7465a831a075b096dbd7f2d6f2c931e509edd]
 8: OFFSET=(0x00554d6d    5590381) ID=(0x000000d9) NAME=[EOI ] SIZE=(         0) SHA1=[da39a3ee5e6b4b0d3255bfef95601890afd80709]
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

    imageFilepath := path.Join(assetsPath, "NDM_8901.jpg")

    cmd := exec.Command(
            "go", "run", appFilepath,
            "--json-list",
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

    if len(result) != 9 {
        t.Fatalf("JPEG segment count not correct: (%d)", len(result))
    }

    for _, s := range result {
        if s.Data != nil {
            t.Fatalf("No segments were supposed to have data but do.")
        }
    }
}

func TestMain_Json_NoData_SegmentIndex(t *testing.T) {
    imageFilepath := path.Join(assetsPath, "NDM_8901.jpg")

    cmd := exec.Command(
            "go", "run", appFilepath,
            "--json-object",
            "--filepath", imageFilepath)

    b := new(bytes.Buffer)
    cmd.Stdout = b
    cmd.Stderr = b

    err := cmd.Run()
    raw := b.Bytes()

    if err != nil {
        fmt.Printf("RAW:\n%s\n", string(raw))
        panic(err)
    }

    result := make(map[string][]JsonResultJpegSegmentIndexItem)

    err = json.Unmarshal(raw, &result)
    log.PanicIf(err)

    if result == nil || len(result) == 0 {
        t.Fatalf("Segment index not returned/populated.")
    }

// TODO(dustin): !! Test actual segments returned in lists and indexes.
}

func TestMain_Json_Data(t *testing.T) {
    imageFilepath := path.Join(assetsPath, "NDM_8901.jpg")

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

    if len(result) != 9 {
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

func init() {
    goPath := os.Getenv("GOPATH")

    assetsPath = path.Join(goPath, "src", "github.com", "dsoprea", "go-jpeg-image-structure", "assets")
    appFilepath = path.Join(goPath, "src", "github.com", "dsoprea", "go-jpeg-image-structure", "tools", "js_dump", "main.go")
}
