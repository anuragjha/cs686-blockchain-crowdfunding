package resource

import (
	"../p5"
	"encoding/json"
	"log"
)

type UserLandingPage struct {
	Pid              p5.PublicIdentity
	FromPid          string
	BTxs             p5.BorrowingTransactions // key - BorrowingTxId
	PromisedInString string
	BB               p5.BalanceBook
	Purse            p5.Wallet
}

type LoginPageStruct struct {
	Phrase  string
	CidJson string
}

func NewLoginPageStruct(phrase string, cidJson string) LoginPageStruct {
	lps := LoginPageStruct{}
	lps.Phrase = phrase
	lps.CidJson = cidJson
	return lps
}

func JsonToLoginPageStruct(str string) LoginPageStruct {
	lps := LoginPageStruct{}
	err := json.Unmarshal([]byte(str), lps)
	if err != nil {
		log.Println("Error in unmarshalling LoginPageStruct - err : ", err)
	}
	return lps
}

func (lps *LoginPageStruct) LoginPageStructToJson() []byte {
	jsonBytes, err := json.Marshal(lps)
	if err != nil {
		log.Println("Error in marshalling LoginPageStruct - err : ", err)
	}

	return jsonBytes
}
