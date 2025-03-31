package handlers

import (
	"net/http"
	"os"

	"SushiSync/config"
)

func ListHandler(w http.ResponseWriter, r *http.Request) {
	// Récupérer le répertoire de destination
	destDir := config.DataDir

	// Ouvrir le répertoire
	dir, err := os.Open(destDir)
	if err != nil {
		http.Error(w, "Erreur en ouvrant le répertoire", http.StatusInternalServerError)
		return
	}
	defer dir.Close()

	// Lire les fichiers dans le répertoire
	files, err := dir.Readdir(-1)
	if err != nil {
		http.Error(w, "Erreur en lisant le répertoire", http.StatusInternalServerError)
		return
	}

	// Créer une liste HTML des fichiers
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("<h1>Liste des fichiers</h1><ul>"))
	for _, file := range files {
		fileName := file.Name()
		w.Write([]byte("<li><a href=\"/download?file=" + fileName + "\">" + fileName + "</a></li>"))
	}
	w.Write([]byte("</ul>"))
}