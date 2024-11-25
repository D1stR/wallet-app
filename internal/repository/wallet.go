package repository

import (
	"WalletApp/internal/domain"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"log"
)

type WalletRepository interface {
	UpdateWalletBalance(tx *sql.Tx, walletID uuid.UUID, operationType domain.OperationType, amount float64) error
	GetWalletBalance(walletID uuid.UUID) (float64, error)
	GetWalletByID(walletID uuid.UUID) (*domain.Wallet, error)
}

type walletRepository struct {
	db *sql.DB
}

func NewWalletRepository(db *sql.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) UpdateWalletBalance(tx *sql.Tx, walletID uuid.UUID, operationType domain.OperationType, amount float64) error {
	var query string
	switch operationType {
	case domain.Deposit:
		query = "UPDATE wallets SET balance = balance + $1 WHERE id = $2"
	case domain.Withdraw:
		query = "UPDATE wallets SET balance = balance - $1 WHERE id = $2"
	default:
		log.Printf("Invalid operation type: %s", operationType)
		return errors.New("invalid operation type")
	}

	if _, err := tx.Exec(query, amount, walletID); err != nil {
		log.Printf("Error executing query: %s, error: %v", query, err)
		return err
	}

	log.Printf("Successfully updated balance for walletID: %s, operationType: %s, amount: %f", walletID, operationType, amount)
	return nil
}

func (r *walletRepository) GetWalletBalance(walletID uuid.UUID) (float64, error) {
	var balance float64
	err := r.db.QueryRow("SELECT balance FROM wallets WHERE id = $1", walletID).Scan(&balance)
	if err != nil {
		log.Printf("Error getting balance for walletID: %s, error: %v", walletID, err)
	}
	return balance, err
}

func (r *walletRepository) GetWalletByID(walletID uuid.UUID) (*domain.Wallet, error) {
	var wallet domain.Wallet
	err := r.db.QueryRow("SELECT id, balance FROM wallets WHERE id = $1", walletID).Scan(&wallet.ID, &wallet.Balance)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Printf("Wallet not found for walletID: %s", walletID)
			return nil, errors.New("wallet not found")
		}
		log.Printf("Error getting wallet by ID: %s, error: %v", walletID, err)
		return nil, err
	}
	log.Printf("Successfully retrieved wallet for walletID: %s", walletID)
	return &wallet, nil
}
