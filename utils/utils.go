// Package utils contains functions to be used across the application.
package utils

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

var logFn = log.Panic

// HandleErr takes error and log it.
func HandleErr(err error) {
	if err != nil {
		logFn(err)
	}
}

//ToBytes takes interface and return encode the byte.
func ToBytes(i interface{}) []byte {
	var aBlockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBlockBuffer)
	HandleErr(encoder.Encode(i))
	return aBlockBuffer.Bytes()
}

//FromBytes takes an interface and data and then will decode the data to the interface.
func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(i)
}

//Hash takes an interface, hashes it and returns the hex encoding of the hash.
func Hash(i interface{}) string {
	s := fmt.Sprintf("%v", i)
	hash := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", hash)
}

//Splitter takes an string ,seperator,index. And it splits by seperator and return what you want it.
func Splitter(s, sep string, i int) string {
	r := strings.Split(s, sep)
	if len(r)-1 < i {
		return ""
	}
	return r[i]

}

//ToJSON takes an interface and it return JSON encoding.
func ToJSON(i interface{}) []byte {
	r, err := json.Marshal(i)
	HandleErr(err)
	return r
}
