[![GoDoc](https://godoc.org/github.com/bored-ton/ton-go?status.svg)](https://godoc.org/github.com/bored-ton/ton-go)

# TON GO Labrary

Library provides a set of tools for working with TON blockchain.

## TON Center Client
### Tasks

- TON Center API client https://toncenter.com/api/v2/
  - [x] auth
  - [x] accounts
  - [ ] blocks
  - [ ] transactions
  - [ ] run methods
  - [ ] send

### Example
```go
package main

import (
	"log"

	tongo "github.com/bored-ton/ton-go"
	"github.com/shopspring/decimal"
)

func main() {
	client := tongo.NewTonCenterAnonimousClient("https://toncenter.com/api/v2/")
	balance, err := client.Balance("EQCi7s7EYcPzzCrmPlQHck5FerXojlPt32f5vRsbjIPtOLkM")
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("balance: %s", balance.Div(decimal.NewFromFloat(1e9)).String())
}
```
