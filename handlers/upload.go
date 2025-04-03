package handlers

import (
	"bytes"
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

	// Si le client a déjà compressé, on traite différemment
	if r.Header.Get("X-Client-Compressed") == "true" {
		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Erreur lors de la lecture du fichier compressé", http.StatusInternalServerError)
			return
		}
		// Décompresser pour obtenir le contenu original
		gr, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			http.Error(w, "Erreur lors de l'ouverture du gzip reader", http.StatusInternalServerError)
			return
		}
		originalData, err := io.ReadAll(gr)
		gr.Close()
		if err != nil {
			http.Error(w, "Erreur lors de la décompression", http.StatusInternalServerError)
			return
		}

		// Calculer le hash SHA256 sur le contenu original
		hash := sha256.Sum256(originalData)
		hashStr := hex.EncodeToString(hash[:])

		// Sauvegarder le fichier compressé (tel que reçu)
		err = os.WriteFile(compressedPath, data, 0644)
		if err != nil {
			http.Error(w, "Erreur en sauvegardant le fichier compressé", http.StatusInternalServerError)
			return
		}
		// Sauvegarder la valeur du hash
		err = os.WriteFile(hashPath, []byte(hashStr), 0644)
		if err != nil {
			http.Error(w, "Erreur en sauvegardant le hash", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Fichier %s uploadé et compressé avec succès", filename)
		return
	}

	// Créer le fichier compressé
	out, err := os.Create(compressedPath)
	if err != nil {
		http.Error(w, "Erreur en sauvegardant le fichier compressé", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// Créer le writer gzip avec compression optimale
	gw, err := gzip.NewWriterLevel(out, gzip.BestCompression)
	if err != nil {
		http.Error(w, "Erreur lors de la création du writer gzip", http.StatusInternalServerError)
		return
	}

	// Créer un hasher pour calculer le hash sur le contenu original
	hasher := sha256.New()
	// Copier depuis file vers le gzip writer et le hasher simultanément
	multiWriter := io.MultiWriter(gw, hasher)
	_, err = io.Copy(multiWriter, file)
	if err != nil {
		http.Error(w, "Erreur lors de l'écriture du fichier compressé", http.StatusInternalServerError)
		return
	}
	if err := gw.Close(); err != nil {
		http.Error(w, "Erreur lors de la clôture du writer gzip", http.StatusInternalServerError)
		return
	}

	// Convertir le hash en chaîne hexadécimale et sauvegarder dans un fichier
	hashStr := hex.EncodeToString(hasher.Sum(nil))
	err = os.WriteFile(hashPath, []byte(hashStr), 0644)
	if err != nil {
		http.Error(w, "Erreur en sauvegardant le hash", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Fichier %s uploadé et compressé avec succès", filename)
}
