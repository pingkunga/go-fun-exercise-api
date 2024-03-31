package postgres

import (
	"database/sql"
	"time"

	"github.com/KKGo-Software-engineering/fun-exercise-api/wallet"
)

type Wallet struct {
	ID         int       `postgres:"id"`
	UserID     int       `postgres:"user_id"`
	UserName   string    `postgres:"user_name"`
	WalletName string    `postgres:"wallet_name"`
	WalletType string    `postgres:"wallet_type"`
	Balance    float64   `postgres:"balance"`
	CreatedAt  time.Time `postgres:"created_at"`
}

func (p *Postgres) Wallets(walletType string) ([]wallet.Wallet, error) {

	var rows *sql.Rows
	var err error
	if walletType != "" {
		rows, err = p.Db.Query("SELECT * FROM user_wallet WHERE wallet_type = $1", walletType)
	} else {
		rows, err = p.Db.Query("SELECT * FROM user_wallet")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []wallet.Wallet
	for rows.Next() {
		var w Wallet
		err := rows.Scan(&w.ID,
			&w.UserID, &w.UserName,
			&w.WalletName, &w.WalletType,
			&w.Balance, &w.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet.Wallet{
			ID:         w.ID,
			UserID:     w.UserID,
			UserName:   w.UserName,
			WalletName: w.WalletName,
			WalletType: w.WalletType,
			Balance:    w.Balance,
			CreatedAt:  w.CreatedAt,
		})
	}
	return wallets, nil
}

func (p *Postgres) CreateWallet(wallet wallet.Wallet) error {
	_, err := p.Db.Exec("INSERT INTO user_wallet (user_id, user_name, wallet_name, wallet_type, balance) VALUES ($1, $2, $3, $4, $5)",
		wallet.UserID, wallet.UserName, wallet.WalletName, wallet.WalletType, wallet.Balance,
	)
	if err != nil {
		return err
	}
	return nil
}

func (p *Postgres) WalletByUserId(userId string) ([]wallet.Wallet, error) {
	rows, err := p.Db.Query("SELECT * FROM user_wallet WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []wallet.Wallet
	for rows.Next() {
		var w Wallet
		err := rows.Scan(&w.ID,
			&w.UserID, &w.UserName,
			&w.WalletName, &w.WalletType,
			&w.Balance, &w.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, wallet.Wallet{
			ID:         w.ID,
			UserID:     w.UserID,
			UserName:   w.UserName,
			WalletName: w.WalletName,
			WalletType: w.WalletType,
			Balance:    w.Balance,
			CreatedAt:  w.CreatedAt,
		})
	}
	return wallets, nil
}
