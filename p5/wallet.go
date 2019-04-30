package p5

//contains funcs for maintaing wallet

type Currency struct {
	Value float32
	Unit  string
}

type Wallet struct {
	Balance map[string]Currency
}

func NewWallet() Wallet {
	balance := make(map[string]Currency, 1)

	curr := Currency{
		Value: 1000,
		Unit:  "anucoin",
	}
	balance[curr.Unit] = curr

	return Wallet{
		Balance: balance,
	}
}
