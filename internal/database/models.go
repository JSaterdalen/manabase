// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Deck struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Name      string
	Commander sql.NullString
	OwnerID   uuid.UUID
}

type Game struct {
	ID         uuid.UUID
	CreatedAt  time.Time
	UpdatedAt  sql.NullTime
	DatePlayed time.Time
	IsTotem    bool
}

type Player struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	Name      string
}

type PlayerDeckGame struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	PlayerID  uuid.UUID
	GameID    uuid.UUID
	DeckID    uuid.UUID
	IsWon     bool
}
