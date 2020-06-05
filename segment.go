package jpegstructure

import (
	"bytes"
	"errors"
	"fmt"

	"crypto/sha1"
	"encoding/hex"

	"github.com/dsoprea/go-exif/v2"
	"github.com/dsoprea/go-logging"
	"github.com/dsoprea/go-utility/image"
)

var (
	xmpPrefix = []byte("http://ns.adobe.com/xap/1.0/\000")

	// exifPrefix is the prefix found at the top of an EXIF slice. This is JPEG-
	// specific.
	exifPrefix = []byte{'E', 'x', 'i', 'f', 0, 0}
)

var (
	// ErrNoXmp is returned if XMP data was requested but no XMP data was found.
	ErrNoXmp = errors.New("no XMP data")
)

// SofSegment has info read from a SOF segment.
type SofSegment struct {
	// BitsPerSample is the bits-per-sample.
	BitsPerSample byte

	// Width is the image width.
	Width uint16

	// Height is the image height.
	Height uint16

	// ComponentCount is the number of color components.
	ComponentCount byte
}

// String returns a string representation of the SOF segment.
func (ss SofSegment) String() string {

	// TODO(dustin): Add test

	return fmt.Sprintf("SOF<BitsPerSample=(%d) Width=(%d) Height=(%d) ComponentCount=(%d)>", ss.BitsPerSample, ss.Width, ss.Height, ss.ComponentCount)
}

// SegmentVisitor describes a segment-visitor struct.
type SegmentVisitor interface {
	// HandleSegment is triggered for each segment encountered as well as the
	// scan-data.
	HandleSegment(markerId byte, markerName string, counter int, lastIsScanData bool) error
}

// SofSegmentVisitor describes a visitor that is only called for each SOF
// segment.
type SofSegmentVisitor interface {
	// HandleSof is called for each encountered SOF segment.
	HandleSof(sof *SofSegment) error
}

// Segment describes a single segment.
type Segment struct {
	MarkerId   byte
	MarkerName string
	Offset     int
	Data       []byte
}

// SetExif encodes and sets EXIF data into this segment.
func (s *Segment) SetExif(ib *exif.IfdBuilder) (err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	ibe := exif.NewIfdByteEncoder()

	exifData, err := ibe.EncodeToExif(ib)
	log.PanicIf(err)

	len_ := len(exifPrefix)

	s.Data = make([]byte, len_+len(exifData))
	copy(s.Data[0:], exifPrefix)
	copy(s.Data[len_:], exifData)

	return nil
}

// Exif returns an `exif.Ifd` instance for the EXIF data we currently have.
func (s *Segment) Exif() (rootIfd *exif.Ifd, data []byte, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	len_ := len(exifPrefix)

	rawExif := s.Data[len_:]

	jpegLogger.Debugf(nil, "Attempting to parse (%d) byte EXIF blob (Exif).", len(rawExif))

	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	_, index, err := exif.Collect(im, ti, rawExif)
	log.PanicIf(err)

	return index.RootIfd, rawExif, nil
}

// FlatExif parses the EXIF data and just returns a list of tags.
func (s *Segment) FlatExif() (exifTags []exif.ExifTag, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// TODO(dustin): Add test

	len_ := len(exifPrefix)

	rawExif := s.Data[len_:]

	jpegLogger.Debugf(nil, "Attempting to parse (%d) byte EXIF blob (FlatExif).", len(rawExif))

	exifTags, err = exif.GetFlatExifData(rawExif)
	log.PanicIf(err)

	return exifTags, nil
}

// EmbeddedString returns a string of properties that can be embedded into an
// longer string of properties.
func (s *Segment) EmbeddedString() string {
	h := sha1.New()
	h.Write(s.Data)

	// TODO(dustin): Add test

	digestString := hex.EncodeToString(h.Sum(nil))

	return fmt.Sprintf("OFFSET=(0x%08x %10d) ID=(0x%02x) NAME=[%-5s] SIZE=(%10d) SHA1=[%s]", s.Offset, s.Offset, s.MarkerId, markerNames[s.MarkerId], len(s.Data), digestString)
}

// String returns a descriptive string.
func (s *Segment) String() string {

	// TODO(dustin): Add test

	return fmt.Sprintf("Segment<%s>", s.EmbeddedString())
}

// IsExif returns true if EXIF data.
func (s *Segment) IsExif() bool {
	if s.MarkerId != MARKER_APP1 {
		return false
	}

	// TODO(dustin): Add test

	len_ := len(exifPrefix)

	if len(s.Data) < len_ {
		return false
	}

	if bytes.Compare(s.Data[:len_], exifPrefix) != 0 {
		return false
	}

	return true
}

// IsXmp returns true if XMP data.
func (s *Segment) IsXmp() bool {
	if s.MarkerId != MARKER_APP1 {
		return false
	}

	// TODO(dustin): Add test

	len_ := len(xmpPrefix)

	if len(s.Data) < len_ {
		return false
	}

	if bytes.Compare(s.Data[:len_], xmpPrefix) != 0 {
		return false
	}

	return true
}

// FormattedXml returns a formatted XML string. This only makes sense for a
// segment comprised of XML data (like XMP).
func (s *Segment) FormattedXmp() (formatted string, err error) {
	defer func() {
		if state := recover(); state != nil {
			err = log.Wrap(state.(error))
		}
	}()

	// TODO(dustin): Add test

	if s.IsXmp() != true {
		log.Panicf("not an XMP segment")
	}

	len_ := len(xmpPrefix)

	raw := string(s.Data[len_:])

	formatted, err = FormatXml(raw)
	log.PanicIf(err)

	return formatted, nil
}

var (
	// Enforce interface conformance.
	_ riimage.MediaContext = new(Segment)
)
