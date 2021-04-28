module github.com/dsoprea/go-jpeg-image-structure/v2

go 1.12

// Development only
// replace github.com/dsoprea/go-utility/v2 => ../../go-utility/v2
// replace github.com/dsoprea/go-logging => ../../go-logging
// replace github.com/dsoprea/go-exif/v3 => ../../go-exif/v3
// replace github.com/dsoprea/go-photoshop-info-format => ../../go-photoshop-info-format
// replace github.com/dsoprea/go-iptc => ../../go-iptc

require (
	github.com/dsoprea/go-exif/v3 v3.0.0-20210428042052-dca55bf8ca15
	github.com/dsoprea/go-iptc v0.0.0-20200609062250-162ae6b44feb
	github.com/dsoprea/go-logging v0.0.0-20200710184922-b02d349568dd
	github.com/dsoprea/go-photoshop-info-format v0.0.0-20200609050348-3db9b63b202c
	github.com/dsoprea/go-utility/v2 v2.0.0-20200717064901-2fccff4aa15e
	github.com/go-xmlfmt/xmlfmt v0.0.0-20191208150333-d5b6f63a941b
	github.com/jessevdk/go-flags v1.4.0
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
)
