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
  "bufio"
  "path/filepath"
  "fmt"
  "os"
  "crypto/rand"
  "math/big"
  "unicode"
  "unicode/utf8"
  "flag"
  "math"
)

func main() {
  var dicts       string = "/usr/share/dict/"
  lang, ppWords, wordMaxLen, numPhrases := setFlags()
  dictPath := dicts + *lang

  if len(*lang) == 0 {
    usage(dicts)
    os.Exit(1)
  }

  _, err := os.Stat(dictPath)
  if os.IsNotExist(err) {
    fmt.Printf("Language \"%s\" is not installed\n\n", *lang)
    usage(dicts)
    os.Exit(1)
  }

  words   := readDict(dictPath, *wordMaxLen)
  entropy := math.Log2(float64(len(words)))*float64(*ppWords)
  fmt.Printf("%d words in dictionary giving a total of %g bit entropy.\n\n", len(words), entropy)
  
  for index := 0; index < *numPhrases; index++ {
    buildPhrases(words, *ppWords)
  }
  fmt.Println("")
}


// buildPhrases generates a random passphrase of ppWords length using the
// provided words slice
func buildPhrases(words []string, ppWords int) {
  dictLength := big.NewInt(int64(len(words)))
  fmt.Print(": ")
  for i := 0; i < ppWords; i++ {
    index, err := rand.Int(rand.Reader, dictLength)
    if err != nil {
      fmt.Println("Something strange happened")
      return
    }
    word := []byte(words[index.Uint64()])
    fmt.Printf("%s", ucFirst(toUtf8(word)))
  }
  fmt.Println("")
}


// toUtf8 convert provided []byte from latin1 to UTF-8
func toUtf8(iso8859_1_buf []byte) string {
  buf := make([]rune, len(iso8859_1_buf))
  for i, b := range iso8859_1_buf {
    buf[i] = rune(b)
  }
  return string(buf)
}


// ucFirst uppercases the first letter of the provided string
func ucFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}


// readDict reads provided dictionary into memory and returns a slice of its words.
func readDict(path string, wordMaxLen int) ([]string) {
  file, err := os.Open(path)
  if err != nil {
    return nil
  }
  defer file.Close()

  var words []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    word := scanner.Text()
    if (len(word) <= wordMaxLen) {
      words = append(words, word)
    }
  }
  return words
}


// Define program arguments and default values
func setFlags() (*string, *int, *int, *int){
  lang        := flag.String("l", "", "Language to use for generating the passphrase")
  numWords    := flag.Int("n", 5, "Number of words to use in the passphrase")
  wordMaxLen  := flag.Int("m", 10, "Max length of words to use in passphrase")
  numPhrases  := flag.Int("c", 10, "Number of passphrases to generate")
  flag.Parse()

  return lang, numWords, wordMaxLen, numPhrases
}


// Print usage and a list of available dictionaries
func usage(dictPath string) {
  fmt.Println("Usage:")
  flag.PrintDefaults()
  fmt.Println("")
  fmt.Println("Available languages")
  fmt.Println("===================")
  filepath.Walk(
    dictPath,
    func(path string, f os.FileInfo, err error) error {
      if !f.IsDir() && f.Mode()&os.ModeSymlink == 0 && f.Name() != "README.select-wordlist"{
        fmt.Printf("  - %s\n", f.Name())
      }
      return nil
    },
  )
}
