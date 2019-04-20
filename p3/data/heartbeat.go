package data

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strconv"

	"../../p1"
	"../../p2/block"
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
		Hops:        3,
	}
}

//PrepareHeartBeatData func prepares  and returns heartbeat
func PrepareHeartBeatData(sbc *SyncBlockChain, selfId int32, peerMapBase64 string, addr string) HeartBeatData {

	makeNew := rand.Int() % 2
	//makeNew := 1
	var ifNewBlock bool
	blockJSON := ""

	if makeNew == 1 {
		ifNewBlock = true
		blockJSON = newBlockJson(sbc) // creating a new block's json
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

//newBlockJson func create a new block, inserts it into blockchain and returns json string
func newBlockJson(sbc *SyncBlockChain) string {
	mpt := p1.MerklePatriciaTrie{}

	for i := 0; i <= 5; i++ {
		random := strconv.Itoa(rand.Int() % 10)
		mpt.Insert("Time Now "+random, strconv.Itoa(rand.Int())+"is humbly yours "+os.Args[1])
	}
	b1 := sbc.GenBlock(mpt)
	sbc.Insert(b1)
	return block.EncodeToJSON(&b1)
}

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
	json.Unmarshal([]byte(heartBeatDatajson), &hbd)
	return hbd
}
