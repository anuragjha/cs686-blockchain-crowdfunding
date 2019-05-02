package p5

import (
	"encoding/hex"
	"encoding/json"
	"golang.org/x/crypto/sha3"
	"log"
	"strconv"
	"time"
)

type Transaction struct {
	Id        string
	From      *PublicIdentity
	To        *PublicIdentity //if To is empty then its a borrowing tx
	Tokens    float64
	Timestamp time.Time
}

func NewTransaction(from *PublicIdentity, to *PublicIdentity, tokens float64) Transaction {
	tx := Transaction{
		From:      from,
		To:        to,
		Tokens:    tokens,
		Timestamp: time.Now(),
	}

	tx.Id = tx.genId()

	return tx
}

func (tx *Transaction) genId() string {
	str := tx.From.PublicKey.N.String() +
		tx.To.PublicKey.N.String() +
		strconv.FormatFloat(float64(tx.Tokens), 'f', -1, 64) +
		tx.Timestamp.String()
	sum := sha3.Sum256([]byte(str))
	return hex.EncodeToString(sum[:])
}

func (tx *Transaction) Show() string {
	str := tx.Id +
		tx.From.PublicKey.N.String() +
		tx.To.PublicKey.N.String() +
		strconv.FormatFloat(float64(tx.Tokens), 'f', -1, 64) +
		tx.Timestamp.String()
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
