package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Parse le fichier envoyé
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erreur en récupérant le fichier", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Crée un fichier sur le serveur
	out, err := os.Create("./data/" + header.Filename)
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

	fmt.Fprintf(w, "Fichier %s uploadé avec succès", header.Filename)
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	// Récupère le nom du fichier depuis l'URL
	filename := r.URL.Query().Get("file")

	// Vérifie que le fichier source existe
	sourcePath := filepath.Join("./data/", filename)
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		http.Error(w, "Erreur en ouvrant le fichier source", http.StatusNotFound)
		return
	}
	defer sourceFile.Close()

	// Crée le fichier dans le dossier output
	outputPath := filepath.Join("./output/", filename)
	outputFile, err := os.Create(outputPath)
	if err != nil {
		http.Error(w, "Erreur en créant le fichier dans le dossier output", http.StatusInternalServerError)
		return
	}
	defer outputFile.Close()

	// Copie le contenu du fichier source vers le fichier output
	_, err = io.Copy(outputFile, sourceFile)
	if err != nil {
		http.Error(w, "Erreur lors de la copie du fichier", http.StatusInternalServerError)
		return
	}

	// Ferme le fichier output pour s'assurer que l'écriture est terminée
	outputFile.Close()

	// Ouvre le fichier output pour le servir
	fileToServe, err := os.Open(outputPath)
	if err != nil {
		http.Error(w, "Erreur en ouvrant le fichier pour le téléchargement", http.StatusInternalServerError)
		return
	}
	defer fileToServe.Close()

	// Définit les headers pour le téléchargement
	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	// Copie le contenu du fichier dans la réponse
	io.Copy(w, fileToServe)
}

func main() {
	// Crée le dossier data s'il n'existe pas
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		os.Mkdir("./data", 0755)
	}

	// Crée le dossier output s'il n'existe pas
	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		os.Mkdir("./output", 0755)
	}

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/download", downloadHandler)

	fmt.Println("Serveur démarré sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Exemple upload : curl -X POST -F "file=@/path/to/file.txt" http://localhost:8080/upload
// Exemple download : curl -X GET http://localhost:8080/download?file=file.txt -o file.txt
