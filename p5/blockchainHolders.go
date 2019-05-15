package p5

import (
	"encoding/json"
	"sync"
)

type BlockChainHolders struct {
	Holders map[string]PublicIdentity
	mux     sync.Mutex
}

func NewBlockChainHolders() BlockChainHolders {
	bch := BlockChainHolders{}
	bch.Holders = make(map[string]PublicIdentity)
	return bch
}

func (BCH *BlockChainHolders) AddBlockChainHolder(addr string, pid PublicIdentity) {
	BCH.mux.Lock()
	defer BCH.mux.Unlock()

	BCH.Holders[addr] = pid
}

func (BCH *BlockChainHolders) DeleteBlockChainHolder(addr string) {
	BCH.mux.Lock()
	defer BCH.mux.Unlock()

	delete(BCH.Holders, addr)
}

func (BCH *BlockChainHolders) Show() string {
	jsonStr, err := json.Marshal(BCH.Holders)
	if err != nil {
		jsonStr = []byte("{}")
	}
	return string(jsonStr)
}
