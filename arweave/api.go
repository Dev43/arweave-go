package arweave

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client struct
type Client struct {
	client *http.Client
	url    string
}

// Dial creates a new arweave client
func Dial(url string) (*Client, error) {
	return &Client{client: new(http.Client), url: url}, nil
}

// GetData requests the data of a transaction
func (c *Client) GetData(txID string) (string, error) {
	body, err := c.get(fmt.Sprintf("tx/%s/data", txID))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// LastTransaction requests the last transaction of an account
func (c *Client) LastTransaction(address string) (string, error) {
	body, err := c.get(fmt.Sprintf("wallet/%s/last_tx", address))
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetTransaction requests the information of a transaction
func (c *Client) GetTransaction(txID string) (*JSONTransaction, error) {
	body, err := c.get(fmt.Sprintf("tx/%s", txID))
	if err != nil {
		return nil, err
	}
	tx := JSONTransaction{}
	err = json.Unmarshal(body, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// GetTransaction requests the information of a transaction
func (c *Client) GetPendingTransactions() ([]string, error) {
	body, err := c.get(fmt.Sprintf("tx/pending"))
	if err != nil {
		return nil, err
	}
	tx := []string{}
	err = json.Unmarshal(body, &tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// GetTransactionField requests the specific field of a specific transaction
func (c *Client) GetTransactionField(txID string, field string) (string, error) {
	_, ok := allowedFields[field]
	if !ok {
		return "", errors.New("field does not exist")
	}
	body, err := c.get(fmt.Sprintf("tx/%s/%s", txID, field))
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetBlockByID requests a block by its id
func (c *Client) GetBlockByID(blockID string) (*Block, error) {
	body, err := c.get(fmt.Sprintf("block/hash/%s", blockID))
	if err != nil {
		return nil, err
	}
	block := Block{}
	err = json.Unmarshal(body, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetBlockByHeight requests a block by its height
func (c *Client) GetBlockByHeight(height int64) (*Block, error) {
	body, err := c.get(fmt.Sprintf("block/height/%d", height))
	if err != nil {
		return nil, err
	}
	block := Block{}
	err = json.Unmarshal(body, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetCurrentBlock requests the latest block of the weave
func (c *Client) GetCurrentBlock() (*Block, error) {
	body, err := c.get("current_block")
	if err != nil {
		return nil, err
	}
	block := Block{}
	err = json.Unmarshal(body, &block)
	if err != nil {
		return nil, err
	}

	return &block, nil
}

// GetReward requests the current network reward
func (c *Client) GetReward(data []byte) (string, error) {
	body, err := c.get(fmt.Sprintf("price/%d", len(data)))
	if err != nil {
		return "", err
	}
	return string(body), nil

}

// GetBalance requests the current balance of an address
func (c *Client) GetBalance(address string) (string, error) {
	body, err := c.get(fmt.Sprintf("wallet/%s/balance", address))
	if err != nil {
		return "", err
	}
	return string(body), nil

}

// GetPeers requests the list of peers of a node
func (c *Client) GetPeers() ([]string, error) {
	body, err := c.get("peers")
	if err != nil {
		return nil, err
	}
	peers := []string{}
	err = json.Unmarshal(body, &peers)
	if err != nil {
		return nil, err
	}

	return peers, nil

}

// GetInfo requests the information of a node
func (c *Client) GetInfo() (*NetworkInfo, error) {
	body, err := c.get("info")
	if err != nil {
		return nil, err
	}
	info := NetworkInfo{}
	json.Unmarshal(body, &info)
	return &info, nil
}

// Commit sends the serialized transaction to the arweave
func (c *Client) Commit(data []byte) (string, error) {
	body, err := c.post(context.TODO(), "tx", data)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (c *Client) get(endpoint string) ([]byte, error) {
	resp, err := c.client.Get(fmt.Sprintf("%s/%s", c.url, endpoint))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("not found")
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api error with body %s", string(b))
	}
	return b, err
}

func (c *Client) post(ctx context.Context, endpoint string, body []byte) ([]byte, error) {
	r := bytes.NewReader(body)
	resp, err := c.client.Post(fmt.Sprintf("%s/%s", c.url, endpoint), "application/json", r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
