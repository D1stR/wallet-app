package service

import (
	"WalletApp/internal/domain"
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"testing"
)

type MockWalletRepository struct {
	wallets map[uuid.UUID]*domain.Wallet
}

func (m *MockWalletRepository) GetWalletByID(walletID uuid.UUID) (*domain.Wallet, error) {
	wallet, exists := m.wallets[walletID]
	if !exists {
		return nil, ErrWalletNotFound
	}
	return wallet, nil
}

func (m *MockWalletRepository) UpdateWalletBalance(tx *sql.Tx, walletID uuid.UUID, operationType domain.OperationType, amount float64) error {
	wallet, exists := m.wallets[walletID]
	if !exists {
		return ErrWalletNotFound
	}
	if operationType == domain.Withdraw && wallet.Balance < amount {
		return ErrInsufficientFunds
	}
	if operationType == domain.Withdraw {
		wallet.Balance -= amount
	} else {
		wallet.Balance += amount
	}
	return nil
}

func (m *MockWalletRepository) GetWalletBalance(walletID uuid.UUID) (float64, error) {
	wallet, exists := m.wallets[walletID]
	if !exists {
		return 0, ErrWalletNotFound
	}
	return wallet.Balance, nil
}

func TestWalletService_UpdateWalletBalance(t *testing.T) {
	walletID := uuid.New()
	repo := &MockWalletRepository{
		wallets: map[uuid.UUID]*domain.Wallet{
			walletID: {ID: walletID, Balance: 100.0},
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	service := NewWalletService(repo, db)

	tests := []struct {
		name            string
		operationType   domain.OperationType
		amount          float64
		expectedError   error
		expectedBalance float64
	}{
		{
			name:            "Successful deposit",
			operationType:   domain.Deposit,
			amount:          50.0,
			expectedError:   nil,
			expectedBalance: 150.0,
		},
		{
			name:            "Successful withdrawal",
			operationType:   domain.Withdraw,
			amount:          50.0,
			expectedError:   nil,
			expectedBalance: 50.0,
		},
		{
			name:            "Insufficient funds",
			operationType:   domain.Withdraw,
			amount:          200.0,
			expectedError:   ErrInsufficientFunds,
			expectedBalance: 100.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var expectedQuery string
			var args []driver.Value

			switch tt.expectedError {
			case nil:
				expectedQuery = "UPDATE wallets SET balance = ? WHERE id = ?"
				args = []driver.Value{tt.expectedBalance, walletID}
			default:
				expectedQuery = ""
				args = nil
			}

			if tt.name == "Wallet not found" {
				repo.wallets = map[uuid.UUID]*domain.Wallet{} // Удаляем кошелёк из репозитория.
			} else if tt.expectedError == nil {
				mock.ExpectBegin()                                                                         // Начинаем транзакцию.
				mock.ExpectExec(expectedQuery).WithArgs(args...).WillReturnResult(sqlmock.NewResult(1, 1)) // Выполняем запрос обновления.
				mock.ExpectCommit()                                                                        // Завершаем транзакцию.
			}

			err := service.UpdateWalletBalance(walletID, tt.operationType, tt.amount)
			if err != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}

			if tt.name != "Wallet not found" {
				wallet, _ := repo.GetWalletByID(walletID)
				if wallet.Balance != tt.expectedBalance {
					t.Errorf("Expected balance %f, got %f", tt.expectedBalance, wallet.Balance)
				}
			}
		})
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestWalletService_GetWalletByID(t *testing.T) {
	walletID := uuid.New()
	repo := &MockWalletRepository{
		wallets: map[uuid.UUID]*domain.Wallet{
			walletID: {ID: walletID, Balance: 100.0},
		},
	}
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	service := NewWalletService(repo, db)

	tests := []struct {
		name          string
		walletID      uuid.UUID
		expectedError error
	}{
		{
			name:          "Wallet found",
			walletID:      walletID,
			expectedError: nil,
		},
		{
			name:          "Wallet not found",
			walletID:      uuid.New(),
			expectedError: ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.GetWalletByID(tt.walletID)
			if err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
		})
	}
}

func TestWalletService_GetWalletBalance(t *testing.T) {
	walletID := uuid.New()
	repo := &MockWalletRepository{
		wallets: map[uuid.UUID]*domain.Wallet{
			walletID: {ID: walletID, Balance: 100.0},
		},
	}
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	defer db.Close()

	service := NewWalletService(repo, db)

	tests := []struct {
		name            string
		walletID        uuid.UUID
		expectedBalance float64
		expectedError   error
	}{
		{
			name:            "Wallet found",
			walletID:        walletID,
			expectedBalance: 100.0,
			expectedError:   nil,
		},
		{
			name:            "Wallet not found",
			walletID:        uuid.New(),
			expectedBalance: 0.0,
			expectedError:   ErrWalletNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			balance, err := service.GetWalletBalance(tt.walletID)
			if err != tt.expectedError {
				t.Errorf("Expected error %v, got %v", tt.expectedError, err)
			}
			if balance != tt.expectedBalance {
				t.Errorf("Expected balance %f, got %f", tt.expectedBalance, balance)
			}
		})
	}
}
