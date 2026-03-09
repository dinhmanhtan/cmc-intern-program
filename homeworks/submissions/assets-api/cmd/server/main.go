package main

import (
	"log"
	"net/http"

	"assets-api/internal/config"
	"assets-api/internal/database"
	"assets-api/internal/handler"
	"assets-api/internal/repository"
	"assets-api/internal/service"
)

func main() {

	cfg := config.Load()

	db, err := database.ConnectWithRetry(cfg.DB_DSN, 5)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewAssetRepository(db)
	service := service.NewAssetService(repo)

	assetHandler := handler.NewAssetHandler(service)

	http.HandleFunc("/assets/stats", assetHandler.Stats)
	http.HandleFunc("/assets/count", assetHandler.Count)

	http.HandleFunc("/assets/batch", func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodPost {
			assetHandler.BatchCreate(w, r)
			return
		}

		if r.Method == http.MethodDelete {
			assetHandler.BatchDelete(w, r)
			return
		}
	})

	http.HandleFunc("/assets", assetHandler.List)
	http.HandleFunc("/assets/search", assetHandler.Search)

	// gọi đúng HealthHandler
	http.HandleFunc("/health", handler.HealthHandler(db))

	log.Println("Server running on :" + cfg.PORT)

	http.ListenAndServe(":"+cfg.PORT, nil)
}
