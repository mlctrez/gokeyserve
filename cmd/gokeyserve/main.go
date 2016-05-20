package main

import (
	"flag"
	"github.com/mlctrez/gokeyserve/server"
	"log"
	"net"
	"net/rpc"
	"time"
)

// example of using gokeyserve/server

func main() {

	listenAddress := flag.String("listen", ":1234", "the listen address for the server")
	generatorCount := flag.Int("generators", 1, "number of goroutines for generators")
	expireInterval := flag.String("exinterval", "5m", "how often to expire cached keys in memory, 0s for no keys kept in memory")

	flag.Parse()

	intvl, err := time.ParseDuration(*expireInterval)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Start(intvl, *generatorCount)
	if err != nil {
		log.Fatal(err)
	}

	s := new(server.GoKeyServer)
	rpc.Register(s)

	log.Printf("listening on address %v", *listenAddress)

	addr, err := net.ResolveTCPAddr("tcp", *listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		con, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go rpc.ServeConn(con)
	}

}
