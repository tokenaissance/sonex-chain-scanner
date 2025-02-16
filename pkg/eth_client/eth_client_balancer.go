package eth_client

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/sonex-dex/sonex-chain-scanner/pkg/common"
	"github.com/sonex-dex/sonex-chain-scanner/pkg/common/config"
)

type EthClientBalancer struct {
	nodes   []*node
	clients map[string]*EthClient
}

type node struct {
	url          string
	weight       int64
	failureCount int
	lastFailure  time.Time
}

func NewEthClientBalancer(chainConf config.ChainConfig) (*EthClientBalancer, error) {
	var nodes []*node
	for url, weight := range chainConf.Nodes {
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			continue
		}
		_, err := ethclient.Dial(url)
		if err != nil {
			continue
		}
		nodes = append(nodes, &node{url: url, weight: weight})
	}
	if len(nodes) == 0 {
		return nil, fmt.Errorf("%s", common.ERROR_NO_AVAILABLE_NODE)
	}

	return &EthClientBalancer{nodes: nodes}, nil
}

func (ethClientBalancer *EthClientBalancer) getNode() *node {
	// Check availability
	availableNodes := []node{}
	linq.From(ethClientBalancer.nodes).ForEachT(func(item *node) {
		if item.failureCount > 10 && time.Since(item.lastFailure) < time.Minute*5 {
			return
		}
		availableNodes = append(availableNodes, *item)
	})

	// Calculate weight
	totalWeight := linq.From(availableNodes).SelectT(func(item node) int64 { return item.weight }).SumInts()
	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	randValue := rand.Int63n(totalWeight)
	var cumulativeWeight int64
	for _, node := range availableNodes {
		cumulativeWeight += node.weight
		if randValue < cumulativeWeight {
			return &node
		}
	}

	return nil
}

func (ethClientBalancer *EthClientBalancer) increaseFailureCount(url string) {
	nodeA := linq.From(ethClientBalancer.nodes).FirstWithT(func(item *node) bool { return item.url == url })
	if nodeA == nil {
		return
	}
	node := nodeA.(*node)
	node.failureCount++
	node.lastFailure = time.Now()
}

func (ethClientBalancer *EthClientBalancer) GetClient(ctx context.Context) (*EthClient, error) {
	node := ethClientBalancer.getNode()
	if node == nil {
		return nil, fmt.Errorf("%s", common.ERROR_NO_AVAILABLE_NODE)
	}
	client, ok := ethClientBalancer.clients[node.url]
	if ok {
		return client, nil
	}

	tr := transport{
		OnFailure: func(req *http.Request, err error) {
			ethClientBalancer.increaseFailureCount(node.url)
		},
	}
	rpcClient, err := rpc.DialOptions(
		context.Background(), node.url,
		rpc.WithHTTPClient(&http.Client{Transport: &tr}),
	)
	if err != nil {
		return nil, err
	}

	return &EthClient{Client: ethclient.NewClient(rpcClient)}, nil
}
