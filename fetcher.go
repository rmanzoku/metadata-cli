package fetcher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"math/big"
	"net/http"
)

type Metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

func (m *Metadata) DecodeImage() (image.Image, error) {
	return nil, nil
}

func FetchMetadata(ctx context.Context, rpc string, contract string, tokenId *big.Int) (*Metadata, error) {
	fmt.Println(rpc, contract, tokenId.Text(10))
	result, err := call(ctx, rpc, "net_version", nil)
	if err != nil {
		return nil, err
	}
	fmt.Println(result)
	return nil, nil
}

type rpcRequest struct {
	JsonRPC string   `json:"jsonrpc"`
	Method  string   `json:"method"`
	Params  []string `json:"params"`
	ID      uint     `json:"id"`
}

type rpcResponce struct {
	JsonRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      uint   `json:"id"`
}

func call(ctx context.Context, rpc, method string, params []string) (string, error) {
	in := &rpcRequest{
		JsonRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1010101,
	}
	input, err := json.Marshal(in)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", rpc, bytes.NewReader(input))
	if err != nil {
		return "", err
	}

	cli := new(http.Client)
	resp, err := cli.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	res := new(rpcResponce)
	err = json.Unmarshal(b, res)
	if err != nil {
		return "", err
	}

	return res.Result, nil
}
