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

 0: ID=(0xd8) OFFSET=(0x00000000 0)
 1: ID=(0xe1) OFFSET=(0x00000002 2)
 2: ID=(0xe1) OFFSET=(0x000080b4 32948)
 3: ID=(0xdb) OFFSET=(0x00008ab6 35510)
 4: ID=(0xc0) OFFSET=(0x00008b3c 35644)
 5: ID=(0xc4) OFFSET=(0x00008b4f 35663)
 6: ID=(0xda) OFFSET=(0x00008cf3 36083)
 7: ID=(0x00) OFFSET=(0x00008cf5 36085)
 8: ID=(0xd9) OFFSET=(0x00554d6d 5590381)
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
        fmt.Printf(string(raw))
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
