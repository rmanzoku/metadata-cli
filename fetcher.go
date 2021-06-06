package fetcher

import (
	"bytes"
	"context"
	"encoding/hex"
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

func TokenURI(ctx context.Context, rpc string, contract string, tokenId *big.Int) (*Metadata, error) {
	buf := make([]byte, 32)
	tokenId.FillBytes(buf)
	data := "0x" + "c87b56dd" + hex.EncodeToString(buf)

	in := struct {
		To   string `json:"to"`
		Data string `json:"data"`
	}{
		To:   contract,
		Data: data,
	}

	param, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}
	latest, _ := json.Marshal("latest")

	result, err := call(ctx, rpc, "eth_call", []json.RawMessage{param, latest})
	if err != nil {
		return nil, err
	}
	b, err := decodeString(result)
	if err != nil {
		return nil, err
	}
	uri := string(b[64:]) // remove string header
	fmt.Println(uri)
	return nil, nil
}

type rpcRequest struct {
	JsonRPC string            `json:"jsonrpc"`
	Method  string            `json:"method"`
	Params  []json.RawMessage `json:"params"`
	ID      uint              `json:"id"`
}

type rpcResponce struct {
	JsonRPC string `json:"jsonrpc"`
	Result  string `json:"result"`
	ID      uint   `json:"id"`
	Error   *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func call(ctx context.Context, rpc, method string, params []json.RawMessage) (string, error) {
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

	if res.Error != nil {
		return "", fmt.Errorf("code:%d message:%s", res.Error.Code, res.Error.Message)
	}

	return res.Result, nil
}
