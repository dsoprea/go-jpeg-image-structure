package jpegstructure

import (
    "bufio"
    "bytes"
    "errors"
    "io"
    "os"

    "github.com/dsoprea/go-logging"
)

var (
    // ErrJpegParseStoppedEarlier is an error that's usually symptomatic of an
    // image created by an nonstandard or exotic JPEG implementation.
    ErrJpegParseStoppedEarlier = errors.New("processing finished before EOI encountered")
)

type JpegMediaParser struct {
}

func NewJpegMediaParser() *JpegMediaParser {
    return new(JpegMediaParser)
}

func (jmp *JpegMediaParser) Parse(r io.Reader, size int) (intfc interface{}, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    s := bufio.NewScanner(r)

    // Since each segment can be any size, our buffer must allowed to grow as
    // large as the file.
    buffer := []byte{}
    s.Buffer(buffer, size)

    js := NewJpegSplitter(nil)
    s.Split(js.Split)

    for s.Scan() != false {
    }

    log.PanicIf(s.Err())

    // From time to time we encounter images that are nonconformant (or
    // unexpected, at the very least) and disrupt our parser. This will allow
    // us to identify those scenarios.
    if js.MarkerId() != MARKER_EOI {
        return nil, ErrJpegParseStoppedEarlier
    }

    return js.Segments(), nil
}

func (jmp *JpegMediaParser) ParseFile(filepath string) (intfc interface{}, err error) {
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

func (jmp *JpegMediaParser) ParseBytes(data []byte) (intfc interface{}, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    b := bytes.NewBuffer(data)

    sl, err := jmp.Parse(b, len(data))
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
