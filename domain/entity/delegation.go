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

// DelegationRequest represent a query in order to show the delegations
type DelegationRequest struct {
	Limit  int
	Offset int
	Date   time.Time
}
