package service

import (
	"WalletApp/internal/domain"
	"WalletApp/internal/repository"
	"database/sql"
	"github.com/google/uuid"
	"log"
)

type WalletServiceInterface interface {
	UpdateWalletBalance(walletID uuid.UUID, operationType domain.OperationType, amount float64) error
	GetWalletByID(walletID uuid.UUID) (*domain.Wallet, error)
	GetWalletBalance(walletID uuid.UUID) (float64, error)
}

type WalletService struct {
	repo repository.WalletRepository
	db   *sql.DB
}

func NewWalletService(repo repository.WalletRepository, db *sql.DB) WalletServiceInterface {
	return &WalletService{repo: repo, db: db}
}

func (s *WalletService) UpdateWalletBalance(walletID uuid.UUID, operationType domain.OperationType, amount float64) error {
	log.Printf("Starting UpdateWalletBalance for walletID: %s, operationType: %s, amount: %f", walletID, operationType, amount)

	tx, err := s.db.Begin()
	if err != nil {
		log.Printf("Error beginning transaction: %v", err)
		return err
	}

	defer func() {
		if err != nil {
			log.Printf("Transaction rollback due to error: %v", err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("Error rolling back transaction: %v", rollbackErr)
			}
		} else {
			if commitErr := tx.Commit(); commitErr != nil {
				log.Printf("Error committing transaction: %v", commitErr)
				err = commitErr
			}
		}
	}()

	wallet, err := s.repo.GetWalletByID(walletID)
	if err != nil {
		log.Printf("Error getting wallet by ID: %s, error: %v", walletID, err)
		return ErrWalletNotFound
	}

	if operationType == domain.Withdraw {
		if wallet.Balance < amount {
			log.Printf("Insufficient funds for walletID: %s, balance: %f, requested: %f", walletID, wallet.Balance, amount)
			return ErrInsufficientFunds
		}
	}

	if err := s.repo.UpdateWalletBalance(tx, walletID, operationType, amount); err != nil {
		log.Printf("Error updating wallet balance for walletID: %s, error: %v", walletID, err)
		return err
	}

	log.Printf("Successfully updated wallet balance for walletID: %s", walletID)
	return nil
}

func (s *WalletService) GetWalletByID(walletID uuid.UUID) (*domain.Wallet, error) {
	log.Printf("Getting wallet by ID: %s", walletID)
	return s.repo.GetWalletByID(walletID)
}

func (s *WalletService) GetWalletBalance(walletID uuid.UUID) (float64, error) {
	log.Printf("Getting wallet balance for walletID: %s", walletID)
	return s.repo.GetWalletBalance(walletID)
}
