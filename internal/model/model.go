package model

import (
	"time"
)

type Task struct {
	Title     string
	Text      string
	CreatedAt time.Time
	Complite  bool
}

type Event struct {
	RawInput  string
	ErrorText string
	CreatedAt time.Time
}

func NewTask(title string, text string) Task {
	return Task{
		Title:     title,
		Text:      text,
		CreatedAt: time.Now(),
		Complite:  false,
	}
}

func NewEvent(input string, errText string) Event {
	return Event{
		RawInput:  input,
		ErrorText: errText,
		CreatedAt: time.Now(),
	}
}
