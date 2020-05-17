module github.com/dsoprea/go-jpeg-image-structure

go 1.13

// Development only
// replace github.com/dsoprea/go-utility => ../go-utility
// replace github.com/dsoprea/go-exif/v2 => ../go-exif/v2

require (
	github.com/dsoprea/go-exif/v2 v2.0.0-20200517080529-c9be4b30b064
	github.com/dsoprea/go-logging v0.0.0-20200502201358-170ff607885f
	github.com/dsoprea/go-utility v0.0.0-20200512094054-1abbbc781176
	github.com/jessevdk/go-flags v1.4.0
	golang.org/x/net v0.0.0-20200513185701-a91f0712d120 // indirect
	gopkg.in/yaml.v2 v2.3.0 // indirect
)
