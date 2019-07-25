package tx

import "math/big"

// Transaction struct
type Transaction struct {
	id        []byte   // A SHA2-256 hash of the signature
	lastTx    string   // The ID of the last transaction made from the account. If no previous transactions have been made from the address this field is set to an empty string.
	owner     *big.Int // The modulus of the RSA key pair corresponding to the wallet making the transaction
	target    string   // If making a financial transaction this field contains the wallet address of the recipient base64url encoded. If the transaction is not a financial this field is set to an empty string.
	quantity  string   // If making a financial transaction this field contains the amount in Winston to be sent to the receiving wallet. If the transaction is not financial this field is set to the string "0". 1 AR = 1000000000000 (1e+12) Winston
	data      []byte   // If making an archiving transaction this field contains the data to be archived base64url encoded. If the transaction is not archival this field is set to an empty string.
	reward    string   // This field contains the mining reward for the transaction in Winston.
	tags      []Tag    // Transaction tags
	signature []byte   // Signature using the RSA-PSS signature scheme using SHA256 as the MGF1 masking algorithm
}

// Transaction encoded transaction to send to the arweave client
type transactionJSON struct {
	// Id A SHA2-256 hash of the signature, based 64 URL encoded.
	ID string `json:"id"`
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
	Tags []Tag `json:"tags"`
}

// Tag contains any tags the user wants to add to the transaction
type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
