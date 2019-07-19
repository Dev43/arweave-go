package arweave


type Transactor interface {
	CreateTransaction(w *wallet.Wallet, amount string, data []byte, target string) (*tx.Transaction, error)
	SendTransaction(tx *tx.Transaction) (string, error)
	WaitMined(ctx context.Context, tx *tx.Transaction) (*tx.JSONTransaction, error)
}