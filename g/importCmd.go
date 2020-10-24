package g

import (
  "os"
  "fmt"
  "github.com/spf13/cobra"
  "github.com/gookit/color"
  "github.com/cnf/structhash"
)

var csvColType int
var csvColCategory int
var csvColDate int
var csvColValue int
var csvColReference int
var csvColSenderReceiver int

var csvDelimiter string
var csvDateFormat string
var csvValueDecimalSeparator string

func GetValueForMappedField(tx []string, mapping map[string]int, fieldname string) (string) {
  idx := mapping[fieldname]
  txlen := len(tx)
  if idx > 0 && txlen >= idx {
    return tx[(idx-1)]
  }

  return ""
}

func GetTransactionsFromCSV(user string, file string, mapping map[string]int)([]Transaction, error) {
  var transactions []Transaction

  delimiterRunes := []rune(csvDelimiter)
  decimalSeparatorRunes := []rune(csvValueDecimalSeparator)

  txes, err := ImportCSV(file, delimiterRunes[0])
  if err != nil {
    return transactions, err
  }

  for _, tx := range txes {
    transactionSHA1 := structhash.Sha1(tx, 1)

    transaction, err := NewTransaction(
      "",
      GetValueForMappedField(tx, mapping, "type"),
      GetValueForMappedField(tx, mapping, "category"),
      GetValueForMappedField(tx, mapping, "date"),
      csvDateFormat,
      GetValueForMappedField(tx, mapping, "value"),
      decimalSeparatorRunes[0],
      user)
    if err != nil {
      return transactions, err
    }

    transaction.Reference = GetValueForMappedField(tx, mapping, "reference")
    transaction.SenderReceiver = GetValueForMappedField(tx, mapping, "senderReceiver")

    transaction.SHA1 = fmt.Sprintf("%x", transactionSHA1)

    transactions = append(transactions, transaction)
  }

  return transactions, nil
}

func GetTransactionsFromRevolutCSV(user string, file string) ([]Transaction, error) {
  var transactions []Transaction

  rtxes, err := ImportRevolutCSV(file)
  if err != nil {
    return transactions, err
  }

  for _, rtx := range rtxes {
    transactionSHA1 := structhash.Sha1(rtx, 1)

    var rtxType string
    var rtxPaid string
    if rtx.PaidIn == "" && rtx.PaidOut != "" {
      rtxType = TX_TYPE_OUT
      rtxPaid = rtx.PaidOut
    } else if rtx.PaidIn != "" && rtx.PaidOut == "" {
      rtxType = TX_TYPE_IN
      rtxPaid = rtx.PaidIn
    }

    transaction, err := NewTransaction(
      "",
      rtxType,
      rtx.Category,
      rtx.CompletedDate,
      "Jan 2, 2006",
      rtxPaid,
      '.',
      user)
    if err != nil {
      fmt.Printf("%s %+v\n", CharError, err)
      continue
    }

    var rtxExchange *string = nil
    if rtx.ExchangeIn == "" && rtx.ExchangeOut != "" {
      rtxExchange = &rtx.ExchangeOut
    } else if rtx.ExchangeIn != "" && rtx.ExchangeOut == "" {
      rtxExchange = &rtx.ExchangeIn
    }

    if rtxExchange != nil {
      rtxValueExchanged, err := GetDecimalFromValueString(*rtxExchange, '.')
      if err != nil {
        fmt.Printf("%s %+v\n", CharError, err)
        continue
      }

      transaction.ValueExchanged = rtxValueExchanged.StringFixedBank(2)
      transaction.ExchangeRate = rtx.ExchangeRate
    }

    transaction.Reference = rtx.Reference

    transaction.SHA1 = fmt.Sprintf("%x", transactionSHA1)

    transactions = append(transactions, transaction)
  }

  return transactions, nil
}

var importCmd = &cobra.Command{
  Use:   "import ([flags]) [file]",
  Short: "Import transactions",
  Long: "Import transactions from various formats.",
  Args: cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    var transactions []Transaction
    var err error

    user := GetCurrentUser()

    switch(format) {
    case "geld":
      // TODO:
      fmt.Printf("%s not yet implemented\n", CharError)
      os.Exit(1)
    case "csv":
      mapping := make(map[string]int)

      mapping["type"] = csvColType
      mapping["category"] = csvColCategory
      mapping["date"] = csvColDate
      mapping["value"] = csvColValue
      mapping["reference"] = csvColReference
      mapping["senderReceiver"] = csvColSenderReceiver

      transactions, err = GetTransactionsFromCSV(user, args[0], mapping)
      if err != nil {
        fmt.Printf("%s %+v\n", CharError, err)
        os.Exit(1)
      }
    case "revolut":
      transactions, err = GetTransactionsFromRevolutCSV(user, args[0])
      if err != nil {
        fmt.Printf("%s %+v\n", CharError, err)
        os.Exit(1)
      }
    default:
      fmt.Printf("%s specify an import format; see `geld import --help` for more info\n", CharError)
      os.Exit(1)
    }

    sha1List, sha1Err := database.GetImportsSHA1List(user)
    if sha1Err != nil {
        fmt.Printf("%s %+v\n", CharError, sha1Err)
        os.Exit(1)
    }

    for _, transaction := range transactions {
      if id, ok := sha1List[transaction.SHA1]; ok {
        fmt.Printf("%s %s was previously imported as %s; not importing again\n", CharInfo, color.FgLightWhite.Render(transaction.SHA1), color.FgLightWhite.Render(id))
        continue
      }

      importedId, err := database.AddTransaction(user, transaction)
      if err != nil {
        fmt.Printf("%s %s could not be imported: %+v\n", CharError, color.FgLightWhite.Render(transaction.SHA1), color.FgRed.Render(err))
        continue
      }

      fmt.Printf("%s %s was imported as %s\n", CharInfo, color.FgLightWhite.Render(transaction.SHA1), color.FgLightWhite.Render(importedId))
      sha1List[transaction.SHA1] = importedId
    }

    err = database.UpdateImportsSHA1List(user, sha1List)
    if err != nil {
        fmt.Printf("%s %+v\n", CharError, err)
        os.Exit(1)
    }

    return
  },
}

func init() {
  rootCmd.AddCommand(importCmd)
  importCmd.Flags().StringVar(&format, "format", "", "Format to import, possible values: geld, revolut")
  importCmd.Flags().StringVar(&csvDelimiter, "csv-delimiter", ";", "CSV delimiter")
  importCmd.Flags().StringVar(&csvValueDecimalSeparator, "csv-value-decimal-separator", ".", "Decimal separator of the CSV value column")
  importCmd.Flags().StringVar(&csvDateFormat, "csv-format-date", "Jan 2, 2006", "Format of the CSV date column, see https://golang.org/pkg/time/#Parse")
  importCmd.Flags().IntVar(&csvColType, "csv-col-type", 0, "CSV column number of type field")
  importCmd.Flags().IntVar(&csvColCategory, "csv-col-category", 0, "CSV column number of category field")
  importCmd.Flags().IntVar(&csvColDate, "csv-col-date", 0, "CSV column number of date field")
  importCmd.Flags().IntVar(&csvColValue, "csv-col-value", 0, "CSV column number of value field")
  importCmd.Flags().IntVar(&csvColReference, "csv-col-reference", 0, "CSV column number of reference field")
  importCmd.Flags().IntVar(&csvColSenderReceiver, "csv-col-sender-receiver", 0, "CSV column number of sender or receiver field")

  var err error
  database, err = InitDatabase()
  if err != nil {
    fmt.Printf("%s %+v\n", CharError, err)
    os.Exit(1)
  }
}
