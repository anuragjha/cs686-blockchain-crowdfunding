package p2

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"

	"golang.org/x/crypto/sha3"

	"./block"
)

// Blockchain struct defines the Blockchain
type Blockchain struct {
	Chain  map[int32][]block.Block `json:"chain"`  //map, key - Height and Value - list of Blocks at this height
	Length int32                   `json:"length"` //equals to the highest block height
}

// BlockchainJson struct for json
type BlockchainJson struct {
	BlockJsonList []string
}

//Initial func initializes the
func (blockchain *Blockchain) Initial() {
	blockchain.Chain = map[int32][]block.Block{}
	blockchain.Length = 0
}

//NewBlockChain func returns a new Blockchain
func NewBlockchain() Blockchain {
	return Blockchain{
		Chain:  make(map[int32][]block.Block),
		Length: 0,
	}
}

// Get func takes height and returns list of blocks at this height
func (blockchain *Blockchain) Get(height int32) ([]block.Block, bool) {
	if height > 0 && blockchain.Length >= height { //sun
		return blockchain.Chain[height], true
	}
	return nil, false
}

// Insert func takes a block, use the height to insert blockhash , but ignore if hash alrady present
func (blockchain *Blockchain) Insert(block block.Block) {

	blockHeight := block.Header.Height
	isValidBlock := false
	//fmt.Println("in insert of blockchain::\tblock height: ", blockHeight)

	if blockchain.Length == 0 && blockHeight == 1 && block.Header.ParentHash == "genesis" {
		isValidBlock = genesis(blockchain, block, blockHeight)
		//fmt.Println("blockHeight == 1)
	} else if blockHeight > 0 && blockHeight <= blockchain.Length { //adding fork
		isValidBlock = addFork(blockchain, block, blockHeight)
		//fmt.Println("blockHeight < blockchain.Length")
	} else if blockHeight > 0 && blockHeight >= blockchain.Length+1 { //can be any height greater than chain length
		isValidBlock = addLength(blockchain, block, blockHeight)
		//fmt.Println("blockHeight == blockchain.Length")
	}

	if isValidBlock == true {
		fmt.Println("block added to blockchain :-) : ", block.Header.Hash, " at height : ", block.Header.Height)
	} else {
		fmt.Println("block will not be added to blockchain :-( : ", block.Header.Hash, " at height : ", block.Header.Height)
	}

}

//genesis func creates the 1st block of blockchain
func genesis(blockchain *Blockchain, block block.Block, blockHeight int32) bool {
	blockchain.Chain[blockHeight] = append(blockchain.Chain[blockHeight], block)
	blockchain.Length++
	// fmt.Println("GENESIS")
	return true
}

//addLength function adds a block such that it increases the length of block
func addLength(blockchain *Blockchain, block block.Block, blockHeight int32) bool {
	//fmt.Println("in Blockchain.go - addLength")
	blockchain.Chain[blockHeight] = append(blockchain.Chain[blockHeight], block)
	blockchain.Length = blockHeight
	return true

}

//addFork method adds a block at previously known height
func addFork(blockchain *Blockchain, block block.Block, blockHeight int32) bool {
	//fmt.Println("in Blockchain.go - addFork")
	blockList := blockchain.Chain[blockHeight]

	isBlockCorrect := true
	for i := range blockList {
		if reflect.DeepEqual(blockList[i].Header.Hash, block.Header.Hash) {
			isBlockCorrect = false
			//fmt.Println("BLOCK will not be added as blockList[i].Header.Hash == block.Header.Hash - therefore"+
			//	" duplicate - at height : ", blockHeight)
			//fmt.Println("blockList[i].Header.Hash = ", blockList[i].Header.Hash)
			//fmt.Println("block.Header.Hash = ", block.Header.Hash)
			break
		}
	}
	if isBlockCorrect == false {
		return false
	}

	blockList = append(blockList, block)
	blockchain.Chain[blockHeight] = blockList // replacing with new blocklist
	return true

}

// EncodeToJSON func creates a list of jsonString
// by iterating over all blocks, generates block jsonString
// calling function should check if the returned string is empty ""
func EncodeToJSON(blockchain *Blockchain) string {

	if len(blockchain.Chain) == 0 {
		return "[]"
	}

	jsonStringBlockchain := "["
	for i := int32(1); i <= blockchain.Length; i++ {
		height := i
		for _, v := range blockchain.Chain[height] {
			thisBlock := v
			jsonStringBlockchain += block.EncodeToJSON(&thisBlock) + ","
		}
	}
	jsonStringBlockchain = jsonStringBlockchain[:len(jsonStringBlockchain)-1]
	jsonStringBlockchain += "]"
	return jsonStringBlockchain
}

// DecodeFromJSON func takes a blockchain instance and a jsonString as input and
// decodes jsonString into type blockchain
func DecodeFromJSON(blockchain *Blockchain, jsonString string) {
	blockJsonList := []block.BlockJson{}

	jerr := json.Unmarshal([]byte(jsonString), &blockJsonList)
	if jerr == nil {
		for i := range blockJsonList {
			jsonBlockByteArray, jerr := json.Marshal(blockJsonList[i])
			if jerr == nil {
				newBlock := block.DecodeFromJSON(string(jsonBlockByteArray))
				blockchain.Insert(newBlock)
			}

		}
	} else {
		fmt.Println("in blockchain DecodeFromJSON - Error in Marshal/Unmarshal : ", jerr)

	}
}

func (bc *Blockchain) Show() string {
	rs := ""
	var idList []int
	for id := range bc.Chain {
		idList = append(idList, int(id))
	}
	sort.Ints(idList)
	for _, id := range idList {
		var hashs []string
		for _, block := range bc.Chain[int32(id)] {
			hashs = append(hashs, block.Header.Hash+"<="+block.Header.ParentHash)
		}
		sort.Strings(hashs)
		rs += fmt.Sprintf("%v: ", id)
		for _, h := range hashs {
			rs += fmt.Sprintf("%s, ", h)
		}
		rs += "\n"
	}
	sum := sha3.Sum256([]byte(rs))
	rs = fmt.Sprintf("This is the BlockChain: %s\n", hex.EncodeToString(sum[:])) + rs
	return rs
}
