package g

import (
  "os"
  "fmt"
  "github.com/spf13/cobra"
  "github.com/gookit/color"
  "github.com/cnf/structhash"
)

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

    transaction, err := NewTransaction("", rtxType, rtx.Category, rtx.CompletedDate, rtxPaid, user)
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
      rtxValueExchanged, err := GetDecimalFromValueString(*rtxExchange)
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

  var err error
  database, err = InitDatabase()
  if err != nil {
    fmt.Printf("%s %+v\n", CharError, err)
    os.Exit(1)
  }
}
