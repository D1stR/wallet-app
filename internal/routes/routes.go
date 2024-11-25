package routes

import (
	"WalletApp/internal/handler"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, walletHandler *handlers.WalletHandler) {
	v1 := app.Group("/api/v1")

	v1.Post("/wallet", walletHandler.HandleWalletOperation)
	v1.Get("/wallets/:uuid", walletHandler.GetWalletBalance)
}
