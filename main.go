package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

var tasks []Task
var nextID int = 1

const dataFile = "tasks.json"

func add(title string) {
	t := Task{
		ID:        nextID,
		Title:     title,
		Done:      false,
		CreatedAt: time.Now(),
	}
	tasks = append(tasks, t)
	nextID++
	fmt.Println("Tâche ajoutée :", t.Title)
}
func list() {
	if len(tasks) == 0 {
		fmt.Println("Aucune tâche trouvée.")
		return
	}
	for _, t := range tasks {
		status := " "
		if t.Done {
			status = "x"
		}
		fmt.Printf("[%s] %d: %s (créée: %s)\n", status, t.ID, t.Title, t.CreatedAt.Format("01-02-2006 15:04:05"))
	}
}
func save() error {
	b, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dataFile, b, 0644)
}

func load() error {
	b, err := ioutil.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Pas de fichier, pas de tâches à charger
		}
		return err
	}
	if err := json.Unmarshal(b, &tasks); err != nil {
		return err
	}
	// Mettre à jour nextID
	max := 0
	for _, t := range tasks {
		if t.ID > max {
			max = t.ID
		}
	}
	nextID = max + 1
	return nil
}

func main() {
	// Flags simples
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addTitle := addCmd.String("title", "", "Titre de la tâche")

	listCmd := flag.NewFlagSet("list", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Usage : todo-go <command> [--flags]")
		fmt.Println("Commands : add, list")
		os.Exit(1)
	}
	if err := load(); err != nil {
		fmt.Println("Erreur de chargement des tâches :", err)
		os.Exit(1)
	}
	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		if *addTitle == "" {
			fmt.Println("Usage: todo add --title \"Ma tâche\"")
			os.Exit(1)
		}
		add(*addTitle)
		if err := save(); err != nil {
			fmt.Println("Erreur de sauvegarde des tâches :", err)
			os.Exit(1)
		}
	case "list":
		listCmd.Parse(os.Args[2:])
		list()
	default:
		fmt.Println("Commande inconnue:", os.Args[1])
		os.Exit(1)
	}
}
