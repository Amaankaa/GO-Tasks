package main

import (
	"log"
	"task_manager/data"
	"task_manager/router"
)

func main() {
	if err := data.InitMongo(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	r := router.SetupRouter()
	r.Run(":8080")
}
