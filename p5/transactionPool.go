package p5

import (
	"bytes"
	"sync"
)

//this is for mining
type TransactionPool struct {
	Pool map[Transaction]bool //actually a set
	mux  sync.Mutex
}

func NewTransactionPool() TransactionPool {
	return TransactionPool{
		Pool: make(map[Transaction]bool),
	}
}

func (txp *TransactionPool) AddToTransactionPool(tx Transaction) {
	txp.mux.Lock()
	defer txp.mux.Unlock()

	txp.Pool[tx] = false
}

func (txp *TransactionPool) DeleteFromTransactionPool(tx Transaction) {
	txp.mux.Lock()
	defer txp.mux.Unlock()

	delete(txp.Pool, tx)
}

func (txp *TransactionPool) Show() string {
	var byteBuf bytes.Buffer

	for tx := range txp.Pool {
		byteBuf.WriteString(tx.Show() + "\n")
	}

	return byteBuf.String()
}

func (txp *TransactionPool) ReadFromTransactionPool(n int) map[Transaction]bool {
	txp.mux.Lock()
	defer txp.mux.Unlock()

	tempMap := make(map[Transaction]bool, n)
	counter := 0
	for tx := range txp.Pool {
		counter++
		if counter < n {
			break
		}
		tempMap[tx] = false
	}
	return tempMap
}
