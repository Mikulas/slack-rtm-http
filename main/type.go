package main

type Payload struct {
	Message string `json:"message"`
	RollbackValue string `json:"rollback,omitempty"`
}
