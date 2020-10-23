package g

import (
  "os/user"
  "regexp"
  "strings"
  "github.com/shopspring/decimal"
)


func GetCurrentUser() (string) {
  user, err := user.Current()
  if err != nil {
    return "unknown"
  }

  return user.Username
}

func GetIdFromName(name string) string {
  reg, regerr := regexp.Compile("[^a-zA-Z0-9]+")
  if regerr != nil {
      return ""
  }

  id := strings.ToLower(reg.ReplaceAllString(name, ""))

  return id
}

func GetDecimalFromValueString(value string) (decimal.Decimal, error) {
  var regEx *regexp.Regexp
  var dec decimal.Decimal
  var err error

  regEx = regexp.MustCompile("[^\\d.]")
  dec, err = decimal.NewFromFormattedString(value, regEx)
  if err == nil {
    return dec, nil
  }

  return dec, err
}
