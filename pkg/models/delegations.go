package models

import "time"

// Delegations represent the delegations structure can be found in db.
type Delegations struct {
	ID        uint      `db:"id"`
	TezosID   int       `json:"id" db:"id_tezos" gorm:"unique"`
	Timestamp time.Time `json:"timestamp"`
	Amount    int       `json:"amount"`
	Delegator string    `json:"delegator"`
	Level     int       `json:"level"`
}
