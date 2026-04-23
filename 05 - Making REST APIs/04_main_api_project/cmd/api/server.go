package main

import (
	"fmt"
	"log"
	"net/http"
	mw "restapi/internal/api/middlewares"
	"restapi/internal/api/router"
	"restapi/pkg/utils"
	"time"
)

func main() {

	port := "3000"
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

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}
