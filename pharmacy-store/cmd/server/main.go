package main

import (
	"log"
	"pharmacy-store/configs"
	"pharmacy-store/internal/infrastructure/http"
	natsClient "pharmacy-store/internal/infrastructure/messaging/nats"
	"pharmacy-store/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/nats-io/nats.go"
)

func main() {
	config, err := configs.LoadConfig(".")
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	db, err := configs.InitDB(config)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	natsClient, err := natsClient.NewNatsClient(config.NATSUrl)
	if err != nil {
		log.Fatalf("Could not connect to NATS: %v", err)
	}
	defer natsClient.Conn.Close()

	_, err = natsClient.Conn.Subscribe("test.subject", func(m *nats.Msg) {
		log.Printf("Received a message: %s", string(m.Data))
	})
	if err != nil {
		log.Fatalf("Error subscribing to test.subject: %v", err)
	}

	router := http.NewRouter(db, natsClient)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Replace with your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(middleware.Logger())

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Could not run server: %v", err)
	}
}
