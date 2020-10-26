package g

import (
  "fmt"

  "github.com/spf13/cobra"
)

var VERSION string

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Display version of geld",
  Long:  `The version of geld.`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("geld", VERSION)
  },
}
