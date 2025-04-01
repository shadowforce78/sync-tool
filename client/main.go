package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// ServerURL est l'adresse du serveur SushiSync
var ServerURL = "http://localhost:8080"

// FileInfo représente les informations sur un fichier
type FileInfo struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

func main() {
	// Crée une nouvelle application Fyne
	a := app.New()
	w := a.NewWindow("SushiSync Client")
	w.Resize(fyne.NewSize(600, 400))

	// URL du serveur
	serverURLEntry := widget.NewEntry()
	serverURLEntry.SetText(ServerURL)
	serverURLLabel := widget.NewLabel("URL du serveur:")
	serverURLContainer := container.New(layout.NewHBoxLayout(), serverURLLabel, serverURLEntry)

	// Déclare la variable files avant son utilisation
	var files []FileInfo
	var selectedID widget.ListItemID = -1
	var fileList *widget.List

	// Fonction pour rafraîchir la liste des fichiers
	refreshFileList := func() {
		resp, err := http.Get(serverURLEntry.Text + "/list")
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			dialog.ShowError(fmt.Errorf("erreur serveur: %d", resp.StatusCode), w)
			return
		}

		// Réinitialiser la liste des fichiers avant de la remplir à nouveau
		files = []FileInfo{}

		err = json.NewDecoder(resp.Body).Decode(&files)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}

		// Afficher le nombre de fichiers pour le débogage
		fmt.Printf("Nombre de fichiers récupérés: %d\n", len(files))

		// Reset la sélection quand la liste est mise à jour
		selectedID = -1

		// Force la reconstruction complète de la liste
		fileList.Refresh()
	}
	// Liste des fichiers
	fileList = widget.NewList(
		func() int {
			return len(files)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("Fichier très long pour garantir l'espace")
		},
		func(id widget.ListItemID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			if id < len(files) {
				file := files[id]
				label.SetText(fmt.Sprintf("%s (%d bytes)", file.Name, file.Size))
			} else {
				label.SetText("") // Efface le texte si l'ID est hors limites
			}
		},
	)
	fileList.OnSelected = func(id widget.ListItemID) {
		selectedID = id
	}

	// Bouton pour rafraîchir la liste
	refreshBtn := widget.NewButton("Rafraîchir", refreshFileList)

	// Bouton pour télécharger un fichier
	downloadBtn := widget.NewButton("Télécharger", func() {
		if len(files) == 0 || selectedID < 0 {
			dialog.ShowInformation("Information", "Veuillez sélectionner un fichier à télécharger", w)
			return
		}

		fileName := files[selectedID].Name
		if fileName == "" {
			dialog.ShowInformation("Information", "Veuillez sélectionner un fichier à télécharger", w)
			return
		}

		saveDialog := dialog.NewFileSave(
			func(writer fyne.URIWriteCloser, err error) {
				if err != nil || writer == nil {
					return
				}
				defer writer.Close()

				resp, err := http.Get(fmt.Sprintf("%s/download?file=%s", serverURLEntry.Text, fileName))
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					dialog.ShowError(fmt.Errorf("erreur serveur: %d", resp.StatusCode), w)
					return
				}

				_, err = io.Copy(writer, resp.Body)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				dialog.ShowInformation("Succès", "Fichier téléchargé avec succès", w)
			},
			w,
		)
		saveDialog.SetFileName(fileName)
		saveDialog.Show()
	})

	// Bouton pour uploader un fichier
	uploadBtn := widget.NewButton("Uploader", func() {
		openDialog := dialog.NewFileOpen(
			func(reader fyne.URIReadCloser, err error) {
				if err != nil || reader == nil {
					return
				}
				defer reader.Close()

				var requestBody bytes.Buffer
				multipartWriter := multipart.NewWriter(&requestBody)

				fileWriter, err := multipartWriter.CreateFormFile("file", filepath.Base(reader.URI().Path()))
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				_, err = io.Copy(fileWriter, reader)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}

				multipartWriter.Close()

				req, err := http.NewRequest("POST", serverURLEntry.Text+"/upload", &requestBody)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					dialog.ShowError(err, w)
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					bodyBytes, _ := io.ReadAll(resp.Body)
					dialog.ShowError(fmt.Errorf("erreur serveur: %d - %s", resp.StatusCode, string(bodyBytes)), w)
					return
				}

				dialog.ShowInformation("Succès", "Fichier uploadé avec succès", w)
				refreshFileList()
			},
			w,
		)
		openDialog.Show()
	})

	// Organisation des widgets dans la fenêtre
	buttonsContainer := container.New(layout.NewHBoxLayout(), refreshBtn, downloadBtn, uploadBtn)

	// Augmente la taille du conteneur de la liste pour tous les fichiers
	fileListContainer := container.New(layout.NewVBoxLayout(), widget.NewLabel("Fichiers disponibles:"), fileList)
	fileListContainer.Resize(fyne.NewSize(500, 300))

	content := container.NewBorder(
		container.NewVBox(serverURLContainer, buttonsContainer),
		nil, nil, nil,
		fileListContainer,
	)

	w.SetContent(content)

	// Rafraîchit la liste au démarrage
	refreshFileList()

	w.ShowAndRun()
}