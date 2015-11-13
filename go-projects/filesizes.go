// Copyright (c) 2015, Tommy Gildseth
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice,
//    this list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
// ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
// LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
// CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
// SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
// INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
// CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
// ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.
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
