package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"

	"github.com/miranaky/kaengkaengcoin/utils"
)

const (
	walletFile string = "kaengkaengCoin.wallet"
)

type wallet struct {
	privateKey *ecdsa.PrivateKey
	Address    string
}

var w *wallet

func hasWalletFile() bool {
	_, err := os.Stat(walletFile)
	return !os.IsNotExist(err)
}

func createPrivKey() *ecdsa.PrivateKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	utils.HandleErr(err)
	return privKey
}

func persistKey(key *ecdsa.PrivateKey) {
	bytes, err := x509.MarshalECPrivateKey(key)
	utils.HandleErr(err)
	err = os.WriteFile(walletFile, bytes, 0644)
	utils.HandleErr(err)

}

func restoreKey() (key *ecdsa.PrivateKey) {
	keyAsByte, err := os.ReadFile(walletFile)
	utils.HandleErr(err)
	key, err = x509.ParseECPrivateKey(keyAsByte)
	utils.HandleErr(err)
	return
}

func encodeBigInt(a, b []byte) string {
	z := append(a, b...)
	return fmt.Sprintf("%x", z)
}

func aFromK(key *ecdsa.PrivateKey) string {
	return encodeBigInt(key.X.Bytes(), key.Y.Bytes())
}

func sign(payload string, w *wallet) string {
	payloadAsB, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	r, s, err := ecdsa.Sign(rand.Reader, w.privateKey, payloadAsB)
	utils.HandleErr(err)
	return encodeBigInt(r.Bytes(), s.Bytes())
}

func restoreBigInt(payload string) (*big.Int, *big.Int, error) {
	bytes, err := hex.DecodeString(payload)
	if err != nil {
		return nil, nil, err
	}
	firstBytes := bytes[:len(bytes)/2]
	secondBytes := bytes[len(bytes)/2:]
	bigA, bigB := big.Int{}, big.Int{}
	bigA.SetBytes(firstBytes)
	bigB.SetBytes(secondBytes)
	return &bigA, &bigB, nil
}

func generatePublicKey(address string) *ecdsa.PublicKey {
	x, y, err := restoreBigInt(address)
	utils.HandleErr(err)
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	return &publicKey
}

func varify(signature, payload, address string) bool {
	s, r, err := restoreBigInt(signature)
	utils.HandleErr(err)
	publicKey := generatePublicKey(address)
	payloadBytes, err := hex.DecodeString(payload)
	utils.HandleErr(err)
	ok := ecdsa.Verify(publicKey, payloadBytes, r, s)
	return ok
}
func Wallet() *wallet {
	if w == nil {
		w = &wallet{}
		if hasWalletFile() {
			w.privateKey = restoreKey()
		} else {
			key := createPrivKey()
			persistKey(key)
			w.privateKey = key
		}
		w.Address = aFromK(w.privateKey)
	}
	return w
}
