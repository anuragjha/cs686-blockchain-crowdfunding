package p3

import (
	//"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	//"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"

	"../p1"
	"../p2/block"
	//"../p2"
	"../p4"
	"./data"
)

//var TA_SERVER = "http://localhost:6688"
var INIT_SERVER = "http://localhost:6686"

//var REGISTER_SERVER = TA_SERVER + "/peer"
//var REGISTER_SERVER = INIT_SERVER + "/peer"

//SELF_ADDR var BC_DOWNLOAD_SERVER = TA_SERVER + "/upload"
var BC_DOWNLOAD_SERVER = INIT_SERVER + "/upload"

//changes in init for arg of port provided
var SELF_ADDR = "http://localhost:6686"
var SELF_ADDR_PREFIX = "http://localhost:"

// SBC is safe for distributed use
var SBC data.SyncBlockChain

//Peers is the Peer List which is for each node
var Peers data.PeerList

var tryingForHeight int32
var GetNewParent bool

const Difficulty = 5

var ifStarted bool

func init() {
	// This function will be executed before everything else.

	//init coz node not removed from peerlist and receieve heartbeat even before it start()s
	id := Register()
	Peers = data.NewPeerList(id, 32)

	SELF_ADDR = SELF_ADDR_PREFIX + os.Args[1]
	fmt.Println("Node : ", SELF_ADDR)
	//ifStarted = false

}

// Start handler - does Register ID, download BlockChain, start HeartBeat
func Start(w http.ResponseWriter, r *http.Request) {

	if ifStarted == false {
		ifStarted = true

		id := Register()                 //register ID
		Peers = data.NewPeerList(id, 32) //initialize PeerList // 32 sunnit
		SBC = data.NewBlockChain()       //create new Block chain //apr4

		if Peers.GetSelfId() == 6686 {
			mpt := p1.MerklePatriciaTrie{}
			mpt.Insert("First Message", "Anurag's blockchain")
			nonce := p4.FindNonce("genesis", &mpt, Difficulty)
			b1 := SBC.GenBlock(1, "genesis", mpt, nonce)
			SBC.Insert(b1)
		}

		//if Peers.GetSelfId() != 6686 { //download if not 6686
		if SELF_ADDR != INIT_SERVER {

			Peers.Add(INIT_SERVER, int32(6686)) // add Init server to peer list of node
			Download()                          //download BlockChain //apr4 - remove it after testing
		}

		//start HearBeat
		go StartHeartBeat()
		go StartTryingNonces() //pow
	}
	w.WriteHeader(200)
	_, err := w.Write([]byte("started"))
	if err != nil {
		log.Println("Err -  in start - during writing to client")
	}

}

//Show func -  Display peerList and sbc
func Show(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "%s\n%s", Peers.Show(), SBC.Show())
	if err != nil {
		log.Println("Err in show func while writing response")
	}
}

//Canonical func -  Display canonical chain
func Canonical(w http.ResponseWriter, r *http.Request) {

	canonicalChains := p4.GetCanonicalChains(&SBC)

	_, _ = fmt.Fprint(w, "Canonical Chain(s) : \n")
	for i, chain := range canonicalChains {
		_, _ = fmt.Fprint(w, "\nChain #"+strconv.Itoa(i+1))
		_, err := fmt.Fprint(w, "\n", chain.ShowCanonical())
		if err != nil {
			_, _ = fmt.Fprint(w, "ERROr in Canonical")
		}
	}
	//
}

// Register to TA's server, get an ID
func Register() int32 {

	// resp, err := http.Get(REGISTER_SERVER)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer resp.Body.Close()

	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	body := os.Args[1]

	id, err := strconv.Atoi(string(body))
	if err != nil {
		log.Fatal(err)
		return 0
	}

	return int32(id)
}

// Download blockchain from TA server
func Download() {
	resp, err := http.Get(BC_DOWNLOAD_SERVER)
	//resp, err := http.Get("http://localhost:6686/upload/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body) //blockChainJson
	if err != nil {
		log.Fatal(err)
	}

	SBC.UpdateEntireBlockChain(string(body))
}

