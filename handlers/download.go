package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"

	"c:\Users\saumondeluxe\Desktop\DevProject\SushiSync\config"
)

// DownloadHandler gère les requêtes de téléchargement de fichiers
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère le nom du fichier depuis l'URL
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "Nom du fichier manquant", http.StatusBadRequest)
		return
	}

	// Nettoie le chemin pour éviter les attaques
	filename = filepath.Base(filename)
	filePath := filepath.Join(config.DataDir, filename)

	// Vérifie que le fichier existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		http.Error(w, "Fichier introuvable", http.StatusNotFound)
		return
	}

	// Ouvre le fichier pour le télécharger
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Erreur en ouvrant le fichier", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Définit les headers pour le téléchargement direct
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Stream directement le fichier sans le copier ailleurs
	io.Copy(w, file)
}
