package model

import "time"

type EWallet struct {
	ID        int
	UserID    int
	Balance   int64
	Income    int64
	Expense   int64
	CreatedAt time.Time
	UpdatedAt *time.Time
}
