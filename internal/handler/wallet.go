package handlers

import (
	"WalletApp/internal/domain"
	"WalletApp/internal/models"
	"WalletApp/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type WalletHandler struct {
	service service.WalletServiceInterface
}

func NewWalletHandler(service service.WalletServiceInterface) *WalletHandler {
	return &WalletHandler{service: service}
}

func (h *WalletHandler) HandleWalletOperation(c *fiber.Ctx) error {
	var request models.WalletRequest
	if err := c.BodyParser(&request); err != nil {
		log.Printf("Error parsing request body: %v", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	walletID, err := uuid.Parse(request.WalletID)
	if err != nil {
		log.Printf("Invalid wallet ID: %s, error: %v", request.WalletID, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid wallet ID"})
	}

	opType, err := domain.OperationTypeFromString(request.OperationType)
	if err != nil {
		log.Printf("Error parsing operation type: %s, error: %v", request.OperationType, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err = h.service.UpdateWalletBalance(walletID, opType, request.Amount)
	if err != nil {
		log.Printf("Error updating wallet balance for walletID: %s, error: %v", walletID, err)
		if err == service.ErrWalletNotFound {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Wallet not found"})
		}
		if err == service.ErrInsufficientFunds {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Insufficient funds"})
		}
		if err == service.ErrInvalidOperationType {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid operation type"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update wallet"})
	}

	log.Printf("Successfully handled wallet operation for walletID: %s", walletID)
	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Operation successful"})
}

func (h *WalletHandler) GetWalletBalance(c *fiber.Ctx) error {
	walletIDParam := c.Params("uuid")

	walletID, err := uuid.Parse(walletIDParam)
	if err != nil {
		log.Printf("Invalid wallet ID from params: %s, error: %v", walletIDParam, err)
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Invalid wallet ID"})
	}

	balance, err := h.service.GetWalletBalance(walletID)
	if err != nil {
		log.Printf("Error retrieving balance for walletID: %s, error: %v", walletID, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Could not retrieve balance"})
	}

	log.Printf("Successfully retrieved balance for walletID: %s, balance: %f", walletID, balance)
	return c.Status(http.StatusOK).JSON(fiber.Map{"balance": balance})
}
