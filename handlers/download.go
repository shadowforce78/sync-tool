package handlers

import (
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"SushiSync/config"
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
	compressedPath := filePath + ".gz"
	hashPath := filePath + ".sha256"

	// Vérifier que le fichier compressé et le hash existent
	if _, err := os.Stat(compressedPath); os.IsNotExist(err) {
		http.Error(w, "Fichier introuvable", http.StatusNotFound)
		return
	}
	if _, err := os.Stat(hashPath); os.IsNotExist(err) {
		http.Error(w, "Fichier de vérification d'intégrité introuvable", http.StatusNotFound)
		return
	}

	// Si le paramètre "decompress" est true, décompresse et vérifie l'intégrité
	if r.URL.Query().Get("decompress") == "true" {
		f, err := os.Open(compressedPath)
		if err != nil {
			http.Error(w, "Erreur en ouvrant le fichier compressé", http.StatusInternalServerError)
			return
		}
		defer f.Close()

		gz, err := gzip.NewReader(f)
		if err != nil {
			http.Error(w, "Erreur lors de la décompression", http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		var buf bytes.Buffer
		_, err = io.Copy(&buf, gz)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture décompressée", http.StatusInternalServerError)
			return
		}
		decompressedData := buf.Bytes()

		// Calculer le hash du contenu décompressé
		calculatedHash := sha256.Sum256(decompressedData)
		calculatedHashStr := hex.EncodeToString(calculatedHash[:])
		savedHashData, err := os.ReadFile(hashPath)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture du hash sauvegardé", http.StatusInternalServerError)
			return
		}
		savedHashStr := string(savedHashData)
		if calculatedHashStr != savedHashStr {
			http.Error(w, "Erreur d'intégrité du fichier", http.StatusInternalServerError)
			return
		}

		// Envoyer le contenu décompressé
		w.Header().Set("Content-Disposition", "attachment; filename="+filename)
		w.Header().Set("Content-Type", "application/octet-stream")
		_, err = io.Copy(w, bytes.NewReader(decompressedData))
		if err != nil {
			http.Error(w, "Erreur lors de l'envoi du fichier", http.StatusInternalServerError)
			return
		}
		return
	}

	// Sinon, envoyer le fichier compressé tel quel
	f, err := os.Open(compressedPath)
	if err != nil {
		http.Error(w, "Erreur en ouvrant le fichier compressé", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// On peut indiquer que le fichier est compressé pour que le client sache quoi faire
	w.Header().Set("Content-Disposition", "attachment; filename="+filename+".gz")
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, f)
	if err != nil {
		http.Error(w, "Erreur lors de l'envoi du fichier", http.StatusInternalServerError)
		return
	}
}
