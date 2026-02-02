package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"kasir-api/internal/category/handler"
	categoryRepo "kasir-api/internal/category/repository"
	categorySvc "kasir-api/internal/category/service"
	"kasir-api/internal/config"
	"kasir-api/internal/db"
	productHandler "kasir-api/internal/product/handler"
	productRepo "kasir-api/internal/product/repository"
	productSvc "kasir-api/internal/product/service"
)

func main() {
	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("Failed to load config:", err)
		os.Exit(1)
	}

	// 2. Open database connection
	database, err := db.Open(cfg.Database.URL)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		os.Exit(1)
	}
	defer database.Close()

	// 3. Run database migrations
	if err := db.Migrate(database); err != nil {
		fmt.Println("Failed to migrate database:", err)
		os.Exit(1)
	}
	fmt.Println("Database migrated successfully")

	// 4. Initialize category layers with Dependency Injection
	catRepo := categoryRepo.NewPostgres(database)
	catSvc := categorySvc.New(catRepo)
	catHandler := handler.New(catSvc)

	// 5. Initialize product layers with Dependency Injection
	prodRepo := productRepo.NewPostgres(database)
	prodSvc := productSvc.New(prodRepo)
	prodHandler := productHandler.New(prodSvc)

	// 6. Setup HTTP routes
	http.HandleFunc("/health", healthHandler)

	// Category routes
	http.HandleFunc("/categories", catHandler.Handle)
	http.HandleFunc("/categories/", catHandler.HandleWithID)

	// Product routes
	http.HandleFunc("/api/produk", prodHandler.Handle)
	http.HandleFunc("/api/produk/", prodHandler.HandleWithID)

	// 7. Start server
	port := cfg.Server.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Println("Server error:", err)
		os.Exit(1)
	}
}

// healthHandler returns the health status of the API
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
