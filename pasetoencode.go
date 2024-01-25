package pakarbibackend

import (
	"encoding/json"

	"fmt"

	"aidanwoods.dev/go-paseto"
)

// <--- FUNCTION PASETO ENCODER --->
func Decoder(publickey, tokenstr string) (payload Payload, err error) {
	var token *paseto.Token
	var pubKey paseto.V4AsymmetricPublicKey
	pubKey, err = paseto.NewV4AsymmetricPublicKeyFromHex(publickey)
	if err != nil {
		fmt.Println("Decode NewV4AsymmetricPublicKeyFromHex : ", err)
	}
	parser := paseto.NewParser()
	token, err = parser.ParseV4Public(pubKey, tokenstr, nil)
	if err != nil {
		fmt.Println("Decode ParseV4Public : ", err)
	} else {
		json.Unmarshal(token.ClaimsJSON(), &payload)
	}
	return payload, err
}

func DecodeGetParkiran(PublicKey, tokenStr string) (pay string, err error) {
	key, err := Decoder(PublicKey, tokenStr)
	if err != nil {
		fmt.Println("Cannot decode the token", err.Error())
	}
	return key.Parkiran, nil
}

func DecodeGetUser(PublicKey, tokenStr string) (user string, err error) {
	key, err := Decoder(PublicKey, tokenStr)
	if err != nil {
		fmt.Println("Cannot decode the token", err.Error())
		return "", err
	}
	return key.Parkiran, nil
}
