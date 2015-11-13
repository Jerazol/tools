package main

import (
  "bufio"
  "path/filepath"
  "fmt"
  "os"
  "strconv"
  "crypto/rand"
  "math/big"
	"unicode"
	"unicode/utf8"
)

func main() {
  var dict        string = "/usr/share/dict/"
  var lang        string = dict + os.Args[1]
  var words       []string
  var pplen       int64 = 4
  var maxWordLen  int64 = 10
  var err         error
  var dictLength  *big.Int
  //byte[] latin1 = ...
  //byte[] utf8 = new String(latin1, "ISO-8859-1").getBytes("UTF-8");

  if len(os.Args) < 2 {
    fmt.Println("No language specified")
    listAvailDicts(dict)
    os.Exit(1)
  }

  if len(os.Args) == 3 {
    pplen, err = strconv.ParseInt(os.Args[2], 10, 0)
  }

  if len(os.Args) == 4 {
    maxWordLen, err = strconv.ParseInt(os.Args[3], 10, 0)
  }

  _, err = os.Stat(lang)
  if os.IsNotExist(err) {
    fmt.Printf("%s is not installed\n", os.Args[1])
    listAvailDicts(dict)
    os.Exit(1)
  }

  words = readLines(lang, maxWordLen)
  fmt.Printf("  - %d\n", pplen)
  fmt.Printf("  - %d\n", len(words))

  dictLength = big.NewInt(int64(len(words)))
  for i := 0; i < int(pplen); i++ {
    var index *big.Int
    index, err = rand.Int(rand.Reader, dictLength)
    word := []byte(words[index.Uint64()])
    fmt.Printf("%s", upperFirst(toUtf8(word)))
  }
  fmt.Println("")
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

func listAvailDicts(dictPath string) {
  fmt.Println("---")
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
func readLines(path string, maxWordLen int64) ([]string) {
  file, err := os.Open(path)
  if err != nil {
    return nil
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    word := scanner.Text()
    if (int64(len(word)) <= maxWordLen) {
      lines = append(lines, word)
    }
  }
  return lines
}
