package appinfralayer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

const TokenUseAccess string = "access"
const TokenUseID string = "id"
const JWTPartsDelim string = "."
const JWTPartsLength int = 3

type CognitoIDTokenJWTHeader struct {
	KeyID     string `json:"kid"`
	Algorithm string `json:"alg"`
}

func NewCognitoIDTokenJWTHeaderByStr(s string) (CognitoIDTokenJWTHeader, error) {
	value := CognitoIDTokenJWTHeader{}
	headerDec, err := DecodeBase64URL(s)
	if err != nil {
		return value, err
	}
	err = json.Unmarshal(headerDec, &value)
	if err != nil {
		return value, err
	}
	return value, nil
}

type CognitoIDTokenJWTPayload struct {
	Sub                 string `json:"sub"`
	Aud                 string `json:"aud"`
	EmailVerified       bool   `json:"email_verified"`
	TokenUse            string `json:"token_use"`
	AuthTime            int64  `json:"auth_time"`
	Iss                 string `json:"iss"`
	AuthServiceUserName string `json:"cognito:username"`
	Exp                 int64  `json:"exp"`
	GivenName           string `json:"given_name"`
	Iat                 int64  `json:"iat"`
	Email               string `json:"email"`
}

func (value CognitoIDTokenJWTPayload) IsValidExp(ts int64) bool {
	if value.Exp == 0 {
		return false
	}
	return ts <= value.Exp
}

func (value CognitoIDTokenJWTPayload) IsValidIss(region, poolID string) bool {
	return value.Iss == fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, poolID)
}

func (value CognitoIDTokenJWTPayload) IsIDToken() bool {
	return value.TokenUse == TokenUseAccess
}

func NewCognitoIDTokenJWTPayloadByStr(s string) (CognitoIDTokenJWTPayload, error) {
	value := CognitoIDTokenJWTPayload{}
	headerDec, err := DecodeBase64URL(s)
	if err != nil {
		return value, err
	}
	err = json.Unmarshal(headerDec, &value)
	if err != nil {
		return value, err
	}
	return value, nil
}

type CognitoIDTokenJWT string

func (value CognitoIDTokenJWT) GetHeaderAndPayload() (CognitoIDTokenJWTHeader, CognitoIDTokenJWTPayload, error) {
	jwtStr := string(value)
	jwtSlice := strings.Split(jwtStr, JWTPartsDelim)
	if len(jwtSlice) != JWTPartsLength {
		return CognitoIDTokenJWTHeader{}, CognitoIDTokenJWTPayload{}, fmt.Errorf(`length of jwt parts is %d, must be 3.`, len(jwtSlice))
	}
	header, err := NewCognitoIDTokenJWTHeaderByStr(jwtSlice[0])
	if err != nil {
		return header, CognitoIDTokenJWTPayload{}, err
	}

	payload, err := NewCognitoIDTokenJWTPayloadByStr(jwtSlice[1])
	if err != nil {
		return header, payload, err
	}
	return header, payload, nil
}

func DecodeBase64URL(data string) ([]byte, error) {
	data = strings.Replace(data, "-", "+", -1) // 62nd char of encoding
	data = strings.Replace(data, "_", "/", -1) // 63rd char of encoding

	switch len(data) % 4 { // Pad with trailing '='s
	case 0: // no padding
	case 2:
		data += "==" // 2 pad chars
	case 3:
		data += "=" // 1 pad char
	}

	return base64.StdEncoding.DecodeString(data)
}

type JSONWebKeys struct {
	Keys []JSONWebKey `json:"keys"`
}

func (model JSONWebKeys) Have(target string) bool {
	for _, key := range model.Keys {
		if key.Kid == target {
			return true
		}
	}
	return false
}

type JSONWebKey struct {
	Alg string `json:"alg"`
	E   string `json:"e"`
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	N   string `json:"n"`
	Use string `json:"use"`
}
