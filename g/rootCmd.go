package g

import (
  "fmt"
  "github.com/spf13/cobra"
  "os"
)

var database *Database

var txType string
var txCategory string
var txDate string
var txValue string
var txReference string
var txSenderReceiver string

var since string
var until string

var format string

const(
  CharTrack = " ▶"
  CharFinish = " ■"
  CharErase = " ◀"
  CharError = " ▲"
  CharInfo = " ●"
  CharMore = " ◆"
)

var rootCmd = &cobra.Command{
  Use:   "geld",
  Short: "Command line Geldverwaltung",
  Long:  `A command line money tracker.`,
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Printf("%s %+v\n", CharError, err)
    os.Exit(-1)
  }
}

func init() {
  cobra.OnInitialize(initConfig)
}

func initConfig() {
}
