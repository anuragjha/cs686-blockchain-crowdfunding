package p5

import (
	"../p1"
	"../p2"
	b "../p2/block"
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

func (bb *BalanceBook) BuildBalanceBook(chain p2.Blockchain, fromHeight int32) { // not using fromHeight for now
	//loop over the blockchain[0] of chains
	var i int32

	for i = fromHeight; i <= chain.Length; i++ {
		blks, found := chain.Get(i)
		if found && len(blks) > 0 {
			blk := b.Block(blks[0])
			//mpt := p1.MerklePatriciaTrie(blk.Value)
			mpt := p1.MerklePatriciaTrie{}
			mpt.Initial()
			mpt = blk.Value
			keyValuePairs := mpt.GetAllKeyValuePairs() //key - txid value - txJson
			//loop over all key valye pairs and collect borrowing txs
			for _, txjson := range keyValuePairs {
				log.Println("\nIn BuildBalanceBook - txJson to be converted to Tx ---> \n", txjson)
				tx := JsonToTransaction(txjson)
				log.Println("\nIn BuildBalanceBook - txJson to be converted to Tx ---> tx ID\n", tx.Id)
				log.Println("\nIn BuildBalanceBook - txJson to be converted to Tx ---> tx Tokens\n", tx.Tokens)

				bb.UpdateABalanceBookForTx(tx) // updating BalanceBook for transaction

				//blk.Header.//todo  blk here  // add a mined by in the block // for now do without tx fess
				//minerKey :=  bb.GetKey(blk.Header.miner)  // assuming miner - public key of miner
			}
		}
	}

}

//func (bb *BalanceBook) UpdateABalanceBookForBlock() { // update bb based on transaction here // todo
//	toKey 	:=  bb.GetKey(tx.To.PublicKey)
//	fromKey := 	bb.GetKey(tx.From.PublicKey)
//
//	bb.UpdateABalanceInBook(toKey, -tx.Tokens)
//	bb.UpdateABalanceInBook(toKey, -tx.Fees)
//	bb.UpdateABalanceInBook(fromKey, tx.Tokens)
//}

func (bb *BalanceBook) UpdateABalanceBookForTx(tx Transaction) { // update bb based on transaction here

	if tx.TxType == "genesis" {
		//toKey 	:=  bb.GetKey(tx.To.PublicKey)
		//bb.UpdateABalanceInBook(toKey, tx.Tokens)

	} else if tx.TxType == "start" { //default transaction read
		toKey := bb.GetKey(tx.To.PublicKey)
		bb.UpdateABalanceInBook(toKey, tx.Tokens)

	} else if tx.To.Label != "" { // A to B token transfer // pay

		fromKey := bb.GetKey(tx.From.PublicKey)
		bb.UpdateABalanceInBook(fromKey, -tx.Tokens)

		toKey := bb.GetKey(tx.To.PublicKey)
		bb.UpdateABalanceInBook(toKey, tx.Tokens)
		//bb.UpdateABalanceInBook(toKey, -tx.Fees)

	} else if tx.To.Label == "" { // A 's Req Tx // Requirement

	}

}

func (bb *BalanceBook) UpdateABalanceInBook(PublicKeyHashStr string, updateBalanceBy float64) {
	bb.mux.Lock()
	defer bb.mux.Unlock()

	//pubIdHashStr := GetHashOfPublicKey(pubId) //hashOfPublicKey

	currBalance := bb.GetBalanceFromKey(PublicKeyHashStr)

	newBalance := currBalance + updateBalanceBy

	bb.Book.Insert(PublicKeyHashStr, fmt.Sprintf("%f", newBalance))
}

// GetKey func takes in PublicKey and returns Key for Book
func (bb *BalanceBook) GetKey(publicKey *rsa.PublicKey) string {
	hash := sha3.Sum256(publicKey.N.Bytes())
	hashKey := hex.EncodeToString(hash[:])
	return hashKey
}

func (bb *BalanceBook) GetBalanceFromPublicKey(publicKey *rsa.PublicKey) float64 {

	//pubIdHashStr := GetHashOfPublicKey(pubId) //hashOfPublicKey
	PublicKeyHashStr := bb.GetKey(publicKey)
	//balance, err := bb.Book.Get(PublicKeyHashStr)
	balance, err := bb.Book.Get(PublicKeyHashStr)
	if err != nil {
		log.Println("GetBalanceFromPublicKey - Error In GetBalance returning 0 , err : ", err)
		return float64(0)
	}

	bal, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		log.Println("Error in strconv from string to float returning 0, err :", err)
		return float64(0)
	}
	return bal

}

func (bb *BalanceBook) GetBalanceFromKey(PublicKeyHashStr string) float64 {

	//pubIdHashStr := GetHashOfPublicKey(pubId) //hashOfPublicKey
	//PublicKeyHashStr := bb.GetKey(publicKey)
	//balance, err := bb.Book.Get(PublicKeyHashStr)
	balance, err := bb.Book.Get(PublicKeyHashStr)
	if err != nil {
		log.Println("GetBalanceFromKey - Error In GetBalance returning 0 , err : ", err)
		return float64(0)
	}

	bal, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		log.Println("Error in strconv from string to float returning 0, err :", err)
		return float64(0)
	}
	return bal

}

//func (bb *BalanceBook) IsBalanceEnough(PublicKeyHashStr string, balanceNeeded float64) bool {
//	currentBalance := bb.GetBalanceFromKey(PublicKeyHashStr) - bb.GetBalanceFromKey(PublicKeyHashStr, bb.Promised)
//	if currentBalance >= balanceNeeded {
//		return true
//	}
//	return false
//}

func (bb *BalanceBook) Show() string {
	return bb.Book.String()
}

func (bb *BalanceBook) ShowPromised() string {
	return bb.Promised.String()
}
