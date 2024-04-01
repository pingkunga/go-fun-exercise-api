package wallet

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

// Test handler.go
func TestGetWallet(t *testing.T) {
	t.Run("given unable to get wallets should return 500 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubError := StubWallet{err: errors.New("unable to get wallets")}
		handler := New(stubError)
		err := handler.WalletHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		if (rec.Code != http.StatusInternalServerError) && (rec.Body.String() != `{"message":"unable to get wallets"}`) {
			t.Errorf("expected 500 and error message, got %d and %s", rec.Code, rec.Body.String())
		}

	})

	t.Run("given user able to getting wallet should return list of wallets", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallets", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := []Wallet{
			{
				ID:         1,
				UserID:     1,
				UserName:   "pingkunga",
				WalletName: "pingkunga_wallet",
				WalletType: "Savings",
				Balance:    99999,
			},
			{
				ID:         2,
				UserID:     2,
				UserName:   "pingkungb",
				WalletName: "pingkungb_wallet",
				WalletType: "Savings",
				Balance:    99999,
			},
		}
		stubWallet := StubWallet{wallet: expected}
		handler := New(stubWallet)
		err := handler.WalletHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		//Check if the response is 200
		actualBody := rec.Body.String()
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d and %s", rec.Code, actualBody)
		}

		//Convert response to struct wallet
		var got []Wallet
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("expected list of wallets, got %s", actualBody)
		}

		//Check if the response is the same as the expected
		if len(got) != len(expected) {
			t.Errorf("expected list of wallets 2, got empty")
		}

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected list of wallets %v, got %v", expected, got)
		}
	})

}

func TestCreateWallet(t *testing.T) {
	t.Run("given unable to create wallet should return 500 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallets", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubError := StubWallet{err: errors.New("unable to create wallet")}
		handler := New(stubError)
		err := handler.CreateWalletHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		if (rec.Code != http.StatusInternalServerError) && (rec.Body.String() != `{"message":"unable to create wallet"}`) {
			t.Errorf("expected 500 and error message, got %d and %s", rec.Code, rec.Body.String())
		}

	})

	t.Run("given user able to create wallet should return wallet created", func(t *testing.T) {
		e := echo.New()

		//https://stackoverflow.com/questions/76197311/unit-test-for-post-request-is-not-working-in-go
		req := httptest.NewRequest(http.MethodPost, "/api/v1/wallets", io.NopCloser(strings.NewReader(`{"user_id": 1, "user_name": "pingkunga", "wallet_name": "pingkunga_wallet", "wallet_type": "Savings", "Balance": 99999}`)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		//ResponseWriter interface and records all the responses from the handler:
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := Wallet{
			UserID:     1,
			UserName:   "pingkunga",
			WalletName: "pingkunga_wallet",
			WalletType: "Savings",
			Balance:    99999,
		}
		stubWallet := StubWallet{createWallet: expected}
		handler := New(stubWallet)
		err := handler.CreateWalletHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		//Check if the response is 201
		actualBody := rec.Body.String()
		if rec.Code != http.StatusCreated {
			t.Errorf("expected 201, got %d and %s", rec.Code, actualBody)
		}

		//Convert response to struct wallet
		var got Wallet
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("expected wallet created, got %s", actualBody)
		}

		//Check if the response is the same as the expected
		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected wallet created %v, got %v", expected, got)
		}
	})
}

func TestUpdateWallet(t *testing.T) {
	t.Run("given unable to update wallet should return 500 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/wallets", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubError := StubWallet{err: errors.New("unable to update wallet")}
		handler := New(stubError)
		err := handler.UpdateWalletHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		if (rec.Code != http.StatusInternalServerError) && (rec.Body.String() != `{"message":"unable to update wallet"}`) {
			t.Errorf("expected 500 and error message, got %d and %s", rec.Code, rec.Body.String())
		}

	})

	t.Run("given user able to update wallet should return wallet updated", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPut, "/api/v1/wallets", io.NopCloser(strings.NewReader(`{"user_id": 1, "user_name": "pingkunga_updated", "wallet_name": "pingkunga_wallet", "wallet_type": "Savings", "Balance": 99999}`)))
		//req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		req.Header.Add("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := Wallet{
			UserID:     1,
			UserName:   "pingkunga_updated",
			WalletName: "pingkunga_wallet",
			WalletType: "Savings",
			Balance:    99999,
		}
		stubWallet := StubWallet{updateWallet: expected}
		handler := New(stubWallet)
		err := handler.UpdateWalletHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		//Check if the response is 200
		actualBody := rec.Body.String()
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d and %s", rec.Code, actualBody)
		}

		//Convert response to struct wallet
		var got Wallet
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("expected wallet updated, got %s", actualBody)
		}

		//Check if the response is the same as the expected
		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected wallet updated %v, got %v", expected, got)
		}
	})
}

