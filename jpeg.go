package exifjpeg

import (
	"bufio"
	"encoding/binary"
	"io"
	"os"

	"github.com/dsoprea/go-logging"
)

const (
	MARKER_EOI   = 0xD9
	MARKER_SOS   = 0xDA
	MARKER_SOD   = 0x93
	MARKER_DQT   = 0xDB
	MARKER_APP0  = 0xE0
	MARKER_APP1  = 0xE1
	MARKER_APP2  = 0xE2
	MARKER_APP3  = 0xE3
	MARKER_APP4  = 0xE4
	MARKER_APP5  = 0xE5
	MARKER_APP6  = 0xE6
	MARKER_APP7  = 0xE7
	MARKER_APP8  = 0xE8
	MARKER_APP10 = 0xEA
	MARKER_APP12 = 0xEC
	MARKER_APP13 = 0xED
	MARKER_APP14 = 0xEE
	MARKER_APP15 = 0xEF
	MARKER_COM   = 0xFE
	MARKER_CME   = 0x64
	MARKER_SIZ   = 0x51
)

var (
	jpegLogger        = log.NewLogger("exifjpeg.jpeg")
	jpegMagicStandard = []byte{0xff, 0xd8, 0xff}
	jpegMagic2000     = []byte{0xff, 0x4f, 0xff}

	markerLen = map[byte]int{
		0x00: 0,
		0x01: 0,
		0xd0: 0,
		0xd1: 0,
		0xd2: 0,
		0xd3: 0,
		0xd4: 0,
		0xd5: 0,
		0xd6: 0,
		0xd7: 0,
		0xd8: 0,
		0xd9: 0,
		0xda: 0,

		// J2C
		0x30: 0,
		0x31: 0,
		0x32: 0,
		0x33: 0,
		0x34: 0,
		0x35: 0,
		0x36: 0,
		0x37: 0,
		0x38: 0,
		0x39: 0,
		0x3a: 0,
		0x3b: 0,
		0x3c: 0,
		0x3d: 0,
		0x3e: 0,
		0x3f: 0,
		0x4f: 0,
		0x92: 0,
		0x93: 0,

		// J2C extensions
		0x74: 4,
		0x75: 4,
		0x77: 4,
	}
)

type JpegNavigator struct {
	f *os.File
	r *bufio.Reader
}

func NewJpegNavigator(filepath string) *JpegNavigator {
	defer func() {
		if state := recover(); state != nil {
			log.Panic(state.(error))
		}
	}()

	f, err := os.Open(filepath)
	log.PanicIf(err)

	buffer := make([]byte, 3)
	n, err := f.Read(buffer)
	log.PanicIf(err)

	if n != 3 {
		log.Panicf("file not long enough to identify as a JPEG")
	}

	if buffer[0] == jpegMagic2000[0] && buffer[1] == jpegMagic2000[1] && buffer[2] == jpegMagic2000[2] {
		// TODO(dustin): Return to this.
		log.Panicf("JPEG2000 not supported")
	}

	if buffer[0] != jpegMagicStandard[0] || buffer[1] != jpegMagicStandard[1] || buffer[2] != jpegMagicStandard[2] {
		log.Panicf("file does not look like a JPEG: (%X) (%X) (%X)", buffer[0], buffer[1], buffer[2])
	}

	r := bufio.NewReader(f)

	jn := &JpegNavigator{
		f: f,
		r: r,
	}

	// Seek to first segment.

	err = jn.SeekToNextSegment()
	log.PanicIf(err)

	return jn
}

func (jp *JpegNavigator) Close() {
	defer func() {
		if state := recover(); state != nil {
			log.Panic(state.(error))
		}
	}()

	jp.f.Close()
}

type SegmentVisitorCb func(markerId byte) (continue_ bool, err error)

func (jn *JpegNavigator) VisitSegments(cb SegmentVisitorCb) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	for {
		si, err := jn.ReadSegment()
		log.PanicIf(err)

		continue_, err := cb(si.MarkerId)
		log.PanicIf(err)

		if continue_ == false {
			break
		}

		err = jn.SeekToNextSegment()
		if err != nil {
			if log.Is(err, io.EOF) == true {
				break
			}

			log.Panic(err)
		}
	}

	return nil
}

func (jn *JpegNavigator) SeekToNextSegment() (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// Find boundary of next marker.
	_, err = jn.r.ReadBytes(0xff)
	log.PanicIf(err)

	// Seek past additional padding.

	for {
		b, err := jn.r.ReadByte()
		log.PanicIf(err)

		if b != 0xff {
			err = jn.r.UnreadByte()
			log.PanicIf(err)

			break
		}
	}

	return nil
}

type SegmentInfo struct {
	MarkerId   byte
	DataLength uint32
}

func (jn *JpegNavigator) ReadSegment() (si *SegmentInfo, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	markerId, err := jn.r.ReadByte()

	si = &SegmentInfo{
		MarkerId: markerId,
	}

	if len_, found := markerLen[markerId]; found == false {
		// It's not one of the static-length markers. Read the length.
		//
		// The length is an unsigned 16-bit network/big-endian.

		len_ := uint16(0)
		err = binary.Read(jn.r, binary.BigEndian, &len_)
		log.PanicIf(err)

		si.DataLength = uint32(len_)

		// Includes the bytes of the length itself.
		len_ -= 2
	} else if len_ > 0 {
		// Accomodates the non-zero markers in our marker index, which only
		// represent J2C extensions.
		//
		// The length is an unsigned 32-bit network/big-endian.

		buffer := make([]byte, len_)
		n, err := jn.r.Read(buffer)
		log.PanicIf(err)

		if n != len_ {
			log.Panicf("ran out of data for J2C segment")
		}

		// TODO(dustin): We're just assuming that no length is greater than 4 (which is the current case).
		si.DataLength = uint32(len_)
	}

	return si, nil
}
