package fetcher_test

import (
	"context"
	"math/big"
	"os"
	"reflect"
	"testing"

	fetcher "github.com/rmanzoku/nft-metadata-fetcher"
)

var (
	rpc = "https://mainnet.infura.io/v3/" + os.Getenv("INFURA_KEY")
)

type args struct {
	rpc      string
	contract string
	tokenId  *big.Int
}

var argsUniswap = args{
	rpc:      rpc,
	contract: "0xc36442b4a4522e871399cd717abdd847ab11fe88",
	tokenId:  big.NewInt(1),
}

var argsMCHH = args{
	rpc:      rpc,
	contract: "0x273f7f8e6489682df756151f5525576e322d51a3",
	tokenId:  big.NewInt(1),
}

func TestFetchMetadata(t *testing.T) {
	tests := []struct {
		name    string
		args    args
		want    *fetcher.Metadata
		wantErr bool
	}{
		{"uniswap", argsUniswap, nil, false},
		{"mch hero", argsMCHH, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fetcher.FetchMetadata(context.TODO(), tt.args.rpc, tt.args.contract, tt.args.tokenId)
			if (err != nil) != tt.wantErr {
				t.Errorf("FetchMetadata() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FetchMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}
