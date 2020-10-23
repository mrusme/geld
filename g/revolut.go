package g

import (
  "os"
  "io"
  "strings"
  "io/ioutil"
  "encoding/csv"
  "bytes"
)

type RevolutTransaction struct {
  CompletedDate string
  Reference     string
  PaidOut       string
  PaidIn        string
  ExchangeOut   string
  ExchangeIn    string
  Balance       string
  ExchangeRate  string
  Category      string
}

func ImportRevolutCSV(filename string) ([]RevolutTransaction, error) {
  var rtxes []RevolutTransaction

  file, err := os.Open(filename)
  if err != nil {
    return rtxes, err
  }
  defer file.Close()

  fileBytes, err := ioutil.ReadAll(file)
  reader := bytes.NewReader(fileBytes)
  fileBytes, err = ioutil.ReadAll(reader)
  if err != nil {
    return rtxes, err
  }

  r := csv.NewReader(strings.NewReader(string(fileBytes)))
  r.LazyQuotes = true
  r.Comma = ';'
  r.TrimLeadingSpace = true
  for {
    record, err := r.Read()
    if err == io.EOF {
      break
    }
    if err != nil {
      return rtxes, err
    }

    rtx := RevolutTransaction{
      record[0],
      record[1],
      record[2],
      record[3],
      record[4],
      record[5],
      record[6],
      record[7],
      record[8],
    }

    rtxes = append(rtxes, rtx)
  }

  return rtxes, nil
}