// Upload blockchain to whoever called this method, return jsonStr
func Upload(w http.ResponseWriter, r *http.Request) {
	blockChainJson, err := SBC.BlockChainToJson()
	if err != nil {
		//data.PrintError(err, "Upload") // todo
		log.Println("Err - in Upload func")
	}
	fmt.Fprint(w, blockChainJson)

	//remove comments above after testing
	//UploadGenesis(w, r)
}

// Upload genesis blockchain to whoever called this method, return jsonStr
func UploadGenesis(w http.ResponseWriter, r *http.Request) {

	nbc := data.NewBlockChain()
	gbl, _ := SBC.Get(1)
	nbc.Insert(gbl[0])

	blockChainJson, err := nbc.BlockChainToJson()
	if err != nil {
		//data.PrintError(err, "Upload") // todo
		log.Println("in Err of Upload Genesis")
	}
	_, err = fmt.Fprint(w, blockChainJson)
	if err != nil {
		log.Println("in Err of Upload Genesis writing response")
	}
}

// UploadBlock func - Upload a block to whoever called this method, return jsonStr
func UploadBlock(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ubHeight, err := strconv.Atoi(vars["height"])
	if err != nil {
		returnCode500(w, r)
	} else {
		ubHash := vars["hash"]
		//fmt.Println("\nuploading block of -\nubHeight : ", ubHeight)
		//fmt.Println("ubHash : ", ubHash, "\n\n")

		uBlock, found := SBC.GetBlock(int32(ubHeight), ubHash)
		if found == false {
			fmt.Println("Err : in Handlers - UploadBlock - found = false - 204")
			returnCode204(w, r)
		} else {
			fmt.Println("in Handlers - UploadBlock - found = true")
			blockJson := block.EncodeToJSON(&uBlock)
			_, err = fmt.Fprint(w, blockJson)
			if err != nil {
				log.Println("Err : in Handlers - UploadBlock - during writing response")
			}
		}
	}

}

// HeartBeatReceive func - Received a heartbeat in request body
func HeartBeatReceive(w http.ResponseWriter, r *http.Request) {
	// Read body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Err : in HeartBeatReceive - reached err of ioutil.ReadAll -")
		log.Println(err)
	} else {
		defer r.Body.Close()

		//heartBeat := data.DecodeToHeartBeatData(string(body)) // heartBeat struct

		go processHeartBeat(data.DecodeToHeartBeatData(string(body))) // process for the receieved heartbeat

		go forwardHeartBeat(data.DecodeToHeartBeatData(string(body))) // forward the heartBeat // here sunnit
	}

}

// processHeartBeat func updates the peerlist, and IfNewBlock then insert the block in SBC
func processHeartBeat(heartBeat data.HeartBeatData) {

	//use hearBeatData to update peer list and get block if the ifNew is set to true
	updatePeerList(&heartBeat)

	if heartBeat.IfNewBlock { //add block in blockchain

		newBlock := block.DecodeFromJSON(heartBeat.BlockJson)
		mptHash := p1.MerklePatriciaTrie(newBlock.Value).Root

		//y := sha3.Sum256([]byte(newBlock.Header.ParentHash + newBlock.Header.Nonce + mptHash))
		//y1 := hex.EncodeToString(y[:])
		//log.Println("++++++++++++++++++++ receieving MPT ROOT at height : ", mptHash)
		//log.Println("++++++++++++++++++++ receving proof : ", y1)

		if p4.POW(newBlock.Header.ParentHash, newBlock.Header.Nonce, mptHash, Difficulty) {
			//fmt.Println(":::: PROCESSING HeartBeat : in ProcessHeartBeat : newBlock.Value.Root : ", mptHash)

			//apr4
			//hold parent / grandparent / etc blocks to be put once we find the begining block based on nodes local copy
			//var blockHolder []block.Block
			//// apr4

			if SBC.CheckParentHash(newBlock) {

				SBC.Insert(newBlock) // if parentHash exist then directly insert and POW is satisfied

			} else if AskForBlock(newBlock.Header.Height-1, newBlock.Header.ParentHash, make([]block.Block, 0) /*, SBC.GetLength(), newBlock.Header.Height-1*/) {
				//if parent cannot be found then ask for parent block and insert both
				SBC.Insert(newBlock)
				//AskForBlock(newBlock.Header.Height, newBlock.Header.ParentHash, make([]block.Block, 0), SBC.GetLength()+1, newBlock.Header.Height+1)

			}
		}

		//fmt.Println("NOT processing HeartBeat : in ProcessHeartBeat : newBlock.Value.Root : ", newBlock.Value.Root)

	}
}

