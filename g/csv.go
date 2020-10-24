package g

import (
  "os"
  "io"
  "strings"
  "io/ioutil"
  "encoding/csv"
  "bytes"
)

func ImportCSV(filename string, delimiter rune) ([][]string, error) {
  var txes [][]string

  file, err := os.Open(filename)
  if err != nil {
    return txes, err
  }
  defer file.Close()

  fileBytes, err := ioutil.ReadAll(file)
  reader := bytes.NewReader(fileBytes)
  fileBytes, err = ioutil.ReadAll(reader)
  if err != nil {
    return txes, err
  }

  r := csv.NewReader(strings.NewReader(string(fileBytes)))
  r.LazyQuotes = true
  r.Comma = delimiter
  r.TrimLeadingSpace = true
  for {
    record, err := r.Read()
    if err == io.EOF {
      break
    }
    if err != nil {
      return txes, err
    }

    txes = append(txes, record)
  }

  return txes, nil
}
