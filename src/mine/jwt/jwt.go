package jwt

import (
	"time"
	"github.com/gbrlsnchs/jwt"
	"github.com/google/uuid"
	"io/ioutil"
	"encoding/pem"
	"crypto/x509"
	"crypto/rsa"
	"os"
)

var (
	priv_key_location string
	application_id string
)

type CustomPayload struct {
	jwt.Payload
	Application_id string `json:"application_id,omitempty"`
}

func Set_priv_key_location(loc string){
    if _, err := os.Stat(loc); err != nil {
        if os.IsNotExist(err) {
            panic("priv key doesn't exists")
        }
    }
	priv_key_location = loc
}

func Set_app_id(id string){
	application_id = id
}

func Gen() string {

	priv, _ := ioutil.ReadFile(priv_key_location)
	privPem, _ := pem.Decode(priv)
	var privPemBytes []byte
	privPemBytes = privPem.Bytes
	var parsedKey interface{}
	parsedKey, _ = x509.ParsePKCS8PrivateKey(privPemBytes)
	var privateKey *rsa.PrivateKey
	_ = privateKey
	var ok bool
	_ = ok
	privateKey, ok = parsedKey.(*rsa.PrivateKey)
	hs := jwt.RSAPrivateKey(privateKey)

	now := time.Now()
	pl := CustomPayload{
		Payload: jwt.Payload{
			ExpirationTime: jwt.NumericDate(now.Add(24 * 30 * 12 * time.Hour)),
			JWTID:          uuid.New().String(),
			IssuedAt:       jwt.NumericDate(now),
		},
		Application_id: application_id,
	}

	token, err := jwt.Sign(pl, jwt.NewRS256(hs))
	_ = token
	if err != nil {
	}
	
	return string(token)
}