package main

import (
	"time"

	"github.com/google/uuid"
)

type Finding struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Title         string    `json:"title"`
	Severity      string    `json:"severity"`
	SeverityScore int       `json:"severity_score"`
	File          string    `json:"file"`
	AtLine        int       `json:"at_line"`
	Description   string    `json:"description"`
	PentestId     uuid.UUID `json:"pentest_id"`
}
