package utils

import (
	"bytes"
	"encoding/gob"
)

func Convert_int32_to_byte(in int32) []byte {
	var buf bytes.Buffer
	encorder := gob.NewEncoder(&buf)
	err := encorder.Encode(in)
	if err != nil {
		return nil
	}
	return buf.Bytes()
}
