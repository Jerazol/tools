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
)

func main() {
  var dicts       string = "/usr/share/dict/"
  var err         error

  lang, pplen, wordMaxLen := setFlags()
  langPath := dicts + *lang

  if len(*lang) == 0 {
    usage(dicts)
    os.Exit(1)
  }

  _, err = os.Stat(langPath)
  if os.IsNotExist(err) {
    fmt.Printf("Language \"%s\" is not installed\n\n", *lang)
    usage(dicts)
    os.Exit(1)
  }

  words := readLines(langPath, *wordMaxLen)
  fmt.Printf("  - %d\n", *pplen)
  fmt.Printf("  - %d\n", len(words))
  fmt.Printf("  - %d\n", *wordMaxLen)

  dictLength := big.NewInt(int64(len(words)))
  for i := 0; i < *pplen; i++ {
    var index *big.Int
    index, err = rand.Int(rand.Reader, dictLength)
    word := []byte(words[index.Uint64()])
    fmt.Printf("%s", upperFirst(toUtf8(word)))
  }
  fmt.Println("")
}

func setFlags() (*string, *int, *int){
  lang        := flag.String("l", "", "Language to use for generating the passphrase")
  numWords    := flag.Int("n", 5, "Number of words to use in the passphrase")
  wordMaxLen  := flag.Int("m", 10, "Max length of words to use in passphrase")
  flag.Parse()

  return lang, numWords, wordMaxLen
}



func toUtf8(iso8859_1_buf []byte) string {
    buf := make([]rune, len(iso8859_1_buf))
    for i, b := range iso8859_1_buf {
        buf[i] = rune(b)
    }
    return string(buf)
}

func upperFirst(s string) string {
	if s == "" {
		return ""
	}
	r, n := utf8.DecodeRuneInString(s)
	return string(unicode.ToUpper(r)) + s[n:]
}

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

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string, wordMaxLen int) ([]string) {
  file, err := os.Open(path)
  if err != nil {
    return nil
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    word := scanner.Text()
    if (len(word) <= wordMaxLen) {
      lines = append(lines, word)
    }
  }
  return lines
}
