package g

import (
  "errors"
  "strings"
  "time"
  "fmt"
  "github.com/gookit/color"
  "github.com/shopspring/decimal"
)

type Transaction struct {
  ID             string      `json:"-"`
  Type           string      `json:"type,omitempty"`
  Category       string      `json:"category,omitempty"`
  Date           time.Time   `json:"date,omitempty"`
  DateValue      time.Time   `json:"valueDate,omitempty"`
  Value          string      `json:"value,omitempty"`
  ValueExchanged string      `json:"valueExchanged,omitempty"`
  ExchangeRate   string      `json:"exchangeRate,omitempty"`
  Reference      string      `json:"reference,omitempty"`
  User           string      `json:"user,omitempty"`

  SHA1           string      `json:"-"`
}

const(
  TX_TYPE_IN  = "in"
  TX_TYPE_OUT = "out"
)

func NewTransaction(
  id string,
  typ string,
  category string,
  date string,
  value string,
  user string) (Transaction, error) {
  var err error

  newTransaction := Transaction{}

  newTransaction.ID = id
  newTransaction.Type = typ
  newTransaction.Category = category

  if date == "" {
    newTransaction.Date = time.Now()
  } else {
    // TODO: Implement parsing that recognizes different formats
    newTransaction.Date, err = time.Parse("Jan 2, 2006", date)
    if err != nil {
      return newTransaction, err
    }
  }

  newTransaction.DateValue = newTransaction.Date

  dec, err := GetDecimalFromValueString(value)
  if err != nil {
    return newTransaction, err
  }

  if dec.IsNegative() {
    newTransaction.Type = TX_TYPE_OUT
    dec = dec.Mul(decimal.NewFromInt(-1))
  } else if newTransaction.Type == "" {
    newTransaction.Type = TX_TYPE_IN
  }

  newTransaction.Value = dec.StringFixedBank(2)
  newTransaction.ValueExchanged = ""
  newTransaction.ExchangeRate = ""
  newTransaction.Reference = ""
  newTransaction.User = user

  return newTransaction, nil
}

func (transaction *Transaction) SetIDFromDatabaseKey(key string) (error) {
  splitKey := strings.Split(key, ":")

  if len(splitKey) < 3 || len(splitKey) > 3 {
    return errors.New("not a valid database key")
  }

  transaction.ID = splitKey[2]
  return nil
}

func (transaction *Transaction) GetValueDecimal() (decimal.Decimal) {
  dec, err := decimal.NewFromString(transaction.Value)
  if err != nil {
    return decimal.NewFromInt(0)
  }

  return dec
}

func (transaction *Transaction) GetTypeVerb() (string) {
  switch(transaction.Type) {
  case TX_TYPE_IN: return "received"
  case TX_TYPE_OUT: return "spent"
  default: return "???"
  }
}

func (transaction *Transaction) GetOutput(full bool) (string) {
  var output string = ""

  if full == false {
    output = fmt.Sprintf("%s %s %s for %s on %s", color.FgGray.Render(transaction.ID), transaction.GetTypeVerb(), color.FgLightWhite.Render(transaction.Value), color.FgLightWhite.Render(transaction.Category), color.FgLightWhite.Render(transaction.Date.Format("2006-01-02")) )
  } else {
  }

  return output
}

func GetFilteredTransactions(transactions []Transaction, typ string, category string, since time.Time, until time.Time) ([]Transaction, error) {
  var filteredTransactions []Transaction

  for _, transaction := range transactions {
    if typ != "" && GetIdFromName(transaction.Type) != GetIdFromName(typ) {
      continue
    }

    if category != "" && GetIdFromName(transaction.Category) != GetIdFromName(category) {
      continue
    }

    if since.IsZero() == false && since.Before(transaction.Date) == false && since.Equal(transaction.Date) == false {
      continue
    }

    if until.IsZero() == false && until.After(transaction.Date) == false && until.Equal(transaction.Date) == false {
      continue
    }

    filteredTransactions = append(filteredTransactions, transaction)
  }

  return filteredTransactions, nil
}
