package jpegstructure

import (
    "testing"
    "os"
    "path"

    "io/ioutil"

    "github.com/dsoprea/go-logging"
)

func TestParseSegments(t *testing.T) {
    filepath := path.Join(assetsPath, testImageRelFilepath)
    f, err := os.Open(filepath)
    log.PanicIf(err)

    defer f.Close()

    stat, err := f.Stat()
    log.PanicIf(err)

    size := stat.Size()

    sl, err := ParseSegments(f, int(size))
    log.PanicIf(err)

    expected := []*Segment {
        &Segment{
            MarkerId: 0xd8,
            Offset: 0x0,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset: 0x2,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset: 0x000080b4,
        },
        &Segment{
            MarkerId: 0xdb,
            Offset: 0x8ab6,
        },
        &Segment{
            MarkerId: 0xc0,
            Offset: 0x8b3c,
        },
        &Segment{
            MarkerId: 0xc4,
            Offset: 0x8b4f,
        },
        &Segment{
            MarkerId: 0xda,
            Offset: 0x8cf3,
        },
        &Segment{
            MarkerId: 0x0,
            Offset: 0x8cf5,
        },
        &Segment{
            MarkerId: 0xd9,
            Offset: 0x554d6d,
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

func TestParseFileStructure(t *testing.T) {
    filepath := path.Join(assetsPath, testImageRelFilepath)

    sl, err := ParseFileStructure(filepath)
    log.PanicIf(err)

    expected := []*Segment {
        &Segment{
            MarkerId: 0xd8,
            Offset: 0x0,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset: 0x2,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset: 0x000080b4,
        },
        &Segment{
            MarkerId: 0xdb,
            Offset: 0x8ab6,
        },
        &Segment{
            MarkerId: 0xc0,
            Offset: 0x8b3c,
        },
        &Segment{
            MarkerId: 0xc4,
            Offset: 0x8b4f,
        },
        &Segment{
            MarkerId: 0xda,
            Offset: 0x8cf3,
        },
        &Segment{
            MarkerId: 0x0,
            Offset: 0x8cf5,
        },
        &Segment{
            MarkerId: 0xd9,
            Offset: 0x554d6d,
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

func TestParseBytesStructure(t *testing.T) {
    filepath := path.Join(assetsPath, testImageRelFilepath)

    data, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    sl, err := ParseBytesStructure(data)
    log.PanicIf(err)

    expected := []*Segment {
        &Segment{
            MarkerId: 0xd8,
            Offset: 0x0,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset: 0x2,
        },
        &Segment{
            MarkerId: 0xe1,
            Offset: 0x000080b4,
        },
        &Segment{
            MarkerId: 0xdb,
            Offset: 0x8ab6,
        },
        &Segment{
            MarkerId: 0xc0,
            Offset: 0x8b3c,
        },
        &Segment{
            MarkerId: 0xc4,
            Offset: 0x8b4f,
        },
        &Segment{
            MarkerId: 0xda,
            Offset: 0x8cf3,
        },
        &Segment{
            MarkerId: 0x0,
            Offset: 0x8cf5,
        },
        &Segment{
            MarkerId: 0xd9,
            Offset: 0x554d6d,
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

func TestParseBytesStructure_Offsets(t *testing.T) {
    filepath := path.Join(assetsPath, testImageRelFilepath)

    data, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    sl, err := ParseBytesStructure(data)
    log.PanicIf(err)

    err = sl.Validate(data)
    log.PanicIf(err)
}

func TestParseBytesStructure_Offsets_Error(t *testing.T) {
    filepath := path.Join(assetsPath, testImageRelFilepath)

    data, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    segments := []*Segment {
        &Segment{
            MarkerId: 0x0,
            Offset: 0x0,
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
