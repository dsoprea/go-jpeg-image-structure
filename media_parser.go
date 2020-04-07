package jpegstructure

import (
    "bufio"
    "bytes"
    "io"
    "os"

    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-utility/image"
)

type JpegMediaParser struct {
}

func NewJpegMediaParser() *JpegMediaParser {
    return new(JpegMediaParser)
}

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

    log.PanicIf(s.Err())

    return js.Segments(), nil
}

func (jmp *JpegMediaParser) ParseFile(filepath string) (ec riimage.MediaContext, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    f, err := os.Open(filepath)
    log.PanicIf(err)

    stat, err := f.Stat()
    log.PanicIf(err)

    size := stat.Size()

    sl, err := jmp.Parse(f, int(size))
    log.PanicIf(err)

    return sl, nil
}

func (jmp *JpegMediaParser) ParseBytes(data []byte) (ec riimage.MediaContext, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    br := bytes.NewReader(data)

    sl, err := jmp.Parse(br, len(data))
    log.PanicIf(err)

    return sl, nil
}

func (jmp *JpegMediaParser) LooksLikeFormat(data []byte) bool {
    if len(data) < 4 {
        return false
    }

    len_ := len(data)
    if data[0] != 0xff || data[1] != MARKER_SOI || data[len_-2] != 0xff || data[len_-1] != MARKER_EOI {
        return false
    }

    return true
}

var (
    // Enforce that `JpegMediaParser` looks like a `riimage.MediaParser`.
    _ riimage.MediaParser = new(JpegMediaParser)
)
