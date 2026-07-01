package domain

import "time"

type Market struct {
	ID        string
	Symbol    string
	Enabled   bool
	DeletedAt *time.Time
}

