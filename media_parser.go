package jpegstructure

import (
    "os"
    "io"
    "bufio"
    "bytes"

    "github.com/dsoprea/go-logging"
)


type JpegMediaParser struct {
}

func NewJpegMediaParser() *JpegMediaParser {
    return new(JpegMediaParser)
}

func (jmp *JpegMediaParser) Parse(r io.Reader, size int) (sl *SegmentList, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    s := bufio.NewScanner(r)

    // Since each segment can be any size, our buffer must allowed to grow as
    // large as the file.
    buffer := []byte {}
    s.Buffer(buffer, size)

    js := NewJpegSplitter(nil)
    s.Split(js.Split)

    for ; s.Scan() != false; { }
    log.PanicIf(s.Err())

    return js.Segments(), nil
}

func (jmp *JpegMediaParser) ParseFile(filepath string) (sl *SegmentList, err error) {
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

    sl, err = jmp.Parse(f, int(size))
    log.PanicIf(err)

    return sl, nil
}

func (jmp *JpegMediaParser) ParseBytes(data []byte) (sl *SegmentList, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    b := bytes.NewBuffer(data)

    sl, err = jmp.Parse(b, len(data))
    log.PanicIf(err)

    return sl, nil
}

func (jmp *JpegMediaParser) LooksLikeFormat(data []byte) bool {
    if len(data) < 4 {
        return false
    }

    len_ := len(data)
    if data[0] != 0xff || data[1] != MARKER_SOI || data[len_ - 2] != 0xff || data[len_ - 1] != MARKER_EOI {
        return false
    }

    return true
}
