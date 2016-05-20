package gokeyserve

import (
	"crypto/rsa"
	"net/rpc"
)

type Request struct {
}

type Response struct {
	Key *rsa.PrivateKey
}

type Client struct {
	rpcClient *rpc.Client
}

func NewClient(serverAddress string) (client *Client, err error) {
	r, err := rpc.Dial("tcp", serverAddress)
	if err != nil {
		return nil, err
	}
	client = &Client{rpcClient: r}
	return client, nil
}

func (c *Client) NewKey() (key *rsa.PrivateKey, err error) {
	resp := &Response{}
	err = c.rpcClient.Call("GoKeyServer.Generate", Request{}, resp)
	return resp.Key, err
}

func (c *Client) Close() error {
	if c.rpcClient != nil {
		return c.rpcClient.Close()
	}
	return nil
}




