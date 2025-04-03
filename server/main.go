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

	// Variable port pour le serveur
	PORT := "8080"

	// Enregistre les routes
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/download", handlers.DownloadHandler)
	http.HandleFunc("/list", handlers.ListHandler)

	fmt.Println("Serveur démarré sur http://localhost:" + PORT)
	http.ListenAndServe(":"+PORT, nil)
}
