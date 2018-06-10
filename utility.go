package jpegstructure

import (
    "fmt"
    "bytes"

    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-exif"
)

func DumpBytes(data []byte) {
    fmt.Printf("DUMP: ")
    for _, x := range data {
        fmt.Printf("%02x ", x)
    }

    fmt.Printf("\n")
}

func DumpBytesClause(data []byte) {
    fmt.Printf("DUMP: ")

    fmt.Printf("[]byte { ")

    for i, x := range data {
        fmt.Printf("0x%02x", x)

        if i < len(data) - 1 {
            fmt.Printf(", ")
        }
    }

    fmt.Printf(" }\n")
}

func DumpBytesToString(data []byte) string {
    b := new(bytes.Buffer)

    for i, x := range data {
        _, err := b.WriteString(fmt.Sprintf("%02x", x))
        log.PanicIf(err)

        if i < len(data) - 1 {
            _, err := b.WriteRune(' ')
            log.PanicIf(err)
        }
    }

    return b.String()
}

func DumpBytesClauseToString(data []byte) string {
    b := new(bytes.Buffer)

    for i, x := range data {
        _, err := b.WriteString(fmt.Sprintf("0x%02x", x))
        log.PanicIf(err)

        if i < len(data) - 1 {
            _, err := b.WriteString(", ")
            log.PanicIf(err)
        }
    }

    return b.String()
}


type ExifTag struct {
    ParentIfdName string `json:"parent_ifd_name"`
    IfdName string `json:"ifd_name"`

    TagId uint16 `json:"id"`
    TagName string `json:"name"`

    TagTypeId uint16 `json:"type_id"`
    TagTypeName string `json:"type_name"`
    Value interface{} `json:"value"`
    ValueBytes []byte `json:"value_bytes"`

    ChildIfdName string `json:"child_ifd_name"`
}

func (et ExifTag) String() string {
    return fmt.Sprintf("ExifTag<PARENT-IFD=[%s] IFD=[%s] TAG-ID=(0x%02x) TAG-NAME=[%s] TAG-TYPE=[%s] VALUE=[%v] VALUE-BYTES=(%d) CHILD-IFD=[%s]", et.ParentIfdName, et.IfdName, et.TagId, et.TagName, et.TagTypeName, et.Value, len(et.ValueBytes), et.ChildIfdName)
}

func ParseExifData(exifData []byte) (rootIfd *exif.Ifd, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    _, index, err := exif.Collect(exifData)
    log.PanicIf(err)

    return index.RootIfd, nil
}

func GetFlatExifData(exifData []byte) (exifTags []ExifTag, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    rootIfd, err := ParseExifData(exifData)
    log.PanicIf(err)

    q := []*exif.Ifd{ rootIfd }

    exifTags = make([]ExifTag, 0)

    for ; len(q) > 0; {
        var ifd *exif.Ifd
        ifd, q = q[0], q[1:]

        parentIfdName := ""
        if ifd.ParentIfd != nil {
            parentIfdName = ifd.ParentIfd.Identity().IfdName
        }

        ii := ifd.Identity()

        ti := exif.NewTagIndex()
        for _, ite := range ifd.Entries {
            tagName := ""

            it, err := ti.Get(ii, ite.TagId)
            if err != nil {
                // If it's a non-standard tag, just leave the name blank.
                if log.Is(err, exif.ErrTagNotFound) != true {
                    log.PanicIf(err)
                }
            } else {
                tagName = it.Name
            }

            value, err := ifd.TagValue(ite)
            log.PanicIf(err)

            valueBytes, err := ifd.TagValueBytes(ite)
            log.PanicIf(err)

            et := ExifTag{
                ParentIfdName: parentIfdName,
                IfdName: ii.IfdName,
                TagId: ite.TagId,
                TagName: tagName,
                TagTypeId: ite.TagType,
                TagTypeName: exif.TypeNames[ite.TagType],
                Value: value,
                ValueBytes: valueBytes,
                ChildIfdName: ite.ChildIfdName,
            }

            exifTags = append(exifTags, et)
        }

        for _, childIfd := range ifd.Children {
            q = append(q, childIfd)
        }

        if ifd.NextIfd != nil {
            q = append(q, ifd.NextIfd)
        }
    }

    return exifTags, nil
}
