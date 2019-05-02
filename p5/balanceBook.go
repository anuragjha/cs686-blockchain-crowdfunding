package p5

import (
	"../p1"
	"fmt"
	//"go/types"
	"log"
	"strconv"
	"sync"
	//"sync"
)

type BalanceBook struct {
	Book p1.MerklePatriciaTrie //key - hashOfPubKey and Value - balance
	mux  sync.Mutex
}

func NewBalanceBook() BalanceBook {
	book := p1.MerklePatriciaTrie{}
	book.Initial()

	return BalanceBook{
		Book: book,
	}
}

func (bb *BalanceBook) UpdateABalanceInBook(PublicKeyStr string, updateBalanceBy float64) {
	bb.mux.Lock()
	defer bb.mux.Unlock()

	//pubIdHashStr := GetHashOfPublicKey(pubId) //hashOfPublicKey

	currBalance := bb.GetBalance(PublicKeyStr)

	newBalance := currBalance + updateBalanceBy

	bb.Book.Insert(PublicKeyStr, fmt.Sprintf("%f", newBalance))
}

func (bb *BalanceBook) GetBalance(PublicKeyStr string) float64 {

	//pubIdHashStr := GetHashOfPublicKey(pubId) //hashOfPublicKey
	balance, err := bb.Book.Get(PublicKeyStr)
	if err != nil {
		log.Println("Error In GetBalance returning 0 , err : ", err)
		return 0
	}

	bal, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		log.Println("Error in strconv from string to float returning 0, err :", err)
		return 0
	}
	return bal

}

func (bb *BalanceBook) IsBalanceEnough(PublicKeyStr string, balanceNeeded float64) bool {
	currentBalance := bb.GetBalance(PublicKeyStr)
	if currentBalance >= balanceNeeded {
		return true
	}
	return false
}

func (bb *BalanceBook) Show() string {
	return bb.Book.String()
}
