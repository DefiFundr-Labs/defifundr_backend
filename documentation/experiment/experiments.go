package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"math/big"
// 	"net/http"
// 	"strings"
// 	"sync"
// 	"time"
// )

// // Configuration with REAL CU costs from your usage
// const (
// 	MORALIS_API_KEY = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub25jZSI6IjQzMTQxMTNmLWNmZWYtNGE3YS1iMDM4LWZiYThkOTAwNzMyOSIsIm9yZ0lkIjoiMzU1MTg3IiwidXNlcklkIjoiMzY1MDYyIiwidHlwZUlkIjoiMmM3MTgwMjQtNDM1NS00MzdmLTg0YTMtZWMyMDUxMDE5YmVhIiwidHlwZSI6IlBST0pFQ1QiLCJpYXQiOjE2OTMzNjM2ODEsImV4cCI6NDg0OTEyMzY4MX0.xZkNPRNI7hxM-2nvJNt3IkfhacX6IPiP8brQIYIBO_0"
// 	BASE_URL        = "https://deep-index.moralis.io/api/v2.2"
// 	CU_PER_CALL     = 450  // Your REAL CU cost per call (not 1000!)
// 	DAILY_CU_LIMIT  = 40000
// 	MAX_DAILY_CALLS = DAILY_CU_LIMIT / CU_PER_CALL // 88 calls per day max
// )

// // Enhanced caching with different strategies
// const (
// 	FRESH_DATA_CACHE    = 5 * time.Minute   // Fresh API data
// 	STALE_DATA_CACHE    = 2 * time.Hour     // Older but acceptable data  
// 	SKELETON_CACHE      = 24 * time.Hour    // Zero balances skeleton
// 	POPULAR_WALLET_CACHE = 1 * time.Hour    // High-value wallets
// )

// // Complete network configuration - ALL supported chains
// type Network struct {
// 	Name        string
// 	ChainID     string
// 	Priority    int  // 1=highest, 5=lowest for CU optimization
// 	IsTestnet   bool
// 	IsSupported bool // Some networks might not be fully supported yet
// }

// var allNetworks = map[string]Network{
// 	// MAINNET - Priority 1-3
// 	"ethereum": {Name: "Ethereum", ChainID: "0x1", Priority: 1, IsTestnet: false, IsSupported: true},
// 	"polygon":  {Name: "Polygon", ChainID: "0x89", Priority: 1, IsTestnet: false, IsSupported: true},
// 	"base":     {Name: "Base", ChainID: "0x2105", Priority: 2, IsTestnet: false, IsSupported: true},
// 	"bsc":      {Name: "Binance Smart Chain", ChainID: "0x38", Priority: 2, IsTestnet: false, IsSupported: true},
// 	"arbitrum": {Name: "Arbitrum", ChainID: "0xa4b1", Priority: 2, IsTestnet: false, IsSupported: true},
// 	"optimism": {Name: "Optimism", ChainID: "0xa", Priority: 3, IsTestnet: false, IsSupported: true},
// 	"starknet": {Name: "StarkNet", ChainID: "starknet", Priority: 4, IsTestnet: false, IsSupported: false},
	
// 	// TESTNET - Priority 4-5
// 	"sepolia":          {Name: "Ethereum Sepolia", ChainID: "0xaa36a7", Priority: 4, IsTestnet: true, IsSupported: true},
// 	"amoy":             {Name: "Polygon Amoy", ChainID: "0x13882", Priority: 4, IsTestnet: true, IsSupported: true},
// 	"bsc_testnet":      {Name: "BSC Testnet", ChainID: "0x61", Priority: 4, IsTestnet: true, IsSupported: true},
// 	"arbitrum_sepolia": {Name: "Arbitrum Sepolia", ChainID: "421614", Priority: 5, IsTestnet: true, IsSupported: true},
// 	"optimism_sepolia": {Name: "Optimism Sepolia", ChainID: "11155420", Priority: 5, IsTestnet: true, IsSupported: true},
// 	"base_sepolia":     {Name: "Base Sepolia", ChainID: "0x14a34", Priority: 5, IsTestnet: true, IsSupported: true},
// 	"starknet_sepolia": {Name: "StarkNet Sepolia", ChainID: "starknet-sepolia", Priority: 5, IsTestnet: true, IsSupported: false},
// }

// // Complete stablecoin addresses for ALL networks
// var allStablecoinAddresses = map[string]map[string]string{
// 	"ethereum": {
// 		"USDT": "0xdAC17F958D2ee523a2206206994597C13D831ec7",
// 		"USDC": "0xA0b86a33E6Fa0E7834c8fa9a7E5b20D1f9b5d8a7",
// 		"DAI":  "0x6B175474E89094C44Da98b954EedeAC495271d0F",
// 		"USDD": "0x0C10bF8FcB7Bf5412187A595ab97a3609160b5c6",
// 		"LUSD": "0x5f98805A4E8be255a32880FDeC7F6728C6568bA0",
// 		"EURT": "0xC581b735A1688071A1746c968e0798D642EDE491",
// 	},
// 	"polygon": {
// 		"USDT": "0xc2132D05D31c914a87C6611C10748AEb04B58e8F",
// 		"USDC": "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174",
// 		"DAI":  "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063",
// 		"USDD": "0xFFA4D863C96e743A2e1513824EA006B8D0353C57",
// 	},
// 	"base": {
// 		"USDT": "0xfde4C96c8593536E31F229EA8f37b2ADa2699bb2",
// 		"USDC": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
// 		"DAI":  "0x50c5725949A6F0c72E6C4a641F24049A917DB0Cb",
// 	},
// 	"bsc": {
// 		"USDT": "0x55d398326f99059fF775485246999027B3197955",
// 		"USDC": "0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d",
// 		"DAI":  "0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3",
// 		"USDD": "0xd17479997F34dd9156Deef8F95A52D81D265be9c",
// 	},
// 	"arbitrum": {
// 		"USDT": "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9",
// 		"USDC": "0xaf88d065e77c8cC2239327C5EDb3A432268e5831",
// 		"DAI":  "0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1",
// 	},
// 	"optimism": {
// 		"USDT": "0x94b008aA00579c1307B0EF2c499aD98a8ce58e58",
// 		"USDC": "0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85",
// 		"DAI":  "0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1",
// 	},
// 	"starknet": {
// 		"USDT": "0x068f5c6a61780768455de69077e07e89787839bf8166decfbf92b645209c0fb8",
// 		"USDC": "0x053c91253bc9682c04929ca02ed00b3e423f6710d2ee7e0d5ebb06f3ecf368a8",
// 		"DAI":  "0x00da114221cb83fa859dbdb4c44beeaa0bb37c7537ad5ae66fe5e0efd20e6eb3",
// 	},
// 	// Testnets use same addresses as mainnets mostly
// 	// "sepolia":          {"USDT": "0xdAC17F958D2ee523a2206206994597C13D831ec7", "USDC": "0xA0b86a33E6Fa0E7834c8fa9a7E5b20D1f9b5d8a7", "DAI": "0x6B175474E89094C44Da98b954EedeAC495271d0F"},
// 	// "amoy":             {"USDT": "0xc2132D05D31c914a87C6611C10748AEb04B58e8F", "USDC": "0x2791Bca1f2de4661ED88A30C99A7a9449Aa84174", "DAI": "0x8f3Cf7ad23Cd3CaDbD9735AFf958023239c6A063"},
// 	// "bsc_testnet":      {"USDT": "0x55d398326f99059fF775485246999027B3197955", "USDC": "0x8AC76a51cc950d9822D68b83fE1Ad97B32Cd580d", "DAI": "0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3"},
// 	// "arbitrum_sepolia": {"USDT": "0xFd086bC7CD5C481DCC9C85ebE478A1C0b69FCbb9", "USDC": "0xaf88d065e77c8cC2239327C5EDb3A432268e5831", "DAI": "0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1"},
// 	// "optimism_sepolia": {"USDT": "0x94b008aA00579c1307B0EF2c499aD98a8ce58e58", "USDC": "0x0b2C639c533813f4Aa9D7837CAf62653d097Ff85", "DAI": "0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1"},
// 	// "base_sepolia":     {"USDT": "0xfde4C96c8593536E31F229EA8f37b2ADa2699bb2", "USDC": "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", "DAI": "0x50c5725949A6F0c72E6C4a641F24049A917DB0Cb"},
// 	// "starknet_sepolia": {"USDT": "0x068f5c6a61780768455de69077e07e89787839bf8166decfbf92b645209c0fb8", "USDC": "0x053c91253bc9682c04929ca02ed00b3e423f6710d2ee7e0d5ebb06f3ecf368a8", "DAI": "0x00da114221cb83fa859dbdb4c44beeaa0bb37c7537ad5ae66fe5e0efd20e6eb3"},
// }

