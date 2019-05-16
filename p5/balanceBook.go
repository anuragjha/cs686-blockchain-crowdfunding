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
	Book p1.MerklePatriciaTrie //key - hashOfPubKey and Value - balance

	// key - Requirement transaction id -||- value - BorrowingTransaction
	Promised map[string]BorrowingTransaction
	mux      sync.Mutex
}

func NewBalanceBook() BalanceBook {
	book := p1.MerklePatriciaTrie{}
	book.Initial()
	//promised := p1.MerklePatriciaTrie{}
	//promised.Initial()

	promised := make(map[string]BorrowingTransaction)

	return BalanceBook{
		Book:     book,
		Promised: promised,
	}
}

func (bb *BalanceBook) BuildBalanceBook(chain p2.Blockchain, fromHeight int32) { // not using fromHeight for now

	log.Println(">>>>>>>>>>>>>>> In BuildBalanceBook  <<<<<<<<<<<<<<<<")
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
				//log.Println("\nIn BuildBalanceBook - txJson to be converted to Tx ---> \n", txjson)
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

	log.Println(">>>>>>>>>>>>>>> In UpdateABalanceBookForTx  Overall start <<<<<<<<<<<<<<<<")

	if tx.TxType == "default" {
		log.Println("<<<<<<<<<<<<<<<<<<<<<< In UpdateABalanceBookForTx - default !!!!!!!! !!!!!! - tx id - ", tx.Id,
			">>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		fromKey := bb.GetKey(tx.From.PublicKey)
		bb.UpdateABalanceInBook(fromKey, -tx.Tokens)

		toKey := bb.GetKey(tx.To.PublicKey)
		bb.UpdateABalanceInBook(toKey, tx.Tokens)

	} else if tx.TxType == "start" { //default transaction read
		log.Println("<<<<<<<<<<<<<<<<<<<<<< In UpdateABalanceBookForTx - start !!!!!!!! !!!!!! - tx id - ", tx.Id,
			">>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		toKey := bb.GetKey(tx.To.PublicKey)
		bb.UpdateABalanceInBook(toKey, tx.Tokens)

	} else if tx.TxType == "" && tx.To.Label != "" && tx.From.Label != "" && tx.ToTxId == "" { // A to B token transfer // pay
		log.Println("<<<<<<<<<<<<<<<<<<<<<< In UpdateABalanceBookForTx -A to B token transfer !!!!!!!! !!!!!! - tx id - ", tx.Id,
			">>>>>>>>>>>>>>>>>>>>>>>>>>>>")

		fromKey := bb.GetKey(tx.From.PublicKey)
		bb.UpdateABalanceInBook(fromKey, -tx.Tokens)

		toKey := bb.GetKey(tx.To.PublicKey)
		bb.UpdateABalanceInBook(toKey, tx.Tokens)
		//bb.UpdateABalanceInBook(toKey, -tx.Fees)

	} else if tx.TxType == "req" && tx.ToTxId == "" /*tx.To.Label == "" && tx.ToTxId == "" && tx.From.Label != ""*/ { // A 's Req Tx // Requirement

		log.Println("<<<<<<<<<<<<<<<<<<<<<< In UpdateABalanceBookForTx - // A 's Req Tx // Requirement !!!!!!!! !!!!!! - tx id - ", tx.Id,
			">>>>>>>>>>>>>>>>>>>>>>>>>>>>")
		bb.PutTxInPromised(tx)

	} else if tx.TxType == "promise" && tx.ToTxId != "" /*tx.ToTxId != "" && tx.From.Label != "" && tx.To.Label == ""*/ { // B to Req Tx

		log.Println("<<<<<<<<<<<<<<<<<<<<<< In UpdatePromiseBookForTx - // B +++++> Req Tx !!!!!!!! !!!!!! - tx id - ", tx.Id,
			">>>>>>>>>>>>>>>>>>>>>>>>>>>>", tx.Tokens)
		bb.UpdateABalanceInPromised(tx)
	}

}

func (bb *BalanceBook) UpdateABalanceInBook(PublicKeyHashStr string, updateBalanceBy float64) {
	bb.mux.Lock()
	defer bb.mux.Unlock()

	log.Println(">>>>>>>>>>>>>>> In UpdateABalanceInBook  <<<<<<<<<<<<<<<<")
	//pubIdHashStr := GetHashOfPublicKey(pubId) //hashOfPublicKey

	currBalance := bb.GetBalanceFromKey(PublicKeyHashStr)
	//log.Println(">>>>>>>>>>>>>>> In UpdateABalanceInBook - currBalance :", currBalance," <<<<<<<<<<<<<<<<")

	newBalance := currBalance + updateBalanceBy
	//log.Println(">>>>>>>>>>>>>>> In UpdateABalanceInBook - newBalance :", newBalance," <<<<<<<<<<<<<<<<")

	bb.Book.Insert(PublicKeyHashStr, fmt.Sprintf("%f", newBalance))
}

func (bb *BalanceBook) UpdateABalanceInPromised(tx Transaction) {
	bb.mux.Lock()
	defer bb.mux.Unlock()

	log.Println(">>>>>>>>>>>>>>>  In UpdateABalanceInPromised  <<<<<<<<<<<<<<<<")
	log.Println("Transaction being processed : \n", tx.TransactionToJson())
	log.Println("And Promised dataStructure is >---->>>", bb.Promised)
	// Promised --->
	// key - Requirement transaction id -||- value - BorrowingTransaction
	if _, ok := bb.Promised[tx.ToTxId]; !ok {

		btx := bb.Promised[tx.ToTxId]

		//btx.PromisesMade[tx.Id] = tx
		btx.PromisesMade = append(btx.PromisesMade, tx)

		log.Println(">>>>>>>>>>>>>>> In UpdateABalanceInPromised - PromisesMade : ", bb.Promised[tx.ToTxId].PromisesMade)

		enough := btx.CheckForEnoughPromises()
		if enough {
			//transfer token from Promised Tx User -to- Req Tx User
			log.Println("Enough Promises ----------> -----------> ------->", enough)

			bb.TransferPromisesMade(btx)

			delete(bb.Promised, tx.ToTxId)
		}

	}
}

func (bb *BalanceBook) TransferPromisesMade(btx BorrowingTransaction) {

	log.Println(">>>>>>>>>>>>>>> In TransferPromisesMade  <<<<<<<<<<<<<<<<")
	for _, ptx := range btx.PromisesMade {
		// for subtract
		KeyForSub := GetKeyForBook(ptx.From.PublicKey)
		log.Println("KeyForSub (should be present in ShowBlance) ----->", KeyForSub)
		bb.UpdateABalanceInBook(KeyForSub, -ptx.Tokens)

		//for add
		KeyForAdd := GetKeyForBook(btx.BorrowingTx.From.PublicKey)
		log.Println("KeyForAdd (should be present in ShowBlance) ----->", KeyForAdd)
		bb.UpdateABalanceInBook(KeyForAdd, ptx.Tokens)
	}

}

func (btx *BorrowingTransaction) CheckForEnoughPromises() bool {

	log.Println(">>>>>>>>>>>>>>> In CheckForEnoughPromises  <<<<<<<<<<<<<<<<")

	valueNeeded := btx.BorrowingTx.Tokens

	valuePromised := float64(0)

	for _, ptx := range btx.PromisesMade {
		valuePromised += ptx.Tokens
	}

	log.Println("Value Needed - ---- > ", valueNeeded)
	log.Println("Value Promised - ---- > ", valuePromised)
	if valueNeeded <= valuePromised {
		return true
	}
	return false
}

func (bb *BalanceBook) PutTxInPromised(tx Transaction) {
	bb.mux.Lock()
	defer bb.mux.Unlock()

	log.Println(">>>>>>>>>>>>>>> In PutTxInPromised  <<<<<<<<<<<<<<<<")

	// Promised --->
	// key - Requirement transaction in json -||- value - MPT of <To Txs in json, value>
	if _, ok := bb.Promised[tx.Id]; !ok {

		btx := NewBorrowingTransaction(tx)

		bb.Promised[tx.Id] = btx

		log.Println(">>>>>>>>>>>>>>> In PutTxInPromised - bb.Promised for btx : ", bb.Promised[tx.Id].BorrowingTx.Tokens)
	}

}

func GetKeyForBook(publicKey *rsa.PublicKey) string {
	hash := sha3.Sum256(publicKey.N.Bytes())
	hashKey := hex.EncodeToString(hash[:])
	return hashKey
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
	str := ""
	for _, btx := range bb.Promised {
		str += "\nBTX : "
		str += btx.BorrowingTxId
		str += " of amt : "
		str += fmt.Sprintf("%f", btx.BorrowingTx.Tokens)
		for _, promise := range btx.PromisesMade {
			str += "\n\t\tPTX : "
			str += promise.From.Label
			str += "promises : "
			str += fmt.Sprintf("%f", promise.Tokens)
		}
	}
	return str
}
