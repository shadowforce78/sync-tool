package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/saumondeluxe/SushiSync/config"
)

// UploadHandler gère les requêtes d'upload de fichiers
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	// Limite la taille du fichier
	r.Body = http.MaxBytesReader(w, r.Body, config.MaxUploadSize)

	// Parse le fichier envoyé
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erreur en récupérant le fichier", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Nettoie le nom du fichier pour éviter les attaques path traversal
	filename := filepath.Base(header.Filename)
	filePath := filepath.Join(config.DataDir, filename)

	// Crée un fichier sur le serveur
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Erreur en sauvegardant le fichier", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Copie le contenu
	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Erreur lors de l'écriture du fichier", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Fichier %s uploadé avec succès", filename)
}
