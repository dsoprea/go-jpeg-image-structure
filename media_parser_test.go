package jpegstructure

import (
    "os"
    "path"
    "testing"

    "io/ioutil"

    "github.com/dsoprea/go-logging"
)

func TestJpegMediaParser_ParseSegments(t *testing.T) {
    filepath := path.Join(assetsPath, testImageRelFilepath)
    f, err := os.Open(filepath)
    log.PanicIf(err)

    defer f.Close()

    stat, err := f.Stat()
    log.PanicIf(err)

    size := stat.Size()

    jmp := NewJpegMediaParser()

    sl, err := jmp.Parse(f, int(size))
    log.PanicIf(err)

    expected := []*Segment{
        &Segment{
            MarkerId: 0xd8,
            Offset:   0x0,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset:   0x2,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset:   0x000080b4,
        },
        &Segment{
            MarkerId: 0xdb,
            Offset:   0x8ab6,
        },
        &Segment{
            MarkerId: 0xc0,
            Offset:   0x8b3c,
        },
        &Segment{
            MarkerId: 0xc4,
            Offset:   0x8b4f,
        },
        &Segment{
            MarkerId: 0xda,
            Offset:   0x8cf3,
        },
        &Segment{
            MarkerId: 0x0,
            Offset:   0x8cf5,
        },
        &Segment{
            MarkerId: 0xd9,
            Offset:   0x554d6d,
        },
    }

    if len(sl.segments) != len(expected) {
        t.Fatalf("Number of segments is unexpected: (%d) != (%d)", len(sl.segments), len(expected))
    }

    for i, s := range sl.segments {
        if s.MarkerId != expected[i].MarkerId {
            t.Fatalf("Segment (%d) marker-ID not correct: (0x%02x != 0x%02x)", i, s.MarkerId, expected[i].MarkerId)
        } else if s.Offset != expected[i].Offset {
            t.Fatalf("Segment (%d) offset not correct: (0x%08x != 0x%08x)", i, s.Offset, expected[i].Offset)
        }
    }
}

func TestJpegMediaParser_ParseBytesStructure(t *testing.T) {
    filepath := path.Join(assetsPath, testImageRelFilepath)

    data, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    jmp := NewJpegMediaParser()

    sl, err := jmp.ParseBytes(data)
    log.PanicIf(err)

    expectedSegments := []*Segment{
        &Segment{
            MarkerId: 0xd8,
            Offset:   0x0,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset:   0x2,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset:   0x000080b4,
        },
        &Segment{
            MarkerId: 0xdb,
            Offset:   0x8ab6,
        },
        &Segment{
            MarkerId: 0xc0,
            Offset:   0x8b3c,
        },
        &Segment{
            MarkerId: 0xc4,
            Offset:   0x8b4f,
        },
        &Segment{
            MarkerId: 0xda,
            Offset:   0x8cf3,
        },
        &Segment{
            MarkerId: 0x0,
            Offset:   0x8cf5,
        },
        &Segment{
            MarkerId: 0xd9,
            Offset:   0x554d6d,
        },
    }

    expectedSl := NewSegmentList(expectedSegments)

    if sl.OffsetsEqual(expectedSl) != true {
        t.Fatalf("Segments not expected")
    }
}

func TestJpegMediaParser_ParseBytesStructure_Offsets(t *testing.T) {
    filepath := path.Join(assetsPath, testImageRelFilepath)

    data, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    jmp := NewJpegMediaParser()

    sl, err := jmp.ParseBytes(data)
    log.PanicIf(err)

    err = sl.Validate(data)
    log.PanicIf(err)
}

func TestJpegMediaParser_ParseBytesStructure_MultipleEois(t *testing.T) {
    defer func() {
        if state := recover(); state != nil {
            err := log.Wrap(state.(error))
            log.PrintErrorf(err, "Test failure.")
            t.Fatalf("Test failure.")
        }
    }()

    filepath := path.Join(assetsPath, "IMG_6691_Multiple_EOIs.jpg")

    data, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    jmp := NewJpegMediaParser()

    sl, err := jmp.ParseBytes(data)
    log.PanicIf(err)

    expectedSegments := []*Segment{
        &Segment{
            MarkerId: 0xd8,
            Offset:   0x0,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset:   0x00000002,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset:   0x00007002,
        },
        &Segment{
            MarkerId: 0xe2,
            Offset:   0x00007fa0,
        },
        &Segment{
            MarkerId: 0xdb,
            Offset:   0x00008002,
        },
        &Segment{
            MarkerId: 0xc0,
            Offset:   0x00008088,
        },
        &Segment{
            MarkerId: 0xc4,
            Offset:   0x0000809b,
        },
        &Segment{
            MarkerId: 0xda,
            Offset:   0x0000823f,
        },
        &Segment{
            MarkerId: 0x0,
            Offset:   0x00008241,
        },
        &Segment{
            MarkerId: 0xd9,
            Offset:   0x00487540,
        },
    }

    expectedSl := NewSegmentList(expectedSegments)

    if sl.OffsetsEqual(expectedSl) != true {
        t.Fatalf("Segments not expected")
    }
}

func TestJpegMediaParser_LooksLikeFormat(t *testing.T) {
    filepath := path.Join(assetsPath, "NDM_8901.jpg")

    data, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    jmp := NewJpegMediaParser()

    if jmp.LooksLikeFormat(data) != true {
        t.Fatalf("not detected as JPEG")
    }
}
