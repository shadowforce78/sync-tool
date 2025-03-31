package main

import (
	"fmt"
	"net/http"
	"os"

	"SushiSync/config"
	"SushiSync/handlers"
)

func main() {
	// Crée le dossier data s'il n'existe pas
	if _, err := os.Stat(config.DataDir); os.IsNotExist(err) {
		os.Mkdir(config.DataDir, 0755)
	}

	// Enregistre les routes
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/download", handlers.DownloadHandler)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
