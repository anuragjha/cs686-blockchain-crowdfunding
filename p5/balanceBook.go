package p5

import (
	"../p1"
	"crypto/rsa"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"

	//"go/types"
	"log"
	"strconv"
	"sync"
	//"sync"
)

type BalanceBook struct {
	Book     p1.MerklePatriciaTrie //key - hashOfPubKey and Value - balance
	Promised p1.MerklePatriciaTrie //key - hashOfPubKey and Value - promised amount
	mux      sync.Mutex
}

func NewBalanceBook() BalanceBook {
	book := p1.MerklePatriciaTrie{}
	book.Initial()

	promised := p1.MerklePatriciaTrie{}
	promised.Initial()

	return BalanceBook{
		Book:     book,
		Promised: promised,
	}
}

// GetKey func takes in PublicKey and returns Key for Book
func (bb *BalanceBook) GetKey(publicKey *rsa.PublicKey) string {
	hash := sha3.Sum256(publicKey.N.Bytes())
	hashKey := hex.EncodeToString(hash[:])
	return hashKey
}

func (bb *BalanceBook) UpdateABalanceInBook(PublicKeyHashStr string, updateBalanceBy float64) {
	bb.mux.Lock()
	defer bb.mux.Unlock()

	//pubIdHashStr := GetHashOfPublicKey(pubId) //hashOfPublicKey
	currBalance := bb.GetBalance(PublicKeyHashStr, bb.Book)

	newBalance := currBalance + updateBalanceBy

	bb.Book.Insert(PublicKeyHashStr, fmt.Sprintf("%f", newBalance))
}

func (bb *BalanceBook) GetBalance(PublicKeyHashStr string, thisBook p1.MerklePatriciaTrie) float64 {

	//pubIdHashStr := GetHashOfPublicKey(pubId) //hashOfPublicKey
	//balance, err := bb.Book.Get(PublicKeyHashStr)
	balance, err := thisBook.Get(PublicKeyHashStr)
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

func (bb *BalanceBook) IsBalanceEnough(PublicKeyHashStr string, balanceNeeded float64) bool {
	currentBalance := bb.GetBalance(PublicKeyHashStr, bb.Book) - bb.GetBalance(PublicKeyHashStr, bb.Promised)
	if currentBalance >= balanceNeeded {
		return true
	}
	return false
}

func (bb *BalanceBook) Show() string {
	return bb.Book.String()
}

func (bb *BalanceBook) ShowPromised() string {
	return bb.Promised.String()
}
