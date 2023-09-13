package entity

import "time"

// Delegation struct represent a delegation regarding Delegated POS
type Delegation struct {
	Amount    int64
	Block     string
	Id        int64
	Delegator string
	TimeStamp time.Time
}
