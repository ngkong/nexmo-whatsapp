# nexmo-whatsapp

This is my first GOLANG test project. I built this to test the functionality of nexmo whatsapp services.

How to install:
```
git clone https://github.com/ngkong/nexmo-whatsapp
cd nexmo-whatsapp
go install github.com/gbrlsnchs/jwt
go install github.com/google/uuid
go install golang.org/x/crypto/ed25519
```

Setup:
1. copy nexmo private key to ./key/private.key
2. copy conf.json.example to conf.json, adjust the configuration inside

How to run:
```
go run main.go
```

Send message:
```
curl "http://localhost:8008/send?message={message}&to={number}"
```

Setup your nexmo callback URL to:
```
status: http://yourdomain:8008/status
inbound: http://yourdomain:8008/inbound
```