// // Enhanced wallet balance with complete metadata
// type WalletBalance struct {
// 	Network         string    `json:"network"`
// 	NetworkType     string    `json:"network_type"` // "mainnet" or "testnet"
// 	Token           string    `json:"token"`
// 	TokenName       string    `json:"token_name"`
// 	Balance         string    `json:"balance"`
// 	FormattedBalance string   `json:"formatted_balance"`
// 	Decimals        int       `json:"decimals"`
// 	USDValue        float64   `json:"usd_value"`
// 	Address         string    `json:"address"`
// 	Logo            string    `json:"logo"`
// 	IsZero          bool      `json:"is_zero"`
// 	DataSource      string    `json:"data_source"` // "fresh_api", "cached_api", "skeleton", "unsupported"
// 	LastUpdated     time.Time `json:"last_updated"`
// 	CacheAge        string    `json:"cache_age,omitempty"`
// 	Priority        int       `json:"priority"`
// 	IsSupported     bool      `json:"is_supported"`
// 	Error           string    `json:"error,omitempty"`
// }

// // Smart fetching strategy
// type FetchStrategy struct {
// 	UseAPIForPriority1   bool // Always fetch Ethereum, Polygon  
// 	UseAPIForPriority2   bool // Fetch Base, BSC, Arbitrum if CU available
// 	UseAPIForPriority3   bool // Fetch Optimism if CU available
// 	UseAPIForTestnets    bool // Fetch testnets if CU available
// 	ForceRefreshPopular  bool // Force refresh popular wallets
// 	EmergencyModeActive  bool // Only priority 1 networks
// }

// // CU usage tracking and protection
// type CUTracker struct {
// 	dailyUsage      int
// 	lastResetTime   time.Time
// 	callsToday      int
// 	emergencyMode   bool
// 	mutex           sync.Mutex
// 	hourlyUsage     map[int]int // Track usage by hour
// }

// var cuTracker = &CUTracker{
// 	lastResetTime: time.Now(),
// 	hourlyUsage:   make(map[int]int),
// }

// // Enhanced cache with different data types
// type PortfolioCache struct {
// 	freshData    map[string]CacheEntry // Recent API data
// 	staleData    map[string]CacheEntry // Older API data (still usable)
// 	skeletonData map[string]CacheEntry // Zero balance templates
// 	mutex        sync.RWMutex
// }

// type CacheEntry struct {
// 	Data         []WalletBalance
// 	Timestamp    time.Time
// 	AccessCount  int
// 	IsPopular    bool
// 	TotalUSDValue float64
// }

// var portfolioCache = &PortfolioCache{
// 	freshData:    make(map[string]CacheEntry),
// 	staleData:    make(map[string]CacheEntry),
// 	skeletonData: make(map[string]CacheEntry),
// }

// // Complete portfolio response
// type PortfolioResponse struct {
// 	Success         bool            `json:"success"`
// 	WalletAddress   string          `json:"wallet_address"`
// 	TotalNetworks   int             `json:"total_networks"`
// 	TotalTokens     int             `json:"total_tokens"`
// 	MainnetBalances []WalletBalance `json:"mainnet_balances"`
// 	TestnetBalances []WalletBalance `json:"testnet_balances"`
// 	Summary         PortfolioSummary `json:"summary"`
// 	Metadata        ResponseMetadata `json:"metadata"`
// 	Timestamp       time.Time       `json:"timestamp"`
// }

// type PortfolioSummary struct {
// 	TotalMainnetUSD   float64 `json:"total_mainnet_usd"`
// 	TotalTestnetUSD   float64 `json:"total_testnet_usd"`
// 	GrandTotalUSD     float64 `json:"grand_total_usd"`
// 	NetworksWithBalance int   `json:"networks_with_balance"`
// 	TokensWithBalance   int   `json:"tokens_with_balance"`
// 	ZeroBalanceTokens   int   `json:"zero_balance_tokens"`
// }

// type ResponseMetadata struct {
// 	CacheHitRate    float64         `json:"cache_hit_rate"`
// 	CUUsed          int             `json:"cu_used"`
// 	RemainingCU     int             `json:"remaining_cu"`
// 	DataSources     map[string]int  `json:"data_sources"` // Count by source type
// 	FetchStrategy   FetchStrategy   `json:"fetch_strategy"`
// 	ProcessingTime  string          `json:"processing_time"`
// 	NetworkStats    map[string]int  `json:"network_stats"` // Count by network type
// 	Message         string          `json:"message,omitempty"`
// }

