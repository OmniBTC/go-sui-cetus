package modules

type Router struct {
	Pools []PoolInfo
	Path  []string // len(path) = len(pools) + 1
	IsA2B []bool   // len(isa2b) == len(path)
}

func TokenRouter(pools []PoolInfo, coinIn string, coinOut string) []Router {
	routers := make([]Router, 0)
	coin2pools := make(map[string][]PoolInfo)

	// one step
	for _, pool := range pools {
		coin2pools[pool.CoinAAddress] = append(coin2pools[pool.CoinAAddress], pool)
		coin2pools[pool.CoinBAddress] = append(coin2pools[pool.CoinBAddress], pool)

		isPool, isA2b := isPoolMatch(pool, coinIn, coinOut)
		if isPool {
			routers = append(routers, Router{
				Pools: []PoolInfo{pool},
				Path:  []string{coinIn, coinOut},
				IsA2B: []bool{isA2b},
			})
		}
	}

	// two step
	for _, coinInPool := range coin2pools[coinIn] {
		middleCoin := coinInPool.CoinBAddress
		firstPoolIsA2B := true
		if coinIn == middleCoin {
			firstPoolIsA2B = false
			middleCoin = coinInPool.CoinAAddress
		}
		if middleCoin == coinOut {
			continue
		}

		for _, coinOutPool := range coin2pools[middleCoin] {
			if coinOutPool.CoinAAddress == coinOut || coinOutPool.CoinBAddress == coinOut {
				routers = append(routers, Router{
					Pools: []PoolInfo{coinInPool, coinOutPool},
					Path:  []string{coinIn, middleCoin, coinOut},
					IsA2B: []bool{firstPoolIsA2B, coinOutPool.CoinBAddress == coinOut},
				})
			}
		}
	}

	return routers
}

func isPoolMatch(pool PoolInfo, coinA string, coinB string) (isPool bool, isA2b bool) {
	if coinA == pool.CoinAAddress {
		if coinB == pool.CoinBAddress {
			return true, true
		} else {
			return false, false
		}
	} else if coinA == pool.CoinBAddress {
		if coinB == pool.CoinAAddress {
			return true, false
		} else {
			return false, false
		}
	} else {
		return false, false
	}
}
