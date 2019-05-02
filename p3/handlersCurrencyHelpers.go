package p3

import (
	"../p5"
)

func InitBalanceBook() {
	//currency //p5
	BalanceBook = p5.NewBalanceBook()
}

func InitWallet() {
	Wallet = p5.NewWallet()
}

func DefaultTokens() {
	pubKeyStr := ID.GetMyPublicIdentity().PublicKey.N.String()
	value := Wallet.Balance[p5.TOKENUNIT]
	BalanceBook.UpdateABalanceInBook(pubKeyStr, value) //todo p5 todo p5
}