// ForwardHeartBeat func forwards the receieved heartbeat to all its peers
func forwardHeartBeat(heartBeatData data.HeartBeatData) {

	Peers.Rebalance()
	peerMap := Peers.Copy()
	hopCount := heartBeatData.Hops //to forward heartbeat
	if hopCount > 0 {
		heartBeatData.Hops--
		heartBeatData.Id = Peers.GetSelfId()
		heartBeatData.Addr = SELF_ADDR

		//json, _ := json.Marshal(peerMap)
		//heartBeatData.PeerMapJson = string(json)

		//list over peers and send them heartBeat
		if len(peerMap) > 0 {
			for peer := range peerMap {
				_, _ = http.Post(peer+"/heartbeat/receive", "application/json; charset=UTF-8",
					strings.NewReader(heartBeatData.EncodeToJson()))
			}
		}
	}

}

// updatePeerList func updates the existing peerlist with data from received peerMap
func updatePeerList(heartBeat *data.HeartBeatData) {
	Peers.Add(heartBeat.Addr, heartBeat.Id)
	Peers.InjectPeerMapJson(heartBeat.PeerMapJson, SELF_ADDR)
}

// AskForBlock - Ask another server to return a block of certain height and hash
func AskForBlock(height int32, hash string, blockHolder []block.Block) bool {

	//found := false
	Peers.Rebalance()
	peerMap := Peers.Copy()
	//var peersToRemove []string

	//list over peers and send them heartBeat
	//if len(peerMap) > 0 {
	for peer := range peerMap {
		//fmt.Println("\n\nin AskForBlock : req URL : ", peer+"/block/"+strconv.Itoa(int(height))+"/"+hash)
		resp, err := http.Get(peer + "/block/" + strconv.Itoa(int(height)) + "/" + hash)
		if err != nil {
			log.Println("Askblock Err 1 : ", err)
			log.Println("in AskForBlock - deleting peer : ", peer)
			Peers.Delete(peer)
			continue

		} else {
			defer resp.Body.Close() //moved from above err check to here

			body, err := ioutil.ReadAll(resp.Body) //blockJson
			if err != nil {
				log.Println("Askblock Err 2 : ", err)
				continue
			}

			reqBlock := block.DecodeFromJSON(string(body))
			//fmt.Println("\n in AskForBlock : reqBlock", reqBlock, "\n")

			if SBC.CheckParentHash(reqBlock) {
				SBC.Insert(reqBlock)                         // this block
				for i := len(blockHolder) - 1; i >= 0; i-- { // and rest of previous block
					SBC.Insert(blockHolder[i])
				}
				return true
			}

			//if !SBC.CheckParentHash(reqBlock) {  // if parenthash not in local blockchain
			fmt.Println("AskBlock - cannot find parent block for block in height : ", height)
			blockHolder = append(blockHolder, reqBlock) //apr4
			if height <= 1 {
				//continue //apr5
				break
			}
			AskForBlock(height-1, reqBlock.Header.ParentHash, blockHolder) // ask for parents parent block

		} // parsing responsed block
	} // looping peerlist

	return false // if parent block not found

}

