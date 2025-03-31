package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"SushiSync/config"
)

// FileInfo représente les informations sur un fichier
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// ListHandler gère les requêtes de listage des fichiers disponibles
func ListHandler(w http.ResponseWriter, r *http.Request) {
	// Liste tous les fichiers dans le répertoire de données
	files, err := os.ReadDir(config.DataDir)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du répertoire", http.StatusInternalServerError)
		return
	}

	// Prépare la liste des informations de fichiers
	var fileInfos []FileInfo
	for _, file := range files {
		if file.IsDir() {
			continue // Ignore les sous-répertoires
		}

		info, err := file.Info()
		if err != nil {
			continue // Ignore les fichiers problématiques
		}

		fileInfos = append(fileInfos, FileInfo{
			Name: file.Name(),
			Size: info.Size(),
		})
	}

	// Configure l'en-tête pour JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode et renvoie la liste au format JSON
	json.NewEncoder(w).Encode(fileInfos)
}