// // CU Management Functions
// func (tracker *CUTracker) determineFetchStrategy() FetchStrategy {
// 	tracker.mutex.Lock()
// 	defer tracker.mutex.Unlock()

// 	// Reset daily counter if new day
// 	if time.Since(tracker.lastResetTime) > 24*time.Hour {
// 		tracker.dailyUsage = 0
// 		tracker.callsToday = 0
// 		tracker.lastResetTime = time.Now()
// 		tracker.emergencyMode = false
// 		tracker.hourlyUsage = make(map[int]int)
// 	}

// 	remainingCU := DAILY_CU_LIMIT - tracker.dailyUsage
// 	remainingCalls := remainingCU / CU_PER_CALL

// 	// Determine strategy based on remaining CU
// 	strategy := FetchStrategy{}

// 	if remainingCalls >= 12 {
// 		// Full service - all networks (you can afford this!)
// 		strategy.UseAPIForPriority1 = true
// 		strategy.UseAPIForPriority2 = true
// 		strategy.UseAPIForPriority3 = true
// 		strategy.UseAPIForTestnets = true // You can afford testnets now!
// 	} else if remainingCalls >= 8 {
// 		// High service - most networks
// 		strategy.UseAPIForPriority1 = true
// 		strategy.UseAPIForPriority2 = true
// 		strategy.UseAPIForPriority3 = true
// 		strategy.UseAPIForTestnets = false
// 	} else if remainingCalls >= 4 {
// 		// Standard service - priority networks
// 		strategy.UseAPIForPriority1 = true
// 		strategy.UseAPIForPriority2 = true
// 		strategy.UseAPIForPriority3 = false
// 		strategy.UseAPIForTestnets = false
// 	} else if remainingCalls >= 2 {
// 		// Minimal service - top priority only
// 		strategy.UseAPIForPriority1 = true
// 		strategy.UseAPIForPriority2 = false
// 		strategy.UseAPIForPriority3 = false
// 		strategy.UseAPIForTestnets = false
// 	} else {
// 		// Emergency mode - cache only
// 		strategy.UseAPIForPriority1 = false
// 		strategy.UseAPIForPriority2 = false
// 		strategy.UseAPIForPriority3 = false
// 		strategy.UseAPIForTestnets = false
// 		tracker.emergencyMode = true
// 	}

// 	return strategy
// }

// func (tracker *CUTracker) recordUsage(calls int) {
// 	tracker.mutex.Lock()
// 	defer tracker.mutex.Unlock()
	
// 	cuUsed := calls * CU_PER_CALL
// 	tracker.dailyUsage += cuUsed
// 	tracker.callsToday += calls
	
// 	// Track hourly usage
// 	hour := time.Now().Hour()
// 	tracker.hourlyUsage[hour] += calls
// }

// func (tracker *CUTracker) getRemainingCU() int {
// 	tracker.mutex.Lock()
// 	defer tracker.mutex.Unlock()
// 	return DAILY_CU_LIMIT - tracker.dailyUsage
// }

// // Create complete portfolio skeleton (all networks, all tokens, zero balances)
// func createPortfolioSkeleton(walletAddress string) []WalletBalance {
// 	var skeleton []WalletBalance

// 	// Generate skeleton for ALL networks and tokens
// 	for networkKey, network := range allNetworks {
// 		tokens := allStablecoinAddresses[networkKey]
// 		if tokens == nil {
// 			// Use mainnet tokens as fallback for unsupported networks
// 			tokens = allStablecoinAddresses["ethereum"]
// 		}

// 		for tokenSymbol, tokenAddress := range tokens {
// 			networkType := "mainnet"
// 			if network.IsTestnet {
// 				networkType = "testnet"
// 			}

// 			dataSource := "skeleton"
// 			if !network.IsSupported {
// 				dataSource = "unsupported"
// 			}

// 			balance := WalletBalance{
// 				Network:         network.Name,
// 				NetworkType:     networkType,
// 				Token:           tokenSymbol,
// 				TokenName:       tokenSymbol + " (" + network.Name + ")",
// 				Balance:         "0",
// 				FormattedBalance: "0.00",
// 				Decimals:        getDefaultDecimals(tokenSymbol),
// 				USDValue:        0.0,
// 				Address:         tokenAddress,
// 				Logo:            getFallbackLogo(tokenAddress, tokenSymbol),
// 				IsZero:          true,
// 				DataSource:      dataSource,
// 				LastUpdated:     time.Now(),
// 				Priority:        network.Priority,
// 				IsSupported:     network.IsSupported,
// 			}

// 			if !network.IsSupported {
// 				balance.Error = "Network not yet supported - coming soon"
// 			}

// 			skeleton = append(skeleton, balance)
// 		}
// 	}

// 	return skeleton
// }

// // Smart portfolio fetching with complete coverage
// func GetCompletePortfolio(walletAddress, userID string) PortfolioResponse {
// 	startTime := time.Now()
	
// 	// Initialize response with skeleton
// 	response := PortfolioResponse{
// 		Success:       true,
// 		WalletAddress: walletAddress,
// 		Timestamp:     time.Now(),
// 	}

// 	// Create complete skeleton first (shows ALL networks/tokens)
// 	skeletonBalances := createPortfolioSkeleton(walletAddress)
	
// 	// Determine fetch strategy based on CU availability
// 	strategy := cuTracker.determineFetchStrategy()
// 	response.Metadata.FetchStrategy = strategy

// 	// Get cached data if available
// 	freshBalances, staleBalances := getFromCache(walletAddress)
	
// 	// Merge skeleton with any cached data
// 	portfolioBalances := mergeBalances(skeletonBalances, freshBalances, staleBalances)
	
// 	// Determine which networks need fresh API calls
// 	networksToFetch := determineNetworksToFetch(strategy, portfolioBalances)
	
// 	// Make API calls for selected networks
// 	cuUsed := 0
// 	dataSources := make(map[string]int)
	
// 	if len(networksToFetch) > 0 {
// 		freshAPIBalances := fetchNetworkBalances(walletAddress, networksToFetch)
// 		cuUsed = len(networksToFetch)
// 		cuTracker.recordUsage(cuUsed)
		
// 		// Merge fresh API data
// 		portfolioBalances = mergeBalances(portfolioBalances, freshAPIBalances, nil)
		
// 		// Cache the fresh data
// 		cacheBalances(walletAddress, freshAPIBalances, true)
// 	}

// 	// Count data sources
// 	for _, balance := range portfolioBalances {
// 		dataSources[balance.DataSource]++
// 	}

