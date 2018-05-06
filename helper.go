package jpegstructure

import (
    "os"
    "io"
    "bufio"
    "bytes"

    "github.com/dsoprea/go-logging"
)

func ParseSegments(r io.Reader, size int) (sl SegmentList, err error) {
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

func ParseFileStructure(filepath string) (sl SegmentList, err error) {
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

    sl, err = ParseSegments(f, int(size))
    log.PanicIf(err)

    return sl, nil
}

func ParseBytesStructure(data []byte) (sl SegmentList, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    b := bytes.NewBuffer(data)

    sl, err = ParseSegments(b, len(data))
    log.PanicIf(err)

    return sl, nil
}
