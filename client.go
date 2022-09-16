package goclient

import (
	"crypto/sha1"
	"encoding/binary"
	"fmt"
	"time"

	sdk "github.com/lokidb/server/client"
)

const vBucketCount int = 1024
const hashUseLenght int = 4

type NodeAddress struct {
	Host string
	Port int
}

type node struct {
	NodeAddress
	stub *sdk.Client
}

type client struct {
	nodes    map[NodeAddress]node
	vBuckets [vBucketCount]NodeAddress
}

func (na *NodeAddress) fullAddress() string {
	return fmt.Sprintf("%s:%d", na.Host, na.Port)
}

func (na *NodeAddress) hash() int {
	h := sha1.New()
	h.Write([]byte(na.fullAddress()))
	byts := h.Sum(nil)[:hashUseLenght]
	return int(binary.BigEndian.Uint32(byts))
}

func New(nodeAddress []NodeAddress, timeout time.Duration) *client {
	c := new(client)

	for i := 0; i < vBucketCount; i++ {
		min := 9223372036854775807 // Max value
		var current_na NodeAddress

		for _, na := range nodeAddress {
			if na.Host == "localhost" {
				na = NodeAddress{Host: "127.0.0.1", Port: na.Port}
			}

			hash := na.hash()
			mod := ((i + 1) * vBucketCount)
			dis := hash % mod

			if dis < min {
				min = dis
				current_na = na
			}
		}

		c.vBuckets[i] = current_na
	}

	c.nodes = make(map[NodeAddress]node, len(nodeAddress))
	for _, na := range nodeAddress {
		if na.Host == "localhost" {
			na.Host = "127.0.0.1"
		}

		stub := sdk.New(fmt.Sprintf("%s:%d", na.Host, na.Port), timeout)
		n := new(node)
		n.NodeAddress = na
		n.stub = stub
		c.nodes[na] = *n
	}

	return c
}

func (c *client) nodeByKey(key string) *sdk.Client {
	h := sha1.New()
	h.Write([]byte(key))
	byts := h.Sum(nil)[:hashUseLenght]
	num := int(binary.BigEndian.Uint32(byts))
	node_id := c.vBuckets[num%vBucketCount]
	node := c.nodes[node_id]
	return node.stub
}

func (c *client) Get(key string) (string, error) {
	node := c.nodeByKey(key)
	return node.Get(key)
}

func (c *client) Set(key string, value string) error {
	node := c.nodeByKey(key)
	return node.Set(key, value)
}

func (c *client) Del(key string) (bool, error) {
	node := c.nodeByKey(key)
	return node.Del(key)
}

func (c *client) Keys() ([]string, error) {
	keys := make([]string, 0, 1000)
	for _, node := range c.nodes {
		kys, err := node.stub.Keys()
		if err != nil {
			return nil, err
		}

		keys = append(keys, kys...)
	}

	return keys, nil
}

func (c *client) Flush() error {
	for _, node := range c.nodes {
		err := node.stub.Flush()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *client) Close() {
	for _, node := range c.nodes {
		node.stub.Close()
	}
}
