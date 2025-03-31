package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ServerURL est l'adresse du serveur SushiSync
var ServerURL = "http://localhost:8080"

// FileInfo représente les informations sur un fichier
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func main() {
	fmt.Println("=== SushiSync Client CLI ===")
	fmt.Printf("Serveur: %s\n\n", ServerURL)

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("\nCommandes disponibles:")
		fmt.Println("1. Lister les fichiers")
		fmt.Println("2. Télécharger un fichier")
		fmt.Println("3. Uploader un fichier")
		fmt.Println("4. Changer l'URL du serveur")
		fmt.Println("5. Quitter")
		fmt.Print("\nChoix: ")

		scanner.Scan()
		choice := scanner.Text()

		switch choice {
		case "1":
			listFiles()
		case "2":
			downloadFile(scanner)
		case "3":
			uploadFile(scanner)
		case "4":
			changeServerURL(scanner)
		case "5":
			fmt.Println("Au revoir!")
			return
		default:
			fmt.Println("Choix invalide. Veuillez réessayer.")
		}
	}
}

// Affiche la liste des fichiers disponibles sur le serveur
func listFiles() {
	fmt.Println("\n=== Liste des fichiers ===")

	// Appel à l'API pour obtenir la liste des fichiers
	resp, err := http.Get(ServerURL + "/list")
	if err != nil {
		fmt.Printf("Erreur: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Erreur serveur: %d\n", resp.StatusCode)
		return
	}

	// Décode la réponse JSON
	var files []FileInfo
	err = json.NewDecoder(resp.Body).Decode(&files)
	if err != nil {
		fmt.Printf("Erreur de décodage: %v\n", err)
		return
	}

	// Affiche les fichiers
	if len(files) == 0 {
		fmt.Println("Aucun fichier trouvé.")
		return
	}

	fmt.Println("ID | Nom | Taille")
	fmt.Println("---------------")
	for i, file := range files {
		fmt.Printf("%d | %s | %d octets\n", i+1, file.Name, file.Size)
	}
}

// Télécharge un fichier depuis le serveur
func downloadFile(scanner *bufio.Scanner) {
	fmt.Println("\n=== Téléchargement de fichier ===")

	// Récupère la liste des fichiers
	resp, err := http.Get(ServerURL + "/list")
	if err != nil {
		fmt.Printf("Erreur: %v\n", err)
		return
	}

	var files []FileInfo
	err = json.NewDecoder(resp.Body).Decode(&files)
	resp.Body.Close()

	if err != nil {
		fmt.Printf("Erreur de décodage: %v\n", err)
		return
	}

	if len(files) == 0 {
		fmt.Println("Aucun fichier disponible.")
		return
	}

	// Affiche les fichiers
	fmt.Println("Fichiers disponibles:")
	for i, file := range files {
		fmt.Printf("%d. %s (%d octets)\n", i+1, file.Name, file.Size)
	}

	// Demande à l'utilisateur de choisir un fichier
	fmt.Print("\nEntrez l'ID du fichier à télécharger (ou 0 pour annuler): ")
	scanner.Scan()
	idStr := scanner.Text()

	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 || id > len(files) {
		fmt.Println("ID invalide.")
		return
	}

	fileName := files[id-1].Name

	// Demande où enregistrer le fichier
	fmt.Print("Entrez le chemin où enregistrer le fichier (laisser vide pour le dossier courant): ")
	scanner.Scan()
	savePath := scanner.Text()

	if savePath == "" {
		savePath = "."
	}

	filePath := filepath.Join(savePath, fileName)

	// Télécharge le fichier
	fileResp, err := http.Get(fmt.Sprintf("%s/download?file=%s", ServerURL, fileName))
	if err != nil {
		fmt.Printf("Erreur de téléchargement: %v\n", err)
		return
	}
	defer fileResp.Body.Close()

	if fileResp.StatusCode != http.StatusOK {
		fmt.Printf("Erreur serveur: %d\n", fileResp.StatusCode)
		return
	}

	// Crée le fichier local
	out, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("Erreur de création du fichier: %v\n", err)
		return
	}
	defer out.Close()

	// Copie le contenu
	_, err = io.Copy(out, fileResp.Body)
	if err != nil {
		fmt.Printf("Erreur d'écriture: %v\n", err)
		return
	}

	fmt.Printf("Fichier téléchargé avec succès: %s\n", filePath)
}

// Upload un fichier vers le serveur
func uploadFile(scanner *bufio.Scanner) {
	fmt.Println("\n=== Upload de fichier ===")

	// Demande le chemin du fichier à uploader
	fmt.Print("Entrez le chemin du fichier à uploader: ")
	scanner.Scan()
	filePath := scanner.Text()

	if filePath == "" {
		fmt.Println("Chemin invalide.")
		return
	}

	// Vérifie que le fichier existe
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Println("Le fichier n'existe pas.")
		return
	}

	// Ouvre le fichier
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Erreur d'ouverture du fichier: %v\n", err)
		return
	}
	defer file.Close()

	// Prépare la requête multipart
	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	fileWriter, err := multipartWriter.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Printf("Erreur de création du formulaire: %v\n", err)
		return
	}

	// Copie le contenu du fichier dans la requête
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		fmt.Printf("Erreur de copie du fichier: %v\n", err)
		return
	}

	multipartWriter.Close()

	// Envoie la requête au serveur
	req, err := http.NewRequest("POST", ServerURL+"/upload", &requestBody)
	if err != nil {
		fmt.Printf("Erreur de création de la requête: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erreur d'envoi: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Printf("Erreur serveur: %d - %s\n", resp.StatusCode, string(bodyBytes))
		return
	}

	fmt.Println("Fichier uploadé avec succès!")
}

// Change l'URL du serveur
func changeServerURL(scanner *bufio.Scanner) {
	fmt.Printf("\nURL actuelle: %s\n", ServerURL)
	fmt.Print("Nouvelle URL (laisser vide pour annuler): ")
	scanner.Scan()
	newURL := scanner.Text()

	if newURL != "" {
		// Vérifie le format de l'URL
		if !strings.HasPrefix(newURL, "http://") && !strings.HasPrefix(newURL, "https://") {
			newURL = "http://" + newURL
		}
		ServerURL = newURL
		fmt.Printf("URL changée pour: %s\n", ServerURL)
	}
}
