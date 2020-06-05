package jpegstructure

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/dsoprea/go-logging"
	"github.com/go-xmlfmt/xmlfmt"
)

// DumpBytes prints the hex for a given byte-slice.
func DumpBytes(data []byte) {
	fmt.Printf("DUMP: ")
	for _, x := range data {
		fmt.Printf("%02x ", x)
	}

	fmt.Printf("\n")
}

// DumpBytesClause prints a Go-formatted byte-slice expression.
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

// DumpBytesToString returns a string of hex-encoded bytes.
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

// DumpBytesClauseToString returns a string of Go-formatted byte values.
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

// FormatXml prettifies XML data.
func FormatXml(raw string) (formatted string, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	formatted = xmlfmt.FormatXML(raw, "  ", "  ")
	formatted = strings.TrimSpace(formatted)

	return formatted, nil
}
