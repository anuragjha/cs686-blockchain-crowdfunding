package p5

import (
	"fmt"
	"sync"
)

//contains funcs for maintaing wallet

//type Currency struct {
//	Value float64
//	Unit  string
//}

const TOKENUNIT = "pingala"

type Wallet struct {
	Balance map[string]float64
	mux     sync.Mutex
}

func NewWallet() Wallet {
	balance := make(map[string]float64, 1)
	balance[TOKENUNIT] = 0

	//balance[TokenUNIT] = 1001

	return Wallet{
		Balance: balance,
	}
}

func (w *Wallet) Add(value float64) {
	w.mux.Lock()
	defer w.mux.Unlock()

	w.Balance[TOKENUNIT] = w.Balance[TOKENUNIT] + value

}

func (w *Wallet) Subtract(value float64) {
	w.mux.Lock()
	defer w.mux.Unlock()

	w.Balance[TOKENUNIT] = w.Balance[TOKENUNIT] - value

}

func (w *Wallet) Show() string {
	return "Wallet : \n" + fmt.Sprintf("%f", w.Balance[TOKENUNIT]) + TOKENUNIT
	//return showStr
}
