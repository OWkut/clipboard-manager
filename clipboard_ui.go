//go:build ui
// +build ui

package main

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/atotto/clipboard"
)

const historyFile = "clipboard_history.json"

var internalCopy = false

func getActiveWindowID() string {
	cmd := exec.Command("xdotool", "getactivewindow")
	output, err := cmd.Output()
	if err != nil {
		log.Println("Erreur lors de l'obtention de la fenêtre active :", err)
		return ""
	}
	return string(output)
}

func restoreFocus(windowID string) {
	if windowID == "" {
		log.Println("ID de la fenêtre non valide.")
		return
	}
	cmd := exec.Command("xdotool", "windowactivate", windowID)
	err := cmd.Run()
	if err != nil {
		log.Println("Erreur lors du retour du focus :", err)
	}
}

func pasteClipboard() {
	cmd := exec.Command("xdotool", "key", "ctrl+v")
	err := cmd.Run()
	if err != nil {
		log.Println("Erreur lors de la simulation du collage avec xdotool :", err)
	}
}

func removeDuplicates(history []string) []string {
	unique := make(map[string]bool)
	result := []string{}

	for _, item := range history {
		if !unique[item] {
			unique[item] = true
			result = append(result, item)
		}
	}
	return result
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Gestionnaire de Presse-papiers")
	var selectedRow int = -1
	myWindow.SetMaster()

	history := loadHistory()

	// Créer un widget.Table pour afficher l'historique
	table := widget.NewTable(
		func() (int, int) {
			return len(history), 2 // 2 colonnes : index et texte
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.TableCellID, o fyne.CanvasObject) {
			if i.Col == 0 {
				o.(*widget.Label).SetText(strconv.Itoa(i.Row + 1)) // Afficher l'index
			} else {
				o.(*widget.Label).SetText(history[len(history)-1-i.Row]) // Afficher le texte
			}
		},
	)

	previousWindowID := getActiveWindowID()

	// Gérer la sélection dans le tableau
	table.OnSelected = func(id widget.TableCellID) {
		if id.Row >= 0 && id.Row < len(history) {
			text := history[len(history)-1-id.Row]
			err := clipboard.WriteAll(text)
			if err != nil {
				log.Println("Erreur lors de la copie dans le presse-papiers :", err)
			} else {
				internalCopy = true // Indiquer que c'est un copier interne
			}

			myWindow.Hide()
			restoreFocus(previousWindowID) // Force le retour du focus
			time.Sleep(100 * time.Millisecond)
			pasteClipboard()
		}
	}

	// Bouton pour effacer l'historique
	clearButton := widget.NewButton("Effacer l'historique", func() {
		history = []string{}          // Vider l'historique en mémoire
		err := os.Remove(historyFile) // Supprimer le fichier d'historique
		if err != nil {
			log.Println("Erreur lors de la suppression du fichier d'historique :", err)
		}
		table.Refresh()
	})

	// Bouton pour coller la sélection
	pasteButton := widget.NewButton("Coller la sélection", func() {
		if selectedRow >= 0 && selectedRow < len(history) {
			text := history[len(history)-1-selectedRow]
			err := clipboard.WriteAll(text)
			selectedRow = -1
			table.UnselectAll()
			if err != nil {
				log.Println("Erreur lors de la copie dans le presse-papiers :", err)
			}
		} else {
			log.Println("Aucune ligne sélectionnée.")
		}
	})

	// Mettre à jour l'historique en temps réel
	go func() {
		lastText := ""
		for {
			text, err := clipboard.ReadAll()
			if err == nil && text != lastText {
				lastText = text

				// Vérifie si le texte est déjà dans l'historique
				if len(history) == 0 || history[len(history)-1] != text {
					history = append(history, text)
					if len(history) > 20 { // Limiter l'historique à 20 éléments
						history = history[1:]
					}
					saveHistory(removeDuplicates(history))
					table.Refresh()
				}
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	// Organiser l'interface avec un conteneur Border
	myWindow.SetContent(container.NewBorder(
		widget.NewLabel("Historique du presse-papiers :"), // En haut
		container.NewHBox(clearButton, pasteButton),       // En bas
		nil,                        // À gauche (rien)
		nil,                        // À droite (rien)
		container.NewScroll(table), // Au centre (avec scrollbar)
	))

	myWindow.Resize(fyne.NewSize(600, 400)) // Ajustez la taille selon vos besoins
	myWindow.ShowAndRun()
}

func loadHistory() []string {
	if _, err := os.Stat(historyFile); os.IsNotExist(err) {
		return []string{}
	}

	data, err := os.ReadFile(historyFile)
	if err != nil {
		log.Println("Erreur lors de la lecture du fichier d'historique :", err)
		return []string{}
	}

	var history []string
	if err := json.Unmarshal(data, &history); err != nil {
		log.Println("Erreur lors de la désérialisation de l'historique :", err)
		return []string{}
	}

	return history
}

func saveHistory(history []string) {
	if len(history) == 0 {
		// Si l'historique est vide, supprimer le fichier
		err := os.Remove(historyFile)
		if err != nil && !os.IsNotExist(err) {
			log.Println("Erreur lors de la suppression du fichier d'historique :", err)
		}
		return
	}

	data, err := json.Marshal(history)
	if err != nil {
		log.Println("Erreur lors de la sérialisation de l'historique :", err)
		return
	}

	if err := os.WriteFile(historyFile, data, 0644); err != nil {
		log.Println("Erreur lors de l'écriture du fichier d'historique :", err)
	}
}
