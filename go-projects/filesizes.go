package main

import (
	"fmt"
        "path/filepath"
        "os"

)

func visit(path string, f os.FileInfo, err error) error {
        if !f.IsDir() {
                fmt.Printf("Size: %d\n", f.Size())
        } else {
                fmt.Printf("IsDir")
        }
        return nil
}

func main() {
        err := filepath.Walk(os.Args[1], visit)
        fmt.Printf("filepath.Walk() returned %v\n", err)
}