// 	// Separate mainnet and testnet
// 	var mainnetBalances, testnetBalances []WalletBalance
// 	for _, balance := range portfolioBalances {
// 		if balance.NetworkType == "testnet" {
// 			testnetBalances = append(testnetBalances, balance)
// 		} else {
// 			mainnetBalances = append(mainnetBalances, balance)
// 		}
// 	}

// 	// Calculate summary
// 	summary := calculatePortfolioSummary(portfolioBalances)
	
// 	// Build response
// 	response.MainnetBalances = mainnetBalances
// 	response.TestnetBalances = testnetBalances
// 	response.TotalNetworks = len(allNetworks)
// 	response.TotalTokens = len(portfolioBalances)
// 	response.Summary = summary
	
// 	// Calculate cache hit rate
// 	totalItems := len(portfolioBalances)
// 	cacheHits := dataSources["fresh_cache"] + dataSources["stale_cache"] + dataSources["skeleton"]
// 	cacheHitRate := float64(cacheHits) / float64(totalItems) * 100
	
// 	response.Metadata = ResponseMetadata{
// 		CacheHitRate:   cacheHitRate,
// 		CUUsed:         cuUsed * CU_PER_CALL,
// 		RemainingCU:    cuTracker.getRemainingCU(),
// 		DataSources:    dataSources,
// 		FetchStrategy:  strategy,
// 		ProcessingTime: time.Since(startTime).Round(time.Millisecond).String(),
// 		NetworkStats: map[string]int{
// 			"mainnet_networks": len(getNetworksByType(false)),
// 			"testnet_networks": len(getNetworksByType(true)),
// 			"supported_networks": len(getSupportedNetworks()),
// 			"fresh_api_calls": cuUsed,
// 		},
// 	}

// 	// Add informative message
// 	if cuUsed == 0 {
// 		response.Metadata.Message = "All data served from cache - no CU used"
// 	} else if strategy.EmergencyModeActive {
// 		response.Metadata.Message = fmt.Sprintf("Emergency mode: only priority 1 networks checked (%d CU used)", cuUsed*CU_PER_CALL)
// 	} else {
// 		response.Metadata.Message = fmt.Sprintf("Fresh data from %d networks (%d CU used)", cuUsed, cuUsed*CU_PER_CALL)
// 	}

// 	return response
// }

// // Helper functions
// func getFromCache(walletAddress string) (fresh []WalletBalance, stale []WalletBalance) {
// 	portfolioCache.mutex.RLock()
// 	defer portfolioCache.mutex.RUnlock()

// 	if entry, exists := portfolioCache.freshData[walletAddress]; exists {
// 		if time.Since(entry.Timestamp) < FRESH_DATA_CACHE {
// 			fresh = entry.Data
// 		} else if time.Since(entry.Timestamp) < STALE_DATA_CACHE {
// 			// Move to stale cache
// 			go func() {
// 				portfolioCache.mutex.Lock()
// 				portfolioCache.staleData[walletAddress] = entry
// 				delete(portfolioCache.freshData, walletAddress)
// 				portfolioCache.mutex.Unlock()
// 			}()
// 			stale = entry.Data
// 		}
// 	}

// 	if entry, exists := portfolioCache.staleData[walletAddress]; exists {
// 		if time.Since(entry.Timestamp) < STALE_DATA_CACHE {
// 			stale = append(stale, entry.Data...)
// 		}
// 	}

// 	return fresh, stale
// }

// func mergeBalances(skeleton, fresh, stale []WalletBalance) []WalletBalance {
// 	// Create map for easy lookup
// 	balanceMap := make(map[string]WalletBalance)
	
// 	// Start with skeleton (all zeros)
// 	for _, balance := range skeleton {
// 		key := balance.Network + "-" + balance.Token
// 		balanceMap[key] = balance
// 	}
	
// 	// Override with stale data if available
// 	for _, balance := range stale {
// 		key := balance.Network + "-" + balance.Token
// 		balance.DataSource = "stale_cache"
// 		balance.CacheAge = time.Since(balance.LastUpdated).Round(time.Minute).String()
// 		balanceMap[key] = balance
// 	}
	
// 	// Override with fresh data (highest priority)
// 	for _, balance := range fresh {
// 		key := balance.Network + "-" + balance.Token
// 		balance.DataSource = "fresh_cache"
// 		if balance.LastUpdated.IsZero() {
// 			balance.DataSource = "fresh_api"
// 			balance.LastUpdated = time.Now()
// 		}
// 		balanceMap[key] = balance
// 	}
	
// 	// Convert back to slice
// 	var result []WalletBalance
// 	for _, balance := range balanceMap {
// 		result = append(result, balance)
// 	}
	
// 	return result
// }

// func determineNetworksToFetch(strategy FetchStrategy, currentBalances []WalletBalance) []string {
// 	var networksToFetch []string
	
// 	// Check which networks need fresh data
// 	for networkKey, network := range allNetworks {
// 		if !network.IsSupported {
// 			continue // Skip unsupported networks
// 		}
		
// 		shouldFetch := false
		
// 		// Check strategy
// 		if strategy.UseAPIForPriority1 && network.Priority == 1 {
// 			shouldFetch = true
// 		} else if strategy.UseAPIForPriority2 && network.Priority == 2 {
// 			shouldFetch = true
// 		} else if strategy.UseAPIForPriority3 && network.Priority == 3 {
// 			shouldFetch = true
// 		} else if strategy.UseAPIForTestnets && network.IsTestnet {
// 			shouldFetch = true
// 		}
		
// 		if shouldFetch {
// 			// Check if we have recent data for this network
// 			hasRecentData := false
// 			for _, balance := range currentBalances {
// 				if balance.Network == network.Name && 
// 				   (balance.DataSource == "fresh_api" || balance.DataSource == "fresh_cache") &&
// 				   time.Since(balance.LastUpdated) < FRESH_DATA_CACHE {
// 					hasRecentData = true
// 					break
// 				}
// 			}
			
// 			if !hasRecentData {
// 				networksToFetch = append(networksToFetch, networkKey)
// 			}
// 		}
// 	}
	
// 	return networksToFetch
// }

// func fetchNetworkBalances(walletAddress string, networks []string) []WalletBalance {
// 	var allBalances []WalletBalance
	
// 	for _, networkKey := range networks {
// 		network := allNetworks[networkKey]
// 		balances := getNetworkBalances(walletAddress, networkKey, network)
// 		allBalances = append(allBalances, balances...)
// 	}
	
// 	return allBalances
// }

// func getNetworkBalances(walletAddress, networkKey string, network Network) []WalletBalance {
// 	var balances []WalletBalance
	
