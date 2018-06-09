package jpegstructure

import (
    "path"
    "testing"

    "io/ioutil"

    "github.com/dsoprea/go-logging"
)

func TestIsJpeg(t *testing.T) {
    filepath := path.Join(assetsPath, "NDM_8901.jpg")

    data, err := ioutil.ReadFile(filepath)
    log.PanicIf(err)

    if IsJpeg(data) != true {
        t.Fatalf("not detected as JPEG")
    }
}
