package data

import (
	//"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	//"golang.org/x/crypto/sha3"
	//"math/rand"
	//"os"
	//"strconv"
	//"strings"
	//"time"
	//"../../p1"
	//"../../p2/block"
)

//HeartBeatData struct defines the data to be sent between peers perodically
type HeartBeatData struct {
	IfNewBlock  bool   `json:"ifNewBlock"`
	Id          int32  `json:"id"`
	BlockJson   string `json:"blockJson"`
	PeerMapJson string `json:"peerMapJson"`
	Addr        string `json:"addr"`
	Hops        int32  `json:"hops"`
}

//NewHeartBeatData creates new HeartBeatData
func NewHeartBeatData(ifNewBlock bool, id int32, blockJson string, peerMapJson string, addr string) HeartBeatData {

	return HeartBeatData{
		IfNewBlock:  ifNewBlock,
		Id:          id,
		BlockJson:   blockJson,
		PeerMapJson: peerMapJson,
		Addr:        addr,
		Hops:        2, //todo change to 3
	}
}

//PrepareHeartBeatData func prepares  and returns heartbeat
func PrepareHeartBeatData(sbc *SyncBlockChain, selfId int32, peerMapBase64 string, addr string, makingNew bool, newBlockJson string) HeartBeatData {

	//makeNew := rand.Int() % 2
	//makeNew := 1
	var ifNewBlock bool
	blockJSON := "{}"

	if makingNew == true {
		ifNewBlock = true
		blockJSON = newBlockJson //newBlockJson(sbc) // creating a new block's json

	} else {
		ifNewBlock = false
	}

	return NewHeartBeatData(
		ifNewBlock,
		selfId,
		blockJSON,
		peerMapBase64,
		addr,
	)
}

////newBlockJson func create a new block, inserts it into blockchain and returns json string
//func newBlockJson(sbc *SyncBlockChain) string {
//	mpt := p1.MerklePatriciaTrie{}
//	mpt.Initial()
//
//	for i := 0; i <= 5; i++ {
//		random := strconv.Itoa(rand.Int() % 10)
//		mpt.Insert("Time Now "+random, strconv.Itoa(rand.Int())+"is humbly yours "+os.Args[1])
//	}
//	b1 := sbc.GenBlock(mpt, "0")
//	sbc.Insert(b1)
//	return block.EncodeToJSON(&b1)
//}

//EncodeToJson func encodes HeartBeatData to json string
func (data *HeartBeatData) EncodeToJson() string {

	dataEncodedBytearray, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Here in err condition of EncodeToJson of heartbeat.go")
		return "{}"
	}
	return string(dataEncodedBytearray)
}

//DecodeToHeartBeatData func decodes json string to HeartBeatData
func DecodeToHeartBeatData(heartBeatDatajson string) HeartBeatData {
	hbd := HeartBeatData{}
	err := json.Unmarshal([]byte(heartBeatDatajson), &hbd)
	if err != nil {
		log.Println("Err in DecodeToHeartBeatData in heartbeat.go")
	}
	return hbd
}

////// pow
//var Nonce string //pow

//Initial function to start POW
//func StartTryingNonces(SBC *SyncBlockChain) {
//
//	parentHash := SBC.GetLatestBlocks()[rand.Int()%len(SBC.GetLatestBlocks())].Header.ParentHash //random parent from blocks at latest height
//
//	mpt := GenerateRandomMPT()
//
//	tryingNonces(parentHash ,mpt, 7)
//}

//// tryingNonces is a go routine
//func tryingNonces(parentHash string, mpt string, difficulty int) {
//
//	// y = SHA3(parentHash + nonce + mptRootHash)
//
//	Nonce = InitializeNonce(8)
//
//	for {
//
//		if POW(parentHash, Nonce, mpt, difficulty) {
//
//		}
//		Nonce = NextNonce(Nonce)
//	}
//
//}

////check if POW is satisfied - return true or false
//func POW(parentHash string, Nonce string, mptRootHash string, difficulty int) bool {
//
//	y := sha3.Sum256 ([]byte(parentHash+Nonce+mptRootHash))
//	proof := string(y[:difficulty - 1])
//
//	against := "0000000"
//
//	if strings.Compare(proof, against) == 0 {
//		return true
//	}
//	return false
//}

//// func to get NextNonce //
//func NextNonce(previousNonce string) string {
//
//	bytes, err := hex.DecodeString(previousNonce)
//	if err != nil {
//		return InitializeNonce(8)
//	}
//
//
//	return hex.EncodeToString(bytes)
//}
