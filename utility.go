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

    ChildIfdName string `json:"child_ifd_name"`
}

func (et ExifTag) String() string {
    return fmt.Sprintf("ExifTag<PARENT-IFD=[%s] IFD=[%s] TAG-ID=(0x%02x) TAG-NAME=[%s] TAG-TYPE=[%s] VALUE=[%v] CHILD-IFD=[%s]", et.ParentIfdName, et.IfdName, et.TagId, et.TagName, et.TagTypeName, et.Value, et.ChildIfdName)
}

func GetExifData(exifData []byte) (exifTags []ExifTag, err error) {
    defer func() {
        if state := recover(); state != nil {
            err = log.Wrap(state.(error))
        }
    }()

    e := exif.NewExif()

    _, index, err := e.Collect(exifData)
    log.PanicIf(err)

    q := make([]*exif.Ifd, 1)
    q[0] = index.RootIfd

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
            it, err := ti.Get(ii, ite.TagId)
            log.PanicIf(err)

            value, err := ifd.TagValue(ite)

            et := ExifTag{
                ParentIfdName: parentIfdName,
                IfdName: ii.IfdName,
                TagId: ite.TagId,
                TagName: it.Name,
                TagTypeId: ite.TagType,
                TagTypeName: exif.TypeNames[ite.TagType],
                Value: value,
                ChildIfdName: ite.ChildIfdName,
            }

            exifTags = append(exifTags, et)
        }

        for _, childIfd := range ifd.Children {
            q = append(q, childIfd)
        }
    }

    return exifTags, nil
}
