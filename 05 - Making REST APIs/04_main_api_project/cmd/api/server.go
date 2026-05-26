package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	mw "restapi/internal/api/middlewares"
	"restapi/internal/api/router"
	"restapi/pkg/utils"
	"time"

	"github.com/joho/godotenv"
)

//go:embed .env
var envFile embed.FS

func loadEnvFromEmbeddedFile() {
	content, err := envFile.ReadFile(".env")
	if err != nil {
		log.Fatalf("Error reading the .env file")
		return
	}

	// Create a temporary file to store the .env content
	tempFile, err := os.CreateTemp("", ".env")
	if err != nil {
		log.Fatalf("Error creating temporary file for .env content")
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(content)
	if err != nil {
		log.Fatalf("Error writing to temporary .env file")
		return
	}

	err = tempFile.Close()
	if err != nil {
		log.Fatalf("Error closing temporary .env file")
		return
	}

	err = godotenv.Load(tempFile.Name())
	if err != nil {
		log.Fatalf("Error loading .env file from embedded content: %v", err)
	}
}

func main() {

	// Only for development
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading .env file:", err)
	// }

	loadEnvFromEmbeddedFile()

	port := os.Getenv("API_PORT")
	fmt.Println("Server is running on port ", port)
	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class", "firstName", "lastName", "page", "pageSize"},
	}

	rl := mw.NewRateLimiter(5, time.Minute)

	// Middlewares order is first-in first-applied
	secureMux := utils.ApplyMiddlewares(router.Router(),
		mw.SecurityHeaders,
		mw.Compression,
		mw.Hpp(hppOptions),
		mw.ResponseTimeMiddleware,
		rl.RateLimitingMiddleware,
		mw.Cors)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: secureMux,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
