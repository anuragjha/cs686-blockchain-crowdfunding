# cs686_BlockChain_P5 - Lending system


# (1) Crypto 

## security.go
Implements the public-private keys and Hash generation. Also creating signature and verifying signature.
Contains two data structs - Identity and PublicIdentity.
GetPublicIdentity in security.go - for Node to be able to get its PublicId in new variable.

Using golang rsa PKCS1v15 to sign and verify and encrypt and decrypt.

1. new data structs in security.go
  - Identity {
      privateKey,
      PublicKey,
      Label,
      }
  
  - PublicIdentity {
      PublicKey,
      Label,
      }

2. new data struct in PeerList
  secureId - added to contain secureId of dataType Identity.
  peerMapPid - now also contains PeerMapPid
  PeerMapPid is a Map of - Addr of Node (key) - and PublicIdentity (value) of peers.

## heartbeat.go
 HeartBeat changed to now also include
   Pid
   SignForBlockJson - for BlockJson
   PeerMapPidJson  
   ## Todo - sender have to encrypt the heartbeat with public key of receving peers, Receiever have to decrypt the heartbeat with private key of itself.

## handlers.go
Add functionality in creating heartBeat
add signature for blockjson
add Pid of Sender


## API
GET /uploadpids - to be used first time along with download blockchain.

HeartBeat send and Receieve now additionally deals with Signature of sender and its verification by receiever.

>>>>>>>>>>>>
Have to add functionality to send encrypted blockjson and decrypt blockjson - in handler.go (funcs available in security.go)


# (2) Currency

1. Data structs in wallet.go
  - Wallet {
    Balance,
  }
  
2. Data structs in transaction.go
  - Transaction {
    Id
    From
    To
    Tokens
    Timestamp
  }

  - TransactionPool {
     list of transaction
  }

  - TransactionBeat {
      Transaction
      FromPid
      TxSig
  }

## transaction.go
CreateTransaction func -> takes params From public Id, To public Id, Tokens and Timestamp and -> returns Tx.
NewTransactionBeat func -> takes params Tx, From public Id and FromSig and -> returns TransactionBeat.
CreateTransactionBeat func ->  takes params Tx and Identity and -> returns 
AddToTransactionPool func -> 
DeleteFromTransactionPool func ->
ReadFromTransactionPool func ->













