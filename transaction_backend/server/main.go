package main

import (
	"log"
	"net/http"
	"os"

	"transaction-backend/application"
	"transaction-backend/infrastructure/db"
	"transaction-backend/infrastructure/ws"
	"transaction-backend/interfaces/consumer"
)

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://user:pass@localhost:5432/transactions?sslmode=disable"
	}

	repo, err := db.NewPostgresRepository(dsn)
	if err != nil {
		log.Fatal(err)
	}

	hub := ws.NewWebSocketHub(repo)
	service := application.NewTransactionService(repo, hub)

	err = consumer.StartRabbitConsumer(
		"amqp://guest:guest@localhost:5672/",
		"transactions-queue",
		"transactions-exchange",
		"transactions.new",
		service,
	)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws/transactions", hub.HandleConnection)

	log.Println("Backend running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
