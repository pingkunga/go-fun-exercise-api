package wallet

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	store Storer
}

// for implement interface in wallet.go
type Storer interface {
	Wallets(walletType string) ([]Wallet, error)
	CreateWallet(wallet Wallet) error
	UpdateWallet(wallet Wallet) error
	DeleteWalletByUserId(userId string) error
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
//	    @Param			wallet_type	query	string	false	"wallet type" Enums(Savings, Credit Card, Crypto Wallet)
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

// CreateWalletHandler
//
//	@Summary		Create wallet
//	@Description	Create wallet
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Param			wallet	body	Wallet	true	"Wallet object"
//	@Success		201	{object}	Wallet
//	@Router			/api/v1/wallets [post]
//	@Failure		400	{object}	Err
//	@Failure		500	{object}	Err
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

// UpdateWalletHandler
//
//	@Summary		Update wallet
//	@Description	Update wallet
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Param			wallet	body	Wallet	true	"Wallet object"
//	@Success		200	{object}	Wallet
//	@Router			/api/v1/wallets [put]
//	@Failure		400	{object}	Err
//	@Failure		500	{object}	Err
func (h *Handler) UpdateWalletHandler(c echo.Context) error {
	var wallet Wallet
	if err := c.Bind(&wallet); err != nil {
		return c.JSON(http.StatusBadRequest, Err{Message: err.Error()})
	}

	if err := h.store.UpdateWallet(wallet); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallet)
}

// DeleteWalletByUserIdHandler
//
//	@Summary		Delete wallet by user id
//	@Description	Delete wallet by user id
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"user id"
//	@Success		204	{object}	Err
//  @Failure		500	{object}	Err
//  @Router			/api/v1/users/:id/wallets [delete]

func (h *Handler) DeleteWalletByUserIdHandler(c echo.Context) error {
	id := c.Param("id")
	if err := h.store.DeleteWalletByUserId(id); err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// WalletByUserIdHandler
//
//	@Summary		Get wallet by user id
//	@Description	Get wallet by user id
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"user id"
//	@Success		200	{object}	Wallet
//	@Failure		500	{object}	Err
//	@Router			/api/v1/users/:id/wallets [get]

func (h *Handler) WalletByUserIdHandler(c echo.Context) error {
	id := c.Param("id")
	wallet, err := h.store.WalletByUserId(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallet)
}
