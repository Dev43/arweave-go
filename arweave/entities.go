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

//A wallet address is a base64url encoded SHA256 hash of the raw unencoded RSA modulus.
type Transaction struct {
	id        [32]byte                 // A SHA2-256 hash of the signature, based 64 URL encoded.
	lastTx    string                   // The ID of the last transaction made from the same address base64url encoded. If no previous transactions have been made from the address this field is set to an empty string.
	owner     string                   //The modulus of the RSA key pair corresponding to the wallet making the transaction, base64url encoded.
	target    string                   //  If making a financial transaction this field contains the wallet address of the recipient base64url encoded. If the transaction is not a financial this field is set to an empty string.
	quantity  string                   // If making a financial transaction this field contains the amount in Winston to be sent to the receiving wallet. If the transaction is not financial this field is set to the string "0". 1 AR = 1000000000000 (1e+12) Winston
	data      string                   //If making an archiving transaction this field contains the data to be archived base64url encoded. If the transaction is not archival this field is set to an empty string.
	reward    string                   //  This field contains the mining reward for the transaction in Winston.
	tags      []map[string]interface{} // The data for the signature is comprised of previous data from the rest of the transaction.
	signature []byte                   // The data for the signature is comprised of previous data from the rest of the transaction.
}

type JsonTransaction struct {
	Id        string                   `json:"id"`        // A SHA2-256 hash of the signature, based 64 URL encoded.
	LastTx    string                   `json:"last_tx"`   // The ID of the last transaction made from the same address base64url encoded. If no previous transactions have been made from the address this field is set to an empty string.
	Owner     string                   `json:"owner"`     //The modulus of the RSA key pair corresponding to the wallet making the transaction, base64url encoded.
	Target    string                   `json:"target"`    //  If making a financial transaction this field contains the wallet address of the recipient base64url encoded. If the transaction is not a financial this field is set to an empty string.
	Quantity  string                   `json:"quantity"`  // If making a financial transaction this field contains the amount in Winston to be sent to the receiving wallet. If the transaction is not financial this field is set to the string "0". 1 AR = 1000000000000 (1e+12) Winston
	Data      string                   `json:"data"`      //If making an archiving transaction this field contains the data to be archived base64url encoded. If the transaction is not archival this field is set to an empty string.
	Reward    string                   `json:"reward"`    //  This field contains the mining reward for the transaction in Winston.
	Signature string                   `json:"signature"` // The data for the signature is comprised of previous data from the rest of the transaction.
	Tags      []map[string]interface{} `json:"tags"`      // The data for the signature is comprised of previous data from the rest of the transaction.
}
