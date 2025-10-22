package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
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
func suppr(id int) error {
	for i, t := range tasks {
		if t.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)
			return nil
		}
	}
	return errors.New("Tâche non trouvée")
}

func markDone(id int) error {
	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Done = true
			return nil
		}
	}
	return errors.New("Tâche non trouvée")
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
func find(title string) {
	var results []Task
	for _, t := range tasks {
		if strings.Contains(strings.ToLower(t.Title), strings.ToLower(title)) {
			results = append(results, t)
			status := " "
			if t.Done {
				status = "x"
			}
			fmt.Printf("[%s] %d: %s (créée: %s)\n", status, t.ID, t.Title, t.CreatedAt.Format("01-02-2006 15:04:05"))
		}
	}
	if len(results) == 0 {
		fmt.Println("Aucune tâche trouvée avec le titre :", title)
	}
}
func save() error {
	b, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(dataFile, b, 0644)
}

func load() error {
	b, err := os.ReadFile(dataFile)
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
	doneCmd := flag.NewFlagSet("done", flag.ExitOnError)
	doneId := doneCmd.Int("id", 0, "ID de la tâche")
	supprCmd := flag.NewFlagSet("suppr", flag.ExitOnError)
	supprId := supprCmd.Int("id", 0, "ID de la tâche")
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	findCmd := flag.NewFlagSet("find", flag.ExitOnError)
	findTitle := findCmd.String("title", "", "Titre de la tâche")

	if len(os.Args) < 2 {
		fmt.Println("Usage : todo-go <command> [--flags]")
		fmt.Println("Commands : add, list, done, suppr, find")
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
	case "suppr":
		supprCmd.Parse(os.Args[2:])
		if *supprId == 0 {
			fmt.Println("Usage: todo suppr --id N")
			os.Exit(1)
		}
		if err := suppr(*supprId); err != nil {
			fmt.Println("Erreur :", err)
			os.Exit(1)
		}
		if err := save(); err != nil {
			fmt.Println("Erreur de sauvegarde :", err)
			os.Exit(1)
		} else {
			fmt.Println("Tâche supprimée :", *supprId)
		}

	case "done":
		doneCmd.Parse(os.Args[2:])
		if *doneId == 0 {
			fmt.Println("Usage: todo done --id N")
			os.Exit(1)
		}
		if err := markDone(*doneId); err != nil {
			fmt.Println("Erreur :", err)
			os.Exit(1)
		}
		if err := save(); err != nil {
			fmt.Println("Erreur de sauvegarde :", err)
			os.Exit(1)
		} else {
			fmt.Println("Tâche marquée comme faite :", *doneId)
		}
	case "list":
		listCmd.Parse(os.Args[2:])
		list()
	case "find":
		findCmd.Parse(os.Args[2:])
		if *findTitle == "" {
			fmt.Println("Usage: todo find --title \"Ma tâche\"")
		}
		find(*findTitle)
	default:
		fmt.Println("Commande inconnue:", os.Args[1])
		os.Exit(1)
	}
}
