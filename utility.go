package jpegstructure

import (
	"bytes"
	"fmt"

	"github.com/dsoprea/go-exif"
	"github.com/dsoprea/go-logging"
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

		if i < len(data)-1 {
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

		if i < len(data)-1 {
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

		if i < len(data)-1 {
			_, err := b.WriteString(", ")
			log.PanicIf(err)
		}
	}

	return b.String()
}

type ExifTag struct {
	IfdPath string `json:"ifd_path"`

	TagId   uint16 `json:"id"`
	TagName string `json:"name"`

	TagTypeId   uint16      `json:"type_id"`
	TagTypeName string      `json:"type_name"`
	Value       interface{} `json:"value"`
	ValueBytes  []byte      `json:"value_bytes"`

	ChildIfdPath string `json:"child_ifd_path"`
}

func (et ExifTag) String() string {
	return fmt.Sprintf("ExifTag<IFD-PATH=[%s] TAG-ID=(0x%02x) TAG-NAME=[%s] TAG-TYPE=[%s] VALUE=[%v] VALUE-BYTES=(%d) CHILD-IFD-PATH=[%s]", et.IfdPath, et.TagId, et.TagName, et.TagTypeName, et.Value, len(et.ValueBytes), et.ChildIfdPath)
}

func ParseExifData(exifData []byte) (rootIfd *exif.Ifd, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	_, index, err := exif.Collect(im, ti, exifData)
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

	q := []*exif.Ifd{rootIfd}

	exifTags = make([]ExifTag, 0)

	for len(q) > 0 {
		var ifd *exif.Ifd
		ifd, q = q[0], q[1:]

		ti := exif.NewTagIndex()
		for _, ite := range ifd.Entries {
			tagName := ""

			it, err := ti.Get(ifd.IfdPath, ite.TagId)
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
				IfdPath:      ifd.IfdPath,
				TagId:        ite.TagId,
				TagName:      tagName,
				TagTypeId:    ite.TagType,
				TagTypeName:  exif.TypeNames[ite.TagType],
				Value:        value,
				ValueBytes:   valueBytes,
				ChildIfdPath: ite.ChildIfdPath,
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
