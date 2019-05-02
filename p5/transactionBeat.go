package p5

type TransactionBeat struct {
	Tx      Transaction
	FromPid *PublicIdentity
	TxSig   []byte
	Hops    int
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
