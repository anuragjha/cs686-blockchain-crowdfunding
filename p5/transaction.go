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
	Fees      float64
	Timestamp time.Time
}

func NewTransaction(from *PublicIdentity, to *PublicIdentity, tokens float64, fees float64) Transaction {
	tx := Transaction{
		From:      from,
		To:        to,
		Tokens:    tokens,
		Fees:      fees,
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

func VerifyTxSig(fromPid *PublicIdentity, tx Transaction, txSig []byte) bool {

	return VerifySingature(fromPid.PublicKey, TransactionToJsonByteArray(tx), txSig)
}

func TransactionToJsonByteArray(tx Transaction) []byte {
	txJson, err := json.Marshal(tx)
	if err != nil {
		log.Println("in TransactionToJsonByteArray : Error in marshalling Tx : ", err)
	}

	return txJson
}

func (tx *Transaction) TransactionToJson() string {
	txJson, err := json.Marshal(tx)
	if err != nil {
		log.Println("in TransactionToJsonByteArray : Error in marshalling Tx : ", err)
	}

	return string(txJson)
}

func DecodeToTransaction(txJson []byte) Transaction {
	tx := Transaction{}
	err := json.Unmarshal(txJson, &tx)
	if err != nil {
		log.Println("Error in unmarshalling Transaction")
	}

	return tx
}

func IsTransactionValid(tx Transaction, balanceBook BalanceBook) bool {

	//balanceBook.Book <hash of PublicKey, Balance Amount>
	//getting hash of public key of tx.From - to get key for balance.Book
	//hash :=sha3.Sum256(tx.From.PublicKey.N.Bytes())
	//hashKey := hex.EncodeToString(hash[:])
	//using hashKey to get the Balance amount
	//balanceStr, err := balanceBook.Book.Get(hashKey)
	//balance, err := strconv.ParseFloat(balanceStr, 64) // todo ?? if ERR then should i make balance zero ???? !!!
	//if err != nil {
	//	return false
	//}

	//if  balance > tx.Tokens {
	//	return true
	//}
	//return false

	return balanceBook.IsBalanceEnough(balanceBook.GetKey(tx.From.PublicKey), tx.Tokens+tx.Fees)
}
