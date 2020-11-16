package g

import (
  "os"
  "fmt"
  "time"
  "strings"
  "github.com/spf13/cobra"
)

var transactionCmd = &cobra.Command{
  Use:   "transaction ([flags]) [id]",
  Short: "Display or update transaction",
  Long: "Display or update tracked transaction.",
  Args: cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    id := args[0]

    transaction, err := database.GetTransaction(txUser, id)
    if err != nil {
      fmt.Printf("%s %+v\n", CharError, err)
      os.Exit(1)
    }

    if txType != "" || txCategory != "" || txDate != "" || txValue != "" || txReference != "" || txSenderReceiver != "" {
      if txType != "" {
        transaction.Type = txType
      }

      if txCategory != "" {
        transaction.Category = txCategory
      }

      if txDate != "" {
        transaction.Date, err = time.Parse("Jan 2, 2006", txDate)
        if err != nil {
          fmt.Printf("%s %+v\n", CharError, err)
          os.Exit(1)
        }

        transaction.DateValue = transaction.Date
      }

      if txValue != "" {
        newTransaction, err := NewTransaction(
          "",
          "",
          "",
          "",
          "",
          txValue,
          '.',
          txUser)
        if err != nil {
          fmt.Printf("%s %+v\n", CharError, err)
          os.Exit(1)
        }
        transaction.Type = newTransaction.Type
        transaction.Value = newTransaction.Value
      }

      if txReference != "" {
        transaction.Reference = strings.Replace(txReference, "\\n", "\n", -1)
      }

      if txSenderReceiver != "" {
        transaction.SenderReceiver = txSenderReceiver
      }

      _, err = database.UpdateTransaction(txUser, transaction)
      if err != nil {
        fmt.Printf("%s %+v\n", CharError, err)
        os.Exit(1)
      }
    }

    fmt.Printf("%s %s\n", CharInfo, transaction.GetOutput(true))
    return
  },
}

func init() {
  rootCmd.AddCommand(transactionCmd)
  transactionCmd.Flags().StringVarP(&txType, "type", "t", "", "Update type, possible values: in, out")
  transactionCmd.Flags().StringVarP(&txCategory, "category", "c", "", "Update category")
  transactionCmd.Flags().StringVarP(&txDate, "date", "d", "", "Update date/time of transaction\n\nUse 'Jan 2, 2006' format.")
  transactionCmd.Flags().StringVarP(&txValue, "value", "v", "", "Update value, e.g. 12.99")
  transactionCmd.Flags().StringVarP(&txReference, "reference", "r", "", "Update reference")
  transactionCmd.Flags().StringVar(&txSenderReceiver, "sender-receiver", "", "Update sender or receiver")

  var err error
  database, err = InitDatabase()
  if err != nil {
    fmt.Printf("%s %+v\n", CharError, err)
    os.Exit(1)
  }
}