// 	url := fmt.Sprintf("%s/%s/erc20?chain=%s", BASE_URL, walletAddress, network.ChainID)
// 	resp, err := makeAPIRequest(url)
// 	if err != nil {
// 		// Return skeleton with error for this network
// 		tokens := allStablecoinAddresses[networkKey]
// 		for tokenSymbol, tokenAddress := range tokens {
// 			balances = append(balances, WalletBalance{
// 				Network:         network.Name,
// 				NetworkType:     getNetworkType(network.IsTestnet),
// 				Token:           tokenSymbol,
// 				TokenName:       tokenSymbol,
// 				Balance:         "0",
// 				FormattedBalance: "0.00",
// 				Decimals:        getDefaultDecimals(tokenSymbol),
// 				USDValue:        0.0,
// 				Address:         tokenAddress,
// 				Logo:            getFallbackLogo(tokenAddress, tokenSymbol),
// 				IsZero:          true,
// 				DataSource:      "api_error",
// 				LastUpdated:     time.Now(),
// 				Priority:        network.Priority,
// 				IsSupported:     network.IsSupported,
// 				Error:           err.Error(),
// 			})
// 		}
// 		return balances
// 	}

// 	expectedTokens := allStablecoinAddresses[networkKey]
// 	foundTokens := make(map[string]TokenBalance)
	
// 	for _, token := range *resp {
// 		normalizedAddr := strings.ToLower(token.TokenAddress)
// 		foundTokens[normalizedAddr] = token
// 	}

// 	// Create balance entry for each expected token
// 	for tokenSymbol, tokenAddress := range expectedTokens {
// 		normalizedAddr := strings.ToLower(tokenAddress)
// 		networkType := getNetworkType(network.IsTestnet)
		
// 		if foundToken, exists := foundTokens[normalizedAddr]; exists && foundToken.Balance != "0" {
// 			// Non-zero balance found
// 			usdValue := calculateUSDValue(foundToken.Balance, foundToken.Decimals)
// 			formattedBalance := formatBalance(foundToken.Balance, foundToken.Decimals)
			
// 			balances = append(balances, WalletBalance{
// 				Network:         network.Name,
// 				NetworkType:     networkType,
// 				Token:           tokenSymbol,
// 				TokenName:       foundToken.Name,
// 				Balance:         foundToken.Balance,
// 				FormattedBalance: formattedBalance,
// 				Decimals:        foundToken.Decimals,
// 				USDValue:        usdValue,
// 				Address:         tokenAddress,
// 				Logo:            getTokenLogo(foundToken.Logo, tokenAddress, tokenSymbol),
// 				IsZero:          false,
// 				DataSource:      "fresh_api",
// 				LastUpdated:     time.Now(),
// 				Priority:        network.Priority,
// 				IsSupported:     network.IsSupported,
// 			})
// 		} else {
// 			// Zero balance or not found
// 			balances = append(balances, WalletBalance{
// 				Network:         network.Name,
// 				NetworkType:     networkType,
// 				Token:           tokenSymbol,
// 				TokenName:       tokenSymbol,
// 				Balance:         "0",
// 				FormattedBalance: "0.00",
// 				Decimals:        getDefaultDecimals(tokenSymbol),
// 				USDValue:        0.0,
// 				Address:         tokenAddress,
// 				Logo:            getFallbackLogo(tokenAddress, tokenSymbol),
// 				IsZero:          true,
// 				DataSource:      "fresh_api",
// 				LastUpdated:     time.Now(),
// 				Priority:        network.Priority,
// 				IsSupported:     network.IsSupported,
// 			})
// 		}
// 	}
	
// 	return balances
// }

// func cacheBalances(walletAddress string, balances []WalletBalance, isFresh bool) {
// 	portfolioCache.mutex.Lock()
// 	defer portfolioCache.mutex.Unlock()
	
// 	totalUSD := 0.0
// 	for _, balance := range balances {
// 		totalUSD += balance.USDValue
// 	}
	
// 	entry := CacheEntry{
// 		Data:          balances,
// 		Timestamp:     time.Now(),
// 		AccessCount:   1,
// 		IsPopular:     totalUSD > 1000, // Consider wallets with >$1k as popular
// 		TotalUSDValue: totalUSD,
// 	}
	
// 	if isFresh {
// 		portfolioCache.freshData[walletAddress] = entry
// 	} else {
// 		portfolioCache.staleData[walletAddress] = entry
// 	}
// }

// func calculatePortfolioSummary(balances []WalletBalance) PortfolioSummary {
// 	var totalMainnetUSD, totalTestnetUSD float64
// 	var networksWithBalance, tokensWithBalance, zeroBalanceTokens int
	
// 	networkHasBalance := make(map[string]bool)
	
// 	for _, balance := range balances {
// 		if balance.NetworkType == "testnet" {
// 			totalTestnetUSD += balance.USDValue
// 		} else {
// 			totalMainnetUSD += balance.USDValue
// 		}
		
// 		if balance.USDValue > 0 {
// 			tokensWithBalance++
// 			networkHasBalance[balance.Network] = true
// 		} else {
// 			zeroBalanceTokens++
// 		}
// 	}
	
// 	networksWithBalance = len(networkHasBalance)
	
// 	return PortfolioSummary{
// 		TotalMainnetUSD:     totalMainnetUSD,
// 		TotalTestnetUSD:     totalTestnetUSD,
// 		GrandTotalUSD:       totalMainnetUSD + totalTestnetUSD,
// 		NetworksWithBalance: networksWithBalance,
// 		TokensWithBalance:   tokensWithBalance,
// 		ZeroBalanceTokens:   zeroBalanceTokens,
// 	}
// }

// // User-friendly API endpoint with filtering options
// func GetPortfolioWithFilters(walletAddress, userID string, options PortfolioOptions) PortfolioResponse {
// 	// Get complete portfolio
// 	portfolio := GetCompletePortfolio(walletAddress, userID)
	
// 	// Apply filters
// 	if options.ShowOnlyNonZero {
// 		portfolio.MainnetBalances = filterNonZeroBalances(portfolio.MainnetBalances)
// 		portfolio.TestnetBalances = filterNonZeroBalances(portfolio.TestnetBalances)
// 	}
	
// 	if options.ShowOnlyMainnet {
// 		portfolio.TestnetBalances = []WalletBalance{}
// 	}
	
// 	if options.ShowOnlyTestnet {
// 		portfolio.MainnetBalances = []WalletBalance{}
// 	}
	
// 	if len(options.SpecificNetworks) > 0 {
// 		portfolio.MainnetBalances = filterByNetworks(portfolio.MainnetBalances, options.SpecificNetworks)
// 		portfolio.TestnetBalances = filterByNetworks(portfolio.TestnetBalances, options.SpecificNetworks)
// 	}
	
