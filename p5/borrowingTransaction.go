package p5

import (
	"../p1"
	"../p2"
	b "../p2/block"
)

type BorrowingTransaction struct {
	BorrowingTxId string
	BorrowingTx   Transaction
	PromisesMade  []Transaction // key - transaction id (Lending) // todo todo -- changed from map to array -- check start here
}

func NewBorrowingTransaction(tx Transaction) BorrowingTransaction {
	bt := BorrowingTransaction{}
	bt.BorrowingTxId = tx.Id
	bt.BorrowingTx = tx
	bt.PromisesMade = make([]Transaction, 1)

	return bt

}

type BorrowingTransactions struct {
	BorrowingTxs map[string]Transaction // key - BorrowingTxId value - txJson
}

func NewBorrowingTransactions() BorrowingTransactions {
	btxs := BorrowingTransactions{}
	btxs.BorrowingTxs = make(map[string]Transaction)
	return btxs
}

func BuildBorrowingTransactions(chains []p2.Blockchain) BorrowingTransactions {

	btx := NewBorrowingTransactions()
	//loop over the blockchain[00 of chains
	var i int32
	if len(chains) > 0 {
		for i = 1; i <= chains[0].Length; i++ {
			blks, found := chains[0].Get(i)
			if found && len(blks) > 0 {
				blk := b.Block(blks[0])
				mpt := p1.MerklePatriciaTrie(blk.Value)
				keyValuePairs := mpt.GetAllKeyValuePairs() //key - txid value - txJson
				//loop over all key valye pairs and collect borrowing txs
				for _, txjson := range keyValuePairs {
					tx := JsonToTransaction(txjson)
					if tx.To.Label == "" && tx.ToTxId == "" && tx.Tokens > 0 && tx.TxType != "start" && tx.TxType != "default" && tx.From.Label != "" {
						btx.BorrowingTxs[tx.Id] = tx
					}

				}
			}
		}
	}

	return btx

}
