package main

import (
    "os"
    "fmt"

    "io/ioutil"
    "encoding/json"

    "github.com/dsoprea/go-jpeg-image-structure"
    "github.com/dsoprea/go-logging"
    "github.com/jessevdk/go-flags"
)

var (
    options = &struct {
        Filepath string `short:"f" long:"filepath" required:"true" description:"File-path of JPEG image ('-' for STDIN)"`
        JsonAsList bool `short:"l" long:"json-list" description:"Print segments as a JSON object"`
        JsonAsObject bool `short:"o" long:"json-object" description:"Print segments as a JSON object"`
        IncludeData bool `short:"d" long:"data" description:"Include actual JPEG data (only with JSON)"`
    } {}
)

type segmentResult struct {
    MarkerId byte `json:"marker_id"`
    MarkerName string `json:"marker_name"`
    Offset int `json:"offset"`
    Data []byte `json:"data"`
}


type segmentIndexItem struct {
    MarkerName string `json:"marker_name"`
    Offset int `json:"offset"`
    Data []byte `json:"data"`
}


func main() {
    _, err := flags.Parse(options)
    if err != nil {
        os.Exit(-1)
    }

    if options.JsonAsList == true && options.JsonAsObject == true {
        fmt.Println("Only -jsonlist *or* -jsonobject can be chosen.")
        os.Exit(-2)
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

    segments := make([]segmentResult, len(sl.Segments()))
    segmentIndex := make(map[byte][]segmentIndexItem)

    for i, s := range sl.Segments() {
        var data []byte
        if (options.JsonAsList == true || options.JsonAsObject == true) && options.IncludeData == true {
            data = s.Data
        }

        segments[i] = segmentResult{
            MarkerId: s.MarkerId,
            MarkerName: s.MarkerName,
            Offset: s.Offset,
            Data: data,
        }

        sii := segmentIndexItem{
            MarkerName: s.MarkerName,
            Offset: s.Offset,
            Data: data,
        }

        if grouped, found := segmentIndex[s.MarkerId]; found == true {
            segmentIndex[s.MarkerId] = append(grouped, sii)
        } else {
            segmentIndex[s.MarkerId] = []segmentIndexItem { sii }
        }
    }

    if options.JsonAsList == true {
        raw, err := json.MarshalIndent(segments, "  ", "  ")
        log.PanicIf(err)

        fmt.Println(string(raw))
    } else if options.JsonAsObject == true {
        raw, err := json.MarshalIndent(segmentIndex, "  ", "  ")
        log.PanicIf(err)

        fmt.Println(string(raw))
    } else {
        fmt.Printf("JPEG Segments:\n")
        fmt.Printf("\n")

        sl.Print()
    }
}
