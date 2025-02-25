package main

type RequestStatus int

const (
	IN_PROGRESS RequestStatus = iota
	COMPLETE
)

func (s RequestStatus) String() string {
	return [...]string{"in_progress", "complete"}[s]
}

type Request struct {
	IdempotencyKey string // uuid
	Status         RequestStatus
	Response       string // json string
}
