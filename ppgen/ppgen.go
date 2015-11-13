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

  lang, pplen, wordMaxLen := setFlags()
  langPath := dicts + *lang

  if len(*lang) == 0 {
    usage(dicts)
    os.Exit(1)
  }

  _, err := os.Stat(langPath)
  if os.IsNotExist(err) {
    fmt.Printf("Language \"%s\" is not installed\n\n", *lang)
    usage(dicts)
    os.Exit(1)
  }

  words := readDict(langPath, *wordMaxLen)
  combinations := math.Pow(float64(len(words)), float64(*pplen))
  fmt.Printf("\n%d words in dictionary giving a total of %e possible passphrases.\n\n", len(words), combinations)

  dictLength := big.NewInt(int64(len(words)))
  for i := 0; i < *pplen; i++ {
    var index *big.Int
    index, err = rand.Int(rand.Reader, dictLength)
    word := []byte(words[index.Uint64()])
    fmt.Printf("%s", ucFirst(toUtf8(word)))
  }
  fmt.Println("")
}


//Convert provided []byte from latin1 to UTF-8
func toUtf8(iso8859_1_buf []byte) string {
    buf := make([]rune, len(iso8859_1_buf))
    for i, b := range iso8859_1_buf {
        buf[i] = rune(b)
    }
    return string(buf)
}


//Uppercase first letter of provided string
func ucFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}


// reads provided dictionary into memory
// and returns a slice of its words.
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
func setFlags() (*string, *int, *int){
  lang        := flag.String("l", "", "Language to use for generating the passphrase")
  numWords    := flag.Int("n", 5, "Number of words to use in the passphrase")
  wordMaxLen  := flag.Int("m", 10, "Max length of words to use in passphrase")
  flag.Parse()

  return lang, numWords, wordMaxLen
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
