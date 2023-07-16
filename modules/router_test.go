package modules

import (
	"reflect"
	"testing"
)

func TestTokenRouter(t *testing.T) {
	pools := []PoolInfo{
		{
			CoinAAddress: "A",
			CoinBAddress: "B",
		},
		{
			CoinAAddress: "A",
			CoinBAddress: "C",
		},
		{
			CoinAAddress: "D",
			CoinBAddress: "B",
		},
		{
			CoinAAddress: "C",
			CoinBAddress: "D",
		},
		{
			CoinAAddress: "D",
			CoinBAddress: "E",
		},
		{
			CoinAAddress: "C",
			CoinBAddress: "F",
		},
		{
			CoinAAddress: "G",
			CoinBAddress: "H",
		},
	}
	type args struct {
		pools   []PoolInfo
		coinIn  string
		coinOut string
	}
	tests := []struct {
		name string
		args args
		want []Router
	}{
		{
			name: "case a->b",
			args: args{
				pools:   pools,
				coinIn:  "A",
				coinOut: "B",
			},
			want: []Router{
				{
					Pools: []PoolInfo{{
						CoinAAddress: "A",
						CoinBAddress: "B",
					}},
					Path:  []string{"A", "B"},
					IsA2B: []bool{true},
				},
			},
		},
		{
			name: "case b->a",
			args: args{
				pools:   pools,
				coinIn:  "B",
				coinOut: "A",
			},
			want: []Router{
				{
					Pools: []PoolInfo{{
						CoinAAddress: "A",
						CoinBAddress: "B",
					}},
					Path:  []string{"B", "A"},
					IsA2B: []bool{false},
				},
			},
		},
		{
			name: "case a->d",
			args: args{
				pools:   pools,
				coinIn:  "A",
				coinOut: "D",
			},
			want: []Router{
				{
					Pools: []PoolInfo{{
						CoinAAddress: "A",
						CoinBAddress: "B",
					}, {
						CoinAAddress: "D",
						CoinBAddress: "B",
					}},
					Path:  []string{"A", "B", "D"},
					IsA2B: []bool{true, false},
				},
				{
					Pools: []PoolInfo{{
						CoinAAddress: "A",
						CoinBAddress: "C",
					}, {
						CoinAAddress: "C",
						CoinBAddress: "D",
					}},
					Path:  []string{"A", "C", "D"},
					IsA2B: []bool{true, true},
				},
			},
		},
		{
			name: "case a->e",
			args: args{
				pools:   pools,
				coinIn:  "A",
				coinOut: "E",
			},
			want: []Router{},
		},
		{
			name: "case a->h",
			args: args{
				pools:   pools,
				coinIn:  "A",
				coinOut: "H",
			},
			want: []Router{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TokenRouter(tt.args.pools, tt.args.coinIn, tt.args.coinOut); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TokenRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}
