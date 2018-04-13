package exifjpeg

import (
    "fmt"
)

func PrintBytes(data []byte) {
    for _, b := range data {
        fmt.Printf("%02X ", b)
    }

    fmt.Printf("\n")
}
