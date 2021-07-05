package utils

import (
	"bytes"
	"encoding/gob"
	"log"
)

func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToBytes(i interface{}) []byte {
	var aBlockBuffer bytes.Buffer
	encoder := gob.NewEncoder(&aBlockBuffer)
	HandleErr(encoder.Encode(i))
	return aBlockBuffer.Bytes()
}

func FromBytes(i interface{}, data []byte) {
	decoder := gob.NewDecoder(bytes.NewReader(data))
	decoder.Decode(i)
}