// 	if len(options.SpecificTokens) > 0 {
// 		portfolio.MainnetBalances = filterByTokens(portfolio.MainnetBalances, options.SpecificTokens)
// 		portfolio.TestnetBalances = filterByTokens(portfolio.TestnetBalances, options.SpecificTokens)
// 	}
	
// 	// Recalculate summary after filtering
// 	allBalances := append(portfolio.MainnetBalances, portfolio.TestnetBalances...)
// 	portfolio.Summary = calculatePortfolioSummary(allBalances)
// 	portfolio.TotalTokens = len(allBalances)
	
// 	return portfolio
// }

// type PortfolioOptions struct {
// 	ShowOnlyNonZero    bool     `json:"show_only_non_zero"`
// 	ShowOnlyMainnet    bool     `json:"show_only_mainnet"`
// 	ShowOnlyTestnet    bool     `json:"show_only_testnet"`
// 	SpecificNetworks   []string `json:"specific_networks,omitempty"`
// 	SpecificTokens     []string `json:"specific_tokens,omitempty"`
// 	IncludeMetadata    bool     `json:"include_metadata"`
// 	ForceRefresh       bool     `json:"force_refresh"`
// }

// // CU Management and Monitoring
// func GetSystemStatus() map[string]interface{} {
// 	// Get CU usage stats without holding lock too long
// 	cuTracker.mutex.Lock()
// 	dailyUsage := cuTracker.dailyUsage
// 	callsToday := cuTracker.callsToday
// 	emergencyMode := cuTracker.emergencyMode
// 	hourlyUsage := make(map[int]int)
// 	for k, v := range cuTracker.hourlyUsage {
// 		hourlyUsage[k] = v
// 	}
// 	cuTracker.mutex.Unlock()
	
// 	remainingCU := DAILY_CU_LIMIT - dailyUsage
// 	remainingCalls := remainingCU / CU_PER_CALL
// 	usagePercentage := float64(dailyUsage) / float64(DAILY_CU_LIMIT) * 100
	
// 	// Cache statistics
// 	portfolioCache.mutex.RLock()
// 	freshCacheSize := len(portfolioCache.freshData)
// 	staleCacheSize := len(portfolioCache.staleData)
// 	skeletonCacheSize := len(portfolioCache.skeletonData)
// 	portfolioCache.mutex.RUnlock()
	
// 	// Strategy for next requests (call separately to avoid deadlock)
// 	currentStrategy := determineFetchStrategyForStatus(remainingCalls, emergencyMode)
	
// 	return map[string]interface{}{
// 		"cu_usage": map[string]interface{}{
// 			"daily_cu_used":      dailyUsage,
// 			"daily_cu_limit":     DAILY_CU_LIMIT,
// 			"remaining_cu":       remainingCU,
// 			"calls_made_today":   callsToday,
// 			"remaining_calls":    remainingCalls,
// 			"usage_percentage":   usagePercentage,
// 			"emergency_mode":     emergencyMode,
// 			"hourly_breakdown":   hourlyUsage,
// 		},
// 		"cache_stats": map[string]interface{}{
// 			"fresh_cache_size":    freshCacheSize,
// 			"stale_cache_size":    staleCacheSize,
// 			"skeleton_cache_size": skeletonCacheSize,
// 			"total_cached":        freshCacheSize + staleCacheSize + skeletonCacheSize,
// 		},
// 		"network_coverage": map[string]interface{}{
// 			"total_networks":       len(allNetworks),
// 			"supported_networks":   len(getSupportedNetworks()),
// 			"mainnet_networks":     len(getNetworksByType(false)),
// 			"testnet_networks":     len(getNetworksByType(true)),
// 			"priority_1_networks":  len(getNetworksByPriority(1)),
// 			"priority_2_networks":  len(getNetworksByPriority(2)),
// 			"priority_3_networks":  len(getNetworksByPriority(3)),
// 		},
// 		"current_strategy": map[string]interface{}{
// 			"priority_1_enabled":   currentStrategy.UseAPIForPriority1,
// 			"priority_2_enabled":   currentStrategy.UseAPIForPriority2,
// 			"priority_3_enabled":   currentStrategy.UseAPIForPriority3,
// 			"testnets_enabled":     currentStrategy.UseAPIForTestnets,
// 			"emergency_mode":       currentStrategy.EmergencyModeActive,
// 		},
// 		"recommendations": generateRecommendations(usagePercentage, remainingCalls),
// 	}
// }

// // Helper function to determine strategy without mutex conflicts
// func determineFetchStrategyForStatus(remainingCalls int, emergencyMode bool) FetchStrategy {
// 	strategy := FetchStrategy{}

// 	if remainingCalls >= 12 {
// 		// Full service - all networks (you can afford this!)
// 		strategy.UseAPIForPriority1 = true
// 		strategy.UseAPIForPriority2 = true
// 		strategy.UseAPIForPriority3 = true
// 		strategy.UseAPIForTestnets = true // You can afford testnets now!
// 	} else if remainingCalls >= 8 {
// 		// High service - most networks
// 		strategy.UseAPIForPriority1 = true
// 		strategy.UseAPIForPriority2 = true
// 		strategy.UseAPIForPriority3 = true
// 		strategy.UseAPIForTestnets = false
// 	} else if remainingCalls >= 4 {
// 		// Standard service - priority networks
// 		strategy.UseAPIForPriority1 = true
// 		strategy.UseAPIForPriority2 = true
// 		strategy.UseAPIForPriority3 = false
// 		strategy.UseAPIForTestnets = false
// 	} else if remainingCalls >= 2 {
// 		// Minimal service - top priority only
// 		strategy.UseAPIForPriority1 = true
// 		strategy.UseAPIForPriority2 = false
// 		strategy.UseAPIForPriority3 = false
// 		strategy.UseAPIForTestnets = false
// 	} else {
// 		// Emergency mode - cache only
// 		strategy.UseAPIForPriority1 = false
// 		strategy.UseAPIForPriority2 = false
// 		strategy.UseAPIForPriority3 = false
// 		strategy.UseAPIForTestnets = false
// 	}
	
// 	strategy.EmergencyModeActive = emergencyMode
// 	return strategy
// }

// func generateRecommendations(usagePercentage float64, remainingCalls int) []string {
// 	var recommendations []string
	
