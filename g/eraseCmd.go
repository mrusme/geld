package g

import (
  "os"
  "fmt"
  "github.com/spf13/cobra"
  "github.com/gookit/color"
)

var eraseCmd = &cobra.Command{
  Use:   "erase ([flags]) [id]",
  Short: "Erase transaction",
  Long: "Erase tracked transaction.",
  Args: cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    user := GetCurrentUser()
    id := args[0]

    err := database.EraseTransaction(user, id)
    if err != nil {
      fmt.Printf("%s %+v\n", CharError, err)
      os.Exit(1)
    }

    fmt.Printf("%s erased %s\n", CharInfo, color.FgLightWhite.Render(id))
    return
  },
}

func init() {
  rootCmd.AddCommand(eraseCmd)

  var err error
  database, err = InitDatabase()
  if err != nil {
    fmt.Printf("%s %+v\n", CharError, err)
    os.Exit(1)
  }
}
