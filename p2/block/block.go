package block

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"../../p1"
	"golang.org/x/crypto/sha3"
)

// Header struct defines the header of each block
type Header struct {
	Height     int32  //`json:"height"`
	Timestamp  int64  //`json:"timestamp"`
	Hash       string //`json:"hash"`
	ParentHash string //`json:"parenthash"`
	Size       int32  // `json:"parenthash"`
}

// Block struct defines the block
type Block struct {
	Header Header                //`json:"header"`
	Value  p1.MerklePatriciaTrie //`json:"merklepatriciatrie"`
}

// BlockJson is a block struct for json
type BlockJson struct {
	Height     int32             `json:"height"`
	Timestamp  int64             `json:"timeStamp"`
	Hash       string            `json:"hash"`
	ParentHash string            `json:"parentHash"`
	Size       int32             `json:"size"`
	MPT        map[string]string `json:"mpt"`
}

// Initial function a Block initializes the block for height, parentHash and Value
func (block *Block) Initial(height int32, parentHash string, value p1.MerklePatriciaTrie) {

	block.Header.Timestamp = time.Now().Unix()
	block.Header.Height = height
	block.Header.ParentHash = parentHash
	block.Value = value
	block.Header.Size = int32(len([]byte(block.Value.String()))) // mpt converted to string and then to byte array
	block.Header.Hash = block.Hash()

}

// DecodeFromJSON func takes json string of type blockJson and converts it into a Block // proxy for : DecodeFromJson
func DecodeFromJSON(jsonString string) Block {

	// block := Block{}

	blockJson := BlockJson{}
	jerr := json.Unmarshal([]byte(jsonString), &blockJson)
	if jerr != nil {
		fmt.Println("block Err : ", jerr)
		return Block{}
	}
	return DecodeToBlock(blockJson.Height,
		blockJson.Timestamp,
		blockJson.Hash,
		blockJson.ParentHash,
		blockJson.Size,
		blockJson.MPT)
}

// DecodeToBlock func creates a type block from from all given parameters
func DecodeToBlock(height int32, timestamp int64, hash string, parentHash string, size int32, keyValueMap map[string]string) Block {

	block := Block{}
	block.Header.Height = height
	block.Header.Timestamp = timestamp
	block.Header.Hash = hash
	block.Header.ParentHash = parentHash
	block.Header.Size = size

	//creating mpt from key - value pairs
	blockMPT := p1.MerklePatriciaTrie{}
	blockMPT.Initial()
	for k, v := range keyValueMap {
		blockMPT.Insert(k, v)
	}
	block.Value = blockMPT

	return block
}

// EncodeToJSON func takes type Block and converts it into json string
func EncodeToJSON(block *Block) string {

	blockForJson := BlockJson{
		Height:     block.Header.Height,
		Timestamp:  block.Header.Timestamp,
		Hash:       block.Header.Hash,
		ParentHash: block.Header.ParentHash,
		Size:       block.Header.Size,
		MPT:        block.Value.GetAllKeyValuePairs(),
	}

	jsonByteArray, jerr := json.Marshal(blockForJson)
	jsonString := "{}"
	if jerr == nil {
		jsonString = string(jsonByteArray)
	}
	return jsonString //empty jsonString if not encoded else some value
}

//Hash func takes an instance of block and hashes it
//hash_str := string(b.Header.Height) + string(b.Header.Timestamp) + b.Header.ParentHash +
//     b.Value.Root + string(b.Header.Size)
func (block *Block) Hash() string {
	var hashStr string

	hashStr = string(block.Header.Height) + string(block.Header.Timestamp) + string(block.Header.ParentHash) +
		string(block.Value.Root) + string(block.Header.Size)

	sum := sha3.Sum256([]byte(hashStr))
	return "HashStart_" + hex.EncodeToString(sum[:]) + "_HashEnd"
}
