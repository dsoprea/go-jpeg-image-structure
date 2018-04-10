package exifjpeg

import (
    "fmt"

    "github.com/dsoprea/go-logging"
)

const (
    BigEndianByteOrder = iota
    LittleEndianByteOrder = iota
)

var (
    ifdLogger = log.NewLogger("exifjpeg.ifd")
)

type IfdByteOrder int

func (ibo IfdByteOrder) IsBigEndian() bool {
    return ibo == BigEndianByteOrder
}

func (ibo IfdByteOrder) IsLittleEndian() bool {
    return ibo == LittleEndianByteOrder
}

type Ifd struct {
    data []byte
    byteOrder IfdByteOrder
}

func NewIfd(data []byte, byteOrder IfdByteOrder) *Ifd {
    return &Ifd{
        data: data,
        byteOrder: byteOrder,
    }
}


type IfdVisitor func() error

func (ifd *Ifd) Scan(v IfdVisitor) (err error) {
    fmt.Printf("IFD: Scanning.\n")

    // content := string(data)
    // fmt.Printf("APPDATA (%02X): LEN=(%d)\n>>>>>>>>>>\n%s\n<<<<<<<<<<\n", markerId, len(data), content)

    return nil
}
