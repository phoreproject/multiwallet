package util

import "github.com/phoreproject/wallet-interface"

// All implemented coins currently have 100m satoshis per coin
func SatoshisPerCoin(coinType wallet.CoinType) float64 {
	return 100000000
}
