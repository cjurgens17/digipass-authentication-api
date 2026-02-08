package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"time"
)

const SecondsInDay uint32 = 86400

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type PublicClaims struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Exp int64  `json:"exp"`
	Iat int64  `json:"iat"`
	Jti string `json:"jti"`
}

type PrivateClaims struct {
	Data any `json:"data"`
}

type Payload struct {
	Public  PublicClaims  `json:"public"`
	Private PrivateClaims `json:"private"`
}

type ClientInfo struct {
	Tenant string
	Url    string
}

func base64URLEncode(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func GenerateHeader() string {
	alg := os.Getenv("JWT_ALGO")
	if alg == "" {
		alg = "HS256"
	}

	header := Header{
		Alg: alg,
		Typ: "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		log.Fatal("Error Attempting to serialize JSON")
	}

	return base64URLEncode(headerJSON)
}

func DummyDBService(clientId string) ClientInfo {
	return ClientInfo{
		Tenant: "SomberLake",
		Url:    "https://railway.digipass.authapi.dev.com",
	}
}

func GenerateUnixExpiration(value ...uint32) int64 {
	exp := SecondsInDay

	if len(value) > 1 {
		log.Fatal("Value Can only accept one argument")
	}

	if len(value) > 0 {
		exp = value[0]
	}

	return time.Now().Add(time.Duration(exp) * time.Second).Unix()
}

func GetCurrentUnixTimestamp() int64 {
	return time.Now().Unix()
}

func DummyGenerateJTI() string {
	return "jti:abc:001"
}

func GeneratePayload(clientId string, userId string) string {
	clientInfo := DummyDBService(clientId)
	exp := GenerateUnixExpiration(86400)
	iat := GetCurrentUnixTimestamp()
	jti := DummyGenerateJTI()
	payload := Payload{
		Public: PublicClaims{
			Iss: clientInfo.Tenant,
			Sub: userId,
			Aud: clientInfo.Url,
			Exp: exp,
			Iat: iat,
			Jti: jti,
		},
		Private: PrivateClaims{
			Data: nil,
		},
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		log.Fatal("Error Trying to serialize payload to JSON")
	}

	return base64URLEncode(payloadJSON)
}

func GenerateSignature(encodedHeader string, encodedPayload string) string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("Secret not available, exiting to prevent vulnerability")
	}

	signatureBase := encodedHeader + "." + encodedPayload

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(signatureBase))
	signature := h.Sum(nil)

	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	return encodedSignature
}

func GenerateJWT(clientId string, userId string) string {
	encodedHeader := GenerateHeader()
	encodedPayload := GeneratePayload(clientId, userId)
	signature := GenerateSignature(encodedHeader, encodedPayload)

	dot := "."

	token := encodedHeader + dot + encodedPayload + dot + signature
	return token
}
