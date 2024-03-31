package wallet

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

// for implement interface in wqllet.go
type Storer interface {
	Wallets(walletType string) ([]Wallet, error)
	CreateWallet(wallet Wallet) error
	WalletByUserId(userId string) ([]Wallet, error)
}

func New(db Storer) *Handler {
	return &Handler{store: db}
}

type Err struct {
	Message string `json:"message"`
}

// WalletHandler
//
//		@Summary		Get all wallets
//		@Description	Get all wallets
//		@Tags			wallet
//	 @Param			wallet_type	query	string	false	"wallet type" Enums(Savings, Credit Card, Crypto Wallet)
//		@Accept			json
//		@Produce		json
//		@Success		200	{object}	Wallet
//		@Router			/api/v1/wallets [get]
//		@Failure		500	{object}	Err
func (h *Handler) WalletHandler(c echo.Context) error {
	walletType := c.QueryParam("wallet_type")
	wallets, err := h.store.Wallets(walletType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallets)
}

func (h *Handler) CreateWalletHandler(c echo.Context) error {
	var wallet Wallet
	if err := c.Bind(&wallet); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}
	if err := h.store.CreateWallet(wallet); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, wallet)
}

// WalletByUserIdHandler
//
//	@Summary		Get wallet by user id
//	@Description	Get wallet by user id
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"User ID"
//	@Success		200	{object}	Wallet
//	@Router			/api/v1/users/{id}/wallets [get]
//	@Failure		500	{object}	Err

func (h *Handler) WalletByUserIdHandler(c echo.Context) error {
	id := c.Param("id")
	wallet, err := h.store.WalletByUserId(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallet)
}