func TestDeleteWalletByUserId(t *testing.T) {
	t.Run("given unable to delete wallet should return 500 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/99/wallets", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubError := StubWallet{err: errors.New("unable to delete wallet")}
		handler := New(stubError)
		err := handler.DeleteWalletByUserIdHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		if (rec.Code != http.StatusInternalServerError) && (rec.Body.String() != `{"message":"unable to delete wallet"}`) {
			t.Errorf("expected 500 and error message, got %d and %s", rec.Code, rec.Body.String())
		}

	})

	t.Run("given user able to delete wallet should return wallet deleted", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/99/wallets", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubWallet := StubWallet{deleteWallet: "wallet deleted"}
		handler := New(stubWallet)
		err := handler.DeleteWalletByUserIdHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		//Check if the response is 200
		actualBody := rec.Body.String()
		if rec.Code != http.StatusNoContent {
			t.Errorf("expected 204, got %d and %s", rec.Code, actualBody)
		}
	})

}

func TestGetWalletByUserId(t *testing.T) {
	t.Run("given unable to get wallet by user id should return 500 and error message", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/99/wallets", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		stubError := StubWallet{err: errors.New("unable to get wallet by user id")}
		handler := New(stubError)
		err := handler.WalletByUserIdHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		if (rec.Code != http.StatusInternalServerError) && (rec.Body.String() != `{"message":"unable to get wallet by user id"}`) {
			t.Errorf("expected 500 and error message, got %d and %s", rec.Code, rec.Body.String())
		}

	})

	t.Run("given user able to get wallet by user id should return list of wallets", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/99/wallets", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		expected := []Wallet{
			{
				ID:         1,
				UserID:     99,
				UserName:   "pingkunga",
				WalletName: "pingkunga_wallet",
				WalletType: "Savings",
				Balance:    99999,
			},
			{
				ID:         2,
				UserID:     99,
				UserName:   "pingkungb",
				WalletName: "pingkungb_wallet",
				WalletType: "Savings",
				Balance:    99999,
			},
		}

		stubWallet := StubWallet{wallet: expected}
		handler := New(stubWallet)
		err := handler.WalletByUserIdHandler(c)

		if err != nil {
			t.Errorf("got some error %v", err)
		}

		//Check if the response is 200
		actualBody := rec.Body.String()
		if rec.Code != http.StatusOK {
			t.Errorf("expected 200, got %d and %s", rec.Code, actualBody)
		}

		//Convert response to struct wallet
		var got []Wallet
		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
			t.Errorf("expected list of wallets, got %s", actualBody)
		}

		//Check if the response is the same as the expected
		if len(got) != len(expected) {
			t.Errorf("expected list of wallets 2, got empty")
		}

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected list of wallets %v, got %v", expected, got)
		}
	})
}

// Struct from postgres/wallet.go
type StubWallet struct {
	wallet       []Wallet
	createWallet Wallet
	updateWallet Wallet
	deleteWallet string
	err          error
}

// ล้อกับ type Storer interface in handler.go
func (s StubWallet) Wallets(walletType string) ([]Wallet, error) {
	return s.wallet, s.err
}

func (s StubWallet) CreateWallet(wallet Wallet) error {
	return s.err
}

func (s StubWallet) UpdateWallet(wallet Wallet) error {
	return s.err
}

func (s StubWallet) DeleteWalletByUserId(userId string) error {
	return s.err
}

func (s StubWallet) WalletByUserId(userId string) ([]Wallet, error) {
	return s.wallet, s.err
}
