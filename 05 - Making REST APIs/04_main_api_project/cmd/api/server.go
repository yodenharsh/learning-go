package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	mw "restapi/internal/api/middlewares"
	"restapi/internal/api/router"
	"restapi/internal/repository/sqlconnect"
	"restapi/pkg/utils"
	"time"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	sqlconnect.ConnectDb()

	port := os.Getenv("API_PORT")
	fmt.Println("Server is running on port ", port)
	hppOptions := mw.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		Whitelist:                   []string{"sortBy", "sortOrder", "name", "age", "class", "firstName", "lastName"},
	}

	rl := mw.NewRateLimiter(5, time.Minute)

	// Middlewares order is first-in first-applied
	secureMux := utils.ApplyMiddlewares(router.Router(),
		mw.Hpp(hppOptions),
		mw.Compression,
		mw.SecurityHeaders,
		mw.ResponseTimeMiddleware,
		rl.RateLimitingMiddleware,
		mw.Cors)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: secureMux,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
