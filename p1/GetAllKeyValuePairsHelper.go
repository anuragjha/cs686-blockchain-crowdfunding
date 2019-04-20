package p1

import (
	"errors"
)

func (mpt *MerklePatriciaTrie) GetAllKeyValuePairsHelper(mptKeyValuePairs map[string]string, thisNode Node, hexPath []uint8) (map[string]string, error) {
	currentHexPath := hexPath

	switch {
	case thisNode.node_type == 1:
		for i := 0; i < 16; i++ {
			if thisNode.branch_value[i] != "" {
				newcurrentHexPath := append(currentHexPath, uint8(i)) //int should be treated as part of ascii path
				mpt.GetAllKeyValuePairsHelper(mptKeyValuePairs, mpt.db[thisNode.branch_value[i]], newcurrentHexPath)
			}
		}
		if thisNode.branch_value[16] != "" {
			key := HexArraytoString(currentHexPath)
			mptKeyValuePairs[key] = thisNode.branch_value[16]
		}

	case thisNode.node_type == 2 && is_ext_node(thisNode.flag_value.encoded_prefix) == true:
		thisNodePath := compact_decode(thisNode.flag_value.encoded_prefix)
		currentHexPath := append(currentHexPath, thisNodePath...) //int should be treated as part of ascii path
		mpt.GetAllKeyValuePairsHelper(mptKeyValuePairs, mpt.db[thisNode.flag_value.value], currentHexPath)

	case thisNode.node_type == 2 && is_ext_node(thisNode.flag_value.encoded_prefix) == false:
		thisNodePath := compact_decode(thisNode.flag_value.encoded_prefix)
		currentHexPath := append(currentHexPath, thisNodePath...)
		key := HexArraytoString(currentHexPath)
		mptKeyValuePairs[key] = thisNode.flag_value.value
	default:
		return nil, errors.New("Error in contructing key Value map from MPT")

	}

	return mptKeyValuePairs, nil

}
