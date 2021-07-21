package wallet

import (
	"crypto/x509"
	"encoding/hex"
	"io/fs"
	"reflect"
	"testing"
)

const (
	testKey     string = "3077020101042084978155e02faee35816ddb5c498ecc135afd0b2228f3bad00d698add87c9ed9a00a06082a8648ce3d030107a1440342000436e6cd5de38c5ebf39a54256ae143a344fbd853634dac3cf02e88d0a6a84a4182d4769f88add0aaf276ae867ceac88a26779ca294ad29909817223ebc245c064"
	testPayload string = "000609a1d5586ca6ff02d21077dc1e9033959ddd1ba53288721b1b3adc98e6e3"
	testSig     string = "a69c03ac6a226d3d7989e8cb61ff2b6945ca67a35b728b41a9abcca89340424990a08b241720fc9b0e59eb2a31d3eaa8f854edf044d9ce9f08b456d1aa2e9749"
)

type fakeLayer struct {
	fakeHasWalletFile func() bool
}

func (f fakeLayer) hasWalletFile() bool {
	return f.fakeHasWalletFile()
}

func (fakeLayer) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return nil
}

func (fakeLayer) ReadFile(name string) ([]byte, error) {
	return x509.MarshalECPrivateKey(makeTestWallet().privateKey)
}

func TestWallet(t *testing.T) {
	t.Run("New wallet is created.", func(t *testing.T) {
		files = fakeLayer{
			fakeHasWalletFile: func() bool {
				t.Log("I have been called.")
				return false
			},
		}
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("New wallet should return a new wallet instance.")
		}
	})
	t.Run("Wallet is restored.", func(t *testing.T) {
		w = nil
		files = fakeLayer{
			fakeHasWalletFile: func() bool {
				t.Log("I have been called.")
				return true
			},
		}
		tw := Wallet()
		if reflect.TypeOf(tw) != reflect.TypeOf(&wallet{}) {
			t.Error("Exist wallet should restore from file.")
		}
	})

}

func makeTestWallet() *wallet {
	w := &wallet{}
	b, _ := hex.DecodeString(testKey)
	key, _ := x509.ParseECPrivateKey(b)
	w.privateKey = key
	w.Address = aFromK(key)
	return w

}

func TestSign(t *testing.T) {
	s := Sign(testPayload, makeTestWallet())
	_, err := hex.DecodeString(s)
	if err != nil {
		t.Errorf("Sign() should return a hex encoded string, got %s", s)
	}
}

func TestVerify(t *testing.T) {
	type test struct {
		input string
		ok    bool
	}
	tests := []test{
		{input: testPayload, ok: true},
		{input: "000409a1d5586ca6ff02d21077dc1e9033959ddd1ba53288721b1b3adc98e6e3", ok: false},
	}
	for _, tc := range tests {
		w := makeTestWallet()
		ok := Verify(testSig, tc.input, w.Address)
		if ok != tc.ok {
			t.Error("Verify() could not verify testSignature and testPayload.")
		}

	}
}

func TestRestoreBigInts(t *testing.T) {
	_, _, err := restoreBigInts("xx")
	if err == nil {
		t.Error("restoreBigInts() should return error when payload is not hex.")
	}
}
