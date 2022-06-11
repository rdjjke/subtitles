package main

import (
    "fmt"
)

func main() {
    args := ParseArgs()
    fmt.Printf("%+v", args)
}
