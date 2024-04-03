//go:build integration

package wallet

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Response is a wrapper around http.Response that provides a way to check for
type Response struct {
	*http.Response
	err error
}

func clientRequest(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	//req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}
	defer r.Body.Close()

	//NewDecoder == json.unmarshal
	//เอา r.Body มา decode แล้วเก็บใน v >>
	// - v เป็น Struct ในที่นี้ user
	return json.NewDecoder(r.Body).Decode(v)
}

func uri(paths ...string) string {
	baseURL := os.Getenv("TEST_URL")

	if baseURL == "" {
		baseURL = "http://localhost:1323/api/v1"
	}

	if paths == nil {
		return baseURL
	}
	return baseURL + "/" + strings.Join(paths, "/")
}

// ================================================================
// Test Method Section
// ================================================================
func TestITGetWallets(t *testing.T) {
	//Arrange
	var result []Wallet

	//Act
	res := clientRequest(http.MethodGet, uri("wallets"), nil)
	err := res.Decode(&result)

	//Assert
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	//ดูว่ามีของคืนมาไหม ไม่สนใจว่า Value ถูกไหม
	assert.Greater(t, len(result), 0)
}

func TestITGetWalletByWallerType(t *testing.T) {
	//Arrange
	var result []Wallet

	//Act
	res := clientRequest(http.MethodGet, uri("wallets?wallet_type=Savings"), nil)
	err := res.Decode(&result)

	//Assert
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)

	//Initial Data มี 2 ตัว
	assert.Equal(t, len(result), 2)
}

func TestITCreateWallet(t *testing.T) {
	//Arrange
	wallet := Wallet{
		UserID:     100,
		UserName:   "PingkungA",
		WalletName: "PingkungA Wallet",
		WalletType: "Savings",
		Balance:    1000,
	}

	//Act
	res := clientRequest(http.MethodPost, uri("wallets"), strings.NewReader(`{
		"user_id": 100,
		"user_name": "PingkungA",
		"wallet_name": "PingkungA Wallet",
		"wallet_type": "Savings",
		"balance": 1000
	}`))
	var result Wallet
	err := res.Decode(&result)

	//Assert
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, wallet.UserID, result.UserID)
	assert.Equal(t, wallet.UserName, result.UserName)
	assert.Equal(t, wallet.WalletName, result.WalletName)
	assert.Equal(t, wallet.WalletType, result.WalletType)
	assert.Equal(t, wallet.Balance, result.Balance)
}

func TestITUpdateWallet(t *testing.T) {
	//Arrange
	wallet := seedWallet(t)
	wallet.Balance = 2000

	//Act
	res := clientRequest(http.MethodPut, uri("wallets"), strings.NewReader(`{
		"user_id": 190,
		"user_name": "PingkungB",
		"wallet_name": "PingkungB Wallet",
		"wallet_type": "Savings",
		"balance": 2000
	}`))
	var result Wallet
	err := res.Decode(&result)

	//Assert
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, wallet.UserID, result.UserID)
	assert.Equal(t, wallet.UserName, result.UserName)
	assert.Equal(t, wallet.WalletName, result.WalletName)
	assert.Equal(t, wallet.WalletType, result.WalletType)
	assert.Equal(t, wallet.Balance, result.Balance)
}

func TestITDeleteWalletByUserID(t *testing.T) {
	//Arrange UserId = 190
	seedWallet(t)

	//Act
	res := clientRequest(http.MethodDelete, uri("users/190/wallets"), nil)

	//Assert
	assert.EqualValues(t, http.StatusNoContent, res.StatusCode)
}

func seedWallet(t *testing.T) Wallet {
	var walletEntry Wallet
	body := bytes.NewBufferString(`{
		"user_id": 190,
		"user_name": "PingkungB",
		"wallet_name": "PingkungB Wallet",
		"wallet_type": "Savings",
		"balance": 1000
	}`)
	err := clientRequest(http.MethodPost, uri("wallets"), body).Decode(&walletEntry)
	if err != nil {
		t.Fatal("can't create uomer:", err)
	}
	return walletEntry
}

func TestGetWalletByUserID(t *testing.T) {
	//Arrange
	var result []Wallet

	//Act
	res := clientRequest(http.MethodGet, uri("users/1/wallets"), nil)
	err := res.Decode(&result)

	//Assert
	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)

	//Initial Data มี 3 ตัว
	assert.Equal(t, len(result), 3)
}
