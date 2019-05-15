package p3

import (
	"../p1"
	"../p5"
)

var BalanceBook p5.BalanceBook

func InitBalanceBook() {
	//currency //p5
	BalanceBook = p5.NewBalanceBook()
}

//func DefaultTokens() {
//	pubKeyStr := ID.GetMyPublicIdentity().PublicKey.N.String()
//	value := Wallet.Balance[p5.TOKENUNIT]
//	BalanceBook.UpdateABalanceInBook(pubKeyStr, value) //todo p5 todo p5
//}

// func to generate transactionMPT
func GenerateTransactionsMPT() p1.MerklePatriciaTrie {
	mpt := p1.MerklePatriciaTrie{}
	mpt.Initial()

	random := 5 //int((time.Now().UnixNano() / 100000 % 5))

	txs := TxPool.ReadFromTransactionPool(random)
	for _, tx := range txs {

		//todo - check if transaction valid

		TxPool.DeleteFromTransactionPool(tx.Id) //delete from TransactionPool

		mpt.Insert(tx.Id, tx.TransactionToJson())
	}

	return mpt
}

//func MarkTxInTxPoolAsUsed(mpt p1.MerklePatriciaTrie) {
//	usedTxPool := mpt.GetAllKeyValuePairs()
//	for _, txJson := range usedTxPool {
//		tx := p5.DecodeToTransaction([]byte(txJson))
//		TxPool.Pool[tx] = true
//	}
//}
