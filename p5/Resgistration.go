package p5

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"

	//"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/sha3"
	"log"
)

func Hello() {

	//message := "Hello World !!!"
	//fmt.Println(message)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("Couldnt generate keys - err :", err)
	}
	fmt.Println("Private Key is : ", priv.D)

	out := x509.MarshalPKCS1PrivateKey(priv) // out is a byte array
	//fmt.Println("out : ", out)

	var p = &pem.Block{
		Type:  "RSA Private Key",
		Bytes: out,
	}
	fmt.Println("pem.Block : ", pem.EncodeToMemory(p))

	//publicKey := &privateKey.PublicKey
	////fmt.Println("Public Key is : ", publicKey.N)
	//
	//hash1 := sha256.New()
	//encrypted, err := rsa.EncryptOAEP(hash1, rand.Reader, publicKey, []byte(message), []byte("what?"))
	//if err != nil {
	//	fmt.Println("Error in Encryption")
	//}
	////fmt.Println("encrypted : ",string(encrypted))
	//
	//
	//decrypted, err := rsa.DecryptOAEP(hash1, rand.Reader, privateKey, []byte(encrypted), []byte("what?"))
	//fmt.Println("Decrypted : ", string(decrypted))

}

func generateHash(mbytes []byte) string {
	sum := sha3.Sum256(mbytes)
	return hex.EncodeToString(sum[:])
}
