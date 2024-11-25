package main

import (
	"WalletApp/internal/config"
	"WalletApp/internal/handler"
	"WalletApp/internal/repository"
	"WalletApp/internal/routes"
	"WalletApp/internal/service"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
)

func main() {
	if err := SetupViper(); err != nil {
		log.Fatal(err.Error())
	}

	cfg := config.LoadConfig(viper.GetViper())
	app := fiber.New()

	db := initDatabase(&cfg)
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	repo := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(repo, db)
	walletHandler := handlers.NewWalletHandler(walletService)

	routes.SetupRoutes(app, walletHandler)
	if err := app.Listen(":8080"); err != nil {
		log.Fatal("can`t run app", err)
	}
}

func SetupViper() error {
	viper.SetConfigName("config")
	viper.SetConfigType("env")
	viper.AddConfigPath("./configs")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}
	return nil
}

func initDatabase(cfg *config.Config) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}

	return db
}
