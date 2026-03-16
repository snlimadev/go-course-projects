package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"example.com/personal-notes/note"
	"example.com/personal-notes/todo"
)

type saver interface {
	Save() error
}

type outputtable interface {
	saver
	Display()
}

func main() {
	todoText := getUserInput("Todo:")
	todo, err := todo.New(todoText)

	if err != nil {
		fmt.Println("\nError:", err)
		return
	}

	title, content := getNoteData()
	userNote, err := note.New(title, content)

	if err != nil {
		fmt.Println("\nError:", err)
		return
	}

	err = outputData(todo)

	if err != nil {
		return
	}

	outputData(userNote)

	fmt.Println("\nPress ENTER to exit...")
	fmt.Scanln()
}

func getUserInput(prompt string) string {
	fmt.Printf("%v ", prompt)

	reader := bufio.NewReader(os.Stdin)
	text, err := reader.ReadString('\n')

	if err != nil {
		return ""
	}

	text = strings.TrimSuffix(text, "\n")
	text = strings.TrimSuffix(text, "\r")

	return text
}

func getNoteData() (string, string) {
	title := getUserInput("\nNote title:")
	content := getUserInput("Note content:")

	return title, content
}

func saveData(data saver) error {
	err := data.Save()

	if err != nil {
		fmt.Println("\nError saving:", err)
		return err
	}

	return nil
}

func outputData(data outputtable) error {
	data.Display()
	return saveData(data)
}
