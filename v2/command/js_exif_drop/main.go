package main

import (
	"fmt"
	"os"

	"io/ioutil"

	"github.com/dsoprea/go-jpeg-image-structure/v2"
	"github.com/dsoprea/go-logging"
	"github.com/jessevdk/go-flags"
)

var (
	options = &struct {
		InputFilepath  string `short:"f" long:"input-filepath" required:"true" description:"File-path of JPEG image to read"`
		OutputFilepath string `short:"o" long:"output-filepath" description:"File-path of JPEG image to write (if not provided, then the input JPEG will be used)"`
	}{}
)

func main() {
	_, err := flags.Parse(options)
	if err != nil {
		os.Exit(-1)
	}

	data, err := ioutil.ReadFile(options.InputFilepath)
	log.PanicIf(err)

	jmp := jpegstructure.NewJpegMediaParser()

	intfc, err := jmp.ParseBytes(data)
	log.PanicIf(err)

	sl := intfc.(*jpegstructure.SegmentList)

	wasDropped, err := sl.DropExif()
	log.PanicIf(err)

	fmt.Printf("%v\n", wasDropped)

	if wasDropped == false {
		os.Exit(10)
	}

	outputFilepath := options.OutputFilepath
	if outputFilepath == "" {
		outputFilepath = options.InputFilepath
	}

	f, err := os.OpenFile(outputFilepath, os.O_CREATE|os.O_WRONLY, 0644)
	log.PanicIf(err)

	defer f.Close()

	err = sl.Write(f)
	log.PanicIf(err)
}
