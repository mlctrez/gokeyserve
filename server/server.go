package server

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"time"
)

// GoKeyServer is the server struct
type GoKeyServer struct {
	keyChan          chan *rsa.PrivateKey
	rotationInterval time.Duration
	generatorCount   int
}

// GenerateKey creates a rsa.PrivateKey of bits length or returns nil on error.
func GenerateKey(bits int) (privatekey *rsa.PrivateKey) {
	privatekey, _ = rsa.GenerateKey(rand.Reader, bits)
	return privatekey
}

// GetGeneratedKey retrieves a pre generated key from the channel.
func (s *GoKeyServer) GetGeneratedKey() *rsa.PrivateKey {
	return <-s.keyChan
}

// New is used to initialize the channel generator and key rotation
// seconds determines how often in memory keys are rotated or zero for no rotation.
func New(rotationInterval time.Duration, generatorCount int) (s *GoKeyServer, err error) {

	if generatorCount < 1 {
		return nil, errors.New("generatorCount must be 1 or greater")
	}

	s = &GoKeyServer{
		rotationInterval: rotationInterval,
		generatorCount:   generatorCount,
		keyChan:          make(chan *rsa.PrivateKey, generatorCount),
	}

	for i := 0; i < generatorCount; i++ {
		go s.generate()
	}

	if rotationInterval > 0 {
		go s.rotate()
	}

	return s, nil
}

func (s *GoKeyServer) generate() {
	for {
		s.keyChan <- GenerateKey(2048)
	}
}

// ensures that the any keys held in memory are rotated every rotationSeconds
func (s *GoKeyServer) rotate() {
	for {
		time.Sleep(s.rotationInterval)
		for i := 0; i < s.generatorCount; i++ {
			_ = s.GetGeneratedKey()
		}
	}
}
