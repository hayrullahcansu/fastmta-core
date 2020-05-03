package dto

import "time"

type Broker struct {
	ID          string `json:"id"`
	Name        string
	StartedDate time.Time
	Status      string
}
