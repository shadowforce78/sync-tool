package handlers

import (
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"SushiSync/config"
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
	// Définir les chemins de sauvegarde pour le fichier compressé et le hash
	filePath := filepath.Join(config.DataDir, filename)
	compressedPath := filePath + ".gz"
	hashPath := filePath + ".sha256"

	// Lire le contenu complet du fichier pour le compresser et calculer son hash
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Erreur lors de la lecture du fichier", http.StatusInternalServerError)
		return
	}

	// Calculer le hash SHA256 du contenu original
	hash := sha256.Sum256(data)
	hashStr := hex.EncodeToString(hash[:])

	// Créer le fichier compressé
	out, err := os.Create(compressedPath)
	if err != nil {
		http.Error(w, "Erreur en sauvegardant le fichier compressé", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Utiliser le niveau de compression optimal
	gw, err := gzip.NewWriterLevel(out, gzip.BestCompression)
	if err != nil {
		http.Error(w, "Erreur lors de la création du writer gzip", http.StatusInternalServerError)
		return
	}
	_, err = gw.Write(data)
	if err != nil {
		http.Error(w, "Erreur lors de la compression", http.StatusInternalServerError)
		return
	}
	gw.Close()

	// Sauvegarder la valeur du hash dans un fichier
	err = os.WriteFile(hashPath, []byte(hashStr), 0644)
	if err != nil {
		http.Error(w, "Erreur en sauvegardant le hash", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Fichier %s uploadé et compressé avec succès", filename)
}
