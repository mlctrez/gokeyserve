package server

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"log"
	"time"
)

var keychan chan *rsa.PrivateKey
var rotationInterval time.Duration
var generatorCount int

// Request is the first RPC argument. It contains no data.
type Request struct {
}

// Response is the return from the RPC call containing the private key.
type Response struct {
	Key *rsa.PrivateKey
}

// GoKeyServer is the rpc type.
type GoKeyServer int

// Generate is the RPC method to generate a private key.
func (t *GoKeyServer) Generate(req *Request, res *Response) error {
	res.Key = GetGeneratedKey()
	return nil
}

// GenerateKey creates a rsa.PrivateKey of bits length or returns nil on error.
func GenerateKey(bits int) (privatekey *rsa.PrivateKey) {
	privatekey, _ = rsa.GenerateKey(rand.Reader, bits)
	return privatekey
}

// GetGeneratedKey retrieves a pre generated key from the channel.
func GetGeneratedKey() *rsa.PrivateKey {
	return <-keychan
}

// Start is used to initialize the channel generator and key rotation
// seconds determines how often in memory keys are rotated or zero for no rotation.
func Start(rotateEvery time.Duration, generators int) error {

	if generators < 1 {
		return errors.New("generators must be 1 or greater")
	}

	rotationInterval = rotateEvery
	generatorCount = generators
	keychan = make(chan *rsa.PrivateKey, generatorCount-1)

	log.Printf("generators: %v", generatorCount)
	for i := 0; i < generatorCount; i++ {
		go generate()
	}

	if rotationInterval > 0*time.Second {
		log.Printf("rotating every: %v", rotationInterval)
		go rotate()
	}

	return nil
}

func generate() {
	log.Println("new generator")
	for {
		log.Println("generating key")
		keychan <- GenerateKey(2048)
	}
}

// ensures that the any keys held in memory are rotated every rotationSeconds
func rotate() {
	for {
		time.Sleep(rotationInterval)
		for i := 0; i < generatorCount; i++ {
			log.Println("rotating key")
			GetGeneratedKey()
		}
	}
}
