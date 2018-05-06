package main

import (
    "os"
    "fmt"

    "io/ioutil"
    "encoding/json"

    "github.com/dsoprea/go-jpeg-image-structure"
    "github.com/dsoprea/go-logging"
    "github.com/dsoprea/go-exif"
    "github.com/jessevdk/go-flags"
)

var (
    options = &struct {
        Filepath string `short:"f" long:"filepath" required:"true" description:"File-path of JPEG image ('-' for STDIN)"`
        Json bool `short:"j" long:"json" description:"Print as JSON"`
    } {}
)


type SegmentResult struct {
    MarkerId byte `json:"marker_id"`
    MarkerName string `json:"marker_name"`
    Offset int `json:"offset"`
    Data []byte `json:"data"`
}


type SegmentIndexItem struct {
    MarkerName string `json:"marker_name"`
    Offset int `json:"offset"`
    Data []byte `json:"data"`
}


func main() {
    _, err := flags.Parse(options)
    if err != nil {
        os.Exit(-1)
    }

    var data []byte
    if options.Filepath == "-" {
        var err error
        data, err = ioutil.ReadAll(os.Stdin)
        log.PanicIf(err)
    } else {
        var err error
        data, err = ioutil.ReadFile(options.Filepath)
        log.PanicIf(err)
    }

    sl, err := jpegstructure.ParseBytesStructure(data)
    log.PanicIf(err)

    var exifTags []jpegstructure.ExifTag
    for _, s := range sl {
        if s.MarkerId < jpegstructure.MARKER_APP0 || s.MarkerId > jpegstructure.MARKER_APP15 {
            continue
        }

        var err error

        exifTags, err = jpegstructure.GetExifData(s.Data)
        if err != nil {
            if log.Is(err, exif.ErrNotExif) == true {
                continue
            }

            log.Panic(err)
        }

        break
    }

    if exifTags == nil {
        log.Panicf("exif data not found")
    }

    if options.Json == true {
        raw, err := json.MarshalIndent(exifTags, "  ", "  ")
        log.PanicIf(err)

        fmt.Println(string(raw))
    } else {
        for i, et := range exifTags {
            fmt.Printf("%2d: IFD=[%s] ID=(0x%02x) NAME=[%s] TYPE=(%d):[%s] VALUE=[%v]", i, et.IfdName, et.TagId, et.TagName, et.TagTypeId, et.TagTypeName, et.Value)

            if et.ChildIfdName != "" {
                fmt.Printf(" CHILD-IFD=[%s]", et.ChildIfdName)
            }

            fmt.Printf("\n")
        }
    }
}
