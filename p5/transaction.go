package p5

import (
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"log"
	"strconv"
	"sync"
	"time"
)

type Transaction struct {
	Id        string
	From      *PublicIdentity
	To        *PublicIdentity
	Tokens    float32
	Timestamp time.Time
}

//this is for mining
type TransactionPool struct {
	Pool map[Transaction]bool //actually a set
	mux  sync.Mutex
}

type TransactionBeat struct {
	Tx      Transaction
	FromPid *PublicIdentity
	TxSig   []byte
}

func NewTransaction(from *PublicIdentity, to *PublicIdentity, tokens float32) Transaction {
	tx := Transaction{
		From:      from,
		To:        to,
		Tokens:    tokens,
		Timestamp: time.Now(),
	}

	sum := sha3.Sum256([]byte(tx.ToString()))
	tx.Id = hex.EncodeToString(sum[:])

	return tx
}

func (tx *Transaction) ToString() string {
	str := tx.From.PublicKey.N.String() + tx.From.PublicKey.N.String() +
		strconv.FormatFloat(float64(tx.Tokens), 'f', -1, 64) + time.Now().String()
	return str
}

func CreateTxSig(tx Transaction, fromSid *Identity) []byte {
	return fromSid.GenSignature(TransactionToJsonByteArray(tx))
}

func TransactionToJsonByteArray(tx Transaction) []byte {
	txJson, err := json.Marshal(tx)
	if err != nil {
		log.Println("in TransactionToJsonByteArray : Error in marshalling Tx : ", err)
	}

	return txJson
}

func NewTransactionBeat(tx Transaction, fromPid *PublicIdentity, fromSig []byte) TransactionBeat {
	return TransactionBeat{
		Tx:      tx,
		FromPid: fromPid,
		TxSig:   fromSig,
	}
}

func CreateTransactionBeat(tx Transaction, sid *Identity) TransactionBeat {

	pid := sid.GetMyPublicIdentity()

	return TransactionBeat{
		Tx:      tx,
		FromPid: &pid,
		TxSig:   CreateTxSig(tx, sid),
	}
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