// 	if usagePercentage > 90 {
// 		recommendations = append(recommendations, "⚠️ Daily CU limit almost exceeded - consider upgrading to Pro plan")
// 		recommendations = append(recommendations, "💾 Increase cache duration to reduce API calls")
// 	} else if usagePercentage > 70 {
// 		recommendations = append(recommendations, "🔧 Enable emergency mode if usage spikes")
// 		recommendations = append(recommendations, "📊 Monitor usage patterns for optimization")
// 	} else if remainingCalls > 20 {
// 		recommendations = append(recommendations, "✅ CU usage healthy - all features available")
// 		recommendations = append(recommendations, "🚀 Consider enabling testnet fetching")
// 	}
	
// 	if remainingCalls < 5 {
// 		recommendations = append(recommendations, "🚨 Switch to cache-only mode to preserve remaining CU")
// 	}
	
// 	return recommendations
// }

// // Helper functions for filtering and utilities
// func filterNonZeroBalances(balances []WalletBalance) []WalletBalance {
// 	var filtered []WalletBalance
// 	for _, balance := range balances {
// 		if !balance.IsZero {
// 			filtered = append(filtered, balance)
// 		}
// 	}
// 	return filtered
// }

// func filterByNetworks(balances []WalletBalance, networks []string) []WalletBalance {
// 	networkSet := make(map[string]bool)
// 	for _, network := range networks {
// 		networkSet[strings.ToLower(network)] = true
// 	}
	
// 	var filtered []WalletBalance
// 	for _, balance := range balances {
// 		if networkSet[strings.ToLower(balance.Network)] {
// 			filtered = append(filtered, balance)
// 		}
// 	}
// 	return filtered
// }

// func filterByTokens(balances []WalletBalance, tokens []string) []WalletBalance {
// 	tokenSet := make(map[string]bool)
// 	for _, token := range tokens {
// 		tokenSet[strings.ToUpper(token)] = true
// 	}
	
// 	var filtered []WalletBalance
// 	for _, balance := range balances {
// 		if tokenSet[strings.ToUpper(balance.Token)] {
// 			filtered = append(filtered, balance)
// 		}
// 	}
// 	return filtered
// }

// func getNetworksByType(isTestnet bool) []Network {
// 	var networks []Network
// 	for _, network := range allNetworks {
// 		if network.IsTestnet == isTestnet {
// 			networks = append(networks, network)
// 		}
// 	}
// 	return networks
// }

// func getNetworksByPriority(priority int) []Network {
// 	var networks []Network
// 	for _, network := range allNetworks {
// 		if network.Priority == priority {
// 			networks = append(networks, network)
// 		}
// 	}
// 	return networks
// }

// func getSupportedNetworks() []Network {
// 	var networks []Network
// 	for _, network := range allNetworks {
// 		if network.IsSupported {
// 			networks = append(networks, network)
// 		}
// 	}
// 	return networks
// }

// func getNetworkType(isTestnet bool) string {
// 	if isTestnet {
// 		return "testnet"
// 	}
// 	return "mainnet"
// }

// func getTokenLogo(moralisLogo, tokenAddress, symbol string) string {
// 	if moralisLogo != "" {
// 		return moralisLogo
// 	}
// 	return getFallbackLogo(tokenAddress, symbol)
// }

// func getFallbackLogo(tokenAddress, symbol string) string {
// 	// Fallback logo sources
// 	logoMap := map[string]string{
// 		"USDT": "https://cryptologos.cc/logos/tether-usdt-logo.png",
// 		"USDC": "https://cryptologos.cc/logos/usd-coin-usdc-logo.png",
// 		"DAI":  "https://cryptologos.cc/logos/multi-collateral-dai-dai-logo.png",
// 		"USDD": "https://cryptologos.cc/logos/usdd-usdd-logo.png",
// 		"LUSD": "https://cryptologos.cc/logos/liquity-usd-lusd-logo.png",
// 		"EURT": "https://cryptologos.cc/logos/tether-eurt-eurt-logo.png",
// 	}
	
// 	if logo, exists := logoMap[symbol]; exists {
// 		return logo
// 	}
	
// 	return fmt.Sprintf("https://via.placeholder.com/64x64/007bff/ffffff?text=%s", symbol)
// }

// func getDefaultDecimals(tokenSymbol string) int {
// 	switch tokenSymbol {
// 	case "USDT", "USDC", "EURT":
// 		return 6
// 	case "DAI", "USDD", "LUSD":
// 		return 18
// 	default:
// 		return 18
// 	}
// }

// func formatBalance(balance string, decimals int) string {
// 	if balance == "0" || balance == "" {
// 		return "0.00"
// 	}
	
// 	balanceBig := new(big.Int)
// 	balanceBig.SetString(balance, 10)
	
// 	if balanceBig.Cmp(big.NewInt(0)) == 0 {
// 		return "0.00"
// 	}
	
// 	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
// 	balanceFloat := new(big.Float).SetInt(balanceBig)
// 	divisorFloat := new(big.Float).SetInt(divisor)
// 	result := new(big.Float).Quo(balanceFloat, divisorFloat)
	
// 	formatted, _ := result.Float64()
// 	return fmt.Sprintf("%.2f", formatted)
// }

// func calculateUSDValue(balance string, decimals int) float64 {
// 	if balance == "0" || balance == "" {
// 		return 0.0
// 	}
	
// 	balanceBig := new(big.Int)
// 	balanceBig.SetString(balance, 10)
	
// 	if balanceBig.Cmp(big.NewInt(0)) == 0 {
// 		return 0.0
// 	}
	
// 	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)
// 	balanceFloat := new(big.Float).SetInt(balanceBig)
// 	divisorFloat := new(big.Float).SetInt(divisor)
// 	result := new(big.Float).Quo(balanceFloat, divisorFloat)
	
// 	// Assume stablecoins are pegged to $1 USD
// 	usdValue, _ := result.Float64()
// 	return usdValue
// }

// // Placeholder types and functions
// type TokenBalance struct {
// 	TokenAddress string `json:"token_address"`
// 	Name         string `json:"name"`
// 	Symbol       string `json:"symbol"`
// 	Logo         string `json:"logo"`
// 	Balance      string `json:"balance"`
// 	Decimals     int    `json:"decimals"`
// }

// type BalanceResponse []TokenBalance

// func makeAPIRequest(url string) (*BalanceResponse, error) {
// 	client := &http.Client{Timeout: 30 * time.Second}
	
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return nil, err
// 	}

