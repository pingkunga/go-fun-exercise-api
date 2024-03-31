package wallet

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestWallet(t *testing.T) {
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
				WalletType: "A",
				Balance:    99999,
			},
			{
				ID:         2,
				UserID:     2,
				UserName:   "pingkungb",
				WalletName: "pingkungb_wallet",
				WalletType: "B",
				Balance:    99999,
			},
		}
		stubWallet := StubWallet{Wallet: expected}
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
		if len(got) == len(expected) {
			t.Errorf("expected list of wallets 2, got empty")
		}

		if !reflect.DeepEqual(expected, got) {
			t.Errorf("expected list of wallets %v, got %v", expected, got)
		}
	})
}

// Struct from postgres/wallet.go
type StubWallet struct {
	Wallet []Wallet
	err    error
}

func (s StubWallet) Wallets() ([]Wallet, error) {
	return s.Wallet, s.err
}
