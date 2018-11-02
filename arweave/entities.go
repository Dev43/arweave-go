package arweave

type NetworkInfo struct {
	Network          string `json:"network"`
	Version          int    `json:"version"`
	Release          int    `json:"release"`
	Height           int    `json:"height"`
	Current          string `json:"current"`
	Blocks           int    `json:"blocks"`
	Peers            int    `json:"peers"`
	QueueLength      int    `json:"queue_length"`
	NodeStateLatency int    `json:"node_state_latency"`
}

type Block struct {
	HashList      []string      `json:"hash_list"`
	Nonce         string        `json:"nonce"`
	PreviousBlock string        `json:"previous_block"`
	Timestamp     int           `json:"timestamp"`
	LastRetarget  int           `json:"last_retarget"`
	Diff          int           `json:"diff"`
	Height        int           `json:"height"`
	Hash          string        `json:"hash"`
	IndepHash     string        `json:"indep_hash"`
	Txs           []interface{} `json:"txs"`
	WalletList    []struct {
		Wallet   string `json:"wallet"`
		Quantity int64  `json:"quantity"`
		LastTx   string `json:"last_tx"`
	} `json:"wallet_list"`
	RewardAddr string        `json:"reward_addr"`
	Tags       []interface{} `json:"tags"`
	RewardPool int           `json:"reward_pool"`
	WeaveSize  int           `json:"weave_size"`
	BlockSize  int           `json:"block_size"`
}
