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
	Wallets() ([]Wallet, error)
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
//	@Summary		Get all wallets
//	@Description	Get all wallets
//	@Tags			wallet
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	Wallet
//	@Router			/api/v1/wallets [get]
//	@Failure		500	{object}	Err
func (h *Handler) WalletHandler(c echo.Context) error {
	wallets, err := h.store.Wallets()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Err{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, wallets)
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