//StartHeartBeat func periodically sends heartbeatdata to peers
func StartHeartBeat() {

	for true {
		Peers.Rebalance()
		peerMap := Peers.Copy()
		//PeerMapJson, _ := Peers.PeerMapToJson() //
		PeerMapJson, _ := data.PeerMapToJson(peerMap) //apr4

		//selfAddr := "http://localhost:" + os.Args[1] // SELF_ADDR
		heartBeat := data.PrepareHeartBeatData(&SBC, Peers.GetSelfId(), PeerMapJson, SELF_ADDR, false, "{}")

		//list over peers and send them heartBeat
		if len(peerMap) > 0 {
			for peer := range peerMap {
				_, err := http.Post(peer+"/heartbeat/receive", "application/json; charset=UTF-8",
					strings.NewReader(heartBeat.EncodeToJson())) //apr4
				if err != nil {
					Peers.Delete(peer)
					//fmt.Println("deleting peer : ", peer)
				}
			}
		}
		time.Sleep(10 * time.Second)
	}
}

////          pow ///

//StartTryingNonces func sends heartbeatdata with new block information to peers
func StartTryingNonces() {

	tryingNonces( /*parentHash ,&mpt, */ Difficulty)

}

// tryingNonces func tries to create a new block
func tryingNonces( /*parentHash string, mpt *p1.MerklePatriciaTrie, */ difficulty int) {

	// y = SHA3(parentHash + nonce + mptRootHash)

	var parentBlock block.Block
	var parentHash string
	var mpt p1.MerklePatriciaTrie

	GetNewParent = true
	var nonce string

	for {

		if GetNewParent == true {
			parentBlock = SBC.GetLatestBlocks()[0] //[rand.Int()%len(SBC.GetLatestBlocks())]//random parent from blocks at latest height
			parentHash = parentBlock.Header.Hash
			tryingForHeight = parentBlock.Header.Height + 1
			fmt.Println("in tryingNonces : parentHash : ", parentHash)
			mpt = p1.GenerateRandomMPT()
			nonce = p4.InitializeNonce(8)

			GetNewParent = false
		}

		if p4.POW(parentHash, nonce, mpt.Root, difficulty) {
			//generate block send heartbeat (blockBeat)
			SendBlockBeat(tryingForHeight, parentHash, nonce, mpt)

			GetNewParent = true
		}
		nonce = p4.InitializeNonce(8) //NextNonce(Nonce)

	}

}

// SendBlockBeat func prepares heartbeat data and sends across to peers
func SendBlockBeat(height int32, parentHash string, nonce string, mpt p1.MerklePatriciaTrie) {

	//log.Println("-------------------- sending MPT ROOT at height : ", mpt.Root)
	//y := sha3.Sum256([]byte(parentHash + nonce + mpt.Root))
	//y1 := hex.EncodeToString(y[:])
	//log.Println("-------------------- sending proof : ", y1)

	Peers.Rebalance()
	peerMap := Peers.Copy()
	PeerMapJson, _ := Peers.PeerMapToJson()

	b1 := SBC.GenBlock(height, parentHash, mpt, nonce)
	SBC.Insert(b1)
	blockJson := block.EncodeToJSON(&b1)

	heartBeat := data.PrepareHeartBeatData(&SBC, Peers.GetSelfId(), PeerMapJson, SELF_ADDR, true, blockJson)

	//list over peers and send them heartBeat
	if len(peerMap) > 0 {
		for peer := range peerMap {
			_, err := http.Post(peer+"/heartbeat/receive", "application/json; charset=UTF-8",
				strings.NewReader(heartBeat.EncodeToJson())) //apr4
			if err != nil {
				Peers.Delete(peer)
				//fmt.Println("deleting peer : ", peer)
			}
		}
	}
}

//for return code 500
func returnCode500(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Server Error", http.StatusInternalServerError)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//for return code 204
func returnCode204(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Block does not exists", http.StatusNoContent)
	http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
}
