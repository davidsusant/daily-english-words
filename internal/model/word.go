package model

import "time"

// Word represents a vocabulary entry in the database
type Word struct {
	ID				int 		`json:"id"`
	Word			string		`json:"word"`
	PartOfSpeech	string		`json:"part_of_speech"`
	Definition		string		`json:"definition"`
	Example			string		`json:"example"`
	Difficulty		string		`json:"difficulty"`
	AssignedDate	*time.Time	`json:"assigned_date,omitempty"`
	CreatedAt		time.Time	`json:"created_at"`
}