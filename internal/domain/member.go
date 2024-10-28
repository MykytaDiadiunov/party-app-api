package domain

type Member struct {
	PartyId uint64
	UserId  uint64
}

func (p Member) GetUserId() uint64 {
	return p.UserId
}
