# cs686_BlockChain_P5 - Lending system


# Registration
1. new data structs in security.go (*1*)
  - Identity {
      privateKey,
      PublicKey,
      HashForKey,
      Label,
      }
  
  - PublicIdentity {
      PublicKey,
      HashForKey,
      Label,
      }

2. new data struct in PeerList
  - added to contain secureId of dataType Identity.
  - now also contains PeerMapPid
  PeerMapPid is a Map of Addr (key) and PublicIdentity (value) of peers.

3. HeartBeat changed to now also include  a)sender Pid   b)PeerMapPid

4. New methods
GetPublicIdentity in security.go - for Node to be able to get its PublicId in new variable



(*1*)
# security.go
Implements the public-private keys and Hash generation. Also creating signature and verifying signature.
Contains two data structs - Identity and PublicIdentity (should have json converion and back only for PublicIdentity)



# API
GET /uploadpids - to be used first time along with download blockchain.

HeartBeat send and Receieve now additionally deals with Signature of sender and its verification by receiever.
