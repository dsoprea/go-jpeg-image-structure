package jpegstructure

import (
	"fmt"
	"os"
	"path"
	"testing"

	"io/ioutil"

	log "github.com/dsoprea/go-logging"
)

func TestJpegMediaParser_Parse(t *testing.T) {
	filepath := GetTestImageFilepath()

	f, err := os.Open(filepath)
	log.PanicIf(err)

	defer f.Close()

	stat, err := f.Stat()
	log.PanicIf(err)

	size := stat.Size()

	jmp := NewJpegMediaParser()

	intfc, err := jmp.Parse(f, int(size))
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	expected := []*Segment{
		{
			MarkerId: 0xd8,
			Offset:   0x0,
		},
		{
			MarkerId: 0xe1,
			Offset:   0x2,
		},
		{
			MarkerId: 0xe1,
			Offset:   0x000080b4,
		},
		{
			MarkerId: 0xe1,
			Offset:   0x80b8,
		},
		{
			MarkerId: 0xdb,
			Offset:   0x8aba,
		},
		{
			MarkerId: 0xc0,
			Offset:   0x8b40,
		},
		{
			MarkerId: 0xc4,
			Offset:   0x8b53,
		},
		{
			MarkerId: 0xda,
			Offset:   0x8cf7,
		},
		{
			MarkerId: 0x0,
			Offset:   0x8cf9,
		},
		{
			MarkerId: 0xd9,
			Offset:   0x554d71,
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

func TestJpegMediaParser_ParseBytes(t *testing.T) {
	filepath := GetTestImageFilepath()

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	jmp := NewJpegMediaParser()

	intfc, err := jmp.ParseBytes(data)
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	expectedSegments := []*Segment{
		{
			MarkerId: 0xd8,
			Offset:   0x0,
		},
		{
			MarkerId: 0xe1,
			Offset:   0x2,
		},
		{
			MarkerId: 0xe1,
			Offset:   0x000080b4,
		},
		{
			MarkerId: 0xe1,
			Offset:   0x80b8,
		},
		{
			MarkerId: 0xdb,
			Offset:   0x8aba,
		},
		{
			MarkerId: 0xc0,
			Offset:   0x8b40,
		},
		{
			MarkerId: 0xc4,
			Offset:   0x8b53,
		},
		{
			MarkerId: 0xda,
			Offset:   0x8cf7,
		},
		{
			MarkerId: 0x0,
			Offset:   0x8cf9,
		},
		{
			MarkerId: 0xd9,
			Offset:   0x554d71,
		},
	}

	expectedSl := NewSegmentList(expectedSegments)

	if sl.OffsetsEqual(expectedSl) != true {
		t.Fatalf("Segments not expected")
	}
}

func TestJpegMediaParser_ParseBytes_Offsets(t *testing.T) {
	filepath := GetTestImageFilepath()

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	jmp := NewJpegMediaParser()

	intfc, err := jmp.ParseBytes(data)
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	err = sl.Validate(data)
	log.PanicIf(err)
}

func TestJpegMediaParser_ParseBytes_MultipleEois(t *testing.T) {
	defer func() {
		if state := recover(); state != nil {
			err := log.Wrap(state.(error))
			log.PrintErrorf(err, "Test failure.")
			t.Fatalf("Test failure.")
		}
	}()

	assetsPath := GetTestAssetsPath()
	filepath := path.Join(assetsPath, "IMG_6691_Multiple_EOIs.jpg")

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	jmp := NewJpegMediaParser()

	intfc, err := jmp.ParseBytes(data)
	log.PanicIf(err)

	sl := intfc.(*SegmentList)

	expectedSegments := []*Segment{
		{
			MarkerId: 0xd8,
			Offset:   0x0,
		},
		{
			MarkerId: 0xe1,
			Offset:   0x00000002,
		},
		{
			MarkerId: 0xe1,
			Offset:   0x00007002,
		},
		{
			MarkerId: 0xe2,
			Offset:   0x00007fa0,
		},
		{
			MarkerId: 0xdb,
			Offset:   0x00008002,
		},
		{
			MarkerId: 0xc0,
			Offset:   0x00008088,
		},
		{
			MarkerId: 0xc4,
			Offset:   0x0000809b,
		},
		{
			MarkerId: 0xda,
			Offset:   0x0000823f,
		},
		{
			MarkerId: 0x0,
			Offset:   0x00008241,
		},
		{
			MarkerId: 0xd9,
			Offset:   0x003f24db,
		},
	}

	expectedSl := NewSegmentList(expectedSegments)

	if sl.OffsetsEqual(expectedSl) != true {
		for i, segment := range sl.segments {
			fmt.Printf("%d: ACTUAL: MARKER=(%02x) OFF=(%10x)\n", i, segment.MarkerId, segment.Offset)
		}

		for i, segment := range expectedSl.segments {
			fmt.Printf("%d: EXPECTED: MARKER=(%02x) OFF=(%10x)\n", i, segment.MarkerId, segment.Offset)
		}

		t.Fatalf("Segments not expected")
	}
}

func TestJpegMediaParser_LooksLikeFormat(t *testing.T) {
	filepath := GetTestImageFilepath()

	data, err := ioutil.ReadFile(filepath)
	log.PanicIf(err)

	jmp := NewJpegMediaParser()

	if jmp.LooksLikeFormat(data) != true {
		t.Fatalf("not detected as JPEG")
	}
}
