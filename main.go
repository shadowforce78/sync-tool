package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	maxUploadSize = 10 << 20 // 10MB
	dataDir       = "./data"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Limite la taille du fichier
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// Parse le fichier envoyé
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erreur en récupérant le fichier", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Nettoie le nom du fichier pour éviter les attaques path traversal
	filename := filepath.Base(header.Filename)
	filePath := filepath.Join(dataDir, filename)

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
		http.Error(w, "Erreur lors de l’écriture du fichier", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Fichier %s uploadé avec succès", filename)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère le nom du fichier depuis l'URL
	filename := r.URL.Query().Get("file")
	if filename == "" {
		http.Error(w, "Nom du fichier manquant", http.StatusBadRequest)
		return
	}

	// Nettoie le chemin pour éviter les attaques
	filename = filepath.Base(filename)
	filePath := filepath.Join(dataDir, filename)

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

func main() {
	// Crée le dossier data s'il n'existe pas
	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.Mkdir(dataDir, 0755)
	}

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Exemple upload : curl -X POST -F "file=@/path/to/file.txt" http://localhost:8080/upload
// Exemple download : curl -X GET "http://localhost:8080/download?file=file.txt" -o file.txt
