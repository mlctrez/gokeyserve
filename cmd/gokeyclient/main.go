package main

import (
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"github.com/mlctrez/gokeyserve"
	"log"
)

func main() {

	serverAddress := flag.String("server", "127.0.0.1:1234", "the address of the server")
	flag.Parse()

	c, err := gokeyserve.NewClient(*serverAddress)
	defer c.Close()

	if err != nil {
		log.Fatal(err)
	}

	k, err := c.NewKey()
	if err != nil {
		log.Fatal(err)
	}

	// prints the private key in pem format
	block := pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	fmt.Print(string(pem.EncodeToMemory(&block)))

}
