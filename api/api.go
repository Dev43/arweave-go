package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/Dev43/arweave-go/tx"
)

var v2BlockFormatHeader = map[string]string{"X-Block-Format": "2"}


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
func (c *Client) GetData(ctx context.Context, txID string) (string, error) {
	body, err := c.get(ctx, fmt.Sprintf("tx/%s/data", txID), nil)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// LastTransaction requests the last transaction of an account
func (c *Client) LastTransaction(ctx context.Context, address string) (string, error) {
	body, err := c.get(ctx, fmt.Sprintf("wallet/%s/last_tx", address), nil)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetTransaction requests the information of a transaction
func (c *Client) GetTransaction(ctx context.Context, txID string) (*tx.Transaction, error) {
	body, err := c.get(ctx, fmt.Sprintf("tx/%s", txID), nil)
	if err != nil {
		return nil, err
	}
	// If it sends us a pending message, return a nil receipt and error
	if string(body) == "Pending" {
		return nil, nil
	}
	tx := tx.Transaction{}
	err = json.Unmarshal(body, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// GetTransaction requests the information of a transaction
func (c *Client) GetPendingTransactions(ctx context.Context) ([]string, error) {
	body, err := c.get(ctx, fmt.Sprintf("tx/pending"), nil)
	if err != nil {
		return nil, err
	}
	var tran []string
	err = json.Unmarshal(body, &tran)
	if err != nil {
		return nil, err
	}
	return tran, nil
}

// GetTransactionField requests the specific field of a specific transaction
func (c *Client) GetTransactionField(ctx context.Context, txID string, field string) (string, error) {
	_, ok := allowedFields[field]
	if !ok {
		return "", errors.New("field does not exist")
	}
	body, err := c.get(ctx, fmt.Sprintf("tx/%s/%s", txID, field), nil)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// GetBlockByID requests a block by its id
func (c *Client) GetBlockByID(ctx context.Context, blockID string) (*Block, error) {
	body, err := c.get(ctx, fmt.Sprintf("block/hash/%s", blockID), v2BlockFormatHeader)
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
func (c *Client) GetBlockByHeight(ctx context.Context, height int64) (*Block, error) {
	body, err := c.get(ctx, fmt.Sprintf("block/height/%d", height), v2BlockFormatHeader)
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
func (c *Client) GetCurrentBlock(ctx context.Context) (*Block, error) {
	body, err := c.get(ctx, "current_block", v2BlockFormatHeader)
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
func (c *Client) GetReward(ctx context.Context, data []byte) (string, error) {
	body, err := c.get(ctx, fmt.Sprintf("price/%d", len(data)), nil)
	if err != nil {
		return "", err
	}
	return string(body), nil

}

// GetBalance requests the current balance of an address
func (c *Client) GetBalance(ctx context.Context, address string) (string, error) {
	body, err := c.get(ctx, fmt.Sprintf("wallet/%s/balance", address), nil)
	if err != nil {
		return "", err
	}
	return string(body), nil

}

// GetPeers requests the list of peers of a node
func (c *Client) GetPeers(ctx context.Context) ([]string, error) {
	body, err := c.get(ctx, "peers", nil)
	if err != nil {
		return nil, err
	}
	var peers []string
	err = json.Unmarshal(body, &peers)
	if err != nil {
		return nil, err
	}

	return peers, nil

}

// GetInfo requests the information of a node
func (c *Client) GetInfo(ctx context.Context) (*NetworkInfo, error) {
	body, err := c.get(ctx, "info",nil)
	if err != nil {
		return nil, err
	}
	info := NetworkInfo{}
	err = json.Unmarshal(body, &info)
	if err != nil {
		return nil, err
	}

	return &info, nil
}

// Commit sends a transaction to the weave with a context
func (c *Client) Commit(ctx context.Context, data []byte) (string, error) {
	body, err := c.post(ctx, "tx", data, nil)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func getResponse(resp io.ReadCloser, returnedError error) ([]byte, error) {
	if resp != nil {
		defer resp.Close()
	}
	if returnedError != nil {
		return handleHTTPError(resp, returnedError)
	}

	b, err := ioutil.ReadAll(resp)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func handleHTTPError(resp io.Reader, returnedError error) ([]byte, error) {
	if resp != nil {
		buf := new(bytes.Buffer)
		if _, err := buf.ReadFrom(resp); err == nil {
			return nil, fmt.Errorf("%v %v", returnedError, buf.String())
		}
	}
	return nil, returnedError
}

func (c *Client) requestWithContext(ctx context.Context, method string, url string, body []byte, headers map[string]string) (io.ReadCloser, error) {
	req, err := http.NewRequest(method, url, ioutil.NopCloser(bytes.NewReader(body)))
	if err != nil {
		return nil, err
	}
	reqWithContext := req.WithContext(ctx)
	reqWithContext.ContentLength = int64(len(body))

	if len(headers) > 0 {
		for k, v := range headers {
			reqWithContext.Header.Set(k, v)
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return resp.Body, errors.New(resp.Status)
	}
	return resp.Body, nil
}

func (c *Client) post(ctx context.Context, endpoint string, body []byte, headers map[string]string) ([]byte, error) {
	headers["Content-type"] = "application/json"
	resp, err := c.requestWithContext(ctx, "POST", c.formatURL(endpoint), body, headers)
	return getResponse(resp, err)
}

func (c *Client) get(ctx context.Context, endpoint string, headers map[string]string) ([]byte, error) {
	resp, err := c.requestWithContext(ctx, "GET", c.formatURL(endpoint), nil, headers)
	return getResponse(resp, err)
}

func (c *Client) formatURL(endpoint string) string {
	return fmt.Sprintf("%s/%s", c.url, endpoint)
}
