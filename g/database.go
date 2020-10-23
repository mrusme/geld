package g

import (
  "os"
  "sort"
  "strings"
  "log"
  "errors"
  "encoding/json"
  "github.com/tidwall/buntdb"
  "github.com/google/uuid"
)

type Database struct {
  DB *buntdb.DB
}

func InitDatabase() (*Database, error) {
  dbfile, ok := os.LookupEnv("GELD_DB")
  if ok == false || dbfile == "" {
    return nil, errors.New("please `export GELD_DB` to the location the geld database should be stored at")
  }

  db, err := buntdb.Open(dbfile)
  if err != nil {
    return nil, err
  }

  db.CreateIndex("category", "*", buntdb.IndexJSON("category"))

  database := Database{db}
  return &database, nil
}

func (database *Database) NewID() (string) {
  id, err := uuid.NewRandom()
  if err != nil {
    log.Fatalln("could not generate UUID: %+v", err)
  }
  return id.String()
}

func (database *Database) AddTransaction(user string, transaction Transaction) (string, error) {
  id := database.NewID()

  transactionJson, jsonerr := json.Marshal(transaction)
  if jsonerr != nil {
    return id, jsonerr
  }

  dberr := database.DB.Update(func(tx *buntdb.Tx) error {
    _, _, seterr := tx.Set(user + ":transaction:" + id, string(transactionJson), nil)
    if seterr != nil {
      return seterr
    }

    return nil
  })

  return id, dberr
}

func (database *Database) GetTransaction(user string, transactionId string) (Transaction, error) {
  var transaction Transaction

  dberr := database.DB.View(func(tx *buntdb.Tx) error {
    value, err := tx.Get(user + ":transaction:" + transactionId, false)
    if err != nil {
      return nil
    }

    json.Unmarshal([]byte(value), &transaction)

    return nil
  })

  return transaction, dberr
}

func (database *Database) UpdateTransaction(user string, transaction Transaction) (string, error) {
  transactionJson, jsonerr := json.Marshal(transaction)
  if jsonerr != nil {
    return transaction.ID, jsonerr
  }

  dberr := database.DB.Update(func(tx *buntdb.Tx) error {
    _, _, seerr := tx.Set(user + ":transaction:" + transaction.ID, string(transactionJson), nil)
    if seerr != nil {
      return seerr
    }

    return nil
  })

  return transaction.ID, dberr
}

func (database *Database) EraseTransaction(user string, id string) (error) {
  dberr := database.DB.Update(func(tx *buntdb.Tx) error {
    _, delerr := tx.Delete(user + ":transaction:" + id)
    if delerr != nil {
      return delerr
    }

    return nil
  })

  return dberr
}

func (database *Database) ListTransactions(user string) ([]Transaction, error) {
  var transactions []Transaction

  dberr := database.DB.View(func(tx *buntdb.Tx) error {
    tx.AscendKeys(user + ":transaction:*", func(key, value string) bool {
      var transaction Transaction
      json.Unmarshal([]byte(value), &transaction)

      transaction.SetIDFromDatabaseKey(key)

      transactions = append(transactions, transaction)
      return true
    })

    return nil
  })

  sort.Slice(transactions, func(i, j int) bool { return transactions[i].Date.Before(transactions[j].Date) })
  return transactions, dberr
}

func (database *Database) GetImportsSHA1List(user string) (map[string]string, error) {
  var sha1List = make(map[string]string)

  dberr := database.DB.View(func(tx *buntdb.Tx) error {
    value, err := tx.Get(user + ":imports:sha1", false)
    if err != nil {
      return nil
    }

    sha1Transactions := strings.Split(value, ",")

    for _, sha1Transaction := range sha1Transactions {
      sha1TransactionSplit := strings.Split(sha1Transaction, ":")
      sha1 := sha1TransactionSplit[0]
      id := sha1TransactionSplit[1]
      sha1List[sha1] = id
    }

    return nil
  })

  return sha1List, dberr
}

func (database *Database) UpdateImportsSHA1List(user string, sha1List map[string]string) (error) {
    var sha1Transactions []string

    for sha1, id := range sha1List {
      sha1Transactions = append(sha1Transactions, sha1 + ":" + id)
    }

    value := strings.Join(sha1Transactions, ",")

    dberr := database.DB.Update(func(tx *buntdb.Tx) error {
      _, _, seterr := tx.Set(user + ":imports:sha1", value, nil)
      if seterr != nil {
        return seterr
      }

      return nil
    })

    return dberr
}

func (database *Database) UpdateCategory(user string, categoryName string, category Category) (error) {
  categoryJson, jsonerr := json.Marshal(category)
  if jsonerr != nil {
    return jsonerr
  }

  categoryId := GetIdFromName(categoryName)

  dberr := database.DB.Update(func(tx *buntdb.Tx) error {
    _, _, sperr := tx.Set(user + ":category:" + categoryId, string(categoryJson), nil)
    if sperr != nil {
      return sperr
    }

    return nil
  })

  return dberr
}

func (database *Database) GetCategory(user string, categoryName string) (Category, error) {
  var category Category
  categoryId := GetIdFromName(categoryName)

  dberr := database.DB.View(func(tx *buntdb.Tx) error {
    value, err := tx.Get(user + ":category:" + categoryId, false)
    if err != nil {
      return nil
    }

    json.Unmarshal([]byte(value), &category)

    return nil
  })

  return category, dberr
}
