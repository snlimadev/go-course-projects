package note

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type Note struct {
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

func New(title, content string) (Note, error) {
	if title == "" || content == "" {
		return Note{}, errors.New("invalid input")
	}

	return Note{
		Title:     title,
		Content:   content,
		CreatedAt: time.Now(),
	}, nil
}

func (note Note) Display() {
	fmt.Printf("\nYour note titled %v has the following content:\n%v\n", note.Title, note.Content)
}

func (note Note) Save() error {
	fileName := strings.ToLower(note.Title) + ".json"
	invalidChars := []string{" ", "<", ">", ":", "\"", "/", "\\", "|", "?", "*"}

	for _, char := range invalidChars {
		fileName = strings.ReplaceAll(fileName, char, "_")
	}

	json, err := json.Marshal(note)

	if err != nil {
		return err
	}

	return os.WriteFile(fileName, json, 0644)
}
