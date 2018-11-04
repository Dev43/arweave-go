package arweave

import "math/big"

// NetworkInfo struct
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

// Block struct
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

// Transaction struct
type Transaction struct {
	id        [32]byte                 // A SHA2-256 hash of the signature
	lastTx    string                   // The ID of the last transaction made from the account. If no previous transactions have been made from the address this field is set to an empty string.
	owner     *big.Int                 // The modulus of the RSA key pair corresponding to the wallet making the transaction
	target    string                   // If making a financial transaction this field contains the wallet address of the recipient base64url encoded. If the transaction is not a financial this field is set to an empty string.
	quantity  string                   // If making a financial transaction this field contains the amount in Winston to be sent to the receiving wallet. If the transaction is not financial this field is set to the string "0". 1 AR = 1000000000000 (1e+12) Winston
	data      string                   // If making an archiving transaction this field contains the data to be archived base64url encoded. If the transaction is not archival this field is set to an empty string.
	reward    string                   // This field contains the mining reward for the transaction in Winston.
	tags      []map[string]interface{} // Transaction tags
	signature []byte                   // Signature using the RSA-PSS signature scheme using SHA256 as the MGF1 masking algorithm
}

type JsonTransaction struct {
	// Id A SHA2-256 hash of the signature, based 64 URL encoded.
	Id string `json:"id"`
	// LastTx represents the ID of the last transaction made from the same address base64url encoded. If no previous transactions have been made from the address this field is set to an empty string.
	LastTx string `json:"last_tx"`
	//Owner is the modulus of the RSA key pair corresponding to the wallet making the transaction, base64url encoded.
	Owner string `json:"owner"`
	// Target if making a financial transaction this field contains the wallet address of the recipient base64url encoded. If the transaction is not a financial this field is set to an empty string.
	Target string `json:"target"`
	// Quantity If making a financial transaction this field contains the amount in Winston to be sent to the receiving wallet. If the transaction is not financial this field is set to the string "0". 1 AR = 1000000000000 (1e+12) Winston
	Quantity string `json:"quantity"`
	// Data If making an archiving transaction this field contains the data to be archived base64url encoded. If the transaction is not archival this field is set to an empty string.
	Data string `json:"data"`
	// Reward This field contains the mining reward for the transaction in Winston.
	Reward string `json:"reward"`
	//  Signature using the RSA-PSS signature scheme using SHA256 as the MGF1 masking algorithm
	Signature string `json:"signature"`
	// Tags Transaction tags
	Tags []map[string]interface{} `json:"tags"`
}

var allowedFields = map[string]bool{
	"id":        true,
	"last_tx":   true,
	"owner":     true,
	"target":    true,
	"quantity":  true,
	"type":      true,
	"data":      true,
	"reward":    true,
	"signature": true,
	"data.html": true,
}
