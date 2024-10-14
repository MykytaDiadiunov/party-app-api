package domain

import "time"

type Party struct {
	Id          uint64
	Title       string
	Description string
	Image       string
	Price       int32
	StartDate   time.Time
	CreatorId   uint64
}

type Parties struct {
	Parties     []Party
	Total       uint64
	CurrentPage int32
	LastPage    int32
}

func (p Party) GetUserId() uint64 {
	return p.CreatorId
}
