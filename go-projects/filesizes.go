package main

import (
  "fmt"
  "path/filepath"
  "os"
  "math"

)

func main() {
  var sizeRange [64]int
  var maxSize int = 0

  if len(os.Args) < 2 {
    fmt.Println("No path specified")
    os.Exit(1)
  }

  _, err := os.Stat(os.Args[1])
  if os.IsNotExist(err) {
    fmt.Printf("%s does not exist\n", os.Args[1])
    os.Exit(1)
  }

  filepath.Walk(
    os.Args[1],
    func(path string, f os.FileInfo, err error) error {
      if !f.IsDir() {
        tmp := int(findSize(f.Size(), 0))
        if maxSize < tmp { maxSize = tmp }
        sizeRange[tmp]++
      }
      return nil
    },
  )

  formatTable(sizeRange, maxSize)
}


func findSize(number int64, exp float64) float64 {
  if number < (1024*int64(math.Pow(2, exp))) || number == 64 {
    return exp
  }

  return findSize(number, exp+1)
}

func formatTable(sizeRange [64]int, maxSize int) {
  units := []string{"KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

  for i := 0; i <= maxSize; i++ {
    size := math.Pow(2,float64(i))
    index := math.Floor(math.Log(size) / math.Log(1024))

    fmt.Printf(
      "# <%4d%s: %d\n",
      int(size/math.Pow(1024,index)),
      units[int(index)],
      sizeRange[i],
    )
  }
}
