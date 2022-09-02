package jpegstructure

import (
	"bufio"
	"bytes"
	"io"
	"os"

	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/image"
)

// JpegMediaParser is a `riimage.MediaParser` that knows how to parse JPEG
// images.
type JpegMediaParser struct {
}

// NewJpegMediaParser returns a new JpegMediaParser.
func NewJpegMediaParser() *JpegMediaParser {

	// TODO(dustin): Add test

	return new(JpegMediaParser)
}

// Parse parses a JPEG uses an `io.ReadSeeker`. Even if it fails, it will return
// the list of segments encountered prior to the failure.
func (jmp *JpegMediaParser) Parse(rs io.ReadSeeker, size int) (ec riimage.MediaContext, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	s := bufio.NewScanner(rs)

	// Since each segment can be any size, our buffer must allowed to grow as
	// large as the file.
	buffer := []byte{}
	s.Buffer(buffer, size)

	js := NewJpegSplitter(nil)
	s.Split(js.Split)

	for s.Scan() != false {
	}

	// Always return the segments that were parsed, at least until there was an
	// error.
	ec = js.Segments()

	log.PanicIf(s.Err())

	return ec, nil
}

// ParseFile parses a JPEG file. Even if it fails, it will return the list of
// segments encountered prior to the failure.
func (jmp *JpegMediaParser) ParseFile(filepath string) (ec riimage.MediaContext, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// TODO(dustin): Add test

	f, err := os.Open(filepath)
	log.PanicIf(err)

	defer f.Close()

	stat, err := f.Stat()
	log.PanicIf(err)

	size := stat.Size()

	sl, err := jmp.Parse(f, int(size))

	// Always return the segments that were parsed, at least until there was an
	// error.
	ec = sl

	log.PanicIf(err)

	return ec, nil
}

// ParseBytes parses a JPEG byte-slice. Even if it fails, it will return the
// list of segments encountered prior to the failure.
func (jmp *JpegMediaParser) ParseBytes(data []byte) (ec riimage.MediaContext, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	br := bytes.NewReader(data)

	sl, err := jmp.Parse(br, len(data))

	// Always return the segments that were parsed, at least until there was an
	// error.
	ec = sl

	log.PanicIf(err)

	return ec, nil
}

// LooksLikeFormat indicates whether the data looks like a JPEG image.
func (jmp *JpegMediaParser) LooksLikeFormat(data []byte) bool {
	if len(data) < 4 {
		return false
	}

	// https://cs.opensource.google/go/go/+/master:src/net/http/sniff.go;l=126;drc=8b364451e2e2f2f816ed877a4639d9342279f299
	return bytes.HasPrefix(data, []byte("\xFF\xD8\xFF"))
}

var (
	// Enforce interface conformance.
	_ riimage.MediaParser = new(JpegMediaParser)
)
