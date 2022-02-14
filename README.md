[![GoDoc](https://godoc.org/github.com/bored-ton/ton-go?status.svg)](https://godoc.org/github.com/bored-ton/ton-go)

# TON GO

Library provides a set of tools for working with TON blockchain.

## TON Center Client
### Tasks

- TON Center API client https://toncenter.com/api/v2/
  - [x] auth
  - [x] accounts
  - [x] blocks
  - [x] transactions
  - [x] run get method
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

```go
package main

import (
	"flag"
	"log"

	tongo "github.com/bored-ton/ton-go"
)

var tonCenterToken = flag.String("token", "", "TonCenter token")

func init() {
	flag.Parse()
	if *tonCenterToken == "" {
		log.Fatalln("TonCenter token is required")
	}
}

func main() {
	client := tongo.NewTonCenterClient("https://toncenter.com/api/v2/", *tonCenterToken)

	req := tongo.RunMethodRequest{
		Address: "EQAZC9dW9sDnlI3CbLR5aDIxd_sgNE-PmlCRjK-H7LNLeUXN",
		Method:  "getstdperiod",
		Stack:   [][]string{},
	}

	output, err := client.RunGetMethod(req)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("ouput: %+v", output)
}
```