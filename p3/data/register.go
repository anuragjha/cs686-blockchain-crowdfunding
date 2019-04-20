package data

import "encoding/json"

//RegisterData struct defines the AssignedId and PeerMapJson to a struct
type RegisterData struct {
	AssignedId  int32  `json:"assignedId"`
	PeerMapJson string `json:"peerMapJson"`
}

//NewRegisterData func creates a new RegsterData
func NewRegisterData(id int32, peerMapJson string) RegisterData {
	return RegisterData{
		AssignedId:  id,
		PeerMapJson: peerMapJson,
	}
}

//EncodeToJson func encodes RegisterData to json string
func (data *RegisterData) EncodeToJson() (string, error) {

	dataEncodedBytearray, err := json.Marshal(data)
	return string(dataEncodedBytearray), err
}
