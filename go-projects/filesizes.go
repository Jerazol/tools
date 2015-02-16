package main

import (
	"fmt"
        "path/filepath"
        "os"
        "math"

)

func main() {
        if len(os.Args) < 1 {
                fmt.Println("Not enough arguments")
                os.Exit(1)
        }

        _, err := os.Stat(os.Args[1])
        if os.IsNotExist(err) {
                fmt.Printf("%s does not exist\n", os.Args[1])
                os.Exit(1)
        }
        
        var sizeRange [64]int
        var maxCount, maxSize int = 0, 0

        filepath.Walk(
                os.Args[1],
                func(path string, f os.FileInfo, err error) error {
                        if !f.IsDir() {
                                tmp := int(findSize(f.Size(), 0))
                                if maxSize < tmp { maxSize = tmp }
                                sizeRange[tmp]++
                                if sizeRange[tmp] > maxCount { maxCount = sizeRange[tmp] }
                        }
                        return nil
                },
        )

        formatTable(sizeRange, maxCount, maxSize)
}


func findSize(number int64, exp float64) float64 {
  if number < (1024*int64(math.Pow(2, exp))) || number == 64 {
    return exp
  }
  
  return findSize(number, exp+1)
}

func formatTable(sizeRange [64]int, maxCount int, maxSize int) {
        for i := 0; i <= maxSize; i++ {
                fmt.Printf("# <%6s: %d\n", formatNum(math.Pow(2,float64(i)), 0), sizeRange[i])
        }
}

func formatNum(size float64, unit int) string {
        units := []string{"KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}
        if size >= 1024 {
                return formatNum(size / 1024, unit+1)
        }

        return fmt.Sprintf("%d%s", int(size), units[unit])
        
}