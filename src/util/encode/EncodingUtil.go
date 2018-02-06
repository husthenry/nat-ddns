package encode

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
)

func Base64Encode(byts []byte) []byte {
	res := make([]byte, base64.StdEncoding.EncodedLen(len(byts)))
	base64.StdEncoding.Encode(res, byts)

	return res
}

func Base64Decode(byts []byte) []byte {
	res := make([]byte, base64.StdEncoding.DecodedLen(len(byts)))
	base64.StdEncoding.Decode(res, byts)

	return res
}

func HexEncode(byts []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(byts)))
	hex.Encode(dst, byts)

	return dst
}

func HexDecode(byts []byte) []byte {
	dst := make([]byte, hex.DecodedLen(len(byts)))
	n, err := hex.Decode(dst, byts)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	fmt.Printf("%s\n", dst[:n])
	return dst
}
