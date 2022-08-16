package main

import (
	"bytes"
	"crypto/rand"
	"golang.org/x/crypto/nacl/sign"
	"os"
)

// GenKeys generate key pair and save it on disk
func GenKeys(privateFile string, publicFile string) error {
	priv, pub, err := sign.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	// private key
	var privateBuffer bytes.Buffer
	pr, err := os.OpenFile(privateFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if _, err := privateBuffer.Write(priv[:]); err != nil {
		return err
	}
	_, err = privateBuffer.WriteTo(pr)
	if err != nil {
		return err
	}
	//

	// public key
	var publicBuffer bytes.Buffer
	pu, err := os.OpenFile(publicFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	if _, err := publicBuffer.Write(pub[:]); err != nil {
		return err
	}
	_, err = publicBuffer.WriteTo(pu)
	if err != nil {
		return err
	}
	//

	return nil
}

// Sign take a msg as arg and sign it with the private key
func Sign(privateKey *[64]byte, data []byte) []byte {
	signedMsg := sign.Sign(nil, data, privateKey)
	return signedMsg
}

// Verify verify a signature
func Verify(publicKey *[32]byte, signature []byte) ([]byte, bool) {
	return sign.Open(nil, signature, publicKey)
}
