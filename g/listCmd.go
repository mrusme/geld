package g

import (
  "os"
  "fmt"
  "time"
  "github.com/spf13/cobra"
  "github.com/shopspring/decimal"
)

var listTotalAmount bool

var listCmd = &cobra.Command{
  Use:   "list",
  Short: "List transactions",
  Long: "List all tracked transactions.",
  Run: func(cmd *cobra.Command, args []string) {
    user := GetCurrentUser()

    transactions, err := database.ListTransactions(user)
    if err != nil {
      fmt.Printf("%s %+v\n", CharError, err)
      os.Exit(1)
    }

    var sinceTime time.Time
    var untilTime time.Time

    if since != "" {
      sinceTime, err = time.Parse("Jan 2, 2006", since)
      if err != nil {
        fmt.Printf("%s %+v\n", CharError, err)
        os.Exit(1)
      }
    }

    if until != "" {
      untilTime, err = time.Parse("Jan 2, 2006", until)
      if err != nil {
        fmt.Printf("%s %+v\n", CharError, err)
        os.Exit(1)
      }
    }

    var filteredTransactions []Transaction
    filteredTransactions, err = GetFilteredTransactions(transactions, txType, txCategory, sinceTime, untilTime)
    if err != nil {
      fmt.Printf("%s %+v\n", CharError, err)
      os.Exit(1)
    }

    totalAmount := decimal.NewFromInt(0);
    for _, transaction := range filteredTransactions {
      if transaction.Type == TX_TYPE_IN {
        totalAmount = totalAmount.Add(transaction.GetValueDecimal())
      } else if transaction.Type == TX_TYPE_OUT {
        totalAmount = totalAmount.Sub(transaction.GetValueDecimal())
      }
      fmt.Printf("%s\n", transaction.GetOutput(false))
    }

    if listTotalAmount == true {
      fmt.Printf("\nTOTAL: %s\n\n", totalAmount.StringFixedBank(2))
    }
    return
  },
}

func init() {
  rootCmd.AddCommand(listCmd)
  listCmd.Flags().StringVar(&since, "since", "", "Date to start the list from, e.g. 'Oct 1, 2019'")
  listCmd.Flags().StringVar(&until, "until", "", "Date to list until, e.g. 'Nov 15, 2019'")
  listCmd.Flags().StringVarP(&txType, "type", "t", "", "Type to be listed, possible values: in, out")
  listCmd.Flags().StringVarP(&txCategory, "category", "c", "", "Category to be listed")
  listCmd.Flags().BoolVar(&listTotalAmount, "total", false, "Show total amount for listed transactions")

  var err error
  database, err = InitDatabase()
  if err != nil {
    fmt.Printf("%s %+v\n", CharError, err)
    os.Exit(1)
  }
}
