//go:build daemon
// +build daemon

package main

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/atotto/clipboard"
)

const historyFile = "clipboard_history.json"

func main() {
	history := loadHistory()

	lastText := ""
	for {
		text, err := clipboard.ReadAll()
		if err == nil && text != lastText {
			lastText = text
			history = append(history, text)
			if len(history) > 20 { // Limite l'historique à 20 éléments
				history = history[1:]
			}
			saveHistory(history)
		}
		time.Sleep(500 * time.Millisecond)
	}
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
