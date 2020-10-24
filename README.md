![geld](documentation/geld.png)
-------------------------------

Geld, zählen. A command line tool for tracking money & budgets.

[Download the latest version for macOS, Linux, FreeBSD, NetBSD, OpenBSD & Plan9 here](https://github.com/mrusme/geld/releases/latest).


## Build

```sh
make
```

**Info**: This will build using the version 0.0.0. You can prefix the `make` 
command with `VERSION=x.y.z` and set `x`, `y` and `z` accordingly if you want 
the version in `geld --help` to be a different one.


## Usage

Please make sure to `export GELD_DB=~/.config/geld.db` (or whatever location 
you would like to have the geld database at).


### List transactions

```sh
geld list --help
```

#### Examples:

List all transactions:

```sh
geld list
```

List all transactions since a specific date:

```sh
geld list --since "Oct 2, 2020"
```

List all transactions and add the total amount:

```sh
geld list --total
```


### Import transactions

During import, *geld* will create SHA1 sums for every transaction, which allows 
it to identify every imported transaction. This way *geld* won't import the 
exact same transaction twice. This means that if you import periodical exports 
that might contain duplicate transactions you don't need to worry about them 
messing up your database.

```sh
geld import --help
```

The following formats are supported as of right now:


#### `csv`: Generic CSV

It is possible to import generic CSV exports by specifying in which columns
`geld` can find individual properties.

#### Examples:

Import a generic CSV export:

```sh
geld import --format csv ./generic.csv \
  --csv-col-date 1 \
  --csv-col-value 8 \
  --csv-col-reference 5 \
  --csv-col-sender-receiver 3 \
  --csv-format-date "2.1.2006" \
  --csv-delimiter ";" \
  --csv-value-decimal-separator ","
```


#### `revolut`: Revolut CSV

It is possible to import CSV exports from [Revolut](https://revolut.com). To
export a CSV, open the Revolut iOS or Android app, select your account, click
the `...` button (right next to the `+ Add money` and `-> Send` buttons),
choose `Statement` from the popup menu and select `Excel` on the top of the
statement screen.

#### Examples:

Import a Revolut CSV export:

```sh
geld import --format revolut ./revolut.csv
```
