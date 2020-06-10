package main

import (
    "fmt"
    "os"

    "encoding/json"
    "io/ioutil"

    "github.com/dsoprea/go-iptc"
    "github.com/dsoprea/go-logging"
    "github.com/jessevdk/go-flags"

    "github.com/dsoprea/go-jpeg-image-structure"
)

// TODO(dustin): Add comments to all of these structs.

var (
    options = &struct {
        Filepath       string `short:"f" long:"filepath" required:"true" description:"File-path of JPEG image ('-' for STDIN)"`
        JsonAsList     bool   `short:"l" long:"json-list" description:"Print segments as a JSON list"`
        JsonAsObject   bool   `short:"o" long:"json-object" description:"Print segments as a JSON object"`
        IncludeData    bool   `short:"d" long:"data" description:"Include actual JPEG data (only with JSON)"`
        Verbose        bool   `short:"v" long:"verbose" description:"Enable logging verbosity"`
        JustXmp        bool   `short:"x" long:"just-xmp" description:"Just print raw XMP XML. Fails if not present."`
        JustFullIptc   bool   `short:"i" long:"just-full-iptc" description:"Just print raw IPTC data. Fails if not present."`
        JustSimpleIptc bool   `short:"s" long:"just-simple-iptc" description:"Just print raw IPTC data. Omit non-standard tags, omit non-human-readable text, omit repeated tags). Fails if not present."`
    }{}
)

type segmentResult struct {
    MarkerId   byte   `json:"marker_id"`
    MarkerName string `json:"marker_name"`
    Offset     int    `json:"offset"`
    Data       []byte `json:"data"`
    Length     int    `json:"length"`
}

type segmentIndexItem struct {
    Offset int    `json:"offset"`
    Data   []byte `json:"data"`
    Length int    `json:"length"`
}

func main() {
    _, err := flags.Parse(options)
    if err != nil {
        os.Exit(-1)
    }

    if options.Verbose == true {
        scp := log.NewStaticConfigurationProvider()
        scp.SetLevelName(log.LevelNameDebug)

        log.LoadConfiguration(scp)

        cla := log.NewConsoleLogAdapter()
        log.AddAdapter("console", cla)
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

    jmp := jpegstructure.NewJpegMediaParser()

    intfc, parseErr := jmp.ParseBytes(data)

    // If there was an error *and* we got back some segments, print the segments
    // before panicing.
    if intfc == nil && parseErr != nil {
        log.Panic(parseErr)
    }

    sl := intfc.(*jpegstructure.SegmentList)

    if options.JustXmp == true {
        _, s, err := sl.FindXmp()
        log.PanicIf(err)

        xml, err := s.FormattedXmp()
        log.PanicIf(err)

        fmt.Println(xml)

        os.Exit(0)
    }

    if options.JustSimpleIptc == true {
        tags, err := sl.Iptc()
        log.PanicIf(err)

        distilled := iptc.GetSimpleDictionaryFromParsedTags(tags)
        sorted := jpegstructure.SortStringStringMap(distilled)

        for _, pair := range sorted {
            fmt.Printf("%s: %s\n", pair[0], pair[1])
        }

        os.Exit(0)
    } else if options.JustFullIptc == true {
        tags, err := sl.Iptc()
        log.PanicIf(err)

        distilled := iptc.GetDictionaryFromParsedTags(tags)
        sorted := jpegstructure.SortStringStringMap(distilled)

        for _, pair := range sorted {
            fmt.Printf("%s: %s\n", pair[0], pair[1])
        }

        os.Exit(0)
    }

    segments := make([]segmentResult, len(sl.Segments()))
    segmentIndex := make(map[string][]segmentIndexItem)

    for i, s := range sl.Segments() {
        var data []byte
        if (options.JsonAsList == true || options.JsonAsObject == true) && options.IncludeData == true {
            data = s.Data
        }

        segments[i] = segmentResult{
            MarkerId:   s.MarkerId,
            MarkerName: s.MarkerName,
            Offset:     s.Offset,
            Length:     len(s.Data),
            Data:       data,
        }

        sii := segmentIndexItem{
            Offset: s.Offset,
            Length: len(s.Data),
            Data:   data,
        }

        if grouped, found := segmentIndex[s.MarkerName]; found == true {
            segmentIndex[s.MarkerName] = append(grouped, sii)
        } else {
            segmentIndex[s.MarkerName] = []segmentIndexItem{sii}
        }
    }

    if parseErr != nil {
        fmt.Printf("JPEG Segments (incomplete due to error):\n")
        fmt.Printf("\n")

        sl.Print()

        fmt.Printf("\n")

        log.Panic(parseErr)
    }

    if options.JsonAsList == true {
        raw, err := json.MarshalIndent(segments, "", "  ")
        log.PanicIf(err)

        fmt.Println(string(raw))
    } else if options.JsonAsObject == true {
        raw, err := json.MarshalIndent(segmentIndex, "", "  ")
        log.PanicIf(err)

        fmt.Println(string(raw))
    } else {
        fmt.Printf("JPEG Segments:\n")
        fmt.Printf("\n")

        sl.Print()

        sl.FindXmp()
    }
}