// 	req.Header.Add("Accept", "application/json")
// 	req.Header.Add("X-API-Key", MORALIS_API_KEY)

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode != 200 {
// 		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
// 	}

// 	var balanceResp BalanceResponse
// 	err = json.Unmarshal(body, &balanceResp)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
// 	}

// 	return &balanceResp, nil
// }

// // Example usage and testing
// func main() {
// 	fmt.Println("🚀 Complete Portfolio Tracker with Smart CU Management")
// 	fmt.Printf("📊 Daily Budget: %d CU (%d calls max)\n", DAILY_CU_LIMIT, MAX_DAILY_CALLS)
// 	fmt.Printf("💰 Cost per call: %d CU (REAL usage-based)\n", CU_PER_CALL)
// 	fmt.Printf("🌐 Networks supported: %d total (%d mainnet, %d testnet)\n", 
// 		len(allNetworks), len(getNetworksByType(false)), len(getNetworksByType(true)))
// 	fmt.Printf("💡 FREE PLAN CAPACITY: ~500-2000 users/day with caching\n")
// 	fmt.Println()

// 	// Example wallet address
// 	walletAddress := "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"
// 	userID := "user123"

// 	// Get complete portfolio (shows ALL networks and tokens)
// 	fmt.Println("🔍 Getting complete portfolio...")
// 	portfolio := GetCompletePortfolio(walletAddress, userID)
	
// 	fmt.Printf("✅ Portfolio loaded for %s:\n", walletAddress)
// 	fmt.Printf("📊 Total Networks: %d\n", portfolio.TotalNetworks)
// 	fmt.Printf("🪙 Total Tokens: %d\n", portfolio.TotalTokens)
// 	fmt.Printf("💰 Total Value: $%.2f USD\n", portfolio.Summary.GrandTotalUSD)
// 	fmt.Printf("⚡ CU Used: %d\n", portfolio.Metadata.CUUsed)
// 	fmt.Printf("📈 Cache Hit Rate: %.1f%%\n", portfolio.Metadata.CacheHitRate)
// 	fmt.Printf("⏱️ Processing Time: %s\n", portfolio.Metadata.ProcessingTime)
// 	fmt.Println()

// 	// Show sample balances
// 	fmt.Println("🏦 MAINNET BALANCES (showing first 5):")
// 	for i, balance := range portfolio.MainnetBalances {
// 		if i >= 5 { break }
// 		status := "💰"
// 		if balance.IsZero { status = "⭕" }
// 		fmt.Printf("  %s %s - %s: %s ($%.2f) [%s]\n", 
// 			status, balance.Network, balance.Token, 
// 			balance.FormattedBalance, balance.USDValue, balance.DataSource)
// 	}
	
// 	fmt.Println("\n🧪 TESTNET BALANCES (showing first 3):")
// 	for i, balance := range portfolio.TestnetBalances {
// 		if i >= 3 { break }
// 		status := "💰"
// 		if balance.IsZero { status = "⭕" }
// 		fmt.Printf("  %s %s - %s: %s ($%.2f) [%s]\n", 
// 			status, balance.Network, balance.Token, 
// 			balance.FormattedBalance, balance.USDValue, balance.DataSource)
// 	}

// 	// Example: Get filtered portfolio (only non-zero balances)
// 	fmt.Println("\n🎯 Getting filtered portfolio (non-zero only)...")
// 	filteredPortfolio := GetPortfolioWithFilters(walletAddress, userID, PortfolioOptions{
// 		ShowOnlyNonZero: true,
// 		IncludeMetadata: true,
// 	})
	
// 	nonZeroCount := len(filteredPortfolio.MainnetBalances) + len(filteredPortfolio.TestnetBalances)
// 	fmt.Printf("💎 Non-zero balances found: %d\n", nonZeroCount)
	
// 	if nonZeroCount == 0 {
// 		fmt.Println("  📝 This wallet appears to have no stablecoin balances")
// 		fmt.Println("  💡 All balances are zero or the wallet is new")
// 	} else {
// 		for _, balance := range filteredPortfolio.MainnetBalances {
// 			fmt.Printf("  💰 %s - %s: %s ($%.2f)\n", 
// 				balance.Network, balance.Token, balance.FormattedBalance, balance.USDValue)
// 		}
// 		for _, balance := range filteredPortfolio.TestnetBalances {
// 			fmt.Printf("  🧪 %s - %s: %s ($%.2f)\n", 
// 				balance.Network, balance.Token, balance.FormattedBalance, balance.USDValue)
// 		}
// 	}

// 	// Show system status (simplified to avoid hanging)
// 	fmt.Println("\n📊 SYSTEM STATUS:")
// 	status := GetSystemStatus()
	
// 	// Print status in a safe way to avoid JSON marshaling issues
// 	fmt.Printf("CU Usage:\n")
// 	if cuUsage, ok := status["cu_usage"].(map[string]interface{}); ok {
// 		fmt.Printf("  Daily CU Used: %v\n", cuUsage["daily_cu_used"])
// 		fmt.Printf("  Remaining CU: %v\n", cuUsage["remaining_cu"])
// 		fmt.Printf("  Usage Percentage: %.1f%%\n", cuUsage["usage_percentage"])
// 		fmt.Printf("  Emergency Mode: %v\n", cuUsage["emergency_mode"])
// 	}
	
// 	fmt.Printf("\nCache Stats:\n")
// 	if cacheStats, ok := status["cache_stats"].(map[string]interface{}); ok {
// 		fmt.Printf("  Fresh Cache: %v entries\n", cacheStats["fresh_cache_size"])
// 		fmt.Printf("  Stale Cache: %v entries\n", cacheStats["stale_cache_size"])
// 		fmt.Printf("  Total Cached: %v entries\n", cacheStats["total_cached"])
// 	}
	
// 	fmt.Printf("\nCurrent Strategy:\n")
// 	if strategy, ok := status["current_strategy"].(map[string]interface{}); ok {
// 		fmt.Printf("  Priority 1 Networks: %v\n", strategy["priority_1_enabled"])
// 		fmt.Printf("  Priority 2 Networks: %v\n", strategy["priority_2_enabled"])
// 		fmt.Printf("  Emergency Mode: %v\n", strategy["emergency_mode"])
// 	}
	
// 	fmt.Printf("\nRecommendations:\n")
// 	if recommendations, ok := status["recommendations"].([]string); ok {
// 		for _, rec := range recommendations {
// 			fmt.Printf("  %s\n", rec)
// 		}
// 	}
	
// 	fmt.Println("\n✅ System analysis complete!")
// 	fmt.Printf("💡 Total portfolio value: $%.2f across %d networks\n", 
// 		portfolio.Summary.GrandTotalUSD, portfolio.TotalNetworks)
// 	fmt.Printf("🎯 Optimization: %.1f%% cache efficiency, %d CU saved\n", 
// 		portfolio.Metadata.CacheHitRate, DAILY_CU_LIMIT - portfolio.Metadata.CUUsed)
// }