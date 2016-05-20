package gokeyserve

import (
	"crypto/rsa"
	"net/rpc"
)

// Request is the first RPC argument. It contains no data.
type Request struct {
}

// Response is the return from the RPC call containing the private key.
type Response struct {
	Key *rsa.PrivateKey
}

// Client contains the rpc client.
type Client struct {
	rpcClient *rpc.Client
}

// NewClient constructs a new gokeyserve rpc client.
func NewClient(serverAddress string) (client *Client, err error) {
	r, err := rpc.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	client = &Client{rpcClient: r}
	return client, nil
}

// NewKey calls the remote server to generate the key.
func (c *Client) NewKey() (key *rsa.PrivateKey, err error) {
	resp := &Response{}
	err = c.rpcClient.Call("GoKeyServer.Generate", Request{}, resp)
	key = resp.Key
	return key, err
}

// Close closes the underlying rpc client.
func (c *Client) Close() error {
	if c.rpcClient != nil {
		return c.rpcClient.Close()
	}
	return nil
}
